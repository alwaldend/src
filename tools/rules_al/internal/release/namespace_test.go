package release

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// TestMain supplies isolated SSH and rsync endpoints to the built release CLI.
func TestMain(m *testing.M) {
	if os.Getenv("RELEASE_TRANSPORT_HELPER") == "1" {
		if err := transportFixture(os.Args[1:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	os.Exit(m.Run())
}

// TestPublicationAgentCLI authenticates with a public identity and an agent-held key.
func TestPublicationAgentCLI(t *testing.T) {
	sshPath, err := exec.LookPath("ssh")
	if err != nil {
		t.Fatalf("find OpenSSH used by publication: %v", err)
	}
	dir := t.TempDir()
	t.Chdir(dir)
	// A relative Unix socket avoids platform path-length limits in Bazel scratch.
	t.Setenv("SSH_AUTH_SOCK", "agent.sock")
	address, publicKey := publicationSSHFixture(t)
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		t.Fatal(err)
	}
	identity, knownHosts := filepath.Join(dir, "identity.pub"), filepath.Join(dir, "known_hosts")
	writeBundleFixture(t, identity, string(ssh.MarshalAuthorizedKey(publicKey)))
	writeBundleFixture(t, knownHosts, "["+host+"]:"+port+" "+string(ssh.MarshalAuthorizedKey(publicKey)))
	status, payload, bundle := filepath.Join(dir, "version"), filepath.Join(dir, "payload"), filepath.Join(dir, "bundle")
	writeBundleFixture(t, status, "STABLE_VERSION 2026.44.0\nSTABLE_VERSION_CHANNEL release\nSTABLE_GIT_COMMIT 0123456789012345678901234567890123456789\nSTABLE_GIT_TREE_STATE clean\n")
	writeBundleFixture(t, payload, "fixture\n")
	if output, err := releaseBundleCommand("generate", "--project", "projects/fixture", "--version_file", status, "--add_file", payload, "--output_dir", bundle).CombinedOutput(); err != nil {
		t.Fatalf("generate agent fixture: %v\n%s", err, output)
	}
	args := []string{"deploy", "--release_dir", bundle, "--environment", "fixture", "--ssh_host", host, "--ssh_port", port, "--ssh_user", "fixture", "--ssh_config", "/dev/null", "--ssh_identity", identity, "--ssh_known_hosts", knownHosts, "--ssh_path", sshPath, "--steps", "link", "--timeout", "10s"}
	if output, err := releaseBundleCommand(args...).CombinedOutput(); err != nil {
		t.Fatalf("authenticate with agent-held identity: %v\n%s", err, output)
	}
	if output, err := releaseBundleCommand(append(args, "--ssh_disable_agent")...).CombinedOutput(); err == nil {
		t.Fatalf("agent-only identity succeeded with agent disabled: %s", output)
	} else if !bytes.Contains(output, []byte("Permission denied")) {
		t.Fatalf("expected authentication rejection with agent disabled: %v\n%s", err, output)
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		artifact, err := json.MarshalIndent(map[string]any{"agent_authentication": true, "disabled_agent_rejected": true, "identity": ssh.FingerprintSHA256(publicKey)}, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "ssh-agent.json"), artifact, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// publicationSSHFixture serves an isolated agent and SSH command endpoint.
// Its generated private key stays in memory and its listeners close with the test.
func publicationSSHFixture(t *testing.T) (string, ssh.PublicKey) {
	t.Helper()
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		t.Fatal(err)
	}
	keyring := agent.NewKeyring()
	if err := keyring.Add(agent.AddedKey{PrivateKey: privateKey}); err != nil {
		t.Fatal(err)
	}
	agentListener, err := net.Listen("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = agentListener.Close() })
	go func() {
		for {
			conn, err := agentListener.Accept()
			if err != nil {
				return
			}
			go func() {
				defer conn.Close()
				_ = agent.ServeAgent(keyring, conn)
			}()
		}
	}()
	config := &ssh.ServerConfig{PublicKeyCallback: func(_ ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
		if !bytes.Equal(key.Marshal(), signer.PublicKey().Marshal()) {
			return nil, fmt.Errorf("unrecognized fixture identity")
		}
		return nil, nil
	}}
	config.AddHostKey(signer)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = listener.Close() })
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go func() {
				defer conn.Close()
				if err := conn.SetDeadline(time.Now().Add(15 * time.Second)); err != nil {
					return
				}
				server, channels, requests, err := ssh.NewServerConn(conn, config)
				if err != nil {
					return // Authentication failure is expected when the agent is disabled.
				}
				defer server.Close()
				go ssh.DiscardRequests(requests)
				for request := range channels {
					channel, requests, err := request.Accept()
					if err != nil {
						return
					}
					for request := range requests {
						if request.Type == "exec" {
							_ = request.Reply(true, nil)
							_, _ = channel.SendRequest("exit-status", false, ssh.Marshal(struct {
								// Status reports successful execution to OpenSSH.
								Status uint32
							}{0}))
							_ = channel.Close()
							break
						}
						_ = request.Reply(false, nil)
					}
				}
			}()
		}
	}()
	return listener.Addr().String(), signer.PublicKey()
}

// transportFixture executes the SSH shell protocol or copies selected payloads.
func transportFixture(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("fixture transport has no arguments")
	}
	if args[0] != "--recursive" {
		cmd := exec.Command("/bin/sh", "-c", args[len(args)-1])
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("fixture remote shell: %w", err)
		}
		return nil
	}
	source, destination := args[len(args)-2], args[len(args)-1]
	if _, remote, ok := strings.Cut(source, ":"); ok {
		source = remote
	}
	if _, remote, ok := strings.Cut(destination, ":"); ok {
		destination = remote
	}
	info, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("inspect fixture source: %w", err)
	}
	if !info.IsDir() {
		body, err := os.ReadFile(source)
		if err != nil {
			return fmt.Errorf("read fixture download: %w", err)
		}
		if err := os.WriteFile(destination, body, 0o644); err != nil {
			return fmt.Errorf("write fixture download: %w", err)
		}
		return nil
	}
	if err := os.MkdirAll(destination, 0o755); err != nil {
		return fmt.Errorf("create fixture destination: %w", err)
	}
	if err := os.CopyFS(destination, os.DirFS(source)); err != nil {
		return fmt.Errorf("transfer fixture directory: %w", err)
	}

	return nil
}

// TestPublicationNamespaceCLI checks both namespaces and rejects path escapes.
func TestPublicationNamespaceCLI(t *testing.T) {
	for _, namespace := range []string{"projects", "third_party"} {
		t.Run(namespace, func(t *testing.T) {
			dir := t.TempDir()
			status, payload := filepath.Join(dir, "version"), filepath.Join(dir, "openapi.json")
			version := "2026.44.0"
			if namespace == "third_party" {
				version = "2.173-4377d1791c01"
			}
			writeBundleFixture(t, status, "STABLE_VERSION "+version+"\nSTABLE_VERSION_CHANNEL release\nSTABLE_GIT_COMMIT 0123456789012345678901234567890123456789\nSTABLE_GIT_TREE_STATE clean\n")
			writeBundleFixture(t, payload, "schema fixture\n")
			deployment := filepath.Join(dir, "deployment.json")
			configuredNamespace := namespace
			if namespace == "projects" {
				configuredNamespace = ""
			}
			writeBundleFixture(t, deployment, fmt.Sprintf(`{"ssh":{"environment":"fixture","host":"fixture.invalid","namespace":%q,"public_url":"https://download.example"}}`, configuredNamespace))
			bundle := filepath.Join(dir, "bundle")
			cmd := releaseBundleCommand("generate", "--project", namespace+"/com_x_api", "--version_file", status, "--add_file", payload, "--deployment", deployment, "--output_dir", bundle)
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("generate namespace bundle: %v\n%s", err, output)
			}
			manifest, err := os.ReadFile(filepath.Join(bundle, "release.json"))
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Contains(manifest, []byte("https://download.example/"+namespace+"/com_x_api/releases/"+version+"/openapi.json")) {
				t.Fatalf("incorrect namespace URL in %s", manifest)
			}
			helper, err := os.Executable()
			if err != nil {
				t.Fatal(err)
			}
			remote := filepath.Join(dir, "remote")
			if err := os.MkdirAll(filepath.Join(remote, "staging"), 0o755); err != nil {
				t.Fatal(err)
			}
			args := []string{"deploy", "--release_dir", bundle, "--ssh_deployment", deployment, "--environment", "fixture", "--ssh_root", remote, "--ssh_path", helper, "--rsync_path", helper}
			cmd = releaseBundleCommand(args...)
			cmd.Env = append(os.Environ(), "RELEASE_TRANSPORT_HELPER=1")
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("publish namespace bundle: %v\n%s", err, output)
			}
			body, err := os.ReadFile(filepath.Join(remote, namespace, "com_x_api", "releases", version, "openapi.json"))
			if err != nil || string(body) != "schema fixture\n" {
				t.Fatalf("published bytes: %q, error=%v", body, err)
			}
			writeBundleFixture(t, payload, "replacement fixture\n")
			rival := filepath.Join(dir, "rival")
			if output, err := releaseBundleCommand("generate", "--project", namespace+"/com_x_api", "--version_file", status, "--add_file", payload, "--output_dir", rival).CombinedOutput(); err != nil {
				t.Fatalf("generate conflicting namespace bundle: %v\n%s", err, output)
			}
			cmd = releaseBundleCommand(append(args, "--release_dir", rival)...)
			cmd.Env = append(os.Environ(), "RELEASE_TRANSPORT_HELPER=1")
			if output, err := cmd.CombinedOutput(); err == nil {
				t.Fatalf("namespace collision was accepted: %s", output)
			}
			body, err = os.ReadFile(filepath.Join(remote, namespace, "com_x_api", "releases", version, "openapi.json"))
			if err != nil || string(body) != "schema fixture\n" {
				t.Fatalf("namespace collision changed bytes: %q, error=%v", body, err)
			}
			for _, invalid := range []string{"../outside", "sites", "staging", "third_party/other"} {
				cmd = releaseBundleCommand(append(args, "--namespace", invalid)...)
				cmd.Env = append(os.Environ(), "RELEASE_TRANSPORT_HELPER=1")
				if output, err := cmd.CombinedOutput(); err == nil {
					t.Fatalf("unsafe namespace %q accepted: %s", invalid, output)
				}
			}
			if namespace == "third_party" {
				cmd = releaseBundleCommand(append(args, "--steps", "link")...)
				cmd.Env = append(os.Environ(), "RELEASE_TRANSPORT_HELPER=1")
				if output, err := cmd.CombinedOutput(); err == nil {
					t.Fatalf("third-party website selection accepted: %s", output)
				}
			}
			var artifact any
			if err := json.Unmarshal(manifest, &artifact); err != nil {
				t.Fatal(err)
			}
			if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
				if err := os.WriteFile(filepath.Join(output, namespace+".json"), manifest, 0o644); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// TestMain supplies isolated native-tool fixtures without shell scripts.
func TestMain(m *testing.M) {
	switch filepath.Base(os.Args[0]) {
	case "p11tool":
		os.Exit(serveDiscoveryFixture())
	case "curl":
		os.Exit(serveCurlFixture())
	default:
		os.Exit(m.Run())
	}
}

// serveDiscoveryFixture exposes public certificates on hardware, not trust stores.
func serveDiscoveryFixture() int {
	args := strings.Join(os.Args[1:], " ")
	switch {
	case strings.Contains(args, "--list-tokens"):
		fmt.Fprint(os.Stdout, "Token 0:\n URL: pkcs11:token=system\n Type: Trust module\n"+
			"Token 1:\n URL: pkcs11:token=fixture\n Type: Hardware token\n"+
			"Token 2:\n URL: pkcs11:token=last-trust\n Type: Trust module\n")
	case strings.Contains(args, "--list-all-certs"):
		if os.Args[len(os.Args)-1] != "pkcs11:token=fixture" {
			return 96
		}
		fmt.Fprintln(os.Stdout, os.Getenv("FIXTURE_URIS"))
	case strings.Contains(args, "--export"):
		var certificates map[string]string
		if err := json.Unmarshal([]byte(os.Getenv("FIXTURE_CERTIFICATES")), &certificates); err != nil {
			return 97
		}
		certificate := certificates[os.Args[len(os.Args)-1]]
		if certificate == "" {
			return 98
		}
		fmt.Fprint(os.Stdout, certificate)
	default:
		return 99
	}
	return 0
}

// fixtureCertificate creates public synthetic certificate data for CLI discovery.
func fixtureCertificate(t *testing.T, commonName string, usage x509.ExtKeyUsage, expired bool) string {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	end := now.Add(time.Hour)
	if expired {
		end = now.Add(-time.Minute)
	}
	certificate := &x509.Certificate{
		SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: commonName},
		NotBefore: now.Add(-time.Hour), NotAfter: end,
		KeyUsage: x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{usage},
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, certificate, certificate, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))
}

// serveCurlFixture checks the hardware transport boundary and returns test data.
func serveCurlFixture() int {
	if err := os.WriteFile(os.Getenv("FIXTURE_CALLED"), []byte("called"), 0o600); err != nil {
		return 90
	}
	for _, option := range []string{"--cert-type", "--key-type"} {
		if argumentValue(option) != os.Getenv("FIXTURE_TYPE") {
			return 91
		}
	}
	if argumentValue("--engine") != "pkcs11" || os.Args[1] != "--disable" {
		return 92
	}
	if argumentValue("--cert") != os.Getenv("FIXTURE_SELECTED_URI") {
		return 93
	}
	if argumentValue("--key") != os.Getenv("FIXTURE_SELECTED_KEY") {
		return 95
	}
	var request map[string]string
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil || request["name"] != `fixture"role` {
		return 94
	}
	if os.Getenv("FIXTURE_MODE") == "timeout" {
		time.Sleep(5 * time.Second)
	}
	if os.Getenv("FIXTURE_MODE") == "failure" {
		fmt.Fprintln(os.Stderr, "curl: (58) fixture TLS certificate loading error")
		return 58
	}
	if os.Getenv("FIXTURE_MODE") == "oversized" {
		_, _ = io.WriteString(os.Stdout, strings.Repeat("x", 2<<20))
		return 0
	}
	_, _ = io.WriteString(os.Stdout, os.Getenv("FIXTURE_RESPONSE"))
	return 0
}

// argumentValue retrieves a native fixture option without accepting shell syntax.
func argumentValue(name string) string {
	for i, arg := range os.Args[1:] {
		if arg == name && i+2 < len(os.Args) {
			return os.Args[i+2]
		}
	}
	return ""
}

// TestLoginCommand exercises the built command and real private file replacement.
func TestLoginCommand(t *testing.T) {
	binary, err := runfiles.Rlocation(os.Getenv("LOGIN_BINARY"))
	if err != nil {
		t.Fatal(err)
	}
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	// certificateURI selects synthetic hardware without accessing a real token.
	const certificateURI = "pkcs11:token=fixture;object=Authentication%20certificate;id=%01;type=cert"
	certificates, err := json.Marshal(map[string]string{
		certificateURI: fixtureCertificate(t, `fixture"role.users.example`, x509.ExtKeyUsageClientAuth, false),
		"pkcs11:token=fixture;id=other;type=cert":   fixtureCertificate(t, "other.users.example", x509.ExtKeyUsageClientAuth, false),
		"pkcs11:token=fixture;id=expired;type=cert": fixtureCertificate(t, `fixture"role.users.example`, x509.ExtKeyUsageClientAuth, true),
		"pkcs11:token=fixture;id=server;type=cert":  fixtureCertificate(t, `fixture"role.users.example`, x509.ExtKeyUsageServerAuth, false),
	})
	if err != nil {
		t.Fatal(err)
	}
	results := make(map[string]string)
	for _, scenario := range []struct {
		// name identifies the command-level behavior.
		name string
		// response is a synthetic Vault response body.
		response string
		// mode controls native transport failure injection.
		mode string
		// uris overrides discovery output when nonempty.
		uris string
		// selectedURI identifies the expected automatic or explicit selection.
		selectedURI string
		// args adds typed command options.
		args []string
		// address overrides the TLS server address when nonempty.
		address string
		// success allows the token file to change.
		success bool
		// called requires a transport invocation after input validation.
		called bool
		// symlink requests a symbolic-link destination.
		symlink bool
	}{
		{name: "provider-login", response: `{"auth":{"client_token":"fixture-new-token"}}`, success: true, called: true},
		{name: "engine-login", response: `{"auth":{"client_token":"fixture-new-token"}}`, args: []string{"--certificate-type", "ENG"}, success: true, called: true},
		{name: "malformed-response", response: "fixture-private-response", called: true},
		{name: "missing-token", response: `{"auth":{}}`, called: true},
		{name: "whitespace-token", response: `{"auth":{"client_token":"invalid token"}}`, called: true},
		{name: "trailing-json", response: `{"auth":{"client_token":"fixture-new-token"}} {}`, called: true},
		{name: "transport-failure", mode: "failure", called: true},
		{name: "transport-timeout", mode: "timeout", args: []string{"--timeout", "100ms"}, called: true},
		{name: "oversized-response", mode: "oversized", called: true},
		{name: "role-matching-discovery", uris: certificateURI + "\npkcs11:token=fixture;id=other;type=cert", response: `{"auth":{"client_token":"fixture-new-token"}}`, success: true, called: true},
		{name: "expired-certificate", uris: "pkcs11:token=fixture;id=expired;type=cert"},
		{name: "server-only-certificate", uris: "pkcs11:token=fixture;id=server;type=cert"},
		{name: "skip-expired-certificate", uris: "pkcs11:token=fixture;id=expired;type=cert\n" + certificateURI, response: `{"auth":{"client_token":"fixture-new-token"}}`, success: true, called: true},
		{name: "empty-discovery", uris: "No tokens found"},
		{name: "explicit-selection", uris: "No tokens found", args: []string{"--certificate-uri", certificateURI}, response: `{"auth":{"client_token":"fixture-new-token"}}`, success: true, called: true},
		{name: "plaintext-address", address: "http://vault.invalid"},
		{name: "credential-address", address: "https://user:fixture-private-response@vault.invalid"},
		{name: "pin-bearing-uri", args: []string{"--certificate-uri", "pkcs11:token=fixture?pin-value=fixture-private-response"}},
		{name: "invalid-mode", args: []string{"--certificate-type", "PEM"}},
		{name: "symlink-destination", response: `{"auth":{"client_token":"fixture-new-token"}}`, symlink: true},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			dir := t.TempDir()
			for _, name := range []string{"curl", "p11tool"} {
				if err := os.Symlink(self, filepath.Join(dir, name)); err != nil {
					t.Fatal(err)
				}
			}
			tokenPath := filepath.Join(dir, "token")
			if err := os.WriteFile(tokenPath, []byte("previous-token\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			if scenario.symlink {
				link := filepath.Join(dir, "token-link")
				if err := os.Symlink(tokenPath, link); err != nil {
					t.Fatal(err)
				}
				tokenPath = link
			}
			ca := filepath.Join(dir, "ca.pem")
			if err := os.WriteFile(ca, []byte("fixture CA"), 0o600); err != nil {
				t.Fatal(err)
			}
			uris := scenario.uris
			if uris == "" {
				uris = "URL: " + certificateURI
			}
			selectedURI, selectedKey := certificateURI, "pkcs11:token=fixture;id=%01;type=private"
			if scenario.selectedURI != "" {
				selectedURI, selectedKey = scenario.selectedURI, scenario.selectedURI
			}
			address := scenario.address
			if address == "" {
				address = "https://vault.invalid"
			}
			certificateType := "PROV"
			if scenario.name == "engine-login" {
				certificateType = "ENG"
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			args := append([]string{`fixture"role`, "--token-file", tokenPath}, scenario.args...)
			command := exec.CommandContext(ctx, binary, args...)
			calledPath := filepath.Join(dir, "called")
			command.Env = append(os.Environ(), "PATH="+dir, "VAULT_ADDR="+address, "VAULT_CACERT="+ca,
				"FIXTURE_URIS="+uris, "FIXTURE_RESPONSE="+scenario.response, "FIXTURE_MODE="+scenario.mode,
				"FIXTURE_CERTIFICATES="+string(certificates),
				"FIXTURE_TYPE="+certificateType, "FIXTURE_SELECTED_URI="+selectedURI,
				"FIXTURE_SELECTED_KEY="+selectedKey, "FIXTURE_CALLED="+calledPath)
			output, err := command.CombinedOutput()
			if (err == nil) != scenario.success {
				t.Fatalf("success=%v, error=%v, output=%s", scenario.success, err, output)
			}
			if scenario.mode == "failure" && !strings.Contains(string(output), "fixture TLS certificate loading error") {
				t.Fatal("curl's useful TLS diagnostic was hidden")
			}
			for _, secret := range []string{"fixture-private-response", "fixture-new-token", "previous-token"} {
				if strings.Contains(string(output), secret) {
					t.Fatal("secret fixture data leaked into command output")
				}
			}
			_, statErr := os.Stat(calledPath)
			if (statErr == nil) != scenario.called {
				t.Fatalf("transport invocation=%v, want %v", statErr == nil, scenario.called)
			}
			got, err := os.ReadFile(tokenPath)
			if err != nil {
				t.Fatal(err)
			}
			want := "previous-token\n"
			if scenario.success {
				want = "fixture-new-token\n"
				info, err := os.Stat(tokenPath)
				if err != nil || info.Mode().Perm() != 0o600 {
					t.Fatal("successful replacement must have mode 0600")
				}
			}
			if string(got) != want {
				t.Fatal("token file differs from the expected successful replacement or preserved original")
			}
			files, err := filepath.Glob(filepath.Join(dir, ".vault-token-*"))
			if err != nil || len(files) != 0 {
				t.Fatal("temporary token file was not cleaned up")
			}
			results[scenario.name] = "passed"
		})
	}
	if outputDir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); outputDir != "" {
		data, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(outputDir, "login-results.json"), append(data, '\n'), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

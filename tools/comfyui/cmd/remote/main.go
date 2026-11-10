package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
	"github.com/spf13/cobra"
)

// main reports command failures without rendering injected credentials.
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err := command().ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			os.Exit(exit.ExitCode())
		}
		os.Exit(1)
	}
}

// command adapts wrapper flags and passes arguments after -- to the CLI.
func command() *cobra.Command {
	var binary, endpoint, endpointCommand, caFile string
	var endpointArgs []string
	cmd := &cobra.Command{
		Use: "remote [bridge options] -- [comfy options and command]", SilenceUsage: true, SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateRouting(args); err != nil {
				return fmt.Errorf("validate CLI routing: %w", err)
			}
			if endpoint == "" {
				endpoint = os.Getenv("COMFY_HTTPS_URL")
			}
			if endpoint == "" && endpointCommand != "" {
				path, err := executable(endpointCommand)
				if err != nil {
					return fmt.Errorf("resolve endpoint command: %w", err)
				}
				child := exec.CommandContext(cmd.Context(), path, endpointArgs...)
				child.Env = childEnvironment("")
				output, err := child.Output()
				if err != nil {
					return fmt.Errorf("read rental endpoint: check Vault access and rent/refresh connection outputs: %w", err)
				}
				// The owning Terraform runner initializes its backend before output.
				// With -no-color, the requested raw value is the final output line.
				lines := strings.Split(strings.TrimSpace(string(output)), "\n")
				endpoint = strings.TrimSpace(lines[len(lines)-1])
			}
			if err := run(cmd.Context(), binary, endpoint, caFile, args); err != nil {
				return fmt.Errorf("execute authenticated CLI: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&binary, "cli", "", "Packaged CLI runfile")
	cmd.Flags().StringVar(&endpoint, "url", "", "HTTPS endpoint; defaults to COMFY_HTTPS_URL or the endpoint command")
	cmd.Flags().StringVar(&endpointCommand, "endpoint-command", "", "Read-only endpoint command runfile")
	cmd.Flags().StringArrayVar(&endpointArgs, "endpoint-arg", nil, "Endpoint command argument")
	cmd.Flags().StringVar(&caFile, "ca-file", "", "Trusted PEM or DER certificate runfile")
	return cmd
}

// executable resolves an absolute filesystem path or declared Bazel runfile.
func executable(name string) (string, error) {
	if name == "" {
		return "", errors.New("missing declared executable or certificate")
	}
	if filepath.IsAbs(name) {
		return name, nil
	}
	path, err := runfiles.Rlocation(name)
	if err != nil {
		return "", fmt.Errorf("resolve declared runfile: %w", err)
	}
	return path, nil
}

// validateRouting rejects arguments that could bypass the authenticated target.
func validateRouting(args []string) error {
	for i, arg := range args {
		name, value, inline := strings.Cut(arg, "=")
		if name == "--host" || name == "--port" {
			return errors.New("host and port are owned by the authenticated target")
		}
		if name == "--where" {
			if !inline && i+1 < len(args) {
				value = args[i+1]
			}
			if value != "local" {
				return errors.New("the authenticated target requires --where local")
			}
		}
	}
	return nil
}

// childEnvironment selects the local target and removes injected credentials.
// An empty address retains backend credentials for the owning output command.
func childEnvironment(address string) []string {
	env := make([]string, 0, len(os.Environ())+3)
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		upper := strings.ToUpper(key)
		if upper == "HTTP_PROXY" || upper == "HTTPS_PROXY" || upper == "ALL_PROXY" || upper == "NO_PROXY" {
			continue
		}
		if key == "HF_TOKEN" || key == "HUGGING_FACE_HUB_TOKEN" || strings.HasPrefix(key, "TF_VAR_") || strings.HasPrefix(key, "IMAGEGEN_") || strings.HasPrefix(key, "VAULT_") || strings.HasPrefix(key, "VAST_") || (strings.HasPrefix(key, "COMFY_") && key != "COMFY_REPOSITORY_PROFILE") {
			continue
		}
		if address != "" && strings.HasPrefix(key, "TF_HTTP_") {
			continue
		}
		env = append(env, entry)
	}
	return append(env, "COMFY_WHERE=local", "COMFY_LOCAL_URL=http://"+address, "COMFY_NO_WATCH=1", "DO_NOT_TRACK=1", "NO_PROXY=127.0.0.1,localhost")
}

// run owns verified transport, the local listener and its foreground CLI child.
func run(ctx context.Context, binary, endpoint, caFile string, args []string) (resultErr error) {
	target, err := url.Parse(endpoint)
	if err != nil || target.Scheme != "https" || target.Host == "" || target.User != nil || target.RawQuery != "" || target.Fragment != "" || (target.Path != "" && target.Path != "/") {
		return errors.New("no valid root HTTPS rental endpoint; rent and refresh outputs first")
	}
	username, password := os.Getenv("IMAGEGEN_WEB_USERNAME"), os.Getenv("IMAGEGEN_WEB_PASSWORD")
	if username == "" || password == "" || strings.Contains(username, ":") {
		return errors.New("missing or invalid Vault-injected web login")
	}
	path, err := executable(caFile)
	if err != nil {
		return fmt.Errorf("resolve trusted certificate: %w", err)
	}
	certificate, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read trusted certificate: %w", err)
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(certificate) {
		cert, err := x509.ParseCertificate(certificate)
		if err != nil {
			return fmt.Errorf("parse trusted certificate: %w", err)
		}
		roots.AppendCertsFromPEM(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}))
	}
	transport := &http.Transport{TLSClientConfig: &tls.Config{RootCAs: roots, MinVersion: tls.VersionTLS12}, TLSHandshakeTimeout: 15 * time.Second, ResponseHeaderTimeout: 60 * time.Second}
	defer transport.CloseIdleConnections()
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("bind private CLI bridge: %w", err)
	}
	proxy := &httputil.ReverseProxy{
		Rewrite: func(request *httputil.ProxyRequest) {
			request.SetURL(target)
			request.Out.Host = target.Host
			request.Out.Header.Del("Cookie")
			if request.Out.Header.Get("Origin") != "" {
				request.Out.Header.Set("Origin", "https://"+target.Host)
			}
			request.Out.SetBasicAuth(username, password)
		},
		Transport: transport,
		ErrorLog:  log.New(io.Discard, "", 0),
		ErrorHandler: func(w http.ResponseWriter, _ *http.Request, _ error) {
			http.Error(w, "verified HTTPS API request failed", http.StatusBadGateway)
		},
		ModifyResponse: func(response *http.Response) error {
			if response.StatusCode >= 300 && response.StatusCode < 400 {
				return errors.New("rental API redirects are refused")
			}
			return nil
		},
	}
	address := listener.Addr().String()
	server := &http.Server{ReadHeaderTimeout: 10 * time.Second, Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if r.Host != address || (origin != "" && !(r.Header.Get("Upgrade") == "websocket" && origin == "http://"+address)) || !sameUser(r.RemoteAddr, address) {
			http.Error(w, "private CLI connection required", http.StatusForbidden)
			return
		}
		proxy.ServeHTTP(w, r)
	})}
	defer server.Close()
	go func() { _ = server.Serve(listener) }()
	cli, err := executable(binary)
	if err != nil {
		return fmt.Errorf("resolve packaged CLI: %w", err)
	}
	child := exec.CommandContext(ctx, cli, args...)
	workspace := os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	if workspace == "" {
		workspace, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("locate CLI task storage: %w", err)
		}
	}
	profileRoot := os.Getenv("COMFY_REPOSITORY_PROFILE")
	if profileRoot == "" {
		profileRoot = filepath.Join(workspace, "out", "comfyui", "cli")
	}
	if err := os.MkdirAll(profileRoot, 0o700); err != nil {
		return fmt.Errorf("prepare CLI task profile: %w", err)
	}
	profile, err := os.MkdirTemp(profileRoot, "invocation-")
	if err != nil {
		return fmt.Errorf("create private CLI profile: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(profile); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("remove private CLI profile: %w", err))
		}
	}()
	child.Env = append(childEnvironment(address), "COMFY_REPOSITORY_PROFILE="+profile)
	child.Stdin, child.Stdout, child.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := child.Run(); err != nil {
		return fmt.Errorf("run packaged ComfyUI CLI: %w", err)
	}
	return nil
}

// sameUser checks the established client's Linux TCP socket ownership, failing closed.
func sameUser(peer, local string) bool {
	peerHost, peerPort, err := net.SplitHostPort(peer)
	if err != nil || peerHost != "127.0.0.1" {
		return false
	}
	_, localPort, err := net.SplitHostPort(local)
	if err != nil {
		return false
	}
	remoteNumber, err := strconv.Atoi(peerPort)
	if err != nil {
		return false
	}
	localNumber, err := strconv.Atoi(localPort)
	if err != nil {
		return false
	}
	file, err := os.Open("/proc/net/tcp")
	if err != nil {
		return false
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 8 && fields[1] == fmt.Sprintf("0100007F:%04X", remoteNumber) && fields[2] == fmt.Sprintf("0100007F:%04X", localNumber) && fields[3] == "01" {
			return fields[7] == strconv.Itoa(os.Getuid())
		}
	}
	return false
}

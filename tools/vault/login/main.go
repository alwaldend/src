package main

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode"

	"github.com/spf13/cobra"
)

// options contains the command's typed hardware, timeout, and storage settings.
type options struct {
	// certificateURI selects a certificate without automatic token discovery.
	certificateURI string
	// certificateType selects OpenSSL provider or legacy engine loading.
	certificateType string
	// tokenFile overrides the default token destination beneath the user's home.
	tokenFile string
	// timeout bounds discovery and the login request together.
	timeout time.Duration
}

// privateError preserves a cause without exposing secret-bearing diagnostic text.
type privateError struct {
	// cause is available through errors.Unwrap without appearing in Error.
	cause error
}

// Error omits potentially sensitive parse or subprocess output.
func (e privateError) Error() string { return "details withheld" }

// Unwrap preserves the original failure for programmatic inspection.
func (e privateError) Unwrap() error { return e.cause }

// boundedOutput limits retained native-tool output, including private responses.
type boundedOutput struct {
	// buffer holds bytes within the configured limit.
	buffer bytes.Buffer
	// limit is the maximum number of retained bytes.
	limit int
}

// Write rejects excess output instead of allocating an unbounded response buffer.
func (b *boundedOutput) Write(p []byte) (int, error) {
	if len(p) > b.limit-b.buffer.Len() {
		return 0, errors.New("native tool output exceeds the size limit")
	}
	n, err := b.buffer.Write(p)
	if err != nil {
		return n, fmt.Errorf("buffer native tool output: %w", err)
	}
	return n, nil
}

// main owns cancellation and reports a single sanitized command failure.
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := newCommand().ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "vault login:", err)
		os.Exit(1)
	}
}

// newCommand declares the supported interface without arbitrary curl arguments.
func newCommand() *cobra.Command {
	settings := options{}
	command := &cobra.Command{
		Use:           "login <role-name>",
		Short:         "Log in to Vault using a PKCS#11 client certificate",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(command *cobra.Command, args []string) error {
			if err := settings.login(command.Context(), args[0], command.ErrOrStderr()); err != nil {
				return fmt.Errorf("certificate login: %w", err)
			}
			if _, err := fmt.Fprintln(command.OutOrStdout(), "Vault login succeeded."); err != nil {
				return fmt.Errorf("report login success: %w", err)
			}
			return nil
		},
	}
	command.Flags().StringVar(&settings.certificateURI, "certificate-uri", "", "Explicit PKCS#11 certificate URI (without a PIN)")
	command.Flags().StringVar(&settings.certificateType, "certificate-type", "PROV", "OpenSSL certificate and key type: PROV or ENG")
	command.Flags().StringVar(&settings.tokenFile, "token-file", "", "Token destination (default: ~/.vault-token)")
	command.Flags().DurationVar(&settings.timeout, "timeout", time.Minute, "Maximum discovery and login duration")
	return command
}

// login validates configuration, authenticates once, and then saves the token.
func (o options) login(parent context.Context, role string, diagnostics io.Writer) error {
	if strings.TrimSpace(role) == "" || o.timeout <= 0 {
		return errors.New("a role name and a positive timeout are required")
	}
	if o.certificateType != "PROV" && o.certificateType != "ENG" {
		return errors.New("certificate type must be PROV or ENG")
	}
	address, err := url.Parse(os.Getenv("VAULT_ADDR"))
	if err != nil {
		return fmt.Errorf("parse VAULT_ADDR: %w", privateError{cause: err})
	}
	if address.Scheme != "https" || address.Hostname() == "" || address.User != nil || address.RawQuery != "" || address.Fragment != "" {
		return errors.New("VAULT_ADDR must be an HTTPS URL without credentials, query, or fragment")
	}
	ca := os.Getenv("VAULT_CACERT")
	if ca == "" {
		return errors.New("VAULT_CACERT is required")
	}
	if info, err := os.Stat(ca); err != nil {
		return fmt.Errorf("inspect Vault CA file: %w", err)
	} else if !info.Mode().IsRegular() {
		return errors.New("VAULT_CACERT must identify a regular file")
	}
	if o.tokenFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("find home directory for token storage: %w", err)
		}
		o.tokenFile = filepath.Join(home, ".vault-token")
	}
	if err := validateDestination(o.tokenFile); err != nil {
		return fmt.Errorf("validate token destination: %w", err)
	}
	ctx, cancel := context.WithTimeout(parent, o.timeout)
	defer cancel()
	uri := o.certificateURI
	if uri == "" {
		uri, err = discoverCertificate(ctx, role)
		if err != nil {
			return fmt.Errorf("select certificate: %w", err)
		}
	}
	decodedURI, err := url.PathUnescape(uri)
	if err != nil {
		return fmt.Errorf("parse certificate URI: %w", privateError{cause: err})
	}
	if !strings.HasPrefix(uri, "pkcs11:") || len(uri) == len("pkcs11:") ||
		strings.Contains(strings.ToLower(decodedURI), "pin-value") || strings.Contains(strings.ToLower(decodedURI), "pin-source") {
		return errors.New("certificate URI must use pkcs11: without embedded PIN options")
	}
	payload, err := json.Marshal(map[string]string{"name": role})
	if err != nil {
		return fmt.Errorf("encode certificate login request: %w", err)
	}
	args := []string{
		"--disable", "--silent", "--show-error", "--fail", "--proto", "=https",
		"--engine", "pkcs11", "--cert-type", o.certificateType, "--key-type", o.certificateType,
		"--cert", uri, "--key", privateKeyURI(uri), "--cacert", ca,
		"--max-time", strconv.FormatFloat(o.timeout.Seconds(), 'f', 3, 64),
		"--header", "Content-Type: application/json", "--request", "POST", "--data-binary", "@-",
		strings.TrimRight(address.String(), "/") + "/v1/auth/cert/login",
	}
	response, err := runNative(ctx, "curl", args, payload, 1<<20, diagnostics)
	if err != nil {
		return fmt.Errorf("request certificate login: %w", err)
	}
	var envelope struct {
		// Auth contains Vault's authentication result.
		Auth struct {
			// ClientToken is retained privately until it is persisted.
			ClientToken string `json:"client_token"`
		} `json:"auth"`
	}
	decoder := json.NewDecoder(bytes.NewReader(response))
	if err := decoder.Decode(&envelope); err != nil {
		return fmt.Errorf("decode login response: %w", privateError{cause: err})
	}
	if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
		return errors.New("login response must contain exactly one JSON document")
	}
	token := envelope.Auth.ClientToken
	if token == "" || strings.IndexFunc(token, func(r rune) bool { return unicode.IsSpace(r) || unicode.IsControl(r) }) >= 0 {
		return errors.New("login response did not contain a valid client token")
	}
	if err := saveToken(o.tokenFile, token); err != nil {
		return fmt.Errorf("save Vault token: %w", err)
	}
	return nil
}

// privateKeyURI matches keys by certificate ID without copying certificate labels.
func privateKeyURI(uri string) string {
	path, query, hasQuery := strings.Cut(strings.TrimPrefix(uri, "pkcs11:"), "?")
	attributes := strings.Split(path, ";")
	hasID := false
	for _, attribute := range attributes {
		if strings.HasPrefix(attribute, "id=") && len(attribute) > len("id=") {
			hasID = true
		}
	}
	keyAttributes := make([]string, 0, len(attributes))
	for _, attribute := range attributes {
		if hasID && strings.HasPrefix(attribute, "object=") {
			continue
		}
		if attribute == "type=cert" {
			attribute = "type=private"
		}
		keyAttributes = append(keyAttributes, attribute)
	}
	result := "pkcs11:" + strings.Join(keyAttributes, ";")
	if hasQuery {
		result += "?" + query
	}
	return result
}

// certificateCandidate keeps public identity metadata private during selection.
type certificateCandidate struct {
	// uri identifies the hardware certificate and its associated private key.
	uri string
	// certificate supplies validity, usage, and identity metadata.
	certificate *x509.Certificate
}

// discoverCertificate selects a usable hardware client certificate for the role.
func discoverCertificate(ctx context.Context, role string) (string, error) {
	output, err := runNative(ctx, "p11tool", []string{"--list-tokens", "--batch"}, nil, 64<<10, nil)
	if err != nil {
		return "", fmt.Errorf("list PKCS#11 tokens: %w", err)
	}
	tokens := hardwareTokenURIs(string(output))
	if len(tokens) > 32 {
		return "", errors.New("too many hardware tokens for automatic discovery")
	}
	var candidates, matching []certificateCandidate
	seen := make(map[string]bool)
	now := time.Now()
	for _, token := range tokens {
		certificates, err := runNative(ctx, "p11tool", []string{"--only-urls", "--list-all-certs", "--batch", token}, nil, 64<<10, nil)
		if err != nil {
			return "", fmt.Errorf("list hardware certificates: %w", err)
		}
		for _, uri := range strings.Fields(string(certificates)) {
			if !strings.HasPrefix(uri, "pkcs11:") || seen[uri] {
				continue
			}
			seen[uri] = true
			if len(seen) > 64 {
				return "", errors.New("too many hardware certificates for automatic discovery")
			}
			encoded, err := runNative(ctx, "p11tool", []string{"--export", "--batch", uri}, nil, 64<<10, nil)
			if err != nil {
				return "", fmt.Errorf("read public hardware certificate: %w", err)
			}
			if block, _ := pem.Decode(encoded); block != nil {
				encoded = block.Bytes
			}
			certificate, err := x509.ParseCertificate(encoded)
			if err != nil {
				return "", fmt.Errorf("parse hardware certificate: %w", privateError{cause: err})
			}
			if !usableClientCertificate(certificate, now) {
				continue
			}
			candidate := certificateCandidate{uri: uri, certificate: certificate}
			candidates = append(candidates, candidate)
			identities := append([]string{certificate.Subject.CommonName}, certificate.DNSNames...)
			for _, identity := range identities {
				if identity == role || strings.HasPrefix(identity, role+".") {
					matching = append(matching, candidate)
					break
				}
			}
		}
	}
	if len(matching) != 0 {
		candidates = matching
	} else if len(candidates) > 1 {
		return "", errors.New("multiple usable hardware certificates do not match the role identity; use --certificate-uri")
	}
	if len(candidates) == 0 {
		return "", errors.New("no currently valid hardware client certificate was found")
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].certificate.NotBefore.After(candidates[j].certificate.NotBefore)
	})
	return candidates[0].uri, nil
}

// hardwareTokenURIs ignores trust modules regardless of their enumeration order.
func hardwareTokenURIs(output string) []string {
	var result []string
	var uri string
	hardware := false
	for _, line := range append(strings.Split(output, "\n"), "Token end:") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Token ") {
			if hardware && uri != "" {
				result = append(result, uri)
			}
			uri, hardware = "", false
		} else if strings.HasPrefix(line, "URL:") {
			uri = strings.TrimSpace(strings.TrimPrefix(line, "URL:"))
		} else if line == "Type: Hardware token" {
			hardware = true
		}
	}
	return result
}

// usableClientCertificate excludes trust anchors, expired keys, and server-only use.
func usableClientCertificate(certificate *x509.Certificate, now time.Time) bool {
	if certificate.IsCA || now.Before(certificate.NotBefore) || now.After(certificate.NotAfter) ||
		(certificate.KeyUsage != 0 && certificate.KeyUsage&x509.KeyUsageDigitalSignature == 0) {
		return false
	}
	if len(certificate.ExtKeyUsage) == 0 && len(certificate.UnknownExtKeyUsage) == 0 {
		return true
	}
	for _, usage := range certificate.ExtKeyUsage {
		if usage == x509.ExtKeyUsageClientAuth || usage == x509.ExtKeyUsageAny {
			return true
		}
	}
	return false
}

// runNative connects argv and private bounded output without invoking a shell.
func runNative(ctx context.Context, name string, args []string, input []byte, limit int, diagnostics io.Writer) ([]byte, error) {
	command := exec.CommandContext(ctx, name, args...)
	command.Stdin = bytes.NewReader(input)
	command.Stderr = io.Discard
	if diagnostics != nil {
		command.Stderr = diagnostics
	}
	command.WaitDelay = time.Second
	output := boundedOutput{limit: limit}
	command.Stdout = &output
	if err := command.Run(); err != nil {
		return nil, fmt.Errorf("execute %s: %w", name, err)
	}
	return output.buffer.Bytes(), nil
}

// validateDestination rejects nonregular existing files before authentication.
func validateDestination(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("inspect token file: %w", err)
	}
	if !info.Mode().IsRegular() {
		return errors.New("existing token destination must be a regular file, not a symbolic link")
	}
	return nil
}

// saveToken atomically replaces the destination using a same-directory 0600 file.
func saveToken(path, token string) (result error) {
	if err := validateDestination(path); err != nil {
		return fmt.Errorf("recheck token destination: %w", err)
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".vault-token-*")
	if err != nil {
		return fmt.Errorf("create private token file: %w", err)
	}
	closed, installed := false, false
	defer func() {
		if !closed {
			if err := file.Close(); err != nil {
				result = errors.Join(result, fmt.Errorf("close temporary token file: %w", err))
			}
		}
		if !installed {
			if err := os.Remove(file.Name()); err != nil {
				result = errors.Join(result, fmt.Errorf("remove temporary token file: %w", err))
			}
		}
	}()
	if err := file.Chmod(0o600); err != nil {
		return fmt.Errorf("set private token permissions: %w", err)
	}
	if _, err := io.WriteString(file, token+"\n"); err != nil {
		return fmt.Errorf("write private token file: %w", err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("synchronize private token file: %w", err)
	}
	err = file.Close()
	closed = true
	if err != nil {
		return fmt.Errorf("close private token file: %w", err)
	}
	if err := os.Rename(file.Name(), path); err != nil {
		return fmt.Errorf("replace token file: %w", err)
	}
	installed = true
	return nil
}

package main

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	vault "github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

// withheldError preserves an upstream cause while excluding secret-bearing text.
type withheldError struct {
	// cause is available to errors.Is/As but never rendered in command output.
	cause error
}

// Error reports an operational failure without response bodies or credentials.
func (e withheldError) Error() string { return "upstream details withheld" }

// Unwrap preserves error identity for callers.
func (e withheldError) Unwrap() error { return e.cause }

// validRecord verifies the complete cached Ed25519 pair without exposing it.
func validRecord(data map[string]any) error {
	private, _ := data["private_key"].(string)
	public, _ := data["public_key"].(string)
	fingerprint, _ := data["fingerprint"].(string)
	if data["algorithm"] != "ssh-ed25519" {
		return errors.New("cached SSH algorithm is invalid")
	}
	parsed, err := ssh.ParseRawPrivateKey([]byte(private))
	if err != nil {
		return fmt.Errorf("parse cached private key: %w", withheldError{err})
	}
	key, ok := parsed.(*ed25519.PrivateKey)
	if !ok {
		return errors.New("cached private key is not Ed25519")
	}
	if !bytes.Equal(*key, ed25519.NewKeyFromSeed(key.Seed())) {
		return errors.New("cached private key is internally inconsistent")
	}
	derived, err := ssh.NewPublicKey(key.Public())
	if err != nil {
		return fmt.Errorf("derive cached public key: %w", withheldError{err})
	}
	if public != string(ssh.MarshalAuthorizedKey(derived)) || fingerprint != ssh.FingerprintSHA256(derived) {
		return errors.New("cached SSH pair or fingerprint does not match")
	}
	return nil
}

// ensure reuses an intact SSH key record or creates a new pair exactly once.
func ensure(ctx context.Context) (bool, error) {
	mount, path := os.Getenv("VAULT_SSH_KEY_MOUNT"), os.Getenv("VAULT_SSH_KEY_PATH")
	pattern := regexp.MustCompile(`^[A-Za-z0-9_-][A-Za-z0-9_.-]*(/[A-Za-z0-9_-][A-Za-z0-9_.-]*)*$`)
	if !pattern.MatchString(mount) || strings.Contains(mount, "/") || !pattern.MatchString(path) {
		return false, errors.New("set a valid VAULT_SSH_KEY_MOUNT and VAULT_SSH_KEY_PATH")
	}
	address, err := url.Parse(os.Getenv("VAULT_ADDR"))
	if err != nil || address.Scheme != "https" || address.Host == "" || address.User != nil || address.RawQuery != "" || address.Fragment != "" {
		return false, errors.New("VAULT_ADDR must be an HTTPS URL")
	}
	if token := os.Getenv("VAULT_TOKEN"); token == "" || strings.ContainsAny(token, "\r\n") {
		return false, errors.New("missing or invalid VAULT_TOKEN")
	}
	if skip := os.Getenv("VAULT_SKIP_VERIFY"); skip != "" && skip != "false" && skip != "0" {
		return false, errors.New("Vault TLS verification must be enabled")
	}
	configuration := vault.DefaultConfig()
	configuration.DisableRedirects = true
	configuration.HttpClient.CheckRedirect = func(*http.Request, []*http.Request) error { return errors.New("Vault redirect rejected") }
	client, err := vault.NewClient(configuration)
	if err != nil {
		return false, fmt.Errorf("configure Vault: %w", withheldError{err})
	}
	client.SetMaxRetries(0)
	kv := client.KVv2(mount)
	stored, err := kv.Get(ctx, path)
	if err == nil {
		if err := validRecord(stored.Data); err != nil {
			return false, fmt.Errorf("validate cached SSH identity: %w", err)
		}
		return true, nil
	}
	if !errors.Is(err, vault.ErrSecretNotFound) {
		return false, fmt.Errorf("read SSH key cache: %w", withheldError{err})
	}
	metadata, err := client.Logical().ReadWithContext(ctx, mount+"/metadata/"+path)
	if err != nil {
		return false, fmt.Errorf("check SSH key history: %w", withheldError{err})
	}
	if stored != nil || metadata != nil {
		return false, errors.New("SSH key history exists without usable data; restore the original identity")
	}
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return false, fmt.Errorf("generate SSH key pair: %w", err)
	}
	encoded, err := ssh.MarshalPrivateKey(private, "Vault-managed Ed25519 key")
	if err != nil {
		return false, fmt.Errorf("encode SSH private key: %w", withheldError{err})
	}
	publicKey, err := ssh.NewPublicKey(public)
	if err != nil {
		return false, fmt.Errorf("encode SSH public key: %w", withheldError{err})
	}
	data := map[string]any{"algorithm": "ssh-ed25519", "private_key": string(pem.EncodeToMemory(encoded)), "public_key": string(ssh.MarshalAuthorizedKey(publicKey)), "fingerprint": ssh.FingerprintSHA256(publicKey)}
	if _, err := kv.Put(ctx, path, data, vault.WithCheckAndSet(0)); err != nil {
		return false, fmt.Errorf("store SSH key pair with CAS=0: %w", withheldError{err})
	}
	return false, nil
}

// command supplies a bounded, cancellation-aware IaC bootstrap entry point.
func command() *cobra.Command {
	root := &cobra.Command{Use: "ssh_key", Short: "Create or reuse an Ed25519 SSH key pair in Vault", SilenceErrors: true, SilenceUsage: true}
	root.AddCommand(&cobra.Command{Use: "ensure", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, cancel := context.WithTimeout(cmd.Context(), time.Minute)
		defer cancel()
		cached, err := ensure(ctx)
		if err != nil {
			return fmt.Errorf("ensure SSH key pair: %w", err)
		}
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "SSH key pair stored in Vault (cached=%t).\n", cached); err != nil {
			return fmt.Errorf("write completion: %w", err)
		}
		return nil
	}})
	return root
}

// main reports sanitized failures and propagates cancellation.
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if err := command().ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

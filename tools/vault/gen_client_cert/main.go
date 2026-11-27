package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"git.alwaldend.com/alwaldend/src/tools/go/pkg/sensitiveerr"
	"git.alwaldend.com/alwaldend/src/tools/vault/internal/privatefile"
	"git.alwaldend.com/alwaldend/src/tools/vault/pkg/vaultcmd"
	"github.com/spf13/cobra"
)

// options selects the PKI request and the local credential destination.
type options struct {
	// vault identifies the selected Vault executable.
	vault string
	// mount identifies the PKI secrets engine.
	mount string
	// role authorizes issuance for the requested identity.
	role string
	// directory contains the issued material.
	directory string
	// user selects a user identity, optionally scoped to a host.
	user string
	// host selects a host certificate or scopes a user identity.
	host string
	// cached selects previously issued material.
	cached bool
	// ttl is the requested Vault certificate lifetime.
	ttl string
}

// main parses PKI options while preserving trailing Vault write arguments.
func main() {
	var o options
	command := &cobra.Command{
		Use: "gen_client_cert [vault-write-arguments...]", Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.generate(cmd.Context(), args); err != nil {
				return fmt.Errorf("generate client certificate: %w", err)
			}
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&o.vault, "vault", "", "Vault executable")
	flags.StringVar(&o.mount, "path", "pki/ica_clients", "Mount path")
	flags.StringVar(&o.role, "role", "pki_ica_clients_users", "Role name")
	flags.StringVar(&o.directory, "output_dir", "", "Output directory")
	flags.StringVar(&o.user, "user", "", "Username")
	flags.StringVar(&o.host, "host", "", "Hostname")
	flags.BoolVar(&o.cached, "no_create", false, "Use a previously issued certificate")
	flags.StringVar(&o.ttl, "ttl", "28927206", "TTL")
	if err := vaultcmd.Execute(command, "output_dir"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// generate issues or loads a certificate and writes its PEM and PKCS#12 representations.
func (o options) generate(ctx context.Context, passthrough []string) error {
	if o.user == "" && o.host == "" {
		return errors.New("either --user or --host is required")
	}
	if o.directory == "" || o.mount == "" || o.role == "" || o.ttl == "" {
		return errors.New("output directory, mount, role and TTL must be nonempty")
	}
	target := o.user
	commonName := o.user + ".users.alwaldend.com"
	if o.host != "" {
		target = o.host
		if o.user == "" {
			commonName = o.host + ".infra.alwaldend.com"
		} else {
			commonName = o.host + "." + commonName
		}
	}
	if target == "." || !filepath.IsLocal(target) || filepath.Base(target) != target {
		return errors.New("identity must be a single file name")
	}
	if !o.cached && o.vault == "" {
		return errors.New("--vault is required when requesting a new certificate")
	}
	if err := os.MkdirAll(o.directory, 0o700); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	base := filepath.Join(o.directory, target)
	var raw []byte
	var err error
	if o.cached {
		raw, err = readRecord(base + ".json")
		if err != nil {
			return fmt.Errorf("load cached certificate: %w", err)
		}
	} else {
		args := append([]string{"write", "-format=json", o.mount + "/issue/" + o.role, "common_name=" + commonName, "ttl=" + o.ttl}, passthrough...)
		raw, err = vaultcmd.Output(ctx, o.vault, nil, args...)
		if err != nil {
			return fmt.Errorf("issue certificate: %w", err)
		}
	}
	var record struct {
		// Data contains the key pair and chain issued by Vault.
		Data struct {
			// Certificate contains the leaf certificate in PEM.
			Certificate string `json:"certificate"`
			// Chain lists the issuing certificates in PEM.
			Chain []string `json:"ca_chain"`
			// PrivateKey contains the leaf certificate's private key in PEM.
			PrivateKey string `json:"private_key"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &record); err != nil {
		return sensitiveerr.WithContext("decode certificate response", err)
	}
	if record.Data.Certificate == "" || record.Data.PrivateKey == "" {
		return errors.New("certificate record is missing certificate or private key")
	}
	if !o.cached {
		if err := privatefile.Write(base+".json", raw); err != nil {
			return fmt.Errorf("save issued certificate record: %w", err)
		}
	}
	combined := append([]string{record.Data.Certificate}, record.Data.Chain...)
	combined = append(combined, record.Data.PrivateKey, "")
	for _, file := range []struct {
		// suffix selects a certificate representation.
		suffix string
		// data contains the encoded representation.
		data string
	}{
		{".pub", record.Data.Certificate},
		{".key", record.Data.PrivateKey},
		{"_chain.pub", strings.Join(record.Data.Chain, "\n")},
		{".pem", strings.Join(combined, "\n")},
	} {
		if err := privatefile.Write(base+file.suffix, []byte(file.data)); err != nil {
			return fmt.Errorf("save %s file: %w", file.suffix, err)
		}
	}
	// Capture the bundle before replacing its destination so failed conversion
	// preserves any previously issued PKCS#12 file.
	bundle, err := vaultcmd.Output(ctx, "openssl", nil, "pkcs12", "-export", "-certfile", base+"_chain.pub", "-in", base+".pub", "-inkey", base+".key")
	if err != nil {
		return fmt.Errorf("encode PKCS#12 bundle: %w", err)
	}
	if err := privatefile.Write(base+".pfx", bundle); err != nil {
		return fmt.Errorf("save PKCS#12 bundle: %w", err)
	}
	return nil
}

// readRecord bounds memory use when reloading a previously issued certificate.
func readRecord(path string) (data []byte, result error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open certificate record: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("close certificate record: %w", err))
		}
	}()
	data, err = io.ReadAll(io.LimitReader(file, vaultcmd.OutputLimit+1))
	if err != nil {
		return nil, fmt.Errorf("read certificate record: %w", err)
	}
	if len(data) > vaultcmd.OutputLimit {
		return nil, errors.New("certificate record exceeds 1 MiB")
	}
	return data, nil
}

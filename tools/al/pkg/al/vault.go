package al

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"git.alwaldend.com/alwaldend/src/tools/al/api/al_proto"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"github.com/hashicorp/vault/api/tokenhelper"
	"google.golang.org/protobuf/proto"
)

const VAULT_DEFAULT_NAME = "default"

type Vault struct {
	helper  tokenhelper.TokenHelper
	clients map[string]*api.Client
	config  *al_proto.Config
}

func NewVault(config *al_proto.Config) (*Vault, error) {
	helper, err := tokenhelper.NewInternalTokenHelper()
	if err != nil {
		return nil, fmt.Errorf("could not create the internal token helper")
	}
	res := &Vault{
		helper:  helper,
		clients: map[string]*api.Client{},
		config:  config,
	}
	return res, nil
}

func (self *Vault) FetchSecret(ctx context.Context, name string) (map[string]any, error) {
	secret, err := VaultSecretByName(self.config, name)
	if err != nil {
		return nil, fmt.Errorf("could not find secret %s: %w", name, err)
	}
	secret, err = self.normalizeSecret(secret)
	if err != nil {
		return nil, fmt.Errorf("could not normalize secret %s: %w", name, err)
	}
	client, err := self.clientForSecret(ctx, secret)
	if err != nil {
		return nil, fmt.Errorf("could not get a client for the secret %s: %w", secret.Name, err)
	}
	data, err := client.KVv2(secret.Kv.Mount).Get(ctx, secret.Kv.Path)
	if err != nil {
		return nil, fmt.Errorf("could not fetch secret %s: %w", secret.Name, err)
	}
	return data.Data, nil
}

func (self *Vault) DefaultEnv(ctx context.Context) ([]string, error) {
	res, err := self.Env(ctx, VAULT_DEFAULT_NAME, VAULT_DEFAULT_NAME)
	return res, err
}

// Create vault environment variables
// https://developer.hashicorp.com/vault/docs/commands#configure-environment-variables
func (self *Vault) Env(ctx context.Context, vaultName string, authName string) ([]string, error) {
	if vaultName == "" {
		vaultName = VAULT_DEFAULT_NAME
	}
	if authName == "" {
		authName = VAULT_DEFAULT_NAME
	}
	vault, err := VaultByName(self.config, vaultName)
	if err != nil {
		return nil, fmt.Errorf("missing vault: %w", err)
	}
	auth, err := VaultAuthByName(self.config, authName)
	if err != nil {
		return nil, fmt.Errorf("missing auth: %w", err)
	}
	tlsConfig, err := tlsConfig(vault)
	if err != nil {
		return nil, fmt.Errorf("could not create tls config: %w", err)
	}
	res := []string{
		fmt.Sprintf("VAULT_ADDR=%s", vault.Config.Address),
	}
	if tlsConfig.CACert != "" {
		res = append(res, fmt.Sprintf("VAULT_CACERT=%s", tlsConfig.CACert))
	}
	if tlsConfig.ClientKey != "" {
		res = append(res, fmt.Sprintf("VAULT_CLIENT_KEY=%s", tlsConfig.ClientKey))
	}
	if tlsConfig.ClientCert != "" {
		res = append(res, fmt.Sprintf("VAULT_CLIENT_CERT=%s", tlsConfig.ClientCert))
	}
	if !auth.NoAuth {
		client, err := self.client(ctx, vaultName, authName)
		if err != nil {
			return nil, fmt.Errorf("could not create client for %s/%s: %w", vaultName, authName, err)
		}
		res = append(res, fmt.Sprintf("VAULT_TOKEN=%s", client.Token()))
	}
	return res, nil
}

func (self *Vault) normalizeSecret(secret *al_proto.VaultSecret) (*al_proto.VaultSecret, error) {
	secret = proto.CloneOf(secret)
	if secret.Name == "" {
		return nil, fmt.Errorf("secret is missing a name")
	}
	kv := secret.Kv
	if kv == nil {
		return nil, fmt.Errorf("secret %s is missing kv config", secret.Name)
	}
	if kv.Path == "" {
		return nil, fmt.Errorf("secret %s is missing path", secret.Name)
	}
	if secret.Auth == "" {
		secret.Auth = VAULT_DEFAULT_NAME
	}
	if secret.Vault == "" {
		secret.Vault = VAULT_DEFAULT_NAME
	}
	if kv.Mount == "" {
		kv.Mount = "secrets"
	}
	return secret, nil
}

func (self *Vault) client(ctx context.Context, vaultName string, authName string) (*api.Client, error) {
	path := fmt.Sprintf("%s/%s", vaultName, authName)
	client, ok := self.clients[path]
	if ok {
		return client, nil
	}
	auth, err := VaultAuthByName(self.config, authName)
	if err != nil {
		return nil, fmt.Errorf("could not get auth config: %w", err)
	}
	vault, err := VaultByName(self.config, vaultName)
	if err != nil {
		return nil, fmt.Errorf("could not get vault config: %w", err)
	}
	client, err = newVaultClient(ctx, self.helper, vault, auth)
	if err != nil {
		return nil, fmt.Errorf("could not create vault client %s: %w", path, err)
	}
	self.clients[path] = client
	return client, nil
}
func (self *Vault) clientForSecret(ctx context.Context, secret *al_proto.VaultSecret) (*api.Client, error) {
	client, err := self.client(ctx, secret.Vault, secret.Auth)
	if err != nil {
		return nil, fmt.Errorf("could not create a client for secret %s: %w", secret.Name, err)
	}
	return client, nil
}

func vaultAuthDefault(ctx context.Context, config *al_proto.Config) (*api.Client, error) {
	res, err := vaultAuth(ctx, config, VAULT_DEFAULT_NAME, VAULT_DEFAULT_NAME)
	return res, err
}

func vaultAuth(ctx context.Context, config *al_proto.Config, vaultName string, authName string) (*api.Client, error) {
	vault, err := VaultByName(config, vaultName)
	if err != nil {
		return nil, fmt.Errorf("missing vault: %w", err)
	}
	auth, err := VaultAuthByName(config, authName)
	if err != nil {
		return nil, fmt.Errorf("missing auth: %w", err)
	}
	tokenHelper, err := tokenhelper.NewInternalTokenHelper()
	if err != nil {
		return nil, fmt.Errorf("could not create token helper: %w", err)
	}
	client, err := newVaultClient(ctx, tokenHelper, vault, auth)
	if err != nil {
		return nil, fmt.Errorf("could not create vault client: %w", err)
	}
	return client, nil
}

func tlsConfig(vault *al_proto.Vault) (*api.TLSConfig, error) {
	res := &api.TLSConfig{}
	if vault.Tls.CaCert != "" {
		cacert, err := filepath.Abs(os.ExpandEnv(vault.Tls.CaCert))
		if err != nil {
			return nil, fmt.Errorf("could not expand cacert: %w", err)
		}
		res.CACert = cacert
	}
	if vault.Tls.ClientCert != "" {
		clientCert, err := filepath.Abs(os.ExpandEnv(vault.Tls.ClientCert))
		if err != nil {
			return nil, fmt.Errorf("could not expand client cert: %w", err)
		}
		res.ClientCert = clientCert
	}
	if vault.Tls.ClientKey != "" {
		clientKey, err := filepath.Abs(os.ExpandEnv(vault.Tls.ClientKey))
		if err != nil {
			return nil, fmt.Errorf("could not expand client key: %w", err)
		}
		res.ClientKey = clientKey
	}
	return res, nil
}

func newVaultClient(ctx context.Context, tokenHelper tokenhelper.TokenHelper, vault *al_proto.Vault, auth *al_proto.VaultAuth) (*api.Client, error) {
	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = vault.Config.Address
	tlsConfig, err := tlsConfig(vault)
	if err != nil {
		return nil, fmt.Errorf("could not create tls config: %w", err)
	}
	err = vaultConfig.ConfigureTLS(tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("could not configure tls: %w", err)
	}
	client, err := api.NewClient(vaultConfig)
	if err != nil {
		return nil, fmt.Errorf("could not create vault client: %w", err)
	}
	token, err := tokenHelper.Get()
	if err != nil {
		return nil, fmt.Errorf("could not get token from the token helper: %w", err)
	}
	client.SetToken(token)
	if auth.Approle != nil {
		approleData, err := client.Logical().WriteWithContext(ctx, fmt.Sprintf("auth/approle/role/%s/secret-id", auth.Approle.Name), nil)
		if err != nil {
			return nil, fmt.Errorf("could not create secret id for the approle: %w", err)
		}
		secretId, ok := approleData.Data["secret_id"].(string)
		if !ok {
			return nil, fmt.Errorf("missing secret_id for some reason")
		}
		appRoleAuth, err := approle.NewAppRoleAuth(
			auth.Approle.Name,
			&approle.SecretID{FromString: secretId},
		)
		if err != nil {
			return nil, fmt.Errorf("could not create approle auth: %w", err)
		}
		_, err = client.Auth().Login(ctx, appRoleAuth)
		if err != nil {
			return nil, fmt.Errorf("could not auth using the approle: %w", err)
		}
	}
	return client, nil
}

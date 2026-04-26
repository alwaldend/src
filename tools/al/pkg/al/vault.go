package al

import (
	"context"
	"fmt"
	"path/filepath"

	"git.alwaldend.com/alwaldend/src/tools/al/api/al_proto"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"github.com/hashicorp/vault/api/tokenhelper"
)

const VAULT_DEFAULT_NAME = "default"

func VaultFetchSecrets(ctx context.Context, config *al_proto.Config) (map[string]map[string]any, error) {
	var err error
	clients := make(map[string]*api.Client)
	tokenHelper, err := tokenhelper.NewInternalTokenHelper()
	if err != nil {
		return nil, fmt.Errorf("could not create token helper: %w", err)
	}
	res := make(map[string]map[string]any)
	for i, secret := range config.Secrets {
		if secret.Name == "" {
			return nil, fmt.Errorf("secret %d is missing name", i)
		}
		kv := secret.Kv
		if kv == nil {
			return nil, fmt.Errorf("secret %s is missing kv config", secret.Name)
		}
		if kv.Path == "" {
			return nil, fmt.Errorf("secret %s is missing path", secret.Name)
		}
		auth := secret.Auth
		if auth == "" {
			auth = "default"
		}
		vault := secret.Vault
		if vault == "" {
			vault = "default"
		}
		mount := kv.Mount
		if mount == "" {
			mount = "secrets"
		}
		path := fmt.Sprintf("%s/%s", vault, auth)
		client, ok := clients[path]
		if !ok {
			var vaultCfg *al_proto.Vault
			var vaultAuthCfg *al_proto.VaultAuth
			for _, cfg := range config.Vaults {
				if cfg.Name == vault {
					vaultCfg = cfg
					break
				}
			}
			if vaultCfg == nil {
				return nil, fmt.Errorf("missing config for vault %s", vault)
			}
			for _, cfg := range config.Auth {
				if cfg.Name == auth {
					vaultAuthCfg = cfg
					break
				}
			}
			if vaultAuthCfg == nil {
				return nil, fmt.Errorf("missing auth config with name %s", auth)
			}
			client, err = NewVaultClient(ctx, tokenHelper, vaultCfg, vaultAuthCfg)
			if err != nil {
				return nil, fmt.Errorf("could not create vault client %s: %w", path, err)
			}
		}
		data, err := client.KVv2(mount).Get(ctx, kv.Path)
		if err != nil {
			return nil, fmt.Errorf("could not fetch secret %s: %w", secret.Name, err)
		}
		res[secret.Name] = data.Data
	}
	return res, nil
}

func VaultDefaultEnv(config *al_proto.Config, client *api.Client) ([]string, error) {
	res, err := VaultEnv(config, client, VAULT_DEFAULT_NAME)
	return res, err
}

func VaultEnv(config *al_proto.Config, client *api.Client, vaultName string) ([]string, error) {
	vault, err := VaultByName(config, vaultName)
	if err != nil {
		return nil, fmt.Errorf("missing vault: %w", err)
	}
	cacert, err := filepath.Abs( vault.Tls.CaCert)
	if err != nil {
		return nil, fmt.Errorf("could not create absolute path: %w", err)
	}
	res := []string{
		fmt.Sprintf("VAULT_TOKEN=%s", client.Token()),
		fmt.Sprintf("VAULT_ADDR=%s", client.Address()),
		fmt.Sprintf("VAULT_CACERT=%s", cacert),
	}
	return res, nil
}

func VaultByName(config *al_proto.Config, name string) (*al_proto.Vault, error) {
	var vault *al_proto.Vault
	for _, curVault := range config.Vaults {
		if curVault.Name == name {
			vault = curVault
			break
		}
	}
	if vault == nil {
		return nil, fmt.Errorf("missing vault with name %s", name)
	}
	return vault, nil
}

func VaultAuthByName(config *al_proto.Config, name string) (*al_proto.VaultAuth, error) {
	var auth *al_proto.VaultAuth
	for _, curAuth := range config.Auth {
		if curAuth.Name == name {
			auth = curAuth
			break
		}
	}
	if auth == nil {
		return nil, fmt.Errorf("missing vault auth with name %s", name)
	}
	return auth, nil
}

func VaultAuthDefault(ctx context.Context, config *al_proto.Config) (*api.Client, error) {
	res, err := VaultAuth(ctx, config, VAULT_DEFAULT_NAME, VAULT_DEFAULT_NAME)
	return res, err
}

func VaultAuth(ctx context.Context, config *al_proto.Config, vaultName string, authName string) (*api.Client, error) {
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
	client, err := NewVaultClient(ctx, tokenHelper, vault, auth)
	if err != nil {
		return nil, fmt.Errorf("could not create vault client: %w", err)
	}
	return client, nil
}

func NewVaultClient(ctx context.Context, tokenHelper tokenhelper.TokenHelper, vault *al_proto.Vault, auth *al_proto.VaultAuth) (*api.Client, error) {
	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = vault.Config.Address
	err := vaultConfig.ConfigureTLS(&api.TLSConfig{CACert: vault.Tls.CaCert})
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
	return client, nil
}

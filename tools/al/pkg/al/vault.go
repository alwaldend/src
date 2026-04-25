package al

import (
	"context"
	"fmt"

	"git.alwaldend.com/alwaldend/src/tools/al/api/al_proto"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"github.com/hashicorp/vault/api/tokenhelper"
)

func FetchSecrets(ctx context.Context, config *al_proto.Config) (map[string]map[string]any, error) {
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

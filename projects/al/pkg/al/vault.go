package al

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"github.com/hashicorp/vault/api/tokenhelper"
	"sync"
)

const VAULT_DEFAULT_NAME = "default"

type VaultStoreItem struct {
	Helper tokenhelper.TokenHelper
	Client *api.Client
}

type VaultStore struct {
	clients map[string]*VaultStoreItem
	config  *al_proto.Config
	mx      *sync.RWMutex
}

func NewVault(config *al_proto.Config) *VaultStore {
	return &VaultStore{
		clients: map[string]*VaultStoreItem{},
		config:  config,
		mx:      &sync.RWMutex{},
	}

}

func (self *VaultStore) DefaultEnv(ctx context.Context) fp.Either[[]string] {
	return self.Env(ctx, VAULT_DEFAULT_NAME, VAULT_DEFAULT_NAME)
}

// Create vault environment variables
// https://developer.hashicorp.com/vault/docs/commands#configure-environment-variables
func (self *VaultStore) Env(ctx context.Context, vaultName string, authName string) fp.Either[[]string] {
	if vaultName == "" {
		vaultName = VAULT_DEFAULT_NAME
	}
	if authName == "" {
		authName = VAULT_DEFAULT_NAME
	}
	vault, err := VaultByName(self.config, vaultName)
	if err != nil {
		return fp.Left[[]string](fmt.Errorf("missing vault: %w", err))
	}
	auth, err := VaultAuthByName(self.config, authName)
	if err != nil {
		return fp.Left[[]string](fmt.Errorf("missing auth: %w", err))
	}
	tlsConfig, err := fp.Get(tlsConfig(vault))
	if err != nil {
		return fp.Left[[]string](fmt.Errorf("could not create tls config: %w", err))
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
		client, err := fp.Get(self.Client(ctx, vaultName, authName))
		if err != nil {
			return fp.Left[[]string](fmt.Errorf("could not create client for %s/%s: %w", vaultName, authName, err))
		}
		res = append(res, fmt.Sprintf("VAULT_TOKEN=%s", client.Client.Token()))
	}
	return fp.Right(res)
}

func (self *VaultStore) clientCache(path string) (*VaultStoreItem, bool) {
	self.mx.RLock()
	defer self.mx.RUnlock()
	client, ok := self.clients[path]
	if ok {
		return client, true
	}
	return nil, false
}

func (self *VaultStore) Client(ctx context.Context, conn string, authName string) fp.Either[*VaultStoreItem] {
	if conn == "" {
		conn = VAULT_DEFAULT_NAME
	}
	if authName == "" {
		authName = VAULT_DEFAULT_NAME
	}
	path := fmt.Sprintf("%s/%s", conn, authName)
	c, ok := self.clientCache(path)
	if ok {
		return fp.Right(c)
	}
	auth, err := VaultAuthByName(self.config, authName)
	if err != nil {
		return fp.Left[*VaultStoreItem](fmt.Errorf("could not get auth config: %w", err))
	}
	vault, err := VaultByName(self.config, conn)
	if err != nil {
		return fp.Left[*VaultStoreItem](fmt.Errorf("could not get vault config: %w", err))
	}
	client, err := newVaultClient(ctx, vault, auth).Get()
	if err != nil {
		return fp.Left[*VaultStoreItem](fmt.Errorf("could not create vault client %s: %w", path, err))
	}
	self.clients[path] = client
	return fp.Right(client)
}

func tlsConfig(vault *al_proto.VaultConn) fp.Either[*api.TLSConfig] {
	res := &api.TLSConfig{}
	if vault.Tls.CaCert != "" {
		cacert, err := filepath.Abs(os.ExpandEnv(vault.Tls.CaCert))
		if err != nil {
			return fp.Left[*api.TLSConfig](fmt.Errorf("could not expand cacert: %w", err))
		}
		res.CACert = cacert
	}
	if vault.Tls.ClientCert != "" {
		clientCert, err := filepath.Abs(os.ExpandEnv(vault.Tls.ClientCert))
		if err != nil {
			return fp.Left[*api.TLSConfig](fmt.Errorf("could not expand client cert: %w", err))
		}
		res.ClientCert = clientCert
	}
	if vault.Tls.ClientKey != "" {
		clientKey, err := filepath.Abs(os.ExpandEnv(vault.Tls.ClientKey))
		if err != nil {
			return fp.Left[*api.TLSConfig](fmt.Errorf("could not expand client key: %w", err))
		}
		res.ClientKey = clientKey
	}
	return fp.Right(res)
}

func newVaultClient(ctx context.Context, vault *al_proto.VaultConn, auth *al_proto.VaultAuth) fp.Either[*VaultStoreItem] {
	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = vault.Config.Address
	tlsConfig, err := tlsConfig(vault).Get()
	if err != nil {
		return fp.Left[*VaultStoreItem](fmt.Errorf("could not create tls config: %w", err))
	}
	err = vaultConfig.ConfigureTLS(tlsConfig)
	if err != nil {
		return fp.Left[*VaultStoreItem](fmt.Errorf("could not configure tls: %w", err))
	}
	client, err := api.NewClient(vaultConfig)
	if err != nil {
		return fp.Left[*VaultStoreItem](fmt.Errorf("could not create vault client: %w", err))
	}
	helper, err := tokenhelper.NewInternalTokenHelper()
	if err != nil {
		return fp.Left[*VaultStoreItem](fmt.Errorf("could not create the token helper: %w", err))
	}
	token, err := helper.Get()
	if err != nil {
		return fp.Left[*VaultStoreItem](fmt.Errorf("could not get token from the token helper: %w", err))
	}
	client.SetToken(token)
	if auth.Approle != nil {
		approleData, err := client.Logical().WriteWithContext(ctx, fmt.Sprintf("auth/approle/role/%s/secret-id", auth.Approle.Name), nil)
		if err != nil {
			return fp.Left[*VaultStoreItem](fmt.Errorf("could not create secret id for the approle: %w", err))
		}
		secretId, ok := approleData.Data["secret_id"].(string)
		if !ok {
			return fp.Left[*VaultStoreItem](fmt.Errorf("missing secret_id for some reason"))
		}
		appRoleAuth, err := approle.NewAppRoleAuth(
			auth.Approle.Name,
			&approle.SecretID{FromString: secretId},
		)
		if err != nil {
			return fp.Left[*VaultStoreItem](fmt.Errorf("could not create approle auth: %w", err))
		}
		_, err = client.Auth().Login(ctx, appRoleAuth)
		if err != nil {
			return fp.Left[*VaultStoreItem](fmt.Errorf("could not auth using the approle: %w", err))
		}
	}
	return fp.Right(&VaultStoreItem{
		Client: client,
		Helper: helper,
	})
}

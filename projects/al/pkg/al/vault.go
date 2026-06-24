package al

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"sync"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"github.com/hashicorp/vault/api/tokenhelper"
	"regexp"
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

type VaultOidc struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

func NewVault(config *al_proto.Config) *VaultStore {
	return &VaultStore{
		clients: map[string]*VaultStoreItem{},
		config:  config,
		mx:      &sync.RWMutex{},
	}

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

func (self *VaultStore) OidcLogin(client *api.Client, oidcUrl *url.URL) (*VaultOidc, error) {
	regex, err := regexp.Compile("^https://([^/]+)/ui/vault/identity/oidc/provider/([^/]+)/authorize?(.+)$")
	if err != nil {
		return nil, fmt.Errorf("could not compile regex: %w", err)
	}
	parts := regex.FindAllStringSubmatch(oidcUrl.String(), -1)
	reqData := map[string]string{}
	for key, valueSlice := range oidcUrl.Query() {
		for _, value := range valueSlice {
			reqData[key] = value
		}
	}
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("could not marshal data: %w", err)
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://%s/v1/identity/oidc/provider/%s/authorize", parts[0][1], parts[0][2]),
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.Token()))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute the request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid response code: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read the body: %w", err)
	}
	res := &VaultOidc{}
	if err = json.Unmarshal(body, res); err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}
	return res, nil
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

func VaultTlsConfig(vault *al_proto.VaultConn) fp.Either[*api.TLSConfig] {
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
	tlsConfig, err := VaultTlsConfig(vault).Get()
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

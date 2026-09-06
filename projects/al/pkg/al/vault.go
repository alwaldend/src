package al

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"github.com/hashicorp/vault/api/tokenhelper"
)

const VAULT_DEFAULT_NAME = "default"

// SanitizeVaultError keeps cancellation and HTTP status diagnostics without
// logging server-controlled response bodies, request URLs, or credential values.
func SanitizeVaultError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) {
		return context.Canceled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return context.DeadlineExceeded
	}
	var response *api.ResponseError
	if errors.As(err, &response) {
		return fmt.Errorf("Vault request failed (HTTP %d)", response.StatusCode)
	}
	return fmt.Errorf("Vault request failed (%T)", err)
}

type VaultStoreItem struct {
	Helper tokenhelper.TokenHelper
	Client *api.Client
	owned  bool
}

type VaultStore struct {
	clients map[string]*VaultStoreItem
	config  *al_proto.Config
	mx      *sync.RWMutex
	closed  bool
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

func (self *VaultStore) Start(context.Context) error { return nil }

// Stop revokes only tokens created by this store. The user's token helper is
// neither erased nor revoked. Failed revocations remain available for retry.
func (self *VaultStore) Stop(ctx context.Context) error {
	self.mx.Lock()
	defer self.mx.Unlock()
	self.closed = true
	var errs []error
	for name, item := range self.clients {
		if item.owned {
			if _, err := item.Client.Logical().WriteWithContext(ctx, "auth/token/revoke-self", nil); err != nil {
				errs = append(errs, fmt.Errorf("could not revoke owned Vault token for %s", name))
				continue
			}
		}
		item.Client.ClearToken()
		delete(self.clients, name)
	}
	return errors.Join(errs...)
}

func (self *VaultStore) OidcLogin(client *api.Client, oidcUrl *url.URL) (*VaultOidc, error) {
	return self.OidcLoginContext(context.Background(), client, oidcUrl)
}

// OidcLoginContext authenticates only to the configured Vault HTTPS origin.
// Service-supplied URLs and redirects must never choose a bearer-token recipient.
func (self *VaultStore) OidcLoginContext(ctx context.Context, client *api.Client, oidcUrl *url.URL) (*VaultOidc, error) {
	origin, err := url.Parse(client.Address())
	if err != nil || oidcUrl == nil || origin.Scheme != "https" || oidcUrl.Scheme != "https" ||
		origin.User != nil || oidcUrl.User != nil || !strings.EqualFold(origin.Host, oidcUrl.Host) {
		return nil, fmt.Errorf("OIDC authorization URL must match the configured Vault HTTPS origin")
	}
	parts := regexp.MustCompile(`^/ui/vault/identity/oidc/provider/([^/]+)/authorize$`).FindStringSubmatch(oidcUrl.Path)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid Vault OIDC authorization path")
	}
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
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	endpoint := *origin
	endpoint.Path = "/v1/identity/oidc/provider/" + parts[1] + "/authorize"
	endpoint.RawPath, endpoint.RawQuery, endpoint.Fragment = "", "", ""
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint.String(),
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.Token()))
	httpClient := *client.CloneConfig().HttpClient
	httpClient.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute Vault OIDC authorization request")
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
	self.mx.Lock()
	defer self.mx.Unlock()
	if self.closed {
		return fp.Left[*VaultStoreItem](fmt.Errorf("Vault store is closed"))
	}
	if conn == "" {
		conn = VAULT_DEFAULT_NAME
	}
	if authName == "" {
		authName = VAULT_DEFAULT_NAME
	}
	path := fmt.Sprintf("%s/%s", conn, authName)
	c, ok := self.clients[path]
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
	if vault == nil {
		return fp.Left[*api.TLSConfig](fmt.Errorf("missing Vault connection"))
	}
	if vault.Tls == nil {
		return fp.Right(res)
	}
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

func newVaultClient(ctx context.Context, vault *al_proto.VaultConn, auth *al_proto.VaultAuth) (result fp.Either[*VaultStoreItem]) {
	if vault == nil || vault.Config == nil || auth == nil {
		return fp.Left[*VaultStoreItem](fmt.Errorf("missing Vault connection or authentication configuration"))
	}
	if auth.NoAuth && (auth.Approle != nil || auth.TokenHelper != nil) {
		return fp.Left[*VaultStoreItem](fmt.Errorf("no_auth cannot be combined with another authentication method"))
	}
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
	// api.NewClient may load VAULT_TOKEN from the environment.
	client.ClearToken()
	if auth.NoAuth {
		return fp.Right(&VaultStoreItem{Client: client})
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
		approleData, err := client.Logical().WriteWithContext(ctx, fmt.Sprintf("auth/approle/role/%s/secret-id", auth.Approle.Name), map[string]any{"num_uses": 1})
		if err != nil {
			return fp.Left[*VaultStoreItem](fmt.Errorf("could not create secret id for the approle: %w", SanitizeVaultError(err)))
		}
		if approleData == nil {
			return fp.Left[*VaultStoreItem](fmt.Errorf("missing AppRole SecretID response"))
		}
		accessor, _ := approleData.Data["secret_id_accessor"].(string)
		loggedIn := false
		// A failed or cancelled login must not leave a reusable SecretID behind.
		// Keep the bootstrap client/token separate from the issued role token.
		defer func() {
			if !loggedIn && accessor != "" {
				cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_, cleanupErr := client.Logical().WriteWithContext(cleanupCtx,
					fmt.Sprintf("auth/approle/role/%s/secret-id-accessor/destroy", auth.Approle.Name),
					map[string]any{"secret_id_accessor": accessor})
				if cleanupErr != nil {
					_, originalErr := result.Get()
					result = fp.Left[*VaultStoreItem](errors.Join(originalErr, fmt.Errorf("could not destroy unused AppRole SecretID")))
				}
			}
		}()
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
			return fp.Left[*VaultStoreItem](fmt.Errorf("could not auth using the approle: %w", SanitizeVaultError(err)))
		}
		loggedIn = true
	}
	return fp.Right(&VaultStoreItem{
		Client: client,
		Helper: helper,
		owned:  auth.Approle != nil,
	})
}

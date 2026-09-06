package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/tools/vault/injector/injector_proto"
)

type OidcFetcher struct {
	vault      *al.VaultStore
	httpClient *http.Client
}

func NewOidcFetcher(vault *al.VaultStore) *OidcFetcher {
	return &OidcFetcher{vault: vault}
}

var _ ResourceFetcher = (*OidcFetcher)(nil)

func (self *OidcFetcher) String() string {
	return "com.alwaldend.src.tools.vault.injector.OidcFetcher"
}

func (self *OidcFetcher) Get(ctx context.Context, r *injector_proto.Resource, d []*ResourceResult) (*ResourceResult, error) {
	config := r.GetOidc()
	if r.GetOidc() == nil {
		return nil, fmt.Errorf("missing OIDC config")
	}
	client, err := self.vault.Client(ctx, r.VaultConn, r.VaultAuth).Get()
	if err != nil {
		return nil, fmt.Errorf("could not create vault client: %w", al.SanitizeVaultError(err))
	}
	transport := self.httpClient
	if transport == nil {
		transport = client.Client.CloneConfig().HttpClient
	}
	httpClient := *transport
	// Never forward the bearer token or token request body through redirects.
	httpClient.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	data := map[string]any{
		"scope":         config.Scope,
		"response_type": "code",
		"client_id":     config.ClientId,
		"redirect_uri":  config.RedirectUri,
	}
	reqBody, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("could not marshal data: %w", al.SanitizeVaultError(err))
	}
	token := client.Client.Token()
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/v1/identity/oidc/provider/%s/authorize", client.Client.Address(), config.Name),
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create authorization request: %w", al.SanitizeVaultError(err))
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute the authorization request: %w", al.SanitizeVaultError(err))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read the body: %w", al.SanitizeVaultError(err))
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid authorization request response code: %d", resp.StatusCode)
	}
	oidc := &al.VaultOidc{}
	if err = json.Unmarshal(body, oidc); err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", al.SanitizeVaultError(err))
	}
	clientResp, err := client.Client.Logical().ReadWithContext(ctx, fmt.Sprintf("identity/oidc/client/%s", config.Name))
	if err != nil {
		return nil, fmt.Errorf("could not get oidc client info: %w", al.SanitizeVaultError(err))
	}
	if clientResp == nil {
		return nil, fmt.Errorf("missing oidc client info")
	}
	clientSecret, ok := clientResp.Data["client_secret"]
	if !ok {
		return nil, fmt.Errorf("client data missing client_secret")
	}
	tokenData := map[string]any{
		"code":          oidc.Code,
		"grant_type":    "authorization_code",
		"client_id":     config.ClientId,
		"client_secret": clientSecret,
		"redirect_uri":  config.RedirectUri,
	}
	tokenBody, err := json.Marshal(tokenData)
	if err != nil {
		return nil, fmt.Errorf("could not marshal token data: %w", al.SanitizeVaultError(err))
	}
	tokenReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/v1/identity/oidc/provider/%s/token", client.Client.Address(), config.Name),
		bytes.NewBuffer(tokenBody),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create token request: %w", al.SanitizeVaultError(err))
	}
	tokenReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	tokenResp, err := httpClient.Do(tokenReq)
	if err != nil {
		return nil, fmt.Errorf("could not execute the token request: %w", al.SanitizeVaultError(err))
	}
	defer tokenResp.Body.Close()
	tokenRespBody, err := io.ReadAll(tokenResp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read the token response body: %w", al.SanitizeVaultError(err))
	}
	if tokenResp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid token request response code: %d", tokenResp.StatusCode)
	}
	tokenRespData := &struct {
		IdToken string `json:"id_token"`
	}{}
	if err = json.Unmarshal(tokenRespBody, tokenRespData); err != nil {
		return nil, fmt.Errorf("could not unmarshal token response body: %w", al.SanitizeVaultError(err))
	}
	res := &ResourceResult{
		Name: r.Name,
		Data: map[string]any{
			"id_token": tokenRespData.IdToken,
		},
	}
	return res, nil
}

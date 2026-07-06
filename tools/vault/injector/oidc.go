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
	vault *al.VaultStore
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
		return nil, fmt.Errorf("could not create vault client: %w", err)
	}
	data := map[string]any{
		"scope":         config.Scope,
		"response_type": "code",
		"client_id":     config.ClientId,
		"redirect_uri":  config.RedirectUri,
	}
	reqBody, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("could not marshal data: %w", err)
	}
	token := client.Client.Token()
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/v1/identity/oidc/provider/%s/authorize", client.Client.Address(), config.Name),
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create authorization request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute the authorization request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read the body: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid authorization request response code: %s: %s", resp.Status, string(body))
	}
	oidc := &al.VaultOidc{}
	if err = json.Unmarshal(body, oidc); err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}
	clientResp, err := client.Client.Logical().Read(fmt.Sprintf("identity/oidc/client/%s", config.Name))
	if err != nil {
		return nil, fmt.Errorf("could not get oidc client info: %w", err)
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
		return nil, fmt.Errorf("could not marshal token data: %w", err)
	}
	tokenReq, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/v1/identity/oidc/provider/%s/token", client.Client.Address(), config.Name),
		bytes.NewBuffer(tokenBody),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create token request: %w", err)
	}
	tokenReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	tokenResp, err := http.DefaultClient.Do(tokenReq)
	if err != nil {
		return nil, fmt.Errorf("could not execute the token request: %w", err)
	}
	defer tokenReq.Body.Close()
	tokenRespBody, err := io.ReadAll(tokenResp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read the token response body: %w", err)
	}
	if tokenResp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid token request response code: %s: %s", tokenResp.Status, string(tokenRespBody))
	}
	tokenRespData := &struct {
		IdToken string `json:"id_token"`
	}{}
	if err = json.Unmarshal(tokenRespBody, tokenRespData); err != nil {
		return nil, fmt.Errorf("could not unmarshal token response body: %w", err)
	}
	res := &ResourceResult{
		Name: r.Name,
		Data: map[string]any{
			"id_token": tokenRespData.IdToken,
		},
	}
	return res, nil
}

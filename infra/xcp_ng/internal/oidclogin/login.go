package oidclogin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// authorizeURL checks every redirect destination before any Vault credential is
// attached. The configured discovery URL fixes the provider, not the redirect.
func authorizeURL(location string, vault, discovery, callback *url.URL) (*url.URL, map[string]string, error) {
	u, err := url.Parse(location)
	providerPath := strings.TrimSuffix(discovery.Path, "/.well-known/openid-configuration")
	const prefix = "/v1/identity/oidc/provider/"
	provider := strings.TrimPrefix(providerPath, prefix)
	if !strings.HasPrefix(providerPath, prefix) || provider == "" || strings.Contains(provider, "/") || discovery.Scheme != vault.Scheme || discovery.Host != vault.Host {
		return nil, nil, errors.New("OIDC discovery URL must identify a provider at VAULT_ADDR")
	}
	if err != nil || u.Scheme != vault.Scheme || u.Host != vault.Host || u.User != nil || u.Fragment != "" || u.Path != "/ui/vault/identity/oidc/provider/"+provider+"/authorize" {
		return nil, nil, errors.New("unexpected OIDC authorization destination")
	}
	query, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, nil, errors.New("invalid OIDC authorization query")
	}
	params := map[string]string{}
	for k, v := range query {
		if len(v) != 1 {
			return nil, nil, errors.New("duplicate OIDC authorization parameter")
		}
		params[k] = v[0]
	}
	if params["redirect_uri"] != callback.String() || params["state"] == "" || params["response_type"] != "code" {
		return nil, nil, errors.New("unexpected OIDC authorization parameters")
	}
	target := *vault
	target.Path = providerPath + "/authorize"
	target.RawQuery = ""
	target.Fragment = ""
	return &target, params, nil
}

func Bootstrap(xoClient, vaultClient *http.Client, xo, vault, discovery *url.URL, token string) error {
	start := *xo
	start.Path = "/signin/oidc"
	start.RawQuery = ""
	callback := *xo
	callback.Path = "/signin/oidc/callback"
	callback.RawQuery = ""
	response, err := xoClient.Get(start.String())
	if err != nil {
		return errors.New("starting OIDC login failed")
	}
	response.Body.Close()
	if response.StatusCode != http.StatusFound && response.StatusCode != http.StatusSeeOther {
		return errors.New("XO did not redirect to its OIDC provider")
	}
	target, params, err := authorizeURL(response.Header.Get("Location"), vault, discovery, &callback)
	if err != nil {
		return err
	}
	body, err := json.Marshal(params)
	if err != nil {
		return errors.New("encoding OIDC authorization failed")
	}
	req, err := http.NewRequest(http.MethodPost, target.String(), bytes.NewReader(body))
	if err != nil {
		return errors.New("constructing Vault authorization failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Vault-Token", token)
	response, err = vaultClient.Do(req)
	if err != nil {
		return errors.New("Vault OIDC authorization failed")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Vault OIDC authorization failed (HTTP %d)", response.StatusCode)
	}
	var grant struct {
		Code  string `json:"code"`
		State string `json:"state"`
	}
	if json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&grant) != nil || grant.Code == "" || grant.State != params["state"] {
		return errors.New("invalid Vault OIDC authorization result")
	}
	callback.RawQuery = url.Values{"code": {grant.Code}, "state": {grant.State}}.Encode()
	response, err = xoClient.Get(callback.String())
	if err != nil {
		return errors.New("XO OIDC callback failed")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusFound && response.StatusCode != http.StatusSeeOther {
		return fmt.Errorf("XO OIDC callback did not complete (HTTP %d)", response.StatusCode)
	}
	redirect, err := callback.Parse(response.Header.Get("Location"))
	if err != nil || redirect.Scheme != xo.Scheme || redirect.Host != xo.Host || redirect.Path != "/" || redirect.RawQuery != "" || redirect.Fragment != "" {
		return errors.New("XO OIDC callback returned an unexpected destination")
	}
	return nil
}

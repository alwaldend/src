package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"

	"git.alwaldend.com/alwaldend/src/infra/xcp_ng/internal/oidclogin"
)

type authorization struct {
	AdminGroup string `json:"admin_group"`
	UserGroup  string `json:"user_group"`
}

// XO's synchronized groups expose the plugin provider, but not an issuer.
// Reject duplicate OIDC names rather than selecting an arbitrary group.
func oidcGroups(r *rpc) (map[string]string, error) {
	var groups []struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Provider string `json:"provider"`
	}
	if err := r.call("group.getAll", map[string]any{}, &groups); err != nil {
		return nil, err
	}
	ids := map[string]string{}
	for _, group := range groups {
		if group.Provider != "oidc" {
			continue
		}
		if group.ID == "" || group.Name == "" || ids[group.Name] != "" {
			return nil, errors.New("invalid or duplicate synchronized OIDC group")
		}
		ids[group.Name] = group.ID
	}
	return ids, nil
}

func applyAuthorization(r *rpc, raw, configuration string) (resultErr error) {
	var desired authorization
	if json.Unmarshal([]byte(raw), &desired) != nil || desired.AdminGroup == "" || desired.UserGroup == "" || desired.AdminGroup == desired.UserGroup {
		return errors.New("invalid XO_OIDC_AUTHORIZATION groups")
	}
	var conf struct {
		DiscoveryURL string `json:"discoveryURL"`
	}
	if json.Unmarshal([]byte(configuration), &conf) != nil {
		return errors.New("invalid OIDC configuration")
	}
	xo, err := url.Parse(os.Getenv("XOA_URL"))
	if err != nil {
		return errors.New("invalid XOA_URL")
	}
	xo.Scheme = "https"
	xo.Path = "/"
	xo.RawQuery = ""
	vault, err := url.Parse(os.Getenv("VAULT_ADDR"))
	if err != nil || vault.Scheme != "https" || vault.Host == "" || vault.User != nil || (vault.Path != "" && vault.Path != "/") || vault.RawQuery != "" || vault.Fragment != "" {
		return errors.New("VAULT_ADDR must be an HTTPS origin")
	}
	discovery, err := url.Parse(conf.DiscoveryURL)
	if err != nil {
		return errors.New("invalid OIDC discovery URL")
	}
	token := os.Getenv("VAULT_TOKEN")
	if token == "" {
		return errors.New("VAULT_TOKEN is required for OIDC group bootstrap")
	}
	vaultTLS := &tls.Config{MinVersion: tls.VersionTLS12}
	if path := os.Getenv("VAULT_CACERT"); path != "" {
		pem, err := os.ReadFile(path)
		if err != nil {
			return errors.New("reading VAULT_CACERT failed")
		}
		roots, err := x509.SystemCertPool()
		if err != nil {
			roots = x509.NewCertPool()
		}
		if !roots.AppendCertsFromPEM(pem) {
			return errors.New("invalid VAULT_CACERT")
		}
		vaultTLS.RootCAs = roots
	}
	noRedirect := func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	jar, _ := cookiejar.New(nil)
	xoTransport := http.DefaultTransport.(*http.Transport).Clone()
	xoTransport.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: os.Getenv("XOA_INSECURE") == "true"}
	vaultTransport := http.DefaultTransport.(*http.Transport).Clone()
	vaultTransport.TLSClientConfig = vaultTLS
	defer xoTransport.CloseIdleConnections()
	defer vaultTransport.CloseIdleConnections()
	xoClient := &http.Client{Transport: xoTransport, Jar: jar, Timeout: 45 * time.Second, CheckRedirect: noRedirect}
	defer func() {
		for _, cookie := range jar.Cookies(xo) {
			if cookie.Name == "token" && cookie.Value != "" {
				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()
				resultErr = errors.Join(resultErr, oidclogin.Revoke(ctx, xo, cookie.Value))
				return
			}
		}
	}()
	vaultClient := &http.Client{Transport: vaultTransport, Timeout: 45 * time.Second, CheckRedirect: noRedirect}
	if err := oidclogin.Bootstrap(xoClient, vaultClient, xo, vault, discovery, token); err != nil {
		return err
	}
	groups, err := oidcGroups(r)
	if err != nil {
		return err
	}
	if groups[desired.AdminGroup] == "" || groups[desired.UserGroup] == "" {
		return errors.New("expected synchronized OIDC groups missing after login; check bootstrap Vault identity group membership")
	}
	return nil
}

// externalOIDCIdentities implements Terraform external's string-map protocol.
// Authentication stays in the process environment, outside the query/state.
func externalOIDCIdentities(in io.Reader, out io.Writer) error {
	var query map[string]string
	if json.NewDecoder(io.LimitReader(in, 1<<16)).Decode(&query) != nil || query["url"] == "" || query["issuer"] == "" || (query["insecure"] != "true" && query["insecure"] != "false") {
		return errors.New("invalid OIDC identity query")
	}
	if err := os.Setenv("XOA_INSECURE", query["insecure"]); err != nil {
		return errors.New("setting XO TLS mode failed")
	}
	r, err := connect(query["url"], os.Getenv("XOA_TOKEN"))
	if err != nil {
		return err
	}
	defer r.conn.Close()
	return writeOIDCIdentities(r, "oidc:"+query["issuer"], query["admin_group"], out)
}

// Bind the immutable Vault OIDC subject, not a mutable username or inherited group.
func writeOIDCIdentities(r *rpc, provider, adminGroup string, out io.Writer) error {
	var users []struct {
		ID            string `json:"id"`
		AuthProviders map[string]struct {
			ID string `json:"id"`
		} `json:"authProviders"`
	}
	if err := r.call("user.getAll", map[string]any{}, &users); err != nil {
		return err
	}
	ids := map[string]string{}
	for _, user := range users {
		subject := user.AuthProviders[provider].ID
		if subject == "" {
			continue
		}
		if user.ID == "" || ids[subject] != "" {
			return errors.New("invalid or duplicate synchronized OIDC identity")
		}
		ids[subject] = user.ID
	}
	data, err := json.Marshal(ids)
	if err != nil {
		return err
	}
	groups, err := oidcGroups(r)
	if err != nil {
		return err
	}
	groupData, err := json.Marshal(groups)
	if err != nil {
		return err
	}
	aclIDs := map[string]string{}
	if adminID := groups[adminGroup]; adminGroup != "" && adminID != "" {
		var acls []struct {
			ID      string `json:"id"`
			Subject string `json:"subject"`
			Object  string `json:"object"`
			Action  string `json:"action"`
		}
		if err := r.call("acl.get", map[string]any{}, &acls); err != nil {
			return err
		}
		for _, entry := range acls {
			if entry.Subject != adminID || entry.Action != "admin" {
				continue
			}
			if entry.ID == "" || entry.Object == "" || aclIDs[entry.Object] != "" {
				return errors.New("invalid or duplicate OIDC administrator ACL")
			}
			aclIDs[entry.Object] = entry.ID
		}
	}
	aclData, err := json.Marshal(aclIDs)
	if err != nil {
		return err
	}
	return json.NewEncoder(out).Encode(map[string]string{"users": string(data), "groups": string(groupData), "pool_admin_acl_ids": string(aclData)})
}

package main

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
	"os/signal"
	"strings"
	"syscall"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al_plugin"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
	"git.alwaldend.com/alwaldend/src/tools/go/pkg/sensitiveerr"
	"git.alwaldend.com/alwaldend/src/tools/vault/pve_login/pve_login_proto"
	"github.com/google/uuid"
)

// oidcURL decodes the authorization URL returned by the Proxmox realm.
type oidcURL struct {
	// Data holds the envelope payload returned by the API.
	Data string `json:"data"`
}

// pveTicket decodes the authenticated session used to issue and remove the invocation token.
type pveTicket struct {
	// Data holds the envelope payload returned by the API.
	Data struct {
		// CSRFPreventionToken authenticates state-changing requests made with the session ticket.
		CSRFPreventionToken string `json:"CSRFPreventionToken"`
		// Username identifies the Proxmox account that owns the token.
		Username string `json:"username"`
		// Ticket supplies the PVEAuthCookie used by the token management endpoints.
		Ticket string `json:"ticket"`
	} `json:"data"`
}

// pveToken decodes the identifier and secret returned by token creation.
type pveToken struct {
	// Data holds the envelope payload returned by the API.
	Data struct {
		// TokenID identifies the token together with its owning account.
		TokenID string `json:"full-tokenid"`
		// TokenSecret contains the credential returned only at creation.
		TokenSecret string `json:"value"`
	} `json:"data"`
}

// readResponse reads one complete API reply within the one MiB limit.
func readResponse(input io.Reader) ([]byte, error) {
	const limit = 1 << 20
	data, err := io.ReadAll(io.LimitReader(input, limit+1))
	if err != nil {
		return nil, sensitiveerr.WithContext("read Proxmox response", err)
	}
	if len(data) > limit {
		return nil, errors.New("Proxmox response exceeds one MiB")
	}
	return data, nil
}

// createOIDCRequest obtains the authorization URL for the configured identity provider.
func createOIDCRequest(ctx context.Context, config *pve_login_proto.Config) (*url.URL, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api2/json/access/openid/auth-url", config.PveBaseUrl),
		strings.NewReader("{}"),
	)
	if err != nil {
		return nil, sensitiveerr.WithContext("could not create request", err)
	}
	query := req.URL.Query()
	query.Set("realm", config.PveRealm)
	query.Set("redirect-url", config.PveRedirectUrl)
	req.URL.RawQuery = query.Encode()
	req.Header.Add("Content-Type", "application/json")
	resp, err := pveHTTPClient.Do(req)
	if err != nil {
		return nil, sensitiveerr.WithContext("could not execute request", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}
	body, err := readResponse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}
	data := &oidcURL{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return nil, sensitiveerr.WithContext("could not unmarshal response", err)
	}
	res, err := url.Parse(data.Data)
	if err != nil {
		return nil, sensitiveerr.WithContext("could not parse authorization URL", err)
	}
	return res, nil
}

// createProxmoxTicket exchanges a Vault OIDC authorization for a Proxmox login ticket.
func createProxmoxTicket(ctx context.Context, config *pve_login_proto.Config, oidc *al.VaultOidc) (*pveTicket, error) {
	body, err := json.Marshal(map[string]string{"state": oidc.State, "code": oidc.Code, "redirect-url": config.PveRedirectUrl})
	if err != nil {
		return nil, sensitiveerr.WithContext("could not encode ticket request", err)
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api2/json/access/openid/login", config.PveBaseUrl),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, sensitiveerr.WithContext("could not create request", err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := pveHTTPClient.Do(req)
	if err != nil {
		return nil, sensitiveerr.WithContext("could not execute the request", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid response code: %d", resp.StatusCode)
	}
	body, err = readResponse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}
	res := &pveTicket{}
	if err := json.Unmarshal(body, res); err != nil {
		return nil, sensitiveerr.WithContext("could not unmarshal response body", err)
	}
	if res.Data.Username == "" || res.Data.Ticket == "" || res.Data.CSRFPreventionToken == "" {
		return nil, fmt.Errorf("incomplete Proxmox ticket response")
	}
	return res, nil
}

// createProxmoxToken creates a named API token with a one-hour expiration.
func createProxmoxToken(ctx context.Context, config *pve_login_proto.Config, pveTicket *pveTicket, name string) (*pveToken, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf(
			"%s/api2/json/access/users/%s/token/%s",
			config.PveBaseUrl,
			url.PathEscape(pveTicket.Data.Username),
			name,
		),
		strings.NewReader("{}"),
	)
	if err != nil {
		return nil, sensitiveerr.WithContext("could not create request", err)
	}
	query := req.URL.Query()
	query.Add("comment", "Automatically created by //tools/vault/pve_login")
	query.Add("expire", fmt.Sprintf("%d", time.Now().Add(time.Hour).Unix()))
	query.Add("privsep", "0")
	req.URL.RawQuery = query.Encode()
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("CSRFPreventionToken", pveTicket.Data.CSRFPreventionToken)
	req.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: pveTicket.Data.Ticket})
	resp, err := pveHTTPClient.Do(req)
	if err != nil {
		return nil, sensitiveerr.WithContext("could not execute the request", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid response code: %d", resp.StatusCode)
	}
	body, err := readResponse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}
	res := &pveToken{}
	if err := json.Unmarshal(body, res); err != nil {
		return nil, sensitiveerr.WithContext("could not unmarshal response body", err)
	}
	if res.Data.TokenID != pveTicket.Data.Username+"!"+name || res.Data.TokenSecret == "" {
		return nil, fmt.Errorf("incomplete or mismatched Proxmox token response")
	}
	return res, nil
}

// pveHTTPClient bounds request duration and prevents automatic redirects during authentication.
var pveHTTPClient = &http.Client{
	Timeout:       30 * time.Second,
	CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
}

// deleteProxmoxToken uses the retained session to revoke an invocation's token.
func deleteProxmoxToken(ctx context.Context, config *pve_login_proto.Config, ticket *pveTicket, name string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete,
		fmt.Sprintf("%s/api2/json/access/users/%s/token/%s", config.PveBaseUrl, url.PathEscape(ticket.Data.Username), url.PathEscape(name)), nil)
	if err != nil {
		return sensitiveerr.WithContext("could not create token deletion request", err)
	}
	req.Header.Set("CSRFPreventionToken", ticket.Data.CSRFPreventionToken)
	req.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: ticket.Data.Ticket})
	resp, err := pveHTTPClient.Do(req)
	if err != nil {
		return sensitiveerr.WithContext("could not delete Proxmox token", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Proxmox token deletion returned HTTP %d", resp.StatusCode)
	}
	return nil
}

// plugin serves AL login requests and owns cleanup of invocation-scoped credentials.
type plugin struct {
	// lc owns the cleanup actions run during shutdown.
	lc lifecycle.Manager
}

// Start starts the AL plugin's lifecycle manager.
func (p *plugin) Start(ctx context.Context) error {
	if err := p.lc.Start(ctx); err != nil {
		return fmt.Errorf("start Proxmox login plugin: %w", err)
	}
	return nil
}

// Stop releases resources registered by plugin startup in reverse order.
func (p *plugin) Stop(ctx context.Context) error {
	if err := p.lc.Stop(ctx); err != nil {
		return fmt.Errorf("stop Proxmox login plugin: %w", err)
	}
	return nil
}

// PluginStart authenticates the invocation, registers cleanup, and returns credentials to AL.
func (p *plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (_ *al_proto.PluginStartResponse, retErr error) {
	defer func() {
		if retErr != nil {
			cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			retErr = errors.Join(retErr, p.Stop(cleanupCtx))
		}
	}()
	if req.GetPlugin() == nil {
		return nil, errors.New("missing plugin configuration")
	}
	config := &pve_login_proto.Config{}
	if _, err := al.FromPbJsonToPb(req.GetPlugin().GetData(), config).Get(); err != nil {
		return nil, sensitiveerr.WithContext("could not parse plugin data", err)
	}
	vault := al.NewVault(req.Config)
	if err := p.lc.AddState(lifecycle.StateStarted, vault); err != nil {
		return nil, fmt.Errorf("could not register Vault cleanup: %w", err)
	}
	client, err := vault.Client(ctx, config.VaultConn, config.VaultAuth).Get()
	if err != nil {
		return nil, sensitiveerr.WithContext("could not create the vault client", err)
	}
	oidcUrl, err := createOIDCRequest(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("could not create OIDC request: %w", err)
	}
	vaultOidc, err := vault.OidcLoginContext(ctx, client.Client, oidcUrl)
	if err != nil {
		return nil, sensitiveerr.WithContext("could not login to OIDC provider", err)
	}
	pveTicket, err := createProxmoxTicket(ctx, config, vaultOidc)
	if err != nil {
		return nil, fmt.Errorf("could not create PVE ticket: %w", err)
	}
	name := fmt.Sprintf("tools-vault-pve-login-%s", uuid.New().String())
	// Register by the locally generated name before issuing the creation request:
	// even a lost or malformed response may have created the token.
	if err := p.lc.AddState(lifecycle.StateStarted, lifecycle.StoppableFunc(func(ctx context.Context) error {
		if err := deleteProxmoxToken(ctx, config, pveTicket, name); err != nil {
			return fmt.Errorf("remove invocation token: %w", err)
		}
		return nil
	})); err != nil {
		return nil, fmt.Errorf("could not register token cleanup: %w", err)
	}
	pveToken, err := createProxmoxToken(ctx, config, pveTicket, name)
	if err != nil {
		return nil, fmt.Errorf("could not create proxmox token: %w", err)
	}
	res := &al_proto.ResourceResult{
		Env: map[string]string{
			"PM_API_TOKEN_ID":     pveToken.Data.TokenID,
			"TOKEN_ID":            pveToken.Data.TokenID,
			"PM_API_TOKEN_SECRET": pveToken.Data.TokenSecret,
			"SECRET":              pveToken.Data.TokenSecret,
			"PM_API_URL":          fmt.Sprintf("%s/api2/json", config.PveBaseUrl),
			"API_URL":             fmt.Sprintf("%s/api2/json", config.PveBaseUrl),
		},
	}
	return &al_proto.PluginStartResponse{Result: res}, nil
}

// run serves plugin requests until cancellation or a lifecycle failure.
func run(ctx *al.CmdCtx) error {
	var lc lifecycle.Manager
	provider := &plugin{}
	server := al_plugin.NewPluginServer(ctx, provider)
	lc.Add(provider, server)
	if err := lc.Run(ctx.Ctx, 10*time.Second); err != nil {
		return fmt.Errorf("could not run: %w", err)
	}
	return nil
}

// main binds plugin shutdown to process signals and reports terminal failures.
func main() {
	shutdownCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	ctx := al.NewCmdCtx(shutdownCtx, "com.alwaldend.src.tools.vault.pve_login ")
	if err := run(ctx); err != nil {
		ctx.Logger.Printf("failed: %s", err)
		os.Exit(1)
	}
}

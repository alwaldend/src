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
	"git.alwaldend.com/alwaldend/src/tools/vault/pve_login/pve_login_proto"
	"github.com/google/uuid"
)

type oidcUrl struct {
	Data string `json:"data"`
}

type pveTicket struct {
	Data struct {
		CSRFPreventionToken string `json:"CSRFPreventionToken"`
		Username            string `json:"username"`
		Ticket              string `json:"ticket"`
	} `json:"data"`
}

type pveToken struct {
	Data struct {
		TokenId     string `json:"full-tokenid"`
		TokenSecret string `json:"value"`
	} `json:"data"`
}

func createOIDCRequest(ctx context.Context, config *pve_login_proto.Config) (*url.URL, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api2/json/access/openid/auth-url", config.PveBaseUrl),
		strings.NewReader("{}"),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create request")
	}
	query := req.URL.Query()
	query.Set("realm", config.PveRealm)
	query.Set("redirect-url", config.PveRedirectUrl)
	req.URL.RawQuery = query.Encode()
	req.Header.Add("Content-Type", "application/json")
	resp, err := pveHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute request: %w", safeHTTPError(err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}
	data := &oidcUrl{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response")
	}
	res, err := url.Parse(data.Data)
	if err != nil {
		return nil, fmt.Errorf("could not parse authorization URL")
	}
	return res, nil
}

func createProxmoxTicket(ctx context.Context, config *pve_login_proto.Config, oidc *al.VaultOidc) (*pveTicket, error) {
	body, err := json.Marshal(map[string]string{"state": oidc.State, "code": oidc.Code, "redirect-url": config.PveRedirectUrl})
	if err != nil {
		return nil, fmt.Errorf("could not encode ticket request")
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api2/json/access/openid/login", config.PveBaseUrl),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create request")
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := pveHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute the request: %w", safeHTTPError(err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid response code: %d", resp.StatusCode)
	}
	body, err = io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}
	res := &pveTicket{}
	if err := json.Unmarshal(body, res); err != nil {
		return nil, fmt.Errorf("could not unmarshal response body")
	}
	if res.Data.Username == "" || res.Data.Ticket == "" || res.Data.CSRFPreventionToken == "" {
		return nil, fmt.Errorf("incomplete Proxmox ticket response")
	}
	return res, nil
}

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
		return nil, fmt.Errorf("could not create request")
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
		return nil, fmt.Errorf("could not execute the request: %w", safeHTTPError(err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid response code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}
	res := &pveToken{}
	if err := json.Unmarshal(body, res); err != nil {
		return nil, fmt.Errorf("could not unmarshal response body")
	}
	if res.Data.TokenId != pveTicket.Data.Username+"!"+name || res.Data.TokenSecret == "" {
		return nil, fmt.Errorf("incomplete or mismatched Proxmox token response")
	}
	return res, nil
}

var pveHTTPClient = &http.Client{
	Timeout:       30 * time.Second,
	CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
}

func safeHTTPError(err error) error {
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return urlErr.Err
	}
	return err
}

func deleteProxmoxToken(ctx context.Context, config *pve_login_proto.Config, ticket *pveTicket, name string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete,
		fmt.Sprintf("%s/api2/json/access/users/%s/token/%s", config.PveBaseUrl, url.PathEscape(ticket.Data.Username), url.PathEscape(name)), nil)
	if err != nil {
		return fmt.Errorf("could not create token deletion request")
	}
	req.Header.Set("CSRFPreventionToken", ticket.Data.CSRFPreventionToken)
	req.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: ticket.Data.Ticket})
	resp, err := pveHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not delete Proxmox token: %w", safeHTTPError(err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Proxmox token deletion returned HTTP %d", resp.StatusCode)
	}
	return nil
}

type Plugin struct{ lc lifecycle.Manager }

func (self *Plugin) Start(ctx context.Context) error { return self.lc.Start(ctx) }
func (self *Plugin) Stop(ctx context.Context) error  { return self.lc.Stop(ctx) }

func (self *Plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (_ *al_proto.PluginStartResponse, retErr error) {
	defer func() {
		if retErr != nil {
			cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			retErr = errors.Join(retErr, self.Stop(cleanupCtx))
		}
	}()
	config := &pve_login_proto.Config{}
	if _, err := al.FromPbJsonToPb(req.Plugin.Data, config).Get(); err != nil {
		return nil, fmt.Errorf("could not parse plugin data")
	}
	vault := al.NewVault(req.Config)
	if err := self.lc.AddState(lifecycle.StateStarted, vault); err != nil {
		return nil, fmt.Errorf("could not register Vault cleanup: %w", err)
	}
	client, err := vault.Client(ctx, config.VaultConn, config.VaultAuth).Get()
	if err != nil {
		return nil, fmt.Errorf("could not create the vault client: %w", err)
	}
	oidcUrl, err := createOIDCRequest(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("could not create OIDC request: %w", err)
	}
	vaultOidc, err := vault.OidcLoginContext(ctx, client.Client, oidcUrl)
	if err != nil {
		return nil, fmt.Errorf("could not login to OIDC provider: %w", err)
	}
	pveTicket, err := createProxmoxTicket(ctx, config, vaultOidc)
	if err != nil {
		return nil, fmt.Errorf("could not create PVE ticket: %w", err)
	}
	name := fmt.Sprintf("tools-vault-pve-login-%s", uuid.New().String())
	// Register by the locally generated name before issuing the creation request:
	// even a lost or malformed response may have created the token.
	if err := self.lc.AddState(lifecycle.StateStarted, lifecycle.StoppableFunc(func(ctx context.Context) error {
		return deleteProxmoxToken(ctx, config, pveTicket, name)
	})); err != nil {
		return nil, fmt.Errorf("could not register token cleanup: %w", err)
	}
	pveToken, err := createProxmoxToken(ctx, config, pveTicket, name)
	if err != nil {
		return nil, fmt.Errorf("could not create proxmox token: %w", err)
	}
	res := &al_proto.PluginStartResponse{
		Env: map[string]string{
			"PM_API_TOKEN_ID":     pveToken.Data.TokenId,
			"TOKEN_ID":            pveToken.Data.TokenId,
			"PM_API_TOKEN_SECRET": pveToken.Data.TokenSecret,
			"SECRET":              pveToken.Data.TokenSecret,
			"PM_API_URL":          fmt.Sprintf("%s/api2/json", config.PveBaseUrl),
			"API_URL":             fmt.Sprintf("%s/api2/json", config.PveBaseUrl),
		},
	}
	return res, nil
}

func run(ctx *al.CmdCtx) error {
	var lc lifecycle.Manager
	plugin := &Plugin{}
	server := al_plugin.NewPluginServer(ctx, plugin)
	lc.Add(plugin, server)
	if err := lc.Run(ctx.Ctx, time.Second*10); err != nil {
		return fmt.Errorf("could not run: %w", err)
	}
	return nil
}

func main() {
	shutdownCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	ctx := al.NewCmdCtx(shutdownCtx, "com.alwaldend.src.tools.vault.pve_login ")
	if err := run(ctx); err != nil {
		ctx.Logger.Printf("failed: %s", err)
		os.Exit(1)
	}
}

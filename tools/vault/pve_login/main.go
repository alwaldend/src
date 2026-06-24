package main

import (
	"context"
	"encoding/json"
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

func createOIDCRequest(config *pve_login_proto.Config) (*url.URL, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/api2/json/access/openid/auth-url?realm=%s&redirect-url=%s", config.PveBaseUrl, config.PveRealm, config.PveRedirectUrl),
		strings.NewReader("{}"),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}
	data := &oidcUrl{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response: %w", err)
	}
	res, err := url.Parse(data.Data)
	if err != nil {
		return nil, fmt.Errorf("could not parse the url: %w", err)
	}
	return res, nil
}

func createProxmoxTicket(config *pve_login_proto.Config, oidc *al.VaultOidc) (*pveTicket, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/api2/json/access/openid/login", config.PveBaseUrl),
		strings.NewReader("{}"),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	query := req.URL.Query()
	query.Add("state", oidc.State)
	query.Add("code", oidc.Code)
	query.Add("redirect-url", config.PveRedirectUrl)
	req.URL.RawQuery = query.Encode()
	req.Header.Add("Content-Type", "application/json")
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
		return nil, fmt.Errorf("could not read response body: %w", err)
	}
	res := &pveTicket{}
	if err := json.Unmarshal(body, res); err != nil {
		return nil, fmt.Errorf("could not unmarshel response body: %w", err)
	}
	return res, nil
}

func createProxmoxToken(config *pve_login_proto.Config, pveTicket *pveTicket) (*pveToken, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"%s/api2/json/access/users/%s/token/%s",
			config.PveBaseUrl,
			pveTicket.Data.Username,
			fmt.Sprintf("tools-vault-pve-login-%s", uuid.New().String()),
		),
		strings.NewReader("{}"),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	query := req.URL.Query()
	query.Add("comment", "Automatically created by //tools/vault/pve_login")
	query.Add("expire", fmt.Sprintf("%d", time.Now().Add(time.Hour).Unix()))
	query.Add("privsep", "0")
	req.URL.RawQuery = query.Encode()
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("CSRFPreventionToken", pveTicket.Data.CSRFPreventionToken)
	req.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: pveTicket.Data.Ticket})
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
		return nil, fmt.Errorf("could not read response body: %w", err)
	}
	res := &pveToken{}
	if err := json.Unmarshal(body, res); err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}
	return res, nil
}

type Plugin struct {
}

func (self *Plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (*al_proto.PluginStartResponse, error) {
	config := &pve_login_proto.Config{}
	if _, err := al.FromPbJsonToPb(req.Plugin.Data, config).Get(); err != nil {
		return nil, fmt.Errorf("could not parse plugin data: %w", err)
	}
	vault := al.NewVault(req.Config)
	client, err := vault.Client(ctx, config.VaultConn, config.VaultAuth).Get()
	if err != nil {
		return nil, fmt.Errorf("could not create the vault client: %w", err)
	}
	oidcUrl, err := createOIDCRequest(config)
	if err != nil {
		return nil, fmt.Errorf("could not create OIDC request: %w", err)
	}
	vaultOidc, err := vault.OidcLogin(client.Client, oidcUrl)
	if err != nil {
		return nil, fmt.Errorf("could not login to OIDC provider: %w", err)
	}
	pveTicket, err := createProxmoxTicket(config, vaultOidc)
	if err != nil {
		return nil, fmt.Errorf("could not create PVE ticket: %w", err)
	}
	pveToken, err := createProxmoxToken(config, pveTicket)
	if err != nil {
		return nil, fmt.Errorf("could not create proxmox token: %w", err)
	}
	res := &al_proto.PluginStartResponse{
		Env: map[string]string{
			"PM_API_TOKEN_ID":     pveToken.Data.TokenId,
			"PM_API_TOKEN_SECRET": pveToken.Data.TokenSecret,
			"PM_API_URL":          fmt.Sprintf("%s/api2/json", config.PveBaseUrl),
		},
	}
	return res, nil
}

func run(ctx *al.CmdCtx) error {
	var lc lifecycle.Manager
	server := al_plugin.NewPluginServer(ctx, &Plugin{})
	lc.Add(server)
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

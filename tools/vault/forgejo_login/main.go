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
	"regexp"
	"strings"
	"syscall"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al_plugin"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
	"git.alwaldend.com/alwaldend/src/tools/vault/forgejo_login/forgejo_login_proto"
	"github.com/google/uuid"
)

type vaultOidc struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

func createOIDCRequest(config *forgejo_login_proto.Config) (*url.URL, error) {
	redirectErr := errors.New("redirect")
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return redirectErr
		},
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/user/oauth2/vault", config.ForgejoUrl),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil && !errors.Is(err, redirectErr) {
		return nil, fmt.Errorf("could not execute the request: %w", err)
	}
	defer resp.Body.Close()
	if err == nil {
		return nil, fmt.Errorf("did not hit a redirect for some reason, status: %s", resp.Status)
	}
	res, err := resp.Location()
	if err != nil {
		return nil, fmt.Errorf("could not extract location: %w", err)
	}
	return res, nil
}

func loginToOIDCProvider(oidcUrl *url.URL, vaultToken string) (*vaultOidc, error) {
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
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", vaultToken))
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
	res := &vaultOidc{}
	if err = json.Unmarshal(body, res); err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}
	return res, nil
}

func createProxmoxTicket(config *forgejo_login_proto.Config, oidc *vaultOidc) (*pveTicket, error) {
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

func createForgejoToken(config *forgejo_login_proto.Config, pveTicket *pveTicket) (*pveToken, error) {
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
	query.Add("comment", "Automatically created by //tools/vault/forgejo_login")
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
	config := &forgejo_login_proto.Config{}
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
	vaultOidc, err := loginToOIDCProvider(oidcUrl, client.Client.Token())
	if err != nil {
		return nil, fmt.Errorf("could not login to OIDC provider: %w", err)
	}
	pveTicket, err := createProxmoxTicket(config, vaultOidc)
	if err != nil {
		return nil, fmt.Errorf("could not create PVE ticket: %w", err)
	}
	pveToken, err := createForgejoToken(config, pveTicket)
	if err != nil {
		return nil, fmt.Errorf("could not create proxmox token: %w", err)
	}
	res := &al_proto.PluginStartResponse{
		Env: map[string]string{
			"FORGEJO_API_TOKEN": pveToken,
			"FORGEJO_HOST":      config.ForgejoUrl,
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
	ctx := al.NewCmdCtx(shutdownCtx, "com.alwaldend.src.tools.vault.forgejo_login ")
	if err := run(ctx); err != nil {
		ctx.Logger.Printf("failed: %s", err)
		os.Exit(1)
	}
}

package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	reqUrl "net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al_plugin"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
	"git.alwaldend.com/alwaldend/src/tools/vault/harbor_login/harbor_login_proto"
)

type login struct {
	ctx      *al.CmdCtx
	config   *harbor_login_proto.Config
	client   *http.Client
	redirect *reqUrl.URL
	cookies  []*http.Cookie
}

func newLogin(ctx *al.CmdCtx, config *harbor_login_proto.Config) (*login, error) {
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return nil, fmt.Errorf("could not create cookiejar: %w", err)
	}
	res := &login{
		ctx:    ctx,
		config: config,
	}
	client := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			loc, err := req.Response.Location()
			if err != nil {
				return fmt.Errorf("response is missing the location for some reason: %w", err)
			}
			res.redirect = loc
			res.cookies = append(res.cookies, req.Response.Cookies()...)
			return nil
		},
	}
	res.client = client
	return res, nil
}

func (self *login) createOIDCRequest() (*reqUrl.URL, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/c/oidc/login", self.config.HarborUrl),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	resp, err := self.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute the request: %w", err)
	}
	defer resp.Body.Close()
	if self.redirect == nil {
		return nil, fmt.Errorf("did not hit redirect for some reason: %w", err)
	}
	return self.redirect, nil
}

func (self *login) authorizeSession(oidc *al.VaultOidc) (string, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/c/oidc/callback", self.config.HarborUrl),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("could not create request: %w", err)
	}
	query := req.URL.Query()
	query.Add("state", oidc.State)
	query.Add("code", oidc.Code)
	req.URL.RawQuery = query.Encode()
	resp, err := self.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not execute the request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("invalid response code: %s", resp.Status)
	}
	res := ""
	for _, cookie := range self.cookies {
		if cookie.Name == "sid" {
			res = cookie.Value
		}
	}
	if res == "" {
		return "", fmt.Errorf("mising sid cookie")
	}
	return res, nil
}

type Plugin struct {
	ctx *al.CmdCtx
}

func (self *Plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (*al_proto.PluginStartResponse, error) {
	self.ctx.Logger.Printf("init")
	config := &harbor_login_proto.Config{}
	if _, err := al.FromPbJsonToPb(req.Plugin.Data, config).Get(); err != nil {
		return nil, fmt.Errorf("could not parse plugin data: %w", err)
	}
	vault := al.NewVault(req.Config)
	client, err := vault.Client(ctx, config.VaultConn, config.VaultAuth).Get()
	if err != nil {
		return nil, fmt.Errorf("could not create the vault client: %w", err)
	}
	self.ctx.Logger.Printf("creating OIDC request")
	login, err := newLogin(self.ctx, config)
	if err != nil {
		return nil, fmt.Errorf("could not create login: %w", err)
	}
	oidcUrl, err := login.createOIDCRequest()
	if err != nil {
		return nil, fmt.Errorf("could not create OIDC request: %w", err)
	}
	self.ctx.Logger.Printf("authorizing the OIDC request")
	vaultOidc, err := vault.OidcLogin(client.Client, oidcUrl)
	if err != nil {
		return nil, fmt.Errorf("could not login to OIDC provider: %w", err)
	}
	self.ctx.Logger.Printf("authorizing the session")
	sessionId, err := login.authorizeSession(vaultOidc)
	if err != nil {
		return nil, fmt.Errorf("could not authorize the session using OIDC code: %w", err)
	}
	res := &al_proto.PluginStartResponse{
		Env: map[string]string{
			"HARBOR_SESSION_ID": sessionId,
			"HARBOR_URL":        config.HarborUrl,
		},
	}
	return res, nil
}

func run(ctx *al.CmdCtx) error {
	var lc lifecycle.Manager
	plugin := &Plugin{ctx: ctx}
	server := al_plugin.NewPluginServer(ctx, plugin)
	lc.Add(server)
	if err := lc.Run(ctx.Ctx, time.Second*10); err != nil {
		return fmt.Errorf("could not run: %w", err)
	}
	return nil
}

func main() {
	shutdownCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	ctx := al.NewCmdCtx(shutdownCtx, "com.alwaldend.src.tools.vault.harbor_login ")
	if err := run(ctx); err != nil {
		ctx.Logger.Printf("failed: %s", err)
		os.Exit(1)
	}
}

package main

import (
	"context"
	"errors"
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
	session  string
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
		Jar:     jar,
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			loc, err := req.Response.Location()
			if err != nil {
				return fmt.Errorf("invalid redirect location")
			}
			res.redirect = loc
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			if req.URL.Scheme != via[0].URL.Scheme || req.URL.Host != via[0].URL.Host {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
	res.client = client
	return res, nil
}

func (self *login) createOIDCRequest(ctx context.Context) (*reqUrl.URL, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/c/oidc/login", self.config.HarborUrl),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create request")
	}
	resp, err := self.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute the request: %w", safeHTTPError(err))
	}
	defer resp.Body.Close()
	if self.redirect == nil {
		return nil, fmt.Errorf("OIDC authorization redirect missing")
	}
	return self.redirect, nil
}

func (self *login) authorizeSession(ctx context.Context, oidc *al.VaultOidc) (string, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/c/oidc/callback", self.config.HarborUrl),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("could not create request")
	}
	query := req.URL.Query()
	query.Add("state", oidc.State)
	query.Add("code", oidc.Code)
	req.URL.RawQuery = query.Encode()
	resp, err := self.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not execute the request: %w", safeHTTPError(err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("invalid response code: %d", resp.StatusCode)
	}
	res := ""
	for _, cookie := range self.client.Jar.Cookies(req.URL) {
		if cookie.Name == "sid" {
			res = cookie.Value
		}
	}
	if res == "" {
		return "", fmt.Errorf("mising sid cookie")
	}
	self.session = res
	return res, nil
}

type Plugin struct {
	ctx *al.CmdCtx
	lc  lifecycle.Manager
}

func safeHTTPError(err error) error {
	var urlErr *reqUrl.Error
	if errors.As(err, &urlErr) {
		return urlErr.Err
	}
	return err
}

// Logout destroys this Harbor session without following the provider logout
// redirect, which could terminate an unrelated shared identity-provider session.
func (self *login) logout(ctx context.Context) error {
	base, err := reqUrl.Parse(self.config.HarborUrl)
	if err != nil {
		return fmt.Errorf("invalid Harbor URL")
	}
	session := self.session
	if session == "" {
		for _, cookie := range self.client.Jar.Cookies(base) {
			if cookie.Name == "sid" {
				session = cookie.Value
			}
		}
	}
	if session == "" {
		return nil
	}
	client := *self.client
	client.Jar = nil // Keep the original SID to verify server-side invalidation.
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, self.config.HarborUrl+"/c/oidc/logout", nil)
	if err != nil {
		return fmt.Errorf("could not create Harbor logout request")
	}
	req.AddCookie(&http.Cookie{Name: "sid", Value: session})
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Harbor logout failed: %w", safeHTTPError(err))
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		return fmt.Errorf("Harbor logout returned HTTP %d", resp.StatusCode)
	}
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, self.config.HarborUrl+"/api/v2.0/users/current", nil)
	if err != nil {
		return fmt.Errorf("could not create Harbor session verification request")
	}
	req.AddCookie(&http.Cookie{Name: "sid", Value: session})
	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("Harbor session verification failed: %w", safeHTTPError(err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		return fmt.Errorf("Harbor session invalidation unconfirmed (HTTP %d)", resp.StatusCode)
	}
	self.session = ""
	self.client.Jar.SetCookies(base, []*http.Cookie{{Name: "sid", Value: "", Path: "/", MaxAge: -1}})
	return nil
}

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
	self.ctx.Logger.Printf("init")
	config := &harbor_login_proto.Config{}
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
	self.ctx.Logger.Printf("creating OIDC request")
	login, err := newLogin(self.ctx, config)
	if err != nil {
		return nil, fmt.Errorf("could not create login: %w", err)
	}
	oidcUrl, err := login.createOIDCRequest(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create OIDC request: %w", err)
	}
	self.ctx.Logger.Printf("authorizing the OIDC request")
	vaultOidc, err := vault.OidcLoginContext(ctx, client.Client, oidcUrl)
	if err != nil {
		return nil, fmt.Errorf("could not login to OIDC provider: %w", err)
	}
	self.ctx.Logger.Printf("authorizing the session")
	if err := self.lc.AddState(lifecycle.StateStarted, lifecycle.StoppableFunc(login.logout)); err != nil {
		return nil, fmt.Errorf("could not register session cleanup: %w", err)
	}
	sessionId, err := login.authorizeSession(ctx, vaultOidc)
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
	lc.Add(plugin, server)
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

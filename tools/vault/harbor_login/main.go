package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al_plugin"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
	"git.alwaldend.com/alwaldend/src/tools/go/pkg/sensitiveerr"
	"git.alwaldend.com/alwaldend/src/tools/vault/harbor_login/harbor_login_proto"
)

// login owns an invocation's browser session and OIDC redirect state.
type login struct {
	// config selects the login endpoint and authentication settings.
	config *harbor_login_proto.Config
	// client retains cookies and enforces request and redirect limits.
	client *http.Client
	// redirect records the authorization destination advertised by the login endpoint.
	redirect *url.URL
	// session retains the original session ID for logout verification.
	session string
}

// newLogin configures a cookie jar and redirect policy for one login.
func newLogin(config *harbor_login_proto.Config) (*login, error) {
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return nil, fmt.Errorf("could not create cookiejar: %w", err)
	}
	res := &login{
		config: config,
	}
	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			loc, err := req.Response.Location()
			if err != nil {
				return sensitiveerr.WithContext("invalid redirect location", err)
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

// createOIDCRequest obtains the authorization URL for the configured identity provider.
func (l *login) createOIDCRequest(ctx context.Context) (*url.URL, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/c/oidc/login", l.config.HarborUrl),
		nil,
	)
	if err != nil {
		return nil, sensitiveerr.WithContext("could not create request", err)
	}
	resp, err := l.client.Do(req)
	if err != nil {
		return nil, sensitiveerr.WithContext("could not execute the request", err)
	}
	defer resp.Body.Close()
	if l.redirect == nil {
		return nil, fmt.Errorf("OIDC authorization redirect missing")
	}
	return l.redirect, nil
}

// authorizeSession exchanges a Vault-issued OIDC code for an authenticated browser session.
func (l *login) authorizeSession(ctx context.Context, oidc *al.VaultOidc) (string, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/c/oidc/callback", l.config.HarborUrl),
		nil,
	)
	if err != nil {
		return "", sensitiveerr.WithContext("could not create request", err)
	}
	query := req.URL.Query()
	query.Add("state", oidc.State)
	query.Add("code", oidc.Code)
	req.URL.RawQuery = query.Encode()
	resp, err := l.client.Do(req)
	if err != nil {
		return "", sensitiveerr.WithContext("could not execute the request", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("invalid response code: %d", resp.StatusCode)
	}
	res := ""
	for _, cookie := range l.client.Jar.Cookies(req.URL) {
		if cookie.Name == "sid" {
			res = cookie.Value
		}
	}
	if res == "" {
		return "", fmt.Errorf("missing sid cookie")
	}
	l.session = res
	return res, nil
}

// plugin serves AL login requests and owns cleanup of invocation-scoped credentials.
type plugin struct {
	// ctx provides command-scoped diagnostics for startup.
	ctx *al.CmdCtx
	// lc owns the cleanup actions run during shutdown.
	lc lifecycle.Manager
}

// logout ends the browser session and checks that its retained credentials are rejected.
func (l *login) logout(ctx context.Context) error {
	base, err := url.Parse(l.config.HarborUrl)
	if err != nil {
		return sensitiveerr.WithContext("invalid Harbor URL", err)
	}
	session := l.session
	if session == "" {
		for _, cookie := range l.client.Jar.Cookies(base) {
			if cookie.Name == "sid" {
				session = cookie.Value
			}
		}
	}
	if session == "" {
		return nil
	}
	client := *l.client
	client.Jar = nil // Keep the original SID to verify server-side invalidation.
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, l.config.HarborUrl+"/c/oidc/logout", nil)
	if err != nil {
		return sensitiveerr.WithContext("could not create Harbor logout request", err)
	}
	req.AddCookie(&http.Cookie{Name: "sid", Value: session})
	resp, err := client.Do(req)
	if err != nil {
		return sensitiveerr.WithContext("Harbor logout failed", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		return fmt.Errorf("Harbor logout returned HTTP %d", resp.StatusCode)
	}
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, l.config.HarborUrl+"/api/v2.0/users/current", nil)
	if err != nil {
		return sensitiveerr.WithContext("could not create Harbor session verification request", err)
	}
	req.AddCookie(&http.Cookie{Name: "sid", Value: session})
	resp, err = client.Do(req)
	if err != nil {
		return sensitiveerr.WithContext("Harbor session verification failed", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		return fmt.Errorf("Harbor session invalidation unconfirmed (HTTP %d)", resp.StatusCode)
	}
	l.session = ""
	l.client.Jar.SetCookies(base, []*http.Cookie{{Name: "sid", Value: "", Path: "/", MaxAge: -1}})
	return nil
}

// Start starts the AL plugin's lifecycle manager.
func (p *plugin) Start(ctx context.Context) error {
	if err := p.lc.Start(ctx); err != nil {
		return fmt.Errorf("start Harbor login plugin: %w", err)
	}
	return nil
}

// Stop releases resources registered by plugin startup in reverse order.
func (p *plugin) Stop(ctx context.Context) error {
	if err := p.lc.Stop(ctx); err != nil {
		return fmt.Errorf("stop Harbor login plugin: %w", err)
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
	p.ctx.Logger.Printf("init")
	if req.GetPlugin() == nil {
		return nil, errors.New("missing plugin configuration")
	}
	config := &harbor_login_proto.Config{}
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
	p.ctx.Logger.Printf("creating OIDC request")
	login, err := newLogin(config)
	if err != nil {
		return nil, fmt.Errorf("could not create login: %w", err)
	}
	oidcUrl, err := login.createOIDCRequest(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create OIDC request: %w", err)
	}
	p.ctx.Logger.Printf("authorizing the OIDC request")
	vaultOidc, err := vault.OidcLoginContext(ctx, client.Client, oidcUrl)
	if err != nil {
		return nil, sensitiveerr.WithContext("could not login to OIDC provider", err)
	}
	p.ctx.Logger.Printf("authorizing the session")
	if err := p.lc.AddState(lifecycle.StateStarted, lifecycle.StoppableFunc(login.logout)); err != nil {
		return nil, fmt.Errorf("could not register session cleanup: %w", err)
	}
	sessionID, err := login.authorizeSession(ctx, vaultOidc)
	if err != nil {
		return nil, fmt.Errorf("could not authorize the session using OIDC code: %w", err)
	}
	res := &al_proto.ResourceResult{
		Env: map[string]string{
			"HARBOR_SESSION_ID": sessionID,
			"HARBOR_URL":        config.HarborUrl,
		},
	}
	return &al_proto.PluginStartResponse{Result: res}, nil
}

// run serves plugin requests until cancellation or a lifecycle failure.
func run(ctx *al.CmdCtx) error {
	var lc lifecycle.Manager
	provider := &plugin{ctx: ctx}
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
	ctx := al.NewCmdCtx(shutdownCtx, "com.alwaldend.src.tools.vault.harbor_login ")
	if err := run(ctx); err != nil {
		ctx.Logger.Printf("failed: %s", err)
		os.Exit(1)
	}
}

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al_plugin"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
	"git.alwaldend.com/alwaldend/src/tools/go/pkg/sensitiveerr"
	"git.alwaldend.com/alwaldend/src/tools/vault/forgejo_login/forgejo_login_proto"
	"github.com/google/uuid"
	"golang.org/x/net/html"
	"google.golang.org/protobuf/proto"
)

// login owns an invocation's browser session and OIDC redirect state.
type login struct {
	// config selects the login endpoint and authentication settings.
	config *forgejo_login_proto.Config
	// client retains cookies and enforces request and redirect limits.
	client *http.Client
	// redirect records the authorization destination advertised by the login endpoint.
	redirect *url.URL
	// cookies captures redirect cookies while a token is being issued.
	cookies []*http.Cookie
	// allowExternalRedirect permits the initial browser redirect to the identity provider.
	allowExternalRedirect bool
}

// newLogin configures a cookie jar and redirect policy for one login.
func newLogin(config *forgejo_login_proto.Config) (*login, error) {
	origin, err := url.Parse(config.ForgejoUrl)
	if err != nil {
		return nil, sensitiveerr.WithContext("parse Forgejo URL", err)
	}
	if origin.Host == "" || (origin.Scheme != "http" && origin.Scheme != "https") || origin.User != nil {
		return nil, errors.New("invalid Forgejo URL")
	}
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return nil, fmt.Errorf("could not create cookiejar: %w", err)
	}
	res := &login{
		config: config,
	}
	client := &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !res.allowExternalRedirect && (!strings.EqualFold(req.URL.Host, origin.Host) || !strings.EqualFold(req.URL.Scheme, origin.Scheme)) {
				return errors.New("cross-origin Forgejo redirect rejected")
			}
			if len(via) >= 10 {
				return errors.New("too many redirects")
			}
			loc, err := req.Response.Location()
			if err != nil {
				return sensitiveerr.WithContext("invalid redirect location", err)
			}
			res.redirect = loc
			res.cookies = append(res.cookies, req.Response.Cookies()...)
			return nil
		},
	}
	res.client = client
	return res, nil
}

// createOIDCRequest obtains the authorization URL for the configured identity provider.
func (l *login) createOIDCRequest(ctx context.Context) (*url.URL, error) {
	l.allowExternalRedirect = true
	defer func() { l.allowExternalRedirect = false }()
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/user/oauth2/%s", l.config.ForgejoUrl, url.PathEscape(l.config.ForgejoOauthName)),
		nil,
	)
	if err != nil {
		return nil, sensitiveerr.WithContext("Forgejo request failed", err)
	}
	resp, err := l.client.Do(req)
	if err != nil {
		return nil, sensitiveerr.WithContext("Forgejo request failed", err)
	}
	defer resp.Body.Close()
	if l.redirect == nil {
		return nil, errors.New("OIDC redirect missing")
	}
	return l.redirect, nil
}

// authorizeSession exchanges a Vault-issued OIDC code for an authenticated browser session.
func (l *login) authorizeSession(ctx context.Context, oidc *al.VaultOidc) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/user/oauth2/%s/callback", l.config.ForgejoUrl, l.config.ForgejoOauthName),
		nil,
	)
	if err != nil {
		return sensitiveerr.WithContext("Forgejo request failed", err)
	}
	query := req.URL.Query()
	query.Add("state", oidc.State)
	query.Add("code", oidc.Code)
	req.URL.RawQuery = query.Encode()
	resp, err := l.client.Do(req)
	if err != nil {
		return sensitiveerr.WithContext("Forgejo request failed", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("invalid response code: %d", resp.StatusCode)
	}
	return nil
}

// createForgejoToken issues a named access token and reads its value from the redirect's flash cookie.
func (l *login) createForgejoToken(ctx context.Context, name string) (string, error) {
	l.cookies = nil
	data := url.Values{
		"name":        []string{name},
		"resource":    []string{"all"},
		"repo_search": []string{""},
		"page":        []string{"1"},
		"scope": []string{
			"write:activitypub",
			"write:admin",
			"write:issue",
			"write:misc",
			"write:notification",
			"write:organization",
			"write:package",
			"write:repository",
			"write:user",
		},
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/user/settings/applications/tokens/new/", l.config.ForgejoUrl),
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return "", sensitiveerr.WithContext("Forgejo request failed", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	l.addCSRF(req)
	resp, err := l.client.Do(req)
	if err != nil {
		return "", sensitiveerr.WithContext("Forgejo request failed", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token creation status: %d", resp.StatusCode)
	}
	for _, cookie := range l.cookies {
		if cookie.Name == "flash" {
			queryVal, err := url.QueryUnescape(cookie.Value)
			if err != nil {
				return "", sensitiveerr.WithContext("could not unescape the flash cookie", err)
			}
			query, err := url.ParseQuery(queryVal)
			if err != nil {
				return "", sensitiveerr.WithContext("could not parse the token cookie", err)
			}
			token, ok := query["info"]
			if !ok {
				return "", errors.New("token cookie is missing info")
			}
			if len(token) == 0 {
				return "", errors.New("token cookie info is empty")
			}
			if token[len(token)-1] == "" {
				return "", errors.New("empty token cookie")
			}
			return token[len(token)-1], nil
		}
	}
	return "", fmt.Errorf("missing flash cookie")
}

// attr returns the first matching HTML attribute or the empty string.
func attr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

// class matches a single token in a node's CSS class list.
func class(n *html.Node, name string) bool {
	for _, c := range strings.Fields(attr(n, "class")) {
		if c == name {
			return true
		}
	}
	return false
}

// walk visits the root and its descendants in document order.
func walk(n *html.Node, fn func(*html.Node)) {
	fn(n)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walk(c, fn)
	}
}

// nodeText collects descendant text nodes in document order.
func nodeText(n *html.Node) string {
	var b strings.Builder
	walk(n, func(n *html.Node) {
		if n.Type == html.TextNode {
			b.WriteString(n.Data)
		}
	})
	return b.String()
}

// parseTokenPage requires a complete applications settings page and identifies the requested token row.
func parseTokenPage(body io.Reader, settingsPath, name string) (string, error) {
	raw, err := io.ReadAll(io.LimitReader(body, (4<<20)+1))
	if err != nil {
		return "", sensitiveerr.WithContext("read Forgejo settings page", err)
	}
	if len(raw) > 4<<20 {
		return "", errors.New("could not read complete Forgejo settings page")
	}
	doc, err := html.Parse(strings.NewReader(string(raw)))
	if err != nil {
		return "", sensitiveerr.WithContext("invalid Forgejo settings page", err)
	}
	validPage, newLink := false, false
	id := ""
	var parseErr error
	walk(doc, func(n *html.Node) {
		if class(n, "page-content") && class(n, "user") && class(n, "settings") && class(n, "applications") {
			validPage = true
		}
		if n.Data == "a" && attr(n, "href") == settingsPath+"/tokens/new" {
			newLink = true
		}
		if !class(n, "flex-item") {
			return
		}
		title, tokenID := "", ""
		walk(n, func(c *html.Node) {
			if c.Data == "span" && class(c, "flex-item-title") {
				title = nodeText(c)
			}
			if c.Data == "button" && attr(c, "data-modal-id") == "delete-token" && attr(c, "data-url") == settingsPath+"/tokens/delete" {
				tokenID = attr(c, "data-id")
			}
		})
		if name == "" || title != name {
			return
		}
		numeric, e := strconv.ParseUint(tokenID, 10, 64)
		if e != nil {
			parseErr = sensitiveerr.WithContext("parse Forgejo token identifier", e)
			return
		}
		if numeric == 0 || id != "" {
			parseErr = errors.New("ambiguous or invalid Forgejo token row")
			return
		}
		id = tokenID
	})
	if !validPage || !newLink {
		return "", errors.New("Forgejo session did not return token settings")
	}
	if parseErr != nil {
		return "", fmt.Errorf("inspect token row: %w", parseErr)
	}
	return id, nil
}

// getTokenID looks up an invocation's token on the authenticated settings page.
func (l *login) getTokenID(ctx context.Context, tokenName string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, l.config.ForgejoUrl+"/user/settings/applications", nil)
	if err != nil {
		return "", sensitiveerr.WithContext("could not create token settings request", err)
	}
	resp, err := l.client.Do(req)
	if err != nil {
		return "", sensitiveerr.WithContext("token settings request failed", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || resp.Request.URL.Path != req.URL.Path {
		return "", errors.New("Forgejo session did not return token settings")
	}
	id, err := parseTokenPage(resp.Body, req.URL.Path, tokenName)
	if err != nil {
		return "", fmt.Errorf("inspect token settings: %w", err)
	}
	return id, nil
}

// logout ends the browser session and checks that its retained credentials are rejected.
func (l *login) logout(ctx context.Context) error {
	endpoint, err := url.Parse(l.config.ForgejoUrl + "/user/settings/applications")
	if err != nil {
		return sensitiveerr.WithContext("invalid Forgejo settings URL", err)
	}
	// Keep the original cookie value to verify server-side invalidation.
	originalCookies := l.client.Jar.Cookies(endpoint)
	client := *l.client
	client.Jar = nil
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	request := func(method, path string) (*http.Response, error) {
		req, err := http.NewRequestWithContext(ctx, method, l.config.ForgejoUrl+path, nil)
		if err != nil {
			return nil, sensitiveerr.WithContext("could not create Forgejo session request", err)
		}
		for _, cookie := range originalCookies {
			req.AddCookie(cookie)
		}
		l.addCSRF(req)
		resp, err := client.Do(req)
		if err != nil {
			return nil, sensitiveerr.WithContext("Forgejo session request failed", err)
		}
		return resp, nil
	}
	invalidated := func() (bool, error) {
		resp, err := request(http.MethodGet, "/user/settings/applications")
		if err != nil {
			return false, fmt.Errorf("request authenticated settings: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusSeeOther || resp.StatusCode == http.StatusTemporaryRedirect {
			target, err := resp.Location()
			if err != nil {
				return false, sensitiveerr.WithContext("invalid Forgejo session redirect", err)
			}
			return target.Scheme == endpoint.Scheme && target.Host == endpoint.Host && target.Path == strings.TrimSuffix(endpoint.Path, "/settings/applications")+"/login", nil
		}
		if resp.StatusCode != http.StatusOK {
			return false, fmt.Errorf("verify Forgejo session: HTTP %d", resp.StatusCode)
		}
		if _, err := parseTokenPage(resp.Body, endpoint.Path, ""); err != nil {
			return false, fmt.Errorf("inspect authenticated settings: %w", err)
		}
		return false, nil
	}
	done, err := invalidated()
	if err != nil {
		return fmt.Errorf("check session before logout: %w", err)
	}
	if done {
		return nil
	}
	resp, err := request(http.MethodPost, "/user/logout")
	if err != nil {
		return fmt.Errorf("request Forgejo logout: %w", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusSeeOther {
		return fmt.Errorf("Forgejo logout status: %d", resp.StatusCode)
	}
	done, err = invalidated()
	if err != nil {
		return fmt.Errorf("verify session after logout: %w", err)
	}
	if !done {
		return errors.New("Forgejo session remains authenticated after logout")
	}

	return nil
}

// deleteForgejoToken removes the named token and verifies its absence from account settings.
func (l *login) deleteForgejoToken(ctx context.Context, tokenName string) error {
	id, err := l.getTokenID(ctx, tokenName)
	if err != nil {
		return fmt.Errorf("locate token before deletion: %w", err)
	}
	if id == "" {
		return nil
	}
	data := url.Values{"id": {id}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, l.config.ForgejoUrl+"/user/settings/applications/tokens/delete", strings.NewReader(data.Encode()))
	if err != nil {
		return sensitiveerr.WithContext("could not create token deletion request", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	l.addCSRF(req)
	resp, err := l.client.Do(req)
	if err != nil {
		return sensitiveerr.WithContext("token deletion request failed", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token deletion status: %d", resp.StatusCode)
	}
	id, err = l.getTokenID(ctx, tokenName)
	if err != nil {
		return fmt.Errorf("verify token removal: %w", err)
	}
	if id != "" {
		return errors.New("Forgejo token still exists after deletion")
	}
	return nil
}

// addCSRF copies the current session's CSRF cookie into the request header.
func (l *login) addCSRF(req *http.Request) {
	for _, cookie := range l.client.Jar.Cookies(req.URL) {
		if cookie.Name == "_csrf" {
			req.Header.Set("X-Csrf-Token", cookie.Value)
		}
	}
}

// plugin serves AL login requests and owns cleanup of invocation-scoped credentials.
type plugin struct {
	// ctx provides command-scoped diagnostics for startup.
	ctx *al.CmdCtx
	// lc owns the cleanup actions run during shutdown.
	lc lifecycle.Manager
}

// Start starts the AL plugin's lifecycle manager.
func (p *plugin) Start(ctx context.Context) error {
	if err := p.lc.Start(ctx); err != nil {
		return fmt.Errorf("start Forgejo login plugin: %w", err)
	}
	return nil
}

// Stop releases resources registered by plugin startup in reverse order.
func (p *plugin) Stop(ctx context.Context) error {
	if err := p.lc.Stop(ctx); err != nil {
		return fmt.Errorf("stop Forgejo login plugin: %w", err)
	}
	return nil
}

// parsePluginConfig merges defaults with configuration from at most one explicitly configured call.
func parsePluginConfig(plugin *al_proto.PluginConfig) (*forgejo_login_proto.Config, error) {
	if plugin == nil {
		return nil, errors.New("missing plugin configuration")
	}
	config := &forgejo_login_proto.Config{}
	if _, err := al.FromPbJsonToPb(plugin.Data, config).Get(); err != nil {
		return nil, sensitiveerr.WithContext("could not parse plugin data", err)
	}
	var selected *forgejo_login_proto.Config
	for _, call := range plugin.Calls {
		overrides := &forgejo_login_proto.Config{}
		if _, err := al.FromPbJsonToPb(call.GetData(), overrides).Get(); err != nil {
			return nil, sensitiveerr.WithContext("could not parse plugin call data", err)
		}
		if proto.Size(overrides) == 0 {
			continue
		}
		if selected != nil {
			return nil, errors.New("multiple Forgejo login calls specify configuration overrides")
		}
		selected = overrides
	}
	if selected != nil {
		proto.Merge(config, selected)
	}
	return config, nil
}

// PluginStart authenticates the invocation, registers cleanup, and returns credentials to AL.
func (p *plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (res *al_proto.PluginStartResponse, err error) {
	defer func() {
		if err != nil {
			cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err = errors.Join(err, p.lc.Stop(cleanupCtx))
		}
	}()
	p.ctx.Logger.Printf("init")
	config, err := parsePluginConfig(req.GetPlugin())
	if err != nil {
		return nil, fmt.Errorf("read Forgejo configuration: %w", err)
	}
	vault := al.NewVault(req.Config)
	if err := p.lc.AddState(lifecycle.StateStarted, lifecycle.StoppableFunc(vault.Stop)); err != nil {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return nil, fmt.Errorf("register Vault cleanup: %w", errors.Join(err, vault.Stop(cleanupCtx)))
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
	if err := p.lc.AddState(lifecycle.StateStarted, lifecycle.StoppableFunc(login.logout)); err != nil {
		return nil, fmt.Errorf("register session cleanup: %w", err)
	}
	p.ctx.Logger.Printf("authorizing the session")
	if err := login.authorizeSession(ctx, vaultOidc); err != nil {
		return nil, fmt.Errorf("could not authorize the session using OIDC code: %w", err)
	}
	tokenName := fmt.Sprintf("src_tools_vault_forgejo_login_%s", uuid.New().String())
	if _, err := login.getTokenID(ctx, tokenName); err != nil {
		return nil, fmt.Errorf("cannot establish token cleanup access: %w", err)
	}
	// Register by name before issuance: a lost response may still create a token.
	if err := p.lc.AddState(lifecycle.StateStarted, lifecycle.StoppableFunc(func(ctx context.Context) error {
		if err := login.deleteForgejoToken(ctx, tokenName); err != nil {
			return fmt.Errorf("remove invocation token: %w", err)
		}
		return nil
	})); err != nil {
		return nil, fmt.Errorf("register token cleanup: %w", err)
	}
	token, err := login.createForgejoToken(ctx, tokenName)
	if err != nil {
		return nil, fmt.Errorf("could not create Forgejo token: %w", err)
	}
	res = &al_proto.PluginStartResponse{Result: &al_proto.ResourceResult{
		Env: map[string]string{
			"FORGEJO_API_TOKEN": token,
			"FORGEJO_HOST":      config.ForgejoUrl,
		},
	}}
	return res, nil
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
	ctx := al.NewCmdCtx(shutdownCtx, "com.alwaldend.src.tools.vault.forgejo_login ")
	if err := run(ctx); err != nil {
		ctx.Logger.Printf("failed: %s", err)
		os.Exit(1)
	}
}

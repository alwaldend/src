package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	reqUrl "net/url"
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
	"git.alwaldend.com/alwaldend/src/tools/vault/forgejo_login/forgejo_login_proto"
	"github.com/google/uuid"
	"golang.org/x/net/html"
)

type login struct {
	ctx                   *al.CmdCtx
	config                *forgejo_login_proto.Config
	client                *http.Client
	redirect              *reqUrl.URL
	cookies               []*http.Cookie
	allowExternalRedirect bool
}

func newLogin(ctx *al.CmdCtx, config *forgejo_login_proto.Config) (*login, error) {
	origin, err := reqUrl.Parse(config.ForgejoUrl)
	if err != nil || origin.Host == "" || (origin.Scheme != "http" && origin.Scheme != "https") || origin.User != nil {
		return nil, errors.New("invalid Forgejo URL")
	}
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
				return errors.New("invalid redirect location")
			}
			res.redirect = loc
			res.cookies = append(res.cookies, req.Response.Cookies()...)
			return nil
		},
	}
	res.client = client
	return res, nil
}

func (self *login) createOIDCRequest(ctx context.Context) (*reqUrl.URL, error) {
	self.allowExternalRedirect = true
	defer func() { self.allowExternalRedirect = false }()
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/user/oauth2/%s", self.config.ForgejoUrl, reqUrl.PathEscape(self.config.ForgejoOauthName)),
		nil,
	)
	if err != nil {
		return nil, errors.New("Forgejo request failed")
	}
	resp, err := self.client.Do(req)
	if err != nil {
		return nil, errors.New("Forgejo request failed")
	}
	defer resp.Body.Close()
	if self.redirect == nil {
		return nil, errors.New("OIDC redirect missing")
	}
	return self.redirect, nil
}

func (self *login) authorizeSession(ctx context.Context, oidc *al.VaultOidc) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/user/oauth2/%s/callback", self.config.ForgejoUrl, self.config.ForgejoOauthName),
		nil,
	)
	if err != nil {
		return errors.New("Forgejo request failed")
	}
	query := req.URL.Query()
	query.Add("state", oidc.State)
	query.Add("code", oidc.Code)
	req.URL.RawQuery = query.Encode()
	resp, err := self.client.Do(req)
	if err != nil {
		return errors.New("Forgejo request failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("invalid response code: %d", resp.StatusCode)
	}
	return nil
}

func (self *login) createForgejoToken(ctx context.Context, name string) (string, error) {
	self.cookies = nil
	data := reqUrl.Values{
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
		fmt.Sprintf("%s/user/settings/applications/tokens/new/", self.config.ForgejoUrl),
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return "", errors.New("Forgejo request failed")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	self.addCSRF(req)
	resp, err := self.client.Do(req)
	if err != nil {
		return "", errors.New("Forgejo request failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token creation status: %d", resp.StatusCode)
	}
	for _, cookie := range self.cookies {
		if cookie.Name == "flash" {
			queryVal, err := reqUrl.QueryUnescape(cookie.Value)
			if err != nil {
				return "", errors.New("could not unescape the flash cookie")
			}
			query, err := reqUrl.ParseQuery(queryVal)
			if err != nil {
				return "", errors.New("could not parse the token cookie")
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

func attr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func class(n *html.Node, name string) bool {
	for _, c := range strings.Fields(attr(n, "class")) {
		if c == name {
			return true
		}
	}
	return false
}

func walk(n *html.Node, fn func(*html.Node)) {
	fn(n)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walk(c, fn)
	}
}

func nodeText(n *html.Node) string {
	var b strings.Builder
	walk(n, func(n *html.Node) {
		if n.Type == html.TextNode {
			b.WriteString(n.Data)
		}
	})
	return b.String()
}

// Forgejo 15 lists every token on this browser-session page without pagination.
func parseTokenPage(body io.Reader, settingsPath, name string) (string, error) {
	raw, err := io.ReadAll(io.LimitReader(body, (4<<20)+1))
	if err != nil || len(raw) > 4<<20 {
		return "", errors.New("could not read complete Forgejo settings page")
	}
	doc, err := html.Parse(strings.NewReader(string(raw)))
	if err != nil {
		return "", errors.New("invalid Forgejo settings page")
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
		if e != nil || numeric == 0 || id != "" {
			parseErr = errors.New("ambiguous or invalid Forgejo token row")
			return
		}
		id = tokenID
	})
	if !validPage || !newLink {
		return "", errors.New("Forgejo session did not return token settings")
	}
	return id, parseErr
}

func (self *login) getTokenId(ctx context.Context, tokenName string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, self.config.ForgejoUrl+"/user/settings/applications", nil)
	if err != nil {
		return "", errors.New("could not create token settings request")
	}
	resp, err := self.client.Do(req)
	if err != nil {
		return "", errors.New("token settings request failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || resp.Request.URL.Path != req.URL.Path {
		return "", errors.New("Forgejo session did not return token settings")
	}
	return parseTokenPage(resp.Body, req.URL.Path, tokenName)
}

func (self *login) logout(ctx context.Context) error {
	endpoint, err := reqUrl.Parse(self.config.ForgejoUrl + "/user/settings/applications")
	if err != nil {
		return errors.New("invalid Forgejo settings URL")
	}
	// Keep the original cookie value to verify server-side invalidation.
	originalCookies := self.client.Jar.Cookies(endpoint)
	client := *self.client
	client.Jar = nil
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	request := func(method, path string) (*http.Response, error) {
		req, err := http.NewRequestWithContext(ctx, method, self.config.ForgejoUrl+path, nil)
		if err != nil {
			return nil, errors.New("could not create Forgejo session request")
		}
		for _, cookie := range originalCookies {
			req.AddCookie(cookie)
		}
		self.addCSRF(req)
		resp, err := client.Do(req)
		if err != nil {
			return nil, errors.New("Forgejo session request failed")
		}
		return resp, nil
	}
	invalidated := func() (bool, error) {
		resp, err := request(http.MethodGet, "/user/settings/applications")
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusSeeOther || resp.StatusCode == http.StatusTemporaryRedirect {
			target, err := resp.Location()
			if err != nil {
				return false, errors.New("invalid Forgejo session redirect")
			}
			return target.Scheme == endpoint.Scheme && target.Host == endpoint.Host && target.Path == strings.TrimSuffix(endpoint.Path, "/settings/applications")+"/login", nil
		}
		if resp.StatusCode != http.StatusOK {
			return false, errors.New("could not verify Forgejo session")
		}
		_, err = parseTokenPage(resp.Body, endpoint.Path, "")
		return false, err
	}
	done, err := invalidated()
	if err != nil || done {
		return err
	}
	resp, err := request(http.MethodPost, "/user/logout")
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Forgejo logout status: %d", resp.StatusCode)
	}
	done, err = invalidated()
	if err != nil {
		return err
	}
	if !done {
		return errors.New("Forgejo session remains authenticated after logout")
	}

	return nil
}

func (self *login) deleteForgejoToken(ctx context.Context, tokenName string) error {
	id, err := self.getTokenId(ctx, tokenName)
	if err != nil {
		return err
	}
	if id == "" {
		return nil
	}
	data := reqUrl.Values{"id": {id}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, self.config.ForgejoUrl+"/user/settings/applications/tokens/delete", strings.NewReader(data.Encode()))
	if err != nil {
		return errors.New("could not create token deletion request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	self.addCSRF(req)
	resp, err := self.client.Do(req)
	if err != nil {
		return errors.New("token deletion request failed")
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token deletion status: %d", resp.StatusCode)
	}
	id, err = self.getTokenId(ctx, tokenName)
	if err != nil {
		return err
	}
	if id != "" {
		return errors.New("Forgejo token still exists after deletion")
	}
	return nil
}

func (self *login) addCSRF(req *http.Request) {
	for _, cookie := range self.client.Jar.Cookies(req.URL) {
		if cookie.Name == "_csrf" {
			req.Header.Set("X-Csrf-Token", cookie.Value)
		}
	}
}

type Plugin struct {
	ctx *al.CmdCtx
	lc  lifecycle.Manager
}

func (self *Plugin) Start(ctx context.Context) error {
	return self.lc.Start(ctx)
}

func (self *Plugin) Stop(ctx context.Context) error {
	return self.lc.Stop(ctx)
}

func (self *Plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (res *al_proto.PluginStartResponse, err error) {
	defer func() {
		if err != nil {
			cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err = errors.Join(err, self.lc.Stop(cleanupCtx))
		}
	}()
	self.ctx.Logger.Printf("init")
	config := &forgejo_login_proto.Config{}
	if _, err := al.FromPbJsonToPb(req.Plugin.Data, config).Get(); err != nil {
		return nil, errors.New("could not parse plugin data")
	}
	vault := al.NewVault(req.Config)
	if err := self.lc.AddState(lifecycle.StateStarted, lifecycle.StoppableFunc(vault.Stop)); err != nil {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return nil, errors.Join(err, vault.Stop(cleanupCtx))
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
	if err := self.lc.AddState(lifecycle.StateStarted, lifecycle.StoppableFunc(login.logout)); err != nil {
		return nil, err
	}
	self.ctx.Logger.Printf("authorizing the session")
	if err := login.authorizeSession(ctx, vaultOidc); err != nil {
		return nil, fmt.Errorf("could not authorize the session using OIDC code: %w", err)
	}
	tokenName := fmt.Sprintf("src_tools_vault_forgejo_login_%s", uuid.New().String())
	if _, err := login.getTokenId(ctx, tokenName); err != nil {
		return nil, fmt.Errorf("cannot establish token cleanup access: %w", err)
	}
	// Register by name before issuance: a lost response may still create a token.
	if err := self.lc.AddState(lifecycle.StateStarted, lifecycle.StoppableFunc(func(ctx context.Context) error {
		return login.deleteForgejoToken(ctx, tokenName)
	})); err != nil {
		return nil, err
	}
	token, err := login.createForgejoToken(ctx, tokenName)
	if err != nil {
		return nil, fmt.Errorf("could not create Forgejo token: %w", err)
	}
	res = &al_proto.PluginStartResponse{
		Env: map[string]string{
			"FORGEJO_API_TOKEN": token,
			"FORGEJO_HOST":      config.ForgejoUrl,
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
	ctx := al.NewCmdCtx(shutdownCtx, "com.alwaldend.src.tools.vault.forgejo_login ")
	if err := run(ctx); err != nil {
		ctx.Logger.Printf("failed: %s", err)
		os.Exit(1)
	}
}

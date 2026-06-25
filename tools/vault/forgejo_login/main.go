package main

import (
	"context"
	"encoding/json"
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
)

type login struct {
	ctx      *al.CmdCtx
	config   *forgejo_login_proto.Config
	client   *http.Client
	redirect *reqUrl.URL
	cookies  []*http.Cookie
}

func newLogin(ctx *al.CmdCtx, config *forgejo_login_proto.Config) (*login, error) {
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
		fmt.Sprintf("%s/user/oauth2/vault", self.config.ForgejoUrl),
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
	return self.redirect, nil
}

func (self *login) authorizeSession(oidc *al.VaultOidc) error {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/user/oauth2/%s/callback", self.config.ForgejoUrl, self.config.ForgejoOauthName),
		nil,
	)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}
	query := req.URL.Query()
	query.Add("state", oidc.State)
	query.Add("code", oidc.Code)
	req.URL.RawQuery = query.Encode()
	resp, err := self.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not execute the request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("invalid response code: %s", resp.Status)
	}
	return nil
}

func (self *login) createForgejoToken() (string, string, error) {
	name := fmt.Sprintf("src_tools_vault_forgejo_login_%s", uuid.New().String())
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
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/user/settings/applications/tokens/new/", self.config.ForgejoUrl),
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return "", "", fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := self.client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("could not execute the request: %w", err)
	}
	defer resp.Body.Close()
	for _, cookie := range self.cookies {
		if cookie.Name == "flash" {
			queryVal, err := reqUrl.QueryUnescape(cookie.Value)
			if err != nil {
				return "", "", fmt.Errorf("could not unescape the flash cookie: %w", err)
			}
			query, err := reqUrl.ParseQuery(queryVal)
			if err != nil {
				return "", "", fmt.Errorf("could not parse the query from the token cookie: %w", err)
			}
			token, ok := query["info"]
			if !ok {
				return "", "", fmt.Errorf("token cookie is missing info: %w", err)
			}
			if len(token) == 0 {
				return "", "", fmt.Errorf("field info of the token cookie is empty: %w", err)
			}
			return name, token[len(token)-1], nil
		}
	}
	return "", "", fmt.Errorf("missing flash cookie")
}

func (self *login) getUser(ctx context.Context, token string) (string, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/api/v1/user", self.config.ForgejoUrl),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not execute request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("invalid status code: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not read response body: %w", err)
	}
	data := struct {
		Login string `json:"login"`
	}{}
	if err = json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("could not unmarshal response: %w", err)
	}
	return data.Login, nil
}

func (self *login) getTokenId(ctx context.Context, user string, tokenName string, token string) (string, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/api/v1/users/%s/tokens", self.config.ForgejoUrl, user),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not execute request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("invalid status code: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not read response body: %w", err)
	}
	data := []struct {
		Name string `json:"name"`
		Id   int    `json:"id"`
	}{}
	if err = json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("could not unmarshal response: %w", err)
	}
	for _, tokenData := range data {
		if tokenData.Name == tokenName {
			return strconv.Itoa(tokenData.Id), nil
		}
	}
	return "", fmt.Errorf("could not find token %s in %v", tokenName, data)
}

func (self *login) deleteForgejoToken(ctx context.Context, tokenName string, token string) error {
	user, err := self.getUser(ctx, token)
	if err != nil {
		return fmt.Errorf("could not get user info: %w", err)
	}
	tokenId, err := self.getTokenId(ctx, user, tokenName, token)
	if err != nil {
		return fmt.Errorf("could not find token id: %w", err)
	}
	reqBody := fmt.Sprintf(`
------geckoformboundary4d997fff18b0a002bfca28644f1b05d9
Content-Disposition: form-data; name="id"

%s
------geckoformboundary4d997fff18b0a002bfca28644f1b05d9--
`, tokenId)
	reqUrl := fmt.Sprintf("%s/user/settings/applications/tokens/delete", self.config.ForgejoUrl)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		reqUrl,
		strings.NewReader(reqBody),
	)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}
	resp, err := self.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not execute request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("could not read response body: %w", err)
		}
		return fmt.Errorf("invalid status code %s: %s", resp.Status, string(body))
	}
	return nil

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

func (self *Plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (*al_proto.PluginStartResponse, error) {
	self.ctx.Logger.Printf("init")
	config := &forgejo_login_proto.Config{}
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
	if err := login.authorizeSession(vaultOidc); err != nil {
		return nil, fmt.Errorf("could not authorize the session using OIDC code: %w", err)
	}
	self.ctx.Logger.Printf("creating the token")
	tokenName, token, err := login.createForgejoToken()
	if err != nil {
		return nil, fmt.Errorf("could not create forgejo token: %w", err)
	}
	self.lc.AddState(lifecycle.StateStarted, lifecycle.StoppableFunc(func(ctx context.Context) error {
		self.ctx.Logger.Printf("deleting token %s", tokenName)
		if err := login.deleteForgejoToken(ctx, tokenName, token); err != nil {
			return fmt.Errorf("could not delete forgejo token %s: %w", tokenName, err)
		}
		return nil
	}))
	res := &al_proto.PluginStartResponse{
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
	lc.Add(server, plugin)
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

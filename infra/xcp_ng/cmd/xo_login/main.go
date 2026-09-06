// xo_login authenticates with the caller's Vault identity for an AL invocation.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.alwaldend.com/alwaldend/src/infra/xcp_ng/internal/oidclogin"
	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al_plugin"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
)

type config struct {
	URL           string `json:"xoa_url"`
	Discovery     string `json:"discovery_url"`
	VaultConn     string `json:"vault_conn"`
	VaultAuth     string `json:"vault_auth"`
	IssuerAuth    string `json:"issuer_auth"`
	TargetApprole string `json:"target_approle"`
}

type plugin struct{ lc lifecycle.Manager }

func (p *plugin) Start(ctx context.Context) error { return p.lc.Start(ctx) }
func (p *plugin) Stop(ctx context.Context) error  { return p.lc.Stop(ctx) }

func origin(raw string) (*url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || (u.Path != "" && u.Path != "/") || u.RawQuery != "" || u.Fragment != "" {
		return nil, errors.New("XO login requires an HTTPS origin")
	}
	u.Path = "/"
	return u, nil
}

func (p *plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (res *al_proto.PluginStartResponse, err error) {
	defer func() {
		if err != nil {
			cleanup, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			err = errors.Join(err, p.Stop(cleanup))
		}
	}()
	value, err := al.FromPbJson(req.Plugin.Data).Get()
	if err != nil {
		return nil, errors.New("invalid XO login configuration")
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, errors.New("invalid XO login configuration")
	}
	var conf config
	if json.Unmarshal(raw, &conf) != nil {
		return nil, errors.New("invalid XO login configuration")
	}
	xo, err := origin(conf.URL)
	if err != nil {
		return nil, err
	}
	discovery, err := url.Parse(conf.Discovery)
	if err != nil {
		return nil, errors.New("invalid discovery URL")
	}
	vault := al.NewVault(req.Config)
	if err := p.lc.AddState(lifecycle.StateStarted, lifecycle.StoppableFunc(vault.Stop)); err != nil {
		return nil, err
	}
	authName := conf.VaultAuth
	if conf.IssuerAuth != "" {
		authName = conf.IssuerAuth
	}
	item, err := vault.Client(ctx, conf.VaultConn, authName).Get()
	if err != nil {
		return nil, errors.New("authenticating own Vault identity failed")
	}
	client := item.Client
	if conf.IssuerAuth != "" {
		client, err = issueRole(ctx, item.Client, conf.TargetApprole)
		if err != nil {
			return nil, err
		}
		if err := p.lc.AddState(lifecycle.StateStarted, lifecycle.StoppableFunc(func(cleanup context.Context) error {
			if _, err := client.Logical().WriteWithContext(cleanup, "auth/token/revoke-self", nil); err != nil {
				return errors.New("revoking target AppRole token failed")
			}
			client.ClearToken()
			return nil
		})); err != nil {
			return nil, err
		}
	}
	vaultOrigin, err := origin(client.Address())
	if err != nil {
		return nil, err
	}
	noRedirect := func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	vaultClient := *client.CloneConfig().HttpClient
	vaultClient.CheckRedirect = noRedirect
	vaultClient.Timeout = 45 * time.Second
	jar, _ := cookiejar.New(nil)
	xoClient := &http.Client{Jar: jar, Timeout: 45 * time.Second, CheckRedirect: noRedirect}
	// Register cleanup before the callback: a failed response may still set a token.
	if err := p.lc.AddState(lifecycle.StateStarted, lifecycle.StoppableFunc(func(cleanup context.Context) error {
		for _, cookie := range jar.Cookies(xo) {
			if cookie.Name == "token" && cookie.Value != "" {
				return oidclogin.Revoke(cleanup, xo, cookie.Value)
			}
		}
		return nil
	})); err != nil {
		return nil, err
	}
	if err := oidclogin.Bootstrap(xoClient, &vaultClient, xo, vaultOrigin, discovery, client.Token()); err != nil {
		return nil, err
	}
	for _, cookie := range jar.Cookies(xo) {
		if cookie.Name == "token" && cookie.Value != "" {
			endpoint := *xo
			endpoint.Scheme = "wss"
			return &al_proto.PluginStartResponse{Env: map[string]string{"XOA_TOKEN": cookie.Value, "XOA_URL": endpoint.String(), "XOA_INSECURE": "false"}}, nil
		}
	}
	return nil, errors.New("XO login did not issue a session token")
}

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--report" {
		if os.Getenv("XOA_TOKEN") == "" || os.Getenv("XOA_INSECURE") != "false" {
			fmt.Fprintln(os.Stderr, "XO login session is unavailable")
			os.Exit(1)
		}
		if expected := os.Getenv("XO_EXPECT_RESOURCE_SET"); expected != "" {
			endpoint, err := url.Parse(os.Getenv("XOA_URL"))
			if err == nil {
				endpoint.Scheme = "https"
				endpoint, err = origin(endpoint.String())
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, "invalid XO verification origin")
				os.Exit(1)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			report, err := oidclogin.VerifyScope(ctx, endpoint, os.Getenv("XOA_TOKEN"), expected, os.Getenv("XO_EXPECT_VM_ID"))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			if json.NewEncoder(os.Stdout).Encode(report) != nil {
				os.Exit(1)
			}
			return
		}
		if os.Getenv("XO_EXPECT_VM_ID") != "" {
			fmt.Fprintln(os.Stderr, "XO_EXPECT_RESOURCE_SET is required with XO_EXPECT_VM_ID")
			os.Exit(1)
		}
		fmt.Println(`{"oidc_callback_completed":true}`)
		return
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	cmd := al.NewCmdCtx(ctx, "xo_login ")
	p := &plugin{}
	var lc lifecycle.Manager
	lc.Add(p, al_plugin.NewPluginServer(cmd, p))
	if err := lc.Run(ctx, 15*time.Second); err != nil {
		cmd.Logger.Printf("XO login failed: %s", err)
		os.Exit(1)
	}
}

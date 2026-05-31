package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al_plugin"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"git.alwaldend.com/alwaldend/src/tools/vault/pve_login/pve_login_proto"
	"github.com/google/uuid"
	"github.com/hashicorp/vault/api"
)

var logger = log.New(os.Stderr, "com.alwaldend.src.tools.vault.pve_login ", log.Flags())

type state struct {
	ctx     context.Context
	oidcUrl struct {
		Data string `json:"data"`
	}
	vaultToken string
	vaultOidc  struct {
		Code  string `json:"code"`
		State string `json:"state"`
	}
	pveTicket struct {
		Data struct {
			CSRFPreventionToken string `json:"CSRFPreventionToken"`
			Username            string `json:"username"`
			Ticket              string `json:"ticket"`
		} `json:"data"`
	}
	pveToken struct {
		Data struct {
			TokenId     string `json:"full-tokenid"`
			TokenSecret string `json:"value"`
		} `json:"data"`
	}
	config *pve_login_proto.Config
	req    *al_proto.PluginStartRequest
	resp   *al_proto.PluginStartResponse
}

func (self *state) left(err error) fp.Either[*state] {
	return fp.Left[*state](err)
}

func (self *state) right() fp.Either[*state] {
	return fp.Right(self)
}

type stateMonad = fp.Either[*state]

func createVaultToken(s *state) stateMonad {
	return fp.Pipe3(
		al.NewVault,
		func(v *al.Vault) fp.Either[*api.Client] {
			return v.Client(s.ctx, s.config.VaultConn, s.config.VaultAuth)
		},
		func(v *api.Client) stateMonad {
			s.vaultToken = v.Token()
			return fp.Right(s)
		},
	)(s.req.Config)
}

func createOIDCRequest(s *state) stateMonad {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/api2/json/access/openid/auth-url?realm=%s&redirect-url=%s", s.config.PveBaseUrl, s.config.PveRealm, s.config.PveRedirectUrl),
		strings.NewReader("{}"),
	)
	if err != nil {
		return s.left(fmt.Errorf("could not create request: %w", err))
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return s.left(fmt.Errorf("could not execute request: %w", err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return s.left(fmt.Errorf("invalid status code: %s", resp.Status))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return s.left(fmt.Errorf("could not read response body"))
	}
	err = json.Unmarshal(body, &s.oidcUrl)
	if err != nil {
		return s.left(fmt.Errorf("could not unmarshal response: %w", err))
	}
	return s.right()
}

func loginToOIDCProvider(s *state) stateMonad {
	regex, err := regexp.Compile("^https://([^/]+)/ui/vault/identity/oidc/provider/([^/]+)/authorize?(.+)$")
	if err != nil {
		return s.left(fmt.Errorf("could not compile regex: %w", err))
	}
	parts := regex.FindAllStringSubmatch(s.oidcUrl.Data, -1)
	urlParsed, err := url.Parse(s.oidcUrl.Data)
	if err != nil {
		return s.left(fmt.Errorf("could not parse url: %w", err))
	}
	reqData := map[string]string{}
	for key, valueSlice := range urlParsed.Query() {
		for _, value := range valueSlice {
			reqData[key] = value
		}
	}
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return s.left(fmt.Errorf("could not marshal data: %w", err))
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://%s/v1/identity/oidc/provider/%s/authorize", parts[0][1], parts[0][2]),
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return s.left(fmt.Errorf("could not create request: %w", err))
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.vaultToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return s.left(fmt.Errorf("could not execute the request: %w", err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return s.left(fmt.Errorf("invalid response code: %s", resp.Status))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return s.left(fmt.Errorf("could not read the body: %w", err))
	}
	err = json.Unmarshal(body, &s.vaultOidc)
	if err != nil {
		return s.left(fmt.Errorf("could not unmarshal response body: %w", err))
	}
	return s.right()
}

func createProxmoxTicket(s *state) stateMonad {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/api2/json/access/openid/login", s.config.PveBaseUrl),
		strings.NewReader("{}"),
	)
	if err != nil {
		return s.left(fmt.Errorf("could not create request: %w", err))
	}
	query := req.URL.Query()
	query.Add("state", s.vaultOidc.State)
	query.Add("code", s.vaultOidc.Code)
	query.Add("redirect-url", s.config.PveRedirectUrl)
	req.URL.RawQuery = query.Encode()
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return s.left(fmt.Errorf("could not execute the request: %w", err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return s.left(fmt.Errorf("invalid response code: %s", resp.Status))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return s.left(fmt.Errorf("could not read response body: %w", err))
	}
	err = json.Unmarshal(body, &s.pveTicket)
	if err != nil {
		return s.left(fmt.Errorf("could not unmarshel response body: %w", err))
	}
	return s.right()
}

func createProxmoxToken(s *state) stateMonad {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"%s/api2/json/access/users/%s/token/%s",
			s.config.PveBaseUrl,
			s.pveTicket.Data.Username,
			fmt.Sprintf("tools-vault-pve-login-%s", uuid.New().String()),
		),
		strings.NewReader("{}"),
	)
	if err != nil {
		return s.left(fmt.Errorf("could not create request: %w", err))
	}
	query := req.URL.Query()
	query.Add("comment", "Automatically created by //tools/vault/pve_login")
	query.Add("expire", fmt.Sprintf("%d", time.Now().Add(time.Hour).Unix()))
	query.Add("privsep", "0")
	req.URL.RawQuery = query.Encode()
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("CSRFPreventionToken", s.pveTicket.Data.CSRFPreventionToken)
	req.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: s.pveTicket.Data.Ticket})
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return s.left(fmt.Errorf("could not execute the request: %w", err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return s.left(fmt.Errorf("invalid response code: %s", resp.Status))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return s.left(fmt.Errorf("could not read response body: %w", err))
	}
	err = json.Unmarshal(body, &s.pveToken)
	if err != nil {
		return s.left(fmt.Errorf("could not unmarshal response body: %w", err))
	}
	return s.right()
}

func parseConfig(s *state) stateMonad {
	_, err := fp.Get(al_plugin.ParseConfig(s.req.Plugin, s.config))
	if err != nil {
		return s.left(fmt.Errorf("could not parse plugin config: %w", err))
	}
	return s.right()
}

func createResponse(s *state) fp.Either[*al_proto.PluginStartResponse] {
	resp := &al_proto.PluginStartResponse{
		Env: map[string]string{
			"PM_API_TOKEN_ID":     s.pveToken.Data.TokenId,
			"PM_API_TOKEN_SECRET": s.pveToken.Data.TokenSecret,
			"PM_API_URL":          fmt.Sprintf("%s/api2/json", s.config.PveBaseUrl),
		},
	}
	return fp.Right(resp)
}

type Plugin struct {
}

func (self Plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (*al_proto.PluginStartResponse, error) {
	return fp.Get(fp.Pipe7(
		parseConfig,
		createVaultToken,
		createOIDCRequest,
		loginToOIDCProvider,
		createProxmoxTicket,
		createProxmoxToken,
		createResponse,
	)(&state{
		ctx:    ctx,
		req:    req,
		config: &pve_login_proto.Config{},
	}))
}

func main() {
	fp.GetOrExit[struct{}](os.Stderr, 1)(al_plugin.Serve(context.Background(), os.Stdin, os.Stdout, Plugin{}))
}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
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
	"github.com/google/uuid"
)

var pveBaseUrl = flag.String("pve_base_url", "", "Base url")
var pveRedirectUrl = flag.String("pve_redirect_url", "", "Proxmox redirect url")
var pveRealm = flag.String("pve_realm", "", "Proxmox realm")
var vaultConn = flag.String("vault_conn", "", "Vault connection")
var vaultAuth = flag.String("vault_auth", "", "Vault auth")
var logger = log.New(os.Stderr, "com.alwaldend.src.tools.vault.pve_login ", log.Flags())

func oidcData() (string, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/api2/json/access/openid/auth-url?realm=%s&redirect-url=%s", *pveBaseUrl, *pveRealm, *pveRedirectUrl),
		strings.NewReader("{}"),
	)
	if err != nil {
		return "", fmt.Errorf("could not create request: %w", err)
	}
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
		return "", fmt.Errorf("could not read response body")
	}
	data := struct {
		Data string `json:"data"`
	}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal response: %w", err)
	}
	return data.Data, nil
}

func oidcLogin(oidcUrl string, vaultToken string) (string, string, error) {
	regex, err := regexp.Compile("^https://([^/]+)/ui/vault/identity/oidc/provider/([^/]+)/authorize?(.+)$")
	if err != nil {
		return "", "", fmt.Errorf("could not compile regex: %w", err)
	}
	parts := regex.FindAllStringSubmatch(oidcUrl, -1)
	urlParsed, err := url.Parse(oidcUrl)
	if err != nil {
		return "", "", fmt.Errorf("could not parse url: %w", err)
	}
	reqData := map[string]string{}
	for key, valueSlice := range urlParsed.Query() {
		for _, value := range valueSlice {
			reqData[key] = value
		}
	}
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return "", "", fmt.Errorf("could not marshal data: %w", err)
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://%s/v1/identity/oidc/provider/%s/authorize", parts[0][1], parts[0][2]),
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return "", "", fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", vaultToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("could not execute the request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("invalid response code: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("could not read the body: %w", err)
	}
	respData := struct {
		Code  string `json:"code"`
		State string `json:"state"`
	}{}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return "", "", fmt.Errorf("could not unmarshal response body: %w", err)
	}
	return respData.Code, respData.State, nil
}

func createTicket(code string, state string) (string, string, string, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/api2/json/access/openid/login", *pveBaseUrl),
		strings.NewReader("{}"),
	)
	if err != nil {
		return "", "", "", fmt.Errorf("could not create request: %w", err)
	}
	query := req.URL.Query()
	query.Add("state", state)
	query.Add("code", code)
	query.Add("redirect-url", *pveRedirectUrl)
	req.URL.RawQuery = query.Encode()
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", "", fmt.Errorf("could not execute the request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", "", "", fmt.Errorf("invalid response code: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", fmt.Errorf("could not read response body: %w", err)
	}
	data := struct {
		Data struct {
			CSRFPreventionToken string
			Username            string `json:"username"`
			Ticket              string `json:"ticket"`
		} `json:"data"`
	}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", "", "", fmt.Errorf("could not unmarshel response body: %w", err)
	}
	return data.Data.Username, data.Data.Ticket, data.Data.CSRFPreventionToken, nil
}

func createToken(username string, ticket string, csrfToken string) (string, string, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"%s/api2/json/access/users/%s/token/%s",
			*pveBaseUrl,
			username,
			fmt.Sprintf("tools-vault-pve-login-%s", uuid.New().String()),
		),
		strings.NewReader("{}"),
	)
	if err != nil {
		return "", "", fmt.Errorf("could not create request: %w", err)
	}
	query := req.URL.Query()
	query.Add("comment", "Automatically created by //tools/vault/pve_login")
	query.Add("expire", fmt.Sprintf("%d", time.Now().Add(time.Hour).Unix()))
	query.Add("privsep", "0")
	req.URL.RawQuery = query.Encode()
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("CSRFPreventionToken", csrfToken)
	req.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: ticket})
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("could not execute the request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("invalid response code: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("could not read response body: %w", err)
	}
	data := struct {
		Data struct {
			TokenId     string `json:"full-tokenid"`
			TokenSecret string `json:"value"`
		} `json:"data"`
	}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", "", fmt.Errorf("could not unmarshal response body: %w", err)
	}
	return data.Data.TokenId, data.Data.TokenSecret, nil
}

type Plugin struct {
}

func (self Plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (res *al_proto.PluginStartResponse, err error) {
	logger.Println("Creating a vault token")
	vault, err := al.NewVault(req.Config)
	if err != nil {
		return nil, fmt.Errorf("could not create vault client: %w", err)
	}
	vault_client, err := vault.Client(ctx, *vaultConn, *vaultAuth)
	if err != nil {
		return nil, fmt.Errorf("could not authorize vault client: %w", err)
	}
	logger.Println("Created a vault token")
	logger.Println("Requesting OIDC url")
	oidcUrl, err := oidcData()
	if err != nil {
		return nil, fmt.Errorf("could not create oidc url: %w", err)
	}
	logger.Println("Got OIDC url")
	logger.Println("Authenticating to the OIDC provider")
	code, state, err := oidcLogin(oidcUrl, vault_client.Token())
	if err != nil {
		return nil, fmt.Errorf("could not authenticate to the oidc url: %w", err)
	}
	logger.Println("Authenticated to the OIDC url")
	logger.Println("Creating a ticket")
	username, ticket, csrfToken, err := createTicket(code, state)
	if err != nil {
		return nil, fmt.Errorf("could not create Proxmox ticket: %w", err)
	}
	logger.Println("Created a ticket")
	logger.Println("Creating a token")
	tokenId, tokenSecret, err := createToken(username, ticket, csrfToken)
	if err != nil {
		return nil, fmt.Errorf("could not create token: %w", err)
	}
	logger.Println("Created a token")

	res = &al_proto.PluginStartResponse{
		Env: map[string]string{
			"PM_API_TOKEN_ID":     tokenId,
			"PM_API_TOKEN_SECRET": tokenSecret,
			"PM_API_URL":          fmt.Sprintf("%s/api2/json", *pveBaseUrl),
		},
	}
	return res, nil
}

func main() {
	flag.Parse()
	err := al_plugin.Serve(context.Background(), os.Stdin, os.Stdout, Plugin{})
	if err != nil {
		logger.Printf("failed to run plugin: %s", err.Error())
		os.Exit(1)
	}
}

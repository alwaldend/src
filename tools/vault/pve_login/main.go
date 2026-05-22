package main

import (
	"bytes"
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

	"git.alwaldend.com/alwaldend/src/tools/al/pkg/al"
	"github.com/google/uuid"
)

var pveBaseUrl = flag.String("pve_base_url", "https://host1.pve1.dc1.alwaldend.com:8006", "Base url")
var pveRedirectUrl = flag.String("pve_redirect_url", "https://host1.pve1.dc1.alwaldend.com:8006", "Proxmox redirect url")
var pveRealm = flag.String("pve_realm", "src_infra_dc1_vault", "Proxmox realm")

func main() {
	flag.Parse()

	req1 := al.Must(
		http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("%s/api2/json/access/openid/auth-url?realm=%s&redirect-url=%s", *pveBaseUrl, *pveRealm, *pveRedirectUrl),
			strings.NewReader("{}"),
		),
	)
	req1.Header.Add("Content-Type", "application/json")
	resp1 := al.Must(http.DefaultClient.Do(req1))
	defer resp1.Body.Close()
	body1 := al.Must(io.ReadAll(resp1.Body))
	log.Println("Created OIDC url", string(body1))
	data1 := struct {
		Data string `json:"data"`
	}{}
	al.Check(json.Unmarshal(body1, &data1))

	regex := regexp.MustCompile("^https://([^/]+)/ui/vault/identity/oidc/provider/([^/]+)/authorize?(.+)$")
	req2Parts := regex.FindAllStringSubmatch(data1.Data, -1)
	req2BodyData := map[string]string{}
	for key, valueSlice := range al.Must(url.Parse(data1.Data)).Query() {
		for _, value := range valueSlice {
			req2BodyData[key] = value
		}
	}
	req2Body := al.Must(json.Marshal(req2BodyData))
	req2 := al.Must(
		http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("https://%s/v1/identity/oidc/provider/%s/authorize", req2Parts[0][1], req2Parts[0][2]),
			bytes.NewBuffer(req2Body),
		),
	)
	req2.Header.Add("Content-Type", "application/json")
	req2.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("VAULT_TOKEN")))
	resp2 := al.Must(http.DefaultClient.Do(req2))
	defer resp2.Body.Close()
	body2 := al.Must(io.ReadAll(resp2.Body))
	log.Println("Authorized using the OIDC url", string(body2))
	data2 := struct {
		Code  string `json:"code"`
		State string `json:"state"`
	}{}
	al.Check(json.Unmarshal(body2, &data2))

	req3 := al.Must(
		http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("%s/api2/json/access/openid/login", *pveBaseUrl),
			strings.NewReader("{}"),
		),
	)
	req3Query := req3.URL.Query()
	req3Query.Add("state", data2.State)
	req3Query.Add("code", data2.Code)
	req3Query.Add("redirect-url", *pveRedirectUrl)
	req3.URL.RawQuery = req3Query.Encode()
	req3.Header.Add("Content-Type", "application/json")
	resp3 := al.Must(http.DefaultClient.Do(req3))
	defer resp3.Body.Close()
	body3 := al.Must(io.ReadAll(resp3.Body))
	log.Println("Exchanged OIDC token for a ticket", string(body3))
	data3 := struct {
		Data struct {
			CSRFPreventionToken string
			Username            string `json:"username"`
			Ticket              string `json:"ticket"`
		} `json:"data"`
	}{}
	al.Check(json.Unmarshal(body3, &data3))

	req4 := al.Must(
		http.NewRequest(
			http.MethodPost,
			fmt.Sprintf(
				"%s/api2/json/access/users/%s/token/%s",
				*pveBaseUrl,
				data3.Data.Username,
				fmt.Sprintf("tools-vault-pve-login-%s", uuid.New().String()),
			),
			strings.NewReader("{}"),
		),
	)
	req4Query := req4.URL.Query()
	req4Query.Add("comment", "Automatically created by //tools/vault/pve_login")
	req4Query.Add("expire", fmt.Sprintf("%d", time.Now().Add(time.Hour).Unix()))
	req4Query.Add("privsep", "0")
	req4.URL.RawQuery = req4Query.Encode()
	req4.Header.Add("Content-Type", "application/json")
	req4.Header.Add("Content-Type", "application/json")
	req4.Header.Add("CSRFPreventionToken", data3.Data.CSRFPreventionToken)
	req4.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: data3.Data.Ticket})
	resp4 := al.Must(http.DefaultClient.Do(req4))
	defer resp4.Body.Close()
	body4 := al.Must(io.ReadAll(resp4.Body))
	log.Println("Created a token", string(body4))
	data4 := struct {
		Data struct {
			TokenId     string `json:"full-tokenid"`
			TokenSecret string `json:"value"`
		} `json:"data"`
	}{}
	al.Check(json.Unmarshal(body4, &data4))
	if resp4.StatusCode != 200 {
		panic(fmt.Sprintf("invalid status code: %d", resp4.StatusCode))
	}

	os.Stdout.WriteString(fmt.Sprintf("PM_API_TOKEN_ID=%s", data4.Data.TokenId))
	os.Stdout.WriteString("\n")
	os.Stdout.WriteString(fmt.Sprintf("PM_API_TOKEN_SECRET=%s", data4.Data.TokenSecret))
	os.Stdout.WriteString("\n")
	os.Stdout.WriteString(fmt.Sprintf("PM_API_URL=%s/api2/json", *pveBaseUrl))
	os.Stdout.WriteString("\n")
}

package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
)

func TestIssuerCredentialsStayIsolated(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "failed_target_login"}[fail], func(t *testing.T) {
			issued, login, destroyed := 0, 0, 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var body map[string]any
				if json.NewDecoder(r.Body).Decode(&body) != nil {
					t.Error("invalid request")
				}
				switch r.URL.Path {
				case "/v1/auth/approle/role/src_infra_flux_git/secret-id":
					issued++
					if r.Header.Get("X-Vault-Token") != "issuer-token" || body["num_uses"] != float64(1) {
						t.Error("SecretID issuance scope wrong")
					}
					_, _ = w.Write([]byte(`{"data":{"secret_id":"private-secret","secret_id_accessor":"private-accessor"}}`))
				case "/v1/auth/approle/login":
					login++
					if r.Header.Get("X-Vault-Token") != "" || body["role_id"] != "src_infra_flux_git" || body["secret_id"] != "private-secret" {
						t.Error("target login crossed credential boundary")
					}
					if fail {
						w.WriteHeader(403)
						_, _ = w.Write([]byte(`{"errors":["private-secret"]}`))
						return
					}
					_, _ = w.Write([]byte(`{"auth":{"client_token":"target-token"}}`))
				case "/v1/auth/approle/role/src_infra_flux_git/secret-id-accessor/destroy":
					destroyed++
					if r.Header.Get("X-Vault-Token") != "issuer-token" || body["secret_id_accessor"] != "private-accessor" {
						t.Error("SecretID cleanup exceeded owned scope")
					}
					w.WriteHeader(204)
				default:
					t.Error("unexpected Vault operation")
					w.WriteHeader(404)
				}
			}))
			defer server.Close()
			issuer, err := api.NewClient(&api.Config{Address: server.URL, HttpClient: server.Client()})
			if err != nil {
				t.Fatal(err)
			}
			issuer.SetToken("issuer-token")
			client, err := issueRole(context.Background(), issuer, "src_infra_flux_git")
			if fail {
				if err == nil || client != nil || destroyed != 1 || strings.Contains(err.Error(), "private") {
					t.Fatal("target failure did not clean up safely")
				}
			} else if err != nil || client.Token() != "target-token" || destroyed != 0 {
				t.Fatal("target login failed")
			}
			if issuer.Token() != "issuer-token" || issued != 1 || login != 1 {
				t.Fatal("issuer mutated or repeated issuance")
			}
		})
	}
}

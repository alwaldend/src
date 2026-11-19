package al

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"github.com/hashicorp/vault/api"
)

func testVaultConfig(address string, auth *al_proto.VaultAuth) *al_proto.Config {
	return &al_proto.Config{
		VaultConn: []*al_proto.VaultConn{{Name: "default", Config: &al_proto.VaultConfig{Address: address}}},
		VaultAuth: []*al_proto.VaultAuth{auth},
	}
}

func TestNoAuthDoesNotReadOrForwardCredentials(t *testing.T) {
	// A nonexistent HOME makes accidental token-helper use fail.
	t.Setenv("HOME", filepath.Join(t.TempDir(), "missing"))
	t.Setenv("VAULT_TOKEN", "synthetic-ambient-token")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != "" {
			t.Error("no_auth sent an authentication token")
		}
		io.WriteString(w, `{"data":{}}`)
	}))
	defer server.Close()
	store := NewVault(testVaultConfig(server.URL, &al_proto.VaultAuth{Name: "default", NoAuth: true}))
	item, err := store.Client(context.Background(), "", "").Get()
	if err != nil {
		t.Fatal(err)
	}
	if item.Helper != nil || item.Client.Token() != "" {
		t.Fatal("no_auth loaded credentials")
	}
	if _, err := item.Client.Logical().ReadWithContext(context.Background(), "fixture"); err != nil {
		t.Fatal(err)
	}
	if err := store.Stop(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Client(context.Background(), "", "").Get(); err == nil {
		t.Fatal("closed store accepted authentication")
	}
}

func TestVaultOwnedCredentialCleanup(t *testing.T) {
	// Vault's token helper caches home-directory discovery process-wide.
	home := t.TempDir()
	t.Setenv("HOME", home)
	tokenPath := filepath.Join(home, ".vault-token")
	if err := os.WriteFile(tokenPath, []byte("bootstrap-token"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, failLogin := range []bool{false, true} {
		t.Run(map[bool]string{false: "successful-login", true: "failed-login"}[failLogin], func(t *testing.T) {
			var issued, revoked, destroyed atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/v1/auth/approle/role/fixture/secret-id":
					if r.Header.Get("X-Vault-Token") != "bootstrap-token" {
						t.Error("SecretID creation lost bootstrap authentication")
					}
					var body map[string]any
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Error(err)
					}
					if body["num_uses"] != float64(1) {
						t.Error("SecretID must be single-use")
					}
					issued.Add(1)
					io.WriteString(w, `{"data":{"secret_id":"generated-secret","secret_id_accessor":"owned-accessor"}}`)
				case "/v1/auth/approle/login":
					if failLogin {
						w.WriteHeader(http.StatusForbidden)
						io.WriteString(w, `{"errors":["denied"]}`)
						return
					}
					io.WriteString(w, `{"auth":{"client_token":"role-token","lease_duration":600}}`)
				case "/v1/auth/token/revoke-self":
					if r.Header.Get("X-Vault-Token") != "role-token" {
						t.Error("revoking a token not owned by this invocation")
					}
					revoked.Add(1)
					w.WriteHeader(http.StatusNoContent)
				case "/v1/auth/approle/role/fixture/secret-id-accessor/destroy":
					if r.Header.Get("X-Vault-Token") != "bootstrap-token" {
						t.Error("SecretID cleanup lost bootstrap authentication")
					}
					destroyed.Add(1)
					w.WriteHeader(http.StatusNoContent)
				default:
					t.Errorf("unexpected Vault path: %s", r.URL.Path)
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer server.Close()
			store := NewVault(testVaultConfig(server.URL, &al_proto.VaultAuth{Name: "default", Approle: &al_proto.VaultAuthApprole{Name: "fixture"}}))
			_, err := store.Client(context.Background(), "", "").Get()
			if (err != nil) != failLogin {
				t.Fatalf("login error = %v", err)
			}
			if !failLogin {
				if _, err := store.Client(context.Background(), "", "").Get(); err != nil {
					t.Fatal(err)
				}
			}
			if err := store.Stop(context.Background()); err != nil {
				t.Fatal(err)
			}
			if err := store.Stop(context.Background()); err != nil {
				t.Fatal(err)
			}
			if issued.Load() != 1 {
				t.Fatal("cache created duplicate credentials")
			}
			if failLogin && (destroyed.Load() != 1 || revoked.Load() != 0) {
				t.Fatal("failed login did not clean only the SecretID")
			}
			if !failLogin && (revoked.Load() != 1 || destroyed.Load() != 0) {
				t.Fatal("successful login did not revoke exactly the issued token")
			}
			data, err := os.ReadFile(tokenPath)
			if err != nil || string(data) != "bootstrap-token" {
				t.Fatal("user token helper was changed")
			}
		})
	}
}

func TestSanitizeVaultError(t *testing.T) {
	err := SanitizeVaultError(&api.ResponseError{StatusCode: 403, Errors: []string{"synthetic-secret"}, URL: "https://fixture.invalid/secret"})
	if strings.Contains(err.Error(), "synthetic-secret") || strings.Contains(err.Error(), "fixture.invalid") || !strings.Contains(err.Error(), "403") {
		t.Fatal("unsafe or unhelpful Vault error")
	}
}

func TestOIDCOriginAndRedirectProtection(t *testing.T) {
	var trustedCalls, foreignCalls atomic.Int32
	foreign := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { foreignCalls.Add(1) }))
	defer foreign.Close()
	redirect := false
	trusted := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trustedCalls.Add(1)
		if r.Header.Get("Authorization") != "Bearer fixture-token" {
			t.Error("missing expected fixture auth")
		}
		if redirect {
			http.Redirect(w, r, foreign.URL, http.StatusTemporaryRedirect)
			return
		}
		io.WriteString(w, `{"code":"fixture-code","state":"fixture-state"}`)
	}))
	defer trusted.Close()
	client, err := api.NewClient(&api.Config{Address: trusted.URL, HttpClient: trusted.Client()})
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken("fixture-token")
	store := NewVault(nil)
	for _, address := range []string{foreign.URL + "/ui/vault/identity/oidc/provider/fixture/authorize?state=x", trusted.URL + "/wrong-path", strings.Replace(trusted.URL, "https:", "http:", 1) + "/ui/vault/identity/oidc/provider/fixture/authorize"} {
		parsed, err := url.Parse(address)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := store.OidcLoginContext(context.Background(), client, parsed); err == nil {
			t.Fatal("accepted invalid authorization URL")
		}
	}
	if trustedCalls.Load() != 0 || foreignCalls.Load() != 0 {
		t.Fatal("invalid URL caused network access")
	}
	parsed, _ := url.Parse(trusted.URL + "/ui/vault/identity/oidc/provider/fixture/authorize?state=x")
	result, err := store.OidcLoginContext(context.Background(), client, parsed)
	if err != nil || result.Code != "fixture-code" {
		t.Fatalf("valid authorization failed: %v", err)
	}
	redirect = true
	if _, err := store.OidcLoginContext(context.Background(), client, parsed); err == nil {
		t.Fatal("accepted redirect")
	}
	if foreignCalls.Load() != 0 {
		t.Fatal("forwarded Vault token to foreign origin")
	}
}

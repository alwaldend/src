package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	reqUrl "net/url"
	"strings"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/tools/vault/forgejo_login/forgejo_login_proto"
	"google.golang.org/protobuf/proto"
)

func pluginData(t *testing.T, value any) *al_proto.Json {
	t.Helper()
	data, err := al.ToPbJson(value).Get()
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestPluginConfigInheritsDefaultsAndSelectsNamedAuthentication(t *testing.T) {
	defaults := map[string]any{
		"forgejo_url":        "https://forgejo.example",
		"forgejo_oauth_name": "vault",
		"vault_conn":         "primary",
		"vault_auth":         "default",
	}
	for _, test := range []struct {
		name     string
		calls    []*al_proto.PluginCall
		wantAuth string
	}{
		{name: "no selected calls", wantAuth: "default"},
		{
			name: "empty calls preserve defaults",
			calls: []*al_proto.PluginCall{
				{Name: "first"},
				{Name: "second", Data: pluginData(t, map[string]any{})},
			},
			wantAuth: "default",
		},
		{
			name: "named authentication inherits service configuration",
			calls: []*al_proto.PluginCall{
				{Name: "ordinary"},
				{Name: "registration", Data: pluginData(t, map[string]any{"vault_auth": "forgejo_registration"})},
			},
			wantAuth: "forgejo_registration",
		},
		{
			name: "empty fields preserve defaults",
			calls: []*al_proto.PluginCall{
				{Data: pluginData(t, map[string]any{"vault_auth": "", "forgejo_url": ""})},
			},
			wantAuth: "default",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			data := pluginData(t, defaults)
			original := proto.Clone(data)
			got, err := parsePluginConfig(&al_proto.PluginConfig{Data: data, Calls: test.calls})
			if err != nil {
				t.Fatal(err)
			}
			want := &forgejo_login_proto.Config{
				ForgejoUrl:       "https://forgejo.example",
				ForgejoOauthName: "vault",
				VaultConn:        "primary",
				VaultAuth:        test.wantAuth,
			}
			if !proto.Equal(got, want) {
				t.Fatalf("configuration = %v, want %v", got, want)
			}
			if !proto.Equal(data, original) {
				t.Fatal("plugin defaults were mutated")
			}
		})
	}
}

func TestPluginConfigRejectsAmbiguousOrMalformedData(t *testing.T) {
	for _, test := range []struct {
		name   string
		plugin *al_proto.PluginConfig
	}{
		{name: "missing plugin"},
		{
			name: "malformed defaults",
			plugin: &al_proto.PluginConfig{
				Data: pluginData(t, "private-default-data"),
			},
		},
		{
			name: "malformed call",
			plugin: &al_proto.PluginConfig{Calls: []*al_proto.PluginCall{
				{Data: pluginData(t, map[string]any{"vault_auth": map[string]any{"private": true}})},
			}},
		},
		{
			name: "unknown call field",
			plugin: &al_proto.PluginConfig{Calls: []*al_proto.PluginCall{
				{Data: pluginData(t, map[string]any{"private_unknown_field": "value"})},
			}},
		},
		{
			name: "ambiguous selected identities",
			plugin: &al_proto.PluginConfig{Calls: []*al_proto.PluginCall{
				{Data: pluginData(t, map[string]any{"vault_auth": "first"})},
				{Data: pluginData(t, map[string]any{"vault_auth": "second"})},
			}},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			config, err := parsePluginConfig(test.plugin)
			if err == nil || config != nil {
				t.Fatalf("invalid configuration accepted: %v", config)
			}
			if strings.Contains(err.Error(), "private") {
				t.Fatal("configuration values disclosed in error")
			}
		})
	}
}

type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func response(r *http.Request, body string) *http.Response {
	return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: r}
}

func TestDeleteVerifiesAbsence(t *testing.T) {
	for _, removed := range []bool{false, true} {
		t.Run(fmt.Sprint(removed), func(t *testing.T) {
			l, err := newLogin(nil, &forgejo_login_proto.Config{ForgejoUrl: "https://forgejo.example"})
			if err != nil {
				t.Fatal(err)
			}
			posted := false
			l.client.Transport = transportFunc(func(r *http.Request) (*http.Response, error) {
				if r.Method == http.MethodPost {
					posted = true
					return response(r, ""), nil
				}
				if posted && removed {
					return response(r, tokenPage("")), nil
				}
				return response(r, tokenPage(tokenRow("owned-token", "7"))), nil
			})
			err = l.deleteForgejoToken(context.Background(), "owned-token")
			if !posted {
				t.Fatal("deletion never attempted")
			}
			if (err == nil) != removed {
				t.Fatalf("removed=%v, error=%v", removed, err)
			}
		})
	}
}

func tokenRow(name, id string) string {
	return `<div class="flex-item"><div class="flex-item-main"><details><summary><span class="flex-item-title">` + name + `</span></summary></details></div><div class="flex-item-trailing"><button data-modal-id="regenerate-token" data-url="/user/settings/applications/tokens/regenerate" data-id="99"></button><button data-modal-id="delete-token" data-url="/user/settings/applications/tokens/delete" data-id="` + id + `"></button></div></div>`
}

func tokenPage(rows string) string {
	return `<html><body><div class="page-content user settings applications"><a href="/user/settings/applications/tokens/new">New token</a><div class="flex-list">` + rows + `</div></div></body></html>`
}

func TestTokenPageSelectsOnlyOwnedToken(t *testing.T) {
	id, err := parseTokenPage(strings.NewReader(tokenPage(tokenRow("unrelated", "1")+tokenRow("owned-token", "7"))), "/user/settings/applications", "owned-token")
	if err != nil || id != "7" {
		t.Fatalf("id=%q err=%v", id, err)
	}
}

func TestTokenPageFailsClosed(t *testing.T) {
	for _, body := range []string{`<html><form>Login</form></html>`, tokenPage(tokenRow("owned-token", "bad")), tokenPage(tokenRow("owned-token", "1") + tokenRow("owned-token", "2")), tokenPage("") + strings.Repeat(" ", 4<<20)} {
		if _, err := parseTokenPage(strings.NewReader(body), "/user/settings/applications", "owned-token"); err == nil {
			t.Fatal("invalid page accepted")
		}
	}
}

func TestTokenCookieErrorIsRedacted(t *testing.T) {
	l, _ := newLogin(nil, &forgejo_login_proto.Config{ForgejoUrl: "https://forgejo.example"})
	l.client.Transport = transportFunc(func(r *http.Request) (*http.Response, error) {
		l.cookies = []*http.Cookie{{Name: "flash", Value: "private-cookie-%zz"}}
		return response(r, "private-body"), nil
	})
	_, err := l.createForgejoToken(context.Background(), "owned-token")
	if err == nil || strings.Contains(err.Error(), "private") {
		t.Fatalf("unsafe cookie error: %v", err)
	}
}

func TestCleanupAfterUnusableCreationResponse(t *testing.T) {
	l, _ := newLogin(nil, &forgejo_login_proto.Config{ForgejoUrl: "https://forgejo.example"})
	created, deleted := false, false
	l.client.Transport = transportFunc(func(r *http.Request) (*http.Response, error) {
		if strings.HasSuffix(r.URL.Path, "/new/") {
			created = true
			return response(r, ""), nil
		}
		if r.Method == http.MethodPost {
			deleted = true
			return response(r, ""), nil
		}
		if created && !deleted {
			return response(r, tokenPage(tokenRow("owned-token", "7"))), nil
		}
		return response(r, tokenPage("")), nil
	})
	if _, err := l.createForgejoToken(context.Background(), "owned-token"); err == nil {
		t.Fatal("missing cookie accepted")
	}
	if err := l.deleteForgejoToken(context.Background(), "owned-token"); err != nil {
		t.Fatal(err)
	}
	if !deleted {
		t.Fatal("token was left behind")
	}
}

func TestSessionRedirectRejectsExternalOrigin(t *testing.T) {
	l, _ := newLogin(nil, &forgejo_login_proto.Config{ForgejoUrl: "https://forgejo.example"})
	calls := 0
	l.client.Transport = transportFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		resp := response(r, "")
		resp.StatusCode = http.StatusFound
		resp.Header.Set("Location", "https://external.example/private-code")
		return resp, nil
	})
	_, err := l.getTokenId(context.Background(), "owned-token")
	if err == nil {
		t.Fatal("external redirect allowed")
	}
	if calls != 1 {
		t.Fatalf("made %d requests", calls)
	}
	if strings.Contains(err.Error(), "private-code") {
		t.Fatal("redirect URL disclosed")
	}
}

func TestLogoutVerifiesRetainedSession(t *testing.T) {
	for _, invalidated := range []bool{false, true} {
		t.Run(fmt.Sprint(invalidated), func(t *testing.T) {
			l, _ := newLogin(nil, &forgejo_login_proto.Config{ForgejoUrl: "https://forgejo.example"})
			endpoint, _ := reqUrl.Parse("https://forgejo.example")
			l.client.Jar.SetCookies(endpoint, []*http.Cookie{{Name: "session", Value: "placeholder-session"}})
			posted := false
			l.client.Transport = transportFunc(func(r *http.Request) (*http.Response, error) {
				cookie, err := r.Cookie("session")
				if err != nil || cookie.Value != "placeholder-session" {
					t.Fatal("did not verify original session")
				}
				if r.Method == http.MethodPost {
					posted = true
					return response(r, `{"redirect":"/"}`), nil
				}
				if posted && invalidated {
					resp := response(r, "")
					resp.StatusCode = http.StatusSeeOther
					resp.Header.Set("Location", "/user/login")
					return resp, nil
				}
				return response(r, tokenPage("")), nil
			})
			err := l.logout(context.Background())
			if !posted || (err == nil) != invalidated {
				t.Fatalf("posted=%v invalidated=%v error=%v", posted, invalidated, err)
			}
		})
	}
}

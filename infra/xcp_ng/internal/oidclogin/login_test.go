package oidclogin

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"testing"
)

func TestAuthorizeDestination(t *testing.T) {
	vault, _ := url.Parse("https://vault.test")
	discovery, _ := url.Parse("https://vault.test/v1/identity/oidc/provider/xcp")
	callback, _ := url.Parse("https://xo.test/signin/oidc/callback")
	query := url.Values{"redirect_uri": {callback.String()}, "state": {"state"}, "response_type": {"code"}, "nonce": {"nonce"}}
	valid := "https://vault.test/ui/vault/identity/oidc/provider/xcp/authorize?" + query.Encode()
	target, params, err := authorizeURL(valid, vault, discovery, callback)
	if err != nil || target.String() != "https://vault.test/v1/identity/oidc/provider/xcp/authorize" || params["nonce"] != "nonce" {
		t.Fatalf("valid redirect rejected: %v", err)
	}
	for _, bad := range []string{strings.Replace(valid, "vault.test", "evil.test", 1), strings.Replace(valid, "/xcp/", "/other/", 1), valid + "&state=duplicate", strings.Replace(valid, "https://", "http://", 1), strings.Replace(valid, "xo.test", "evil.test", 1)} {
		if _, _, err := authorizeURL(bad, vault, discovery, callback); err == nil {
			t.Errorf("accepted unsafe redirect %s", bad)
		}
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestBootstrapLoginBoundaries(t *testing.T) {
	for _, test := range []struct {
		name     string
		badState bool
		status   int
		location string
		wantErr  string
	}{
		{name: "found", status: 302, location: "/"},
		{name: "see_other", status: 303, location: "/"},
		{name: "state_mismatch", badState: true, status: 303, location: "/", wantErr: "invalid Vault OIDC authorization result"},
		{name: "callback_error", status: 500, wantErr: "XO OIDC callback did not complete (HTTP 500)"},
		{name: "html_is_not_success", status: 200, wantErr: "XO OIDC callback did not complete (HTTP 200)"},
		{name: "rejected_login", status: 303, location: "/signin", wantErr: "XO OIDC callback returned an unexpected destination"},
	} {
		t.Run(test.name, func(t *testing.T) {
			xo, _ := url.Parse("https://xo.test")
			vault, _ := url.Parse("https://vault.test")
			discovery, _ := url.Parse("https://vault.test/v1/identity/oidc/provider/xcp")
			query := url.Values{"redirect_uri": {"https://xo.test/signin/oidc/callback"}, "state": {"state"}, "response_type": {"code"}}
			redirect := "https://vault.test/ui/vault/identity/oidc/provider/xcp/authorize?" + query.Encode()
			jar, _ := cookiejar.New(nil)
			callbacks := 0
			xoClient := &http.Client{Jar: jar, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }, Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				if req.URL.Host != "xo.test" || req.Header.Get("X-Vault-Token") != "" {
					t.Fatal("credential destination violation")
				}
				res := &http.Response{StatusCode: 302, Header: http.Header{}, Body: io.NopCloser(strings.NewReader("")), Request: req}
				if req.URL.Path == "/signin/oidc" {
					res.Header.Set("Location", redirect)
					res.Header.Set("Set-Cookie", "session=opaque; Path=/; Secure")
				} else {
					callbacks++
					cookie, err := req.Cookie("session")
					if err != nil || cookie.Value != "opaque" {
						t.Fatal("session cookie missing")
					}
					if req.URL.Query().Get("code") != "private-code" || req.URL.Query().Get("state") != "state" {
						t.Fatal("invalid callback")
					}
					res.StatusCode = test.status
					res.Header.Set("Location", test.location)
					res.Body = io.NopCloser(strings.NewReader("private-response private-code private-token"))
				}
				return res, nil
			})}
			vaultClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				if req.URL.String() != "https://vault.test/v1/identity/oidc/provider/xcp/authorize" || req.Header.Get("X-Vault-Token") != "private-token" || req.Header.Get("Cookie") != "" {
					t.Fatal("Vault request boundary violated")
				}
				state := "state"
				if test.badState {
					state = "wrong"
				}
				return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"code":"private-code","state":"` + state + `"}`)), Header: http.Header{}, Request: req}, nil
			})}
			err := Bootstrap(xoClient, vaultClient, xo, vault, discovery, "private-token")
			if test.wantErr == "" {
				if err != nil {
					t.Fatalf("bootstrap failed: %v", err)
				}
			} else if err == nil || err.Error() != test.wantErr {
				t.Fatalf("unexpected bootstrap error: %v", err)
			}
			wantCallbacks := 1
			if test.badState {
				wantCallbacks = 0
			}
			if callbacks != wantCallbacks {
				t.Fatalf("callback count = %d, want %d", callbacks, wantCallbacks)
			}
		})
	}
}

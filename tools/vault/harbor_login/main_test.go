package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
	"git.alwaldend.com/alwaldend/src/tools/vault/harbor_login/harbor_login_proto"
)

// TestLogoutVerifiesOriginalSessionAndDoesNotFollowProvider verifies local logout while preserving the identity-provider session.
func TestLogoutVerifiesOriginalSessionAndDoesNotFollowProvider(t *testing.T) {
	for _, invalidated := range []bool{true, false} {
		t.Run(fmt.Sprint(invalidated), func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				cookie, err := r.Cookie("sid")
				if err != nil || cookie.Value != "test-session" {
					t.Error("original session was not sent")
				}
				switch r.URL.Path {
				case "/c/oidc/logout":
					http.SetCookie(w, &http.Cookie{Name: "sid", MaxAge: -1})
					w.Header().Set("Location", "http://127.0.0.1:1/provider-logout")
					w.WriteHeader(http.StatusFound)
				case "/api/v2.0/users/current":
					if invalidated {
						w.WriteHeader(http.StatusUnauthorized)
					} else {
						w.WriteHeader(http.StatusOK)
					}
				default:
					t.Errorf("unexpected request %s", r.URL.Path)
				}
			}))
			defer server.Close()
			login, err := newLogin(&harbor_login_proto.Config{HarborUrl: server.URL})
			if err != nil {
				t.Fatal(err)
			}
			login.session = "test-session"
			plugin := &plugin{}
			if err := plugin.lc.AddState(lifecycle.StateStarted, lifecycle.StoppableFunc(login.logout)); err != nil {
				t.Fatal(err)
			}
			err = plugin.Stop(context.Background())
			if (err == nil) != invalidated {
				t.Fatalf("unexpected cleanup result: %v", err)
			}
			if calls != 2 {
				t.Fatalf("requests = %d", calls)
			}
		})
	}
}

// TestCancelledLogout checks that cancellation remains identifiable through session cleanup.
func TestCancelledLogout(t *testing.T) {
	login, err := newLogin(&harbor_login_proto.Config{HarborUrl: "http://127.0.0.1:1"})
	if err != nil {
		t.Fatal(err)
	}
	login.session = "test-session"
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := login.logout(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("logout error = %v, want cancellation", err)
	}
}

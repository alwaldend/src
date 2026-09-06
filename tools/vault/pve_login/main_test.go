package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
	"git.alwaldend.com/alwaldend/src/tools/vault/pve_login/pve_login_proto"
)

func TestTokenCleanupAfterMalformedCreationResponse(t *testing.T) {
	var created, deleted bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api2/json/access/users/example@oidc/token/task" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		c, err := r.Cookie("PVEAuthCookie")
		if err != nil || c.Value != "test-ticket" || r.Header.Get("CSRFPreventionToken") != "test-csrf" {
			t.Error("missing authentication")
		}
		switch r.Method {
		case http.MethodPost:
			created = true
			fmt.Fprint(w, `{invalid`)
		case http.MethodDelete:
			deleted = true
			fmt.Fprint(w, `{"data":null}`)
		default:
			t.Errorf("unexpected method %s", r.Method)
		}
	}))
	defer server.Close()
	config := &pve_login_proto.Config{PveBaseUrl: server.URL}
	ticket := &pveTicket{}
	ticket.Data.Username = "example@oidc"
	ticket.Data.Ticket = "test-ticket"
	ticket.Data.CSRFPreventionToken = "test-csrf"
	plugin := &Plugin{}
	plugin.lc.AddState(lifecycle.StateStarted, lifecycle.StoppableFunc(func(ctx context.Context) error { return deleteProxmoxToken(ctx, config, ticket, "task") }))
	if _, err := createProxmoxToken(context.Background(), config, ticket, "task"); err == nil {
		t.Fatal("malformed response accepted")
	}
	if err := plugin.Stop(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !created || !deleted {
		t.Fatal("token not created and deleted")
	}
}

func TestDeleteDoesNotFollowRedirectOrDiscloseResponse(t *testing.T) {
	hits := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.Header().Set("Location", "/login?token=test-sensitive")
		w.WriteHeader(http.StatusFound)
		fmt.Fprint(w, "test-sensitive")
	}))
	defer server.Close()
	ticket := &pveTicket{}
	ticket.Data.Username = "example@oidc"
	err := deleteProxmoxToken(context.Background(), &pve_login_proto.Config{PveBaseUrl: server.URL}, ticket, "task")
	if err == nil || strings.Contains(err.Error(), "test-sensitive") || hits != 1 {
		t.Fatalf("unexpected deletion result: %v, requests %d", err, hits)
	}
}

func TestCancelledDelete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ticket := &pveTicket{}
	if err := deleteProxmoxToken(ctx, &pve_login_proto.Config{PveBaseUrl: "http://127.0.0.1:1"}, ticket, "task"); err == nil {
		t.Fatal("cancelled cleanup succeeded")
	}
}

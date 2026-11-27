package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
)

// TestServerStopClosesListenerWithoutServeError checks normal shutdown and subsequent connection refusal.
func TestServerStopClosesListenerWithoutServeError(t *testing.T) {
	server := newServer(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
	if err := server.start(t.Context()); err != nil {
		t.Fatal(err)
	}
	address, err := server.address()
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{Timeout: time.Second}
	response, err := client.Get("http://" + address)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	if err := server.stop(ctx); err != nil {
		t.Fatalf("normal shutdown failed: %v", err)
	}
	if response, err := client.Get("http://" + address); err == nil {
		response.Body.Close()
		t.Fatal("listener still open")
	}
}

// TestServerStopExpiredContextForcesConnectionClose checks that a canceled drain closes active connections.
func TestServerStopExpiredContextForcesConnectionClose(t *testing.T) {
	started := make(chan struct{})
	exited := make(chan struct{})
	server := newServer(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-r.Context().Done()
		close(exited)
	})
	if err := server.start(t.Context()); err != nil {
		t.Fatal(err)
	}
	address, err := server.address()
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		client := &http.Client{Timeout: 5 * time.Second}
		response, err := client.Get("http://" + address)
		if err == nil {
			response.Body.Close()
		}
	}()
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("request not started")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if err := server.stop(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Stop = %v", err)
	}
	select {
	case <-exited:
	case <-time.After(time.Second):
		t.Fatal("handler still running")
	}
	<-done
}

// startTestPlugin obtains a Terraform HTTP endpoint using a local Vault fixture.
func startTestPlugin(t *testing.T, vaultAddress string) (*plugin, map[string]string) {
	t.Helper()
	ctx := al.NewCmdCtx(t.Context(), "test ")
	plugin := newPlugin(ctx)
	if err := plugin.Start(ctx.Ctx); err != nil {
		t.Fatal(err)
	}
	data, err := al.ToPbJson(map[string]any{"vault_secret": "fixture", "vault_secret_mount": "fixture"}).Get()
	if err != nil {
		t.Fatal(err)
	}
	response, err := plugin.PluginStart(ctx.Ctx, &al_proto.PluginStartRequest{
		Config: &al_proto.Config{
			VaultConn: map[string]*al_proto.VaultConn{"default": {Config: &al_proto.VaultConfig{Address: vaultAddress}, Tls: &al_proto.VaultTLS{}}},
			VaultAuth: map[string]*al_proto.VaultAuth{"default": {NoAuth: true}},
		},
		Plugin: &al_proto.PluginConfig{Calls: map[string]*al_proto.PluginCall{"fixture": {Data: data}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		stopCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := plugin.Stop(stopCtx); err != nil {
			t.Errorf("stop plugin: %v", err)
		}
	})
	return plugin, response.GetResult().GetEnv()
}

// TestPluginStopsAlreadyRunningBackend checks that plugin shutdown closes a registered backend.
func TestPluginStopsAlreadyRunningBackend(t *testing.T) {
	plugin, env := startTestPlugin(t, "http://127.0.0.1:1")
	endpoint := env["TF_HTTP_ADDRESS"]
	if endpoint == "" {
		t.Fatal("missing backend address")
	}
	client := &http.Client{Timeout: time.Second}
	// Missing basic auth reaches the handler without contacting Vault.
	before, err := client.Get(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	before.Body.Close()
	stopCtx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	if err := plugin.Stop(stopCtx); err != nil {
		t.Fatal(err)
	}
	if after, err := client.Get(endpoint); err == nil {
		after.Body.Close()
		t.Fatal("backend listener survived plugin shutdown")
	}
}

// TestConcurrentStateInitialization preserves the state another client creates after a missing read.
func TestConcurrentStateInitialization(t *testing.T) {
	var stateExists atomic.Bool
	var createAttempted atomic.Bool
	vault := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/fixture/data/fixture/lock":
			io.WriteString(w, `{"data":{"data":{},"metadata":{"version":1}}}`)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/fixture/data/fixture/state":
			if stateExists.Swap(true) {
				io.WriteString(w, `{"data":{"data":{"state":{"serial":17}},"metadata":{"version":1}}}`)
			} else {
				w.WriteHeader(http.StatusNotFound)
				io.WriteString(w, `{"errors":[]}`)
			}
		case r.Method == http.MethodPut && r.URL.Path == "/v1/fixture/data/fixture/state":
			var payload struct {
				// Options carries the comparison version for the KV v2 write.
				Options map[string]int `json:"options"`
			}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Errorf("decode initialization request: %v", err)
			}
			if cas, present := payload.Options["cas"]; !present || cas != 0 {
				t.Errorf("state initialization used CAS %d (present %t)", cas, present)
			}
			createAttempted.Store(true)
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, `{"errors":["check-and-set conflict"]}`)
		default:
			t.Errorf("unexpected Vault request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer vault.Close()
	_, env := startTestPlugin(t, vault.URL)
	client := &http.Client{Timeout: 5 * time.Second}
	for _, payload := range []string{"   ", ""} {
		request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, env["TF_HTTP_ADDRESS"], strings.NewReader(payload))
		if err != nil {
			t.Fatal(err)
		}
		request.SetBasicAuth(env["TF_HTTP_USERNAME"], env["TF_HTTP_PASSWORD"])
		response, err := client.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		body, readErr := io.ReadAll(response.Body)
		if err := errors.Join(readErr, response.Body.Close()); err != nil {
			t.Fatal(err)
		}
		if payload != "" {
			if response.StatusCode != http.StatusInternalServerError || stateExists.Load() {
				t.Fatal("whitespace-only request was not rejected before Vault access")
			}
			continue
		}
		if response.StatusCode != http.StatusOK || string(body) != `{"serial":17}` {
			t.Fatalf("GET state after competing initialization = %d %q", response.StatusCode, body)
		}
	}
	if !createAttempted.Load() {
		t.Fatal("backend did not attempt a create-only initialization")
	}
}

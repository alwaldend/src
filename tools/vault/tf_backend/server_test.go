package main

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
)

func TestServerStopClosesListenerWithoutServeError(t *testing.T) {
	server := NewServer(nil, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
	if err := server.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	address, err := server.Address()
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{Timeout: time.Second}
	response, err := client.Get("http://" + address)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := server.Stop(ctx); err != nil {
		t.Fatalf("normal shutdown failed: %v", err)
	}
	if response, err := client.Get("http://" + address); err == nil {
		response.Body.Close()
		t.Fatal("listener still open")
	}
}

func TestServerStopExpiredContextForcesConnectionClose(t *testing.T) {
	started := make(chan struct{})
	exited := make(chan struct{})
	server := NewServer(nil, func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-r.Context().Done()
		close(exited)
	})
	if err := server.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	address, err := server.Address()
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		response, err := http.Get("http://" + address)
		if err == nil {
			response.Body.Close()
		}
	}()
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("request not started")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := server.Stop(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Stop = %v", err)
	}
	select {
	case <-exited:
	case <-time.After(time.Second):
		t.Fatal("handler still running")
	}
	<-done
}

func TestPluginStopsAlreadyRunningBackend(t *testing.T) {
	ctx := al.NewCmdCtx(context.Background(), "test ")
	plugin := NewPlugin(ctx)
	if err := plugin.Start(ctx.Ctx); err != nil {
		t.Fatal(err)
	}
	data, err := al.ToPbJson(map[string]any{"vault_secret": "fixture", "vault_secret_mount": "fixture"}).Get()
	if err != nil {
		t.Fatal(err)
	}
	response, err := plugin.PluginStart(ctx.Ctx, &al_proto.PluginStartRequest{
		Config: &al_proto.Config{
			VaultConn: []*al_proto.VaultConn{{Name: "default", Config: &al_proto.VaultConfig{Address: "http://127.0.0.1:1"}, Tls: &al_proto.VaultTLS{}}},
			VaultAuth: []*al_proto.VaultAuth{{Name: "default", NoAuth: true}},
		},
		Plugin: &al_proto.PluginConfig{Calls: []*al_proto.PluginCall{{Name: "fixture", Data: data}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	endpoint := response.Env["TF_HTTP_ADDRESS"]
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
	stopCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := plugin.Stop(stopCtx); err != nil {
		t.Fatal(err)
	}
	if after, err := client.Get(endpoint); err == nil {
		after.Body.Close()
		t.Fatal("backend listener survived plugin shutdown")
	}
}

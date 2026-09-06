package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/net/websocket"
)

func fakeRPC(t *testing.T, answer func(string, json.RawMessage) (any, any)) *rpc {
	t.Helper()
	server := httptest.NewServer(websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			var req struct {
				ID     int             `json:"id"`
				Method string          `json:"method"`
				Params json.RawMessage `json:"params"`
			}
			if websocket.JSON.Receive(ws, &req) != nil {
				return
			}
			result, rpcErr := answer(req.Method, req.Params)
			// Real servers can send asynchronous notifications between responses.
			websocket.JSON.Send(ws, map[string]any{"method": "notification", "params": map[string]any{}})
			if websocket.JSON.Send(ws, map[string]any{"jsonrpc": "2.0", "id": req.ID, "result": result, "error": rpcErr}) != nil {
				return
			}
		}
	}))
	t.Cleanup(server.Close)
	conn, err := websocket.Dial("ws"+strings.TrimPrefix(server.URL, "http"), "", server.URL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	return &rpc{conn: conn}
}

func TestOIDCApply(t *testing.T) {
	for _, loaded := range []bool{false, true} {
		t.Run(map[bool]string{false: "unloaded", true: "loaded"}[loaded], func(t *testing.T) {
			var calls []string
			state := plugin{ID: "auth-oidc", Loaded: loaded}
			r := fakeRPC(t, func(method string, params json.RawMessage) (any, any) {
				calls = append(calls, method)
				switch method {
				case "plugin.get":
					return []plugin{state}, nil
				case "plugin.configure":
					if !strings.Contains(string(params), "private-secret") {
						t.Error("missing configuration")
					}
				case "plugin.enableAutoload":
					state.Autoload = true
				case "plugin.load":
					if state.Loaded {
						t.Error("already loaded")
					}
					state.Loaded = true
				default:
					t.Errorf("unexpected method %s", method)
				}
				return nil, nil
			})
			var out bytes.Buffer
			if err := applyOIDC(r, `{"clientID":"client","clientSecret":"private-secret","discoveryURL":"https://issuer"}`, &out); err != nil {
				t.Fatal(err)
			}
			want := []string{"plugin.get", "plugin.configure", "plugin.enableAutoload"}
			if !loaded {
				want = append(want, "plugin.load")
			}
			want = append(want, "plugin.get")
			if !reflect.DeepEqual(calls, want) {
				t.Fatalf("calls %v, want %v", calls, want)
			}
			if strings.Contains(out.String(), "private-secret") {
				t.Fatal("secret leaked")
			}
		})
	}
}

func TestRPCErrorRedacted(t *testing.T) {
	r := fakeRPC(t, func(string, json.RawMessage) (any, any) {
		return nil, map[string]any{"code": 1, "message": "private-secret", "data": "token"}
	})
	err := r.call("plugin.configure", map[string]any{}, nil)
	if err == nil || strings.Contains(err.Error(), "private-secret") || strings.Contains(err.Error(), "token") {
		t.Fatalf("unsafe error: %v", err)
	}
}

func TestOIDCPostcondition(t *testing.T) {
	r := fakeRPC(t, func(method string, _ json.RawMessage) (any, any) {
		if method == "plugin.get" {
			return []plugin{{ID: "auth-oidc"}}, nil
		}
		return nil, nil
	})
	err := applyOIDC(r, `{"clientID":"client","clientSecret":"secret","discoveryURL":"https://issuer"}`, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "postcondition") {
		t.Fatalf("expected failed postcondition: %v", err)
	}
}

func TestInspectAllowlist(t *testing.T) {
	r := fakeRPC(t, func(method string, _ json.RawMessage) (any, any) {
		switch method {
		case "xo.getAllObjects":
			return map[string]any{"pool": map[string]any{"id": "pool", "type": "pool", "name_label": "test", "password": "private-secret"}, "vm": map[string]any{"id": "vm", "type": "VM", "cloudConfig": "private-secret"}}, nil
		case "plugin.get":
			return []map[string]any{{"id": "auth-oidc", "configuration": map[string]any{"clientSecret": "private-secret"}}}, nil
		case "group.getAll":
			return []map[string]any{{"id": "group", "name": "admins", "users": []string{"personal-user"}}}, nil
		case "server.getAll":
			return []map[string]any{{"id": "server", "host": "192.0.2.1", "password": "private-secret"}}, nil
		}
		return nil, nil
	})
	var out bytes.Buffer
	if err := inspect(r, &out); err != nil {
		t.Fatal(err)
	}
	for _, secret := range []string{"private-secret", "personal-user", "cloudConfig"} {
		if strings.Contains(out.String(), secret) {
			t.Fatalf("unsafe output: %s", out.String())
		}
	}
	if !strings.Contains(out.String(), "192.0.2.1") {
		t.Fatal("missing server host")
	}
}

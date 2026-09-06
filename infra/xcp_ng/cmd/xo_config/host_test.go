package main

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestRegisterHost(t *testing.T) {
	for _, existing := range []bool{false, true} {
		t.Run(map[bool]string{false: "new", true: "existing"}[existing], func(t *testing.T) {
			var calls []string
			var servers []hostServer
			unrelated := hostServer{ID: "other", Host: "192.0.2.2", Enabled: true, Status: "connected", PoolID: "other-pool"}
			servers = append(servers, unrelated)
			if existing {
				servers = append(servers, hostServer{ID: "target", Host: "192.0.2.1", Enabled: true, Status: "connected", PoolID: "pool"})
			}
			r := fakeRPC(t, func(method string, raw json.RawMessage) (any, any) {
				calls = append(calls, method)
				var params map[string]any
				if err := json.Unmarshal(raw, &params); err != nil {
					t.Error(err)
				}
				switch method {
				case "server.getAll":
					return servers, nil
				case "server.add":
					if params["host"] != "192.0.2.1" || params["password"] != "private-password" || params["autoConnect"] != true || params["allowUnauthorized"] != false {
						t.Error("incorrect registration parameters")
					}
					servers = append(servers, hostServer{ID: "target", Host: "192.0.2.1", Enabled: true, Status: "connected", PoolID: "pool"})
					return "target", nil
				case "server.disable", "server.set", "server.enable":
					if params["id"] != "target" {
						t.Error("unrelated host mutation")
					}
					if method == "server.set" && params["password"] != "private-password" {
						t.Error("missing password")
					}
				default:
					t.Errorf("unexpected method %s", method)
				}
				return nil, nil
			})
			var out bytes.Buffer
			if err := registerHost(r, "192.0.2.1", "root", "private-password", false, &out, func(time.Duration) {}); err != nil {
				t.Fatal(err)
			}
			want := []string{"server.getAll", "server.add", "server.getAll"}
			if existing {
				want = []string{"server.getAll", "server.getAll"}
			}
			if !reflect.DeepEqual(calls, want) {
				t.Fatalf("calls %v, want %v", calls, want)
			}
			if strings.Contains(out.String(), "private-password") || !strings.Contains(out.String(), "pool") {
				t.Fatal("unsafe or incomplete output")
			}
		})
	}
}

func TestRegisterHostRejectsDuplicates(t *testing.T) {
	r := fakeRPC(t, func(method string, _ json.RawMessage) (any, any) {
		if method != "server.getAll" {
			t.Error("mutation after duplicate matches")
		}
		return []hostServer{{ID: "one", Host: "192.0.2.1"}, {ID: "two", Host: "192.0.2.1"}}, nil
	})
	if err := registerHost(r, "192.0.2.1", "root", "password", false, &bytes.Buffer{}, func(time.Duration) {}); err == nil {
		t.Fatal("duplicate hosts accepted")
	}
}

func TestRegisterHostConnectionBound(t *testing.T) {
	reads := 0
	r := fakeRPC(t, func(method string, _ json.RawMessage) (any, any) {
		if method == "server.getAll" {
			reads++
			return []hostServer{{ID: "target", Host: "192.0.2.1", Status: "disconnected"}}, nil
		}
		return nil, nil
	})
	err := registerHost(r, "192.0.2.1", "root", "password", false, &bytes.Buffer{}, func(time.Duration) {})
	if err == nil || !strings.Contains(err.Error(), "bounded") || reads != 31 {
		t.Fatalf("unexpected bound: reads=%d error=%v", reads, err)
	}
}

func TestRegisterHostFailureRedacted(t *testing.T) {
	r := fakeRPC(t, func(method string, _ json.RawMessage) (any, any) {
		if method == "server.getAll" {
			return []hostServer{}, nil
		}
		return nil, map[string]any{"message": "private-password", "data": "private-username"}
	})
	err := registerHost(r, "192.0.2.1", "private-username", "private-password", false, &bytes.Buffer{}, func(time.Duration) {})
	if err == nil || strings.Contains(err.Error(), "private-") {
		t.Fatalf("unsafe failure: %v", err)
	}
}

func TestRegisterHostDisconnectedAndConnecting(t *testing.T) {
	for _, status := range []string{"disconnected", "connecting"} {
		t.Run(status, func(t *testing.T) {
			server := hostServer{ID: "target", Host: "192.0.2.1", Enabled: true, Status: status}
			var calls []string
			reads := 0
			r := fakeRPC(t, func(method string, _ json.RawMessage) (any, any) {
				calls = append(calls, method)
				switch method {
				case "server.getAll":
					reads++
					if status == "connecting" && reads > 1 {
						server.Status = "connected"
						server.PoolID = "pool"
					}
					return []hostServer{server}, nil
				case "server.disable":
					server.Enabled = false
				case "server.set":
					if server.Enabled {
						t.Error("credentials updated while automatic retries enabled")
					}
				case "server.enable":
					server.Enabled = true
					server.Status = "connected"
					server.PoolID = "pool"
				default:
					t.Errorf("unexpected method %s", method)
				}
				return nil, nil
			})
			if err := registerHost(r, "192.0.2.1", "root", "password", false, &bytes.Buffer{}, func(time.Duration) {}); err != nil {
				t.Fatal(err)
			}
			want := []string{"server.getAll", "server.getAll"}
			if status == "disconnected" {
				want = []string{"server.getAll", "server.disable", "server.set", "server.enable", "server.getAll"}
			}
			if !reflect.DeepEqual(calls, want) {
				t.Fatalf("calls %v, want %v", calls, want)
			}
		})
	}
}

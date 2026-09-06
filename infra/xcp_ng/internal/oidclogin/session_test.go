package oidclogin

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/net/websocket"
)

func TestCleanupDeletesOnlyIssuedToken(t *testing.T) {
	for _, reject := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "rejected_auth"}[reject], func(t *testing.T) {
			server := httptest.NewServer(websocket.Handler(func(conn *websocket.Conn) {
				defer conn.Close()
				for id, method := range []string{"session.signInWithToken", "token.deleteOwn"} {
					var request struct {
						ID     int             `json:"id"`
						Method string          `json:"method"`
						Params json.RawMessage `json:"params"`
					}
					if websocket.JSON.Receive(conn, &request) != nil {
						t.Error("missing cleanup request")
						return
					}
					if request.ID != id+1 || request.Method != method {
						t.Error("unexpected cleanup operation")
						return
					}
					if id == 0 {
						var params map[string]string
						if json.Unmarshal(request.Params, &params) != nil || len(params) != 1 || params["token"] != "private-issued-token" {
							t.Error("unexpected authentication scope")
						}
					} else {
						var params map[string][]string
						if json.Unmarshal(request.Params, &params) != nil || len(params) != 1 || len(params["tokens"]) != 1 || params["tokens"][0] != "private-issued-token" {
							t.Error("cleanup exceeded owned token scope")
						}
					}
					response := map[string]any{"id": request.ID, "result": nil}
					if reject {
						response["error"] = "private-issued-token"
					}
					if websocket.JSON.Send(conn, response) != nil {
						t.Error("sending fixture response failed")
					}
					if reject {
						return
					}
				}
			}))
			defer server.Close()
			conn, err := websocket.Dial(strings.Replace(server.URL, "http:", "ws:", 1), "", server.URL)
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			err = revokeConnection(context.Background(), conn, "private-issued-token")
			if reject {
				if err == nil || err.Error() != "XO rejected session cleanup" {
					t.Fatal("cleanup failure was not safely reported")
				}
			} else if err != nil {
				t.Fatal(err)
			}
		})
	}
}

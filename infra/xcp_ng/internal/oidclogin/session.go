package oidclogin

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/url"
	"time"

	"golang.org/x/net/websocket"
)

// Revoke deletes only the session token issued for this invocation.
func Revoke(ctx context.Context, origin *url.URL, token string) error {
	endpoint := *origin
	endpoint.Scheme = "wss"
	endpoint.Path = "/api/"
	conf, err := websocket.NewConfig(endpoint.String(), origin.String())
	if err != nil {
		return errors.New("invalid XO cleanup endpoint")
	}
	conf.Dialer = &net.Dialer{Timeout: 10 * time.Second}
	conn, err := websocket.DialConfig(conf)
	if err != nil {
		return errors.New("connecting for XO token cleanup failed")
	}
	defer conn.Close()
	return revokeConnection(ctx, conn, token)
}

func revokeConnection(ctx context.Context, conn *websocket.Conn, token string) error {
	conn.MaxPayloadBytes = 1 << 20
	deadline := time.Now().Add(10 * time.Second)
	if end, ok := ctx.Deadline(); ok && end.Before(deadline) {
		deadline = end
	}
	if err := conn.SetDeadline(deadline); err != nil {
		return errors.New("setting XO cleanup deadline failed")
	}
	for i, request := range []struct {
		method string
		params any
	}{
		{"session.signInWithToken", map[string]string{"token": token}},
		{"token.deleteOwn", map[string][]string{"tokens": {token}}},
	} {
		if websocket.JSON.Send(conn, map[string]any{"jsonrpc": "2.0", "id": i + 1, "method": request.method, "params": request.params}) != nil {
			return errors.New("sending XO cleanup request failed")
		}
		matched := false
		for n := 0; n < 1000; n++ {
			var response struct {
				ID    int             `json:"id"`
				Error json.RawMessage `json:"error"`
			}
			if websocket.JSON.Receive(conn, &response) != nil {
				return errors.New("receiving XO cleanup response failed")
			}
			if response.ID != i+1 {
				continue
			}
			if len(response.Error) != 0 && string(response.Error) != "null" {
				return errors.New("XO rejected session cleanup")
			}
			matched = true
			break
		}
		if !matched {
			return errors.New("XO cleanup notification limit exceeded")
		}
	}
	return nil
}

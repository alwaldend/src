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

type ScopeReport struct {
	ScopeVerified         bool `json:"scope_verified"`
	ResourceSetCount      int  `json:"resource_set_count"`
	PermissionObjectCount int  `json:"permission_object_count"`
	OwnVMAdmin            bool `json:"own_vm_admin"`
}

type scopeCall func(string, any, any) error

// VerifyScope uses successful API reads as evidence; errors never prove denial.
func VerifyScope(ctx context.Context, origin *url.URL, token, expectedSet, vmID string) (*ScopeReport, error) {
	endpoint := *origin
	endpoint.Scheme = "wss"
	endpoint.Path = "/api/"
	conf, err := websocket.NewConfig(endpoint.String(), origin.String())
	if err != nil {
		return nil, errors.New("invalid XO verification endpoint")
	}
	conf.Dialer = &net.Dialer{Timeout: 10 * time.Second}
	conn, err := websocket.DialConfig(conf)
	if err != nil {
		return nil, errors.New("connecting for XO scope verification failed")
	}
	defer conn.Close()
	conn.MaxPayloadBytes = 4 << 20
	deadline := time.Now().Add(30 * time.Second)
	if end, ok := ctx.Deadline(); ok && end.Before(deadline) {
		deadline = end
	}
	if conn.SetDeadline(deadline) != nil {
		return nil, errors.New("setting XO verification deadline failed")
	}
	id := 0
	call := func(method string, params, result any) error {
		id++
		if websocket.JSON.Send(conn, map[string]any{"jsonrpc": "2.0", "id": id, "method": method, "params": params}) != nil {
			return errors.New("sending XO scope verification failed")
		}
		for n := 0; n < 1000; n++ {
			var response struct {
				ID     int             `json:"id"`
				Result json.RawMessage `json:"result"`
				Error  json.RawMessage `json:"error"`
			}
			if websocket.JSON.Receive(conn, &response) != nil {
				return errors.New("receiving XO scope verification failed")
			}
			if response.ID != id {
				continue
			}
			if len(response.Error) > 0 && string(response.Error) != "null" {
				return errors.New("XO rejected scope verification")
			}
			if result != nil && json.Unmarshal(response.Result, result) != nil {
				return errors.New("invalid XO scope verification response")
			}
			return nil
		}
		return errors.New("XO scope verification notification limit exceeded")
	}
	if err := call("session.signInWithToken", map[string]string{"token": token}, nil); err != nil {
		return nil, err
	}
	return verifyScope(call, expectedSet, vmID)
}

func verifyScope(call scopeCall, expectedSet, vmID string) (*ScopeReport, error) {
	if expectedSet == "" {
		return nil, errors.New("expected resource set name is required")
	}
	var user struct {
		ID         string `json:"id"`
		Permission string `json:"permission"`
	}
	if err := call("session.getUser", map[string]any{}, &user); err != nil {
		return nil, errors.New("reading authenticated XO identity failed")
	}
	if user.ID == "" || user.Permission == "admin" {
		return nil, errors.New("XO scope verification requires a non-administrator identity")
	}
	var sets []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := call("resourceSet.getAll", map[string]any{}, &sets); err != nil {
		return nil, errors.New("reading own resource sets failed")
	}
	if len(sets) != 1 || sets[0].Name != expectedSet || sets[0].ID == "" {
		return nil, errors.New("XO resource set scope does not match expectation")
	}
	var permissions map[string]map[string]int
	if err := call("acl.getCurrentPermissions", map[string]any{}, &permissions); err != nil {
		return nil, errors.New("reading own XO permissions failed")
	}
	if permissions == nil {
		return nil, errors.New("XO permissions response is missing")
	}
	count := 0
	for object, actions := range permissions {
		for _, enabled := range actions {
			if enabled != 0 && object != vmID {
				return nil, errors.New("XO identity has permissions outside its own VM")
			}
		}
		if len(actions) > 0 {
			count++
		}
	}
	ownVMAdmin := false
	if vmID != "" {
		if permissions[vmID]["administrate"] != 1 {
			return nil, errors.New("XO identity cannot administer its own VM")
		}
		var objects map[string]struct {
			ID          string `json:"id"`
			Type        string `json:"type"`
			ResourceSet string `json:"resourceSet"`
		}
		if err := call("xo.getAllObjects", map[string]any{"filter": map[string]string{"id": vmID}}, &objects); err != nil {
			return nil, errors.New("reading own XO VM failed")
		}
		vm, ok := objects[vmID]
		if !ok || vm.ID != vmID || vm.Type != "VM" || vm.ResourceSet != sets[0].ID {
			return nil, errors.New("XO VM does not belong to the expected resource set")
		}
		ownVMAdmin = true
	}
	return &ScopeReport{ScopeVerified: true, ResourceSetCount: len(sets), PermissionObjectCount: count, OwnVMAdmin: ownVMAdmin}, nil
}

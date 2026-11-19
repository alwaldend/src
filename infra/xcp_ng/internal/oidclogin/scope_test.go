package oidclogin

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestScopeVerificationFailsClosed(t *testing.T) {
	for _, scenario := range []string{"valid", "no_vm", "global_admin", "other_set", "other_pool_acl", "other_vm_acl", "missing_vm_acl", "wrong_vm_set", "transport_failure", "rpc_rejection"} {
		t.Run(scenario, func(t *testing.T) {
			vmID := "own-vm"
			if scenario == "no_vm" {
				vmID = ""
			}
			calls := 0
			call := func(method string, params, result any) error {
				calls++
				var value any
				switch method {
				case "session.getUser":
					permission := "none"
					if scenario == "global_admin" {
						permission = "admin"
					}
					value = map[string]string{"id": "private-user-id", "permission": permission}
				case "resourceSet.getAll":
					sets := []map[string]string{{"id": "own-set-id", "name": "own-set"}}
					if scenario == "other_set" {
						sets = append(sets, map[string]string{"id": "other-id", "name": "other-set"})
					}
					value = sets
				case "acl.getCurrentPermissions":
					if scenario == "transport_failure" || scenario == "rpc_rejection" {
						return errors.New("private-server-error")
					}
					permissions := map[string]map[string]int{}
					if scenario != "no_vm" && scenario != "missing_vm_acl" {
						permissions["own-vm"] = map[string]int{"view": 1, "operate": 1, "administrate": 1}
					}
					if scenario == "other_pool_acl" {
						permissions["other-pool"] = map[string]int{"view": 1}
					}
					if scenario == "other_vm_acl" {
						permissions["other-vm"] = map[string]int{"administrate": 1}
					}
					value = permissions
				case "xo.getAllObjects":
					filter, ok := params.(map[string]any)["filter"].(map[string]string)
					if !ok || len(filter) != 1 || filter["id"] != "own-vm" {
						t.Fatal("VM inspection was not restricted to requested VM")
					}
					set := "own-set-id"
					if scenario == "wrong_vm_set" {
						set = "other-set-id"
					}
					value = map[string]any{"own-vm": map[string]string{"id": "own-vm", "type": "VM", "resourceSet": set}}
				default:
					t.Fatal("unexpected API operation")
				}
				raw, err := json.Marshal(value)
				if err != nil {
					t.Fatal(err)
				}
				return json.Unmarshal(raw, result)
			}
			report, err := verifyScope(call, "own-set", vmID)
			if scenario == "valid" || scenario == "no_vm" {
				if err != nil || !report.ScopeVerified || report.ResourceSetCount != 1 || report.OwnVMAdmin != (vmID != "") {
					t.Fatalf("valid scope rejected: %v", err)
				}
			} else if err == nil || report != nil {
				t.Fatal("unsafe or unverifiable scope accepted")
			}
			if err != nil && strings.Contains(err.Error(), "private") {
				t.Fatal("verification disclosed server response")
			}
			if calls == 0 {
				t.Fatal("scope was not inspected")
			}
		})
	}
}

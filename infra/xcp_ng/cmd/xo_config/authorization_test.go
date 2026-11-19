package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestOIDCIdentityDiscovery(t *testing.T) {
	for _, duplicate := range []bool{false, true} {
		users := []map[string]any{
			{"id": "oidc-id", "email": "mutable-name", "authProviders": map[string]any{"oidc:https://vault.test/provider": map[string]string{"id": "vault-entity-id"}}},
			{"id": "plain-id", "email": "mutable-name"},
			{"id": "foreign-oidc-id", "authProviders": map[string]any{"oidc:https://other.test/provider": map[string]string{"id": "vault-entity-id"}}},
		}
		if duplicate {
			users = append(users, users[0])
		}
		r := fakeRPC(t, func(method string, params json.RawMessage) (any, any) {
			switch method {
			case "user.getAll":
				return users, nil
			case "acl.get":
				return []map[string]string{
					{"id": "owned-acl", "subject": "admin-id", "object": "pool", "action": "admin"},
					{"id": "unrelated-acl", "subject": "someone-else", "object": "pool", "action": "admin"},
					{"id": "view-acl", "subject": "admin-id", "object": "other", "action": "view"},
				}, nil
			case "group.getAll":
				return []map[string]string{{"id": "admin-id", "name": "admins", "provider": "oidc"}, {"id": "plain-id", "name": "admins"}}, nil
			default:
				t.Fatalf("unexpected mutation or read: %s", method)
				return nil, nil
			}
		})
		var out bytes.Buffer
		err := writeOIDCIdentities(r, "oidc:https://vault.test/provider", "admins", &out)
		if duplicate {
			if err == nil || out.Len() != 0 {
				t.Fatal("ambiguous identity discovery did not fail closed")
			}
			continue
		}
		if err != nil {
			t.Fatal(err)
		}
		var result map[string]string
		if json.Unmarshal(out.Bytes(), &result) != nil {
			t.Fatal("invalid external protocol")
		}
		var ids map[string]string
		if json.Unmarshal([]byte(result["users"]), &ids) != nil || len(ids) != 1 || ids["vault-entity-id"] != "oidc-id" {
			t.Fatal("identity discovery did not use immutable OIDC subject")
		}
		var aclIDs map[string]string
		if json.Unmarshal([]byte(result["pool_admin_acl_ids"]), &aclIDs) != nil || len(aclIDs) != 1 || aclIDs["pool"] != "owned-acl" {
			t.Fatal("ACL discovery selected an unrelated grant")
		}
		var groups map[string]string
		if json.Unmarshal([]byte(result["groups"]), &groups) != nil || len(groups) != 1 || groups["admins"] != "admin-id" {
			t.Fatal("identity discovery accepted an ordinary group")
		}
	}
}

func TestOIDCGroupDiscoveryRejectsAmbiguity(t *testing.T) {
	r := fakeRPC(t, func(method string, params json.RawMessage) (any, any) {
		if method != "group.getAll" {
			t.Fatalf("unexpected mutation or read: %s", method)
		}
		return []map[string]string{
			{"id": "first", "name": "admins", "provider": "oidc"},
			{"id": "second", "name": "admins", "provider": "oidc"},
		}, nil
	})
	if _, err := oidcGroups(r); err == nil {
		t.Fatal("ambiguous synchronized group accepted")
	}
}

func TestOIDCACLDiscoveryRejectsAmbiguity(t *testing.T) {
	r := fakeRPC(t, func(method string, params json.RawMessage) (any, any) {
		switch method {
		case "user.getAll":
			return []any{}, nil
		case "group.getAll":
			return []map[string]string{{"id": "admin-id", "name": "admins", "provider": "oidc"}}, nil
		case "acl.get":
			return []map[string]string{
				{"id": "first", "subject": "admin-id", "object": "pool", "action": "admin"},
				{"id": "second", "subject": "admin-id", "object": "pool", "action": "admin"},
			}, nil
		default:
			t.Fatalf("unexpected mutation or read: %s", method)
			return nil, nil
		}
	})
	var out bytes.Buffer
	if err := writeOIDCIdentities(r, "oidc:https://vault.test/provider", "admins", &out); err == nil || out.Len() != 0 {
		t.Fatal("ambiguous ACL adoption did not fail closed")
	}
}

package main

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/hashicorp/vault/api"
)

var roleName = regexp.MustCompile(`^(src_[a-z0-9_]+|user_simeonwarren)$`)

// issueRole uses only an explicitly selected issuer's existing permissions.
func issueRole(ctx context.Context, issuer *api.Client, role string) (client *api.Client, err error) {
	if !roleName.MatchString(role) {
		return nil, errors.New("invalid target AppRole")
	}
	secret, err := issuer.Logical().WriteWithContext(ctx, "auth/approle/role/"+role+"/secret-id", map[string]any{"num_uses": 1})
	if err != nil || secret == nil {
		return nil, errors.New("issuer cannot create target AppRole SecretID")
	}
	accessor, _ := secret.Data["secret_id_accessor"].(string)
	consumed := false
	defer func() {
		if !consumed && accessor != "" {
			cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if _, cleanupErr := issuer.Logical().WriteWithContext(cleanup, "auth/approle/role/"+role+"/secret-id-accessor/destroy", map[string]any{"secret_id_accessor": accessor}); cleanupErr != nil {
				err = errors.Join(err, errors.New("cleaning unused target AppRole SecretID failed"))
			}
		}
	}()
	secretID, _ := secret.Data["secret_id"].(string)
	if secretID == "" {
		return nil, errors.New("target AppRole SecretID missing")
	}
	client, err = api.NewClient(issuer.CloneConfig())
	if err != nil {
		return nil, errors.New("constructing target Vault client failed")
	}
	client.ClearToken()
	login, err := client.Logical().WriteWithContext(ctx, "auth/approle/login", map[string]any{"role_id": role, "secret_id": secretID})
	if err != nil || login == nil || login.Auth == nil || login.Auth.ClientToken == "" {
		return nil, errors.New("authenticating target AppRole failed")
	}
	consumed = true
	client.SetToken(login.Auth.ClientToken)
	return client, nil
}

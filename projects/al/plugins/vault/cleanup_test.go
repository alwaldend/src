package main

import (
	"context"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/plugins/vault/vault_proto"
)

func fixtureVaultConfig(address string) *al_proto.Config {
	return &al_proto.Config{
		VaultConn: map[string]*al_proto.VaultConn{"default": {Config: &al_proto.VaultConfig{Address: address}, Tls: &al_proto.VaultTLS{}}},
		VaultAuth: map[string]*al_proto.VaultAuth{"default": {NoAuth: true}},
	}
}

func TestNoAuthOverridesInheritedVaultToken(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "synthetic-inherited-token")
	config := fixtureVaultConfig("http://127.0.0.1:1")
	fetcher := newVaultEnvFetcher(config, al.NewVault(config))
	result, err := fetcher.Get(context.Background(), &vault_proto.Resource{Res: &vault_proto.Resource_VaultEnv{VaultEnv: &vault_proto.VaultEnv{}}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	value, ok := result.Env["VAULT_TOKEN"]
	if !ok || value != "" {
		t.Fatal("no_auth must explicitly override the inherited token")
	}
}

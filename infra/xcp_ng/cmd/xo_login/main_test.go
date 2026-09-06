package main

import (
	"context"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
)

func TestBootstrapConfigurationSelectsOnlyRequestedRole(t *testing.T) {
	t.Setenv("XO_BOOTSTRAP_ISSUER_APPROLE", "")
	t.Setenv("XO_BOOTSTRAP_APPROLE", "src_infra_dc1_forgejo1")
	cfg, err := al.LoadConfigs(context.Background(), "bootstrap.lua")
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.VaultAuth) != 1 || cfg.VaultAuth[0].Name != "xo_bootstrap" || cfg.VaultAuth[0].Approle.Name != "src_infra_dc1_forgejo1" {
		t.Fatal("bootstrap did not isolate the selected role")
	}
	if len(cfg.Plugins) != 1 || cfg.Plugins[0].Name != "xo_login" || len(cfg.PluginCalls) != 0 {
		t.Fatal("bootstrap loaded unexpected plugins")
	}
	t.Setenv("XO_BOOTSTRAP_APPROLE", "user_simeonwarren")
	if cfg, err := al.LoadConfigs(context.Background(), "bootstrap.lua"); err != nil || cfg.VaultAuth[0].Approle.Name != "user_simeonwarren" {
		t.Fatal("approved user AppRole rejected")
	}
	for _, bad := range []string{"", "../other", "src_role/secret-id", "src_role\nother", "user_other"} {
		t.Setenv("XO_BOOTSTRAP_APPROLE", bad)
		if _, err := al.LoadConfigs(context.Background(), "bootstrap.lua"); err == nil {
			t.Fatal("unsafe role accepted")
		}
	}
}

func TestOriginRejectsCredentialsAndAlternateDestinations(t *testing.T) {
	for _, raw := range []string{"http://xo.test", "https://user:private@xo.test", "https://xo.test/path", "https://xo.test/?token=private", "https://xo.test/#fragment", "wss://xo.test"} {
		if _, err := origin(raw); err == nil {
			t.Fatal("accepted unsafe origin")
		} else if err.Error() != "XO login requires an HTTPS origin" {
			t.Fatal("origin error must not disclose input")
		}
	}
	if u, err := origin("https://xo.test"); err != nil || u.String() != "https://xo.test/" {
		t.Fatal("valid origin rejected")
	}
}

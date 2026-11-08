package main

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"github.com/bazelbuild/rules_go/go/runfiles"
)

func TestBootstrapConfigurationSelectsOnlyRequestedRole(t *testing.T) {
	t.Setenv("XO_BOOTSTRAP_ISSUER_APPROLE", "")
	t.Setenv("XO_BOOTSTRAP_APPROLE", "src_infra_dc1_forgejo1")
	load := func() (*al_proto.Config, config, error) {
		path, err := runfiles.Rlocation(os.Getenv("BOOTSTRAP_CONFIG"))
		if err != nil {
			return nil, config{}, err
		}
		cfg, err := al.LoadConfigs(context.Background(), path)
		if err != nil {
			return nil, config{}, err
		}
		if len(cfg.Plugins) != 1 || cfg.Plugins[0].Name != "xo_login" || len(cfg.PluginCalls) != 0 {
			t.Fatal("bootstrap loaded unexpected plugins")
		}
		value, err := al.FromPbJson(cfg.Plugins[0].Data).Get()
		if err != nil {
			return nil, config{}, err
		}
		raw, err := json.Marshal(value)
		if err != nil {
			return nil, config{}, err
		}
		var conf config
		if err := json.Unmarshal(raw, &conf); err != nil {
			return nil, conf, err
		}
		if !conf.Bootstrap {
			t.Fatal("bootstrap mode not configured")
		}
		err = configureBootstrap(&conf, cfg)
		return cfg, conf, err
	}
	cfg, conf, err := load()
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.VaultAuth) != 1 || cfg.VaultAuth[0].Name != "xo_bootstrap" || cfg.VaultAuth[0].Approle.Name != "src_infra_dc1_forgejo1" || conf.VaultAuth != "xo_bootstrap" {
		t.Fatal("bootstrap did not isolate the selected role")
	}
	t.Setenv("XO_BOOTSTRAP_APPROLE", "user_simeonwarren")
	if cfg, _, err := load(); err != nil || cfg.VaultAuth[0].Approle.Name != "user_simeonwarren" {
		t.Fatal("approved user AppRole rejected")
	}
	t.Setenv("XO_BOOTSTRAP_ISSUER_APPROLE", "src_infra_flux")
	if cfg, conf, err := load(); err != nil || cfg.VaultAuth[0].Approle.Name != "src_infra_flux" || conf.IssuerAuth != "xo_bootstrap" || conf.TargetApprole != "user_simeonwarren" {
		t.Fatal("explicit issuer not selected")
	}
	for _, bad := range []string{"../other", "src_role/secret-id", "src_role\nother", "user_simeonwarren"} {
		t.Setenv("XO_BOOTSTRAP_ISSUER_APPROLE", bad)
		if _, _, err := load(); err == nil {
			t.Fatal("unsafe issuer accepted")
		}
	}
	t.Setenv("XO_BOOTSTRAP_ISSUER_APPROLE", "")
	for _, bad := range []string{"", "../other", "src_role/secret-id", "src_role\nother", "user_other"} {
		t.Setenv("XO_BOOTSTRAP_APPROLE", bad)
		if _, _, err := load(); err == nil {
			t.Fatal("unsafe role accepted")
		}
	}
}

func TestOriginRejectsCredentialsAndAlternateDestinations(t *testing.T) {
	for _, raw := range []string{"http://xo.test", "https://user:private@xo.test", "https://xo.test/path", "https://xo.test/?token=private", "https://xo.test/#fragment", "wss://xo.test"} { // trufflehog:ignore Synthetic origin-rejection fixtures.
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

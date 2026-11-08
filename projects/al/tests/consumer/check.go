package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/api/contracts"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

// main runs the independent consumer outside Bazel actions and retains verified outputs.
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	var source, output string
	command := &cobra.Command{Use: "consumer-check", SilenceUsage: true, Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		if err := checkConsumer(cmd.Context(), source, output); err != nil {
			return fmt.Errorf("validate independent consumer: %w", err)
		}
		return nil
	}}
	command.Flags().StringVar(&source, "source_root", "", "Repository root containing the first-party modules")
	command.Flags().StringVar(&output, "output", "", "Absent or empty directory for the consumer and result.json")
	command.SetContext(ctx)
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

// checkConsumer builds with independent configuration and verifies the public output contract.
func checkConsumer(ctx context.Context, source, output string) error {
	if !filepath.IsAbs(source) || !filepath.IsAbs(output) {
		return fmt.Errorf("--source_root and --output must be absolute paths")
	}
	if err := os.MkdirAll(output, 0o700); err != nil {
		return fmt.Errorf("create consumer output: %w", err)
	}
	entries, err := os.ReadDir(output)
	if err != nil {
		return fmt.Errorf("inspect consumer output: %w", err)
	}
	if len(entries) != 0 {
		return fmt.Errorf("consumer output must be empty")
	}
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"data":{"verified":true}}`)); err != nil {
			return
		}
	}))
	defer server.Close()
	certificate := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: server.Certificate().Raw})
	key, err := x509.MarshalPKCS8PrivateKey(server.TLS.Certificates[0].PrivateKey)
	if err != nil {
		return fmt.Errorf("encode synthetic TLS fixture key: %w", err)
	}
	for name, data := range map[string][]byte{
		"fixture-ca.crt":     certificate,
		"fixture-client.crt": certificate,
		"fixture-client.key": pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: key}),
	} {
		if err := os.WriteFile(filepath.Join(output, name), data, 0o600); err != nil {
			return fmt.Errorf("write synthetic TLS fixture %q: %w", name, err)
		}
	}
	module := filepath.Join(source, "projects/al")
	inputs := map[string]string{
		"MODULE.bazel":                   filepath.Join(module, "tests/consumer/MODULE.bazel.in"),
		"duplicate_resource/BUILD.bazel": filepath.Join(module, "tests/consumer/duplicate_resource/BUILD.bazel.in"),
		"scalar_env/BUILD.bazel":         filepath.Join(module, "tests/consumer/scalar_env/BUILD.bazel.in"),
		"BUILD.bazel":                    filepath.Join(module, "tests/consumer/BUILD.bazel.in"),
		"payload.txt":                    filepath.Join(module, "tests/payload.txt"),
		"status.txt":                     filepath.Join(module, "tests/status.txt"),
		".bazeliskrc":                    filepath.Join(source, ".bazeliskrc"),
	}
	for name, input := range inputs {
		data, err := os.ReadFile(input)
		if err != nil {
			return fmt.Errorf("read consumer input %q: %w", name, err)
		}
		if err := os.MkdirAll(filepath.Dir(filepath.Join(output, name)), 0o700); err != nil {
			return fmt.Errorf("create consumer input directory: %w", err)
		}
		if err := os.WriteFile(filepath.Join(output, name), data, 0o600); err != nil {
			return fmt.Errorf("write consumer input %q: %w", name, err)
		}
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()
	args := []string{"build", "//:bundle", "//:config_runner", "//:layered_runner", "//:flat_config", "@rules_al//:al"}
	for _, name := range []string{"rules_al", "rules_docs", "rules_skills", "rules_promptfoo", "rules_openspec"} {
		location := filepath.Join(source, "tools", name)
		if name == "rules_al" {
			location = module
		}
		args = append(args, "--override_module="+name+"="+location)
	}
	build := exec.CommandContext(ctx, "bazel", args...)
	build.Dir = output
	build.Stdout, build.Stderr = os.Stdout, os.Stderr
	// Bazel run variables identify the caller, whereas this command owns a new workspace.
	for _, value := range os.Environ() {
		if len(value) >= 6 && value[:6] == "BUILD_" {
			continue
		}
		build.Env = append(build.Env, value)
	}
	if err := build.Run(); err != nil {
		return fmt.Errorf("build external consumer: %w", err)
	}
	for _, failure := range []struct{ target, diagnostic, log string }{
		{"//:ambiguous", "Conflicting vault_auth configuration", "ambiguity.log"},
		{"//duplicate_resource:invalid", "Duplicate injector resource", "duplicate-resource.log"},
		{"//scalar_env:invalid", "Injector resource settings must be dictionaries", "scalar-env.log"},
	} {
		negativeArgs := []string{"build", failure.target}
		for _, argument := range args {
			if strings.HasPrefix(argument, "--override_module=") {
				negativeArgs = append(negativeArgs, argument)
			}
		}
		negative := exec.CommandContext(ctx, "bazel", negativeArgs...)
		negative.Dir, negative.Env = output, build.Env
		var diagnostic bytes.Buffer
		negative.Stdout, negative.Stderr = &diagnostic, &diagnostic
		negativeErr := negative.Run()
		if negativeErr == nil {
			return fmt.Errorf("consumer accepted invalid target %s", failure.target)
		}
		if !strings.Contains(diagnostic.String(), failure.diagnostic) {
			return fmt.Errorf("consumer failed before checking %s: %w", failure.target, negativeErr)
		}
		if err := os.WriteFile(filepath.Join(output, failure.log), diagnostic.Bytes(), 0o600); err != nil {
			return fmt.Errorf("retain consumer rejection diagnostic: %w", err)
		}
	}
	bundle := filepath.Join(output, "bazel-bin/bundle")
	manifest, err := os.ReadFile(filepath.Join(bundle, "release.json"))
	if err != nil {
		return fmt.Errorf("read consumer manifest: %w", err)
	}
	release := &contracts.Release{}
	if err := protojson.Unmarshal(manifest, release); err != nil {
		return fmt.Errorf("decode consumer manifest: %w", err)
	}
	if release.Name != "2026.44.0" || release.SourceCommit != "0123456789012345678901234567890123456789" || release.Channel != "release" || release.TreeState != "clean" || release.GetProject().GetSubdir() != "projects/consumer" || len(release.Items) != 1 {
		return fmt.Errorf("external consumer emitted an incorrect source identity")
	}
	file := release.Items[0].GetFile()
	if file.GetName() != "payload.txt" {
		return fmt.Errorf("external consumer emitted an incorrect payload set")
	}
	payload, err := os.ReadFile(filepath.Join(bundle, "files/payload.txt"))
	if err != nil {
		return fmt.Errorf("read consumer payload: %w", err)
	}
	expected, err := os.ReadFile(filepath.Join(output, "payload.txt"))
	if err != nil {
		return fmt.Errorf("read expected payload: %w", err)
	}
	digest := fmt.Sprintf("%x", sha256.Sum256(payload))
	verified := false
	for _, hash := range file.Hashes {
		if hash.Algo == "SHA-256" && hash.Content == digest {
			verified = true
		}
	}
	if !bytes.Equal(payload, expected) || file.Size != int64(len(payload)) || !verified {
		return fmt.Errorf("consumer payload bytes, size, or digest do not match")
	}
	configContent, err := os.ReadFile(filepath.Join(output, "bazel-bin/config_/al.json"))
	if err != nil {
		return fmt.Errorf("read consumer configuration: %w", err)
	}
	config := &al_proto.Config{}
	if err := protojson.Unmarshal(configContent, config); err != nil {
		return fmt.Errorf("decode consumer configuration: %w", err)
	}
	if len(config.VaultAuth) != 1 || !config.VaultAuth["default"].NoAuth || len(config.VaultConn) != 1 || len(config.PluginConfigs) != 1 || len(config.Plugins) != 1 {
		return fmt.Errorf("consumer configuration lost declarations or explicit extension precedence")
	}
	if !reflect.DeepEqual(config.PluginConfigs["plugin"].Args, []string{"--help"}) || config.Plugins["call"].Labels["selected"] != "fixture" || config.VaultConn["default"].Config.Address != "https://vault.example.invalid" {
		return fmt.Errorf("consumer configuration lost labels or connection settings")
	}
	layeredContent, err := os.ReadFile(filepath.Join(output, "bazel-bin/layered_runner.config_/al.json"))
	if err != nil {
		return fmt.Errorf("read wrapper-composed configuration: %w", err)
	}
	layered := &al_proto.Config{}
	if err := protojson.Unmarshal(layeredContent, layered); err != nil {
		return fmt.Errorf("decode wrapper-composed configuration: %w", err)
	}
	if len(layered.VaultAuth) != 1 || !layered.VaultAuth["default"].NoAuth {
		return fmt.Errorf("wrapper composition lost explicit extension precedence")
	}
	parameters, err := al.FromPbJson(config.Plugins["call"].Data).Get()
	if err != nil {
		return fmt.Errorf("decode consumer plugin parameters: %w", err)
	}
	expectedParameters := map[string]any{"nested": []any{true, int64(7), float64(1.5), map[string]any{}, []any{}, nil, strings.TrimSpace(string(expected))}}
	if !reflect.DeepEqual(parameters, expectedParameters) {
		return fmt.Errorf("consumer configuration changed nested parameters or public file substitution")
	}
	flatContent, err := os.ReadFile(filepath.Join(output, "bazel-bin/flat_config_/al.json"))
	if err != nil {
		return fmt.Errorf("read flat injector configuration: %w", err)
	}
	flat := &al_proto.Config{}
	if err := protojson.Unmarshal(flatContent, flat); err != nil {
		return fmt.Errorf("decode flat injector configuration: %w", err)
	}
	calls := flat.Plugins
	if len(calls) != 3 || calls["dict_plugin"].GetPlugin() != "dict_plugin" || calls["default_call"].GetPlugin() != "default_call" || calls["explicit_call_name"].GetPlugin() != "target_plugin" {
		return fmt.Errorf("consumer lost target-derived plugin/call names or explicit call overrides")
	}
	injections := flat.Inject
	if len(injections) != 4 || injections["typed_resources"] == nil || injections["custom_resources"] == nil || injections["dict_resources"] == nil || injections["other_namespace"] == nil {
		return fmt.Errorf("consumer lost dedicated injector declarations")
	}
	flatResources := map[string]any{}
	for name, resource := range injections["typed_resources"].Res {
		args, err := al.FromPbJson(resource.Args).Get()
		if err != nil {
			return fmt.Errorf("decode injector resource %q: %w", name, err)
		}
		definition := map[string]any{"type": resource.Type, "args": args}
		if resource.Data != nil {
			data, err := al.FromPbJson(resource.Data).Get()
			if err != nil {
				return fmt.Errorf("decode injector resource %q data: %w", name, err)
			}
			definition["data"] = data
		}
		if len(resource.Deps) > 0 {
			deps := []any{}
			for _, dep := range resource.Deps {
				deps = append(deps, dep)
			}
			definition["deps"] = deps
		}
		flatResources[name] = definition
	}
	dictionary, err := al.FromPbJson(injections["dict_resources"].Res["DICT_VALUE"].Args).Get()
	if err != nil {
		return fmt.Errorf("decode dictionary injector resource: %w", err)
	}
	if !reflect.DeepEqual(dictionary, map[string]any{"value": "{{ .Data.value }}"}) {
		return fmt.Errorf("consumer lost dictionary injector settings")
	}
	dictionaryData, err := al.FromPbJson(injections["dict_resources"].Res["DICT_VALUE"].Data).Get()
	if err != nil || !reflect.DeepEqual(dictionaryData, map[string]any{"value": "dictionary"}) {
		return fmt.Errorf("consumer lost dictionary injector data")
	}
	other, err := al.FromPbJson(injections["other_namespace"].Res["DICT_VALUE"].Args).Get()
	if err != nil || !reflect.DeepEqual(other, map[string]any{"value": "unrelated namespace"}) {
		return fmt.Errorf("merging configuration mixed resource namespaces")
	}
	expectedResources := map[string]any{
		"EXAMPLE":      map[string]any{"type": "env", "deps": []any{"key", "payload_file"}, "args": map[string]any{"value": "fixture"}, "data": map[string]any{"nested": []any{true, nil, int64(9007199254740993)}}},
		"payload_file": map[string]any{"type": "file", "args": map[string]any{"value": strings.TrimSpace(string(expected))}, "data": map[string]any{"title": "payload"}},
		"key":          map[string]any{"type": "kv", "args": map[string]any{"vault_conn": "fixture_conn", "vault_auth": "fixture_auth", "mount": "secrets", "path": "fixture/key"}},
		"operation":    map[string]any{"type": "op", "args": map[string]any{"method": "read", "path": "fixture/operation"}, "data": map[string]any{"field": "value"}},
		"token":        map[string]any{"type": "oidc", "args": map[string]any{"name": "fixture-provider", "scope": "openid"}},
		"vault":        map[string]any{"type": "vault_env", "args": map[string]any{"auth": "fixture_auth", "conn": "fixture_conn"}},
		"ssh":          map[string]any{"type": "vault_ssh", "args": map[string]any{"backend": "ssh/clients/sign/fixture", "ttl": int64(60)}},
		"helper":       map[string]any{"type": "process", "args": map[string]any{"name": "fixture-tool", "args": []any{"--fixture"}}},
	}
	if !reflect.DeepEqual(flatResources, expectedResources) {
		return fmt.Errorf("flat injector groups changed typed settings, dependency order or Vault metadata")
	}
	for _, runtimePath := range []string{config.PluginConfigs["plugin"].Bin, config.VaultConn["default"].Tls.CaCert, config.VaultConn["default"].Tls.ClientCert, config.VaultConn["default"].Tls.ClientKey} {
		if _, err := os.Stat(filepath.Join(output, "bazel-bin/config_runner.exe.runfiles", runtimePath)); err != nil {
			return fmt.Errorf("inspect propagated consumer runtime input %q: %w", runtimePath, err)
		}
	}
	if err := verifyConsumerTLS(ctx, config, output, server.URL); err != nil {
		return fmt.Errorf("verify consumer TLS runtime: %w", err)
	}
	result, err := json.MarshalIndent(map[string]any{
		"schema": "release/consumer_check/v1", "verified": true, "configuration_verified": true, "tls_runtime_verified": true,
		"extensions_verified": true, "wrapper_extensions_verified": true, "ambiguity_rejected": true, "flat_injector_verified": true, "flat_resource_errors_rejected": true, "default_plugin_names_verified": true,
		"manifest": json.RawMessage(manifest), "payload_sha256": digest, "payload_size": len(payload),
	}, "", "  ")
	if err != nil {
		return fmt.Errorf("encode consumer result: %w", err)
	}
	if err := os.WriteFile(filepath.Join(output, "result.json"), append(result, '\n'), 0o600); err != nil {
		return fmt.Errorf("write consumer result: %w", err)
	}
	fmt.Println("Verified external consumer: " + strconv.Itoa(len(payload)) + " bytes, SHA-256 " + digest)
	return nil
}

// verifyConsumerTLS exercises the actual Vault client with generated runfile
// identifiers from an unrelated working directory in both runfiles modes.
func verifyConsumerTLS(ctx context.Context, config *al_proto.Config, output, address string) (resultErr error) {
	directory, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("read consumer working directory: %w", err)
	}
	saved := map[string]*string{}
	for _, key := range []string{"RUNFILES_DIR", "RUNFILES_MANIFEST_FILE"} {
		if value, ok := os.LookupEnv(key); ok {
			saved[key] = &value
		} else {
			saved[key] = nil
		}
	}
	defer func() {
		resultErr = errors.Join(resultErr, os.Chdir(directory))
		for key, value := range saved {
			if value == nil {
				resultErr = errors.Join(resultErr, os.Unsetenv(key))
			} else {
				resultErr = errors.Join(resultErr, os.Setenv(key, *value))
			}
		}
	}()
	if err := os.Chdir(output); err != nil {
		return fmt.Errorf("select unrelated TLS working directory: %w", err)
	}
	runtimeRoot := filepath.Join(output, "bazel-bin/config_runner.exe.runfiles")
	for _, mode := range []string{"directory", "manifest"} {
		runtimeDirectory, runtimeManifest := runtimeRoot, ""
		if mode == "manifest" {
			runtimeDirectory, runtimeManifest = "", runtimeRoot+"/MANIFEST"
		}
		for key, value := range map[string]string{"RUNFILES_DIR": runtimeDirectory, "RUNFILES_MANIFEST_FILE": runtimeManifest} {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("select %s runfiles: %w", mode, err)
			}
		}
		connection := *config.VaultConn["default"]
		connection.Config = &al_proto.VaultConfig{Address: address}
		tlsConfig, err := al.VaultTlsConfig(&connection).Get()
		if err != nil {
			return fmt.Errorf("resolve %s TLS configuration: %w", mode, err)
		}
		for name, resolved := range map[string]string{"fixture-ca.crt": tlsConfig.CACert, "fixture-client.crt": tlsConfig.ClientCert, "fixture-client.key": tlsConfig.ClientKey} {
			actual, err := os.Stat(resolved)
			if err != nil {
				return fmt.Errorf("open %s TLS runfile %q: %w", mode, name, err)
			}
			expected, err := os.Stat(filepath.Join(output, name))
			if err != nil {
				return fmt.Errorf("inspect synthetic TLS input %q: %w", name, err)
			}
			if !os.SameFile(actual, expected) {
				return fmt.Errorf("%s TLS lookup selected a different input for %q", mode, name)
			}
		}
		store := al.NewVault(&al_proto.Config{VaultConn: map[string]*al_proto.VaultConn{"default": &connection}, VaultAuth: config.VaultAuth})
		client, err := store.Client(ctx, "default", "default").Get()
		if err != nil {
			return fmt.Errorf("construct %s Vault TLS client: %w", mode, err)
		}
		response, readErr := client.Client.Logical().ReadWithContext(ctx, "fixture")
		stopErr := store.Stop(ctx)
		if err := errors.Join(readErr, stopErr); err != nil {
			return fmt.Errorf("request %s local Vault TLS fixture: %w", mode, err)
		}
		if response == nil || response.Data["verified"] != true {
			return fmt.Errorf("%s Vault TLS fixture returned an unexpected result", mode)
		}
	}
	return nil
}

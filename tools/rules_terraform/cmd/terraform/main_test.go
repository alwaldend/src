package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

type mapResolver map[string]string

func (r mapResolver) Rlocation(name string) (string, error) {
	if resolved, found := r[name]; found {
		return resolved, nil
	}
	return "", os.ErrNotExist
}

func writeFile(t *testing.T, name, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(name), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(name, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
}

func fixtureProvider() providerPackage {
	platform := runtime.GOOS + "_" + runtime.GOARCH
	return providerPackage{
		Source: "registry.terraform.io/example/fixture", Version: "1.2.3-rc1", Platform: platform,
		Archive:    "fixture_provider/provider.zip",
		MirrorPath: path.Join("registry.terraform.io/example/fixture", "terraform-provider-fixture_1.2.3-rc1_"+platform+".zip"),
	}
}

func TestValidateArgsAcceptsSavedPlan(t *testing.T) {
	chdir := t.TempDir()
	planPath := filepath.Join(chdir, "reviewed.tfplan")
	writeFile(t, planPath, "plan contents are validated by Terraform")
	for _, name := range []string{planPath, "reviewed.tfplan"} {
		if err := validateArgs([]string{"apply", name}, chdir, true); err != nil {
			t.Fatalf("saved plan %q rejected: %v", name, err)
		}
	}
}

func TestValidateArgsRejectsUnsavedApply(t *testing.T) {
	chdir := t.TempDir()
	plan := filepath.Join(chdir, "reviewed.tfplan")
	writeFile(t, plan, "fixture")
	for _, args := range [][]string{
		nil,
		{"apply"},
		{"apply", ""},
		{"plan", plan},
		{"destroy", plan},
		{"apply", "-auto-approve"},
		{"apply", "-target=module.dns"},
		{"apply", plan, "-auto-approve"},
		{"apply", plan, plan},
		{"apply", "missing.tfplan"},
		{"apply", chdir},
	} {
		if err := validateArgs(args, chdir, true); err == nil {
			t.Errorf("accepted command without exactly one saved plan: %q", args)
		}
	}
}

func TestValidateArgsPreservesStandardCommands(t *testing.T) {
	for _, args := range [][]string{nil, {"apply"}, {"plan", "-target=module.dns"}, {"apply", "-auto-approve"}, {"apply", "missing.tfplan"}} {
		if err := validateArgs(args, "missing-directory", false); err != nil {
			t.Errorf("standard command rejected: %v", err)
		}
	}
}

func TestValidateEnvironmentRejectsArgumentsAndProviderReattachment(t *testing.T) {
	for _, name := range []string{"TF_CLI_ARGS", "TF_CLI_ARGS_apply", "TF_CLI_ARGS_init", "TF_CLI_ARGS_future_operation", "TF_REATTACH_PROVIDERS"} {
		for _, guard := range []bool{false, true} {
			err := validateEnvironment([]string{name + "=secret-bearing injected arguments"}, guard)
			if err == nil || !strings.Contains(err.Error(), name) || strings.Contains(err.Error(), "secret-bearing") {
				t.Errorf("%s: diagnostic must reject override and identify only variable name: %v", name, err)
			}
		}
	}
	if err := validateEnvironment([]string{"TF_CLI_ARGS_apply= "}, true); err == nil {
		t.Fatal("accepted nonempty whitespace value")
	}
	if err := validateEnvironment([]string{"TF_CLI_ARGS=", "TF_CLI_ARGS_apply=", "TF_REATTACH_PROVIDERS=", "TF_CLI_ARGSOTHER=-backup", "OTHER_TF_CLI_ARGS=-backup"}, true); err != nil {
		t.Fatal(err)
	}
}

func TestParseOptionsStopsAtTerraformCommand(t *testing.T) {
	opts, err := parseOptions([]string{"--direct", "--require-saved-plan", "--chdir", "root", "apply", "plan", "--direct"})
	if err != nil || !opts.direct || !opts.requireSavedPlan || opts.chdir != "root" || !reflect.DeepEqual(opts.args, []string{"apply", "plan", "--direct"}) {
		t.Fatalf("options = %+v, error = %v", opts, err)
	}
	for _, args := range [][]string{{"--chdir"}, {"--chdir", ""}, {"--chdir="}} {
		if _, err := parseOptions(args); err == nil {
			t.Errorf("accepted missing chdir: %q", args)
		}
	}
}

func TestProviderArgumentsCannotDownloadOrReplaceProviders(t *testing.T) {
	for _, args := range [][]string{
		{"init", "-plugin-dir", "other"},
		{"init", "-plugin-dir=other"},
		{"init", "--plugin-dir=other"},
		{"providers", "mirror", "other"},
		{"providers", "lock"},
	} {
		if err := validateProviderArguments(args); err == nil {
			t.Errorf("accepted provider-installation override %q", args)
		}
	}
	if err := validateProviderArguments([]string{"init", "-backend=false", "-input=false"}); err != nil {
		t.Fatal(err)
	}
}

func TestRuntimeUsesOnlyPackagedProvidersAndPreservesBackendEnvironment(t *testing.T) {
	root := t.TempDir()
	workingDirectory := filepath.Join(root, "runfiles", "_main", "infra", "fixture")
	writeFile(t, filepath.Join(workingDirectory, "main.tf"), "terraform {}")
	provider := fixtureProvider()
	archive := filepath.Join(root, "provider.zip")
	installed := filepath.Join(workingDirectory, ".terraform", "providers", filepath.FromSlash(provider.Source), provider.Version, provider.Platform)
	writeProviderPackage(t, archive, installed, map[string]string{"terraform-provider-fixture_v1.2.3-rc1": "declared executable"})
	config := configuration{
		Terraform: "terraform/terraform", Workspace: "_main", Chdir: "infra/fixture",
		WorkingDirectory: "_main/infra/fixture", Args: []string{"show", "-json", "reviewed.tfplan"},
		Providers: []providerPackage{provider},
	}
	var stdout, stderr bytes.Buffer
	var calls [][]string
	var generatedConfig string
	call := invocation{
		program: "fixture",
		environ: []string{
			"RUNFILES_DIR=" + filepath.Join(root, "runfiles"), "TMPDIR=" + root,
			"TF_CLI_CONFIG_FILE=/host/config", "TF_PLUGIN_CACHE_DIR=/host/providers",
			"TF_PLUGIN_CACHE_MAY_BREAK_DEPENDENCY_LOCK_FILE=1", "TF_HTTP_ADDRESS=http://loopback/state",
			"AL_TF_BACKEND_CONFIG_fixture=address=literal value with spaces",
		},
		runfiles: mapResolver{"terraform/terraform": filepath.Join(root, "terraform"), provider.Archive: archive},
		stdin:    strings.NewReader(""), stdout: &stdout, stderr: &stderr,
		execute: func(cmd *exec.Cmd) error {
			calls = append(calls, append([]string{}, cmd.Args[1:]...))
			if environmentValue(cmd.Env, "TF_PLUGIN_CACHE_DIR") != "" || environmentValue(cmd.Env, "TF_PLUGIN_CACHE_MAY_BREAK_DEPENDENCY_LOCK_FILE") != "" {
				t.Fatal("inherited provider cache escaped isolation")
			}
			if environmentValue(cmd.Env, "TF_HTTP_ADDRESS") != "http://loopback/state" {
				t.Fatal("HTTP backend environment changed")
			}
			generatedConfig = environmentValue(cmd.Env, "TF_CLI_CONFIG_FILE")
			content, err := os.ReadFile(generatedConfig)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(content), "filesystem_mirror") || strings.Contains(string(content), "direct") || strings.Contains(string(content), "network_mirror") || !strings.Contains(string(content), "disable_checkpoint = true") {
				t.Fatalf("provider configuration permits a network installation: %s", content)
			}
			packagedArchive := filepath.Join(filepath.Dir(generatedConfig), "providers", filepath.FromSlash(provider.MirrorPath))
			resolved, err := filepath.EvalSymlinks(packagedArchive)
			if err != nil || resolved != archive {
				t.Fatalf("provider mirror did not resolve packaged archive: %q, %v", resolved, err)
			}
			if len(calls) == 1 {
				fmt.Fprint(cmd.Stdout, "initialization diagnostic\n")
			} else if cmd.Args[2] == "version" {
				fmt.Fprintf(cmd.Stdout, `{"provider_selections":{%q:%q}}`, provider.Source, provider.Version)
			} else {
				fmt.Fprint(cmd.Stdout, `{ "format_version": "1.2" }`)
			}
			return nil
		},
	}
	code, err := call.run(config, nil)
	if err != nil || code != 0 {
		t.Fatalf("run = %d, %v", code, err)
	}
	want := [][]string{
		{"-chdir=" + workingDirectory, "init", "--backend-config", "address=literal value with spaces"},
		{"-chdir=" + workingDirectory, "version", "-json"},
		{"-chdir=" + workingDirectory, "show", "-json", "reviewed.tfplan"},
	}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("commands = %q, want %q", calls, want)
	}
	if stdout.String() != `{ "format_version": "1.2" }` || stderr.String() != "initialization diagnostic\n" {
		t.Fatalf("initialization corrupted command output: stdout %q, stderr %q", stdout.String(), stderr.String())
	}
	if _, err := os.Stat(generatedConfig); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("invocation CLI configuration was not cleaned: %v", err)
	}
	if _, err := os.Stat(filepath.Join(workingDirectory, "main.tf")); err != nil {
		t.Fatal("cleanup removed the Terraform workspace")
	}
}

func TestDirectInitSkipsAutomaticInitAndKeepsLiteralBackendArgument(t *testing.T) {
	root := t.TempDir()
	var calls [][]string
	call := invocation{
		program: "fixture", environ: []string{"TMPDIR=" + root, "AL_TF_BACKEND_CONFIG_test=key=value"},
		runfiles: mapResolver{"terraform/terraform": filepath.Join(root, "terraform")},
		stdin:    strings.NewReader(""), stdout: io.Discard, stderr: io.Discard,
		execute: func(cmd *exec.Cmd) error {
			calls = append(calls, cmd.Args[1:])
			if cmd.Args[2] == "version" {
				fmt.Fprint(cmd.Stdout, `{"provider_selections":{}}`)
			}
			return nil
		},
	}
	config := configuration{Terraform: "terraform/terraform", Args: []string{"--direct", "--chdir", root, "init", "-backend=false"}}
	if code, err := call.run(config, nil); code != 0 || err != nil {
		t.Fatalf("run = %d, %v", code, err)
	}
	want := [][]string{
		{"-chdir=" + root, "init", "-backend=false", "--backend-config", "key=value"},
		{"-chdir=" + root, "version", "-json"},
	}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("commands = %q, want %q", calls, want)
	}
}

func TestRunRejectsSavedPlanBeforeInitialization(t *testing.T) {
	root := t.TempDir()
	called := false
	call := invocation{
		environ: []string{"TMPDIR=" + root}, runfiles: mapResolver{},
		stdout: io.Discard, stderr: io.Discard,
		execute: func(*exec.Cmd) error { called = true; return nil },
	}
	config := configuration{Args: []string{"--chdir", root, "--require-saved-plan", "apply"}}
	if code, err := call.run(config, nil); code != 2 || err == nil || called {
		t.Fatalf("unsafe apply reached initialization: code %d, error %v, called %v", code, err, called)
	}
}

func TestProviderMirrorUsesDeclaredRunfilesLayout(t *testing.T) {
	root := t.TempDir()
	provider := fixtureProvider()
	mirror := filepath.Join(root, "runfiles", "_main", "fixture.providers")
	writeFile(t, filepath.Join(mirror, filepath.FromSlash(provider.MirrorPath)), "archive")
	call := invocation{environ: []string{"RUNFILES_DIR=" + filepath.Join(root, "runfiles")}, runfiles: mapResolver{}}
	config := configuration{Mirror: "_main/fixture.providers", Providers: []providerPackage{provider}}
	got, err := call.providerMirror(config, filepath.Join(root, "unused"))
	if err != nil || got != mirror {
		t.Fatalf("mirror = %q, %v; want %q", got, err, mirror)
	}
	if err := os.Remove(filepath.Join(mirror, filepath.FromSlash(provider.MirrorPath))); err != nil {
		t.Fatal(err)
	}
	if _, err := call.providerMirror(config, root); err == nil {
		t.Fatal("missing declared runfiles provider was not rejected")
	}
}

func TestProviderMirrorRejectsMissingUnsupportedAndConflictingPackages(t *testing.T) {
	provider := fixtureProvider()
	call := invocation{runfiles: mapResolver{}}
	if _, err := call.providerMirror(configuration{Providers: []providerPackage{provider}}, t.TempDir()); err == nil {
		t.Fatal("missing provider archive was accepted")
	}
	unsupported := provider
	unsupported.Platform = "unsupported_platform"
	if _, err := call.providerMirror(configuration{Providers: []providerPackage{unsupported}}, t.TempDir()); err == nil {
		t.Fatal("unsupported provider platform was accepted")
	}
	conflict := provider
	conflict.Version = "2.0.0"
	conflict.MirrorPath = strings.Replace(conflict.MirrorPath, "1.2.3-rc1", "2.0.0", 1)
	if _, err := call.providerMirror(configuration{Providers: []providerPackage{provider, conflict}}, t.TempDir()); err == nil || !strings.Contains(err.Error(), "multiple versions") {
		t.Fatalf("conflicting provider versions were accepted: %v", err)
	}
	for _, candidate := range []string{"../escape", "/absolute", "registry.terraform.io/Example/fixture"} {
		invalid := provider
		invalid.Source = candidate
		if err := validateProviderPackage(invalid, provider.Platform); err == nil {
			t.Errorf("invalid provider source %q accepted", candidate)
		}
	}
}

func TestWorkingDirectoryUsesRunfilesInsteadOfSource(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source", "infra", "fixture", "main.tf")
	writeFile(t, source, "terraform {}")
	runfilesDirectory := filepath.Join(root, "runfiles")
	workingDirectory := filepath.Join(runfilesDirectory, "_main", "infra", "fixture")
	if err := os.MkdirAll(workingDirectory, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(source, filepath.Join(workingDirectory, "main.tf")); err != nil {
		t.Fatal(err)
	}
	call := invocation{environ: []string{"RUNFILES_DIR=" + runfilesDirectory, "BUILD_WORKSPACE_DIRECTORY=" + filepath.Join(root, "source")}}
	got, err := call.workingDirectory(configuration{WorkingDirectory: "_main/infra/fixture"}, "")
	if err != nil || got != workingDirectory {
		t.Fatalf("working directory = %q, %v", got, err)
	}
}

func TestManifestOnlyTestsMaterializeDeclaredFilesWithoutSourceLocks(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source", "main.tf")
	writeFile(t, source, "terraform {}")
	manifest := filepath.Join(root, "MANIFEST")
	writeFile(t, manifest, "_main/infra/fixture/main.tf "+source+"\n")
	rf, err := runfiles.New(runfiles.ManifestFile(manifest), runfiles.SourceRepo(""))
	if err != nil {
		t.Fatal(err)
	}
	call := invocation{environ: []string{"TEST_TMPDIR=" + root}, runfiles: rf, stderr: io.Discard}
	config := configuration{
		Workspace: "_main", Chdir: "infra/fixture", WorkingDirectory: "_main/infra/fixture",
		Files: []workspaceFile{{Runfile: "_main/infra/fixture/main.tf", Path: "infra/fixture/main.tf"}},
	}
	workingDirectory, err := call.workingDirectory(config, "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(workingDirectory, filepath.Join(root, "terraform-workspace-")) {
		t.Fatalf("workspace escaped test temporary directory: %s", workingDirectory)
	}
	writeFile(t, filepath.Join(workingDirectory, ".terraform.lock.hcl"), "generated lock")
	if _, err := os.Stat(filepath.Join(filepath.Dir(source), ".terraform.lock.hcl")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("test wrote a source lockfile: %v", err)
	}
	call.environ = nil
	if _, err := call.workingDirectory(config, ""); err == nil || !strings.Contains(err.Error(), "--enable_runfiles") {
		t.Fatalf("manifest execution silently relocated persistent state: %v", err)
	}
}

func TestRuntimeTemporaryRootNeverFallsBackToSystemDirectory(t *testing.T) {
	if _, err := runtimeTemporaryRoot(nil); err == nil {
		t.Fatal("missing Bazel and TMPDIR context used a system directory")
	}
	root := t.TempDir()
	got, err := runtimeTemporaryRoot([]string{"BUILD_WORKSPACE_DIRECTORY=" + root})
	if err != nil || got != filepath.Join(root, "out", "rules_terraform", "runtime") {
		t.Fatalf("scratch root = %q, %v", got, err)
	}
	got, err = runtimeTemporaryRoot([]string{"TMPDIR=" + filepath.Join(root, "caller"), "TEST_TMPDIR=" + root})
	if err != nil || got != root {
		t.Fatalf("test scratch priority = %q, %v", got, err)
	}
}

func TestHCLStringEscapesTemplateOpeners(t *testing.T) {
	got := hclString(`some/${path}/%{directive}/"quoted"`)
	if got != `"some/$${path}/%%{directive}/\"quoted\""` {
		t.Fatalf("unsafe HCL literal: %s", got)
	}
}

func TestReadConfigurationSidecarAndManifestFallback(t *testing.T) {
	root := t.TempDir()
	config := configuration{Terraform: "terraform/terraform", Workspace: "_main"}
	encoded, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	program := filepath.Join(root, "target")
	writeFile(t, program+".json", string(encoded))
	if got, err := readConfiguration(program, nil, mapResolver{}); err != nil || !reflect.DeepEqual(got, config) {
		t.Fatalf("adjacent sidecar = %+v, %v", got, err)
	}
	program = filepath.Join(root, "target.runfiles", "_main", "infra", "fixture")
	rf := mapResolver{"_main/infra/fixture.json": filepath.Join(root, "target.json")}
	if got, err := readConfiguration(program, nil, rf); err != nil || !reflect.DeepEqual(got, config) {
		t.Fatalf("manifest sidecar = %+v, %v", got, err)
	}
}

func TestExitCodePreservesTerraformDetailedExitCode(t *testing.T) {
	command := exec.Command(os.Args[0], "-test.run=TestExitHelper")
	command.Env = append(os.Environ(), "RULES_TERRAFORM_EXIT_HELPER=1")
	err := command.Run()
	if got := commandExitCode(err); got != 2 {
		t.Fatalf("detailed exit code = %d, want 2; error %v", got, err)
	}
}

func TestExitHelper(t *testing.T) {
	if os.Getenv("RULES_TERRAFORM_EXIT_HELPER") == "1" {
		os.Exit(2)
	}
}

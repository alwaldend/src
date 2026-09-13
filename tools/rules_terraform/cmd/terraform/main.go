package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

type providerPackage struct {
	Source     string `json:"source"`
	Version    string `json:"version"`
	Platform   string `json:"platform"`
	Archive    string `json:"archive"`
	MirrorPath string `json:"mirror_path"`
}

type workspaceFile struct {
	Runfile string `json:"runfile"`
	Path    string `json:"path"`
}

type configuration struct {
	Terraform        string            `json:"terraform"`
	Workspace        string            `json:"workspace"`
	Chdir            string            `json:"chdir"`
	WorkingDirectory string            `json:"working_directory"`
	Args             []string          `json:"args"`
	Mirror           string            `json:"mirror"`
	Providers        []providerPackage `json:"providers"`
	Files            []workspaceFile   `json:"files"`
}

type options struct {
	direct           bool
	requireSavedPlan bool
	chdir            string
	args             []string
}

type resolver interface {
	Rlocation(string) (string, error)
}

type invocation struct {
	program  string
	environ  []string
	runfiles resolver
	stdin    io.Reader
	stdout   io.Writer
	stderr   io.Writer
	execute  func(*exec.Cmd) error
}

func main() {
	os.Exit(run())
}

func run() int {
	rf, err := runfiles.New(runfiles.ProgramName(os.Args[0]), runfiles.SourceRepo(""))
	if err != nil {
		fmt.Fprintln(os.Stderr, "rules_terraform: locate runfiles:", err)
		return 2
	}
	environ := replaceEnvironment(os.Environ(), rf.Env())
	config, err := readConfiguration(os.Args[0], environ, rf)
	if err != nil {
		fmt.Fprintln(os.Stderr, "rules_terraform: load target configuration:", err)
		return 2
	}
	call := invocation{
		program: os.Args[0], environ: environ, runfiles: rf,
		stdin: os.Stdin, stdout: os.Stdout, stderr: os.Stderr,
		execute: func(cmd *exec.Cmd) error { return cmd.Run() },
	}
	code, err := call.run(config, os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "rules_terraform:", err)
	}
	return code
}

func readConfiguration(program string, environ []string, rf resolver) (configuration, error) {
	candidates := []string{program + ".json"}
	if filepath.Base(program) == program {
		if executable, err := exec.LookPath(program); err == nil {
			candidates = append(candidates, executable+".json")
		}
	}
	// A wrapper may be invoked through a runfiles-relative name while its
	// sidecar is available only through the manifest.
	key := filepath.ToSlash(program) + ".json"
	if _, suffix, ok := strings.Cut(key, ".runfiles/"); ok {
		key = suffix
	} else if filepath.IsAbs(program) {
		if _, suffix, ok := strings.Cut(key, "/bin/"); ok {
			if strings.HasPrefix(suffix, "external/") {
				key = strings.TrimPrefix(suffix, "external/")
			} else {
				workspace := environmentValue(environ, "TEST_WORKSPACE")
				if workspace == "" {
					workspace = "_main"
				}
				key = path.Join(workspace, suffix)
			}
		}
	}
	if !filepath.IsAbs(key) {
		if resolved, err := rf.Rlocation(key); err == nil {
			candidates = append(candidates, resolved)
		}
	}
	for _, candidate := range candidates {
		content, err := os.ReadFile(candidate)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return configuration{}, err
		}
		var config configuration
		if err := json.Unmarshal(content, &config); err != nil {
			return configuration{}, fmt.Errorf("decode %s: %w", candidate, err)
		}
		return config, nil
	}
	return configuration{}, fmt.Errorf("target sidecar %s.json is missing from outputs and runfiles", program)
}

func parseOptions(args []string) (options, error) {
	var opts options
	for len(args) > 0 {
		switch args[0] {
		case "--direct":
			opts.direct = true
		case "--require-saved-plan":
			opts.requireSavedPlan = true
		case "--chdir":
			if len(args) < 2 || args[1] == "" {
				return options{}, errors.New("--chdir requires a directory")
			}
			opts.chdir = args[1]
			args = args[1:]
		case "--":
			opts.args = args[1:]
			return opts, nil
		default:
			if strings.HasPrefix(args[0], "--chdir=") {
				opts.chdir = strings.TrimPrefix(args[0], "--chdir=")
				if opts.chdir == "" {
					return options{}, errors.New("--chdir requires a directory")
				}
			} else {
				opts.args = args
				return opts, nil
			}
		}
		args = args[1:]
	}
	return opts, nil
}

func validateArgs(args []string, chdir string, requireSavedPlan bool) error {
	if !requireSavedPlan {
		return nil
	}
	if len(args) != 2 || args[0] != "apply" || args[1] == "" || strings.HasPrefix(args[1], "-") {
		return errors.New("--require-saved-plan requires apply <saved-plan-file> with no other arguments")
	}
	planPath := args[1]
	if !filepath.IsAbs(planPath) {
		planPath = filepath.Join(chdir, planPath)
	}
	info, err := os.Stat(planPath)
	if err != nil {
		return fmt.Errorf("could not inspect saved plan file: %w", err)
	}
	if !info.Mode().IsRegular() {
		return errors.New("saved plan path must refer to a regular file")
	}
	return nil
}

func validateEnvironment(environ []string, requireSavedPlan bool) error {
	for _, entry := range environ {
		name, value, _ := strings.Cut(entry, "=")
		if value == "" {
			continue
		}
		if name == "TF_CLI_ARGS" || strings.HasPrefix(name, "TF_CLI_ARGS_") {
			if requireSavedPlan {
				return fmt.Errorf("--require-saved-plan rejects nonempty %s", name)
			}
			return fmt.Errorf("packaged provider execution rejects nonempty %s; pass arguments explicitly", name)
		}
		if name == "TF_REATTACH_PROVIDERS" {
			return fmt.Errorf("packaged provider execution rejects nonempty %s", name)
		}
	}
	return nil
}

func validateProviderArguments(args []string) error {
	for _, arg := range args {
		option, _, _ := strings.Cut(arg, "=")
		if option == "-plugin-dir" || option == "--plugin-dir" {
			return errors.New("-plugin-dir cannot replace the providers declared by the Bazel target")
		}
	}
	if len(args) >= 2 && args[0] == "providers" && (args[1] == "mirror" || args[1] == "lock") {
		return errors.New("provider download commands are not supported at runtime; update the Bazel provider extension declarations")
	}
	return nil
}

func (call invocation) run(config configuration, extra []string) (int, error) {
	args := append(append([]string{}, config.Args...), extra...)
	opts, err := parseOptions(args)
	if err != nil {
		return 2, err
	}
	if err := validateEnvironment(call.environ, opts.requireSavedPlan); err != nil {
		return 2, err
	}
	if err := validateProviderArguments(opts.args); err != nil {
		return 2, err
	}
	chdir, err := call.workingDirectory(config, opts.chdir)
	if err != nil {
		return 2, err
	}
	if err := validateArgs(opts.args, chdir, opts.requireSavedPlan); err != nil {
		return 2, err
	}
	terraform, err := call.runfiles.Rlocation(config.Terraform)
	if err != nil {
		return 2, fmt.Errorf("locate declared Terraform executable: %w", err)
	}
	terraform, err = filepath.Abs(terraform)
	if err != nil {
		return 2, err
	}
	temporaryRoot, err := runtimeTemporaryRoot(call.environ)
	if err != nil {
		return 2, err
	}
	temporary, err := os.MkdirTemp(temporaryRoot, "terraform-runtime-")
	if err != nil {
		return 2, fmt.Errorf("create invocation scratch: %w", err)
	}
	defer os.RemoveAll(temporary)
	mirror, err := call.providerMirror(config, temporary)
	if err != nil {
		return 2, err
	}
	cliConfig := filepath.Join(temporary, "terraform.rc")
	content := fmt.Sprintf("disable_checkpoint = true\nprovider_installation {\n  filesystem_mirror {\n    path = %s\n  }\n}\n", hclString(filepath.ToSlash(mirror)))
	if err := os.WriteFile(cliConfig, []byte(content), 0o600); err != nil {
		return 2, fmt.Errorf("write isolated Terraform CLI configuration: %w", err)
	}
	environ := providerEnvironment(call.environ, cliConfig)
	backendArgs := []string{}
	for _, entry := range environ {
		name, value, found := strings.Cut(entry, "=")
		if found && strings.HasPrefix(name, "AL_TF_BACKEND_CONFIG") {
			backendArgs = append(backendArgs, "--backend-config", value)
		}
	}
	runCommand := func(args []string, stdout io.Writer) error {
		cmd := exec.Command(terraform, append([]string{"-chdir=" + chdir}, args...)...)
		cmd.Env = environ
		cmd.Stdin = call.stdin
		cmd.Stdout = stdout
		cmd.Stderr = call.stderr
		return call.execute(cmd)
	}
	if !opts.direct {
		if err := runCommand(append([]string{"init"}, backendArgs...), call.stderr); err != nil {
			return commandExitCode(err), fmt.Errorf("Terraform initialization failed: %w", err)
		}
	}
	if commandUsesProviders(opts.args) {
		if err := call.verifyProviderSelections(config, chdir, environ, runCommand); err != nil {
			return 2, err
		}
	}
	if len(opts.args) > 0 && opts.args[0] == "init" {
		opts.args = append(opts.args, backendArgs...)
	}
	if err := runCommand(opts.args, call.stdout); err != nil {
		return commandExitCode(err), fmt.Errorf("Terraform command failed: %w", err)
	}
	if len(opts.args) > 0 && opts.args[0] == "init" {
		if err := call.verifyProviderSelections(config, chdir, environ, runCommand); err != nil {
			return 2, err
		}
	}
	return 0, nil
}

func (call invocation) workingDirectory(config configuration, override string) (string, error) {
	if override != "" {
		return existingDirectory(override)
	}
	key := config.WorkingDirectory
	if key == "" {
		key = path.Join(config.Workspace, config.Chdir)
	}
	if !validRunfilePath(key) {
		return "", fmt.Errorf("invalid declared working directory %q", key)
	}
	for _, directory := range call.runfilesDirectories() {
		candidate := filepath.Join(directory, filepath.FromSlash(key))
		if result, err := existingDirectory(candidate); err == nil {
			return result, nil
		}
	}
	// Manifest entries may resolve individual files directly into source. Do
	// not use their parent directories: init would then write source locks.
	testTemporary := environmentValue(call.environ, "TEST_TMPDIR")
	if testTemporary == "" {
		return "", errors.New("Terraform requires a directory runfiles tree to preserve its working directory and local state; enable Bazel runfiles (--enable_runfiles), or explicitly pass --chdir before the Terraform command for a caller-owned directory")
	}
	if len(config.Files) == 0 {
		return "", errors.New("cannot materialize a manifest-only test workspace without declared files")
	}
	workspace, err := os.MkdirTemp(testTemporary, "terraform-workspace-")
	if err != nil {
		return "", err
	}
	// This workspace can contain plans or local state. It remains owned by
	// Bazel's test temporary directory; invocation cleanup removes only CLI
	// configuration and provider-mirror scratch.
	for _, file := range config.Files {
		if !validRunfilePath(file.Path) {
			return "", fmt.Errorf("invalid declared workspace path %q", file.Path)
		}
		source, err := call.runfiles.Rlocation(file.Runfile)
		if err != nil {
			return "", fmt.Errorf("locate workspace input %s: %w", file.Runfile, err)
		}
		source, err = filepath.Abs(source)
		if err != nil {
			return "", err
		}
		destination := filepath.Join(workspace, filepath.FromSlash(file.Path))
		if err := os.MkdirAll(filepath.Dir(destination), 0o700); err != nil {
			return "", err
		}
		if err := os.Symlink(source, destination); err != nil {
			return "", fmt.Errorf("materialize workspace input %s: %w", file.Path, err)
		}
	}
	if config.Chdir != "." && config.Chdir != "" && !validRunfilePath(config.Chdir) {
		return "", fmt.Errorf("invalid declared chdir %q", config.Chdir)
	}
	workingDirectory := filepath.Join(workspace, filepath.FromSlash(config.Chdir))
	if err := os.MkdirAll(workingDirectory, 0o700); err != nil {
		return "", err
	}
	fmt.Fprintln(call.stderr, "rules_terraform: manifest-only test workspace retained at", workingDirectory)
	return workingDirectory, nil
}

func (call invocation) runfilesDirectories() []string {
	candidates := []string{
		environmentValue(call.environ, "RUNFILES_DIR"),
		environmentValue(call.environ, "TEST_SRCDIR"),
		call.program + ".runfiles",
	}
	if manifest := environmentValue(call.environ, "RUNFILES_MANIFEST_FILE"); strings.HasSuffix(manifest, ".runfiles_manifest") {
		candidates = append(candidates, strings.TrimSuffix(manifest, "_manifest"))
	} else if filepath.Base(manifest) == "MANIFEST" {
		candidates = append(candidates, filepath.Dir(manifest))
	}
	var result []string
	for _, candidate := range candidates {
		if candidate != "" {
			if directory, err := existingDirectory(candidate); err == nil {
				result = append(result, directory)
			}
		}
	}
	return result
}

func (call invocation) providerMirror(config configuration, temporary string) (string, error) {
	packages := make(map[string]string)
	platform := runtime.GOOS + "_" + runtime.GOARCH
	for _, provider := range config.Providers {
		if err := validateProviderPackage(provider, platform); err != nil {
			return "", err
		}
		if version, found := packages[provider.Source]; found {
			if version != provider.Version {
				return "", fmt.Errorf("multiple versions declared for provider %s", provider.Source)
			}
			return "", fmt.Errorf("duplicate provider package %s", provider.Source)
		}
		packages[provider.Source] = provider.Version
	}
	if config.Mirror != "" {
		if !validRunfilePath(config.Mirror) {
			return "", errors.New("invalid declared provider mirror runfile path")
		}
		for _, directory := range call.runfilesDirectories() {
			mirror := filepath.Join(directory, filepath.FromSlash(config.Mirror))
			if _, err := existingDirectory(mirror); err != nil {
				continue
			}
			for _, provider := range config.Providers {
				if err := regularFile(filepath.Join(mirror, filepath.FromSlash(provider.MirrorPath))); err != nil {
					return "", fmt.Errorf("declared provider %s is missing from the runfiles mirror: %w", provider.Source, err)
				}
			}
			return mirror, nil
		}
	}
	mirror := filepath.Join(temporary, "providers")
	if err := os.Mkdir(mirror, 0o700); err != nil {
		return "", err
	}
	for _, provider := range config.Providers {
		archive, err := call.runfiles.Rlocation(provider.Archive)
		if err != nil {
			return "", fmt.Errorf("locate declared provider %s: %w", provider.Source, err)
		}
		archive, err = filepath.Abs(archive)
		if err != nil {
			return "", err
		}
		if err := regularFile(archive); err != nil {
			return "", fmt.Errorf("declared provider %s archive is unavailable: %w", provider.Source, err)
		}
		destination := filepath.Join(mirror, filepath.FromSlash(provider.MirrorPath))
		if err := os.MkdirAll(filepath.Dir(destination), 0o700); err != nil {
			return "", err
		}
		if err := os.Symlink(archive, destination); err != nil {
			return "", err
		}
	}
	return mirror, nil
}

func validateProviderPackage(provider providerPackage, platform string) error {
	parts := strings.Split(provider.Source, "/")
	if len(parts) != 3 || !validRunfilePath(provider.Source) || provider.Source != strings.ToLower(provider.Source) {
		return fmt.Errorf("provider source %q must be a canonical hostname/namespace/name", provider.Source)
	}
	if provider.Version == "" || strings.ContainsAny(provider.Version, "/\\") {
		return fmt.Errorf("provider %s has an invalid version", provider.Source)
	}
	if provider.Platform != platform {
		return fmt.Errorf("provider %s targets %s, but Terraform is running on %s; declare the matching Bazel provider package", provider.Source, provider.Platform, platform)
	}
	expected := path.Join(provider.Source, "terraform-provider-"+parts[2]+"_"+provider.Version+"_"+provider.Platform+".zip")
	if provider.MirrorPath != expected {
		return fmt.Errorf("provider %s has an invalid packed mirror path", provider.Source)
	}
	if !validRunfilePath(provider.Archive) {
		return fmt.Errorf("provider %s has an invalid archive runfile path", provider.Source)
	}
	return nil
}

func validRunfilePath(value string) bool {
	return value != "" && value != "." && !strings.HasPrefix(value, "/") && !strings.Contains(value, "\\") && path.Clean(value) == value && value != ".." && !strings.HasPrefix(value, "../")
}

func regularFile(name string) error {
	info, err := os.Stat(name)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("%s must be a regular file", name)
	}
	return nil
}

func existingDirectory(name string) (string, error) {
	info, err := os.Stat(name)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s must be a directory", name)
	}
	return filepath.Abs(name)
}

func runtimeTemporaryRoot(environ []string) (string, error) {
	for _, name := range []string{"TEST_TMPDIR", "TMPDIR"} {
		if value := environmentValue(environ, name); value != "" {
			if err := os.MkdirAll(value, 0o700); err != nil {
				return "", fmt.Errorf("create %s directory: %w", name, err)
			}
			return filepath.Abs(value)
		}
	}
	if workspace := environmentValue(environ, "BUILD_WORKSPACE_DIRECTORY"); workspace != "" {
		directory := filepath.Join(workspace, "out", "rules_terraform", "runtime")
		if err := os.MkdirAll(directory, 0o700); err != nil {
			return "", err
		}
		return directory, nil
	}
	return "", errors.New("set TMPDIR to a caller-owned scratch directory when running outside Bazel")
}

func environmentValue(environ []string, name string) string {
	var value string
	for _, entry := range environ {
		if key, candidate, _ := strings.Cut(entry, "="); key == name {
			value = candidate
		}
	}
	return value
}

func replaceEnvironment(environ, replacements []string) []string {
	keys := make(map[string]bool)
	for _, entry := range replacements {
		key, _, _ := strings.Cut(entry, "=")
		keys[key] = true
	}
	result := make([]string, 0, len(environ)+len(replacements))
	for _, entry := range environ {
		key, _, _ := strings.Cut(entry, "=")
		if !keys[key] {
			result = append(result, entry)
		}
	}
	return append(result, replacements...)
}

func providerEnvironment(environ []string, config string) []string {
	result := make([]string, 0, len(environ)+3)
	for _, entry := range environ {
		name, _, _ := strings.Cut(entry, "=")
		switch name {
		case "TF_CLI_CONFIG_FILE", "TF_PLUGIN_CACHE_DIR", "TF_PLUGIN_CACHE_MAY_BREAK_DEPENDENCY_LOCK_FILE", "CHECKPOINT_DISABLE", "TF_IN_AUTOMATION":
			continue
		}
		result = append(result, entry)
	}
	return append(result, "TF_CLI_CONFIG_FILE="+config, "CHECKPOINT_DISABLE=1", "TF_IN_AUTOMATION=1")
}

func hclString(value string) string {
	encoded, _ := json.Marshal(value)
	// HCL quoted strings interpret template openers; filesystem paths do not.
	return strings.ReplaceAll(strings.ReplaceAll(string(encoded), "${", "$${"), "%{", "%%{")
}

func commandExitCode(err error) int {
	var exitError *exec.ExitError
	if errors.As(err, &exitError) && exitError.ExitCode() > 0 {
		return exitError.ExitCode()
	}
	return 1
}

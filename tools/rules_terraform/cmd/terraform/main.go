package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// providerPackage identifies a single declared provider archive and its mirror location.
type providerPackage struct {
	// Source is the provider address used in the dependency lock.
	Source string `json:"source"`
	// Version is the version of the declared archive.
	Version string `json:"version"`
	// Platform selects the operating system and architecture for the archive.
	Platform string `json:"platform"`
	// Archive identifies the Bazel runfile containing the provider ZIP.
	Archive string `json:"archive"`
	// MirrorPath locates the provider ZIP within the packed mirror.
	MirrorPath string `json:"mirror_path"`
}

// workspaceFile maps a runfile into a manifest-only test workspace.
type workspaceFile struct {
	// Runfile locates the declared workspace input.
	Runfile string `json:"runfile"`
	// Path places the file relative to the materialized workspace.
	Path string `json:"path"`
}

// configuration decodes the Bazel-generated target sidecar.
type configuration struct {
	// Terraform identifies the executable runfile.
	Terraform string `json:"terraform"`
	// Workspace names the runfiles repository containing the configuration.
	Workspace string `json:"workspace"`
	// Chdir selects the package directory within the workspace.
	Chdir string `json:"chdir"`
	// WorkingDirectory can supply the full runfiles key of the working directory.
	WorkingDirectory string `json:"working_directory"`
	// Args precedes the arguments supplied by the caller.
	Args []string `json:"args"`
	// Mirror locates the generated provider mirror in directory runfiles.
	Mirror string `json:"mirror"`
	// Providers declares the permitted provider archives.
	Providers []providerPackage `json:"providers"`
	// Files supplies manifest-only tests with their workspace inputs.
	Files []workspaceFile `json:"files"`
}

// options separates launcher controls from Terraform arguments.
type options struct {
	// direct skips automatic initialization when set.
	direct bool
	// requireSavedPlan restricts apply to a single existing plan.
	requireSavedPlan bool
	// chdir overrides the runfiles working directory.
	chdir string
	// args retains the literal Terraform argument suffix.
	args []string
}

// resolver resolves Bazel runfile keys for the active invocation.
type resolver interface {
	// Rlocation maps a canonical runfile name to its filesystem path.
	Rlocation(string) (string, error)
}

// invocation owns the environment and borrowed process streams for one launch.
type invocation struct {
	// context cancels initialization, provenance inspection, and Terraform execution.
	context context.Context
	// program is the launcher path used to discover its runfiles tree.
	program string
	// environ contains the inherited process and Bazel runfiles environment.
	environ []string
	// runfiles resolves canonical keys in the active runfiles layout.
	runfiles resolver
	// stdin supplies Terraform’s interactive input.
	stdin io.Reader
	// stdout receives the requested command’s normal output.
	stdout io.Writer
	// stderr receives diagnostics and initialization output.
	stderr io.Writer
	// execute runs a configured child and returns its exit error.
	execute func(*exec.Cmd) error
}

// maxSidecarBytes limits the generated target configuration read by one invocation.
const maxSidecarBytes = 64 << 20

// readSidecar reads a single generated JSON configuration before decoding.
func readSidecar(path string) (content []byte, result error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open sidecar %q: %w", path, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("close sidecar %q: %w", path, err))
		}
	}()
	content, err = io.ReadAll(io.LimitReader(file, maxSidecarBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read sidecar %q: %w", path, err)
	}
	if len(content) > maxSidecarBytes {
		return nil, fmt.Errorf("sidecar %q exceeds %d bytes", path, maxSidecarBytes)
	}
	return content, nil
}

// main exits with the status of the target invocation.
func main() {
	os.Exit(run())
}

// run loads the rule-generated sidecar and connects the process environment.
func run() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
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
		context: ctx,
		program: os.Args[0], environ: environ, runfiles: rf,
		stdin: os.Stdin, stdout: os.Stdout, stderr: os.Stderr,
		execute: func(cmd *exec.Cmd) error {
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("execute Terraform: %w", err)
			}
			return nil
		},
	}
	code, err := call.run(config, os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "rules_terraform:", err)
	}
	return code
}

// readConfiguration loads the adjacent or manifest-resolved target sidecar.
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
		content, err := readSidecar(candidate)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return configuration{}, fmt.Errorf("read target sidecar %q: %w", candidate, err)
		}
		var config configuration
		if err := json.Unmarshal(content, &config); err != nil {
			return configuration{}, fmt.Errorf("decode %s: %w", candidate, err)
		}
		return config, nil
	}
	return configuration{}, fmt.Errorf("target sidecar %s.json is missing from outputs and runfiles", program)
}

// parseOptions consumes the launcher prefix and preserves the entire Terraform suffix.
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

// validateArgs checks the single saved-plan argument when the target requires one.
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

// validateEnvironment rejects inherited settings that bypass declared provider execution.
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

// validateProviderArguments rejects commands and options that change provider installation sources.
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

// run prepares invocation-owned files and executes Terraform with the declared providers.
func (call invocation) run(config configuration, extra []string) (code int, result error) {
	args := append(append([]string{}, config.Args...), extra...)
	opts, err := parseOptions(args)
	if err != nil {
		return 2, fmt.Errorf("parse launcher options: %w", err)
	}
	if err := validateEnvironment(call.environ, opts.requireSavedPlan); err != nil {
		return 2, fmt.Errorf("check Terraform environment: %w", err)
	}
	if err := validateProviderArguments(opts.args); err != nil {
		return 2, fmt.Errorf("validate provider arguments: %w", err)
	}
	chdir, err := call.workingDirectory(config, opts.chdir)
	if err != nil {
		return 2, fmt.Errorf("resolve working directory: %w", err)
	}
	if err := validateArgs(opts.args, chdir, opts.requireSavedPlan); err != nil {
		return 2, fmt.Errorf("validate Terraform arguments: %w", err)
	}
	terraform, err := call.runfiles.Rlocation(config.Terraform)
	if err != nil {
		return 2, fmt.Errorf("locate declared Terraform executable: %w", err)
	}
	terraform, err = filepath.Abs(terraform)
	if err != nil {
		return 2, fmt.Errorf("resolve Terraform executable path: %w", err)
	}
	temporaryRoot, err := runtimeTemporaryRoot(call.environ)
	if err != nil {
		return 2, fmt.Errorf("select invocation scratch: %w", err)
	}
	temporary, err := os.MkdirTemp(temporaryRoot, "terraform-runtime-")
	if err != nil {
		return 2, fmt.Errorf("create invocation scratch: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(temporary); err != nil {
			result = errors.Join(result, fmt.Errorf("remove invocation scratch %q: %w", temporary, err))
			if code == 0 {
				code = 2
			}
		}
	}()
	mirror, err := call.providerMirror(config, temporary)
	if err != nil {
		return 2, fmt.Errorf("prepare provider mirror: %w", err)
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
		ctx := call.context
		if ctx == nil {
			ctx = context.Background()
		}
		cmd := exec.CommandContext(ctx, terraform, append([]string{"-chdir=" + chdir}, args...)...)
		// Terraform handles SIGINT by stopping provider operations and saving state.
		// Keep the launcher alive until that child exits before its caller cleans up.
		cmd.Cancel = func() error {
			if err := cmd.Process.Signal(os.Interrupt); err != nil {
				return fmt.Errorf("interrupt Terraform: %w", err)
			}
			return nil
		}
		cmd.WaitDelay = 10 * time.Second
		cmd.Env = environ
		cmd.Stdin = call.stdin
		cmd.Stdout = stdout
		cmd.Stderr = call.stderr
		if err := call.execute(cmd); err != nil {
			return fmt.Errorf("run Terraform subprocess: %w", err)
		}
		return nil
	}
	if !opts.direct {
		if err := runCommand(append([]string{"init"}, backendArgs...), call.stderr); err != nil {
			return commandExitCode(err), fmt.Errorf("Terraform initialization failed: %w", err)
		}
	}
	if commandUsesProviders(opts.args) {
		if err := verifyProviderSelections(call.runfiles, config, chdir, environ, runCommand); err != nil {
			return 2, fmt.Errorf("verify providers before command: %w", err)
		}
	}
	if len(opts.args) > 0 && opts.args[0] == "init" {
		opts.args = append(opts.args, backendArgs...)
	}
	if err := runCommand(opts.args, call.stdout); err != nil {
		return commandExitCode(err), fmt.Errorf("Terraform command failed: %w", err)
	}
	if len(opts.args) > 0 && opts.args[0] == "init" {
		if err := verifyProviderSelections(call.runfiles, config, chdir, environ, runCommand); err != nil {
			return 2, fmt.Errorf("verify providers after initialization: %w", err)
		}
	}
	return 0, nil
}

// workingDirectory resolves a directory runfile or materializes a test-owned workspace.
func (call invocation) workingDirectory(config configuration, override string) (string, error) {
	if override != "" {
		directory, err := existingDirectory(override)
		if err != nil {
			return "", fmt.Errorf("resolve explicit working directory: %w", err)
		}
		return directory, nil
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
		return "", fmt.Errorf("create test workspace: %w", err)
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
			return "", fmt.Errorf("resolve workspace input %q: %w", file.Runfile, err)
		}
		destination := filepath.Join(workspace, filepath.FromSlash(file.Path))
		if err := os.MkdirAll(filepath.Dir(destination), 0o700); err != nil {
			return "", fmt.Errorf("create parent for %q: %w", file.Path, err)
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
		return "", fmt.Errorf("create working directory %q: %w", workingDirectory, err)
	}
	fmt.Fprintln(call.stderr, "rules_terraform: manifest-only test workspace retained at", workingDirectory)
	return workingDirectory, nil
}

// runfilesDirectories lists existing candidate runfiles roots in resolution order.
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

// providerMirror uses the declared mirror or materializes links to declared archives.
func (call invocation) providerMirror(config configuration, temporary string) (string, error) {
	packages := make(map[string]string, len(config.Providers))
	platform := runtime.GOOS + "_" + runtime.GOARCH
	for _, provider := range config.Providers {
		if err := validateProviderPackage(provider, platform); err != nil {
			return "", fmt.Errorf("validate declared provider %q: %w", provider.Source, err)
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
		return "", fmt.Errorf("create provider mirror: %w", err)
	}
	for _, provider := range config.Providers {
		archive, err := call.runfiles.Rlocation(provider.Archive)
		if err != nil {
			return "", fmt.Errorf("locate declared provider %s: %w", provider.Source, err)
		}
		archive, err = filepath.Abs(archive)
		if err != nil {
			return "", fmt.Errorf("resolve archive for %q: %w", provider.Source, err)
		}
		if err := regularFile(archive); err != nil {
			return "", fmt.Errorf("declared provider %s archive is unavailable: %w", provider.Source, err)
		}
		destination := filepath.Join(mirror, filepath.FromSlash(provider.MirrorPath))
		if err := os.MkdirAll(filepath.Dir(destination), 0o700); err != nil {
			return "", fmt.Errorf("create mirror directory for %q: %w", provider.Source, err)
		}
		if err := os.Symlink(archive, destination); err != nil {
			return "", fmt.Errorf("link provider archive at %q: %w", destination, err)
		}
	}
	return mirror, nil
}

// validateProviderPackage checks an archive declaration against the runtime platform and layout.
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

// validRunfilePath recognizes normalized relative slash-separated keys.
func validRunfilePath(value string) bool {
	return value != "" && value != "." && !strings.HasPrefix(value, "/") && !strings.Contains(value, "\\") && path.Clean(value) == value && value != ".." && !strings.HasPrefix(value, "../")
}

// regularFile requires an existing regular file at the supplied path.
func regularFile(name string) error {
	info, err := os.Stat(name)
	if err != nil {
		return fmt.Errorf("inspect file %q: %w", name, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("%s must be a regular file", name)
	}
	return nil
}

// existingDirectory returns the absolute path to an existing directory.
func existingDirectory(name string) (string, error) {
	info, err := os.Stat(name)
	if err != nil {
		return "", fmt.Errorf("inspect directory %q: %w", name, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s must be a directory", name)
	}
	absolute, err := filepath.Abs(name)
	if err != nil {
		return "", fmt.Errorf("resolve directory %q: %w", name, err)
	}
	return absolute, nil
}

// runtimeTemporaryRoot resolves caller-owned scratch according to Bazel and process settings.
func runtimeTemporaryRoot(environ []string) (string, error) {
	for _, name := range []string{"TEST_TMPDIR", "TMPDIR"} {
		if value := environmentValue(environ, name); value != "" {
			if err := os.MkdirAll(value, 0o700); err != nil {
				return "", fmt.Errorf("create %s directory: %w", name, err)
			}
			absolute, err := filepath.Abs(value)
			if err != nil {
				return "", fmt.Errorf("resolve scratch directory %q: %w", value, err)
			}
			return absolute, nil
		}
	}
	if workspace := environmentValue(environ, "BUILD_WORKSPACE_DIRECTORY"); workspace != "" {
		directory := filepath.Join(workspace, "out", "rules_terraform", "runtime")
		if err := os.MkdirAll(directory, 0o700); err != nil {
			return "", fmt.Errorf("create runtime scratch %q: %w", directory, err)
		}
		return directory, nil
	}
	return "", errors.New("set TMPDIR to a caller-owned scratch directory when running outside Bazel")
}

// environmentValue returns the last assignment for a variable.
func environmentValue(environ []string, name string) string {
	var value string
	for _, entry := range environ {
		if key, candidate, _ := strings.Cut(entry, "="); key == name {
			value = candidate
		}
	}
	return value
}

// replaceEnvironment overrides named variables while retaining unrelated assignments.
func replaceEnvironment(environ, replacements []string) []string {
	keys := make(map[string]bool, len(replacements))
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

// providerEnvironment isolates Terraform CLI configuration and disables inherited provider caching.
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

// hclString quotes a literal filesystem path for a Terraform CLI configuration.
func hclString(value string) string {
	encoded, _ := json.Marshal(value)
	// HCL quoted strings interpret template openers; filesystem paths do not.
	return strings.ReplaceAll(strings.ReplaceAll(string(encoded), "${", "$${"), "%{", "%%{")
}

// commandExitCode recovers a positive child exit status through wrapped errors.
func commandExitCode(err error) int {
	var exitError *exec.ExitError
	if errors.As(err, &exitError) && exitError.ExitCode() > 0 {
		return exitError.ExitCode()
	}
	return 1
}

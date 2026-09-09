package al

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"syscall"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"github.com/yuin/gopher-lua"
)

const (
	toolCacheSchema   = "git.alwaldend.com/alwaldend/src/projects/al/tool-cache/v1alpha1"
	defaultToolMethod = "bazel"
)

type ToolOptions struct {
	Arguments      []string
	CacheRoot      string
	ConfigPaths    []string
	Environment    []string
	Stderr         io.Writer
	Stdin          io.Reader
	Stdout         io.Writer
	WorkspaceRoot  string
	ReplaceProcess func(string, []string, []string) error
}

type toolMetadata struct {
	Schema string `json:"schema"`
	Name   string `json:"name"`
	Digest string `json:"digest"`
	Method string `json:"method"`
	Kind   string `json:"kind"`
}

func findWorkspaceRoot(explicit string) (string, error) {
	candidate := explicit
	if candidate == "" {
		candidate = os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	}
	absolute, err := filepath.Abs(candidate)
	if err != nil {
		return "", fmt.Errorf("resolve workspace root: %w", err)
	}
	absolute, err = filepath.EvalSymlinks(absolute)
	if err != nil {
		return "", fmt.Errorf("canonicalize workspace root: %w", err)
	}
	for current := absolute; ; current = filepath.Dir(current) {
		if regularFile(filepath.Join(current, "MODULE.bazel")) && regularFile(filepath.Join(current, ".bazelrc")) {
			return current, nil
		}
		parent := filepath.Dir(current)
		if parent == current || explicit != "" {
			break
		}
	}
	return "", fmt.Errorf("%s is not inside a Bazel workspace", absolute)
}

func regularFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

func resolveToolCacheRoot(explicit string, environment []string) (string, error) {
	candidate := explicit
	if candidate == "" {
		candidate = environmentValue(environment, "AL_TOOL_CACHE")
	}
	if candidate == "" {
		candidate = environmentValue(environment, "XDG_CACHE_HOME")
		if candidate == "" {
			home := environmentValue(environment, "HOME")
			if home == "" {
				return "", fmt.Errorf("HOME is unset and no tool cache root was supplied")
			}
			candidate = filepath.Join(home, ".cache")
		}
		candidate = filepath.Join(candidate, "al", "tools")
	}
	absolute, err := filepath.Abs(candidate)
	if err != nil {
		return "", fmt.Errorf("resolve tool cache root: %w", err)
	}
	return filepath.Clean(absolute), nil
}

func environmentValue(environment []string, name string) string {
	prefix := name + "="
	for _, entry := range environment {
		if value, found := strings.CutPrefix(entry, prefix); found {
			return value
		}
	}
	return ""
}

func ToolByName(config *al_proto.Config, name string) (*al_proto.Tool, error) {
	matches := []*al_proto.Tool{}
	for _, candidate := range config.GetTools() {
		if candidate.GetName() == name {
			matches = append(matches, candidate)
		}
	}
	if len(matches) != 1 {
		if len(matches) == 0 {
			return nil, fmt.Errorf("tool %q is not declared", name)
		}
		return nil, fmt.Errorf("tool %q is declared %d times", name, len(matches))
	}
	return matches[0], nil
}

func ToolMethodByName(config *al_proto.Config, name string) (*al_proto.ToolMethod, error) {
	matches := []*al_proto.ToolMethod{}
	for _, candidate := range config.GetToolMethods() {
		if candidate.GetName() == name {
			matches = append(matches, candidate)
		}
	}
	if len(matches) != 1 {
		if len(matches) == 0 {
			return nil, fmt.Errorf("tool method %q is not declared", name)
		}
		return nil, fmt.Errorf("tool method %q is declared %d times", name, len(matches))
	}
	if matches[0].GetKind() != defaultToolMethod {
		return nil, fmt.Errorf("unsupported tool method %q", matches[0].GetKind())
	}
	if matches[0].GetBazel() == nil || matches[0].GetBazel().GetTarget() == "" || matches[0].GetBazel().GetOutput() == "" {
		return nil, fmt.Errorf("bazel tool method %q requires target and output", name)
	}
	return matches[0], nil
}

func toolMethodDigest(workspace string, method *al_proto.ToolMethod, environment []string) (string, error) {
	hash := sha256.New()
	writeHashField(hash, "schema", toolCacheSchema)
	writeHashField(hash, "platform", runtime.GOOS+"/"+runtime.GOARCH)
	writeHashField(hash, "workspace", workspace)
	writeHashField(hash, "kind", method.GetKind())
	writeHashField(hash, "target", method.GetBazel().GetTarget())
	writeHashField(hash, "output", method.GetBazel().GetOutput())
	for _, assignment := range environment {
		if key, _, found := strings.Cut(assignment, "="); found && strings.HasPrefix(key, "AL_") {
			writeHashField(hash, "environment", assignment)
		}
	}
	relativePaths := []string{
		"al.lua",
		"al_lib.lua",
		".bazeliskrc",
		".bazelrc",
		"MODULE.bazel",
		"MODULE.bazel.lock",
		"user.bazelrc",
	}
	sort.Strings(relativePaths)
	for _, relative := range relativePaths {
		path := filepath.Join(workspace, relative)
		if err := hashPath(hash, path, "workspace/"+relative); err != nil {
			if errors.Is(err, fs.ErrNotExist) && (relative == "user.bazelrc" || relative == ".bazeliskrc" || relative == "MODULE.bazel.lock" || relative == "al_lib.lua") {
				writeHashField(hash, "absent", "workspace/"+relative)
				continue
			}
			return "", err
		}
	}
	for _, relative := range []string{filepath.Join(environmentValue(environment, "HOME"), ".bazelrc"), "/etc/bazel.bazelrc"} {
		if err := hashPath(hash, relative, "host/"+filepath.Base(relative)); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				writeHashField(hash, "absent", "host/"+filepath.Base(relative))
				continue
			}
			return "", err
		}
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func writeHashField(hash io.Writer, name string, value string) {
	fmt.Fprintln(hash, name+"\x00"+value)
}

func hashPath(hash io.Writer, path string, label string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(path)
		if err != nil {
			return err
		}
		writeHashField(hash, "symlink:"+label, target)
		return nil
	}
	if info.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		sort.Slice(entries, func(first, second int) bool {
			return entries[first].Name() < entries[second].Name()
		})
		for _, entry := range entries {
			if err := hashPath(hash, filepath.Join(path, entry.Name()), label+"/"+entry.Name()); err != nil {
				return err
			}
		}
		return nil
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("unsupported file mode %s", info.Mode())
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(hash, file)
	return err
}

func validToolCacheEntry(entry string, digest string) (string, bool) {
	content, err := os.ReadFile(filepath.Join(entry, "metadata.json"))
	if err != nil {
		return "", false
	}
	var metadata toolMetadata
	if json.Unmarshal(content, &metadata) != nil ||
		metadata.Schema != toolCacheSchema ||
		metadata.Digest != digest ||
		metadata.Kind != defaultToolMethod {
		return "", false
	}
	executable := filepath.Join(entry, "bin")
	info, err := os.Stat(executable)
	if err != nil || !info.Mode().IsRegular() || info.Mode().Perm()&0o111 == 0 {
		return "", false
	}
	return executable, true
}

func ensurePrivateDirectory(path string) error {
	if err := os.MkdirAll(path, 0o700); err != nil {
		return fmt.Errorf("create tool cache directory: %w", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("inspect tool cache directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("tool cache root %s is not a directory", path)
	}
	if info.Mode().Perm()&0o022 != 0 {
		return fmt.Errorf("tool cache root %s must not be writable by group or others", path)
	}
	return nil
}

func ensureToolCacheEntry(
	stderr io.Writer,
	workspace string,
	cacheRoot string,
	method *al_proto.ToolMethod,
	environment []string,
	build func(string, *al_proto.ToolMethod, []string) error,
) (string, string, bool, error) {
	for attempt := 0; attempt < 3; attempt++ {
		executable, digest, installed, err := ensureToolCacheEntryAttempt(
			stderr, workspace, cacheRoot, method, environment, build,
		)
		if !errors.Is(err, errToolInputsChanged) {
			return executable, digest, installed, err
		}
	}
	return "", "", false, fmt.Errorf("%w after three attempts", errToolInputsChanged)
}

var errToolInputsChanged = errors.New("tool inputs changed during cache installation")

func ensureToolCacheEntryAttempt(
	stderr io.Writer,
	workspace string,
	cacheRoot string,
	method *al_proto.ToolMethod,
	environment []string,
	build func(string, *al_proto.ToolMethod, []string) error,
) (string, string, bool, error) {
	if err := ensurePrivateDirectory(cacheRoot); err != nil {
		return "", "", false, err
	}
	digest, err := toolMethodDigest(workspace, method, environment)
	if err != nil {
		return "", "", false, fmt.Errorf("hash tool inputs: %w", err)
	}
	entry := filepath.Join(cacheRoot, digest)
	if executable, ok := validToolCacheEntry(entry, digest); ok {
		return executable, digest, false, nil
	}
	if err := os.MkdirAll(filepath.Join(cacheRoot, "locks"), 0o700); err != nil {
		return "", "", false, fmt.Errorf("create tool cache lock directory: %w", err)
	}
	lock, err := os.OpenFile(filepath.Join(cacheRoot, "locks", digest+".lock"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return "", "", false, fmt.Errorf("open tool cache lock: %w", err)
	}
	defer lock.Close()
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX); err != nil {
		return "", "", false, fmt.Errorf("lock tool cache entry: %w", err)
	}
	defer syscall.Flock(int(lock.Fd()), syscall.LOCK_UN) //nolint:errcheck
	lockedDigest, err := toolMethodDigest(workspace, method, environment)
	if err != nil {
		return "", "", false, fmt.Errorf("rehash tool inputs: %w", err)
	}
	if lockedDigest != digest {
		return "", "", false, errToolInputsChanged
	}
	if executable, ok := validToolCacheEntry(entry, digest); ok {
		return executable, digest, false, nil
	}
	fmt.Fprintf(stderr, "al: tool cache miss sha256:%s; building %s\n", digest, method.GetBazel().GetTarget())
	if err := build(workspace, method, environment); err != nil {
		return "", "", false, fmt.Errorf("build tool %s: %w", method.GetName(), err)
	}
	builtDigest, err := toolMethodDigest(workspace, method, environment)
	if err != nil {
		return "", "", false, fmt.Errorf("rehash built tool inputs: %w", err)
	}
	if builtDigest != digest {
		return "", "", false, errToolInputsChanged
	}
	if err := installToolCacheEntry(workspace, cacheRoot, entry, method, digest); err != nil {
		return "", "", false, fmt.Errorf("install tool %s: %w", method.GetName(), err)
	}
	executable, ok := validToolCacheEntry(entry, digest)
	if !ok {
		return "", "", false, fmt.Errorf("installed cache entry failed validation")
	}
	return executable, digest, true, nil
}

func installToolCacheEntry(
	workspace string,
	cacheRoot string,
	entry string,
	method *al_proto.ToolMethod,
	digest string,
) error {
	if err := os.MkdirAll(filepath.Dir(entry), 0o700); err != nil {
		return err
	}
	if _, ok := validToolCacheEntry(entry, digest); !ok {
		if err := os.RemoveAll(entry); err != nil {
			return fmt.Errorf("remove incomplete cache entry: %w", err)
		}
	}
	temporary, err := os.MkdirTemp(filepath.Dir(entry), "."+digest+".tmp-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(temporary) //nolint:errcheck
	source := filepath.Join(workspace, filepath.FromSlash(method.GetBazel().GetOutput()))
	destination := filepath.Join(temporary, "bin")
	if err := copyFile(source, destination); err != nil {
		return fmt.Errorf("copy executable from %s: %w", source, err)
	}
	for _, suffix := range []string{".runfiles_manifest", ".repo_mapping"} {
		runfilesSource := source + suffix
		runfilesDestination := destination + suffix
		if _, err := os.Stat(runfilesSource); err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return fmt.Errorf("inspect tool runfiles from %s: %w", runfilesSource, err)
			}
			continue
		}
		if err := copyFile(runfilesSource, runfilesDestination); err != nil {
			return fmt.Errorf("copy tool runfiles from %s: %w", runfilesSource, err)
		}
	}
	metadata := toolMetadata{
		Schema: toolCacheSchema,
		Name:   method.GetName(),
		Digest: digest,
		Method: method.GetKind(),
		Kind:   method.GetKind(),
	}
	content, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}
	content = append(content, '\n')
	if err := os.WriteFile(filepath.Join(temporary, "metadata.json"), content, 0o400); err != nil {
		return err
	}
	if err := os.Rename(temporary, entry); err != nil {
		return err
	}
	return nil
}

func copyFile(source string, destination string) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()
	info, err := input.Stat()
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("not a regular file")
	}
	output, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, info.Mode().Perm())
	if err != nil {
		return err
	}
	if _, err := io.Copy(output, input); err != nil {
		output.Close()
		return err
	}
	return output.Close()
}

func buildToolWithBazel(workspace string, method *al_proto.ToolMethod, environment []string) error {
	bazelPath, err := exec.LookPath("bazel")
	if err != nil {
		return fmt.Errorf("find bazel in PATH: %w", err)
	}
	scratchRoot := filepath.Join(workspace, "out", "al-tool-cache", "builds")
	if err := os.MkdirAll(scratchRoot, 0o700); err != nil {
		return fmt.Errorf("create tool build scratch: %w", err)
	}
	scratch, err := os.MkdirTemp(scratchRoot, method.GetName()+"-")
	if err != nil {
		return fmt.Errorf("create tool build run: %w", err)
	}
	defer os.RemoveAll(scratch) //nolint:errcheck
	command := exec.Command(bazelPath, "build", "--config=agent", method.GetBazel().GetTarget())
	command.Dir = workspace
	command.Env = environment
	command.Stdout = os.Stderr
	command.Stderr = os.Stderr
	return command.Run()
}

func RunTool(ctx context.Context, options *ToolOptions) error {
	if options == nil {
		options = &ToolOptions{}
	}
	if len(options.Arguments) == 0 {
		return fmt.Errorf("tool name is required")
	}
	arguments := options.Arguments
	toolName := arguments[0]
	toolArguments := arguments[1:]
	workspace, err := findWorkspaceRoot(options.WorkspaceRoot)
	if err != nil {
		return err
	}
	configPaths := options.ConfigPaths
	if len(configPaths) == 0 {
		configPaths = []string{filepath.Join(workspace, "al.lua")}
	}
	config, err := LoadConfigs(ctx, configPaths...)
	if err != nil {
		return fmt.Errorf("could not load tool configuration: %w", err)
	}
	tool, err := ToolByName(config, toolName)
	if err != nil {
		return err
	}
	method, err := ToolMethodByName(config, tool.GetMethod())
	if err != nil {
		return err
	}
	environment := options.Environment
	if environment == nil {
		environment = os.Environ()
	}
	environment = append(environment, "BUILD_WORKSPACE_DIRECTORY="+workspace)
	cacheRoot, err := resolveToolCacheRoot(options.CacheRoot, environment)
	if err != nil {
		return err
	}
	if options.Stderr == nil {
		options.Stderr = os.Stderr
	}
	if options.Stdin == nil {
		options.Stdin = os.Stdin
	}
	if options.Stdout == nil {
		options.Stdout = os.Stdout
	}
	if options.ReplaceProcess == nil {
		options.ReplaceProcess = func(executable string, arguments []string, environment []string) error {
			command := exec.Command(executable, arguments[1:]...)
			command.Dir = workspace
			command.Env = environment
			command.Stdin = options.Stdin
			command.Stdout = options.Stdout
			command.Stderr = options.Stderr
			return command.Run()
		}
	}
	executable, _, _, err := ensureToolCacheEntry(
		options.Stderr,
		workspace,
		cacheRoot,
		method,
		environment,
		buildToolWithBazel,
	)
	if err != nil {
		return err
	}
	return options.ReplaceProcess(executable, append([]string{executable}, toolArguments...), environment)
}

func parseToolMethod(state *lua.LState) (*al_proto.ToolMethod, error) {
	method := &al_proto.ToolMethod{}
	if err := parseTableArg(state, method); err != nil {
		return nil, fmt.Errorf("could not parse tool_method: %w", err)
	}
	if method.GetName() == "" {
		return nil, fmt.Errorf("tool_method requires name")
	}
	if method.GetKind() == "" {
		return nil, fmt.Errorf("tool_method requires kind")
	}
	if method.GetBazel() == nil {
		return nil, fmt.Errorf("tool_method requires options for kind %q", method.GetKind())
	}
	if method.GetBazel().GetTarget() == "" || method.GetBazel().GetOutput() == "" {
		return nil, fmt.Errorf("tool_method bazel requires target and output")
	}
	return method, nil
}

func parseTool(state *lua.LState) (*al_proto.Tool, error) {
	tool := &al_proto.Tool{}
	if err := parseTableArg(state, tool); err != nil {
		return nil, fmt.Errorf("could not parse tool: %w", err)
	}
	if tool.GetName() == "" || tool.GetMethod() == "" {
		return nil, fmt.Errorf("tool requires name and method")
	}
	return tool, nil
}

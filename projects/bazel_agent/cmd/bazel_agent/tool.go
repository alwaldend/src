package main

import (
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
	"strconv"
	"strings"
	"syscall"
)

const (
	toolCacheSchema = "agents.alwaldend.com/bazel-agent-tool-cache/v1alpha1"
	sharedToolCache = "/var/cache/bazel/tool_cache"
)

var errToolInputsChanged = errors.New("tool inputs changed during cache installation")

type (
	cachedToolSpec struct {
		Name        string
		Target      string
		Output      string
		Entrypoint  string
		Bundle      bool
		SourcePaths []string
		Environment []string
	}
	toolCacheRecord struct {
		Schema     string `json:"schema"`
		Name       string `json:"name"`
		Digest     string `json:"digest"`
		Target     string `json:"target"`
		Entrypoint string `json:"entrypoint"`
	}
	toolOptions struct {
		WorkspaceRoot string
		CacheRoot     string
	}
	bazelRCFile struct {
		Path     string
		Label    string
		Optional bool
	}
	bazelRCImport struct {
		Path     string
		Optional bool
	}
	toolBuildFunc func(string, string, cachedToolSpec, []string) error
)

var cachedToolSpecs = map[string]cachedToolSpec{
	"mcp_cordis": {
		Name:       "mcp_cordis",
		Target:     "//projects/mcp_cordis:mcp_cordis",
		Output:     "bazel-bin/projects/mcp_cordis/mcp_cordis_",
		Entrypoint: "mcp_cordis",
		Bundle:     true,
		SourcePaths: []string{
			"projects/mcp_cordis/BUILD.bazel",
			"projects/mcp_cordis/include.MODULE.bazel",
			"projects/mcp_cordis/package.json",
			"projects/mcp_cordis/pnpm-lock.yaml",
			"projects/mcp_cordis/pnpm-workspace.yaml",
			"projects/mcp_cordis/cmd/mcp_cordis/main.mjs",
			"projects/mcp_cordis/internal",
			"projects/mcp_cordis/patches",
		},
		Environment: []string{"BAZEL_BINDIR=."},
	},
	"repo_delivery": {
		Name:       "repo_delivery",
		Target:     "//tools/repo_delivery:repo_delivery",
		Output:     "bazel-bin/tools/repo_delivery/cmd/repo_delivery/go_/go",
		Entrypoint: "repo_delivery",
		SourcePaths: []string{
			"tools/repo_delivery/BUILD.bazel",
			"tools/repo_delivery/cmd/repo_delivery",
		},
	},
}

var commonToolSourcePaths = []string{
	".bazeliskrc",
	".bazelrc",
	"MODULE.bazel",
	"MODULE.bazel.lock",
	"tools/bazelrc",
}

func toolUsage() string {
	return strings.Join([]string{
		"Usage:",
		"  bazel_agent tool run [options] TOOL [--] [args...]",
		"  bazel_agent tool warm [options] [TOOL...]",
		"  bazel_agent tool path [options] TOOL",
		"",
		"Options:",
		"  --workspace-root PATH",
		"  --cache-root PATH",
	}, "\n")
}

func parseToolOption(args []string, index int, options *toolOptions) (int, error) {
	if index+1 >= len(args) {
		return index, fmt.Errorf("%s requires a value", args[index])
	}
	value := args[index+1]
	switch args[index] {
	case "--workspace-root":
		options.WorkspaceRoot = value
	case "--cache-root":
		options.CacheRoot = value
	default:
		return index, fmt.Errorf("unknown tool option %q", args[index])
	}
	return index + 2, nil
}

func parseToolNames(args []string) (toolOptions, []string, error) {
	var options toolOptions
	index := 0
	for index < len(args) && strings.HasPrefix(args[index], "--") {
		var err error
		index, err = parseToolOption(args, index, &options)
		if err != nil {
			return toolOptions{}, nil, err
		}
	}
	return options, args[index:], nil
}

func parseToolRun(args []string) (toolOptions, string, []string, error) {
	options, remainder, err := parseToolNames(args)
	if err != nil {
		return toolOptions{}, "", nil, err
	}
	if len(remainder) == 0 {
		return toolOptions{}, "", nil, fmt.Errorf("tool run requires a tool name")
	}
	name := remainder[0]
	arguments := remainder[1:]
	if len(arguments) > 0 && arguments[0] == "--" {
		arguments = arguments[1:]
	}
	return options, name, arguments, nil
}

func runTool(
	args []string,
	environment []string,
	replaceProcess replaceProcessFunc,
	build toolBuildFunc,
) error {
	if len(args) == 0 || args[0] == "--help" || args[0] == "help" {
		fmt.Println(toolUsage())
		return nil
	}
	switch args[0] {
	case "run":
		options, name, toolArgs, err := parseToolRun(args[1:])
		if err != nil {
			return err
		}
		executable, _, _, err := ensureCachedTool(
			name,
			options,
			environment,
			build,
		)
		if err != nil {
			return err
		}
		spec, err := cachedTool(name)
		if err != nil {
			return err
		}
		environment = overrideEnvironment(environment, spec.Environment)
		processArgs := append([]string{executable}, toolArgs...)
		if err := replaceProcess(executable, processArgs, environment); err != nil {
			return fmt.Errorf("execute cached tool %s: %w", name, err)
		}
		return nil
	case "warm":
		options, names, err := parseToolNames(args[1:])
		if err != nil {
			return err
		}
		if len(names) == 0 {
			for name := range cachedToolSpecs {
				names = append(names, name)
			}
			sort.Strings(names)
		}
		for _, name := range names {
			_, digest, installed, err := ensureCachedTool(
				name,
				options,
				environment,
				build,
			)
			if err != nil {
				return err
			}
			state := "cached"
			if installed {
				state = "installed"
			}
			fmt.Printf("%s: %s sha256:%s\n", name, state, digest)
		}
		return nil
	case "path":
		options, names, err := parseToolNames(args[1:])
		if err != nil {
			return err
		}
		if len(names) != 1 {
			return fmt.Errorf("tool path requires exactly one tool name")
		}
		executable, _, _, err := ensureCachedTool(
			names[0],
			options,
			environment,
			build,
		)
		if err != nil {
			return err
		}
		fmt.Println(executable)
		return nil
	default:
		return fmt.Errorf("unknown tool command %q\n%s", args[0], toolUsage())
	}
}

func environmentValue(environment []string, name string) string {
	prefix := name + "="
	for _, entry := range environment {
		if strings.HasPrefix(entry, prefix) {
			return strings.TrimPrefix(entry, prefix)
		}
	}
	return ""
}

func overrideEnvironment(environment, overrides []string) []string {
	if len(overrides) == 0 {
		return environment
	}
	keys := make(map[string]struct{}, len(overrides))
	for _, override := range overrides {
		key, _, _ := strings.Cut(override, "=")
		keys[key] = struct{}{}
	}
	result := make([]string, 0, len(environment)+len(overrides))
	for _, entry := range environment {
		key, _, _ := strings.Cut(entry, "=")
		if _, overridden := keys[key]; !overridden {
			result = append(result, entry)
		}
	}
	return append(result, overrides...)
}

func workspaceRoot(candidate string) (string, error) {
	explicit := candidate != ""
	if candidate == "" {
		current, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("read working directory: %w", err)
		}
		candidate = current
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
		if regularFile(filepath.Join(current, "MODULE.bazel")) &&
			regularFile(filepath.Join(current, ".bazelrc")) {
			return current, nil
		}
		parent := filepath.Dir(current)
		if parent == current || explicit {
			break
		}
	}
	return "", fmt.Errorf("%s is not inside a Bazel workspace", absolute)
}

func regularFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

func resolveToolCacheRoot(options toolOptions, environment []string) (string, error) {
	candidate := options.CacheRoot
	if candidate == "" {
		candidate = environmentValue(environment, "BAZEL_AGENT_TOOL_CACHE")
	}
	if candidate == "" {
		if info, err := os.Stat(sharedToolCache); err == nil && info.IsDir() {
			candidate = sharedToolCache
		}
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
		candidate = filepath.Join(candidate, "bazel_agent", "tools")
	}
	absolute, err := filepath.Abs(candidate)
	if err != nil {
		return "", fmt.Errorf("resolve tool cache root: %w", err)
	}
	return filepath.Clean(absolute), nil
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
		return fmt.Errorf(
			"tool cache root %s must not be writable by group or others",
			path,
		)
	}
	return nil
}

func cachedTool(name string) (cachedToolSpec, error) {
	spec, ok := cachedToolSpecs[name]
	if !ok {
		names := make([]string, 0, len(cachedToolSpecs))
		for candidate := range cachedToolSpecs {
			names = append(names, candidate)
		}
		sort.Strings(names)
		return cachedToolSpec{}, fmt.Errorf(
			"unknown cached tool %q; supported tools: %s",
			name,
			strings.Join(names, ", "),
		)
	}
	return spec, nil
}

func ensureCachedTool(
	name string,
	options toolOptions,
	environment []string,
	build toolBuildFunc,
) (string, string, bool, error) {
	for attempt := 0; attempt < 3; attempt++ {
		executable, digest, installed, err := ensureCachedToolAttempt(
			name,
			options,
			environment,
			build,
		)
		if !errors.Is(err, errToolInputsChanged) {
			return executable, digest, installed, err
		}
	}
	return "", "", false, fmt.Errorf(
		"%w after three attempts",
		errToolInputsChanged,
	)
}

func ensureCachedToolAttempt(
	name string,
	options toolOptions,
	environment []string,
	build toolBuildFunc,
) (string, string, bool, error) {
	spec, err := cachedTool(name)
	if err != nil {
		return "", "", false, err
	}
	root, err := workspaceRoot(options.WorkspaceRoot)
	if err != nil {
		return "", "", false, err
	}
	cacheRoot, err := resolveToolCacheRoot(options, environment)
	if err != nil {
		return "", "", false, err
	}
	if err := ensurePrivateDirectory(cacheRoot); err != nil {
		return "", "", false, err
	}
	digest, err := toolSourceDigest(root, spec, environment)
	if err != nil {
		return "", "", false, fmt.Errorf("hash %s inputs: %w", name, err)
	}
	entry := filepath.Join(cacheRoot, name, digest)
	if executable, ok := validToolCacheEntry(entry, spec, digest); ok {
		return executable, digest, false, nil
	}
	lockDirectory := filepath.Join(cacheRoot, "locks")
	if err := os.MkdirAll(lockDirectory, 0o700); err != nil {
		return "", "", false, fmt.Errorf("create tool cache lock directory: %w", err)
	}
	lockPath := filepath.Join(lockDirectory, name+"-"+digest+".lock")
	lock, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return "", "", false, fmt.Errorf("open tool cache lock: %w", err)
	}
	defer lock.Close()
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX); err != nil {
		return "", "", false, fmt.Errorf("lock tool cache entry: %w", err)
	}
	defer syscall.Flock(int(lock.Fd()), syscall.LOCK_UN) //nolint:errcheck
	lockedDigest, err := toolSourceDigest(root, spec, environment)
	if err != nil {
		return "", "", false, fmt.Errorf("rehash %s inputs: %w", name, err)
	}
	if lockedDigest != digest {
		return "", "", false, errToolInputsChanged
	}
	if executable, ok := validToolCacheEntry(entry, spec, digest); ok {
		return executable, digest, false, nil
	}
	fmt.Fprintf(
		os.Stderr,
		"bazel_agent: %s cache miss sha256:%s; building %s\n",
		name,
		digest,
		spec.Target,
	)
	if err := build(root, cacheRoot, spec, environment); err != nil {
		return "", "", false, fmt.Errorf("build cached tool %s: %w", name, err)
	}
	builtDigest, err := toolSourceDigest(root, spec, environment)
	if err != nil {
		return "", "", false, fmt.Errorf("rehash built %s inputs: %w", name, err)
	}
	if builtDigest != digest {
		return "", "", false, errToolInputsChanged
	}
	if err := installToolCacheEntry(root, cacheRoot, entry, spec, digest); err != nil {
		return "", "", false, fmt.Errorf("install cached tool %s: %w", name, err)
	}
	executable, ok := validToolCacheEntry(entry, spec, digest)
	if !ok {
		return "", "", false, fmt.Errorf("installed cache entry failed validation")
	}
	return executable, digest, true, nil
}

func toolSourceDigest(
	workspace string,
	spec cachedToolSpec,
	environment []string,
) (string, error) {
	hash := sha256.New()
	writeHashField(hash, "schema", toolCacheSchema)
	writeHashField(hash, "platform", runtime.GOOS+"/"+runtime.GOARCH)
	runner, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("resolve runner executable: %w", err)
	}
	runnerDigest, err := fileDigest(runner)
	if err != nil {
		return "", fmt.Errorf("digest runner executable: %w", err)
	}
	writeHashField(hash, "runner", runnerDigest)
	writeHashField(hash, "name", spec.Name)
	writeHashField(hash, "target", spec.Target)
	writeHashField(hash, "output", spec.Output)
	writeHashField(hash, "entrypoint", spec.Entrypoint)
	for _, assignment := range spec.Environment {
		writeHashField(hash, "environment", assignment)
	}
	paths := append([]string{}, commonToolSourcePaths...)
	paths = append(paths, spec.SourcePaths...)
	paths = append(paths, "user.bazelrc")
	sort.Strings(paths)
	for _, relative := range paths {
		path := filepath.Join(workspace, filepath.FromSlash(relative))
		if err := hashPath(hash, path, "workspace/"+filepath.ToSlash(relative)); err != nil {
			if errors.Is(err, fs.ErrNotExist) && relative == "user.bazelrc" {
				writeHashField(hash, "absent", "workspace/user.bazelrc")
				continue
			}
			return "", err
		}
	}
	home := environmentValue(environment, "HOME")
	hostInputs := []struct {
		label string
		path  string
	}{
		{label: "host/home.bazelrc", path: filepath.Join(home, ".bazelrc")},
		{label: "host/system.bazelrc", path: "/etc/bazel.bazelrc"},
	}
	for _, input := range hostInputs {
		if home == "" && input.label == "host/home.bazelrc" {
			writeHashField(hash, "absent", input.label)
			continue
		}
		if err := hashPath(hash, input.path, input.label); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				writeHashField(hash, "absent", input.label)
				continue
			}
			return "", err
		}
	}
	rcFiles := []bazelRCFile{
		{
			Path:  filepath.Join(workspace, ".bazelrc"),
			Label: "workspace/.bazelrc",
		},
		{
			Path:     filepath.Join(workspace, "user.bazelrc"),
			Label:    "workspace/user.bazelrc",
			Optional: true,
		},
		{
			Path:     filepath.Join(home, ".bazelrc"),
			Label:    "host/home.bazelrc",
			Optional: true,
		},
		{
			Path:     "/etc/bazel.bazelrc",
			Label:    "host/system.bazelrc",
			Optional: true,
		},
	}
	if home == "" {
		rcFiles[2].Path = ""
	}
	for _, path := range strings.Split(
		environmentValue(environment, "BAZELRC"),
		",",
	) {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}
		path = resolveBazelRCImport(path, workspace, home)
		rcFiles = append(rcFiles, bazelRCFile{
			Path:  path,
			Label: bazelRCInputLabel(path, workspace, home),
		})
	}
	if err := hashBazelRCClosure(hash, workspace, home, rcFiles); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func hashBazelRCClosure(
	hash hashWriter,
	workspace string,
	home string,
	roots []bazelRCFile,
) error {
	visited := make(map[string]struct{})
	var visit func(bazelRCFile) error
	visit = func(input bazelRCFile) error {
		if input.Path == "" {
			writeHashField(hash, "bazelrc-absent", input.Label)
			return nil
		}
		path := filepath.Clean(input.Path)
		identity := path
		if resolved, err := filepath.EvalSymlinks(path); err == nil {
			identity = resolved
		}
		if _, ok := visited[identity]; ok {
			return nil
		}
		visited[identity] = struct{}{}
		content, err := os.ReadFile(path)
		if err != nil {
			if input.Optional {
				writeHashField(hash, "bazelrc-unavailable", input.Label)
				return nil
			}
			return fmt.Errorf("read Bazel rc %s: %w", input.Label, err)
		}
		writeHashField(hash, "bazelrc-path", input.Label)
		writeHashField(hash, "bazelrc-content", string(content))
		imports, err := bazelRCImports(content)
		if err != nil {
			return fmt.Errorf("parse Bazel rc %s: %w", input.Label, err)
		}
		for _, imported := range imports {
			importedPath := resolveBazelRCImport(
				imported.Path,
				workspace,
				home,
			)
			if err := visit(bazelRCFile{
				Path:     importedPath,
				Label:    bazelRCInputLabel(importedPath, workspace, home),
				Optional: imported.Optional,
			}); err != nil {
				return err
			}
		}
		return nil
	}
	for _, root := range roots {
		if err := visit(root); err != nil {
			return err
		}
	}
	return nil
}

func bazelRCImports(content []byte) ([]bazelRCImport, error) {
	var imports []bazelRCImport
	for index, line := range strings.Split(string(content), "\n") {
		words, err := splitBazelRCWords(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", index+1, err)
		}
		if len(words) == 0 {
			continue
		}
		var optional bool
		switch words[0] {
		case "import":
			if len(words) != 2 {
				return nil, fmt.Errorf("line %d: invalid import", index+1)
			}
		case "try-import":
			if len(words) != 2 {
				return nil, fmt.Errorf("line %d: invalid try-import", index+1)
			}
			optional = true
		case "try-import-if-bazel-version":
			if len(words) < 3 {
				return nil, fmt.Errorf(
					"line %d: invalid conditional import",
					index+1,
				)
			}
			optional = true
		default:
			continue
		}
		imports = append(imports, bazelRCImport{
			Path:     words[len(words)-1],
			Optional: optional,
		})
	}
	return imports, nil
}

func splitBazelRCWords(line string) ([]string, error) {
	var words []string
	var word strings.Builder
	var quote rune
	escaped := false
	active := false
	flush := func() {
		if active {
			words = append(words, word.String())
			word.Reset()
			active = false
		}
	}
	for _, character := range line {
		if escaped {
			word.WriteRune(character)
			escaped = false
			active = true
			continue
		}
		if quote != 0 {
			if character == quote {
				quote = 0
			} else if character == '\\' && quote == '"' {
				escaped = true
			} else {
				word.WriteRune(character)
			}
			active = true
			continue
		}
		switch character {
		case '\\':
			escaped = true
			active = true
		case '\'', '"':
			quote = character
			active = true
		case ' ', '\t', '\r':
			flush()
		case '#':
			if !active {
				flush()
				return words, nil
			}
			word.WriteRune(character)
		default:
			word.WriteRune(character)
			active = true
		}
	}
	if escaped || quote != 0 {
		return nil, fmt.Errorf("unterminated quote or escape")
	}
	flush()
	return words, nil
}

func resolveBazelRCImport(path, workspace, home string) string {
	const workspacePrefix = "%workspace%/"
	switch {
	case path == "%workspace%":
		return workspace
	case strings.HasPrefix(path, workspacePrefix):
		return filepath.Join(
			workspace,
			filepath.FromSlash(strings.TrimPrefix(path, workspacePrefix)),
		)
	case path == "~":
		return home
	case strings.HasPrefix(path, "~/"):
		return filepath.Join(home, filepath.FromSlash(strings.TrimPrefix(path, "~/")))
	case filepath.IsAbs(path):
		return filepath.Clean(path)
	default:
		return filepath.Join(workspace, filepath.FromSlash(path))
	}
}

func bazelRCInputLabel(path, workspace, home string) string {
	if relative, ok := relativePathWithin(workspace, path); ok {
		return "workspace/" + filepath.ToSlash(relative)
	}
	if home != "" {
		if relative, ok := relativePathWithin(home, path); ok {
			return "host/home/" + filepath.ToSlash(relative)
		}
	}
	return "host/absolute/" + filepath.ToSlash(filepath.Clean(path))
}

func relativePathWithin(root, path string) (string, bool) {
	relative, err := filepath.Rel(root, path)
	if err != nil || relative == ".." || strings.HasPrefix(
		relative,
		".."+string(filepath.Separator),
	) {
		return "", false
	}
	return relative, true
}

type hashWriter interface {
	Write([]byte) (int, error)
}

func writeHashField(hash hashWriter, kind, value string) {
	hash.Write([]byte(kind))  //nolint:errcheck
	hash.Write([]byte{0})     //nolint:errcheck
	hash.Write([]byte(value)) //nolint:errcheck
	hash.Write([]byte{0})     //nolint:errcheck
}

func hashPath(hash hashWriter, path, label string) error {
	return filepath.WalkDir(path, func(current string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(path, current)
		if err != nil {
			return err
		}
		entryLabel := label
		if relative != "." {
			entryLabel = filepath.ToSlash(filepath.Join(label, relative))
		}
		entryInfo, err := os.Lstat(current)
		if err != nil {
			return err
		}
		writeHashField(hash, "path", entryLabel)
		writeHashField(hash, "mode", strconv.FormatUint(uint64(entryInfo.Mode()), 8))
		if entryInfo.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(current)
			if err != nil {
				return err
			}
			writeHashField(hash, "symlink", target)
			return nil
		}
		if entryInfo.Mode().IsRegular() {
			file, err := os.Open(current)
			if err != nil {
				return err
			}
			_, copyErr := io.Copy(hash, file)
			closeErr := file.Close()
			if copyErr != nil {
				return copyErr
			}
			return closeErr
		}
		return nil
	})
}

func bazelBuildTool(
	workspace string,
	_ string,
	spec cachedToolSpec,
	environment []string,
) error {
	bazelPath, err := exec.LookPath("bazel")
	if err != nil {
		return fmt.Errorf("find bazel in PATH: %w", err)
	}
	scratchRoot := filepath.Join(
		workspace,
		"out",
		"bazel-agent-tool-cache",
		"builds",
	)
	if err := os.MkdirAll(scratchRoot, 0o700); err != nil {
		return fmt.Errorf("create tool build scratch: %w", err)
	}
	scratch, err := os.MkdirTemp(scratchRoot, spec.Name+"-")
	if err != nil {
		return fmt.Errorf("create tool build run: %w", err)
	}
	defer os.RemoveAll(scratch) //nolint:errcheck
	launchScript := filepath.Join(scratch, "launch")
	command := exec.Command(bazelPath, bazelToolBuildArguments(launchScript, spec)...)
	command.Dir = workspace
	command.Env = environment
	command.Stdout = os.Stderr
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("execute Bazel: %w", err)
	}
	return nil
}

func bazelToolBuildArguments(launchScript string, spec cachedToolSpec) []string {
	return []string{
		"run",
		"--config=agent",
		"--script_path=" + launchScript,
		spec.Target,
	}
}

func validToolCacheEntry(
	entry string,
	spec cachedToolSpec,
	digest string,
) (string, bool) {
	content, err := os.ReadFile(filepath.Join(entry, "ready.json"))
	if err != nil {
		return "", false
	}
	var record toolCacheRecord
	if json.Unmarshal(content, &record) != nil ||
		record.Schema != toolCacheSchema ||
		record.Name != spec.Name ||
		record.Digest != digest ||
		record.Target != spec.Target ||
		record.Entrypoint != spec.Entrypoint {
		return "", false
	}
	executable := filepath.Join(entry, filepath.FromSlash(spec.Entrypoint))
	info, err := os.Stat(executable)
	if err != nil || !info.Mode().IsRegular() || info.Mode().Perm()&0o111 == 0 {
		return "", false
	}
	return executable, true
}

func installToolCacheEntry(
	workspace string,
	cacheRoot string,
	entry string,
	spec cachedToolSpec,
	digest string,
) error {
	parent := filepath.Dir(entry)
	if err := os.MkdirAll(parent, 0o700); err != nil {
		return err
	}
	if _, ok := validToolCacheEntry(entry, spec, digest); !ok {
		if err := os.RemoveAll(entry); err != nil {
			return fmt.Errorf("remove incomplete cache entry: %w", err)
		}
	}
	temporary, err := os.MkdirTemp(parent, "."+digest+".tmp-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(temporary) //nolint:errcheck
	source := filepath.Join(workspace, filepath.FromSlash(spec.Output))
	if spec.Bundle {
		if err := copyTree(source, temporary); err != nil {
			return fmt.Errorf("copy bundle from %s: %w", source, err)
		}
	} else {
		destination := filepath.Join(temporary, filepath.FromSlash(spec.Entrypoint))
		if err := copyTree(source, destination); err != nil {
			return fmt.Errorf("copy executable from %s: %w", source, err)
		}
	}
	record := toolCacheRecord{
		Schema:     toolCacheSchema,
		Name:       spec.Name,
		Digest:     digest,
		Target:     spec.Target,
		Entrypoint: spec.Entrypoint,
	}
	content, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	content = append(content, '\n')
	if err := os.WriteFile(filepath.Join(temporary, "ready.json"), content, 0o400); err != nil {
		return err
	}
	if err := os.Rename(temporary, entry); err != nil {
		return err
	}
	return nil
}

func copyTree(source, destination string) error {
	info, err := os.Lstat(source)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(source)
		if err != nil {
			return err
		}
		if !filepath.IsAbs(target) {
			return os.Symlink(target, destination)
		}
		resolved, err := filepath.EvalSymlinks(source)
		if err != nil {
			return err
		}
		return copyTree(resolved, destination)
	}
	if info.IsDir() {
		assemblyMode := info.Mode().Perm() | 0o700
		if err := os.Mkdir(destination, assemblyMode); err != nil &&
			!os.IsExist(err) {
			return err
		}
		entries, err := os.ReadDir(source)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if err := copyTree(
				filepath.Join(source, entry.Name()),
				filepath.Join(destination, entry.Name()),
			); err != nil {
				return err
			}
		}
		return os.Chmod(destination, assemblyMode)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("unsupported file mode %s", info.Mode())
	}
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()
	output, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, info.Mode().Perm())
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(output, input)
	closeErr := output.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}

package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type gitRepository struct {
	directory  string
	executable string
	runner     commandRunner
}

var gitUnsetEnvironment = []string{
	"GIT_ALTERNATE_OBJECT_DIRECTORIES",
	"GIT_AUTHOR_DATE",
	"GIT_AUTHOR_EMAIL",
	"GIT_AUTHOR_NAME",
	"GIT_COMMON_DIR",
	"GIT_COMMITTER_DATE",
	"GIT_COMMITTER_EMAIL",
	"GIT_COMMITTER_NAME",
	"GIT_CONFIG",
	"GIT_CONFIG_COUNT",
	"GIT_CONFIG_GLOBAL",
	"GIT_CONFIG_NOSYSTEM",
	"GIT_CONFIG_PARAMETERS",
	"GIT_CONFIG_SYSTEM",
	"GIT_CURL_VERBOSE",
	"GIT_DIR",
	"GIT_GRAFT_FILE",
	"GIT_INDEX_FILE",
	"GIT_NAMESPACE",
	"GIT_NO_LAZY_FETCH",
	"GIT_NO_REPLACE_OBJECTS",
	"GIT_OBJECT_DIRECTORY",
	"GIT_OPTIONAL_LOCKS",
	"GIT_QUARANTINE_PATH",
	"GIT_REPLACE_REF_BASE",
	"GIT_SHALLOW_FILE",
	"GIT_WORK_TREE",
}

var gitUnsetEnvironmentPrefixes = []string{
	"GCM_TRACE",
	"GIT_CONFIG_KEY_",
	"GIT_CONFIG_VALUE_",
	"GIT_TRACE",
}

var gitOperationMarkers = []string{
	"AM_HEAD",
	"BISECT_START",
	"CHERRY_PICK_HEAD",
	"MERGE_HEAD",
	"REVERT_HEAD",
	"rebase-apply",
	"rebase-merge",
	"sequencer",
}

const (
	gitNoLazyFetchEnvironment = "GIT_NO_LAZY_FETCH=1"
	gitNoReplaceEnvironment   = "GIT_NO_REPLACE_OBJECTS=1"
	// A full --debug index record is substantially larger than an ordinary
	// path listing. Keep the scan bounded, but allow repositories with
	// hundreds of thousands of normally sized paths instead of inheriting the
	// generic 64 KiB diagnostic-output ceiling.
	indexFlagInspectionOutputLimit = 64 * 1024 * 1024
)

var gitForcedEnvironment = []string{
	gitNoLazyFetchEnvironment,
	gitNoReplaceEnvironment,
	"GIT_OPTIONAL_LOCKS=0",
	"GIT_GRAFT_FILE=" + os.DevNull,
}

var gitGlobalArguments = []string{
	"-c",
	"core.hooksPath=/dev/null",
	"-c",
	"advice.graftFileDeprecated=false",
}

func hardenedGitEnvironment(environment []string) []string {
	result := make([]string, 0, len(environment)+len(gitForcedEnvironment))
	for _, value := range environment {
		name, _, found := strings.Cut(value, "=")
		if found && (name == "GIT_NO_LAZY_FETCH" ||
			name == "GIT_NO_REPLACE_OBJECTS" ||
			name == "GIT_OPTIONAL_LOCKS" ||
			name == "GIT_REPLACE_REF_BASE" || name == "GIT_GRAFT_FILE") {
			continue
		}
		result = append(result, value)
	}
	return append(result, gitForcedEnvironment...)
}

func hardenedGitArguments(arguments []string) []string {
	result := make([]string, 0, len(gitGlobalArguments)+len(arguments))
	result = append(result, gitGlobalArguments...)
	return append(result, arguments...)
}

func openGitRepository(
	ctx context.Context,
	directory string,
	executable string,
	getenv func(string) string,
	runner commandRunner,
) (*gitRepository, error) {
	if strings.TrimSpace(executable) == "" {
		executable = "git"
	}
	if strings.TrimSpace(directory) == "" {
		directory = getenv("BUILD_WORKSPACE_DIRECTORY")
	}
	if strings.TrimSpace(directory) == "" {
		var err error
		directory, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("get current directory: %w", err)
		}
	}
	result, err := runner.Run(ctx, command{
		Name: executable,
		Args: hardenedGitArguments(
			[]string{"rev-parse", "--show-toplevel"},
		),
		Dir:              directory,
		Env:              hardenedGitEnvironment(nil),
		UnsetEnv:         gitUnsetEnvironment,
		UnsetEnvPrefixes: gitUnsetEnvironmentPrefixes,
	})
	if err != nil {
		return nil, fmt.Errorf("resolve repository root: %w", err)
	}
	root := strings.TrimSpace(result.Stdout)
	if root == "" {
		return nil, fmt.Errorf("Git returned an empty repository root")
	}
	repository := &gitRepository{
		directory:  filepath.Clean(root),
		executable: executable,
		runner:     runner,
	}
	if err := repository.requireNonShallowUngraftedHistory(ctx); err != nil {
		return nil, err
	}
	return repository, nil
}

func (g *gitRepository) run(
	ctx context.Context,
	args ...string,
) (commandResult, error) {
	return g.runWithOutputLimit(ctx, 0, args...)
}

func (g *gitRepository) runWithOutputLimit(
	ctx context.Context,
	outputLimit int,
	args ...string,
) (commandResult, error) {
	return g.runEnvironmentWithOutputLimit(ctx, nil, outputLimit, args...)
}

func (g *gitRepository) runEnvironmentWithOutputLimit(
	ctx context.Context,
	environment []string,
	outputLimit int,
	args ...string,
) (commandResult, error) {
	return g.runner.Run(ctx, command{
		Name:             g.executable,
		Args:             hardenedGitArguments(args),
		Dir:              g.directory,
		Env:              hardenedGitEnvironment(environment),
		UnsetEnv:         gitUnsetEnvironment,
		UnsetEnvPrefixes: gitUnsetEnvironmentPrefixes,
		OutputLimit:      outputLimit,
	})
}

func (g *gitRepository) runEnvironment(
	ctx context.Context,
	environment []string,
	args ...string,
) (commandResult, error) {
	return g.runEnvironmentWithOutputLimit(ctx, environment, 0, args...)
}

func (g *gitRepository) runInput(
	ctx context.Context,
	stdin string,
	args ...string,
) (commandResult, error) {
	return g.runInputEnv(ctx, stdin, nil, args...)
}

func (g *gitRepository) runInputEnv(
	ctx context.Context,
	stdin string,
	environment []string,
	args ...string,
) (commandResult, error) {
	return g.runner.Run(ctx, command{
		Name:             g.executable,
		Args:             hardenedGitArguments(args),
		Dir:              g.directory,
		Env:              hardenedGitEnvironment(environment),
		UnsetEnv:         gitUnsetEnvironment,
		UnsetEnvPrefixes: gitUnsetEnvironmentPrefixes,
		Stdin:            stdin,
	})
}

func (g *gitRepository) text(
	ctx context.Context,
	args ...string,
) (string, error) {
	result, err := g.run(ctx, args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result.Stdout), nil
}

func commandExitCode(err error) (int, bool) {
	var commandErr *commandError
	if !errors.As(err, &commandErr) {
		return 0, false
	}
	return commandErr.Result.ExitCode, true
}

func (g *gitRepository) currentBranch(ctx context.Context) (string, error) {
	if err := g.requireCompleteHistory(ctx); err != nil {
		return "", err
	}
	branch, err := g.text(ctx, "symbolic-ref", "--quiet", "--short", "HEAD")
	if err != nil {
		return "", fmt.Errorf("HEAD is detached or unreadable: %w", err)
	}
	if err := g.checkBranch(ctx, branch); err != nil {
		return "", err
	}
	return branch, nil
}

func (g *gitRepository) requireCompleteHistory(ctx context.Context) error {
	if err := g.requireNoPartialClone(ctx); err != nil {
		return err
	}
	return g.requireNonShallowUngraftedHistory(ctx)
}

func (g *gitRepository) requireNonShallowUngraftedHistory(
	ctx context.Context,
) error {
	value, err := g.text(ctx, "rev-parse", "--is-shallow-repository")
	if err != nil {
		return fmt.Errorf("inspect repository history completeness: %w", err)
	}
	switch value {
	case "false":
		// Continue to legacy graft inspection below.
	case "true":
		return fmt.Errorf(
			"shallow repositories are not supported; fetch complete history before delivery",
		)
	default:
		return fmt.Errorf(
			"Git returned invalid shallow-repository state %q",
			value,
		)
	}
	commonDirectory, err := g.text(ctx, "rev-parse", "--git-common-dir")
	if err != nil {
		return fmt.Errorf("resolve Git common directory: %w", err)
	}
	if !filepath.IsAbs(commonDirectory) {
		commonDirectory = filepath.Join(g.directory, commonDirectory)
	}
	grafts := filepath.Join(commonDirectory, "info", "grafts")
	if _, err := os.Lstat(grafts); err == nil {
		return fmt.Errorf(
			"legacy Git grafts are not supported; remove info/grafts before delivery",
		)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("inspect legacy Git graft file: %w", err)
	}
	return nil
}

func (g *gitRepository) requireNoPartialClone(ctx context.Context) error {
	presencePatterns := []string{
		`^extensions\.partialclone$`,
		`^remote\..*\.partialclonefilter$`,
	}
	for _, scope := range []string{"--local", "--worktree"} {
		if scope == "--worktree" {
			enabled, err := g.worktreeConfigEnabled(ctx)
			if err != nil {
				return err
			}
			if !enabled {
				continue
			}
		}
		for _, pattern := range presencePatterns {
			keys, err := g.configKeys(ctx, scope, pattern)
			if err != nil {
				return fmt.Errorf("inspect partial-clone configuration: %w", err)
			}
			if len(keys) != 0 {
				return fmt.Errorf(
					"partial or promisor repositories are not supported; use a complete local object database before delivery",
				)
			}
		}
		promisor, err := g.configHasTrue(
			ctx,
			scope,
			`^remote\..*\.promisor$`,
		)
		if err != nil {
			return fmt.Errorf("inspect promisor configuration: %w", err)
		}
		if promisor {
			return fmt.Errorf(
				"partial or promisor repositories are not supported; use a complete local object database before delivery",
			)
		}
	}
	return nil
}

func (g *gitRepository) worktreeConfigEnabled(ctx context.Context) (bool, error) {
	value, err := g.text(
		ctx,
		"config",
		"--local",
		"--no-includes",
		"--type=bool",
		"--get",
		"extensions.worktreeConfig",
	)
	if err != nil {
		if code, ok := commandExitCode(err); ok && code == 1 {
			return false, nil
		}
		return false, fmt.Errorf("inspect worktree configuration mode: %w", err)
	}
	return value == "true", nil
}

func (g *gitRepository) configKeys(
	ctx context.Context,
	scope string,
	pattern string,
) ([]string, error) {
	result, err := g.run(
		ctx,
		"config",
		scope,
		"--null",
		"--name-only",
		"--get-regexp",
		pattern,
	)
	if err != nil {
		if code, ok := commandExitCode(err); ok && code == 1 {
			return nil, nil
		}
		return nil, err
	}
	value := strings.TrimSuffix(result.Stdout, "\x00")
	if value == "" {
		return nil, fmt.Errorf("Git reported an empty matching configuration key")
	}
	keys := strings.Split(value, "\x00")
	for _, key := range keys {
		if key == "" || strings.ContainsAny(key, "\r\n") {
			return nil, fmt.Errorf("Git returned malformed configuration keys")
		}
	}
	return keys, nil
}

func (g *gitRepository) configHasTrue(
	ctx context.Context,
	scope string,
	pattern string,
) (bool, error) {
	result, err := g.run(
		ctx,
		"config",
		scope,
		"--type=bool",
		"--get-regexp",
		pattern,
	)
	if err != nil {
		if code, ok := commandExitCode(err); ok && code == 1 {
			return false, nil
		}
		return false, err
	}
	for _, line := range strings.Split(strings.TrimSuffix(result.Stdout, "\n"), "\n") {
		switch {
		case strings.HasSuffix(line, " true"):
			return true, nil
		case strings.HasSuffix(line, " false"):
			continue
		default:
			return false, fmt.Errorf("Git returned malformed Boolean configuration")
		}
	}
	return false, nil
}

func (g *gitRepository) checkBranch(
	ctx context.Context,
	branch string,
) error {
	if strings.TrimSpace(branch) == "" {
		return fmt.Errorf("branch name is empty")
	}
	if _, err := g.run(ctx, "check-ref-format", "--branch", branch); err != nil {
		return fmt.Errorf("invalid branch %q: %w", branch, err)
	}
	return nil
}

func (g *gitRepository) head(ctx context.Context) (string, error) {
	head, err := g.text(ctx, "rev-parse", "HEAD")
	if err != nil {
		return "", fmt.Errorf("resolve HEAD: %w", err)
	}
	if !isObjectID(head) {
		return "", fmt.Errorf("Git returned invalid HEAD object ID %q", head)
	}
	return head, nil
}

func (g *gitRepository) branchHead(
	ctx context.Context,
) (string, string, error) {
	if err := g.requireCompleteHistory(ctx); err != nil {
		return "", "", err
	}
	result, err := g.run(
		ctx,
		"for-each-ref",
		"--format=%(HEAD)%00%(refname)%00%(objectname)",
		"refs/heads",
	)
	if err != nil {
		return "", "", fmt.Errorf("resolve current branch and HEAD: %w", err)
	}
	var branchRef string
	var head string
	for _, line := range strings.Split(strings.TrimSuffix(result.Stdout, "\n"), "\n") {
		fields := strings.Split(line, "\x00")
		if len(fields) != 3 {
			return "", "", fmt.Errorf("Git returned malformed branch metadata")
		}
		if fields[0] != "*" {
			continue
		}
		if branchRef != "" {
			return "", "", fmt.Errorf("Git reported multiple current branches")
		}
		branchRef = fields[1]
		head = fields[2]
	}
	if !strings.HasPrefix(branchRef, "refs/heads/") || !isObjectID(head) {
		return "", "", fmt.Errorf("HEAD is detached or changed while being read")
	}
	branch := strings.TrimPrefix(branchRef, "refs/heads/")
	if err := g.checkBranch(ctx, branch); err != nil {
		return "", "", err
	}
	return branch, head, nil
}

func (g *gitRepository) requireBranchHead(
	ctx context.Context,
	branch string,
	head string,
) error {
	currentBranch, currentHead, err := g.branchHead(ctx)
	if err != nil {
		return err
	}
	if currentBranch != branch || currentHead != head {
		return fmt.Errorf(
			"current branch or HEAD changed: got %s at %s, expected %s at %s",
			currentBranch,
			currentHead,
			branch,
			head,
		)
	}
	return nil
}

func (g *gitRepository) tree(
	ctx context.Context,
	oid string,
) (string, error) {
	if !isObjectID(oid) {
		return "", fmt.Errorf("tree source is not a full Git object ID")
	}
	tree, err := g.text(ctx, "rev-parse", "--verify", oid+"^{tree}")
	if err != nil {
		return "", fmt.Errorf("resolve commit tree: %w", err)
	}
	if !isObjectID(tree) {
		return "", fmt.Errorf("Git returned invalid tree object ID %q", tree)
	}
	return tree, nil
}

var gitRemoteNamePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`)

func validateRemoteName(remote string) error {
	if !gitRemoteNamePattern.MatchString(remote) {
		return fmt.Errorf(
			"remote name must use only ASCII letters, digits, underscores, and hyphens",
		)
	}
	return nil
}

func (g *gitRepository) remoteURLs(
	ctx context.Context,
	remote string,
) (string, string, error) {
	firstFetch, firstPush, err := g.rawRemoteURLs(ctx, remote)
	if err != nil {
		return "", "", err
	}
	secondFetch, secondPush, err := g.rawRemoteURLs(ctx, remote)
	if err != nil {
		return "", "", err
	}
	if firstFetch != secondFetch || firstPush != secondPush {
		return "", "", fmt.Errorf(
			"selected remote endpoints changed while they were captured",
		)
	}
	return firstFetch, firstPush, nil
}

func (g *gitRepository) rawRemoteURLs(
	ctx context.Context,
	remote string,
) (string, string, error) {
	if err := validateRemoteName(remote); err != nil {
		return "", "", err
	}
	urlKey := "remote." + remote + ".url"
	pushURLKey := "remote." + remote + ".pushurl"
	if enabled, err := g.worktreeConfigEnabled(ctx); err != nil {
		return "", "", err
	} else if enabled {
		for _, key := range []string{urlKey, pushURLKey} {
			values, err := g.configValues(ctx, "--worktree", key)
			if err != nil {
				return "", "", fmt.Errorf(
					"inspect selected worktree remote endpoint: %w",
					err,
				)
			}
			if len(values) != 0 {
				return "", "", fmt.Errorf(
					"selected remote endpoints in worktree config are not supported",
				)
			}
		}
	}
	fetchValues, err := g.configValues(ctx, "--local", urlKey)
	if err != nil {
		return "", "", fmt.Errorf("resolve raw fetch endpoint: %w", err)
	}
	if len(fetchValues) != 1 {
		return "", "", fmt.Errorf(
			"selected remote must have exactly one raw fetch endpoint, found %d",
			len(fetchValues),
		)
	}
	pushValues, err := g.configValues(ctx, "--local", pushURLKey)
	if err != nil {
		return "", "", fmt.Errorf("resolve raw push endpoint: %w", err)
	}
	if len(pushValues) > 1 {
		return "", "", fmt.Errorf(
			"selected remote must have at most one raw push endpoint, found %d",
			len(pushValues),
		)
	}
	fetchURL, err := exactlyOneRemoteURL(fetchValues[0], "fetch")
	if err != nil {
		return "", "", err
	}
	pushURL := fetchURL
	if len(pushValues) == 1 {
		pushURL, err = exactlyOneRemoteURL(pushValues[0], "push")
		if err != nil {
			return "", "", err
		}
	}
	return fetchURL, pushURL, nil
}

func (g *gitRepository) configValues(
	ctx context.Context,
	scope string,
	key string,
) ([]string, error) {
	result, err := g.run(
		ctx,
		"config",
		scope,
		"--no-includes",
		"--null",
		"--get-all",
		key,
	)
	if err != nil {
		if code, ok := commandExitCode(err); ok && code == 1 {
			return nil, nil
		}
		return nil, err
	}
	if result.Stdout == "" || !strings.HasSuffix(result.Stdout, "\x00") {
		return nil, fmt.Errorf("Git returned malformed raw config values")
	}
	values := strings.Split(strings.TrimSuffix(result.Stdout, "\x00"), "\x00")
	for _, value := range values {
		if value == "" || strings.ContainsRune(value, '\x00') {
			return nil, fmt.Errorf("Git returned an empty or malformed raw config value")
		}
	}
	return values, nil
}

func exactlyOneRemoteURL(value string, kind string) (string, error) {
	if value == "" || strings.TrimSpace(value) != value ||
		strings.ContainsAny(value, "\x00\r\n") {
		return "", fmt.Errorf(
			"selected remote must have one canonical raw %s endpoint",
			kind,
		)
	}
	return value, nil
}

type refSnapshot struct {
	BaseRef        string
	BaseOID        string
	RemoteBaseRef  string
	RemoteHeadRef  string
	RemoteHeadOID  string
	privateRefs    []string
	privateRefOIDs map[string]string
	repository     *gitRepository
}

type boundGitRemote struct {
	endpoint string
}

type isolatedGitTransport struct {
	repository  *gitRepository
	directory   string
	objectDir   string
	remote      boundGitRemote
	environment []string
}

var gitTransportUnsetEnvironment = []string{
	"ALL_PROXY",
	"CURL_CA_BUNDLE",
	"FTP_PROXY",
	"GIT_ALLOW_PROTOCOL",
	"GIT_ASKPASS",
	"GIT_EXEC_PATH",
	"GIT_HTTP_PROXY_AUTHMETHOD",
	"GIT_PROTOCOL_FROM_USER",
	"GIT_PROXY_COMMAND",
	"GIT_SSH",
	"GIT_SSH_COMMAND",
	"GIT_SSH_VARIANT",
	"GIT_SSL_CAINFO",
	"GIT_SSL_CAPATH",
	"GIT_SSL_CERT",
	"GIT_SSL_CERT_PASSWORD_PROTECTED",
	"GIT_SSL_CIPHER_LIST",
	"GIT_SSL_KEY",
	"GIT_SSL_NO_VERIFY",
	"GIT_SSL_VERSION",
	"HTTP_PROXY",
	"HTTPS_PROXY",
	"NO_PROXY",
	"RSYNC_PROXY",
	"SSL_CERT_DIR",
	"SSL_CERT_FILE",
	"SSLKEYLOGFILE",
	"SSH_ASKPASS",
	"SSH_ASKPASS_REQUIRE",
	"all_proxy",
	"ftp_proxy",
	"http_proxy",
	"https_proxy",
	"no_proxy",
}

const isolatedSSHCommand = "ssh -F /dev/null" +
	" -oBatchMode=yes" +
	" -oCanonicalizeHostname=no" +
	" -oClearAllForwardings=yes" +
	" -oHostName=%h" +
	" -oKnownHostsCommand=none" +
	" -oLocalCommand=none" +
	" -oPermitLocalCommand=no" +
	" -oProxyCommand=none" +
	" -oProxyJump=none" +
	" -oRequestTTY=no"

var (
	canonicalSSHUserPattern = regexp.MustCompile(
		`^[A-Za-z0-9][A-Za-z0-9._-]*$`,
	)
	canonicalForgePathPattern = regexp.MustCompile(
		`^[A-Za-z0-9](?:[A-Za-z0-9._-]*[A-Za-z0-9])?$`,
	)
)

func canonicalSSHHost(host string) bool {
	if host == "" || host != strings.ToLower(host) || len(host) > 253 ||
		strings.HasSuffix(host, ".") {
		return false
	}
	for _, label := range strings.Split(host, ".") {
		if label == "" || len(label) > 63 ||
			!canonicalForgePathPattern.MatchString(label) ||
			strings.ContainsAny(label, "._") {
			return false
		}
	}
	return true
}

func canonicalForgeRepository(repository remoteRepository) bool {
	return canonicalSSHHost(repository.Host) &&
		canonicalForgePathPattern.MatchString(repository.Owner) &&
		canonicalForgePathPattern.MatchString(repository.Name)
}

func requireCanonicalSSHEndpoint(endpoint string) error {
	if endpoint == "" || strings.TrimSpace(endpoint) != endpoint ||
		strings.ContainsAny(endpoint, "\x00\r\n") {
		return fmt.Errorf(
			"repository delivery v1 requires a canonical SSH Git endpoint",
		)
	}
	if strings.HasPrefix(endpoint, "ssh://") {
		parsed, err := url.Parse(endpoint)
		if err != nil || parsed.Scheme != "ssh" || parsed.User == nil ||
			!canonicalSSHUserPattern.MatchString(parsed.User.Username()) ||
			parsed.User.String() != parsed.User.Username() ||
			!canonicalSSHHost(parsed.Hostname()) {
			return fmt.Errorf(
				"repository delivery v1 requires a canonical SSH Git endpoint",
			)
		}
		if _, hasPassword := parsed.User.Password(); hasPassword {
			return fmt.Errorf(
				"repository delivery v1 does not permit credentials in Git endpoints",
			)
		}
		repository, err := parseRemoteRepository(endpoint)
		if err != nil || !canonicalForgeRepository(repository) {
			return fmt.Errorf(
				"repository delivery v1 requires a canonical SSH Git endpoint",
			)
		}
		return nil
	}
	if strings.Contains(endpoint, "://") {
		return fmt.Errorf(
			"repository delivery v1 requires SSH Git endpoints; HTTPS credential helpers are intentionally not loaded",
		)
	}
	colon := strings.IndexByte(endpoint, ':')
	if colon <= 0 || strings.Count(endpoint[:colon], "@") != 1 {
		return fmt.Errorf(
			"repository delivery v1 requires a canonical SSH Git endpoint",
		)
	}
	user, host, found := strings.Cut(endpoint[:colon], "@")
	if !found || !canonicalSSHUserPattern.MatchString(user) ||
		!canonicalSSHHost(host) {
		return fmt.Errorf(
			"repository delivery v1 requires a canonical SSH Git endpoint",
		)
	}
	repository, err := parseRemoteRepository(endpoint)
	if err != nil || !canonicalForgeRepository(repository) {
		return fmt.Errorf(
			"repository delivery v1 requires a canonical SSH Git endpoint",
		)
	}
	return nil
}

func bindGitRemote(endpoint string) (boundGitRemote, error) {
	if filepath.IsAbs(endpoint) {
		if filepath.Clean(endpoint) != endpoint ||
			strings.ContainsAny(endpoint, "\x00\r\n") {
			return boundGitRemote{}, fmt.Errorf(
				"local test transport endpoint is not canonical",
			)
		}
		return boundGitRemote{endpoint: endpoint}, nil
	}
	if err := requireCanonicalSSHEndpoint(endpoint); err != nil {
		return boundGitRemote{}, err
	}
	return boundGitRemote{endpoint: endpoint}, nil
}

func transportConfiguration(objectFormat string) ([]string, error) {
	repositoryFormat := "0"
	configuration := [][2]string{
		{"core.bare", "true"},
		{"core.gitProxy", "none"},
		{"core.hooksPath", os.DevNull},
		{"core.repositoryFormatVersion", repositoryFormat},
		{"core.sshCommand", isolatedSSHCommand},
		{"ssh.variant", "ssh"},
		{"http.followRedirects", "false"},
		{"http.proxy", ""},
		{"http.sslVerify", "true"},
		{"protocol.allow", "never"},
		{"protocol.ext.allow", "never"},
		{"protocol.file.allow", "always"},
		{"protocol.ssh.allow", "always"},
	}
	switch objectFormat {
	case "sha1":
	case "sha256":
		configuration[3][1] = "1"
		configuration = append(
			configuration,
			[2]string{"extensions.objectFormat", "sha256"},
		)
	default:
		return nil, fmt.Errorf("unsupported Git object format")
	}
	environment := []string{
		"GIT_CONFIG_COUNT=" + strconv.Itoa(len(configuration)),
		"GIT_CONFIG_GLOBAL=" + os.DevNull,
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_SYSTEM=" + os.DevNull,
		"GIT_TERMINAL_PROMPT=0",
	}
	for index, entry := range configuration {
		environment = append(
			environment,
			"GIT_CONFIG_KEY_"+strconv.Itoa(index)+"="+entry[0],
			"GIT_CONFIG_VALUE_"+strconv.Itoa(index)+"="+entry[1],
		)
	}
	return environment, nil
}

func (g *gitRepository) openIsolatedTransport(
	ctx context.Context,
	remote boundGitRemote,
) (*isolatedGitTransport, error) {
	if remote.endpoint == "" {
		return nil, fmt.Errorf("bound remote endpoint is empty")
	}
	objectDir, err := g.text(
		ctx,
		"rev-parse",
		"--path-format=absolute",
		"--git-path",
		"objects",
	)
	if err != nil {
		return nil, fmt.Errorf("resolve repository object directory: %w", err)
	}
	objectDir, err = canonicalExistingPath(objectDir)
	if err != nil {
		return nil, fmt.Errorf("canonicalize repository object directory: %w", err)
	}
	objectInfo, err := os.Stat(objectDir)
	if err != nil || !objectInfo.IsDir() {
		if err != nil {
			return nil, fmt.Errorf("inspect repository object directory: %w", err)
		}
		return nil, fmt.Errorf("repository object path is not a directory")
	}
	objectFormat, err := g.text(ctx, "rev-parse", "--show-object-format")
	if err != nil {
		return nil, fmt.Errorf("resolve repository object format: %w", err)
	}
	environment, err := transportConfiguration(objectFormat)
	if err != nil {
		return nil, err
	}
	suffix, err := randomHex(12)
	if err != nil {
		return nil, fmt.Errorf("create isolated transport suffix: %w", err)
	}
	relativeRoot := filepath.Join("out", "repo_delivery")
	relativePath := filepath.Join(relativeRoot, "transport-"+suffix)
	ignored, err := g.pathIgnored(ctx, filepath.ToSlash(relativePath))
	if err != nil {
		return nil, err
	}
	if !ignored {
		return nil, fmt.Errorf(
			"isolated transport path %q is not ignored by Git",
			relativePath,
		)
	}
	if err := ensurePrivateDirectory(filepath.Join(g.directory, "out")); err != nil {
		return nil, err
	}
	root := filepath.Join(g.directory, relativeRoot)
	if err := ensurePrivateDirectory(root); err != nil {
		return nil, err
	}
	if err := requireMode0700Directory(root); err != nil {
		return nil, err
	}
	directory := filepath.Join(g.directory, relativePath)
	if err := os.Mkdir(directory, 0o700); err != nil {
		return nil, fmt.Errorf("create isolated transport directory: %w", err)
	}
	transport := &isolatedGitTransport{
		repository: g,
		directory:  directory,
		objectDir:  objectDir,
		remote:     remote,
	}
	failed := true
	defer func() {
		if failed {
			_ = transport.close()
		}
	}()
	for _, relative := range []string{
		"objects",
		filepath.Join("objects", "info"),
		filepath.Join("objects", "pack"),
		"refs",
		filepath.Join("refs", "heads"),
		filepath.Join("refs", "tags"),
	} {
		if err := os.Mkdir(filepath.Join(directory, relative), 0o700); err != nil {
			return nil, fmt.Errorf("create isolated bare Git directory: %w", err)
		}
	}
	if err := writeExclusivePrivateFile(
		filepath.Join(directory, "HEAD"),
		[]byte("ref: refs/heads/repo-delivery\n"),
	); err != nil {
		return nil, fmt.Errorf("create isolated bare Git HEAD: %w", err)
	}
	environment = append(
		environment,
		"GIT_DIR="+directory,
		"GIT_OBJECT_DIRECTORY="+objectDir,
	)
	transport.environment = environment
	failed = false
	return transport, nil
}

func requireMode0700Directory(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("inspect private directory %q: %w", path, err)
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 ||
		info.Mode().Perm() != 0o700 {
		return fmt.Errorf("private path %q must be a real mode-0700 directory", path)
	}
	return nil
}

func writeExclusivePrivateFile(path string, contents []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	if _, err := file.Write(contents); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}

func (t *isolatedGitTransport) run(
	ctx context.Context,
	args ...string,
) (commandResult, error) {
	unset := make([]string, 0, len(gitUnsetEnvironment)+
		len(gitTransportUnsetEnvironment))
	unset = append(unset, gitUnsetEnvironment...)
	unset = append(unset, gitTransportUnsetEnvironment...)
	return t.repository.runner.Run(ctx, command{
		Name:             t.repository.executable,
		Args:             hardenedGitArguments(args),
		Dir:              t.directory,
		Env:              hardenedGitEnvironment(t.environment),
		UnsetEnv:         unset,
		UnsetEnvPrefixes: gitUnsetEnvironmentPrefixes,
	})
}

func (t *isolatedGitTransport) text(
	ctx context.Context,
	args ...string,
) (string, error) {
	result, err := t.run(ctx, args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result.Stdout), nil
}

func (t *isolatedGitTransport) close() error {
	root := filepath.Join(t.repository.directory, "out", "repo_delivery")
	relative, err := filepath.Rel(root, t.directory)
	if err != nil || !strings.HasPrefix(relative, "transport-") ||
		strings.ContainsRune(relative, filepath.Separator) {
		return fmt.Errorf("refuse unsafe isolated transport path %q", t.directory)
	}
	info, err := os.Lstat(t.directory)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect isolated transport directory: %w", err)
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 ||
		info.Mode().Perm() != 0o700 {
		return fmt.Errorf("isolated transport path is not an owned mode-0700 directory")
	}
	if err := os.RemoveAll(t.directory); err != nil {
		return fmt.Errorf("remove isolated transport directory: %w", err)
	}
	if _, err := os.Lstat(t.directory); err == nil {
		return fmt.Errorf("isolated transport directory still exists after removal")
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("inspect removed isolated transport directory: %w", err)
	}
	return nil
}

func (t *isolatedGitTransport) remoteObjectIDs(
	ctx context.Context,
	refs ...string,
) (map[string]string, error) {
	wanted := make(map[string]bool, len(refs))
	args := []string{"ls-remote", "--heads", "--", t.remote.endpoint}
	for _, ref := range refs {
		if !strings.HasPrefix(ref, "refs/heads/") || wanted[ref] {
			return nil, fmt.Errorf("invalid or duplicate exact remote ref %q", ref)
		}
		wanted[ref] = true
		args = append(args, ref)
	}
	result, err := t.run(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("inspect exact remote refs: %w", err)
	}
	values := make(map[string]string, len(refs))
	output := strings.TrimSpace(result.Stdout)
	if output == "" {
		return values, nil
	}
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 || !wanted[fields[1]] ||
			!isObjectID(fields[0]) || values[fields[1]] != "" {
			return nil, fmt.Errorf("remote returned malformed exact-ref data")
		}
		values[fields[1]] = fields[0]
	}
	return values, nil
}

func (g *gitRepository) remoteObjectIDs(
	ctx context.Context,
	remote boundGitRemote,
	refs ...string,
) (values map[string]string, returnErr error) {
	transport, err := g.openIsolatedTransport(ctx, remote)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cleanupErr := transport.close(); cleanupErr != nil {
			returnErr = errors.Join(
				returnErr,
				fmt.Errorf("clean isolated Git transport: %w", cleanupErr),
			)
		}
	}()
	return transport.remoteObjectIDs(ctx, refs...)
}

func (g *gitRepository) installSnapshotRef(
	ctx context.Context,
	ref string,
	oid string,
) (bool, error) {
	if !strings.HasPrefix(ref, "refs/repo-delivery/") || !isObjectID(oid) {
		return false, fmt.Errorf("invalid snapshot ref installation")
	}
	zeroOID := strings.Repeat("0", len(oid))
	if _, err := g.run(
		ctx,
		"update-ref",
		"--no-deref",
		"-m",
		"repo_delivery: install exact fetched snapshot",
		ref,
		oid,
		zeroOID,
	); err != nil {
		observedOID, symbolic, inspectErr := g.exactRefState(ctx, ref)
		if inspectErr != nil {
			return false, errors.Join(
				fmt.Errorf("install exact fetched snapshot ref: %w", err),
				fmt.Errorf("reconcile snapshot ref installation: %w", inspectErr),
			)
		}
		if observedOID == oid && !symbolic {
			return true, fmt.Errorf(
				"snapshot ref installation reported failure after installing the exact object: %w",
				err,
			)
		}
		if observedOID == "" {
			return false, fmt.Errorf("install exact fetched snapshot ref: %w", err)
		}
		return false, fmt.Errorf(
			"snapshot ref installation failed and the destination is no longer absent; refusing ownership",
		)
	}
	return true, nil
}

func (g *gitRepository) exactRefState(
	ctx context.Context,
	ref string,
) (string, bool, error) {
	result, err := g.run(
		ctx,
		"for-each-ref",
		"--format=%(refname)%00%(objectname)%00%(symref)",
		ref,
	)
	if err != nil {
		return "", false, err
	}
	value := strings.TrimSuffix(result.Stdout, "\n")
	if value == "" {
		return "", false, nil
	}
	fields := strings.Split(value, "\x00")
	if len(fields) != 3 || fields[0] != ref || !isObjectID(fields[1]) ||
		strings.ContainsAny(fields[2], "\x00\r\n") {
		return "", false, fmt.Errorf("Git returned malformed exact ref state")
	}
	return fields[1], fields[2] != "", nil
}

func (s *refSnapshot) close(ctx context.Context) error {
	var cleanupErrors []error
	for _, ref := range s.privateRefs {
		expectedOID := s.privateRefOIDs[ref]
		if !isObjectID(expectedOID) {
			cleanupErrors = append(
				cleanupErrors,
				fmt.Errorf(
					"temporary ref %q lacks its snapshot-owned object ID",
					ref,
				),
			)
			continue
		}
		result, err := s.repository.run(
			ctx,
			"for-each-ref",
			"--format=%(refname)%00%(objectname)%00%(symref)",
			ref,
		)
		if err != nil {
			cleanupErrors = append(
				cleanupErrors,
				fmt.Errorf("inspect temporary ref %q for cleanup: %w", ref, err),
			)
			continue
		}
		value := strings.TrimSuffix(result.Stdout, "\n")
		if value == "" {
			continue
		}
		fields := strings.Split(value, "\x00")
		if len(fields) != 3 || fields[0] != ref || !isObjectID(fields[1]) ||
			strings.ContainsAny(fields[2], "\x00\r\n") {
			cleanupErrors = append(
				cleanupErrors,
				fmt.Errorf("temporary ref %q has malformed cleanup state", ref),
			)
			continue
		}
		if fields[2] != "" {
			cleanupErrors = append(
				cleanupErrors,
				fmt.Errorf("temporary ref %q became symbolic; refusing cleanup", ref),
			)
			continue
		}
		if fields[1] != expectedOID {
			cleanupErrors = append(
				cleanupErrors,
				fmt.Errorf(
					"temporary ref %q no longer identifies its snapshot-owned object; refusing cleanup",
					ref,
				),
			)
			continue
		}
		if _, err := s.repository.run(
			ctx,
			"update-ref",
			"--no-deref",
			"-d",
			ref,
			expectedOID,
		); err != nil {
			cleanupErrors = append(
				cleanupErrors,
				fmt.Errorf("delete exact temporary ref %q: %w", ref, err),
			)
		}
	}
	return errors.Join(cleanupErrors...)
}

func (g *gitRepository) snapshot(
	ctx context.Context,
	remote string,
	baseBranch string,
	headBranch string,
) (resultSnapshot *refSnapshot, returnErr error) {
	if err := g.requireCompleteHistory(ctx); err != nil {
		return nil, err
	}
	if err := g.checkBranch(ctx, baseBranch); err != nil {
		return nil, fmt.Errorf("invalid base: %w", err)
	}
	if err := g.checkBranch(ctx, headBranch); err != nil {
		return nil, fmt.Errorf("invalid head: %w", err)
	}
	boundRemote, err := bindGitRemote(remote)
	if err != nil {
		return nil, err
	}
	transport, err := g.openIsolatedTransport(ctx, boundRemote)
	if err != nil {
		return nil, err
	}
	defer func() {
		if transport == nil {
			return
		}
		if cleanupErr := transport.close(); cleanupErr != nil {
			returnErr = errors.Join(
				returnErr,
				fmt.Errorf("clean isolated snapshot transport: %w", cleanupErr),
			)
		}
	}()
	baseRemoteRef := "refs/heads/" + baseBranch
	headRemoteRef := "refs/heads/" + headBranch
	remoteOIDs, err := transport.remoteObjectIDs(
		ctx,
		baseRemoteRef,
		headRemoteRef,
	)
	if err != nil {
		return nil, err
	}
	remoteBaseOID := remoteOIDs[baseRemoteRef]
	if !isObjectID(remoteBaseOID) {
		return nil, fmt.Errorf("remote base ref is absent or invalid")
	}
	remoteHeadOID := remoteOIDs[headRemoteRef]

	suffix, err := randomHex(12)
	if err != nil {
		return nil, fmt.Errorf("create private ref suffix: %w", err)
	}
	prefix := "refs/repo-delivery/" + suffix
	basePrivateRef := prefix + "/base"
	headPrivateRef := prefix + "/head"
	snapshot := &refSnapshot{
		BaseRef:        basePrivateRef,
		RemoteBaseRef:  baseRemoteRef,
		RemoteHeadRef:  headRemoteRef,
		privateRefOIDs: make(map[string]string),
		repository:     g,
	}
	failed := true
	defer func() {
		if failed {
			if cleanupErr := snapshot.close(context.Background()); cleanupErr != nil {
				returnErr = errors.Join(
					returnErr,
					fmt.Errorf("clean failed snapshot refs: %w", cleanupErr),
				)
			}
		}
	}()
	args := []string{
		"fetch",
		"--atomic",
		"--no-tags",
		"--no-prune",
		"--no-prune-tags",
		"--no-write-fetch-head",
		"--recurse-submodules=no",
		"--",
		boundRemote.endpoint,
		"+" + remoteBaseOID + ":" + basePrivateRef,
	}
	if remoteHeadOID != "" {
		args = append(args, "+"+remoteHeadOID+":"+headPrivateRef)
	}
	if _, err := transport.run(ctx, args...); err != nil {
		return nil, fmt.Errorf("fetch exact base and head refs: %w", err)
	}
	snapshot.BaseOID, err = transport.text(
		ctx,
		"rev-parse",
		"--verify",
		basePrivateRef+"^{commit}",
	)
	if err != nil {
		return nil, fmt.Errorf("resolve fetched base: %w", err)
	}
	if !isObjectID(snapshot.BaseOID) {
		return nil, fmt.Errorf("Git returned invalid fetched base object ID")
	}
	if snapshot.BaseOID != remoteBaseOID {
		return nil, fmt.Errorf("fetched base differs from the exact advertised object ID")
	}
	if remoteHeadOID != "" {
		snapshot.RemoteHeadOID, err = transport.text(
			ctx,
			"rev-parse",
			"--verify",
			headPrivateRef+"^{commit}",
		)
		if err != nil {
			return nil, fmt.Errorf("resolve fetched head: %w", err)
		}
		if !isObjectID(snapshot.RemoteHeadOID) {
			return nil, fmt.Errorf("Git returned invalid fetched head object ID")
		}
		if snapshot.RemoteHeadOID != remoteHeadOID {
			return nil, fmt.Errorf("fetched head differs from the exact advertised object ID")
		}
	}
	for _, entry := range []struct {
		ref string
		oid string
	}{
		{ref: basePrivateRef, oid: snapshot.BaseOID},
		{ref: headPrivateRef, oid: snapshot.RemoteHeadOID},
	} {
		if entry.oid == "" {
			continue
		}
		installed, installErr := g.installSnapshotRef(
			ctx,
			entry.ref,
			entry.oid,
		)
		if installed {
			snapshot.privateRefs = append(snapshot.privateRefs, entry.ref)
			snapshot.privateRefOIDs[entry.ref] = entry.oid
		}
		if installErr != nil {
			return nil, installErr
		}
		verifiedOID, err := g.text(
			ctx,
			"rev-parse",
			"--verify",
			entry.ref+"^{commit}",
		)
		if err != nil || verifiedOID != entry.oid {
			if err != nil {
				return nil, fmt.Errorf("verify installed snapshot ref: %w", err)
			}
			return nil, fmt.Errorf("installed snapshot ref changed unexpectedly")
		}
	}
	if cleanupErr := transport.close(); cleanupErr != nil {
		return nil, fmt.Errorf("clean isolated snapshot transport: %w", cleanupErr)
	}
	transport = nil
	failed = false
	return snapshot, nil
}

func randomHex(byteCount int) (string, error) {
	value := make([]byte, byteCount)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return hex.EncodeToString(value), nil
}

func isObjectID(value string) bool {
	if len(value) != 40 && len(value) != 64 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}

type repositoryStatus struct {
	Staged    []string `json:"staged"`
	Unstaged  []string `json:"unstaged"`
	Untracked []string `json:"untracked"`
}

func (s repositoryStatus) clean() bool {
	return len(s.Staged) == 0 && len(s.Unstaged) == 0 &&
		len(s.Untracked) == 0
}

func (g *gitRepository) status(ctx context.Context) (repositoryStatus, error) {
	result, err := g.run(
		ctx,
		"status",
		"--porcelain=v2",
		"-z",
		"--untracked-files=all",
		"--ignore-submodules=none",
		"--no-renames",
	)
	if err != nil {
		return repositoryStatus{}, fmt.Errorf("inspect repository status: %w", err)
	}
	return parsePorcelainV2Z(result.Stdout)
}

func parsePorcelainV2Z(value string) (repositoryStatus, error) {
	if value == "" {
		return repositoryStatus{}, nil
	}
	if !strings.HasSuffix(value, "\x00") {
		return repositoryStatus{}, fmt.Errorf(
			"Git porcelain-v2 status lacked its terminal NUL",
		)
	}
	records := strings.Split(strings.TrimSuffix(value, "\x00"), "\x00")
	status := repositoryStatus{}
	for _, record := range records {
		if record == "" {
			return repositoryStatus{}, fmt.Errorf(
				"Git porcelain-v2 status contained an empty record",
			)
		}
		switch {
		case strings.HasPrefix(record, "1 "):
			fields := strings.SplitN(record, " ", 9)
			if len(fields) != 9 || fields[8] == "" ||
				!validPorcelainXY(fields[1]) ||
				!validPorcelainSubmodule(fields[2]) ||
				!validPorcelainModes(fields[3:6]) ||
				!validPorcelainOIDs(fields[6:8]) {
				return repositoryStatus{}, fmt.Errorf(
					"Git returned malformed ordinary porcelain-v2 status",
				)
			}
			path := fields[8]
			if fields[1][0] != '.' {
				status.Staged = append(status.Staged, path)
			}
			if fields[1][1] != '.' || porcelainSubmoduleDirty(fields[2]) {
				status.Unstaged = append(status.Unstaged, path)
			}
		case strings.HasPrefix(record, "u "):
			fields := strings.SplitN(record, " ", 11)
			if len(fields) != 11 || fields[10] == "" ||
				!validPorcelainXY(fields[1]) ||
				!validPorcelainSubmodule(fields[2]) ||
				!validPorcelainModes(fields[3:7]) ||
				!validPorcelainOIDs(fields[7:10]) {
				return repositoryStatus{}, fmt.Errorf(
					"Git returned malformed unmerged porcelain-v2 status",
				)
			}
			status.Staged = append(status.Staged, fields[10])
			status.Unstaged = append(status.Unstaged, fields[10])
		case strings.HasPrefix(record, "? "):
			if len(record) == 2 {
				return repositoryStatus{}, fmt.Errorf(
					"Git returned an empty untracked porcelain-v2 path",
				)
			}
			status.Untracked = append(status.Untracked, record[2:])
		default:
			return repositoryStatus{}, fmt.Errorf(
				"Git returned unsupported porcelain-v2 status record %q",
				record[:min(len(record), 2)],
			)
		}
	}
	status.Staged = sortedDistinct(status.Staged)
	status.Unstaged = sortedDistinct(status.Unstaged)
	status.Untracked = sortedDistinct(status.Untracked)
	return status, nil
}

func validPorcelainXY(value string) bool {
	if len(value) != 2 {
		return false
	}
	const statuses = ".MADRCUT"
	return strings.ContainsRune(statuses, rune(value[0])) &&
		strings.ContainsRune(statuses, rune(value[1]))
}

func validPorcelainSubmodule(value string) bool {
	if value == "N..." {
		return true
	}
	return len(value) == 4 && value[0] == 'S' &&
		(value[1] == '.' || value[1] == 'C') &&
		(value[2] == '.' || value[2] == 'M') &&
		(value[3] == '.' || value[3] == 'U')
}

func validPorcelainModes(values []string) bool {
	for _, value := range values {
		if len(value) != 6 {
			return false
		}
		for _, character := range value {
			if character < '0' || character > '7' {
				return false
			}
		}
	}
	return true
}

func validPorcelainOIDs(values []string) bool {
	for _, value := range values {
		if !isObjectID(value) {
			return false
		}
	}
	return true
}

func porcelainSubmoduleDirty(value string) bool {
	return value != "N..." && value != "S..."
}

func sortedDistinct(values []string) []string {
	sort.Strings(values)
	if len(values) < 2 {
		return values
	}
	result := values[:1]
	for _, value := range values[1:] {
		if value != result[len(result)-1] {
			result = append(result, value)
		}
	}
	return result
}

func (g *gitRepository) nulPaths(
	ctx context.Context,
	args ...string,
) ([]string, error) {
	result, err := g.run(ctx, args...)
	if err != nil {
		return nil, err
	}
	trimmed := strings.TrimSuffix(result.Stdout, "\x00")
	if trimmed == "" {
		return nil, nil
	}
	return strings.Split(trimmed, "\x00"), nil
}

func (g *gitRepository) addPaths(
	ctx context.Context,
	paths []string,
) error {
	args := []string{"--literal-pathspecs", "add", "--"}
	args = append(args, paths...)
	if _, err := g.run(ctx, args...); err != nil {
		return fmt.Errorf("stage explicit task paths: %w", err)
	}
	return nil
}

func (g *gitRepository) diffCheck(ctx context.Context) error {
	if _, err := g.run(ctx, "diff", "--check"); err != nil {
		return fmt.Errorf("git diff --check failed: %w", err)
	}
	if _, err := g.run(ctx, "diff", "--cached", "--check"); err != nil {
		return fmt.Errorf("git diff --cached --check failed: %w", err)
	}
	return nil
}

func (g *gitRepository) changedPaths(
	ctx context.Context,
	oldObject string,
	newObject string,
) ([]string, error) {
	if !isObjectID(oldObject) || !isObjectID(newObject) {
		return nil, fmt.Errorf("tree diff endpoints must be full Git object IDs")
	}
	result, err := g.run(
		ctx,
		"diff",
		"--no-ext-diff",
		"--no-textconv",
		"--ignore-submodules=none",
		"--no-renames",
		"--name-only",
		"-z",
		oldObject,
		newObject,
		"--",
	)
	if err != nil {
		return nil, fmt.Errorf("inspect immutable aggregate paths: %w", err)
	}
	paths, err := strictNULPaths(result.Stdout)
	if err != nil {
		return nil, fmt.Errorf("parse immutable aggregate paths: %w", err)
	}
	sort.Strings(paths)
	return paths, nil
}

func (g *gitRepository) pathEntry(
	ctx context.Context,
	object string,
	path string,
) (string, error) {
	if !isObjectID(object) {
		return "", fmt.Errorf("tree entry object is not a full Git object ID")
	}
	result, err := g.run(
		ctx,
		"--literal-pathspecs",
		"ls-tree",
		"-z",
		object,
		"--",
		path,
	)
	if err != nil {
		return "", fmt.Errorf("inspect tree entry %q: %w", path, err)
	}
	return result.Stdout, nil
}

func strictNULPaths(value string) ([]string, error) {
	if value == "" {
		return nil, nil
	}
	if !strings.HasSuffix(value, "\x00") {
		return nil, fmt.Errorf("path list lacked its terminal NUL")
	}
	paths := strings.Split(strings.TrimSuffix(value, "\x00"), "\x00")
	for _, path := range paths {
		if path == "" {
			return nil, fmt.Errorf("path list contained an empty entry")
		}
	}
	return paths, nil
}

func (g *gitRepository) rangeDiffCheck(
	ctx context.Context,
	base string,
	head string,
) error {
	if err := requireObjectIDRange(base, head); err != nil {
		return err
	}
	result, err := g.run(ctx, "diff", "--check", base, head)
	if err != nil {
		detail := strings.TrimSpace(redactCredentials(result.Stdout))
		if detail != "" {
			return fmt.Errorf(
				"aggregate git diff --check failed: %s: %w",
				detail,
				err,
			)
		}
		return fmt.Errorf("aggregate git diff --check failed: %w", err)
	}
	return nil
}

func (g *gitRepository) hasStagedChanges(
	ctx context.Context,
) (bool, error) {
	_, err := g.run(ctx, "diff", "--cached", "--quiet", "--exit-code")
	if err == nil {
		return false, nil
	}
	if code, ok := commandExitCode(err); ok && code == 1 {
		return true, nil
	}
	return false, fmt.Errorf("inspect staged diff: %w", err)
}

func (g *gitRepository) hasRangeChanges(
	ctx context.Context,
	base string,
	head string,
) (bool, error) {
	if err := requireObjectIDRange(base, head); err != nil {
		return false, err
	}
	_, err := g.run(ctx, "diff", "--quiet", "--exit-code", base, head)
	if err == nil {
		return false, nil
	}
	if code, ok := commandExitCode(err); ok && code == 1 {
		return true, nil
	}
	return false, fmt.Errorf("inspect aggregate feature diff: %w", err)
}

func (g *gitRepository) pathTracked(
	ctx context.Context,
	path string,
) (bool, error) {
	_, err := g.run(ctx, "ls-files", "--error-unmatch", "--", path)
	if err == nil {
		return true, nil
	}
	if code, ok := commandExitCode(err); ok && code == 1 {
		return false, nil
	}
	return false, fmt.Errorf("inspect whether %q is tracked: %w", path, err)
}

func (g *gitRepository) pathIgnored(
	ctx context.Context,
	path string,
) (bool, error) {
	_, err := g.run(ctx, "check-ignore", "--quiet", "--", path)
	if err == nil {
		return true, nil
	}
	if code, ok := commandExitCode(err); ok && code == 1 {
		return false, nil
	}
	return false, fmt.Errorf("inspect whether %q is ignored: %w", path, err)
}

func (g *gitRepository) operationState(
	ctx context.Context,
) ([]string, error) {
	active := make([]string, 0)
	for _, marker := range gitOperationMarkers {
		value, err := g.text(ctx, "rev-parse", "--git-path", marker)
		if err != nil {
			return nil, fmt.Errorf("resolve Git operation marker %s: %w", marker, err)
		}
		if !filepath.IsAbs(value) {
			value = filepath.Join(g.directory, value)
		}
		if _, err := os.Lstat(value); err == nil {
			active = append(active, marker)
		} else if !os.IsNotExist(err) {
			return nil, fmt.Errorf("inspect Git operation marker %s: %w", marker, err)
		}
	}
	return active, nil
}

func (g *gitRepository) requireNoOperation(ctx context.Context) error {
	active, err := g.operationState(ctx)
	if err != nil {
		return err
	}
	if len(active) != 0 {
		return fmt.Errorf(
			"an in-progress Git operation is present: %s",
			strings.Join(active, ", "),
		)
	}
	return nil
}

func (g *gitRepository) indexTree(
	ctx context.Context,
) (string, error) {
	tree, err := g.text(ctx, "write-tree")
	if err != nil {
		return "", fmt.Errorf("write index tree: %w", err)
	}
	if !isObjectID(tree) {
		return "", fmt.Errorf("Git returned invalid index tree object ID %q", tree)
	}
	return tree, nil
}

func (g *gitRepository) requireIndexTree(
	ctx context.Context,
	expectedTree string,
) error {
	if !isObjectID(expectedTree) {
		return fmt.Errorf("expected index tree is not a full Git object ID")
	}
	actualTree, err := g.indexTree(ctx)
	if err != nil {
		return err
	}
	if actualTree != expectedTree {
		return fmt.Errorf(
			"index tree changed: got %s, expected %s",
			actualTree,
			expectedTree,
		)
	}
	return nil
}

var indexDebugRecordPattern = regexp.MustCompile(
	`^  ctime: [0-9]+:[0-9]+\n` +
		`  mtime: [0-9]+:[0-9]+\n` +
		`  dev: [0-9]+\tino: [0-9]+\n` +
		`  uid: [0-9]+\tgid: [0-9]+\n` +
		`  size: [0-9]+\tflags: ([0-9A-Fa-f]+)\n`,
)

// requireDefaultIndexFlags rejects per-entry index state that can make Git's
// ordinary status and diff views omit worktree changes. The --debug format is
// intentionally parsed strictly: if a future Git changes it, delivery fails
// closed instead of silently accepting an index flag it no longer understands.
func (g *gitRepository) requireDefaultIndexFlags(
	ctx context.Context,
	paths []string,
) error {
	return g.requireDefaultIndexFlagsEnvironment(ctx, paths, nil)
}

func (g *gitRepository) requireDefaultIndexFlagsEnvironment(
	ctx context.Context,
	paths []string,
	environment []string,
) error {
	if len(paths) == 0 {
		return fmt.Errorf("index flag inspection requires at least one path")
	}
	args := []string{
		"--literal-pathspecs",
		"ls-files",
		"--cached",
		"--stage",
		"--debug",
		"-z",
		"--full-name",
		"--no-recurse-submodules",
		"--",
	}
	args = append(args, paths...)
	first, err := g.runEnvironmentWithOutputLimit(
		ctx,
		environment,
		indexFlagInspectionOutputLimit,
		args...,
	)
	if err != nil {
		return fmt.Errorf("inspect exact index flags: %w", err)
	}
	second, err := g.runEnvironmentWithOutputLimit(
		ctx,
		environment,
		indexFlagInspectionOutputLimit,
		args...,
	)
	if err != nil {
		return fmt.Errorf("confirm exact index flags: %w", err)
	}
	if first.Stdout != second.Stdout {
		return fmt.Errorf("index entries changed during flag inspection")
	}
	if err := parseDefaultIndexFlags(first.Stdout); err != nil {
		return fmt.Errorf("refuse non-default index state: %w", err)
	}
	return nil
}

func parseDefaultIndexFlags(value string) error {
	remaining := value
	seen := make(map[string]bool)
	previousPath := ""
	for remaining != "" {
		header, rest, found := strings.Cut(remaining, "\x00")
		if !found {
			return fmt.Errorf("Git index metadata lacked a NUL-delimited header")
		}
		metadata, path, found := strings.Cut(header, "\t")
		if !found {
			return fmt.Errorf("Git index metadata lacked its path separator")
		}
		fields := strings.Fields(metadata)
		if len(fields) != 3 {
			return fmt.Errorf("Git returned malformed stage metadata")
		}
		mode, oid, stage := fields[0], fields[1], fields[2]
		if len(mode) != 6 || !validPorcelainModes([]string{mode}) ||
			!isObjectID(oid) || strings.Trim(oid, "0") == "" ||
			stage != "0" || path == "" || seen[path] ||
			(previousPath != "" && path <= previousPath) {
			return fmt.Errorf("Git returned malformed or non-stage-zero index metadata")
		}
		seen[path] = true
		previousPath = path
		remaining = rest
		match := indexDebugRecordPattern.FindStringSubmatchIndex(remaining)
		if len(match) != 4 || match[0] != 0 {
			return fmt.Errorf("Git returned an unsupported --debug index record")
		}
		flags, err := strconv.ParseUint(
			remaining[match[2]:match[3]],
			16,
			64,
		)
		if err != nil {
			return fmt.Errorf("parse index flags for %q: %w", path, err)
		}
		if flags != 0 {
			return fmt.Errorf(
				"path %q has non-default index flags 0x%x",
				path,
				flags,
			)
		}
		remaining = remaining[match[1]:]
	}
	return nil
}

func (g *gitRepository) commitParents(
	ctx context.Context,
	oid string,
) ([]string, error) {
	value, err := g.text(ctx, "show", "--no-patch", "--format=%P", oid)
	if err != nil {
		return nil, fmt.Errorf("read commit parents for %s: %w", oid, err)
	}
	if value == "" {
		return nil, nil
	}
	parents := strings.Fields(value)
	for _, parent := range parents {
		if !isObjectID(parent) {
			return nil, fmt.Errorf("Git returned invalid parent for commit %s", oid)
		}
	}
	return parents, nil
}

func (g *gitRepository) fullCommitMessage(
	ctx context.Context,
	oid string,
) (string, error) {
	result, err := g.run(ctx, "show", "--no-patch", "--format=%B", oid)
	if err != nil {
		return "", fmt.Errorf("read full commit message %s: %w", oid, err)
	}
	return normalizeText(result.Stdout), nil
}

func equalStrings(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

type commitAuthor struct {
	Name  string
	Email string
	Date  string
}

func (g *gitRepository) author(
	ctx context.Context,
	oid string,
) (commitAuthor, error) {
	result, err := g.run(
		ctx,
		"show",
		"--no-patch",
		"--format=%an%x00%ae%x00%aI",
		oid,
	)
	if err != nil {
		return commitAuthor{}, fmt.Errorf("read commit author %s: %w", oid, err)
	}
	fields := strings.Split(normalizeText(result.Stdout), "\x00")
	if len(fields) != 3 || fields[0] == "" || fields[1] == "" || fields[2] == "" {
		return commitAuthor{}, fmt.Errorf("Git returned malformed author for commit %s", oid)
	}
	return commitAuthor{Name: fields[0], Email: fields[1], Date: fields[2]}, nil
}

func (a commitAuthor) environment() []string {
	return []string{
		"GIT_AUTHOR_NAME=" + a.Name,
		"GIT_AUTHOR_EMAIL=" + a.Email,
		"GIT_AUTHOR_DATE=" + a.Date,
	}
}

func (g *gitRepository) commitChecked(
	ctx context.Context,
	message string,
	amend bool,
	expectedHead string,
	expectedTree string,
	branch string,
	flagPathSets ...[]string,
) (string, error) {
	parents := []string{expectedHead}
	authorOID := ""
	signatureSources := []string(nil)
	if amend {
		var err error
		parents, err = g.commitParents(ctx, expectedHead)
		if err != nil {
			return "", err
		}
		authorOID = expectedHead
		signatureSources = []string{expectedHead}
	}
	return g.commitPlannedChecked(
		ctx,
		message,
		expectedHead,
		expectedTree,
		branch,
		parents,
		authorOID,
		signatureSources,
		flagPathSets...,
	)
}

func (g *gitRepository) commitConsolidatedChecked(
	ctx context.Context,
	message string,
	expectedHead string,
	expectedTree string,
	branch string,
	parentOID string,
	authorOID string,
	signatureSources []string,
	flagPathSets ...[]string,
) (string, error) {
	if !isObjectID(parentOID) || !isObjectID(authorOID) {
		return "", fmt.Errorf("consolidated commit plan contains an invalid object ID")
	}
	if len(signatureSources) == 0 {
		return "", fmt.Errorf("consolidated commit plan lacks signature sources")
	}
	for _, oid := range signatureSources {
		if !isObjectID(oid) {
			return "", fmt.Errorf("consolidated signature source is not a full object ID")
		}
	}
	return g.commitPlannedChecked(
		ctx,
		message,
		expectedHead,
		expectedTree,
		branch,
		[]string{parentOID},
		authorOID,
		signatureSources,
		flagPathSets...,
	)
}

func (g *gitRepository) commitPlannedChecked(
	ctx context.Context,
	message string,
	expectedHead string,
	expectedTree string,
	branch string,
	expectedParents []string,
	authorOID string,
	signatureSources []string,
	flagPathSets ...[]string,
) (string, error) {
	if len(flagPathSets) > 1 {
		return "", fmt.Errorf("commit accepts at most one index-flag path set")
	}
	if len(expectedParents) > 1 {
		return "", fmt.Errorf("aggregate commit accepts at most one parent")
	}
	for _, parent := range expectedParents {
		if !isObjectID(parent) {
			return "", fmt.Errorf("aggregate commit parent is not a full object ID")
		}
	}
	if err := g.requireCompleteHistory(ctx); err != nil {
		return "", err
	}
	if err := g.requireNoOperation(ctx); err != nil {
		return "", err
	}
	currentBranch, currentHead, err := g.branchHead(ctx)
	if err != nil {
		return "", err
	}
	if currentHead != expectedHead {
		return "", fmt.Errorf(
			"HEAD changed before commit: got %s, expected %s",
			currentHead,
			expectedHead,
		)
	}
	if currentBranch != branch {
		return "", fmt.Errorf("current branch changed before commit")
	}
	if len(flagPathSets) == 1 {
		if err := g.requireDefaultIndexFlags(ctx, flagPathSets[0]); err != nil {
			return "", fmt.Errorf("commit index flag gate: %w", err)
		}
	}
	indexTree, err := g.indexTree(ctx)
	if err != nil {
		return "", err
	}
	if indexTree != expectedTree {
		return "", fmt.Errorf(
			"index changed before commit: got tree %s, expected %s",
			indexTree,
			expectedTree,
		)
	}
	var expectedAuthor *commitAuthor
	if authorOID != "" {
		if !isObjectID(authorOID) {
			return "", fmt.Errorf("aggregate author source is not a full object ID")
		}
		author, err := g.author(ctx, authorOID)
		if err != nil {
			return "", err
		}
		expectedAuthor = &author
	}
	sign, err := g.signingConfigured(ctx)
	if err != nil {
		return "", err
	}
	for _, source := range signatureSources {
		if sign {
			break
		}
		sign, err = g.commitHasSignature(ctx, source)
		if err != nil {
			return "", err
		}
	}

	args := []string{"commit-tree", expectedTree}
	for _, parent := range expectedParents {
		args = append(args, "-p", parent)
	}
	if sign {
		args = append(args, "--gpg-sign")
	}
	args = append(args, "-F", "-")
	environment := []string(nil)
	if expectedAuthor != nil {
		environment = expectedAuthor.environment()
	}
	result, err := g.runInputEnv(ctx, message, environment, args...)
	if err != nil {
		return "", fmt.Errorf("create exact aggregate commit object: %w", err)
	}
	newHead := strings.TrimSpace(result.Stdout)
	if !isObjectID(newHead) {
		return "", fmt.Errorf(
			"Git returned invalid aggregate commit object ID %q",
			newHead,
		)
	}
	newTree, treeErr := g.text(ctx, "rev-parse", newHead+"^{tree}")
	newParents, parentsErr := g.commitParents(ctx, newHead)
	newMessage, messageErr := g.fullCommitMessage(ctx, newHead)
	authorMatches := true
	if expectedAuthor != nil {
		newAuthor, authorErr := g.author(ctx, newHead)
		authorMatches = authorErr == nil && newAuthor == *expectedAuthor
	}
	valid := treeErr == nil && parentsErr == nil && messageErr == nil &&
		newTree == expectedTree && equalStrings(newParents, expectedParents) &&
		normalizeText(newMessage) == normalizeText(message) && authorMatches
	if !valid {
		return "", fmt.Errorf(
			"created commit object violated the expected tree, parents, author, or message",
		)
	}
	if err := g.requireNoOperation(ctx); err != nil {
		return "", err
	}
	currentHead, err = g.head(ctx)
	if err != nil {
		return "", err
	}
	if currentHead != expectedHead {
		return "", fmt.Errorf("HEAD changed before branch update")
	}
	currentBranch, err = g.currentBranch(ctx)
	if err != nil {
		return "", err
	}
	if currentBranch != branch {
		return "", fmt.Errorf("current branch changed before branch update")
	}
	if len(flagPathSets) == 1 {
		if err := g.requireDefaultIndexFlags(ctx, flagPathSets[0]); err != nil {
			return "", fmt.Errorf("final commit index flag gate: %w", err)
		}
	}
	currentIndexTree, err := g.indexTree(ctx)
	if err != nil {
		return "", err
	}
	if currentIndexTree != expectedTree {
		return "", fmt.Errorf("index changed before branch update")
	}
	if err := g.advanceAttachedBranch(
		ctx,
		branch,
		expectedHead,
		newHead,
		expectedTree,
	); err != nil {
		return "", fmt.Errorf(
			"advance attached branch with exact compare-and-swap: %w",
			err,
		)
	}
	return newHead, nil
}

func (g *gitRepository) commitCount(
	ctx context.Context,
	base string,
	head string,
) (int, error) {
	if err := requireObjectIDRange(base, head); err != nil {
		return 0, err
	}
	value, err := g.text(ctx, "rev-list", "--count", base+".."+head)
	if err != nil {
		return 0, fmt.Errorf("count feature commits: %w", err)
	}
	count, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("parse feature commit count %q: %w", value, err)
	}
	return count, nil
}

func (g *gitRepository) mergeCommits(
	ctx context.Context,
	base string,
	head string,
) ([]string, error) {
	if err := requireObjectIDRange(base, head); err != nil {
		return nil, err
	}
	value, err := g.text(ctx, "rev-list", "--merges", base+".."+head)
	if err != nil {
		return nil, fmt.Errorf("inspect feature merge commits: %w", err)
	}
	if value == "" {
		return nil, nil
	}
	return strings.Split(value, "\n"), nil
}

func (g *gitRepository) isAncestor(
	ctx context.Context,
	ancestor string,
	descendant string,
) (bool, error) {
	if err := requireObjectIDRange(ancestor, descendant); err != nil {
		return false, err
	}
	_, err := g.run(ctx, "merge-base", "--is-ancestor", ancestor, descendant)
	if err == nil {
		return true, nil
	}
	if code, ok := commandExitCode(err); ok && code == 1 {
		return false, nil
	}
	return false, fmt.Errorf("inspect ancestry: %w", err)
}

type commitProjection struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (g *gitRepository) projection(
	ctx context.Context,
	oid string,
) (commitProjection, error) {
	title, err := g.text(ctx, "show", "--no-patch", "--format=%s", oid)
	if err != nil {
		return commitProjection{}, fmt.Errorf("read commit subject %s: %w", oid, err)
	}
	result, err := g.run(ctx, "show", "--no-patch", "--format=%b", oid)
	if err != nil {
		return commitProjection{}, fmt.Errorf("read commit body %s: %w", oid, err)
	}
	body, err := pullRequestBody(result.Stdout)
	if err != nil {
		return commitProjection{}, fmt.Errorf("project commit %s: %w", oid, err)
	}
	return commitProjection{Title: title, Body: body}, nil
}

func (g *gitRepository) commitBody(
	ctx context.Context,
	oid string,
) (string, error) {
	result, err := g.run(ctx, "show", "--no-patch", "--format=%b", oid)
	if err != nil {
		return "", fmt.Errorf("read commit body %s: %w", oid, err)
	}
	return normalizeText(result.Stdout), nil
}

type featureCommit struct {
	OID             string   `json:"oid"`
	Parents         []string `json:"parents"`
	AuthorName      string   `json:"author_name"`
	AuthorEmail     string   `json:"author_email"`
	CommitterName   string   `json:"committer_name"`
	CommitterEmail  string   `json:"committer_email"`
	SignatureStatus string   `json:"signature_status"`
	Title           string   `json:"title"`
	HasDisclaimer   bool     `json:"has_commit_disclaimer"`
}

func (g *gitRepository) featureCommits(
	ctx context.Context,
	base string,
	head string,
) ([]featureCommit, error) {
	if err := requireObjectIDRange(base, head); err != nil {
		return nil, err
	}
	value, err := g.text(ctx, "rev-list", "--reverse", base+".."+head)
	if err != nil {
		return nil, fmt.Errorf("list feature commits: %w", err)
	}
	if value == "" {
		return nil, nil
	}
	oids := strings.Split(value, "\n")
	commits := make([]featureCommit, 0, len(oids))
	for _, oid := range oids {
		result, err := g.run(
			ctx,
			"show",
			"--no-patch",
			"--format=%H%x00%P%x00%an%x00%ae%x00%cn%x00%ce%x00%G?%x00%s",
			oid,
		)
		if err != nil {
			return nil, fmt.Errorf("inspect feature commit %s: %w", oid, err)
		}
		fields := strings.Split(normalizeText(result.Stdout), "\x00")
		if len(fields) != 8 || !isObjectID(fields[0]) {
			return nil, fmt.Errorf("Git returned malformed metadata for feature commit %s", oid)
		}
		parents := strings.Fields(fields[1])
		body, err := g.commitBody(ctx, oid)
		if err != nil {
			return nil, err
		}
		commits = append(commits, featureCommit{
			OID:             fields[0],
			Parents:         parents,
			AuthorName:      fields[2],
			AuthorEmail:     fields[3],
			CommitterName:   fields[4],
			CommitterEmail:  fields[5],
			SignatureStatus: fields[6],
			Title:           fields[7],
			HasDisclaimer:   hasFinalLine(body, commitDisclaimer),
		})
	}
	return commits, nil
}

func (g *gitRepository) signingConfigured(
	ctx context.Context,
) (bool, error) {
	value, err := g.text(ctx, "config", "--get", "--bool", "commit.gpgsign")
	if err == nil {
		return value == "true", nil
	} else if code, ok := commandExitCode(err); !ok || code != 1 {
		return false, fmt.Errorf("inspect commit signing configuration: %w", err)
	}
	return false, nil
}

func (g *gitRepository) commitHasSignature(
	ctx context.Context,
	oid string,
) (bool, error) {
	result, err := g.run(ctx, "cat-file", "commit", oid)
	if err != nil {
		return false, fmt.Errorf("inspect commit signature: %w", err)
	}
	header, _, _ := strings.Cut(result.Stdout, "\n\n")
	return commitHeaderHasSignature(header), nil
}

func commitHeaderHasSignature(header string) bool {
	for _, name := range []string{"gpgsig", "gpgsig-sha256"} {
		if strings.HasPrefix(header, name+" ") ||
			strings.Contains(header, "\n"+name+" ") {
			return true
		}
	}
	return false
}

func (g *gitRepository) commitRequiresSignature(
	ctx context.Context,
	oid string,
) (bool, error) {
	configured, err := g.signingConfigured(ctx)
	if err != nil || configured {
		return configured, err
	}
	return g.commitHasSignature(ctx, oid)
}

func (g *gitRepository) verifyRequiredSignature(
	ctx context.Context,
	oid string,
) error {
	if !isObjectID(oid) {
		return fmt.Errorf("signature target is not a full Git object ID")
	}
	required, err := g.commitRequiresSignature(ctx, oid)
	if err != nil {
		return err
	}
	return g.verifyCommitSignature(ctx, oid, required)
}

func requireObjectIDRange(base string, head string) error {
	if !isObjectID(base) || !isObjectID(head) {
		return fmt.Errorf("commit range endpoints must be full Git object IDs")
	}
	return nil
}

func (g *gitRepository) verifyCommitSignature(
	ctx context.Context,
	oid string,
	required bool,
) error {
	if !required {
		return nil
	}
	present, err := g.commitHasSignature(ctx, oid)
	if err != nil {
		return err
	}
	if !present {
		return fmt.Errorf("final commit does not preserve the required signature")
	}
	status, err := g.text(ctx, "show", "--no-patch", "--format=%G?", oid)
	if err != nil {
		return fmt.Errorf("inspect final commit signature: %w", err)
	}
	if status == "B" {
		return fmt.Errorf("final commit has a cryptographically bad signature")
	}
	return nil
}

func (g *gitRepository) rebase(
	ctx context.Context,
	baseOID string,
	branch string,
	expectedHead string,
	validators ...func(context.Context, *gitRepository, string) error,
) (rebasedHead string, returnErr error) {
	if len(validators) > 1 {
		return "", fmt.Errorf("rebase accepts at most one candidate validator")
	}
	if !isObjectID(baseOID) {
		return "", fmt.Errorf("rebase base is not a full Git object ID")
	}
	if !isObjectID(expectedHead) {
		return "", fmt.Errorf("rebase head is not a full Git object ID")
	}
	if err := g.requireCompleteHistory(ctx); err != nil {
		return "", err
	}
	if err := g.checkBranch(ctx, branch); err != nil {
		return "", err
	}
	if err := g.requireNoOperation(ctx); err != nil {
		return "", err
	}
	if err := g.requireDefaultIndexFlags(ctx, []string{"."}); err != nil {
		return "", fmt.Errorf("rebase index flag gate: %w", err)
	}
	status, err := g.status(ctx)
	if err != nil {
		return "", err
	}
	if !status.clean() {
		return "", fmt.Errorf("rebase requires a clean worktree and index")
	}
	currentBranch, currentHead, err := g.branchHead(ctx)
	if err != nil {
		return "", err
	}
	if currentBranch != branch {
		return "", fmt.Errorf(
			"current branch changed before rebase: got %q, expected %q",
			currentBranch,
			branch,
		)
	}
	if currentHead != expectedHead {
		return "", fmt.Errorf(
			"HEAD changed before rebase: got %s, expected %s",
			currentHead,
			expectedHead,
		)
	}
	expectedTree, err := g.tree(ctx, expectedHead)
	if err != nil {
		return "", err
	}
	if err := g.requireIndexTree(ctx, expectedTree); err != nil {
		return "", fmt.Errorf("rebase index gate: %w", err)
	}
	oldParents, err := g.commitParents(ctx, expectedHead)
	if err != nil {
		return "", err
	}
	if len(oldParents) != 1 {
		return "", fmt.Errorf(
			"rebase source has %d parents, want exactly one",
			len(oldParents),
		)
	}
	signRequired, err := g.commitRequiresSignature(ctx, expectedHead)
	if err != nil {
		return "", err
	}
	isolated, err := g.createIsolatedRebaseWorktree(ctx, expectedHead)
	if err != nil {
		return "", err
	}
	isolatedDirectory := isolated.directory
	defer func() {
		if isolated == nil {
			return
		}
		cleanupErr := g.removeIsolatedRebaseWorktree(
			context.Background(),
			isolatedDirectory,
		)
		if cleanupErr != nil {
			cleanupErr = fmt.Errorf(
				"clean isolated rebase worktree: %w",
				cleanupErr,
			)
			if returnErr == nil {
				returnErr = cleanupErr
			} else {
				returnErr = errors.Join(returnErr, cleanupErr)
			}
		}
	}()

	newHead, err := isolated.replayExactCommit(
		ctx,
		baseOID,
		oldParents[0],
		expectedHead,
		signRequired,
	)
	if err != nil {
		return "", err
	}
	if len(validators) == 1 {
		if err := validators[0](ctx, isolated, newHead); err != nil {
			return "", fmt.Errorf("validate isolated rebased candidate: %w", err)
		}
	}
	// Remove the temporary checkout before touching the user's branch. If
	// cleanup is not provably complete, the original branch and checkout stay
	// unchanged and the deferred retry handles recoverable cleanup failures.
	if err := g.removeIsolatedRebaseWorktree(ctx, isolatedDirectory); err != nil {
		return "", fmt.Errorf("clean isolated rebase worktree: %w", err)
	}
	isolated = nil
	if err := g.reconcileRebasedBranch(
		ctx,
		branch,
		expectedHead,
		newHead,
	); err != nil {
		return "", err
	}
	return newHead, nil
}

func (g *gitRepository) createIsolatedRebaseWorktree(
	ctx context.Context,
	expectedHead string,
) (*gitRepository, error) {
	suffix, err := randomHex(12)
	if err != nil {
		return nil, fmt.Errorf("create isolated rebase suffix: %w", err)
	}
	relativeRoot := filepath.Join("out", "repo_delivery")
	relativePath := filepath.Join(relativeRoot, "rebase-"+suffix)
	ignored, err := g.pathIgnored(ctx, filepath.ToSlash(relativePath))
	if err != nil {
		return nil, err
	}
	if !ignored {
		return nil, fmt.Errorf(
			"isolated rebase path %q is not ignored by Git",
			relativePath,
		)
	}
	root := filepath.Join(g.directory, relativeRoot)
	if err := ensurePrivateDirectory(filepath.Join(g.directory, "out")); err != nil {
		return nil, err
	}
	if err := ensurePrivateDirectory(root); err != nil {
		return nil, err
	}
	path := filepath.Join(g.directory, relativePath)
	if _, err := os.Lstat(path); err == nil {
		return nil, fmt.Errorf("isolated rebase path already exists")
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("inspect isolated rebase path: %w", err)
	}
	if _, err := g.run(
		ctx,
		"worktree",
		"add",
		"--detach",
		"--no-guess-remote",
		path,
		expectedHead,
	); err != nil {
		return nil, fmt.Errorf("create isolated rebase worktree: %w", err)
	}
	return &gitRepository{
		directory:  path,
		executable: g.executable,
		runner:     g.runner,
	}, nil
}

func ensurePrivateDirectory(path string) error {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		if err := os.Mkdir(path, 0o700); err != nil {
			return fmt.Errorf("create private directory %q: %w", path, err)
		}
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect private directory %q: %w", path, err)
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("private path %q is not a real directory", path)
	}
	return nil
}

func (g *gitRepository) removeIsolatedRebaseWorktree(
	ctx context.Context,
	path string,
) error {
	root := filepath.Join(g.directory, "out", "repo_delivery")
	relative, err := filepath.Rel(root, path)
	if err != nil || relative == "." || relative == ".." ||
		strings.HasPrefix(relative, ".."+string(filepath.Separator)) ||
		strings.ContainsRune(relative, filepath.Separator) {
		return fmt.Errorf("refuse unsafe isolated worktree path %q", path)
	}
	if _, err := g.run(ctx, "worktree", "remove", "--force", path); err != nil {
		return err
	}
	if _, err := os.Lstat(path); err == nil {
		return fmt.Errorf("isolated worktree path still exists after removal")
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("inspect removed isolated worktree: %w", err)
	}
	return nil
}

func (g *gitRepository) replayExactCommit(
	ctx context.Context,
	baseOID string,
	oldParentOID string,
	expectedHead string,
	signRequired bool,
) (string, error) {
	if err := g.requireCompleteHistory(ctx); err != nil {
		return "", err
	}
	status, err := g.status(ctx)
	if err != nil {
		return "", err
	}
	if !status.clean() {
		return "", fmt.Errorf("isolated rebase worktree is not clean")
	}
	currentHead, err := g.head(ctx)
	if err != nil {
		return "", err
	}
	if currentHead != expectedHead {
		return "", fmt.Errorf("isolated rebase HEAD differs from the captured commit")
	}
	expectedTree, err := g.tree(ctx, expectedHead)
	if err != nil {
		return "", err
	}
	if err := g.requireIndexTree(ctx, expectedTree); err != nil {
		return "", fmt.Errorf("isolated rebase index gate: %w", err)
	}
	if err := g.requireDefaultIndexFlags(ctx, []string{"."}); err != nil {
		return "", fmt.Errorf("isolated rebase index flag gate: %w", err)
	}
	args := []string{
		"-c",
		"notes.rewrite.rebase=false",
		"-c",
		"rerere.enabled=false",
		"rebase",
		"--onto",
		baseOID,
		"--no-verify",
		"--no-autostash",
		"--no-update-refs",
		"--no-rebase-merges",
		"--no-fork-point",
		"--no-rerere-autoupdate",
		"--no-autosquash",
		"--no-signoff",
		"--reapply-cherry-picks",
		"--empty=stop",
	}
	if signRequired {
		args = append(args, "--gpg-sign")
	} else {
		args = append(args, "--no-gpg-sign")
	}
	args = append(args, oldParentOID, expectedHead)
	if _, err := g.run(ctx, args...); err != nil {
		active, stateErr := g.operationState(ctx)
		if stateErr == nil && len(active) != 0 {
			_, _ = g.run(context.Background(), "rebase", "--abort")
		}
		return "", fmt.Errorf(
			"isolated exact-OID rebase conflicted or failed; no main-worktree rebase state was changed: %w",
			err,
		)
	}
	if err := g.requireNoOperation(ctx); err != nil {
		return "", err
	}
	newHead, err := g.head(ctx)
	if err != nil {
		return "", err
	}
	if newHead == expectedHead {
		return "", fmt.Errorf("rebase did not create a new exact commit")
	}
	if _, err := g.run(ctx, "symbolic-ref", "--quiet", "HEAD"); err == nil {
		return "", fmt.Errorf("exact-OID rebase unexpectedly left HEAD attached")
	} else if code, ok := commandExitCode(err); !ok || code != 1 {
		return "", fmt.Errorf("inspect detached HEAD after rebase: %w", err)
	}
	newParents, err := g.commitParents(ctx, newHead)
	if err != nil {
		return "", err
	}
	if len(newParents) != 1 || newParents[0] != baseOID {
		return "", fmt.Errorf("rebased commit does not have the exact fetched base parent")
	}
	if err := g.verifyCommitSignature(ctx, newHead, signRequired); err != nil {
		return "", err
	}
	status, err = g.status(ctx)
	if err != nil {
		return "", err
	}
	if !status.clean() {
		return "", fmt.Errorf("rebase changed the isolated worktree outside its commit")
	}
	newTree, err := g.tree(ctx, newHead)
	if err != nil {
		return "", err
	}
	if err := g.requireIndexTree(ctx, newTree); err != nil {
		return "", fmt.Errorf("rebased candidate index gate: %w", err)
	}
	return newHead, nil
}

func (g *gitRepository) requireWorktreeMatchesIndex(
	ctx context.Context,
) error {
	if _, err := g.run(
		ctx,
		"diff-files",
		"--quiet",
		"--ignore-submodules=none",
		"--",
	); err != nil {
		return fmt.Errorf("worktree differs from the exact index: %w", err)
	}
	return nil
}

func (g *gitRepository) advanceAttachedBranch(
	ctx context.Context,
	branch string,
	expectedHead string,
	newHead string,
	expectedIndexTree string,
) error {
	if !isObjectID(newHead) {
		return fmt.Errorf("attached branch target is not a full Git object ID")
	}
	observedTree, err := g.indexTree(ctx)
	if err != nil {
		return fmt.Errorf("attached branch index gate: %w", err)
	}
	if observedTree != expectedIndexTree {
		return fmt.Errorf(
			"attached branch index gate: got %s, want %s",
			observedTree,
			expectedIndexTree,
		)
	}
	if err := g.requireBranchHead(ctx, branch, expectedHead); err != nil {
		return err
	}
	if _, err := g.run(
		ctx,
		"update-ref",
		"-m",
		"repo_delivery: advance attached branch",
		"refs/heads/"+branch,
		newHead,
		expectedHead,
	); err != nil {
		return fmt.Errorf("advance attached branch with exact compare-and-swap: %w", err)
	}
	if err := g.requireBranchHead(ctx, branch, newHead); err != nil {
		return fmt.Errorf("attached branch changed during advancement: %w", err)
	}
	return nil
}

func (g *gitRepository) transitionAttachedBranch(
	ctx context.Context,
	branch string,
	oldHead string,
	newHead string,
) error {
	if !isObjectID(newHead) {
		return fmt.Errorf("attached transition target is not a full Git object ID")
	}
	oldTree, err := g.tree(ctx, oldHead)
	if err != nil {
		return err
	}
	newTree, err := g.tree(ctx, newHead)
	if err != nil {
		return err
	}
	if err := g.requireDefaultIndexFlags(ctx, []string{"."}); err != nil {
		return fmt.Errorf("attached transition index flag gate: %w", err)
	}
	observedSourceTree, err := g.indexTree(ctx)
	if err != nil || observedSourceTree != oldTree {
		return fmt.Errorf(
			"attached transition source index gate: got %s, want %s: %w",
			observedSourceTree,
			oldTree,
			err,
		)
	}
	if err := g.requireWorktreeMatchesIndex(ctx); err != nil {
		return fmt.Errorf("attached transition source worktree gate: %w", err)
	}
	if err := g.requireBranchHead(ctx, branch, oldHead); err != nil {
		return err
	}
	if _, err := g.run(
		ctx,
		"read-tree",
		"-m",
		"-u",
		oldHead,
		newHead,
	); err != nil {
		return fmt.Errorf("transition attached branch tree: %w", err)
	}
	if _, err := g.run(
		ctx,
		"update-ref",
		"-m",
		"repo_delivery: transition attached branch",
		"refs/heads/"+branch,
		newHead,
		oldHead,
	); err != nil {
		return fmt.Errorf("transition attached branch with exact compare-and-swap: %w", err)
	}
	if err := g.requireBranchHead(ctx, branch, newHead); err != nil {
		return fmt.Errorf("attached branch changed during transition: %w", err)
	}
	if observed, err := g.indexTree(ctx); err != nil || observed != newTree {
		return fmt.Errorf(
			"attached transition target index gate: got %s, want %s: %w",
			observed,
			newTree,
			err,
		)
	}
	return nil
}

func (g *gitRepository) reconcileRebasedBranch(
	ctx context.Context,
	branch string,
	expectedHead string,
	newHead string,
) error {
	if err := g.requireCompleteHistory(ctx); err != nil {
		return err
	}
	if err := g.requireNoOperation(ctx); err != nil {
		return err
	}
	if err := g.requireDefaultIndexFlags(ctx, []string{"."}); err != nil {
		return fmt.Errorf("main rebase reconciliation index flag gate: %w", err)
	}
	status, err := g.status(ctx)
	if err != nil {
		return err
	}
	if !status.clean() {
		return fmt.Errorf("main worktree changed while the isolated rebase ran")
	}
	currentBranch, currentHead, err := g.branchHead(ctx)
	if err != nil {
		return err
	}
	if currentBranch != branch || currentHead != expectedHead {
		return fmt.Errorf("main branch or HEAD changed while the isolated rebase ran")
	}
	expectedTree, err := g.tree(ctx, expectedHead)
	if err != nil {
		return err
	}
	if err := g.requireIndexTree(ctx, expectedTree); err != nil {
		return fmt.Errorf("main rebase reconciliation index gate: %w", err)
	}
	if err := g.requireDefaultIndexFlags(ctx, []string{"."}); err != nil {
		return fmt.Errorf("final main rebase index flag gate: %w", err)
	}
	if err := g.transitionAttachedBranch(
		ctx,
		branch,
		expectedHead,
		newHead,
	); err != nil {
		return fmt.Errorf("install exact rebased checkout: %w", err)
	}
	return nil
}

func (g *gitRepository) push(
	ctx context.Context,
	remote string,
	localOID string,
	headRef string,
	remoteHeadOID string,
	baseRef string,
	baseOID string,
) error {
	if !isObjectID(localOID) {
		return fmt.Errorf("push source is not a full Git object ID")
	}
	boundRemote, err := bindGitRemote(remote)
	if err != nil {
		return err
	}
	return g.pushBound(
		ctx,
		boundRemote,
		localOID,
		headRef,
		remoteHeadOID,
		baseRef,
		baseOID,
	)
}

func (g *gitRepository) pushBound(
	ctx context.Context,
	remote boundGitRemote,
	localOID string,
	headRef string,
	remoteHeadOID string,
	baseRef string,
	baseOID string,
) (returnErr error) {
	if err := g.requireCompleteHistory(ctx); err != nil {
		return err
	}
	if !isObjectID(localOID) || !isObjectID(baseOID) {
		return fmt.Errorf("push source and base must be full Git object IDs")
	}
	if remoteHeadOID != "" && !isObjectID(remoteHeadOID) {
		return fmt.Errorf("expected remote head must be absent or a full Git object ID")
	}
	for name, ref := range map[string]string{
		"base": baseRef,
		"head": headRef,
	} {
		if !strings.HasPrefix(ref, "refs/heads/") {
			return fmt.Errorf("%s ref %q is not a full branch ref", name, ref)
		}
		if err := g.checkBranch(ctx, strings.TrimPrefix(ref, "refs/heads/")); err != nil {
			return fmt.Errorf("invalid %s ref: %w", name, err)
		}
	}

	transport, err := g.openIsolatedTransport(ctx, remote)
	if err != nil {
		return err
	}
	defer func() {
		if cleanupErr := transport.close(); cleanupErr != nil {
			returnErr = errors.Join(
				returnErr,
				fmt.Errorf("clean isolated push transport: %w", cleanupErr),
			)
		}
	}()
	before, err := transport.remoteObjectIDs(ctx, baseRef, headRef)
	if err != nil {
		return fmt.Errorf("verify remote refs immediately before push: %w", err)
	}
	if before[baseRef] != baseOID || before[headRef] != remoteHeadOID {
		return fmt.Errorf(
			"remote refs changed before push: base=%s head=%s, expected base=%s head=%s",
			objectIDForReport(before[baseRef]),
			objectIDForReport(before[headRef]),
			objectIDForReport(baseOID),
			objectIDForReport(remoteHeadOID),
		)
	}
	args := exactPushArguments(
		remote.endpoint,
		localOID,
		headRef,
		remoteHeadOID,
	)
	_, pushErr := transport.run(ctx, args...)
	after, inspectErr := transport.remoteObjectIDs(ctx, baseRef, headRef)
	if inspectErr != nil {
		if pushErr != nil {
			return fmt.Errorf(
				"push result is ambiguous and exact remote refs could not be reread: push: %v; reread: %w",
				pushErr,
				inspectErr,
			)
		}
		return fmt.Errorf(
			"push reported success, but exact remote refs could not be verified: %w",
			inspectErr,
		)
	}
	if after[headRef] != localOID {
		if pushErr != nil {
			return fmt.Errorf(
				"push feature ref with exact lease failed; observed base=%s head=%s: %w",
				objectIDForReport(after[baseRef]),
				objectIDForReport(after[headRef]),
				pushErr,
			)
		}
		return fmt.Errorf(
			"push reported success but remote head is %s, want %s; base is %s",
			objectIDForReport(after[headRef]),
			objectIDForReport(localOID),
			objectIDForReport(after[baseRef]),
		)
	}
	if after[baseRef] == baseOID {
		// A transport error can be reported after the server has accepted the
		// update. The exact reread reconciles that ambiguous result.
		return nil
	}
	if err := transport.rollbackPushedHead(
		ctx,
		headRef,
		localOID,
		remoteHeadOID,
		baseRef,
	); err != nil {
		return fmt.Errorf(
			"base changed from %s to %s while the feature ref was pushed; %w",
			objectIDForReport(baseOID),
			objectIDForReport(after[baseRef]),
			err,
		)
	}
	return fmt.Errorf(
		"base changed from %s to %s while the feature ref was pushed; the feature ref was restored to %s",
		objectIDForReport(baseOID),
		objectIDForReport(after[baseRef]),
		objectIDForReport(remoteHeadOID),
	)
}

func exactPushArguments(
	endpoint string,
	localOID string,
	headRef string,
	expectedHeadOID string,
) []string {
	return []string{
		"push",
		"--atomic",
		"--no-verify",
		"--no-follow-tags",
		"--no-push-option",
		"--recurse-submodules=no",
		"--force-with-lease=" + headRef + ":" + expectedHeadOID,
		"--",
		endpoint,
		localOID + ":" + headRef,
	}
}

func (t *isolatedGitTransport) rollbackPushedHead(
	ctx context.Context,
	headRef string,
	pushedOID string,
	priorOID string,
	baseRef string,
) error {
	if priorOID != pushedOID {
		refspec := ":" + headRef
		if priorOID != "" {
			refspec = priorOID + ":" + headRef
		}
		args := []string{
			"push",
			"--atomic",
			"--no-verify",
			"--no-follow-tags",
			"--no-push-option",
			"--recurse-submodules=no",
			"--force-with-lease=" + headRef + ":" + pushedOID,
			"--",
			t.remote.endpoint,
			refspec,
		}
		// An error may still mean the server accepted the rollback. Always
		// verify exact remote state below instead of trusting process status.
		_, _ = t.run(ctx, args...)
	}
	after, err := t.remoteObjectIDs(ctx, baseRef, headRef)
	if err != nil {
		return fmt.Errorf("cannot prove rollback state: %w", err)
	}
	if after[headRef] != priorOID {
		return fmt.Errorf(
			"rollback did not overwrite a concurrent head and could not restore the prior value; observed base=%s head=%s, wanted head=%s",
			objectIDForReport(after[baseRef]),
			objectIDForReport(after[headRef]),
			objectIDForReport(priorOID),
		)
	}
	return nil
}

func objectIDForReport(value string) string {
	if value == "" {
		return "<absent>"
	}
	return value
}

func isAllowedPath(value string, allowed []string) bool {
	clean := filepath.ToSlash(filepath.Clean(value))
	for _, candidate := range allowed {
		candidate = filepath.ToSlash(filepath.Clean(candidate))
		if clean == candidate || strings.HasPrefix(clean, candidate+"/") {
			return true
		}
	}
	return false
}

func validateTaskPaths(paths []string) ([]string, error) {
	result := make([]string, 0, len(paths))
	seen := make(map[string]bool)
	for _, value := range paths {
		if value == "" || filepath.IsAbs(value) {
			return nil, fmt.Errorf("task path %q must be a nonempty relative path", value)
		}
		clean := filepath.Clean(value)
		if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
			return nil, fmt.Errorf("task path %q escapes or names the repository root", value)
		}
		clean = filepath.ToSlash(clean)
		if clean == "out" || strings.HasPrefix(clean, "out/") {
			return nil, fmt.Errorf(
				"task path %q names reserved delivery scratch space",
				value,
			)
		}
		if !seen[clean] {
			seen[clean] = true
			result = append(result, clean)
		}
	}
	return result, nil
}

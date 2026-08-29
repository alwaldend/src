package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type command struct {
	Name             string
	Args             []string
	Dir              string
	Env              []string
	UnsetEnv         []string
	UnsetEnvPrefixes []string
	Stdin            string
	OutputLimit      int
}

type commandResult struct {
	Stdout    string
	Stderr    string
	ExitCode  int
	Truncated bool
}

type commandRunner interface {
	Run(context.Context, command) (commandResult, error)
}

type commandError struct {
	Command command
	Result  commandResult
	Err     error
}

func (e *commandError) Error() string {
	name := safeCommandName(e.Command.Name)
	if e.Result.Truncated {
		limit := e.Command.OutputLimit
		if limit <= 0 {
			limit = commandOutputLimit
		}
		return fmt.Sprintf(
			"%s produced more than %d bytes on stdout or stderr; refusing truncated data",
			name,
			limit,
		)
	}
	detail := strings.TrimSpace(e.Result.Stderr)
	if detail != "" {
		return fmt.Sprintf(
			"%s exited with %d: %s",
			name,
			e.Result.ExitCode,
			redactCredentials(detail),
		)
	}
	return fmt.Sprintf(
		"%s exited with %d: %s",
		name,
		e.Result.ExitCode,
		safeCommandFailureDetail(e.Err),
	)
}

func (e *commandError) Is(target error) bool {
	switch classifyCommandFailure(e.Err) {
	case commandFailureCanceled:
		return target == context.Canceled
	case commandFailureDeadline:
		return target == context.DeadlineExceeded
	case commandFailureNotFound:
		return target == exec.ErrNotFound || target == os.ErrNotExist
	case commandFailureWorkingDirectoryNotFound:
		return target == os.ErrNotExist
	case commandFailurePermission:
		return target == os.ErrPermission
	default:
		return false
	}
}

type execRunner struct{}

const commandOutputLimit = 64 * 1024

type limitedBuffer struct {
	buffer    bytes.Buffer
	limit     int
	truncated bool
}

func (b *limitedBuffer) Write(value []byte) (int, error) {
	written := len(value)
	remaining := b.limit - b.buffer.Len()
	if len(value) > remaining {
		b.truncated = true
	}
	if remaining > 0 {
		if remaining > len(value) {
			remaining = len(value)
		}
		_, _ = b.buffer.Write(value[:remaining])
	}
	return written, nil
}

func (b *limitedBuffer) String() string {
	return b.buffer.String()
}

func (r *execRunner) Run(
	ctx context.Context,
	request command,
) (commandResult, error) {
	cmd := exec.CommandContext(ctx, request.Name, request.Args...)
	cmd.Dir = request.Dir
	if len(request.Env) != 0 || len(request.UnsetEnv) != 0 ||
		len(request.UnsetEnvPrefixes) != 0 {
		cmd.Env = mergeEnvironment(
			os.Environ(),
			request.Env,
			request.UnsetEnv,
			request.UnsetEnvPrefixes,
		)
	}
	if request.Stdin != "" {
		cmd.Stdin = strings.NewReader(request.Stdin)
	}
	outputLimit := request.OutputLimit
	if outputLimit <= 0 {
		outputLimit = commandOutputLimit
	}
	stdout := limitedBuffer{limit: outputLimit}
	stderr := limitedBuffer{limit: outputLimit}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil && ctx.Err() != nil {
		err = ctx.Err()
	}
	result := commandResult{
		Stdout:    stdout.String(),
		Stderr:    stderr.String(),
		ExitCode:  0,
		Truncated: stdout.truncated || stderr.truncated,
	}
	if err == nil && !result.Truncated {
		return result, nil
	}
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		result.ExitCode = exitError.ExitCode()
	} else if err != nil {
		result.ExitCode = -1
	}
	if err == nil {
		err = fmt.Errorf("command output exceeded safety limit")
	}
	// Failed subprocess output and errors cross a diagnostic boundary. Keep
	// the successful stdout path unchanged because callers parse it, but do
	// not return credential-bearing failure output or retain the executable
	// path, arguments, environment, stdin, or working directory in the error.
	result.Stdout = redactCredentials(result.Stdout)
	result.Stderr = redactCredentials(result.Stderr)
	return result, &commandError{
		Command: command{
			Name:        safeCommandName(request.Name),
			OutputLimit: outputLimit,
		},
		Result: result,
		Err: &sanitizedCommandCause{
			detail: safeCommandFailureDetail(err),
			kind:   classifyCommandFailure(err),
		},
	}
}

func mergeEnvironment(
	base []string,
	overrides []string,
	unset []string,
	unsetPrefixes []string,
) []string {
	replacements := make(map[string]string, len(overrides))
	removed := make(map[string]bool, len(unset))
	for _, name := range unset {
		removed[name] = true
	}
	for _, value := range overrides {
		name, _, found := strings.Cut(value, "=")
		if found {
			replacements[name] = value
			delete(removed, name)
		}
	}
	result := make([]string, 0, len(base)+len(replacements))
	for _, value := range base {
		name, _, found := strings.Cut(value, "=")
		if !found || removed[name] || hasAnyPrefix(name, unsetPrefixes) {
			continue
		}
		if _, replace := replacements[name]; replace {
			continue
		}
		result = append(result, value)
	}
	for _, value := range overrides {
		result = append(result, value)
	}
	return result
}

func hasAnyPrefix(value string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(value, prefix) {
			return true
		}
	}
	return false
}

type sanitizedCommandCause struct {
	detail string
	kind   commandFailureKind
}

func (e *sanitizedCommandCause) Error() string {
	return e.detail
}

type commandFailureKind uint8

const (
	commandFailureUnknown commandFailureKind = iota
	commandFailureCanceled
	commandFailureDeadline
	commandFailureNotFound
	commandFailureWorkingDirectoryNotFound
	commandFailurePermission
)

func safeCommandName(value string) string {
	// An explicitly selected forge executable can be an arbitrary path. Only
	// retain names whose complete basename identifies a supported CLI; a
	// parent directory or custom wrapper name may itself contain a credential.
	value = strings.ReplaceAll(strings.TrimSpace(value), `\`, "/")
	if index := strings.LastIndex(value, "/"); index >= 0 {
		value = value[index+1:]
	}
	switch strings.ToLower(value) {
	case "git", "git.exe":
		return "git"
	case "gh", "gh.exe":
		return "gh"
	default:
		return "subprocess"
	}
}

func safeCommandFailureDetail(err error) string {
	if err == nil {
		return "command execution failed"
	}
	var sanitized *sanitizedCommandCause
	if errors.As(err, &sanitized) {
		return sanitized.detail
	}
	switch classifyCommandFailure(err) {
	case commandFailureCanceled:
		return "context canceled"
	case commandFailureDeadline:
		return "context deadline exceeded"
	case commandFailureNotFound:
		return "executable not found"
	case commandFailureWorkingDirectoryNotFound:
		return "working directory not found"
	case commandFailurePermission:
		return "permission denied"
	}
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		return fmt.Sprintf("exit status %d", exitError.ExitCode())
	}
	// Unknown error implementations can include an executable path or a
	// wrapped request value in Error(). Do not render untrusted error text.
	return "command execution failed"
}

func classifyCommandFailure(err error) commandFailureKind {
	if err == nil {
		return commandFailureUnknown
	}
	var sanitized *sanitizedCommandCause
	if errors.As(err, &sanitized) {
		return sanitized.kind
	}
	var pathError *os.PathError
	if errors.As(err, &pathError) && pathError.Op == "chdir" &&
		errors.Is(pathError.Err, os.ErrNotExist) {
		return commandFailureWorkingDirectoryNotFound
	}
	switch {
	case errors.Is(err, context.Canceled):
		return commandFailureCanceled
	case errors.Is(err, context.DeadlineExceeded):
		return commandFailureDeadline
	case errors.Is(err, exec.ErrNotFound), errors.Is(err, os.ErrNotExist):
		return commandFailureNotFound
	case errors.Is(err, os.ErrPermission):
		return commandFailurePermission
	default:
		return commandFailureUnknown
	}
}

var credentialHeaderPattern = regexp.MustCompile(
	`(?im)(^|[^a-z0-9_-])((proxy[-_])?authorization"?` +
		`[ \t]*[:=][ \t]*"?)(bearer|basic|token)[ \t]+[^\r\n"]*`,
)

var malformedURLCredentialPattern = regexp.MustCompile(
	`(?i)((git|https?|ssh)://)[^@\r\n"'<>]*:[^@\r\n"'<>]*@`,
)

var scpCredentialPattern = regexp.MustCompile(
	`(?i)(^|[\s"'(<>=,:;?&\[{])[^/@\s"'()<>]+@` +
		`((\[[a-z0-9:._%\-]+\]|[a-z0-9_.-]+):[^\s"'<>]*)`,
)

var credentialAssignmentKeys = map[string]bool{
	"access_token":            true,
	"auth_token":              true,
	"bearer_token":            true,
	"client_secret":           true,
	"gh_enterprise_token":     true,
	"gh_token":                true,
	"github_enterprise_token": true,
	"github_token":            true,
	"git_http_password":       true,
	"git_password":            true,
	"gitlab_token":            true,
	"oauth_token":             true,
	"passwd":                  true,
	"password":                true,
	"secret":                  true,
	"token":                   true,
}

func redactCredentials(value string) string {
	// Keep host/path diagnostics while removing any URL userinfo. This is
	// deliberately conservative and applies to all subprocess stderr.
	for _, scheme := range []string{
		"git://",
		"http://",
		"https://",
		"ssh://",
	} {
		lower := strings.ToLower(value)
		start := 0
		for {
			index := strings.Index(lower[start:], scheme)
			if index < 0 {
				break
			}
			index += start
			authorityStart := index + len(scheme)
			authorityEnd := len(value)
			for offset, character := range value[authorityStart:] {
				if strings.ContainsRune("/?# \r\n\t\"'<>\",)", character) {
					authorityEnd = authorityStart + offset
					break
				}
			}
			authority := value[authorityStart:authorityEnd]
			at := strings.LastIndex(authority, "@")
			if at >= 0 {
				value = value[:authorityStart] + "<redacted>@" +
					authority[at+1:] + value[authorityEnd:]
				lower = strings.ToLower(value)
				start = authorityStart + len("<redacted>@")
			} else {
				start = authorityEnd
			}
		}
	}
	value = malformedURLCredentialPattern.ReplaceAllString(
		value,
		`${1}<redacted>@`,
	)
	// Git also accepts SCP-like SSH endpoints. Although the common username
	// is simply "git", treating all userinfo as sensitive avoids leaking a
	// token supplied in that position while preserving the host and path.
	value = scpCredentialPattern.ReplaceAllString(
		value,
		`${1}<redacted>@${2}`,
	)
	// Debug output from HTTP clients can contain authentication headers, and
	// malformed-provider diagnostics can echo credential-shaped assignments.
	// Preserve the label so the failure remains actionable, but never its
	// value. These replacements are deliberately case-insensitive.
	value = credentialHeaderPattern.ReplaceAllString(
		value,
		`${1}${2}<redacted>`,
	)
	value = redactCredentialAssignments(value)
	return value
}

func redactCredentialAssignments(value string) string {
	var redacted strings.Builder
	written := 0
	for index := 0; index < len(value); {
		if !isCredentialIdentifierByte(value[index]) {
			index++
			continue
		}
		keyStart := index
		for index < len(value) && isCredentialIdentifierByte(value[index]) {
			index++
		}
		key := strings.ToLower(value[keyStart:index])
		key = strings.ReplaceAll(key, "-", "_")
		if !credentialAssignmentKeys[key] {
			continue
		}
		cursor := index
		if cursor < len(value) &&
			(value[cursor] == '\'' || value[cursor] == '"') {
			cursor++
		}
		for cursor < len(value) &&
			(value[cursor] == ' ' || value[cursor] == '\t') {
			cursor++
		}
		if cursor >= len(value) ||
			(value[cursor] != ':' && value[cursor] != '=') {
			continue
		}
		cursor++
		for cursor < len(value) &&
			(value[cursor] == ' ' || value[cursor] == '\t') {
			cursor++
		}
		quote := byte(0)
		if cursor < len(value) &&
			(value[cursor] == '\'' || value[cursor] == '"') {
			quote = value[cursor]
			cursor++
		}
		valueStart := cursor
		valueEnd := credentialAssignmentEnd(value, valueStart, quote)
		if valueEnd == valueStart {
			continue
		}
		redacted.WriteString(value[written:valueStart])
		redacted.WriteString("<redacted>")
		written = valueEnd
		index = valueEnd
	}
	if written == 0 {
		return value
	}
	redacted.WriteString(value[written:])
	return redacted.String()
}

func credentialAssignmentEnd(value string, start int, quote byte) int {
	for index := start; index < len(value); index++ {
		if value[index] == '\n' || value[index] == '\r' {
			return index
		}
		if quote != 0 {
			if value[index] == '\\' && index+1 < len(value) {
				index++
				continue
			}
			if value[index] == quote {
				return index
			}
			continue
		}
		if value[index] == ' ' || value[index] == '\t' ||
			strings.ContainsRune("&,;}]\"'", rune(value[index])) {
			return index
		}
	}
	return len(value)
}

func isCredentialIdentifierByte(value byte) bool {
	return value >= 'a' && value <= 'z' ||
		value >= 'A' && value <= 'Z' ||
		value >= '0' && value <= '9' ||
		value == '_' || value == '-'
}

func writeJSON(output io.Writer, value any) error {
	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

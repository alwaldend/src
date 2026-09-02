package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestRedactCredentialsForSupportedURLSchemes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "git",
			value: "git://token@git.example/owner/repo.git",
			want:  "git://<redacted>@git.example/owner/repo.git",
		},
		{
			name:  "http",
			value: "http://user:secret@git.example/owner/repo.git",
			want:  "http://<redacted>@git.example/owner/repo.git",
		},
		{
			name:  "https mixed case",
			value: "HTTPS://user:secret@git.example/owner/repo.git",
			want:  "HTTPS://<redacted>@git.example/owner/repo.git",
		},
		{
			name:  "ssh mixed case",
			value: "SsH://token@git.example/owner/repo.git",
			want:  "SsH://<redacted>@git.example/owner/repo.git",
		},
		{
			name: "malformed escaped userinfo",
			value: "parse \"https://user:secret%zz@git.example/repo\": " +
				"invalid URL escape \"%zz\"",
			want: "parse \"https://<redacted>@git.example/repo\": " +
				"invalid URL escape \"%zz\"",
		},
		{
			name:  "scp style",
			value: "fatal: user:secret@git.example:owner/repo.git denied",
			want:  "fatal: <redacted>@git.example:owner/repo.git denied",
		},
		{
			name:  "malformed raw slash",
			value: "https://user:secret/path@git.example/repo",
			want:  "https://<redacted>@git.example/repo",
		},
		{
			name:  "query email is not userinfo",
			value: "https://git.example?contact=user@example.com",
			want:  "https://git.example?contact=user@example.com",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := redactCredentials(test.value); got != test.want {
				t.Fatalf("redactCredentials() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestRedactCredentialsPreservesNonsecretAuthorizationDiagnostic(
	t *testing.T,
) {
	t.Parallel()
	value := "reauthorization: failed; authorization: missing repository scope"
	if got := redactCredentials(value); got != value {
		t.Fatalf("redactCredentials() = %q, want %q", got, value)
	}
}

func TestRedactCredentialsForHeadersAndAssignments(t *testing.T) {
	t.Parallel()
	secret := "header-secret"
	value := strings.Join([]string{
		"> Authorization: Bearer " + secret,
		"Proxy-Authorization: Basic " + secret,
		`http2: header field "authorization" = "Bearer ` + secret + `"`,
		`provider error: access_token="` + secret + `" rejected`,
		"GH_TOKEN=" + secret + " request failed",
		"GH_ENTERPRISE_TOKEN='part & two " + secret + "' request failed",
		`{"token":"part, two & ` + secret + `","reason":"rejected"}`,
	}, "\n")
	got := redactCredentials(value)
	if strings.Contains(got, secret) {
		t.Fatalf("redactCredentials() leaked credential in %q", got)
	}
	for _, diagnostic := range []string{
		"Authorization: <redacted>",
		"Proxy-Authorization: <redacted>",
		`header field "authorization" = "<redacted>"`,
		`access_token="<redacted>"`,
		"GH_TOKEN=<redacted>",
		"GH_ENTERPRISE_TOKEN='<redacted>'",
		`{"token":"<redacted>","reason":"rejected"}`,
	} {
		if !strings.Contains(got, diagnostic) {
			t.Errorf("redactCredentials() = %q, missing %q", got, diagnostic)
		}
	}
}

func TestRedactCredentialsForSCPVariants(t *testing.T) {
	t.Parallel()
	secret := "scp-secret"
	tests := map[string]string{
		"assignment": "remote=" + secret + "@git.example:owner/repo.git",
		"IPv6":       secret + "@[2001:db8::1]:owner/repo.git",
	}
	for name, value := range tests {
		name := name
		value := value
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := redactCredentials(value)
			if strings.Contains(got, secret) {
				t.Fatalf("redactCredentials() leaked credential in %q", got)
			}
			if !strings.Contains(got, "<redacted>@") {
				t.Fatalf("redactCredentials() = %q, want redacted userinfo", got)
			}
		})
	}
}

func TestRedactCredentialsHandlesMultipleSupportedURLs(t *testing.T) {
	t.Parallel()
	secret := "shared-secret"
	value := "fetch " +
		"ssh://user:" + secret + "@git.example/owner/one.git and " +
		"GIT://" + secret + "@git.example/owner/two.git failed"
	got := redactCredentials(value)
	if strings.Contains(got, secret) {
		t.Fatalf("redactCredentials() leaked credential in %q", got)
	}
	want := "fetch ssh://<redacted>@git.example/owner/one.git and " +
		"GIT://<redacted>@git.example/owner/two.git failed"
	if got != want {
		t.Fatalf("redactCredentials() = %q, want %q", got, want)
	}
}

func TestCommandErrorSanitizesEveryFormattingBranch(t *testing.T) {
	t.Parallel()
	secret := "name-and-cause-secret"
	tests := []struct {
		name string
		err  *commandError
		want string
	}{
		{
			name: "stderr",
			err: &commandError{
				Command: command{Name: "/credential-" + secret + "/git"},
				Result: commandResult{
					Stderr:   "Authorization: Bearer " + secret,
					ExitCode: 128,
				},
			},
			want: "git exited with 128: Authorization: <redacted>",
		},
		{
			name: "underlying execution error",
			err: &commandError{
				Command: command{Name: "/credential-" + secret + "/gh"},
				Result:  commandResult{ExitCode: -1},
				Err: fmt.Errorf(
					"fork/exec /credential-%s/gh: permission denied",
					secret,
				),
			},
			want: "gh exited with -1: command execution failed",
		},
		{
			name: "truncated output",
			err: &commandError{
				Command: command{
					Name:        "/credential-" + secret + "/custom-forge",
					OutputLimit: 7,
				},
				Result: commandResult{Truncated: true},
			},
			want: "subprocess produced more than 7 bytes on stdout or " +
				"stderr; refusing truncated data",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := test.err.Error()
			if got != test.want {
				t.Fatalf("commandError.Error() = %q, want %q", got, test.want)
			}
			if strings.Contains(got, secret) {
				t.Fatalf("commandError.Error() leaked credential in %q", got)
			}
		})
	}
}

func TestWrappedCommandErrorDoesNotExposeRawCause(t *testing.T) {
	t.Parallel()
	secret := "wrapped-secret"
	commandErr := &commandError{
		Command: command{Name: "/credential-" + secret + "/gh"},
		Result:  commandResult{ExitCode: -1},
		Err: fmt.Errorf(
			"fork/exec /credential-%s/gh: no such file",
			secret,
		),
	}
	wrapped := fmt.Errorf("start forge CLI: %w", commandErr)
	if got := wrapped.Error(); strings.Contains(got, secret) {
		t.Fatalf("wrapped command error leaked credential in %q", got)
	}
	if errors.Unwrap(commandErr) != nil {
		t.Fatal("commandError unexpectedly exposed its raw subprocess cause")
	}
	var recovered *commandError
	if !errors.As(wrapped, &recovered) || recovered != commandErr {
		t.Fatal("wrapped commandError was not recoverable for exit-code checks")
	}
}

func TestCommandErrorPreservesSafeErrorCategories(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		cause  error
		target error
	}{
		{name: "canceled", cause: context.Canceled, target: context.Canceled},
		{
			name:   "deadline",
			cause:  context.DeadlineExceeded,
			target: context.DeadlineExceeded,
		},
		{name: "not found", cause: os.ErrNotExist, target: os.ErrNotExist},
		{name: "exec not found", cause: exec.ErrNotFound, target: exec.ErrNotFound},
		{name: "permission", cause: os.ErrPermission, target: os.ErrPermission},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := &commandError{Err: test.cause}
			if !errors.Is(err, test.target) {
				t.Fatalf("errors.Is(%v) = false", test.target)
			}
		})
	}
}

func TestExecRunnerScrubsFailedCommandDiagnostics(t *testing.T) {
	t.Parallel()
	secret := "runner-secret"
	executable, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable() error = %v", err)
	}
	result, err := (&execRunner{}).Run(context.Background(), command{
		Name: executable,
		Args: []string{
			"-test.run=^TestExecRunnerHelperProcess$",
		},
		Env: []string{
			"REPO_DELIVERY_RUNNER_HELPER=1",
			"REPO_DELIVERY_RUNNER_STDERR=Authorization: Bearer " + secret,
		},
	})
	if err == nil {
		t.Fatal("execRunner.Run() unexpectedly succeeded")
	}
	if strings.Contains(result.Stderr, secret) ||
		strings.Contains(err.Error(), secret) {
		t.Fatalf(
			"failed command leaked credential: result = %#v, error = %q",
			result,
			err,
		)
	}
	if !strings.Contains(result.Stderr, "Authorization: <redacted>") {
		t.Fatalf(
			"stderr = %q, want redacted authorization diagnostic",
			result.Stderr,
		)
	}
	var commandErr *commandError
	if !errors.As(err, &commandErr) {
		t.Fatalf("error type = %T, want *commandError", err)
	}
	if len(commandErr.Command.Args) != 0 || commandErr.Command.Dir != "" ||
		len(commandErr.Command.Env) != 0 || commandErr.Command.Stdin != "" {
		t.Fatalf(
			"commandError retained sensitive invocation fields: %#v",
			commandErr.Command,
		)
	}
}

func TestExecRunnerHelperProcess(t *testing.T) {
	if os.Getenv("REPO_DELIVERY_RUNNER_HELPER") != "1" {
		return
	}
	_, _ = fmt.Fprintln(
		os.Stderr,
		os.Getenv("REPO_DELIVERY_RUNNER_STDERR"),
	)
	os.Exit(19)
}

func TestExecRunnerDoesNotExposeCredentialBearingExecutablePath(t *testing.T) {
	t.Parallel()
	secret := "missing-path-secret"
	_, err := (&execRunner{}).Run(context.Background(), command{
		Name: "/definitely-missing-" + secret + "/gh",
	})
	if err == nil {
		t.Fatal("execRunner.Run() unexpectedly succeeded")
	}
	got := err.Error()
	if strings.Contains(got, secret) {
		t.Fatalf("execRunner.Run() leaked executable path in %q", got)
	}
	if got != "gh exited with -1: executable not found" {
		t.Fatalf("execRunner.Run() error = %q", got)
	}
	if !errors.Is(err, os.ErrNotExist) || !errors.Is(err, exec.ErrNotFound) {
		t.Fatal("missing executable did not preserve its safe not-found category")
	}
}

func TestExecRunnerDistinguishesMissingWorkingDirectory(t *testing.T) {
	t.Parallel()
	executable, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable() error = %v", err)
	}
	secret := "missing-directory-secret"
	_, err = (&execRunner{}).Run(context.Background(), command{
		Name: executable,
		Dir:  "/definitely-missing-" + secret,
	})
	if err == nil {
		t.Fatal("execRunner.Run() unexpectedly succeeded")
	}
	if got := err.Error(); strings.Contains(got, secret) ||
		!strings.Contains(got, "working directory not found") {
		t.Fatalf("execRunner.Run() error = %q", got)
	}
	if !errors.Is(err, os.ErrNotExist) || errors.Is(err, exec.ErrNotFound) {
		t.Fatal("working-directory error did not preserve the safe category")
	}
}

func TestExecRunnerPreservesCanceledContext(t *testing.T) {
	t.Parallel()
	executable, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable() error = %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = (&execRunner{}).Run(ctx, command{Name: executable})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("execRunner.Run() error = %v, want context canceled", err)
	}
}

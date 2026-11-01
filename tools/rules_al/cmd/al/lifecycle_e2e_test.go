package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
	"golang.org/x/sys/unix"
)

// TestLifecycleCLI exercises the shared operator interface through the built AL
// executable and retains the generated bundle and command observations.
func TestLifecycleCLI(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	binary, err := runfiles.Rlocation(os.Getenv("AL_BINARY"))
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	repository := filepath.Join(directory, "repository")
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	runGit := func(args ...string) {
		t.Helper()
		cmd := exec.CommandContext(ctx, git, args...)
		cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null"}
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Git fixture: %v\n%s", err, output)
		}
	}
	runGit("init", "--quiet", "--initial-branch=master", repository)
	runGit("-C", repository, "config", "user.name", "AL Fixture")
	runGit("-C", repository, "config", "user.email", "fixture@example.invalid")
	runGit("-C", repository, "config", "commit.gpgsign", "false")
	payload := filepath.Join(repository, "payload.txt")
	if err := os.WriteFile(payload, []byte("AL lifecycle fixture.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit("-C", repository, "add", "payload.txt")
	runGit("-C", repository, "commit", "--quiet", "-m", "Initialize AL fixture")
	runGit("-C", repository, "remote", "add", "origin", "git@github.com:fixture/repository.git")
	runAL := func(args ...string) []byte {
		t.Helper()
		cmd := exec.CommandContext(ctx, binary, args...)
		cmd.Dir = repository
		cmd.Env = append(os.Environ(), "BUILD_WORKSPACE_DIRECTORY=", "BUILD_WORKING_DIRECTORY=", "VAULT_TOKEN=")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("AL %s: %v\n%s", args[0], err, output)
		}
		return output
	}
	version := runAL("version", "--repo", repository, "show", "--format", "json")
	if !strings.Contains(string(version), "0.0.0-dev") {
		t.Fatalf("incorrect version: %s", version)
	}
	for _, flag := range []string{"--repository", "--repo"} {
		if got := runAL("version", "show", flag, repository, "--format", "json"); !bytes.Equal(got, version) {
			t.Fatalf("version repository flag %s changed output: %s", flag, got)
		}
	}
	if help := runAL("help", "version", "show"); !strings.Contains(string(help), "--repository") {
		t.Fatalf("native version help omitted repository selection: %s", help)
	}
	provider := runAL("delivery", "--repository", repository, "provider")
	var detected struct {
		// AdapterAvailable records whether the local remote identifies a supported forge.
		AdapterAvailable bool `json:"adapter_available"`
	}
	if err := json.Unmarshal(provider, &detected); err != nil {
		t.Fatal(err)
	}
	if !detected.AdapterAvailable {
		t.Fatalf("provider missing: %s", provider)
	}
	var failureOutput, failureDiagnostic bytes.Buffer
	rejected := exec.CommandContext(ctx, binary, "delivery", "--repository", filepath.Join(directory, "missing-repository"), "provider")
	rejected.Stdout, rejected.Stderr = &failureOutput, &failureDiagnostic
	if err := rejected.Run(); err == nil || failureOutput.Len() != 0 || failureDiagnostic.Len() == 0 {
		t.Fatalf("failed delivery must leave stdout empty and report on stderr: %v; stdout=%s; stderr=%s", err, &failureOutput, &failureDiagnostic)
	}
	status := filepath.Join(directory, "status")
	if err := os.WriteFile(status, []byte("STABLE_VERSION 0.0.0-dev\nSTABLE_VERSION_CHANNEL development\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	bundle := filepath.Join(directory, "bundle")
	runAL("release", "generate", "--project", "projects/fixture", "--version_file", status, "--add_file", payload, "--output_dir", bundle)
	data, err := os.ReadFile(filepath.Join(bundle, "files/payload.txt"))
	if err != nil || string(data) != "AL lifecycle fixture.\n" {
		t.Fatalf("bundle payload %q: %v", data, err)
	}
	arguments := filepath.Join(directory, "bazel-arguments")
	stub := filepath.Join(directory, "bin")
	if err := os.Mkdir(stub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(stub, "bazel"), []byte("#!/bin/sh\nprintf '%s\\n' \"$@\" > \"$AL_BAZEL_ARGUMENTS\"\nexit 23\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.CommandContext(ctx, binary, "bazel", "run", "//fixture:command", "--", "literal argument")
	cmd.Env = append(os.Environ(), "PATH="+stub+":"+os.Getenv("PATH"), "AL_BAZEL_ARGUMENTS="+arguments)
	if output, err := cmd.CombinedOutput(); err == nil || cmd.ProcessState.ExitCode() != 23 {
		t.Fatalf("Bazel exit propagation: %v\n%s", err, output)
	}
	forwarded, err := os.ReadFile(arguments)
	if err != nil || string(forwarded) != "run\n--config=agent\n//fixture:command\n--\nliteral argument\n" {
		t.Fatalf("Bazel arguments %q: %v", forwarded, err)
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		manifest, err := os.ReadFile(filepath.Join(bundle, "release.json"))
		if err != nil {
			t.Fatal(err)
		}
		for name, body := range map[string][]byte{"al-version.json": version, "al-provider.json": provider, "al-release.json": manifest, "al-bazel-arguments.txt": forwarded, "al-failure.txt": failureDiagnostic.Bytes()} {
			if err := os.WriteFile(filepath.Join(output, name), body, 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
}

// TestVersionCancellationHelper is the subprocess boundary used by cancellation fixtures.
func TestVersionCancellationHelper(t *testing.T) {
	marker := os.Getenv("AL_CANCELLATION_MARKER")
	if marker == "" {
		return
	}
	if err := os.WriteFile(marker, []byte(strconv.Itoa(os.Getpid())), 0o600); err != nil {
		t.Fatal(err)
	}
	time.Sleep(30 * time.Second)
}

// TestVersionCancellation exercises operating-system signals through AL and its active child.
func TestVersionCancellation(t *testing.T) {
	binary, err := runfiles.Rlocation(os.Getenv("AL_BINARY"))
	if err != nil {
		t.Fatal(err)
	}
	helper, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for _, operation := range []string{"show", "release-publish", "bazel"} {
		for _, signal := range []syscall.Signal{syscall.SIGINT, syscall.SIGTERM} {
			t.Run(fmt.Sprintf("%s-%d", operation, signal), func(t *testing.T) {
				directory := t.TempDir()
				marker := filepath.Join(directory, "child.pid")
				for _, name := range []string{"git", "bazel"} {
					if err := os.WriteFile(filepath.Join(directory, name), []byte("#!/bin/sh\nexec \"$AL_CANCELLATION_TEST\" -test.run=^TestVersionCancellationHelper$\n"), 0o755); err != nil {
						t.Fatal(err)
					}
				}
				args := []string{"version", "--repo", directory, operation}
				if operation == "bazel" {
					args = append(args, "--", "build", "//fixture")
				}
				if operation == "release-publish" {
					plan := filepath.Join(directory, "plan.json")
					body := `{"schema":"versioning/release_ref_plan/v1","kind":"ReleaseRefPlan","version":"2026.44.0","channel":"release","commit":"` + strings.Repeat("a", 40) + `","treeState":"clean","refs":[{"ref":"refs/tags/v2026.44.0","expected":"` + strings.Repeat("a", 40) + `","operation":"create"}],"atomic":true,"lease":"fixture"}`
					if err := os.WriteFile(plan, []byte(body), 0o600); err != nil {
						t.Fatal(err)
					}
					args = append(args, "--plan", plan)
				}
				command := exec.Command(binary, args...)
				command.Env = append(os.Environ(), "PATH="+directory+":"+os.Getenv("PATH"), "BUILD_WORKSPACE_DIRECTORY=", "AL_CANCELLATION_TEST="+helper, "AL_CANCELLATION_MARKER="+marker)
				command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
				var output bytes.Buffer
				command.Stdout, command.Stderr = &output, &output
				if err := command.Start(); err != nil {
					t.Fatal(err)
				}
				finished := make(chan error, 1)
				go func() { finished <- command.Wait() }()
				t.Cleanup(func() { _ = syscall.Kill(-command.Process.Pid, syscall.SIGKILL) })
				deadline := time.Now().Add(5 * time.Second)
				var child int
				for child == 0 && time.Now().Before(deadline) {
					if content, err := os.ReadFile(marker); err == nil {
						child, _ = strconv.Atoi(string(content))
					}
					time.Sleep(10 * time.Millisecond)
				}
				if child == 0 {
					t.Fatal("version subprocess did not become ready")
				}
				if err := command.Process.Signal(signal); err != nil {
					t.Fatal(err)
				}
				select {
				case err := <-finished:
					if err == nil || syscall.Kill(child, 0) == nil {
						t.Fatalf("cancellation left successful AL or a running child: %v; %s", err, &output)
					}
				case <-time.After(3 * time.Second):
					t.Fatal("version operation ignored cancellation")
				}
			})
		}
	}
}

// TestVersionTerminalGit imitates an SSH helper reading an authentication response.
func TestVersionTerminalGit(t *testing.T) {
	if os.Getenv("AL_TERMINAL_GIT") == "" {
		return
	}
	terminal, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer terminal.Close()
	if err := os.WriteFile(os.Getenv("AL_TERMINAL_READY"), []byte(strconv.Itoa(os.Getpid())), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := terminal.WriteString("Authentication response: "); err != nil {
		t.Fatal(err)
	}
	input := make([]byte, 32)
	if _, err := terminal.Read(input); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(os.Getenv("AL_TERMINAL_RESPONSE"), input, 0o600); err != nil {
		t.Fatal(err)
	}
	os.Exit(1) // Never produce a valid Git result or reach a remote operation.
}

// TestVersionTerminalSession owns the fixture terminal and checks its restoration.
func TestVersionTerminalSession(t *testing.T) {
	if os.Getenv("AL_TERMINAL_SESSION") == "" {
		return
	}
	binary := os.Getenv("AL_TERMINAL_BINARY")
	command := exec.Command(binary, "version", "show")
	command.Stdin, command.Stdout, command.Stderr = os.Stdin, os.Stdout, os.Stderr
	command.Env = append(os.Environ(), "AL_TERMINAL_SESSION=", "AL_TERMINAL_GIT=1")
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	if number := os.Getenv("AL_TERMINAL_SIGNAL"); number != "0" {
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			if _, err := os.Stat(os.Getenv("AL_TERMINAL_READY")); err == nil {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		value, err := strconv.Atoi(number)
		if err != nil {
			t.Fatal(err)
		}
		if err := command.Process.Signal(syscall.Signal(value)); err != nil {
			t.Fatal(err)
		}
	}
	if err := command.Wait(); err == nil {
		t.Fatal("fixture Git failure was lost")
	}
	foreground, err := unix.IoctlGetInt(0, unix.TIOCGPGRP)
	if err != nil || foreground != syscall.Getpgrp() {
		t.Fatalf("terminal not restored: %d, %v", foreground, err)
	}
	if err := os.WriteFile(os.Getenv("AL_TERMINAL_RESULT"), []byte("terminal restored\n"), 0o600); err != nil {
		t.Fatal(err)
	}
}

// TestVersionTerminal verifies prompts, signal cancellation, and terminal restoration through AL.
func TestVersionTerminal(t *testing.T) {
	binary, err := runfiles.Rlocation(os.Getenv("AL_BINARY"))
	if err != nil {
		t.Fatal(err)
	}
	helper, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for _, sig := range []syscall.Signal{0, syscall.SIGINT, syscall.SIGTERM} {
		t.Run(strconv.Itoa(int(sig)), func(t *testing.T) {
			directory := t.TempDir()
			master, err := os.OpenFile("/dev/ptmx", os.O_RDWR|syscall.O_NOCTTY, 0)
			if err != nil {
				t.Fatal(err)
			}
			defer master.Close()
			if err := unix.IoctlSetPointerInt(int(master.Fd()), unix.TIOCSPTLCK, 0); err != nil {
				t.Fatal(err)
			}
			number, err := unix.IoctlGetInt(int(master.Fd()), unix.TIOCGPTN)
			if err != nil {
				t.Fatal(err)
			}
			slave, err := os.OpenFile(fmt.Sprintf("/dev/pts/%d", number), os.O_RDWR|syscall.O_NOCTTY, 0)
			if err != nil {
				t.Fatal(err)
			}
			defer slave.Close()
			stub := filepath.Join(directory, "git")
			if err := os.WriteFile(stub, []byte("#!/bin/sh\nexec \"$AL_TERMINAL_HELPER\" -test.run=^TestVersionTerminalGit$\n"), 0o755); err != nil {
				t.Fatal(err)
			}
			ready, response, result := filepath.Join(directory, "ready"), filepath.Join(directory, "response"), filepath.Join(directory, "result")
			ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
			defer cancel()
			command := exec.CommandContext(ctx, helper, "-test.run=^TestVersionTerminalSession$")
			command.Stdin, command.Stdout, command.Stderr = slave, slave, slave
			command.SysProcAttr = &syscall.SysProcAttr{Setsid: true, Setctty: true, Ctty: 0}
			command.Env = append(os.Environ(), "PATH="+directory+":"+os.Getenv("PATH"), "AL_TERMINAL_SESSION=1", "AL_TERMINAL_BINARY="+binary, "AL_TERMINAL_HELPER="+helper, "AL_TERMINAL_READY="+ready, "AL_TERMINAL_RESPONSE="+response, "AL_TERMINAL_RESULT="+result, "AL_TERMINAL_SIGNAL="+strconv.Itoa(int(sig)))
			if err := command.Start(); err != nil {
				t.Fatal(err)
			}
			deadline := time.Now().Add(5 * time.Second)
			var child int
			for child == 0 && time.Now().Before(deadline) {
				if data, err := os.ReadFile(ready); err == nil {
					child, _ = strconv.Atoi(string(data))
				}
				time.Sleep(10 * time.Millisecond)
			}
			defer func() {
				if child > 0 {
					_ = syscall.Kill(-child, syscall.SIGKILL)
				}
				_ = syscall.Kill(-command.Process.Pid, syscall.SIGKILL)
			}()
			if child == 0 {
				t.Fatal("terminal Git helper did not start")
			}
			if sig == 0 {
				if _, err := master.WriteString("confirmed\n"); err != nil {
					t.Fatal(err)
				}
			}
			if err := command.Wait(); err != nil {
				t.Fatalf("terminal session: %v", err)
			}
			if _, err := os.Stat(result); err != nil {
				t.Fatal(err)
			}
			if sig == 0 {
				if data, err := os.ReadFile(response); err != nil || !bytes.Contains(data, []byte("confirmed")) {
					t.Fatalf("terminal response: %q, %v", data, err)
				}
			}
			if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
				if err := os.WriteFile(filepath.Join(output, fmt.Sprintf("terminal-%d.txt", sig)), []byte("prompt or cancellation completed; terminal restored\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

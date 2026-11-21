package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
)

// testOutputHeader starts a failed-test log copied by Bazel's console UI.
const testOutputHeader = "==================== Test output for "

// testOutputFooter ends a failed-test log copied by Bazel's console UI.
const testOutputFooter = "================================================================================"

// routineBazelInfo matches invocation bookkeeping with optional timestamps and colors.
var routineBazelInfo = regexp.MustCompile(`^(?:\([0-9:]+\) )?(?:\x1b\[[0-9;]*m)*INFO: (?:\x1b\[[0-9;]*m)*(?:Invocation ID:|Current date is )`)

// runWithForgejoLogs streams merged Bazel output to the current log group.
func runWithForgejoLogs(cmd *exec.Cmd, output io.Writer, testPhase bool) error {
	logs, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("open Bazel output: %w", err)
	}
	cmd.Stderr = cmd.Stdout
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start Bazel: %w", err)
	}
	streamErr := streamForgejoLogs(logs, output, testPhase)
	if streamErr != nil {
		streamErr = errors.Join(streamErr, logs.Close())
	}
	if err := errors.Join(streamErr, cmd.Wait()); err != nil {
		return fmt.Errorf("finish Bazel: %w", err)
	}
	return nil
}

// streamForgejoLogs folds failed-test diagnostics and omits routine invocation metadata.
func streamForgejoLogs(input io.Reader, output io.Writer, testPhase bool) error {
	reader := bufio.NewReader(input)
	open := false
	for {
		line, readErr := reader.ReadString('\n')
		text := strings.TrimRight(line, "\r\n")
		if target, found := strings.CutPrefix(text, testOutputHeader); testPhase && found && !open && strings.HasSuffix(target, ":") {
			if _, err := fmt.Fprintf(output, "##[group]Failed test %s\n", strings.TrimSuffix(target, ":")); err != nil {
				return fmt.Errorf("begin failed-test log group: %w", err)
			}
			open = true
		}
		if open || !routineBazelInfo.MatchString(text) {
			if _, err := io.WriteString(output, line); err != nil {
				return fmt.Errorf("write Bazel output: %w", err)
			}
		}
		if open && (text == testOutputFooter || readErr != nil) {
			if _, err := fmt.Fprint(output, "\n##[endgroup]\n"); err != nil {
				return fmt.Errorf("end failed-test log group: %w", err)
			}
			open = false
		}
		if errors.Is(readErr, io.EOF) {
			return nil
		}
		if readErr != nil {
			return fmt.Errorf("read Bazel output: %w", readErr)
		}
	}
}

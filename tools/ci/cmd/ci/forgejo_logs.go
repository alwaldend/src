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

// bazelDiagnostics is the group title for routine records and their continuation lines.
const bazelDiagnostics = "Bazel diagnostics"

// bazelConsoleFormatting matches the optional timestamp and color sequences.
var bazelConsoleFormatting = regexp.MustCompile(`^\([0-9:]+\) |\x1b\[[0-9;]*m`)

// routineBazelRecord starts a diagnostic record that can be folded.
var routineBazelRecord = regexp.MustCompile(`^(INFO|DEBUG|WARNING):`)

// visibleBazelRecord starts an error or a test result to show directly in the phase.
var visibleBazelRecord = regexp.MustCompile(`^((ERROR|FATAL|FAIL|FAILED):|Executed [0-9]+|(@[^ ]*)?//[^ ]+ +(?:\(cached\) +)?(PASSED|FAILED|FLAKY|TIMEOUT|INCOMPLETE|NO STATUS|REMOTE FAILURE)( |$))`)

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

// streamForgejoLogs folds routine diagnostics and failed-test output within a phase.
func streamForgejoLogs(input io.Reader, output io.Writer, testPhase bool) error {
	reader := bufio.NewReader(input)
	group := ""
	setGroup := func(title string) error {
		if group == title {
			return nil
		}
		if group != "" {
			if _, err := fmt.Fprint(output, "\n##[endgroup]\n"); err != nil {
				return fmt.Errorf("end Bazel log group: %w", err)
			}
		}
		if title != "" {
			if _, err := fmt.Fprintf(output, "##[group]%s\n", title); err != nil {
				return fmt.Errorf("begin Bazel log group: %w", err)
			}
		}
		group = title
		return nil
	}
	fold := true
	failedTest := false
	for {
		fragment, readErr := reader.ReadSlice('\n')
		line := string(fragment)
		text := bazelConsoleFormatting.ReplaceAllString(strings.TrimRight(line, "\r\n"), "")
		if line != "" && !failedTest {
			title := ""
			if target, found := strings.CutPrefix(text, testOutputHeader); testPhase && found && strings.HasSuffix(target, ":") {
				title = "Failed test " + strings.TrimSuffix(target, ":")
				failedTest = true
			} else {
				switch {
				case visibleBazelRecord.MatchString(text):
					fold = false
				case routineBazelRecord.MatchString(text):
					fold = true
				}
				if fold {
					title = bazelDiagnostics
				}
			}
			if err := setGroup(title); err != nil {
				return fmt.Errorf("group Bazel output: %w", err)
			}
		}
		if _, err := io.WriteString(output, line); err != nil {
			return fmt.Errorf("write Bazel output: %w", err)
		}
		// Only the first fragment of a long line can begin a Bazel record.
		// Stream its continuation without retaining a copy of the entire line.
		for errors.Is(readErr, bufio.ErrBufferFull) {
			fragment, readErr = reader.ReadSlice('\n')
			if _, err := output.Write(fragment); err != nil {
				return fmt.Errorf("write Bazel output continuation: %w", err)
			}
		}
		if readErr != nil || (failedTest && text == testOutputFooter) {
			if err := setGroup(""); err != nil {
				return fmt.Errorf("close Bazel output group: %w", err)
			}
			failedTest = false
		}
		if errors.Is(readErr, io.EOF) {
			return nil
		}
		if readErr != nil {
			return fmt.Errorf("read Bazel output: %w", readErr)
		}
	}
}

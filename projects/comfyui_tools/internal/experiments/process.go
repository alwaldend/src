package experiments

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

// boundedBuffer limits captured subprocess output without leaking it into diagnostics.
type boundedBuffer struct {
	// Buffer retains at most the configured byte limit.
	bytes.Buffer
	// limit is the maximum capture size in bytes.
	limit int
	// exceeded records output truncation even if the child ignores a write error.
	exceeded bool
}

// Write captures bounded private output and rejects oversized subprocess responses.
func (b *boundedBuffer) Write(p []byte) (int, error) {
	if len(p) > b.limit-b.Len() {
		b.exceeded = true
		return 0, fmt.Errorf("subprocess output limit exceeded")
	}
	n, err := b.Buffer.Write(p)
	if err != nil {
		return n, fmt.Errorf("capture subprocess output: %w", err)
	}
	return n, nil
}

// invoke runs literal argv in an owned process group and stores bounded private receipts.
func invoke(ctx context.Context, argv []string, input []byte, receipt string, seconds int) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(seconds)*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil && err != syscall.ESRCH {
			return fmt.Errorf("stop local command group: %w", err)
		}
		return nil
	}
	cmd.WaitDelay = 2 * time.Second
	cmd.Stdin = bytes.NewReader(input)
	output := &boundedBuffer{limit: 32 << 20}
	diagnostic := &boundedBuffer{limit: 1 << 20}
	cmd.Stdout, cmd.Stderr = output, diagnostic
	runErr := cmd.Run()
	if err := writePrivate(receipt+".json", output.Bytes()); err != nil {
		return nil, fmt.Errorf("retain command receipt: %w", err)
	}
	if err := writePrivate(receipt+".log", diagnostic.Bytes()); err != nil {
		return nil, fmt.Errorf("retain private command diagnostic: %w", err)
	}
	if output.exceeded || diagnostic.exceeded {
		return nil, fmt.Errorf("command output limit exceeded")
	}
	if runErr != nil {
		return output.Bytes(), fmt.Errorf("command failed; inspect private %s log: %w", filepath.Base(receipt), runErr)
	}
	return output.Bytes(), nil
}

// envelope is the native ComfyUI CLI JSON protocol, independent of its diagnostics.
type envelope struct {
	// OK identifies a successful native command response.
	OK bool `json:"ok"`
	// Data contains the command-specific native result.
	Data json.RawMessage `json:"data"`
	// Error preserves structured native outcomes from unsuccessful commands.
	Error struct {
		// Code distinguishes execution outcomes from command transport errors.
		Code string `json:"code"`
		// Details retains the command-specific summary without exposing diagnostics.
		Details json.RawMessage `json:"details"`
	} `json:"error"`
}

// cliFailure retains an unsuccessful native envelope and its process failure cause.
type cliFailure struct {
	// code identifies the native failure category.
	code string
	// details contains the private structured native summary.
	details json.RawMessage
	// cause preserves a failed CLI process or wrapper error.
	cause error
}

// Error describes the failed protocol operation without disclosing its private summary.
func (e *cliFailure) Error() string {
	return "CLI returned an unsuccessful envelope; inspect private receipt"
}

// Unwrap preserves the subprocess cause of an unsuccessful native envelope.
func (e *cliFailure) Unwrap() error {
	return e.cause
}

// callCLI accepts a successful native receipt even if its wrapper subsequently fails cleanup.
func callCLI(ctx context.Context, config configuration, receipt string, result any, args ...string) (bool, error) {
	argv := append(append([]string{}, config.CLI...), "--json")
	argv = append(argv, args...)
	data, runErr := invoke(ctx, argv, nil, receipt, config.TimeoutSeconds)
	var response envelope
	if err := decodeJSON(bytes.NewReader(data), &response); err != nil {
		return false, fmt.Errorf("invalid CLI receipt: %w", err)
	}
	if !response.OK {
		return false, &cliFailure{code: response.Error.Code, details: response.Error.Details, cause: runErr}
	}
	if err := decodeJSON(bytes.NewReader(response.Data), result); err != nil {
		return false, fmt.Errorf("decode CLI result: %w", err)
	}
	return runErr != nil, nil
}

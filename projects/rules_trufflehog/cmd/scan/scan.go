package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// findingEvent preserves the original result and its explicit policy disposition.
type findingEvent struct {
	// Type distinguishes finding records from the final summary.
	Type string `json:"type"`
	// Suppressed indicates an exact allowlist match.
	Suppressed bool `json:"suppressed"`
	// Reason is the approved rationale, when a result is suppressed.
	Reason string `json:"reason,omitempty"`
	// Finding is the unredacted original scanner record.
	Finding json.RawMessage `json:"finding"`
	// MetadataError explains an invalid location while retaining the original result.
	MetadataError string `json:"metadata_error,omitempty"`
}

// scanSummary reports complete scans without confusing suppression with detection.
type scanSummary struct {
	// Type identifies the final summary event.
	Type string `json:"type"`
	// Total counts every detected finding, including approved locations.
	Total int `json:"total"`
	// Suppressed counts exact policy matches.
	Suppressed int `json:"suppressed"`
	// Unmatched counts findings that fail the repository gate.
	Unmatched int `json:"unmatched"`
	// Verification records the wrapper's enforced scan-only mode.
	Verification bool `json:"verification"`
}

// scan owns one scanner process and fails closed on process or result errors.
func scan(parent context.Context, opts options, output, diagnostics io.Writer) (resultErr error) {
	approved, err := loadAllowlist(opts.allowlist)
	if err != nil {
		return fmt.Errorf("load finding policy: %w", err)
	}
	base := os.Getenv("TEST_TMPDIR")
	if base == "" {
		base = os.Getenv("TMPDIR")
	}
	if base == "" {
		return errors.New("set TMPDIR for task-owned scanner temporary storage")
	}
	scratch, err := os.MkdirTemp(base, "trufflehog-scan-")
	if err != nil {
		return fmt.Errorf("create scanner temporary directory: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(scratch); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("remove scanner temporary directory: %w", err))
		}
	}()
	ctx, cancel := context.WithTimeout(parent, opts.timeout)
	defer cancel()
	args := []string{"git", "--json", "--no-color", "--no-update", "--no-verification", "--fail", "--fail-on-scan-errors"}
	if !opts.allRefs {
		args = append(args, "--branch", "HEAD", "--skip-additional-refs")
	}
	args = append(args, (&url.URL{Scheme: "file", Path: opts.repository}).String())
	cmd := exec.CommandContext(ctx, opts.scanner, args...)
	cmd.WaitDelay = time.Second
	cmd.Env = append(os.Environ(), "TMPDIR="+scratch, "TMP="+scratch, "TEMP="+scratch)
	if opts.git != "" {
		absolute, err := filepath.Abs(opts.git)
		if err != nil {
			return fmt.Errorf("resolve Git executable: %w", err)
		}
		if err := os.Symlink(absolute, filepath.Join(scratch, "git")); err != nil {
			return fmt.Errorf("prepare scanner Git executable: %w", err)
		}
		cmd.Env = append(cmd.Env, "PATH="+scratch+string(os.PathListSeparator)+os.Getenv("PATH"))
	}
	cmd.Stderr = diagnostics
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("open scanner output: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start scanner: %w", err)
	}
	summary, streamErr := classify(pipe, output, approved)
	if streamErr != nil {
		cancel()
	}
	waitErr := cmd.Wait()
	if streamErr != nil {
		return fmt.Errorf("read scanner results: %w", streamErr)
	}
	if ctx.Err() != nil {
		return fmt.Errorf("scanner canceled: %w", errors.Join(ctx.Err(), waitErr))
	}
	status := 0
	if waitErr != nil {
		var exited *exec.ExitError
		if !errors.As(waitErr, &exited) || exited.ExitCode() != 183 {
			return fmt.Errorf("scanner process failed: %w", errors.Join(waitErr, ctx.Err()))
		}
		status = 183
	}
	if (summary.Total == 0 && status != 0) || (summary.Total > 0 && status != 183) {
		return errors.New("scanner result count contradicts its exit status")
	}
	if err := json.NewEncoder(output).Encode(summary); err != nil {
		return fmt.Errorf("write scan summary: %w", err)
	}
	if summary.Unmatched > 0 {
		return fmt.Errorf("%d results: %w", summary.Unmatched, errFindings)
	}
	return nil
}

// classify streams complete original findings and their exact allowlist matches.
func classify(input io.Reader, output io.Writer, approved map[location]string) (scanSummary, error) {
	summary := scanSummary{Type: "summary"}
	decoder, encoder := json.NewDecoder(input), json.NewEncoder(output)
	for {
		var raw json.RawMessage
		if err := decoder.Decode(&raw); errors.Is(err, io.EOF) {
			return summary, nil
		} else if err != nil {
			return summary, fmt.Errorf("decode scanner JSON: %w", err)
		}
		key, err := findingLocation(raw)
		if err != nil {
			if writeErr := encoder.Encode(findingEvent{Type: "finding", Finding: raw, MetadataError: err.Error()}); writeErr != nil {
				return summary, fmt.Errorf("preserve malformed finding: %w", errors.Join(err, writeErr))
			}
			return summary, fmt.Errorf("validate scanner result: %w", err)
		}
		reason, suppressed := approved[key]
		summary.Total++
		if suppressed {
			summary.Suppressed++
		} else {
			summary.Unmatched++
		}
		if err := encoder.Encode(findingEvent{Type: "finding", Suppressed: suppressed, Reason: reason, Finding: raw}); err != nil {
			return summary, fmt.Errorf("write original finding: %w", err)
		}
	}
}

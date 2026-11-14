package experiments

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// configuration is the serialized execution contract, stored privately for resume.
type configuration struct {
	// CLI is a literal executable and argument prefix for the existing native CLI.
	CLI []string `json:"cli"`
	// Producer is a literal executable and argument prefix supplying new workflows.
	Producer []string `json:"producer"`
	// QueuedBuffer is the desired number of pending jobs beyond one executing job.
	QueuedBuffer int `json:"queued_buffer"`
	// MaxExperiments bounds accepted workflows in a finite run.
	MaxExperiments int `json:"max_experiments"`
	// Continuous explicitly permits an unbounded experiment stream.
	Continuous bool `json:"continuous"`
	// TimeoutSeconds bounds each CLI or producer operation, including jobs wait.
	TimeoutSeconds int `json:"timeout_seconds"`
	// DuplicateAttempts bounds consecutive producer candidates with known identities.
	DuplicateAttempts int `json:"duplicate_attempts"`
}

// candidate is the producer's model-independent JSON response.
type candidate struct {
	// Workflow names an absolute native UI or API workflow file.
	Workflow string `json:"workflow"`
	// Description records the intended difference from earlier experiments.
	Description string `json:"description"`
	// Notes retains caller-owned prompts, settings and review context.
	Notes map[string]any `json:"notes,omitempty"`
	// IdentityExclusions names existing JSON pointers omitted from duplicate identity.
	IdentityExclusions []string `json:"identity_exclusions,omitempty"`
	// Done explicitly ends the producer stream without submitting another job.
	Done bool `json:"done,omitempty"`
}

// media records a validated original file relative to its private run directory.
type media struct {
	// Path identifies the retained media within the run.
	Path string `json:"path"`
	// SHA256 verifies the original bytes.
	SHA256 string `json:"sha256"`
	// Bytes records the original length.
	Bytes int64 `json:"bytes"`
}

// experiment is a durable ownership and provenance record for one submission.
type experiment struct {
	// Sequence is the monotonic producer request that supplied this experiment.
	Sequence int `json:"sequence"`
	// Directory identifies the experiment's private artifact directory.
	Directory string `json:"directory"`
	// Candidate preserves the exact producer response for review.
	Candidate candidate `json:"candidate"`
	// Identity is the canonical workflow digest after validated exclusions.
	Identity string `json:"identity"`
	// WorkflowSHA256 verifies the submitted snapshot's exact bytes.
	WorkflowSHA256 string `json:"workflow_sha256"`
	// PromptID identifies the native job owned by this run.
	PromptID string `json:"prompt_id,omitempty"`
	// State distinguishes live observation failure from confirmed terminal execution failure.
	State string `json:"state"`
	// Created records snapshot creation in UTC.
	Created time.Time `json:"created"`
	// Submitted records receipt of a valid native prompt identifier.
	Submitted time.Time `json:"submitted,omitempty"`
	// Finished records the terminal local observation.
	Finished time.Time `json:"finished,omitempty"`
	// Media retains every validated download, without replacing originals.
	Media []media `json:"media,omitempty"`
	// Failure names a bounded failure phase without exposing raw diagnostics.
	Failure string `json:"failure,omitempty"`
	// CleanupWarning records a successful CLI envelope followed by a process error.
	CleanupWarning bool `json:"cleanup_warning,omitempty"`
}

// journal is the single durable owner of run state and configuration.
type journal struct {
	// Version identifies the journal schema.
	Version int `json:"version"`
	// Config preserves the normalized execution contract.
	Config configuration `json:"config"`
	// NextSequence prevents producer requests from repeating after interruption.
	NextSequence int `json:"next_sequence"`
	// Experiments records accepted workflows and their native ownership.
	Experiments []*experiment `json:"experiments"`
	// Duplicates counts rejected workflow identities.
	Duplicates int `json:"duplicates"`
	// Exhausted records an explicit producer end or reached finite limit.
	Exhausted bool `json:"exhausted"`
	// Failure records a run-level stop condition that requires inspection.
	Failure string `json:"failure,omitempty"`
	// Updated records the most recent durable state write.
	Updated time.Time `json:"updated"`
}

// control is an external request consumed by the controller without taking its lock.
type control struct {
	// Stop requests that the controller drain rather than submit new jobs.
	Stop bool `json:"stop"`
	// CancelOwned additionally requests cancellation of recorded outstanding IDs.
	CancelOwned bool `json:"cancel_owned"`
}

// decodeJSON preserves exact numeric values in opaque notes and native responses.
func decodeJSON(reader io.Reader, value any) error {
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()
	if err := decoder.Decode(value); err != nil {
		return fmt.Errorf("decode JSON value: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err != nil {
			return fmt.Errorf("decode trailing JSON data: %w", err)
		}
		return fmt.Errorf("expected exactly one JSON value")
	}
	return nil
}

// readJSON streams a private record without imposing a limit absent from its writer.
func readJSON(path string, value any) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open record %q: %w", filepath.Base(path), err)
	}
	defer f.Close()
	if err := decodeJSON(f, value); err != nil {
		return fmt.Errorf("decode record %q: %w", filepath.Base(path), err)
	}
	return nil
}

// writeJSON atomically replaces a private durable JSON record.
func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("encode record: %w", err)
	}
	if err := writePrivate(path, append(data, '\n')); err != nil {
		return fmt.Errorf("write JSON record: %w", err)
	}
	return nil
}

// writePrivate durably and atomically replaces private snapshot or receipt bytes.
func writePrivate(path string, data []byte) error {
	f, err := os.CreateTemp(filepath.Dir(path), ".record-*")
	if err != nil {
		return fmt.Errorf("create record temporary: %w", err)
	}
	defer os.Remove(f.Name())
	if _, err := f.Write(data); err != nil {
		f.Close()
		return fmt.Errorf("write record temporary: %w", err)
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return fmt.Errorf("sync record temporary: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("close record temporary: %w", err)
	}
	if err := os.Rename(f.Name(), path); err != nil {
		return fmt.Errorf("replace record: %w", err)
	}
	d, err := os.Open(filepath.Dir(path))
	if err != nil {
		return fmt.Errorf("open record directory: %w", err)
	}
	defer d.Close()
	if err := d.Sync(); err != nil {
		return fmt.Errorf("sync record directory: %w", err)
	}
	return nil
}

// lockRun acquires the single-controller advisory lock for a private run directory.
func lockRun(directory string) (*os.File, error) {
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return nil, fmt.Errorf("create run directory: %w", err)
	}
	if err := os.Chmod(directory, 0o700); err != nil {
		return nil, fmt.Errorf("restrict run directory: %w", err)
	}
	f, err := os.OpenFile(filepath.Join(directory, ".lock"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open controller lock: %w", err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		f.Close()
		return nil, fmt.Errorf("run already controlled or lock unavailable: %w", err)
	}
	return f, nil
}

// Status reads the private journal as raw JSON for the CLI's status subcommand.
func Status(directory string) (json.RawMessage, error) {
	if !filepath.IsAbs(directory) {
		return nil, fmt.Errorf("run directory must be absolute")
	}
	data, err := os.ReadFile(filepath.Join(directory, "journal.json"))
	if err != nil {
		return nil, fmt.Errorf("read run status: %w", err)
	}
	if !json.Valid(data) {
		return nil, fmt.Errorf("run journal is not valid JSON")
	}
	return json.RawMessage(data), nil
}

// Stop requests a drain, or explicit cancellation of only this run's recorded IDs.
// An idle run applies the request when resumed; this function never contacts an endpoint.
func Stop(directory string, cancelOwned bool) error {
	if _, err := Status(directory); err != nil {
		return fmt.Errorf("verify stop target: %w", err)
	}
	if err := writeJSON(filepath.Join(directory, "control.json"), control{Stop: true, CancelOwned: cancelOwned}); err != nil {
		return fmt.Errorf("request run stop: %w", err)
	}
	return nil
}

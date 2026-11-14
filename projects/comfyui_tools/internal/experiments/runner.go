package experiments

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"sync"
	"time"
)

// completion communicates an asynchronous observation or artifact collection.
type completion struct {
	// index identifies the durable experiment, without sharing mutable records.
	index int
	// collected distinguishes artifact completion from native job completion.
	collected bool
	// media contains validated original downloads.
	media []media
	// warning records a successful envelope with wrapper cleanup failure.
	warning bool
	// err records a bounded operation failure for private diagnostics.
	err error
	// observationFailed means transport or observation ended without a confirmed native outcome.
	observationFailed bool
	// nativeStatus records a confirmed terminal native outcome independent of envelope success.
	nativeStatus string
}

// controller owns all mutable journal state; worker goroutines return messages only.
type controller struct {
	// directory owns private run artifacts.
	directory string
	// journal is the only mutable durable run record.
	journal journal
	// events receives asynchronous job and media observations.
	events chan completion
	// workers tracks observers and collectors currently running.
	workers int
	// cancelling prevents repeated cancellation requests in one controller session.
	cancelling bool
	// reports coalesces immutable report projections independently of queue work.
	reports chan []byte
	// reportErrors retains the first report failure without blocking the reporter.
	reportErrors chan error
	// reportDone closes after all requested report writes settle.
	reportDone chan struct{}
	// tasks joins observers and collectors before releasing run ownership.
	tasks sync.WaitGroup
}

// persist durably records a controller transition.
func (c *controller) persist() error {
	c.journal.Updated = time.Now().UTC()
	if err := writeJSON(filepath.Join(c.directory, "journal.json"), &c.journal); err != nil {
		return fmt.Errorf("persist run journal: %w", err)
	}
	return nil
}

// loadConfiguration normalizes the serialized contract and refuses implicit infinite runs.
func loadConfiguration(path string) (configuration, error) {
	var config configuration
	if err := readJSON(path, &config); err != nil {
		return config, fmt.Errorf("read experiment configuration: %w", err)
	}
	if len(config.CLI) == 0 || len(config.Producer) == 0 || config.CLI[0] == "" || config.Producer[0] == "" {
		return config, fmt.Errorf("cli and producer must be nonempty literal argv arrays")
	}
	if config.QueuedBuffer == 0 {
		config.QueuedBuffer = 2
	}
	if config.TimeoutSeconds == 0 {
		config.TimeoutSeconds = 3600
	}
	if config.DuplicateAttempts == 0 {
		config.DuplicateAttempts = 8
	}
	if config.QueuedBuffer < 1 || config.QueuedBuffer > 64 || config.TimeoutSeconds < 1 || config.DuplicateAttempts < 1 || config.DuplicateAttempts > 100 {
		return config, fmt.Errorf("queue buffer, timeout or duplicate bound is outside the supported range")
	}
	if (config.Continuous && config.MaxExperiments != 0) || (!config.Continuous && config.MaxExperiments <= 0) {
		return config, fmt.Errorf("choose continuous:true or a positive max_experiments")
	}
	return config, nil
}

// Run starts or resumes an authorized experiment stream using a private run directory.
// Context cancellation stops local children and preserves remote jobs for resumption.
func Run(ctx context.Context, configurationPath, directory string) (resultErr error) {
	if !filepath.IsAbs(directory) || !filepath.IsAbs(configurationPath) {
		return fmt.Errorf("configuration and run directory must be absolute paths")
	}
	config, err := loadConfiguration(configurationPath)
	if err != nil {
		return fmt.Errorf("configure experiments: %w", err)
	}
	lock, err := lockRun(directory)
	if err != nil {
		return fmt.Errorf("acquire experiment ownership: %w", err)
	}
	defer lock.Close()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	c := &controller{directory: directory, events: make(chan completion, 256), reports: make(chan []byte, 1), reportErrors: make(chan error, 1), reportDone: make(chan struct{})}
	defer func() { cancel(); c.tasks.Wait() }()
	if err := readJSON(filepath.Join(directory, "journal.json"), &c.journal); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("resume experiment journal: %w", err)
		}
		c.journal = journal{Version: 1, Config: config, Experiments: []*experiment{}}
	} else if c.journal.Version != 1 || !reflect.DeepEqual(c.journal.Config, config) {
		return fmt.Errorf("resume requires the same configuration and journal version")
	}
	go c.renderReports()
	defer func() {
		close(c.reports)
		<-c.reportDone
		select {
		case err := <-c.reportErrors:
			resultErr = errors.Join(resultErr, fmt.Errorf("finish local reports: %w", err))
		default:
		}
	}()
	if err := c.recover(ctx); err != nil {
		return fmt.Errorf("recover experiment stream: %w", err)
	}
	if err := c.persist(); err != nil {
		return fmt.Errorf("initialize experiment stream: %w", err)
	}
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	for {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("experiment controller interrupted; remote jobs retained: %w", err)
		}
		var request control
		if err := readJSON(filepath.Join(directory, "control.json"), &request); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("read run control: %w", err)
		}
		if request.CancelOwned && !c.cancelling {
			c.cancelling = true
			if err := c.cancelOwned(ctx); err != nil {
				return fmt.Errorf("cancel owned experiment jobs: %w", err)
			}
		}
		if !request.Stop && c.journal.Failure == "" && !c.journal.Exhausted {
			if err := c.refill(ctx); err != nil {
				if interrupted := ctx.Err(); interrupted != nil && c.journal.Failure != "submission_unknown" {
					return fmt.Errorf("experiment preparation interrupted; remote jobs retained: %w", interrupted)
				}
				if c.journal.Failure == "" {
					c.journal.Failure = "producer_or_submission"
				}
				if persistErr := c.persist(); persistErr != nil {
					return fmt.Errorf("retain refill failure: %w", errors.Join(err, persistErr))
				}
				return fmt.Errorf("refill experiments: %w", err)
			}
		}
		if c.workers == 0 {
			if err := c.report(); err != nil {
				return fmt.Errorf("publish local run report: %w", err)
			}
			if c.journal.Failure != "" {
				return fmt.Errorf("run stopped with %s; inspect private journal and receipts", c.journal.Failure)
			}
			if request.Stop || c.journal.Exhausted {
				return nil
			}
		}
		select {
		case err := <-c.reportErrors:
			return fmt.Errorf("write asynchronous report: %w", err)
		case <-ctx.Done():
			return fmt.Errorf("experiment controller interrupted; remote jobs retained: %w", ctx.Err())
		case event := <-c.events:
			if err := ctx.Err(); err != nil {
				return fmt.Errorf("experiment observation interrupted; remote jobs retained: %w", err)
			}
			if err := c.complete(ctx, event); err != nil {
				return fmt.Errorf("record experiment completion: %w", err)
			}
		case <-ticker.C:
		}
	}
}

// recover adopts recorded native IDs, including valid receipts saved before a crash.
func (c *controller) recover(ctx context.Context) error {
	for i, e := range c.journal.Experiments {
		if e.State == "submitting" {
			var response struct {
				// OK identifies a successful persisted submission response.
				OK bool `json:"ok"`
				// Data retains the native ownership identifier.
				Data struct {
					// PromptID identifies the accepted native job.
					PromptID string `json:"prompt_id"`
				} `json:"data"`
			}
			if err := readJSON(filepath.Join(c.directory, e.Directory, "submit.json"), &response); err == nil && response.OK && response.Data.PromptID != "" {
				e.PromptID, e.State = response.Data.PromptID, "queued"
			} else {
				e.State = "submission_unknown"
			}
		}
		if e.State == "submission_unknown" {
			c.journal.Failure = "submission_unknown"
		}
		if e.State == "observation_failed" {
			e.State, e.Failure, e.Finished = "queued", "", time.Time{}
		}
		if e.State == "queued" || e.State == "cancelling" {
			c.observe(ctx, i)
		} else if e.State == "collecting" {
			c.collect(ctx, i)
		}
	}
	if c.journal.Failure == "native_observation" {
		c.journal.Failure = ""
		for _, e := range c.journal.Experiments {
			if e.Failure != "" {
				c.journal.Failure = e.Failure
				break
			}
		}
	}
	return nil
}

// refill snapshots and submits distinct experiments until the execution lookahead is full.
func (c *controller) refill(ctx context.Context) error {
	config := c.journal.Config
	outstanding := 0
	for _, e := range c.journal.Experiments {
		if e.State == "queued" {
			outstanding++
		}
	}
	duplicates := 0
	for outstanding < 1+config.QueuedBuffer {
		// Recheck external control between potentially slow producer and CLI operations.
		var request control
		if err := readJSON(filepath.Join(c.directory, "control.json"), &request); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("read refill control: %w", err)
		}
		if request.Stop {
			return nil
		}
		if !config.Continuous && len(c.journal.Experiments) >= config.MaxExperiments {
			c.journal.Exhausted = true
			break
		}
		e, err := c.prepare(ctx)
		if err != nil {
			return fmt.Errorf("prepare next workflow: %w", err)
		}
		if c.journal.Exhausted {
			break
		}
		if e == nil {
			duplicates++
			if duplicates >= config.DuplicateAttempts {
				return fmt.Errorf("producer exhausted distinct identities after %d attempts", duplicates)
			}
			continue
		}
		duplicates = 0
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("workflow preparation interrupted before submission: %w", err)
		}
		var afterPreparation control
		if err := readJSON(filepath.Join(c.directory, "control.json"), &afterPreparation); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("read post-preparation control: %w", err)
		}
		if afterPreparation.Stop {
			return nil
		}
		c.journal.Experiments = append(c.journal.Experiments, e)
		if err := c.persist(); err != nil {
			return fmt.Errorf("record pending submission: %w", err)
		}
		var result struct {
			// PromptID is the native acknowledgment of ownership.
			PromptID string `json:"prompt_id"`
		}
		warning, err := callCLI(ctx, config, filepath.Join(c.directory, e.Directory, "submit"), &result, "run", "--workflow", filepath.Join(c.directory, e.Directory, "workflow.json"))
		if err != nil || result.PromptID == "" {
			e.State, e.Failure = "submission_unknown", "submission_acknowledgment"
			c.journal.Failure = "submission_unknown"
			if saveErr := c.persist(); saveErr != nil {
				return fmt.Errorf("record ambiguous submission: %w", saveErr)
			}
			return fmt.Errorf("submission acknowledgment is ambiguous; automatic retry refused")
		}
		e.PromptID, e.State, e.Submitted, e.CleanupWarning = result.PromptID, "queued", time.Now().UTC(), warning
		if err := c.persist(); err != nil {
			return fmt.Errorf("retain native ownership: %w", err)
		}
		c.observe(ctx, len(c.journal.Experiments)-1)
		outstanding++
	}
	if err := c.persist(); err != nil {
		return fmt.Errorf("retain buffer state: %w", err)
	}
	return nil
}

// prepare invokes the opaque producer and snapshots one novel workflow before submission.
func (c *controller) prepare(ctx context.Context) (*experiment, error) {
	sequence := c.journal.NextSequence
	c.journal.NextSequence++
	directory := filepath.Join(c.directory, fmt.Sprintf("%06d", sequence))
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return nil, fmt.Errorf("create experiment directory: %w", err)
	}
	if err := c.persist(); err != nil {
		return nil, fmt.Errorf("retain producer sequence: %w", err)
	}
	input, err := json.Marshal(map[string]any{"sequence": sequence, "directory": directory})
	if err != nil {
		return nil, fmt.Errorf("encode producer request: %w", err)
	}
	data, err := invoke(ctx, c.journal.Config.Producer, input, filepath.Join(directory, "producer"), c.journal.Config.TimeoutSeconds)
	if err != nil {
		return nil, fmt.Errorf("request workflow from producer: %w", err)
	}
	var result candidate
	if err := decodeJSON(bytes.NewReader(data), &result); err != nil {
		return nil, fmt.Errorf("decode producer response: %w", err)
	}
	if result.Done {
		c.journal.Exhausted = true
		return nil, nil
	}
	if !filepath.IsAbs(result.Workflow) || result.Description == "" {
		return nil, fmt.Errorf("producer must supply an absolute workflow and a variation description")
	}
	info, err := os.Stat(result.Workflow)
	if err != nil {
		return nil, fmt.Errorf("inspect producer workflow: %w", err)
	}
	if !info.Mode().IsRegular() || info.Size() > 64<<20 {
		return nil, fmt.Errorf("workflow must be a regular file of at most 64 MiB")
	}
	workflow, err := os.ReadFile(result.Workflow)
	if err != nil {
		return nil, fmt.Errorf("read producer workflow: %w", err)
	}
	identity, err := workflowIdentity(workflow, result.IdentityExclusions)
	if err != nil {
		return nil, fmt.Errorf("identify producer workflow: %w", err)
	}
	for _, prior := range c.journal.Experiments {
		if prior.Identity == identity {
			c.journal.Duplicates++
			return nil, nil
		}
	}
	if err := writePrivate(filepath.Join(directory, "workflow.json"), workflow); err != nil {
		return nil, fmt.Errorf("snapshot workflow: %w", err)
	}
	return &experiment{Sequence: sequence, Directory: filepath.Base(directory), Candidate: result, Identity: identity, WorkflowSHA256: digest(workflow), State: "submitting", Created: time.Now().UTC()}, nil
}

// observe watches one known native ID without blocking replenishment or media retrieval.
func (c *controller) observe(ctx context.Context, index int) {
	e := *c.journal.Experiments[index]
	c.workers++
	c.tasks.Add(1)
	go func() {
		defer c.tasks.Done()
		var result struct {
			// Jobs is the native completion summary for the requested identifier.
			Jobs []struct {
				// PromptID identifies the completed owned job.
				PromptID string `json:"prompt_id"`
				// Status identifies native execution success or failure.
				Status string `json:"status"`
				// OK confirms successful execution independent of process cleanup.
				OK bool `json:"ok"`
			} `json:"jobs"`
		}
		warning, err := callCLI(ctx, c.journal.Config, filepath.Join(c.directory, e.Directory, fmt.Sprintf("wait-%d", time.Now().UnixNano())), &result, "jobs", "wait", e.PromptID, "--timeout", strconv.Itoa(c.journal.Config.TimeoutSeconds), "--poll-interval", "2")
		failedEnvelope := false
		var failure *cliFailure
		if errors.As(err, &failure) && (failure.code == "execution_error" || failure.code == "cancelled" || failure.code == "wait_timeout") {
			failedEnvelope = decodeJSON(bytes.NewReader(failure.details), &result) == nil
		}
		observationFailed := err != nil
		nativeStatus := ""
		if err == nil || failedEnvelope {
			if len(result.Jobs) != 1 || result.Jobs[0].PromptID != e.PromptID {
				err, observationFailed = fmt.Errorf("observer did not identify the requested native job"), true
			} else {
				switch result.Jobs[0].Status {
				case "completed":
					if !result.Jobs[0].OK || failedEnvelope {
						err, observationFailed = fmt.Errorf("observer returned an inconsistent completion"), true
					} else {
						nativeStatus = "completed"
					}
				case "failed", "error", "cancelled":
					if result.Jobs[0].OK {
						err, observationFailed = fmt.Errorf("observer returned an inconsistent failure"), true
					} else {
						err, observationFailed = fmt.Errorf("owned native job ended unsuccessfully"), false
						nativeStatus = result.Jobs[0].Status
					}
				default:
					err, observationFailed = fmt.Errorf("observer ended without a terminal native outcome"), true
				}
			}
		}
		select {
		case c.events <- completion{index: index, warning: warning, err: err, observationFailed: observationFailed, nativeStatus: nativeStatus}:
		case <-ctx.Done():
		}
	}()
}

// complete persists native completion and schedules independent artifact collection.
func (c *controller) complete(ctx context.Context, event completion) error {
	e := c.journal.Experiments[event.index]
	// A reconciled native outcome already released the old observer's slot.
	if !event.collected && e.State != "queued" && e.State != "cancelling" {
		return nil
	}
	c.workers--
	e.CleanupWarning = e.CleanupWarning || event.warning
	if event.err != nil {
		requestedCancellation := e.State == "cancelling"
		e.State, e.Failure, e.Finished = "failed", "native_execution", time.Now().UTC()
		if event.collected {
			e.Failure = "artifact_collection"
		} else if event.observationFailed {
			e.State, e.Failure, e.Finished = "observation_failed", "native_observation", time.Time{}
		}
		if event.nativeStatus == "cancelled" && !event.observationFailed {
			e.State = "cancelled"
			if requestedCancellation {
				e.Failure = ""
			}
		}
		if e.Failure != "" {
			c.journal.Failure = e.Failure
		}
	} else if event.collected {
		e.State, e.Media, e.Finished = "completed", event.media, time.Now().UTC()
	} else {
		e.State = "collecting"
	}
	if err := c.persist(); err != nil {
		return fmt.Errorf("retain completion: %w", err)
	}
	// The next loop refills first. Collection runs concurrently and never owns journal state.
	if e.State == "collecting" {
		c.collect(ctx, event.index)
	}
	if event.collected {
		// Report writing is independent; original records remain in journal.json.
		if err := c.report(); err != nil {
			return fmt.Errorf("update local report: %w", err)
		}
	}
	return nil
}

// cancelOwned targets only this journal's outstanding IDs through the native CLI guard.
func (c *controller) cancelOwned(ctx context.Context) error {
	for i, e := range c.journal.Experiments {
		if (e.State != "queued" && e.State != "observation_failed" && e.State != "cancelling") || e.PromptID == "" {
			continue
		}
		hasObserver := e.State != "observation_failed"
		e.State = "cancelling"
		if err := c.persist(); err != nil {
			return fmt.Errorf("retain cancellation intent: %w", err)
		}
		var acknowledgment struct {
			// PromptID identifies the exact native cancellation target.
			PromptID string `json:"prompt_id"`
			// QueueDeleteOK reports the optional local pending-queue operation.
			QueueDeleteOK *bool `json:"queue_delete_ok"`
			// InterruptOK reports the optional local execution interrupt operation.
			InterruptOK *bool `json:"interrupt_ok"`
		}
		warning, err := callCLI(ctx, c.journal.Config, filepath.Join(c.directory, e.Directory, "cancel"), &acknowledgment, "jobs", "cancel", e.PromptID)
		if err != nil {
			return fmt.Errorf("cancel recorded experiment %d: %w", e.Sequence, err)
		}
		e.CleanupWarning = e.CleanupWarning || warning
		if acknowledgment.PromptID != e.PromptID {
			return fmt.Errorf("cancellation did not identify the recorded native job")
		}
		var outcome struct {
			// PromptID binds the status observation to the recorded native identifier.
			PromptID string `json:"prompt_id"`
			// Status distinguishes confirmed completion from a best-effort cancellation request.
			Status string `json:"status"`
		}
		warning, err = callCLI(ctx, c.journal.Config, filepath.Join(c.directory, e.Directory, fmt.Sprintf("cancel-status-%d", time.Now().UnixNano())), &outcome, "jobs", "status", e.PromptID)
		if err != nil {
			return fmt.Errorf("observe cancellation outcome for experiment %d: %w", e.Sequence, err)
		}
		e.CleanupWarning = e.CleanupWarning || warning
		if outcome.PromptID != e.PromptID {
			return fmt.Errorf("cancellation outcome did not identify the recorded native job")
		}
		partialFailure := (acknowledgment.QueueDeleteOK != nil && !*acknowledgment.QueueDeleteOK) || (acknowledgment.InterruptOK != nil && !*acknowledgment.InterruptOK)
		if partialFailure && outcome.Status != "completed" && outcome.Status != "failed" && outcome.Status != "error" {
			return fmt.Errorf("native cancellation was not confirmed; live ownership retained")
		}
		switch outcome.Status {
		case "completed":
			e.State, e.Failure = "collecting", ""
		case "cancelled":
			e.State, e.Failure, e.Finished = "cancelled", "", time.Now().UTC()
		case "failed", "error":
			e.State, e.Failure, e.Finished = "failed", "native_execution", time.Now().UTC()
			c.journal.Failure = e.Failure
		case "pending", "running", "queued", "cancelling":
			// The observer retains this still-live job until a terminal outcome.
		default:
			return fmt.Errorf("native cancellation outcome is not observable; ownership retained")
		}
		if e.State != "cancelling" && hasObserver {
			c.workers--
		}
		if err := c.persist(); err != nil {
			return fmt.Errorf("retain reconciled cancellation: %w", err)
		}
		if e.State == "collecting" {
			c.collect(ctx, i)
		} else if e.State == "cancelling" && !hasObserver {
			c.observe(ctx, i)
		}
	}
	if err := c.persist(); err != nil {
		return fmt.Errorf("retain cancellation requests: %w", err)
	}
	return nil
}

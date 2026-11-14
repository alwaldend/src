package experiments

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// nativeJob records a simulated native identifier and serial GPU completion time.
type nativeJob struct {
	// ID is the endpoint-owned identifier returned by submission.
	ID string
	// Finish is the simulated serial GPU completion time in Unix nanoseconds.
	Finish int64
	// Cancelled records an explicit identifier-scoped cancellation.
	Cancelled bool
}

// nativeState is shared by fixture subprocesses under a file lock.
type nativeState struct {
	// Mode selects a specified failure case instead of normal operation.
	Mode string
	// Jobs preserves every accepted submission and its identity.
	Jobs []nativeJob
	// Late counts refills after all simulated execution has already finished.
	Late int
	// Cancelled records the exact IDs supplied to cancellation.
	Cancelled []string
	// Foreign remains present throughout owned cancellation acceptance.
	Foreign bool
	// ObserverFailed records a one-shot transport failure independent of native execution.
	ObserverFailed bool
	// DownloadFailed records a one-shot retrieval transport failure.
	DownloadFailed bool
}

// fixture owns command configuration and a private simulated endpoint.
type fixture struct {
	// t handles command-level acceptance failures.
	t *testing.T
	// directory owns the repeatable test inputs and private run data.
	directory string
	// state names the simulated endpoint journal.
	state string
	// config names the runner's execution contract.
	config string
	// run names the runner's private artifact directory.
	run string
	// binary identifies the Bazel-built experiment command under test.
	binary string
}

// writeFixtureJSON writes repeatable fixture inputs and fails at the owning test boundary.
func writeFixtureJSON(t *testing.T, path string, value any) {
	t.Helper()
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

// newFixture prepares a command-level producer and CLI without network or credentials.
func newFixture(t *testing.T, mode string, count int) *fixture {
	t.Helper()
	directory := t.TempDir()
	binary, err := runfiles.Rlocation(os.Getenv("EXPERIMENTS_BINARY"))
	if err != nil {
		t.Fatal(err)
	}
	helper, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	f := &fixture{t: t, directory: directory, state: filepath.Join(directory, "endpoint.json"), config: filepath.Join(directory, "config.json"), run: filepath.Join(directory, "run"), binary: binary}
	t.Cleanup(func() {
		if t.Failed() {
			if data, err := os.ReadFile(filepath.Join(directory, "command.log")); err == nil {
				t.Logf("fixture command diagnostic: %s", data)
			}
			if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
				for _, name := range []string{"command.log", "run/journal.json"} {
					if data, err := os.ReadFile(filepath.Join(directory, name)); err == nil {
						os.WriteFile(filepath.Join(output, strings.ReplaceAll(t.Name()+"-"+name, "/", "-")), data, 0o600)
					}
				}
			}
		}
	})
	writeFixtureJSON(t, f.state, nativeState{Mode: mode, Foreign: true})
	prefix := []string{helper, "-test.run=^TestNativeHelper$", "--"}
	writeFixtureJSON(t, f.config, map[string]any{
		"cli":           append(append([]string{}, prefix...), "comfy", f.state),
		"producer":      append(append([]string{}, prefix...), "producer", mode),
		"queued_buffer": 2, "max_experiments": count, "timeout_seconds": 15, "duplicate_attempts": 3,
	})
	return f
}

// command prepares the actual binary under test with bounded process execution.
func (f *fixture) command(args ...string) *exec.Cmd {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	f.t.Cleanup(cancel)
	cmd := exec.CommandContext(ctx, f.binary, args...)
	cmd.Env = append(os.Environ(), "COMFY_EXPERIMENT_HELPER=1")
	return cmd
}

// start launches a real runner process and retains its diagnostic artifact.
func (f *fixture) start() *exec.Cmd {
	log, err := os.OpenFile(filepath.Join(f.directory, "command.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		f.t.Fatal(err)
	}
	f.t.Cleanup(func() { log.Close() })
	cmd := f.command("run", "--config", f.config, "--directory", f.run)
	cmd.Stdout, cmd.Stderr = log, log
	if err := cmd.Start(); err != nil {
		f.t.Fatal(err)
	}
	f.t.Cleanup(func() {
		if cmd.ProcessState == nil {
			cmd.Process.Kill()
			cmd.Wait()
		}
	})
	return cmd
}

// stateValue reads the shared simulated endpoint while respecting its subprocess lock.
func (f *fixture) stateValue() nativeState {
	var result nativeState
	if err := withState(f.state, func(s *nativeState) error { result = *s; return nil }); err != nil {
		f.t.Fatal(err)
	}
	return result
}

// awaitSubmissions waits for observable asynchronous progress rather than arbitrary sleeps.
func (f *fixture) awaitSubmissions(count int) {
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if len(f.stateValue().Jobs) >= count {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	f.t.Fatalf("did not observe %d submissions", count)
}

// journalValue reads the public serialized journal contract for E2E assertions.
func (f *fixture) journalValue() map[string]any {
	data, err := os.ReadFile(filepath.Join(f.run, "journal.json"))
	if err != nil {
		f.t.Fatal(err)
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		f.t.Fatal(err)
	}
	return result
}

// receipt retains externally verifiable acceptance evidence in Bazel test outputs.
func (f *fixture) receipt(name string) {
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		writeFixtureJSON(f.t, filepath.Join(directory, name+".json"), map[string]any{"endpoint": f.stateValue(), "journal": f.journalValue()})
	}
}

// TestBufferedExecution verifies slow downloads, native cleanup errors and original media.
func TestBufferedExecution(t *testing.T) {
	f := newFixture(t, "cleanup-warning", 8)
	cmd := f.start()
	if err := cmd.Wait(); err != nil {
		t.Fatal(err)
	}
	s := f.stateValue()
	if len(s.Jobs) != 8 || s.Late != 0 {
		t.Fatalf("buffer lost work: jobs=%d late=%d", len(s.Jobs), s.Late)
	}
	j := f.journalValue()
	for _, raw := range j["experiments"].([]any) {
		e := raw.(map[string]any)
		if e["state"] != "completed" || e["cleanup_warning"] != true {
			t.Fatalf("incorrect completion: %v", e["state"])
		}
		if len(e["media"].([]any)) != 1 {
			t.Fatal("original media missing")
		}
		file := e["media"].([]any)[0].(map[string]any)
		path := filepath.Join(f.run, file["path"].(string))
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		hash := sha256.Sum256(data)
		if hex.EncodeToString(hash[:]) != file["sha256"] {
			t.Fatal("original digest does not match receipt")
		}
		original, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		canvas, err := png.Decode(original)
		original.Close()
		if err != nil || canvas.Bounds() != image.Rect(0, 0, 32, 32) {
			t.Fatal("original image content was changed")
		}
	}
	page, err := os.ReadFile(filepath.Join(f.run, "gallery.html"))
	if err != nil || !strings.Contains(string(page), "<details>") || !strings.Contains(string(page), "positive_prompt") {
		t.Fatal("collapsed report missing")
	}
	f.receipt("buffered-execution")
}

// TestDuplicateIdentity verifies seed and output-prefix changes cannot fill the queue.
func TestDuplicateIdentity(t *testing.T) {
	f := newFixture(t, "duplicate", 4)
	cmd := f.start()
	if err := cmd.Wait(); err == nil {
		t.Fatal("duplicate-only producer unexpectedly succeeded")
	}
	if len(f.stateValue().Jobs) != 1 {
		t.Fatal("seed-only duplicate was submitted")
	}
	if f.journalValue()["duplicates"].(float64) != 3 {
		t.Fatal("duplicate retry bound not observed")
	}
	f.receipt("duplicate-identity")
}

// TestResumeAndLock verifies interruption adopts native IDs and rejects a second controller.
func TestResumeAndLock(t *testing.T) {
	f := newFixture(t, "resume", 7)
	first := f.start()
	f.awaitSubmissions(3)
	deadline := time.Now().Add(5 * time.Second)
	for {
		recorded := f.journalValue()["experiments"].([]any)
		if len(recorded) == 3 && recorded[2].(map[string]any)["prompt_id"] != nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("native acknowledgments were not durably recorded")
		}
		time.Sleep(5 * time.Millisecond)
	}
	other := f.command("run", "--config", f.config, "--directory", f.run)
	if err := other.Run(); err == nil {
		t.Fatal("second controller acquired the same run")
	}
	if err := first.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	if err := first.Wait(); err == nil {
		t.Fatal("interrupted controller reported completion")
	}
	if err := withState(f.state, func(s *nativeState) error {
		s.Mode = "normal"
		for i := range s.Jobs {
			s.Jobs[i].Finish = time.Now().Add(time.Duration(i+1) * 150 * time.Millisecond).UnixNano()
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	second := f.start()
	if err := second.Wait(); err != nil {
		t.Fatal(err)
	}
	if len(f.stateValue().Jobs) != 7 {
		t.Fatal("resume duplicated native jobs")
	}
	f.receipt("resume-and-lock")
}

// TestInterruptDuringPreparation verifies local producer cancellation does not poison resume.
func TestInterruptDuringPreparation(t *testing.T) {
	f := newFixture(t, "slow-producer", 4)
	first := f.start()
	marker := filepath.Join(f.run, "000003", "producer-started")
	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, err := os.Stat(marker); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("producer did not reach the interruption boundary")
		}
		time.Sleep(5 * time.Millisecond)
	}
	if err := first.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	if err := first.Wait(); err == nil {
		t.Fatal("interrupted preparation reported success")
	}
	if f.journalValue()["failure"] != nil {
		t.Fatal("local interruption became a permanent remote failure")
	}
	if err := os.WriteFile(filepath.Join(f.run, "producer-release"), []byte("resume"), 0o600); err != nil {
		t.Fatal(err)
	}
	second := f.start()
	if err := second.Wait(); err != nil {
		t.Fatal(err)
	}
	if len(f.stateValue().Jobs) != 4 {
		t.Fatal("preparation interruption duplicated or lost accepted work")
	}
	f.receipt("preparation-interruption")
}

// TestObserverFailureRetainsOwnership verifies recovery and cancellation after a transport failure.
func TestObserverFailureRetainsOwnership(t *testing.T) {
	for _, cancelOwned := range []bool{false, true} {
		t.Run(strconv.FormatBool(cancelOwned), func(t *testing.T) {
			mode := "observer-disconnect"
			if cancelOwned {
				mode += "-long"
			}
			f := newFixture(t, mode, 1)
			first := f.start()
			if err := first.Wait(); err == nil {
				t.Fatal("observer transport failure was not reported")
			}
			initial := f.journalValue()["experiments"].([]any)[0].(map[string]any)
			if initial["state"] != "observation_failed" || initial["prompt_id"] == "" {
				t.Fatal("observer failure discarded live ownership")
			}
			if cancelOwned {
				if err := f.command("stop", "--directory", f.run, "--cancel-owned").Run(); err != nil {
					t.Fatal(err)
				}
			}
			second := f.start()
			if err := second.Wait(); err != nil {
				t.Fatal(err)
			}
			final := f.journalValue()["experiments"].([]any)[0].(map[string]any)
			expected := "completed"
			if cancelOwned {
				expected = "cancelled"
			}
			if final["state"] != expected || final["prompt_id"] != initial["prompt_id"] || len(f.stateValue().Jobs) != 1 {
				t.Fatal("owned job was lost or resubmitted")
			}
			if cancelOwned && len(f.stateValue().Cancelled) != 1 {
				t.Fatal("owned cancellation did not reach the disconnected job")
			}
			f.receipt("observer-recovery-" + strconv.FormatBool(cancelOwned))
		})
	}
}

// TestNativeTerminalEnvelope distinguishes native error summaries from lost observation.
func TestNativeTerminalEnvelope(t *testing.T) {
	for _, mode := range []string{"job-failure", "external-cancel", "wait-timeout"} {
		t.Run(mode, func(t *testing.T) {
			f := newFixture(t, mode, 1)
			if err := f.start().Wait(); err == nil {
				t.Fatal("unsuccessful native wait reported success")
			}
			initial := f.journalValue()["experiments"].([]any)[0].(map[string]any)
			if mode == "wait-timeout" {
				if initial["state"] != "observation_failed" {
					t.Fatal("timeout was mistaken for terminal execution")
				}
				if err := f.start().Wait(); err != nil {
					t.Fatal(err)
				}
				if f.journalValue()["experiments"].([]any)[0].(map[string]any)["state"] != "completed" {
					t.Fatal("timed-out observation did not recover its original job")
				}
			} else {
				expected := "failed"
				if mode == "external-cancel" {
					expected = "cancelled"
				}
				if initial["state"] != expected || f.journalValue()["failure"] != "native_execution" {
					t.Fatal("confirmed terminal outcome was mistaken for observation failure")
				}
				if err := f.start().Wait(); err == nil {
					t.Fatal("restart cleared a confirmed native failure")
				}
			}
			if len(f.stateValue().Jobs) != 1 {
				t.Fatal("native outcome handling resubmitted generation")
			}
			f.receipt(mode + "-envelope")
		})
	}
}

// TestExternalCancellationRecovery preserves a terminal failure after a later observer disconnect.
func TestExternalCancellationRecovery(t *testing.T) {
	f := newFixture(t, "cancel-then-disconnect", 4)
	if err := f.start().Wait(); err == nil {
		t.Fatal("mixed native and observation failure reported success")
	}
	initial := f.journalValue()
	if initial["failure"] != "native_observation" || len(f.stateValue().Jobs) != 3 {
		t.Fatal("fixture did not reach cancellation followed by observer failure")
	}
	if err := f.start().Wait(); err == nil {
		t.Fatal("restart forgot the confirmed external cancellation")
	}
	if len(f.stateValue().Jobs) != 3 || f.journalValue()["failure"] != "native_execution" {
		t.Fatal("recovery replenished work despite a retained terminal failure")
	}
	f.receipt("external-cancellation-recovery")
}

// TestCancellationReconciliation retains execution and outputs after best-effort cancellation.
func TestCancellationReconciliation(t *testing.T) {
	for _, mode := range []string{"cancel-rejected", "cancel-rejected-cached", "cancel-completed"} {
		t.Run(mode, func(t *testing.T) {
			f := newFixture(t, mode, 1)
			cmd := f.start()
			f.awaitSubmissions(1)
			if err := f.command("stop", "--directory", f.run, "--cancel-owned").Run(); err != nil {
				t.Fatal(err)
			}
			err := cmd.Wait()
			if strings.HasPrefix(mode, "cancel-rejected") {
				if err == nil {
					t.Fatal("rejected interrupt reported confirmed cancellation")
				}
				if f.journalValue()["experiments"].([]any)[0].(map[string]any)["state"] == "cancelled" {
					t.Fatal("live job ownership was released")
				}
				if err := withState(f.state, func(s *nativeState) error {
					s.Jobs[0].Finish = time.Now().Add(-time.Second).UnixNano()
					return nil
				}); err != nil {
					t.Fatal(err)
				}
				if err := f.start().Wait(); err != nil {
					t.Fatal(err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			final := f.journalValue()["experiments"].([]any)[0].(map[string]any)
			if final["state"] != "completed" || len(final["media"].([]any)) != 1 || len(f.stateValue().Jobs) != 1 {
				t.Fatal("cancellation discarded completed original media or resubmitted work")
			}
			if !f.stateValue().Foreign {
				t.Fatal("foreign work was affected")
			}
			f.receipt(mode + "-reconciliation")
		})
	}
}

// TestLargeJournalResume verifies valid persisted records remain readable beyond 64 MiB.
func TestLargeJournalResume(t *testing.T) {
	f := newFixture(t, "normal", 3)
	if err := f.start().Wait(); err != nil {
		t.Fatal(err)
	}
	journal := f.journalValue()
	payload := strings.Repeat("x", 22<<20)
	for _, raw := range journal["experiments"].([]any) {
		raw.(map[string]any)["candidate"].(map[string]any)["notes"].(map[string]any)["large_payload"] = payload
	}
	path := filepath.Join(f.run, "journal.json")
	writeFixtureJSON(t, path, journal)
	info, err := os.Stat(path)
	if err != nil || info.Size() <= 64<<20 {
		t.Fatal("large persisted fixture did not cross the former read limit")
	}
	if err := f.start().Wait(); err != nil {
		t.Fatal(err)
	}
	if len(f.stateValue().Jobs) != 3 {
		t.Fatal("large journal resume resubmitted native jobs")
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		writeFixtureJSON(t, filepath.Join(output, "large-journal-resume.json"), map[string]any{"journal_bytes": info.Size(), "native_jobs": 3, "resumed": true})
	}
}

// TestStopDuringPreparation verifies a completed stop prevents a not-yet-submitted candidate.
func TestStopDuringPreparation(t *testing.T) {
	f := newFixture(t, "slow-producer", 4)
	cmd := f.start()
	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, err := os.Stat(filepath.Join(f.run, "000003", "producer-started")); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("slow producer boundary not reached")
		}
		time.Sleep(5 * time.Millisecond)
	}
	if err := f.command("stop", "--directory", f.run).Run(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(f.run, "producer-release"), []byte("release"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		t.Fatal(err)
	}
	if len(f.stateValue().Jobs) != 3 {
		t.Fatal("producer submitted another job after the completed stop request")
	}
	f.receipt("stop-during-preparation")
}

// TestTransientRetrieval verifies a transport retry never resubmits native generation.
func TestTransientRetrieval(t *testing.T) {
	f := newFixture(t, "download-disconnect", 1)
	if err := f.start().Wait(); err != nil {
		t.Fatal(err)
	}
	if len(f.stateValue().Jobs) != 1 || !f.stateValue().DownloadFailed {
		t.Fatal("retrieval failure case was not exercised exactly once")
	}
	if f.journalValue()["experiments"].([]any)[0].(map[string]any)["state"] != "completed" {
		t.Fatal("retrieval retry did not retain the original result")
	}
	f.receipt("transient-retrieval")
}

// TestExactNumericNotes verifies a full-width seed survives execution, resume and reporting.
func TestExactNumericNotes(t *testing.T) {
	f := newFixture(t, "exact-numbers", 1)
	for attempt := 0; attempt < 2; attempt++ {
		if err := f.start().Wait(); err != nil {
			t.Fatal(err)
		}
		for _, name := range []string{"journal.json", "gallery.html"} {
			data, err := os.ReadFile(filepath.Join(f.run, name))
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(data), "18446744073709551615") {
				t.Fatalf("full-width integer seed changed in %s", name)
			}
		}
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		writeFixtureJSON(t, filepath.Join(output, "exact-numeric-notes.json"), map[string]any{"seed": "18446744073709551615", "resumed": true})
	}
}

// TestAmbiguousAcknowledgment verifies accepted but unacknowledged work is never retried.
func TestAmbiguousAcknowledgment(t *testing.T) {
	f := newFixture(t, "lost-ack", 4)
	for i := 0; i < 2; i++ {
		cmd := f.start()
		if err := cmd.Wait(); err == nil {
			t.Fatal("ambiguous submission reported success")
		}
	}
	if len(f.stateValue().Jobs) != 1 {
		t.Fatal("ambiguous submission was repeated")
	}
	if f.journalValue()["failure"] != "submission_unknown" {
		t.Fatal("ambiguity was not preserved")
	}
	f.receipt("ambiguous-acknowledgment")
}

// TestOwnedCancellation verifies stop targets only recorded IDs and leaves foreign work present.
func TestOwnedCancellation(t *testing.T) {
	f := newFixture(t, "long", 10)
	cmd := f.start()
	f.awaitSubmissions(3)
	stop := f.command("stop", "--directory", f.run, "--cancel-owned")
	if err := stop.Run(); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		t.Fatal(err)
	}
	s := f.stateValue()
	if len(s.Jobs) != 3 || len(s.Cancelled) != 3 || !s.Foreign {
		t.Fatal("ownership or stop boundary failed")
	}
	for _, id := range s.Cancelled {
		if !strings.HasPrefix(id, "owned-") {
			t.Fatal("foreign cancellation")
		}
	}
	f.receipt("owned-cancellation")
}

// TestFailedArtifacts verifies failed execution and escaping output cannot become successes.
func TestFailedArtifacts(t *testing.T) {
	for _, mode := range []string{"job-failure", "escaping-media", "invalid-pointer", "producer-done", "trailing-json"} {
		t.Run(mode, func(t *testing.T) {
			f := newFixture(t, mode, 1)
			cmd := f.start()
			err := cmd.Wait()
			if mode == "producer-done" {
				if err != nil || len(f.stateValue().Jobs) != 0 {
					t.Fatal("producer exhaustion was not graceful")
				}
			} else if err == nil {
				t.Fatal("failed experiment reported success")
			}
			for _, raw := range f.journalValue()["experiments"].([]any) {
				if raw.(map[string]any)["state"] == "completed" {
					t.Fatal("invalid output was published")
				}
			}
			f.receipt(mode)
		})
	}
}

// withState serializes a fixture operation across real CLI subprocesses.
func withState(path string, operation func(*nativeState) error) error {
	lock, err := os.OpenFile(path+".lock", os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return fmt.Errorf("open fixture lock: %w", err)
	}
	defer lock.Close()
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX); err != nil {
		return fmt.Errorf("lock fixture: %w", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read fixture state: %w", err)
	}
	var state nativeState
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("decode fixture state: %w", err)
	}
	if err := operation(&state); err != nil {
		return fmt.Errorf("apply fixture operation: %w", err)
	}
	data, err = json.Marshal(state)
	if err != nil {
		return fmt.Errorf("encode fixture state: %w", err)
	}
	temporary, err := os.CreateTemp(filepath.Dir(path), ".endpoint-*")
	if err != nil {
		return fmt.Errorf("create fixture state temporary: %w", err)
	}
	defer os.Remove(temporary.Name())
	if _, err := temporary.Write(data); err != nil {
		temporary.Close()
		return fmt.Errorf("write fixture state: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close fixture state: %w", err)
	}
	if err := os.Rename(temporary.Name(), path); err != nil {
		return fmt.Errorf("replace fixture state: %w", err)
	}
	return nil
}

// emit writes only the native envelope; diagnostics and wrapper failures remain separate.
func emit(data any) error {
	if err := json.NewEncoder(os.Stdout).Encode(map[string]any{"ok": true, "data": data}); err != nil {
		return fmt.Errorf("write fixture envelope: %w", err)
	}
	return nil
}

// emitFailure reproduces the pinned CLI's structured unsuccessful wait envelope.
func emitFailure(code string, details any) error {
	if err := json.NewEncoder(os.Stdout).Encode(map[string]any{"ok": false, "data": nil, "error": map[string]any{"code": code, "details": details}}); err != nil {
		return fmt.Errorf("write fixture error envelope: %w", err)
	}
	return fmt.Errorf("simulated native unsuccessful outcome")
}

// producerHelper supplies neutral arbitrary API workflows through the documented contract.
func producerHelper(mode string) error {
	var request struct {
		// Sequence identifies a unique requested variation.
		Sequence int `json:"sequence"`
		// Directory is the private caller-owned candidate directory.
		Directory string `json:"directory"`
	}
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		return fmt.Errorf("read producer request: %w", err)
	}
	if mode == "slow-producer" && request.Sequence >= 3 {
		if _, err := os.Stat(filepath.Join(filepath.Dir(request.Directory), "producer-release")); os.IsNotExist(err) {
			if err := os.WriteFile(filepath.Join(request.Directory, "producer-started"), []byte("ready"), 0o600); err != nil {
				return fmt.Errorf("record producer boundary: %w", err)
			}
			for {
				if _, err := os.Stat(filepath.Join(filepath.Dir(request.Directory), "producer-release")); err == nil {
					break
				}
				time.Sleep(5 * time.Millisecond)
			}
		}
	}
	if mode == "producer-done" {
		if err := json.NewEncoder(os.Stdout).Encode(map[string]any{"done": true}); err != nil {
			return fmt.Errorf("write producer exhaustion: %w", err)
		}
		return nil
	}
	subject := fmt.Sprintf("geometric still life %d", request.Sequence)
	if mode == "duplicate" {
		subject = "geometric still life"
	}
	workflow := filepath.Join(request.Directory, "candidate.json")
	data, err := json.Marshal(map[string]any{"1": map[string]any{"class_type": "ArbitraryNode", "inputs": map[string]any{"subject": subject, "seed": request.Sequence, "prefix": strconv.Itoa(request.Sequence)}}})
	if err != nil {
		return fmt.Errorf("encode producer workflow: %w", err)
	}
	if err := os.WriteFile(workflow, data, 0o600); err != nil {
		return fmt.Errorf("write producer workflow: %w", err)
	}
	exclusions := []string{"/1/inputs/seed", "/1/inputs/prefix"}
	if mode == "invalid-pointer" {
		exclusions = []string{"/missing"}
	}
	notes := map[string]any{"positive_prompt": subject, "negative_prompt": "blur"}
	if mode == "exact-numbers" {
		notes["seed"] = json.Number("18446744073709551615")
	}
	if err := json.NewEncoder(os.Stdout).Encode(map[string]any{"workflow": workflow, "description": subject, "identity_exclusions": exclusions, "notes": notes}); err != nil {
		return fmt.Errorf("write producer result: %w", err)
	}
	if mode == "trailing-json" {
		if _, err := fmt.Fprintln(os.Stdout, "{}"); err != nil {
			return fmt.Errorf("write extra fixture JSON value: %w", err)
		}
	}
	return nil
}

// comfyHelper implements only owned submission, observation, download and cancellation.
func comfyHelper(path string, args []string) error {
	if len(args) == 0 || args[0] != "--json" {
		return fmt.Errorf("native JSON mode required")
	}
	args = args[1:]
	var state nativeState
	if err := withState(path, func(s *nativeState) error { state = *s; return nil }); err != nil {
		return fmt.Errorf("inspect fixture mode: %w", err)
	}
	switch args[0] {
	case "run":
		id := ""
		if err := withState(path, func(s *nativeState) error {
			now := time.Now().UnixNano()
			start := now
			if len(s.Jobs) > 0 {
				last := s.Jobs[len(s.Jobs)-1].Finish
				if last < now {
					s.Late++
				} else {
					start = last
				}
			}
			duration := 150 * time.Millisecond
			if s.Mode == "long" || s.Mode == "resume" || s.Mode == "observer-disconnect-long" || strings.HasPrefix(s.Mode, "cancel-rejected") || s.Mode == "cancel-completed" {
				duration = 10 * time.Second
			}
			id = fmt.Sprintf("owned-%d", len(s.Jobs))
			s.Jobs = append(s.Jobs, nativeJob{ID: id, Finish: start + int64(duration)})
			return nil
		}); err != nil {
			return fmt.Errorf("accept fixture job: %w", err)
		}
		if state.Mode == "lost-ack" {
			fmt.Fprint(os.Stdout, "acknowledgment lost")
			return fmt.Errorf("acknowledgment lost")
		}
		if err := emit(map[string]any{"prompt_id": id}); err != nil {
			return fmt.Errorf("acknowledge fixture job: %w", err)
		}
	case "jobs":
		id := args[2]
		if args[1] == "cancel" {
			if err := withState(path, func(s *nativeState) error {
				for i := range s.Jobs {
					if s.Jobs[i].ID == id {
						if s.Mode == "cancel-completed" {
							s.Jobs[i].Finish = time.Now().Add(-time.Second).UnixNano()
							return nil
						}
						if strings.HasPrefix(s.Mode, "cancel-rejected") || s.Jobs[i].Finish <= time.Now().UnixNano() {
							return nil
						}
						s.Jobs[i].Cancelled = true
						s.Cancelled = append(s.Cancelled, id)
						return nil
					}
				}
				return fmt.Errorf("unknown or foreign job")
			}); err != nil {
				return fmt.Errorf("cancel fixture job: %w", err)
			}
			if err := emit(map[string]any{"prompt_id": id, "queue_delete_ok": true, "interrupt_ok": !strings.HasPrefix(state.Mode, "cancel-rejected")}); err != nil {
				return fmt.Errorf("acknowledge cancellation: %w", err)
			}
		} else if args[1] == "status" {
			status := ""
			if err := withState(path, func(s *nativeState) error {
				for _, job := range s.Jobs {
					if job.ID == id {
						status = "running"
						if job.Cancelled {
							status = "cancelled"
						} else if time.Now().UnixNano() >= job.Finish {
							status = "completed"
						} else if s.Mode == "cancel-rejected-cached" {
							status = "cancelled"
						}
						return nil
					}
				}
				return fmt.Errorf("unknown status identifier")
			}); err != nil {
				return fmt.Errorf("read cancellation outcome: %w", err)
			}
			if err := emit(map[string]any{"prompt_id": id, "status": status}); err != nil {
				return fmt.Errorf("emit cancellation outcome: %w", err)
			}
		} else if args[1] == "wait" {
			if state.Mode == "wait-timeout" && !state.ObserverFailed {
				if err := withState(path, func(s *nativeState) error { s.ObserverFailed = true; return nil }); err != nil {
					return fmt.Errorf("record native wait timeout: %w", err)
				}
				return emitFailure("wait_timeout", map[string]any{"jobs": []any{map[string]any{"prompt_id": id, "status": "timed_out", "ok": false}}})
			}
			if strings.HasPrefix(state.Mode, "observer-disconnect") {
				fail := false
				if err := withState(path, func(s *nativeState) error {
					if !s.ObserverFailed {
						s.ObserverFailed = true
						fail = true
					}
					return nil
				}); err != nil {
					return fmt.Errorf("inject observer disconnection: %w", err)
				}
				if fail {
					return fmt.Errorf("simulated observer disconnection before terminal native state")
				}
			}
			for {
				var job nativeJob
				if err := withState(path, func(s *nativeState) error {
					for _, j := range s.Jobs {
						if j.ID == id {
							job = j
							return nil
						}
					}
					return fmt.Errorf("unknown job")
				}); err != nil {
					return fmt.Errorf("observe fixture job: %w", err)
				}
				if job.Cancelled || time.Now().UnixNano() >= job.Finish {
					if state.Mode == "cancel-then-disconnect" && id == "owned-1" && !state.ObserverFailed {
						deadline := time.Now().Add(5 * time.Second)
						for {
							data, readErr := os.ReadFile(filepath.Join(filepath.Dir(path), "run", "journal.json"))
							var record map[string]any
							if readErr == nil && json.Unmarshal(data, &record) == nil {
								items, ok := record["experiments"].([]any)
								if ok && len(items) > 0 && items[0].(map[string]any)["state"] == "cancelled" {
									break
								}
							}
							if time.Now().After(deadline) {
								return fmt.Errorf("cancelled outcome was not durably observed")
							}
							time.Sleep(5 * time.Millisecond)
						}
						if err := withState(path, func(s *nativeState) error { s.ObserverFailed = true; return nil }); err != nil {
							return fmt.Errorf("record disconnect after native cancellation: %w", err)
						}
						return fmt.Errorf("simulated observer disconnect after external cancellation")
					}
					// A different CLI profile may never observe a removed pending job.
					// Acknowledged cancellation must not depend on that observer finishing.
					if job.Cancelled && state.Mode == "long" {
						time.Sleep(time.Hour)
					}
					status := "completed"
					if job.Cancelled {
						status = "cancelled"
					} else if state.Mode == "job-failure" {
						status = "error"
					} else if state.Mode == "external-cancel" || (state.Mode == "cancel-then-disconnect" && id == "owned-0") {
						status = "cancelled"
					}
					summary := map[string]any{"jobs": []any{map[string]any{"prompt_id": id, "status": status, "ok": status == "completed"}}}
					if status == "error" {
						return emitFailure("execution_error", summary)
					}
					if status == "cancelled" {
						return emitFailure("cancelled", summary)
					}
					if err := emit(summary); err != nil {
						return fmt.Errorf("complete fixture observation: %w", err)
					}
					break
				}
				time.Sleep(5 * time.Millisecond)
			}
		} else {
			return fmt.Errorf("unsupported fixture jobs operation")
		}
	case "download":
		if state.Mode == "download-disconnect" {
			fail := false
			if err := withState(path, func(s *nativeState) error {
				if !s.DownloadFailed {
					s.DownloadFailed = true
					fail = true
				}
				return nil
			}); err != nil {
				return fmt.Errorf("inject retrieval disconnection: %w", err)
			}
			if fail {
				return fmt.Errorf("simulated transient retrieval transport failure")
			}
		}
		time.Sleep(700 * time.Millisecond)
		directory := args[3]
		if err := os.MkdirAll(directory, 0o700); err != nil {
			return fmt.Errorf("create fixture media: %w", err)
		}
		file := filepath.Join(directory, "original.png")
		if state.Mode == "escaping-media" {
			file = filepath.Join(filepath.Dir(directory), "escape.png")
		}
		f, err := os.Create(file)
		if err != nil {
			return fmt.Errorf("create fixture original: %w", err)
		}
		canvas := image.NewRGBA(image.Rect(0, 0, 32, 32))
		canvas.Set(0, 0, color.RGBA{R: 255, A: 255})
		if err := png.Encode(f, canvas); err != nil {
			f.Close()
			return fmt.Errorf("encode fixture original: %w", err)
		}
		if err := f.Close(); err != nil {
			return fmt.Errorf("close fixture original: %w", err)
		}
		if err := emit(map[string]any{"files": []any{map[string]any{"path": file}}}); err != nil {
			return fmt.Errorf("acknowledge download: %w", err)
		}
	default:
		return fmt.Errorf("unsupported fixture operation")
	}
	if state.Mode == "cleanup-warning" {
		return fmt.Errorf("simulated wrapper cleanup failure after successful envelope")
	}
	return nil
}

// TestNativeHelper becomes a subprocess fixture only when explicitly selected by argv and env.
func TestNativeHelper(t *testing.T) {
	if os.Getenv("COMFY_EXPERIMENT_HELPER") != "1" {
		return
	}
	args := os.Args
	for len(args) > 0 && args[0] != "--" {
		args = args[1:]
	}
	if len(args) < 3 {
		t.Fatal("fixture invocation missing")
	}
	args = args[1:]
	var err error
	if args[0] == "producer" {
		err = producerHelper(args[1])
	} else {
		err = comfyHelper(args[1], args[2:])
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

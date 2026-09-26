package acceptance_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

var runner = flag.String("runner", "", "packaged smoke test executable")

func TestMain(m *testing.M) {
	if dir := os.Getenv("TEST_TMPDIR"); dir != "" {
		if err := os.Setenv("TMPDIR", dir); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	os.Exit(m.Run())
}

type result struct {
	RunID       string            `json:"run_id"`
	RuntimeDir  string            `json:"runtime_dir"`
	Success     bool              `json:"success"`
	Interrupted bool              `json:"interrupted"`
	Accelerator string            `json:"accelerator"`
	Ports       map[string]int    `json:"ports"`
	Inputs      map[string]string `json:"inputs"`
	Tools       map[string]string `json:"tools"`
	Error       string            `json:"error"`
	QEMU        struct {
		PID   int    `json:"pid"`
		Start string `json:"start"`
	} `json:"qemu"`
	Cleanup struct {
		Success bool `json:"success"`
	} `json:"cleanup"`
}

type invocation struct {
	cmd           *exec.Cmd
	done          chan error
	state, output string
}

func start(t *testing.T, mode string, extra ...string) *invocation {
	t.Helper()
	return startConfigured(t, mode, nil, extra...)
}

func startConfigured(t *testing.T, mode string, overrides map[string]any, extra ...string) *invocation {
	t.Helper()
	r, err := runfiles.New()
	if err != nil {
		t.Fatal(err)
	}
	bin, err := r.Rlocation("_main/" + *runner)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	i := &invocation{state: filepath.Join(dir, "state"), output: filepath.Join(dir, "output"), done: make(chan error, 1)}
	args := append([]string{"--state-parent", i.state, "--outputs", i.output, "--acceptance-mode", mode}, extra...)
	if overrides != nil {
		data, err := os.ReadFile(bin + ".json")
		if err != nil {
			t.Fatal(err)
		}
		var config map[string]any
		if err := json.Unmarshal(data, &config); err != nil {
			t.Fatal(err)
		}
		for key, value := range overrides {
			config[key] = value
		}
		data, err = json.Marshal(config)
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(dir, "scenario.json")
		if err := os.WriteFile(path, data, 0o600); err != nil {
			t.Fatal(err)
		}
		args = append(args, "--config", path)
	}
	i.cmd = exec.Command(bin, args...)
	i.cmd.Env = append(os.Environ(), "VAULT_TOKEN=molecule-isolation-sentinel", "AWS_SECRET_ACCESS_KEY=molecule-isolation-sentinel", "SSH_AUTH_SOCK=/no/production/agent", "ANSIBLE_INVENTORY=/no/production/inventory")
	log, err := os.Create(filepath.Join(dir, "supervisor.log"))
	if err != nil {
		t.Fatal(err)
	}
	i.cmd.Stdout, i.cmd.Stderr = log, log
	if err := i.cmd.Start(); err != nil {
		t.Fatal(err)
	}
	go func() { i.done <- i.cmd.Wait(); log.Close() }()
	t.Cleanup(func() {
		_ = i.cmd.Process.Signal(syscall.SIGTERM)
		// The supervisor cleans handled cancellations. Preserve diagnostics outside
		// t.TempDir so Bazel retains them even when an assertion fails.
		if dst := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); dst != "" {
			artifactDir := filepath.Join(dst, t.Name(), filepath.Base(dir))
			_ = os.MkdirAll(artifactDir, 0o700)
			_ = filepath.Walk(i.output, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return err
				}
				data, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				return os.WriteFile(filepath.Join(artifactDir, filepath.Base(path)), data, 0o600)
			})
		}
	})
	return i
}

func finish(t *testing.T, i *invocation, wantSuccess bool) result {
	t.Helper()
	select {
	case err := <-i.done:
		if (err == nil) != wantSuccess {
			t.Fatalf("exit error = %v, want success %v; logs: %s", err, wantSuccess, i.output)
		}
	case <-time.After(15 * time.Minute):
		_ = i.cmd.Process.Signal(syscall.SIGTERM)
		t.Fatal("runner exceeded acceptance deadline")
	}
	b, err := os.ReadFile(filepath.Join(i.output, "result.json"))
	if err != nil {
		t.Fatal(err)
	}
	var r result
	if err := json.Unmarshal(b, &r); err != nil {
		t.Fatal(err)
	}
	if r.Success != wantSuccess || !r.Cleanup.Success {
		t.Fatalf("unexpected result: success=%v cleanup=%v error=%s; artifacts: %s", r.Success, r.Cleanup.Success, r.Error, i.output)
	}
	if r.RunID == "" || r.Accelerator == "" || len(r.Inputs) == 0 || len(r.Tools) == 0 {
		t.Fatal("missing evidence identity")
	}
	if _, err := os.Stat(r.RuntimeDir); !os.IsNotExist(err) {
		t.Fatalf("runtime survived: %s", r.RuntimeDir)
	}
	assertNoOwnedProcesses(t, r.RuntimeDir)
	if data, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", r.QEMU.PID)); err == nil && r.QEMU.PID > 1 {
		end := strings.LastIndex(string(data), ")")
		fields := strings.Fields(string(data)[end+1:])
		if len(fields) > 19 && fields[19] == r.QEMU.Start && fields[0] != "Z" {
			t.Fatal("owned QEMU process survived cleanup")
		}
	}
	_ = filepath.Walk(i.output, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		for _, secret := range []string{"molecule-isolation-sentinel", "BEGIN OPENSSH PRIVATE KEY", "BEGIN PRIVATE KEY", "BEGIN RSA PRIVATE KEY"} {
			if strings.Contains(string(data), secret) {
				t.Errorf("secret retained in %s", path)
			}
		}
		return nil
	})
	return r
}

func TestSmokeAndIndependentConcurrentRuns(t *testing.T) {
	a, b := start(t, ""), start(t, "")
	ra, rb := finish(t, a, true), finish(t, b, true)
	if ra.RunID == rb.RunID || ra.RuntimeDir == rb.RuntimeDir || ra.Ports["22"] == rb.Ports["22"] {
		t.Fatal("concurrent runs shared resources")
	}
	c := finish(t, start(t, ""), true)
	if c.RunID == ra.RunID || c.RunID == rb.RunID {
		t.Fatal("repeat reused a previous run")
	}
}

func TestFailureCleansResources(t *testing.T) {
	i := start(t, "failure")
	finish(t, i, false)
	data, err := os.ReadFile(filepath.Join(i.output, "events.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	assertionFound := false
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		var event struct{ Task, Outcome string }
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			t.Fatal(err)
		}
		if event.Task == "Inject deliberate verification failure" && event.Outcome == "failed" {
			assertionFound = true
		}
	}
	if !assertionFound {
		t.Fatal("injected assertion outcome was not retained")
	}
}

func TestSoftwareAcceleration(t *testing.T) {
	r := finish(t, start(t, "", "--accelerator=tcg"), true)
	if r.Accelerator != "tcg" {
		t.Fatal("software acceleration was not selected")
	}
}

func TestBootDeadlineCleansResources(t *testing.T) {
	started := time.Now()
	r := finish(t, startConfigured(t, "", map[string]any{"boot_timeout": 1, "tcg_boot_timeout": 1}), false)
	if time.Since(started) > 90*time.Second {
		t.Fatal("boot failure exceeded bounded cleanup deadline")
	}
	if r.Success {
		t.Fatal("unbootable deadline unexpectedly succeeded")
	}
}

func TestUnavailablePrerequisiteFailsBeforeGuestCreation(t *testing.T) {
	i := start(t, "", "--accelerator=unavailable")
	select {
	case err := <-i.done:
		if err == nil {
			t.Fatal("invalid accelerator succeeded")
		}
	case <-time.After(30 * time.Second):
		t.Fatal("preflight did not fail promptly")
	}
	var r result
	data, err := os.ReadFile(filepath.Join(i.output, "result.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &r); err != nil {
		t.Fatal(err)
	}
	if r.Success || !r.Cleanup.Success || r.QEMU.PID != 0 {
		t.Fatalf("unexpected prerequisite result: %s", data)
	}
	if !strings.Contains(string(data), "prerequisite: accelerator") {
		t.Fatal("missing actionable prerequisite diagnostic")
	}
	if _, err := os.Stat(r.RuntimeDir); !os.IsNotExist(err) {
		t.Fatal("preflight leaked runtime state")
	}
}

func TestMissingGuestInputFailsBeforeGuestCreation(t *testing.T) {
	i := startConfigured(t, "", map[string]any{"image": "_main/missing-guest-image.qcow2"})
	select {
	case err := <-i.done:
		if err == nil {
			t.Fatal("missing guest input succeeded")
		}
	case <-time.After(30 * time.Second):
		t.Fatal("missing input did not fail promptly")
	}
	var r result
	data, err := os.ReadFile(filepath.Join(i.output, "result.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &r); err != nil {
		t.Fatal(err)
	}
	if r.Success || !r.Cleanup.Success || r.QEMU.PID != 0 {
		t.Fatal("guest input failure did not clean up before boot")
	}
	if !strings.Contains(string(data), "prerequisite") {
		t.Fatal("missing prerequisite diagnostic")
	}
	if _, err := os.Stat(r.RuntimeDir); !os.IsNotExist(err) {
		t.Fatal("missing input leaked runtime state")
	}
}

func TestHandledCancellation(t *testing.T) {
	i := start(t, "cancellation")
	waitPhase(t, i, "verify")
	started := time.Now()
	if err := i.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	r := finish(t, i, false)
	if !r.Interrupted {
		t.Fatal("cancellation not recorded")
	}
	if time.Since(started) > 60*time.Second {
		t.Fatal("handled cancellation exceeded cleanup deadline")
	}
}

func waitPhase(t *testing.T, i *invocation, phase string) string {
	t.Helper()
	deadline := time.Now().Add(10 * time.Minute)
	for time.Now().Before(deadline) {
		matches, _ := filepath.Glob(filepath.Join(i.state, "*", "manifest.json"))
		if len(matches) == 1 {
			data, _ := os.ReadFile(matches[0])
			var state struct {
				Phase string `json:"phase"`
			}
			_ = json.Unmarshal(data, &state)
			if state.Phase == phase {
				return filepath.Dir(matches[0])
			}
		}
		select {
		case err := <-i.done:
			t.Fatalf("runner exited before %s: %v", phase, err)
		default:
		}
		time.Sleep(time.Second)
	}
	t.Fatal("guest never reached verification")
	return ""
}

func TestAbruptRecoveryAndRepeatableDestroy(t *testing.T) {
	i := start(t, "cancellation")
	root := waitPhase(t, i, "verify")
	data, err := os.ReadFile(filepath.Join(root, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var before result
	if err := json.Unmarshal(data, &before); err != nil {
		t.Fatal(err)
	}
	unrelated := exec.Command("/usr/bin/sleep", "300")
	if err := unrelated.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = unrelated.Process.Kill(); _ = unrelated.Wait() }()
	if err := i.cmd.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	if err := <-i.done; err == nil {
		t.Fatal("SIGKILL unexpectedly succeeded")
	}
	for attempt := 0; attempt < 2; attempt++ {
		cmd := exec.Command(i.cmd.Path, "destroy", root)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("recovery failed: %v: %s", err, output)
		}
	}
	if err := unrelated.Process.Signal(syscall.Signal(0)); err != nil {
		t.Fatal("recovery affected unrelated process")
	}
	if _, err := os.Stat(root); !os.IsNotExist(err) {
		t.Fatal("recovery left private runtime state")
	}
	assertNoOwnedProcesses(t, root)
	evidence := map[string]any{
		"run_id":             before.RunID,
		"qemu":               before.QEMU,
		"destroy_attempts":   2,
		"unrelated_survived": true,
		"runtime_removed":    true,
	}
	data, err = json.MarshalIndent(evidence, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(i.output, "recovery.json"), append(data, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
}

func assertNoOwnedProcesses(t *testing.T, root string) {
	t.Helper()
	entries, err := os.ReadDir("/proc")
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		data, err := os.ReadFile(filepath.Join("/proc", entry.Name(), "environ"))
		if err != nil {
			continue
		}
		if strings.Contains("\x00"+string(data), "\x00MOLECULE_RUNNER_STATE="+root+"\x00") {
			t.Errorf("owned subprocess survived cleanup: PID %s", entry.Name())
		}
	}
}

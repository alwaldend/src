package worker_test

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/api/workerpb"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/model"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/test/exportfixture"
	"github.com/bazelbuild/rules_go/go/runfiles"
)

// cleanerBinary is the built command supplied by Bazel.
var cleanerBinary = flag.String("cleaner", "", "cleaner executable runfile")

// writeWorkerArchive publishes synthetic immutable data under the configured source.
func writeWorkerArchive(t *testing.T, socket, key, account string, items []model.Item) {
	t.Helper()
	data, err := exportfixture.Encode(account, items)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(filepath.Dir(socket), key), data, 0o600); err != nil {
		t.Fatal(err)
	}
}

// startWorker creates a private socket directory and returns an idempotent stop.
func startWorker(t *testing.T, timeout, retention string) (string, func()) {
	t.Helper()
	directory, err := os.MkdirTemp("", "cleaner-async-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(directory); err != nil {
			t.Error(err)
		}
	})
	socket := filepath.Join(directory, "worker.sock")
	return socket, startWorkerAt(t, socket, timeout, retention)
}

// startWorkerAt supports restarting the executable on the same private socket.
func startWorkerAt(t *testing.T, socket, timeout, retention string) func() {
	t.Helper()
	binary, err := runfiles.Rlocation(*cleanerBinary)
	if err != nil {
		t.Fatal(err)
	}
	var logs bytes.Buffer
	writeWorkerArchive(t, socket, "empty.zip", "100", []model.Item{})
	process := exec.Command(binary, "worker", "serve", "--socket", socket, "--source", "exports="+filepath.Dir(socket), "--evaluation-timeout", timeout, "--result-ttl", retention)
	process.Env = []string{"PATH=/usr/bin:/bin", "UNRELATED_SECRET=synthetic-worker-marker"}
	process.Stdout = &logs
	process.Stderr = &logs
	if err := process.Start(); err != nil {
		t.Fatal(err)
	}
	finished := make(chan error, 1)
	go func() { finished <- process.Wait() }()
	var once sync.Once
	stop := func() {
		once.Do(func() {
			if err := process.Process.Signal(os.Interrupt); err != nil {
				t.Error(err)
			}
			select {
			case err := <-finished:
				if err != nil {
					t.Errorf("worker stopped: %v %s", err, logs.String())
				}
			case <-time.After(6 * time.Second):
				_ = process.Process.Kill()
				<-finished
				t.Error("worker shutdown timeout")
			}
			if strings.Contains(logs.String(), "synthetic-worker-marker") {
				t.Error("environment leaked to logs")
			}
		})
	}
	t.Cleanup(stop)
	deadline := time.Now().Add(10 * time.Second)
	for {
		if _, err := os.Stat(socket); err == nil {
			return stop
		}
		if time.Now().After(deadline) {
			t.Fatal("worker socket did not appear")
		}
		time.Sleep(5 * time.Millisecond)
	}
}

// waitStatus waits only for a terminal fixture state under the test deadline.
func waitStatus(t *testing.T, ctx context.Context, rpc workerpb.ArchiveWorkerClient, ref *workerpb.EvaluationRef) *workerpb.EvaluationStatus {
	t.Helper()
	for {
		job, err := rpc.GetEvaluation(ctx, ref)
		if err != nil {
			t.Fatal(err)
		}
		if job.State != workerpb.EvaluationState_RUNNING {
			return job
		}
		select {
		case <-ctx.Done():
			t.Fatal(ctx.Err())
		case <-time.After(5 * time.Millisecond):
		}
	}
}

// writeAsyncArtifact preserves assertions established against the built process.
func writeAsyncArtifact(t *testing.T, name string, value any) {
	t.Helper()
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		data, err := json.MarshalIndent(value, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(directory, name), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// datasetZIP builds independent input without calling production archive writers.
func datasetZIP(t *testing.T, entries map[string]string) []byte {
	t.Helper()
	var data bytes.Buffer
	writer := zip.NewWriter(&data)
	for name, value := range entries {
		entry, err := writer.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := entry.Write([]byte(value)); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return data.Bytes()
}

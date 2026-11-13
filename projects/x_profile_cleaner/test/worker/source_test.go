package worker_test

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/api/workerpb"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/archive"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/cleanup"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// sourceCleanerBinary identifies the built executable for this process fixture.
var sourceCleanerBinary = flag.String("source-cleaner", "", "cleaner executable runfile")

// sourceArchive builds an independent standard export fixture without runtime code.
func sourceArchive(t *testing.T, entries map[string]string) []byte {
	t.Helper()
	var data bytes.Buffer
	writer := zip.NewWriter(&data)
	for name, value := range entries {
		entry, err := writer.CreateHeader(&zip.FileHeader{Name: name, Method: zip.Store})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := io.WriteString(entry, value); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return data.Bytes()
}

// sourceEntries supplies account identity, multipart posts and unrelated private data.
func sourceEntries() map[string]string {
	return map[string]string{
		"data/account.js":         `window.YTD.account.part0 = [{"account":{"accountId":"100","username":"fixture"}}];`,
		"data/tweets.js":          `window.YTD.tweets.part0 = [{"tweet":{"id_str":"101","created_at":"Mon Nov 09 12:00:00 +0000 2020","full_text":"root"}}];`,
		"data/tweets-part1.js":    `window.YTD.tweets.part1 = [{"tweet":{"id_str":"2","created_at":"Tue Nov 10 12:00:00 +0000 2020","full_text":"reply","in_reply_to_status_id_str":"101"}}];`,
		"data/direct-messages.js": "synthetic-private-data-must-not-be-ingested",
	}
}

// TestConfiguredArchiveSources submits only references to directly fetched exports.
func TestConfiguredArchiveSources(t *testing.T) {
	directory, err := os.MkdirTemp("", "cleaner-sources-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(directory); err != nil {
			t.Error(err)
		}
	})
	archive := sourceArchive(t, sourceEntries())
	if err := os.WriteFile(filepath.Join(directory, "export.zip"), archive, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("/etc/passwd", filepath.Join(directory, "escape.zip")); err != nil {
		t.Fatal(err)
	}
	bad := sourceEntries()
	bad["../escape"] = "unsafe entry"
	if err := os.WriteFile(filepath.Join(directory, "unsafe.zip"), sourceArchive(t, bad), 0o600); err != nil {
		t.Fatal(err)
	}

	for name, update := range map[string]func(map[string]string){
		"bad-counter.zip": func(entries map[string]string) {
			entries["data/tweets.js"] = strings.Replace(entries["data/tweets.js"], `"full_text":"root"`, `"full_text":"root","favorite_count":"-1"`, 1)
		},
		"bad-entity.zip": func(entries map[string]string) {
			entries["data/tweets.js"] = strings.Replace(entries["data/tweets.js"], `"full_text":"root"`, `"full_text":"root","entities":{"hashtags":[{"indices":[2,1],"text":"bad"}]}`, 1)
		},
		"bad-field.zip": func(entries map[string]string) {
			entries["data/tweets.js"] = strings.Replace(entries["data/tweets.js"], `"full_text":"root"`, `"full_text":42`, 1)
		},
		"executable.zip": func(entries map[string]string) { entries["data/tweets.js"] += "globalThis.executed=true;" },
		"missing-part.zip": func(entries map[string]string) {
			entries["data/tweets-part2.js"] = strings.Replace(entries["data/tweets-part1.js"], "part1", "part2", 1)
			delete(entries, "data/tweets-part1.js")
		},
		"duplicate.zip": func(entries map[string]string) {
			entries["data/tweets-part1.js"] = strings.ReplaceAll(entries["data/tweets-part1.js"], `"id_str":"2"`, `"id_str":"101"`)
		},
		"null-text.zip": func(entries map[string]string) {
			entries["data/tweets.js"] = strings.Replace(entries["data/tweets.js"], `"full_text":"root"`, `"full_text":null`, 1)
		},
		"null-language.zip": func(entries map[string]string) {
			entries["data/tweets.js"] = strings.Replace(entries["data/tweets.js"], `"full_text":"root"`, `"full_text":"root","lang":null`, 1)
		},
	} {
		entries := sourceEntries()
		update(entries)
		if err := os.WriteFile(filepath.Join(directory, name), sourceArchive(t, entries), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	corrupt := bytes.Clone(archive)
	position := bytes.Index(corrupt, []byte(`"full_text":"root"`))
	if position < 0 {
		t.Fatal("fixture missing stored tweet data")
	}
	corrupt[position+13] ^= 1
	if err := os.WriteFile(filepath.Join(directory, "corrupt.zip"), corrupt, 0o600); err != nil {
		t.Fatal(err)
	}
	excessive := bytes.Clone(archive)
	binary.LittleEndian.PutUint16(excessive[len(excessive)-14:], 10001)
	binary.LittleEndian.PutUint16(excessive[len(excessive)-12:], 10001)
	if err := os.WriteFile(filepath.Join(directory, "excessive.zip"), excessive, 0o600); err != nil {
		t.Fatal(err)
	}

	for name, data := range map[string][]byte{"directory-count.zip": forgedDirectory(false), "directory-size.zip": forgedDirectory(true), "directory-locator.zip": directoryLocator()} {
		if err := os.WriteFile(filepath.Join(directory, name), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	invalidLength := bytes.Clone(archive)
	directoryOffset := binary.LittleEndian.Uint32(invalidLength[len(invalidLength)-6:])
	binary.LittleEndian.PutUint16(invalidLength[int(directoryOffset)+28:], 65535)
	if err := os.WriteFile(filepath.Join(directory, "directory-length.zip"), invalidLength, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(filepath.Join(directory, "pipe.zip"), 0o600); err != nil {
		t.Fatal(err)
	}
	binary, err := runfiles.Rlocation(*sourceCleanerBinary)
	if err != nil {
		t.Fatal(err)
	}
	socket := filepath.Join(directory, "worker.sock")
	process := exec.Command(binary, "worker", "serve", "--socket", socket, "--source", "local="+directory, "--source", "secondary="+directory, "--evaluation-timeout", "2s")
	process.Env = []string{"PATH=/usr/bin:/bin", "UNRELATED_SECRET=synthetic-private-environment"}
	var logs bytes.Buffer
	process.Stdout, process.Stderr = &logs, &logs
	if err := process.Start(); err != nil {
		t.Fatal(err)
	}
	finished := make(chan error, 1)
	go func() { finished <- process.Wait() }()
	stopped := false
	t.Cleanup(func() {
		if stopped {
			return
		}
		if err := process.Process.Signal(os.Interrupt); err != nil {
			t.Error(err)
		}
		select {
		case err := <-finished:
			if err != nil {
				t.Errorf("worker shutdown: %v %s", err, logs.String())
			}
		case <-time.After(6 * time.Second):
			_ = process.Process.Kill()
			<-finished
			t.Error("worker shutdown timeout")
		}
	})
	deadline := time.Now().Add(10 * time.Second)
	for {
		if _, err := os.Stat(socket); err == nil {
			break
		}
		select {
		case err := <-finished:
			stopped = true
			t.Fatalf("source worker did not start: %v %s", err, logs.String())
		default:
		}
		if time.Now().After(deadline) {
			t.Fatal("source worker readiness timeout")
		}
		time.Sleep(5 * time.Millisecond)
	}
	connection, err := grpc.NewClient("unix://"+socket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	rpc := workerpb.NewArchiveWorkerClient(connection)
	for _, test := range []struct {
		// name is the operation ID and expected source fixture.
		name string
		// source selects a trusted configured provider.
		source string
		// key is a source-relative object name, never a path or endpoint.
		key string
		// account binds job access and the exported account.
		account string
		// success requires a complete exact result; failures contain no plan.
		success bool
	}{
		{"local-export", "local", "export.zip", "100", true},
		{"secondary-export", "secondary", "export.zip", "100", true},
		{"wrong-account", "local", "export.zip", "200", false},
		{"unsafe-entry", "local", "unsafe.zip", "100", false},
		{"source-escape", "local", "escape.zip", "100", false},
		{"bad-counter", "local", "bad-counter.zip", "100", false},
		{"bad-entity", "local", "bad-entity.zip", "100", false},
		{"bad-field", "local", "bad-field.zip", "100", false},
		{"executable", "local", "executable.zip", "100", false},
		{"missing-part", "local", "missing-part.zip", "100", false},
		{"duplicate", "local", "duplicate.zip", "100", false},
		{"null-text", "local", "null-text.zip", "100", false},
		{"null-language", "local", "null-language.zip", "100", true},
		{"crc", "local", "corrupt.zip", "100", false},
		{"entries", "local", "excessive.zip", "100", false},
		{"pipe", "local", "pipe.zip", "100", false},
		{"directory-count", "local", "directory-count.zip", "100", false},
		{"directory-size", "local", "directory-size.zip", "100", false},
		{"directory-locator", "local", "directory-locator.zip", "100", false},
		{"directory-length", "local", "directory-length.zip", "100", false},
		{"directory-recovery", "local", "export.zip", "100", true},
	} {
		request := &workerpb.Evaluation{OperationId: test.name, AccountId: test.account, EvaluatedAt: timestamppb.Now(), Script: cleanup.Script(`assert(get_post("101").language==nil); delete_post_when(function(p) return true end); preserve_post("101"); force_delete_post("2")`)}
		// An independent unary caller encodes the new source/key wire fields.
		unknown := append([]byte{0x2a, byte(len(test.source))}, []byte(test.source)...)
		unknown = append(unknown, 0x32, byte(len(test.key)))
		unknown = append(unknown, []byte(test.key)...)
		request.ProtoReflect().SetUnknown(unknown)
		var ack workerpb.SubmitEvaluationResponse
		if err := connection.Invoke(ctx, "/x_profile_cleaner.worker.ArchiveWorker/SubmitEvaluation", request, &ack); err != nil {
			t.Fatal(err)
		}
		if ack.JobId != test.name {
			t.Fatalf("submission binding: %s", ack.JobId)
		}
		ref := &workerpb.EvaluationRef{JobId: ack.JobId, AccountId: test.account}
		job := waitSourceStatus(t, ctx, rpc, ref)
		if test.success {
			if job.State != workerpb.EvaluationState_SUCCEEDED || job.Plan == nil || len(sourcePlan(t, job).Operations) != 1 || sourcePlan(t, job).Operations[0].GetDeletePost().PostId != "2" || !sourcePlan(t, job).Operations[0].GetDeletePost().Forced {
				t.Fatalf("archive plan: %+v", job)
			}
		} else if job.State != workerpb.EvaluationState_FAILED || job.Plan != nil {
			t.Fatalf("unsafe archive was accepted: %+v", job)
		}
		if _, err := rpc.GetEvaluation(ctx, &workerpb.EvaluationRef{JobId: ack.JobId, AccountId: "300"}); status.Code(err) != codes.NotFound {
			t.Fatalf("foreign result access: %v", err)
		}
		if _, err := rpc.DeleteEvaluation(ctx, ref); err != nil {
			t.Fatal(err)
		}
	}
	for _, request := range []*workerpb.Evaluation{
		{OperationId: "unknown-source", SourceName: "missing", ArchiveKey: "export.zip"},
		{OperationId: "parent-key", SourceName: "local", ArchiveKey: "../export.zip"},
		{OperationId: "absolute-key", SourceName: "local", ArchiveKey: "/etc/passwd"},
	} {
		request.AccountId, request.Script, request.EvaluatedAt = "100", cleanup.Script(""), timestamppb.Now()
		if _, err := rpc.SubmitEvaluation(ctx, request); status.Code(err) != codes.InvalidArgument {
			t.Fatalf("unsafe reference accepted: %v", err)
		}
	}
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		data, err := json.MarshalIndent(map[string]any{"local": true, "multiple_sources": true, "payload_upload": false, "exact_target": "2", "unsafe_input_rejected": true}, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(directory, "source-archives.json"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// waitSourceStatus observes a detached job until completion under the caller deadline.
func waitSourceStatus(t *testing.T, ctx context.Context, rpc workerpb.ArchiveWorkerClient, ref *workerpb.EvaluationRef) *workerpb.EvaluationStatus {
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

// observedArchive counts source reads to distinguish early rejection from late failure.
type observedArchive struct {
	// source supplies independently encoded ZIP bytes.
	source *bytes.Reader
	// bytes records actual I/O, including failed parsing attempts.
	bytes int
	// prefixLimit identifies file data that directory parsing must not access.
	prefixLimit int64
	// prefixReads counts reads outside the declared metadata region.
	prefixReads int
	// cancel optionally interrupts after the footer has been read.
	cancel context.CancelFunc
}

// ReadAt records bytes and can cancel during directory scanning.
func (r *observedArchive) ReadAt(data []byte, offset int64) (int, error) {
	n, err := r.source.ReadAt(data, offset)
	r.bytes += n
	if offset < r.prefixLimit {
		r.prefixReads++
	}
	if r.cancel != nil && r.bytes > 65557 {
		r.cancel()
	}
	return n, err
}

// forgedDirectory independently encodes too many headers behind a small footer count.
func forgedDirectory(shortSize bool) []byte {
	header := make([]byte, 52)
	binary.LittleEndian.PutUint32(header, 0x02014b50)
	binary.LittleEndian.PutUint16(header[28:], 6)
	copy(header[46:], "unused")
	directory := bytes.Repeat(header, 25000)
	footer := make([]byte, 22)
	binary.LittleEndian.PutUint32(footer, 0x06054b50)
	binary.LittleEndian.PutUint16(footer[8:], 2)
	binary.LittleEndian.PutUint16(footer[10:], 2)
	size := len(directory)
	if shortSize {
		size = 64
	}
	binary.LittleEndian.PutUint32(footer[12:], uint32(size))
	return append(directory, footer...)
}

// directoryLocator encodes a locator in a 65,535-byte ordinary directory comment.
// The pinned Go reader probes ZIP64 for this directory length even with a small count.
func directoryLocator() []byte {
	prefix := make([]byte, 56)
	binary.LittleEndian.PutUint32(prefix, 0x06064b50)
	binary.LittleEndian.PutUint64(prefix[24:], 2000)
	binary.LittleEndian.PutUint64(prefix[32:], 2000)
	directory := make([]byte, 65535)
	binary.LittleEndian.PutUint32(directory, 0x02014b50)
	binary.LittleEndian.PutUint16(directory[28:], 6)
	binary.LittleEndian.PutUint16(directory[32:], uint16(len(directory)-52))
	copy(directory[46:], "unused")
	locator := directory[len(directory)-20:]
	binary.LittleEndian.PutUint32(locator, 0x07064b50)
	binary.LittleEndian.PutUint32(locator[16:], 1)
	footer := make([]byte, 22)
	binary.LittleEndian.PutUint32(footer, 0x06054b50)
	binary.LittleEndian.PutUint16(footer[8:], 1)
	binary.LittleEndian.PutUint16(footer[10:], 1)
	binary.LittleEndian.PutUint32(footer[12:], uint32(len(directory)))
	binary.LittleEndian.PutUint32(footer[16:], uint32(len(prefix)))
	return append(append(prefix, directory...), footer...)
}

// TestArchiveDirectoryReadBounds observes allocation-bound inputs before worker IPC.
// RPC rejection alone cannot distinguish bounded preflight from an oversized parse.
func TestArchiveDirectoryReadBounds(t *testing.T) {
	evidence := map[string]any{}
	for _, shortSize := range []bool{false, true} {
		name := "underreported-count"
		if shortSize {
			name = "underreported-size"
		}
		t.Run(name, func(t *testing.T) {
			data := forgedDirectory(shortSize)
			reader := &observedArchive{source: bytes.NewReader(data)}
			if _, err := archive.Open(context.Background(), reader, int64(len(data))); err == nil {
				t.Fatal("forged directory accepted")
			}
			evidence[name] = reader.bytes
			if reader.bytes > 128<<10 {
				t.Fatalf("directory rejection consumed %d bytes; expected bounded preflight", reader.bytes)
			}
		})
	}
	for _, during := range []bool{false, true} {
		name := "cancel-before-io"
		if during {
			name = "cancel-during-scan"
		}
		t.Run(name, func(t *testing.T) {
			data := forgedDirectory(false)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			reader := &observedArchive{source: bytes.NewReader(data)}
			if during {
				reader.cancel = cancel
			} else {
				cancel()
			}
			_, err := archive.Open(ctx, reader, int64(len(data)))
			if !errors.Is(err, context.Canceled) || reader.bytes > 128<<10 || (!during && reader.bytes != 0) {
				t.Fatalf("cancelled parsing: bytes=%d error=%v", reader.bytes, err)
			}
			evidence[name] = reader.bytes
		})
	}
	t.Run("zip64-locator", func(t *testing.T) {
		data := directoryLocator()
		reader := &observedArchive{source: bytes.NewReader(data), prefixLimit: 56}
		if _, err := archive.Open(context.Background(), reader, int64(len(data))); err == nil || reader.prefixReads != 0 {
			t.Fatalf("ZIP64 probe escaped preflight: prefix reads=%d error=%v", reader.prefixReads, err)
		}
		evidence["zip64-locator-prefix-reads"] = reader.prefixReads
	})

	t.Run("comment-footer", func(t *testing.T) {
		data := sourceArchive(t, sourceEntries())
		forged := forgedDirectory(true)
		comment := append(bytes.Clone(forged[len(forged)-22:]), []byte("trailing comment")...)
		binary.LittleEndian.PutUint16(data[len(data)-2:], uint16(len(comment)))
		data = append(data, comment...)
		export, err := archive.Open(context.Background(), bytes.NewReader(data), int64(len(data)))
		if err != nil {
			t.Fatal(err)
		}
		records, err := export.Tweets(context.Background())
		if err != nil || export.AccountID != "100" || len(records) != 2 {
			t.Fatalf("valid ZIP comment changed directory selection: %v", err)
		}
		evidence["comment-footer"] = "same validated inventory"
	})
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		data, err := json.MarshalIndent(evidence, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(directory, "archive-directory-bounds.json"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// sourcePlan inspects the cleanup application result from the archive process fixture.
func sourcePlan(t *testing.T, job *workerpb.EvaluationStatus) *workerpb.Plan {
	t.Helper()
	err := cleanup.Validate(job.Plan)
	if err != nil {
		t.Fatal(err)
	}
	return job.Plan
}

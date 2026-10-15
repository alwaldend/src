package cli_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// cliBinary identifies the compiled command under test.
var (
	cliBinary = flag.String("cli", "", "CLI runfile")
	// workerBinary identifies the native server for external transport checks.
	workerBinary = flag.String("worker", "", "worker runfile")
)

// executable resolves a build-owned runfile without executable search paths.
func executable(t *testing.T, name string) string {
	t.Helper()
	path, err := runfiles.Rlocation(name)
	if err != nil {
		t.Fatal(err)
	}
	return path
}

// run invokes the built CLI with a bounded lifetime and separate output channels.
func run(t *testing.T, input string, args ...string) ([]byte, []byte, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, executable(t, *cliBinary), args...)
	cmd.Env = append(os.Environ(), "PATH=")
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}

// invoke decodes a successful ProtoJSON response for subsequent native operations.
func invoke(t *testing.T, input string, args ...string) map[string]any {
	t.Helper()
	stdout, stderr, err := run(t, input, args...)
	if err != nil {
		t.Fatalf("%v: %v\n%s", args, err, stderr)
	}
	var result map[string]any
	if err := json.Unmarshal(stdout, &result); err != nil {
		t.Fatalf("invalid JSON %q: %v", stdout, err)
	}
	return result
}

// request encodes synthetic public native requests as JSON.
func request(t *testing.T, value any) string {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// TestCLIWorkflow checks generated help, local persistence, inputs, and native packages.
func TestCLIWorkflow(t *testing.T) {
	stdout, stderr, err := run(t, "", "--worker", "/missing/worker", "anki.notes.NotesService", "GetNote", "--help")
	if err != nil || len(stderr) != 0 || !bytes.Contains(stdout, []byte("anki.notes.NoteId")) || !bytes.Contains(stdout, []byte("anki.notes.Note")) {
		t.Fatalf("method help: %s %s %v", stdout, stderr, err)
	}
	root := t.TempDir()
	database := filepath.Join(root, "collection.anki2")
	local := []string{"--collection", database}
	nt := invoke(t, "", append(local, "anki.notetypes.NotetypesService", "GetNotetypeIdByName", "--request", `{"val":"Basic"}`)...)
	noteRequest := map[string]any{"deck_id": "1", "note": map[string]any{"notetype_id": nt["ntid"], "fields": []string{"Front 中文", "Back"}, "tags": []string{"fixture"}}}
	file := filepath.Join(root, "request.json")
	if err := os.WriteFile(file, []byte(request(t, noteRequest)), 0o600); err != nil {
		t.Fatal(err)
	}
	added := invoke(t, "", append(local, "anki.notes.NotesService", "AddNote", "--request", "@"+file)...)
	nid := added["note_id"]
	getRequest := request(t, map[string]any{"nid": nid})
	note := invoke(t, getRequest, append(local, "anki.notes.NotesService", "GetNote", "--request", "-")...)
	if note["fields"].([]any)[0] != "Front 中文" {
		t.Fatalf("native content: %v", note)
	}
	note["fields"] = []string{"Changed\n中文", "Back"}
	invoke(t, "", append(local, "anki.notes.NotesService", "UpdateNotes", "--request", request(t, map[string]any{"notes": []any{note}}))...)
	persisted := invoke(t, "", append(local, "anki.notes.NotesService", "GetNote", "--request", getRequest)...)
	if persisted["fields"].([]any)[0] != "Changed\n中文" {
		t.Fatalf("persistence: %v", persisted)
	}
	invoke(t, "", append(local, "anki.media.MediaService", "AddMediaFile", "--request", `{"desired_name":"example.txt","data":"cHVibGljIG1lZGlh"}`)...)
	archive := filepath.Join(root, "collection.colpkg")
	invoke(t, "", append(local, "anki.import_export.BackendImportExportService", "ExportCollectionPackage", "--request", request(t, map[string]any{"out_path": archive, "include_media": true}))...)
	imported := filepath.Join(root, "imported.anki2")
	invoke(t, "", "anki.import_export.BackendImportExportService", "ImportCollectionPackage", "--request", request(t, map[string]any{"backup_path": archive, "col_path": imported, "media_folder": imported + ".media", "media_db": imported + ".media.db"}))
	restored := invoke(t, "", "--collection", imported, "anki.notes.NotesService", "GetNote", "--request", getRequest)
	if restored["fields"].([]any)[0] != "Changed\n中文" {
		t.Fatalf("package persistence: %v", restored)
	}
	media, err := os.ReadFile(imported + ".media/example.txt")
	if err != nil || string(media) != "public media" {
		t.Fatalf("media: %q %v", media, err)
	}
	for _, args := range [][]string{
		{"--worker", "/missing/worker", "anki.notes.NotesService", "GetNote", "--request", "{"},
		{"anki.notes.NotesService", "GetNote", "--request", `{"unknown":1}`},
		{"--collection", database, "anki.notes.NotesService", "GetNote", "--request", `{"nid":"1"}`},
		{"--worker", "/missing/worker", "anki.collection.BackendCollectionService", "LatestProgress"},
		{"--timeout", "0s", "anki.collection.BackendCollectionService", "LatestProgress"},
		{"--timeout", "1ns", "--collection", database, "anki.notes.NotesService", "GetNote", "--request", getRequest},
		{"--target", "127.0.0.1:1", "--timeout", "20ms", "anki.collection.BackendCollectionService", "LatestProgress"},
		{"--target", "127.0.0.1:1", "--collection", database, "anki.notes.NotesService", "GetNote"},
	} {
		stdout, stderr, err := run(t, "", args...)
		if err == nil || len(stdout) != 0 || len(stderr) == 0 {
			t.Fatalf("expected failure %v: %q %q %v", args, stdout, stderr, err)
		}
	}
	if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
		for name, data := range map[string][]byte{"note.json": []byte(request(t, restored)), "collection.colpkg": mustRead(t, archive)} {
			if err := os.WriteFile(filepath.Join(artifacts, name), data, 0o600); err != nil {
				t.Fatal(err)
			}
		}
	}
}

// mustRead reads a repeatable native artifact for the Bazel output directory.
func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// TestExternalTransports keeps server state across separate TCP and Unix invocations.
func TestExternalTransports(t *testing.T) {
	for _, transport := range []string{"tcp", "unix"} {
		t.Run(transport, func(t *testing.T) {
			root := t.TempDir()
			listen := "127.0.0.1:0"
			if transport == "unix" {
				// Relative socket names also work in delivery checkouts with long scratch paths.
				t.Chdir(root)
				listen = "anki.sock"
			}
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			cmd := exec.CommandContext(ctx, executable(t, *workerBinary), "--transport", transport, "--listen", listen)
			diagnostics, err := cmd.StderrPipe()
			if err != nil {
				cancel()
				t.Fatal(err)
			}
			if err := cmd.Start(); err != nil {
				cancel()
				t.Fatal(err)
			}
			t.Cleanup(func() { cancel(); _ = cmd.Wait() })
			ready := make(chan string, 1)
			go func() {
				scanner := bufio.NewScanner(diagnostics)
				var lines []string
				for scanner.Scan() {
					if line := scanner.Text(); strings.HasPrefix(line, "Listening on ") {
						ready <- strings.TrimPrefix(line, "Listening on ")
						return
					}
					lines = append(lines, scanner.Text())
				}
				if err := scanner.Err(); err != nil {
					lines = append(lines, err.Error())
				}
				ready <- strings.Join(lines, "\n")
			}()
			var target string
			select {
			case target = <-ready:
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			}
			if !strings.HasPrefix(target, "tcp://") && !strings.HasPrefix(target, "unix://") {
				t.Fatalf("server startup failed: %s", target)
			}
			if transport == "tcp" {
				target = "passthrough:///" + strings.TrimPrefix(target, "tcp://")
			} else {
				target = "unix:anki.sock"
			}
			network := []string{"--target", target}
			database := filepath.Join(root, "collection.anki2")
			invoke(t, "", append(network, "anki.collection.BackendCollectionService", "OpenCollection", "--request", request(t, map[string]any{"collection_path": database, "media_folder_path": database + ".media", "media_db_path": database + ".media.db"}))...)
			nt := invoke(t, "", append(network, "anki.notetypes.NotetypesService", "GetNotetypeIdByName", "--request", `{"val":"Basic"}`)...)
			added := invoke(t, "", append(network, "anki.notes.NotesService", "AddNote", "--request", request(t, map[string]any{"deck_id": "1", "note": map[string]any{"notetype_id": nt["ntid"], "fields": []string{"Network", "Back"}}}))...)
			note := invoke(t, "", append(network, "anki.notes.NotesService", "GetNote", "--request", request(t, map[string]any{"nid": added["note_id"]}))...)
			if note["fields"].([]any)[0] != "Network" {
				t.Fatalf("shared server state: %v", note)
			}
			invoke(t, "", append(network, "anki.collection.BackendCollectionService", "CloseCollection")...)
			result := invoke(t, "", append(network, "anki.collection.BackendCollectionService", "LatestProgress")...)
			if result == nil {
				t.Fatal("external server stopped after collection close")
			}
		})
	}
}

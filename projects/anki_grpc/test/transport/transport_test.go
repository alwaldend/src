package transport

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/ankitects/anki/proto/anki/collection"
	"github.com/ankitects/anki/proto/anki/generic"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

// workerBinary identifies the native executable in Bazel runfiles.
var workerBinary = flag.String("worker", "", "Anki gRPC server runfile")

// workerPath resolves the executable before changing the test working directory.
func workerPath(t *testing.T) string {
	t.Helper()
	path, err := runfiles.Rlocation(*workerBinary)
	if err != nil {
		t.Fatalf("resolve server: %v", err)
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		t.Fatalf("resolve server path: %v", err)
	}
	return absolute
}

// TestCLI verifies discoverable usage and rejection before protocol startup.
func TestCLI(t *testing.T) {
	path := workerPath(t)
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, path, "--help").CombinedOutput()
	if err != nil {
		t.Fatalf("server help: %v; %s", err, output)
	}
	for _, text := range []string{"stdio", "tcp", "unix", "--listen", "default"} {
		if !strings.Contains(string(output), text) {
			t.Errorf("help omits %q: %s", text, output)
		}
	}
	for _, args := range [][]string{
		{"--unknown"},
		{"--transport"},
		{"--transport", "invalid"},
		{"--transport", "tcp"},
		{"--transport", "unix"},
		{"--listen", "127.0.0.1:0"},
		{"--transport", "tcp", "--listen"},
		{"--transport", "tcp", "--listen", "invalid-address"},
	} {
		t.Run(strings.Join(args, "_"), func(t *testing.T) {
			var output, diagnostics bytes.Buffer
			command := exec.CommandContext(ctx, path, args...)
			command.Stdout = &output
			command.Stderr = &diagnostics
			err := command.Run()
			if err == nil {
				t.Fatalf("invalid arguments accepted: %v", args)
			}
			if diagnostics.Len() == 0 {
				t.Fatal("missing argument diagnostic")
			}
			if output.Len() != 0 {
				t.Fatalf("invalid arguments emitted protocol output: %q", output.String())
			}
		})
	}
}

// TestNetworkCollection verifies startup initialization and persistence after reconnect.
func TestNetworkCollection(t *testing.T) {
	executable := workerPath(t)
	for _, transport := range []string{"tcp", "unix"} {
		t.Run(transport, func(t *testing.T) {
			dir := t.TempDir()
			// Relative Unix paths stay within the platform's socket name limit.
			t.Chdir(dir)
			listen := "127.0.0.1:0"
			if transport == "unix" {
				listen = "worker.sock"
			}
			ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
			defer cancel()
			command := exec.CommandContext(ctx, executable, "--transport", transport, "--listen", listen)
			command.Dir = dir
			command.Stdout = io.Discard
			stderr, err := command.StderrPipe()
			if err != nil {
				t.Fatalf("server stderr: %v", err)
			}
			if err := command.Start(); err != nil {
				t.Fatalf("start server: %v", err)
			}
			t.Cleanup(func() {
				_ = command.Process.Kill()
				_ = command.Wait()
			})
			ready := make(chan string, 1)
			go func() {
				scanner := bufio.NewScanner(stderr)
				if scanner.Scan() {
					ready <- scanner.Text()
				} else {
					ready <- ""
				}
				for scanner.Scan() {
				} // Drain later connection diagnostics.
			}()
			var address string
			select {
			case line := <-ready:
				prefix := "Listening on " + transport + "://"
				if !strings.HasPrefix(line, prefix) {
					t.Fatalf("listener did not start: %q", line)
				}
				address = strings.TrimPrefix(line, prefix)
			case <-ctx.Done():
				t.Fatalf("listener startup: %v", ctx.Err())
			}
			db := filepath.Join(dir, "collection.anki2")
			open := &collection.OpenCollectionRequest{CollectionPath: db, MediaFolderPath: db + ".media", MediaDbPath: db + ".media.db"}
			fields := []string{"<b>Transport</b>", transport}
			var noteID int64
			for session := range 2 {
				connection, err := grpc.NewClient("passthrough:///anki", grpc.WithTransportCredentials(insecure.NewCredentials()),
					grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
						connection, err := (&net.Dialer{}).DialContext(ctx, transport, address)
						if err != nil {
							return nil, fmt.Errorf("dial %s %q: %w", transport, address, err)
						}
						return connection, nil
					}))
				if err != nil {
					t.Fatalf("gRPC client: %v", err)
				}
				t.Cleanup(func() { _ = connection.Close() })
				invoke(t, ctx, connection, "anki.collection.BackendCollectionService/OpenCollection", open, &generic.Empty{})
				if session == 0 {
					nt := &notetypes.NotetypeId{}
					invoke(t, ctx, connection, "anki.notetypes.BackendNotetypesService/GetNotetypeIdByName", &generic.String{Val: "Basic"}, nt)
					added := &notes.AddNoteResponse{}
					invoke(t, ctx, connection, "anki.notes.BackendNotesService/AddNote", &notes.AddNoteRequest{Note: &notes.Note{NotetypeId: nt.Ntid, Fields: fields}, DeckId: 1}, added)
					noteID = added.NoteId
					if noteID == 0 {
						t.Fatal("Anki did not allocate a note")
					}
				}
				stored := &notes.Note{}
				invoke(t, ctx, connection, "anki.notes.BackendNotesService/GetNote", &notes.NoteId{Nid: noteID}, stored)
				if !reflect.DeepEqual(stored.Fields, fields) {
					t.Fatalf("persisted fields: %v", stored.Fields)
				}
				invoke(t, ctx, connection, "anki.collection.BackendCollectionService/CloseCollection", &collection.CloseCollectionRequest{}, &generic.Empty{})
				if err := connection.Close(); err != nil {
					t.Fatalf("disconnect client: %v", err)
				}
			}
			output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR")
			if output == "" {
				output = t.TempDir()
			}
			data, err := os.ReadFile(db)
			if err != nil {
				t.Fatalf("read closed collection: %v", err)
			}
			if err := os.WriteFile(filepath.Join(output, transport+"-collection.anki2"), data, 0o600); err != nil {
				t.Fatalf("emit collection: %v", err)
			}
			report, err := json.MarshalIndent(map[string]any{"transport": transport, "note_id": noteID, "fields": fields, "checks": []string{"persisted-note", "reconnect", "startup-initialization"}}, "", "  ")
			if err != nil {
				t.Fatalf("encode verification: %v", err)
			}
			if err := os.WriteFile(filepath.Join(output, transport+"-verification.json"), append(report, '\n'), 0o600); err != nil {
				t.Fatalf("emit verification: %v", err)
			}
		})
	}
}

// invoke executes an official Anki command through the server's gRPC endpoint.
func invoke(t *testing.T, ctx context.Context, connection *grpc.ClientConn, method string, input, output proto.Message) {
	t.Helper()
	if err := connection.Invoke(ctx, "/"+method, input, output); err != nil {
		t.Fatalf("%s: %v", method, err)
	}
}

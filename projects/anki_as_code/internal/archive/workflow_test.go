package archive_test

import (
	"context"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/archive"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/backend"
	"github.com/ankitects/anki/proto/anki/cards"
	ankicol "github.com/ankitects/anki/proto/anki/collection"
	"github.com/ankitects/anki/proto/anki/config"
	"github.com/ankitects/anki/proto/anki/generic"
	"github.com/ankitects/anki/proto/anki/import_export"
	"github.com/ankitects/anki/proto/anki/media"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"github.com/ankitects/anki/proto/anki/scheduler"
	"github.com/ankitects/anki/proto/anki/search"
	"github.com/ankitects/anki/proto/anki/stats"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/grpc"
)

// workerBinary identifies the built native server supplied by Bazel.
var workerBinary = flag.String("worker", "", "Anki worker runfile")

// worker starts a caller-owned backend and registers caller-owned shutdown.
func worker(t *testing.T) backend.Backend {
	t.Helper()
	binary, err := runfiles.Rlocation(*workerBinary)
	if err != nil {
		t.Fatal(err)
	}
	w := backend.New(backend.BackendSTDOptions{Binary: binary})
	if err := w.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := w.Stop(context.Background()); err != nil {
			t.Error(err)
		}
	})
	return w
}

// connection borrows the already-running test backend's service connection.
func connection(t *testing.T, w backend.Backend) grpc.ClientConnInterface {
	t.Helper()
	conn, err := w.Connection()
	if err != nil {
		t.Fatal(err)
	}
	return conn
}

// fixture records native identities and closed source paths for round trips.
type fixture struct {
	// database is the offline native collection.
	database string
	// archive is Anki's full collection package with media.
	archive string
	// note is the reviewed note identity.
	note int64
	// card is its generated, reviewed card identity.
	card int64
}

// createFixture uses native operations to create state, media, and a review.
func createFixture(t *testing.T) fixture {
	t.Helper()
	ctx := context.Background()
	root := t.TempDir()
	f := fixture{database: filepath.Join(root, "collection.anki2"), archive: filepath.Join(root, "source.colpkg")}
	w := worker(t)
	conn, err := w.Connection()
	if err != nil {
		t.Fatal(err)
	}
	collections := ankicol.NewBackendCollectionServiceClient(conn)
	if _, err := collections.OpenCollection(ctx, &ankicol.OpenCollectionRequest{CollectionPath: f.database, MediaFolderPath: f.database + ".media", MediaDbPath: f.database + ".media.db"}); err != nil {
		t.Fatal(err)
	}
	types := notetypes.NewNotetypesServiceClient(conn)
	nt, err := types.GetNotetypeIdByName(ctx, &generic.String{Val: "Basic"})
	if err != nil {
		t.Fatal(err)
	}
	added, err := notes.NewNotesServiceClient(conn).AddNote(ctx, &notes.AddNoteRequest{DeckId: 1, Note: &notes.Note{NotetypeId: nt.Ntid, Fields: []string{`<img src="example.txt">Front`, "Back"}, Tags: []string{"fixture"}}})
	if err != nil {
		t.Fatal(err)
	}
	f.note = added.NoteId
	found, err := search.NewSearchServiceClient(conn).SearchCards(ctx, &search.SearchRequest{})
	if err != nil || len(found.GetIds()) != 1 {
		t.Fatalf("generated card: %v %v", found, err)
	}
	f.card = found.Ids[0]
	scheduling := scheduler.NewSchedulerServiceClient(conn)
	states, err := scheduling.GetSchedulingStates(ctx, &cards.CardId{Cid: f.card})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := scheduling.AnswerCard(ctx, &scheduler.CardAnswer{CardId: f.card, CurrentState: states.Current, NewState: states.Good, Rating: scheduler.CardAnswer_GOOD, AnsweredAtMillis: time.Now().UnixMilli(), MillisecondsTaken: 250}); err != nil {
		t.Fatal(err)
	}
	if _, err := media.NewMediaServiceClient(conn).AddMediaFile(ctx, &media.AddMediaFileRequest{DesiredName: "example.txt", Data: []byte("public media fixture\n")}); err != nil {
		t.Fatal(err)
	}
	if _, err := config.NewConfigServiceClient(conn).SetConfigJson(ctx, &config.SetConfigJsonRequest{Key: "anki_as_code_state", ValueJson: []byte(`{"path":"decks/lesson/notes/greeting"}`)}); err != nil {
		t.Fatal(err)
	}
	if _, err := import_export.NewBackendImportExportServiceClient(conn).ExportCollectionPackage(ctx, &import_export.ExportCollectionPackageRequest{OutPath: f.archive, IncludeMedia: true}); err != nil {
		t.Fatal(err)
	}
	if err := w.Stop(ctx); err != nil {
		t.Fatal(err)
	}
	return f
}

// inspect verifies preserved identity, review, custom state, and media through Anki.
func inspect(t *testing.T, w backend.Backend, f fixture, want string, withMedia bool) {
	t.Helper()
	ctx := context.Background()
	conn, err := w.Connection()
	if err != nil {
		t.Fatal(err)
	}
	note, err := notes.NewNotesServiceClient(conn).GetNote(ctx, &notes.NoteId{Nid: f.note})
	if err != nil {
		t.Fatal(err)
	}
	if note.Fields[1] != want {
		t.Fatalf("note fields: %v", note.Fields)
	}
	history, err := stats.NewStatsServiceClient(conn).GetReviewLogs(ctx, &cards.CardId{Cid: f.card})
	if err != nil || len(history.GetEntries()) != 1 {
		t.Fatalf("review preservation: %v %v", history, err)
	}
	state, err := config.NewConfigServiceClient(conn).GetConfigJson(ctx, &generic.String{Val: "anki_as_code_state"})
	if err != nil || string(state.GetJson()) != `{"path":"decks/lesson/notes/greeting"}` {
		t.Fatalf("state preservation: %s %v", state, err)
	}
	if withMedia {
		path, err := media.NewMediaServiceClient(conn).GetAbsoluteMediaPath(ctx, &generic.String{Val: "example.txt"})
		if err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(path.Val)
		if err != nil || string(data) != "public media fixture\n" {
			t.Fatalf("media preservation: %q %v", data, err)
		}
	}
}

// assertAvailable checks scratch cleanup without retiring the caller-owned connection.
func assertAvailable(t *testing.T, w backend.Backend, scratch string) {
	t.Helper()
	if _, err := w.Connection(); err != nil {
		t.Fatalf("caller-owned backend was stopped: %v", err)
	}
	entries, err := os.ReadDir(scratch)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("scratch remains: %v", entries)
	}
}

// TestNativeArchiveWorkflow edits private resources and reimports native output.
func TestNativeArchiveWorkflow(t *testing.T) {
	f := createFixture(t)
	original, err := archive.FileHash(f.archive)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	scratch := t.TempDir()
	w := worker(t)
	source, err := archive.OpenPackage(ctx, connection(t, w), f.archive, scratch)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := source.Close(); err != nil {
			t.Error(err)
		}
	})
	inspect(t, w, f, "Back", true)
	conn, err := w.Connection()
	if err != nil {
		t.Fatal(err)
	}
	client := notes.NewNotesServiceClient(conn)
	note, err := client.GetNote(ctx, &notes.NoteId{Nid: f.note})
	if err != nil {
		t.Fatal(err)
	}
	note.Fields[1] = "Edited\n中文"
	if _, err := client.UpdateNotes(ctx, &notes.UpdateNotesRequest{Notes: []*notes.Note{note}}); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(t.TempDir(), "rebuilt.colpkg")
	if err := source.Publish(ctx, output, 0o640, ""); err != nil {
		t.Fatal(err)
	}
	assertAvailable(t, w, scratch)
	if err := source.Close(); err != nil {
		t.Fatal(err)
	}
	if err := source.Publish(ctx, output, 0o600, ""); err == nil {
		t.Fatal("consumed source reused")
	}
	info, err := os.Stat(output)
	if err != nil || info.Mode().Perm() != 0o640 {
		t.Fatalf("output permissions: %v %v", info, err)
	}
	after, err := archive.FileHash(f.archive)
	if err != nil || after != original {
		t.Fatalf("source changed: %v", err)
	}
	w2 := w
	rebuilt, err := archive.OpenSource(ctx, connection(t, w2), output, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := rebuilt.Close(); err != nil {
			t.Error(err)
		}
	}()
	inspect(t, w2, f, "Edited\n中文", true)
	if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
		for name, path := range map[string]string{"source.colpkg": f.archive, "rebuilt.colpkg": output} {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(artifacts, name), data, 0o600); err != nil {
				t.Fatal(err)
			}
		}
	}
}

// TestNativeSourceFailures covers raw copies, conflict guards, cancellation, and cleanup.
func TestNativeSourceFailures(t *testing.T) {
	f := createFixture(t)
	ctx := context.Background()
	t.Run("raw publication", func(t *testing.T) {
		w := worker(t)
		scratch := t.TempDir()
		source, err := archive.OpenSource(ctx, connection(t, w), f.database, scratch)
		if err != nil {
			t.Fatal(err)
		}
		inspect(t, w, f, "Back", false)
		output := filepath.Join(t.TempDir(), "copy.anki2")
		if err := source.Publish(ctx, output, 0o600, ""); err != nil {
			t.Fatal(err)
		}
		assertAvailable(t, w, scratch)
		w2 := worker(t)
		copied, err := archive.OpenSource(ctx, connection(t, w2), output, t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := copied.Close(); err != nil {
				t.Error(err)
			}
		}()
		inspect(t, w2, f, "Back", false)
	})
	for _, scenario := range []string{"changed target", "invalid destination", "cancelled publication", "closed source"} {
		t.Run(scenario, func(t *testing.T) {
			w := worker(t)
			scratch := t.TempDir()
			source, err := archive.OpenSource(ctx, connection(t, w), f.archive, scratch)
			if err != nil {
				t.Fatal(err)
			}
			output := filepath.Join(t.TempDir(), "target.colpkg")
			if err := os.WriteFile(output, []byte("existing target"), 0o600); err != nil {
				t.Fatal(err)
			}
			hash, err := archive.FileHash(output)
			if err != nil {
				t.Fatal(err)
			}
			callCtx := ctx
			switch scenario {
			case "changed target":
				hash = "changed"
			case "invalid destination":
				output = filepath.Join(output, "invalid")
			case "cancelled publication":
				var cancel context.CancelFunc
				callCtx, cancel = context.WithCancel(ctx)
				cancel()
			case "closed source":
				if err := source.Close(); err != nil {
					t.Fatal(err)
				}
			}
			if err := source.Publish(callCtx, output, 0o600, hash); err == nil {
				t.Fatal("expected failure")
			}
			assertAvailable(t, w, scratch)
			if err := source.Close(); err != nil {
				t.Fatal(err)
			}
			if scenario != "invalid destination" {
				data, err := os.ReadFile(output)
				if err != nil || string(data) != "existing target" {
					t.Fatalf("target changed: %q %v", data, err)
				}
			}
			matches, err := filepath.Glob(filepath.Join(filepath.Dir(output), ".anki-output-*"))
			if err != nil || len(matches) != 0 {
				t.Fatalf("staging remains: %v %v", matches, err)
			}
		})
	}
	for _, scenario := range []string{"invalid package", "cancelled open", "-wal", "-journal", "raw package"} {
		t.Run(scenario, func(t *testing.T) {
			w := worker(t)
			scratch := t.TempDir()
			input := f.archive
			callCtx := ctx
			if scenario == "invalid package" {
				input = filepath.Join(t.TempDir(), "bad.colpkg")
				if err := os.WriteFile(input, []byte("invalid package fixture"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			if scenario == "cancelled open" {
				var cancel context.CancelFunc
				callCtx, cancel = context.WithCancel(ctx)
				cancel()
			}
			if scenario == "-wal" || scenario == "-journal" {
				input = filepath.Join(t.TempDir(), "sidecar.anki2")
				data, err := os.ReadFile(f.database)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(input, data, 0o600); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(input+scenario, nil, 0o600); err != nil {
					t.Fatal(err)
				}
			}
			var err error
			if scenario == "raw package" {
				_, err = archive.OpenPackage(callCtx, connection(t, w), f.database, scratch)
			} else {
				_, err = archive.OpenSource(callCtx, connection(t, w), input, scratch)
			}
			if err == nil {
				t.Fatal("expected open failure")
			}
			if scenario == "cancelled open" && !errors.Is(err, context.Canceled) {
				t.Fatalf("cancellation cause: %v", err)
			}
			assertAvailable(t, w, scratch)
		})
	}
}

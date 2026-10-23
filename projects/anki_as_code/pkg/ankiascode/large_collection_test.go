package ankiascode

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"
	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"github.com/ankitects/anki/proto/anki/cards"
	"github.com/ankitects/anki/proto/anki/generic"
	"github.com/ankitects/anki/proto/anki/media"
	"github.com/ankitects/anki/proto/anki/notes"
	"google.golang.org/protobuf/proto"
)

// TestLargeNoteFieldWorkflow verifies a large native protobuf response and clean command streams.
func TestLargeNoteFieldWorkflow(t *testing.T) {
	fixture := newFixture(t)
	root := t.TempDir()
	text := filepath.Join(root, "text")
	if err := prepareFixtureConfigs(t, context.Background(), fixture.base, text); err != nil {
		t.Fatal(err)
	}
	path := fixtureNotePath(text)
	note := &ankiascodepb.NoteConfig{}
	readFixtureTOML(t, path, note)
	note.Fields["hanzi"] = strings.Repeat("x", 17<<20)
	writeFixtureTOML(t, path, note)
	result := filepath.Join(root, "large-field.colpkg")
	runFixtureCLI(t, "col", "build", "--base", fixture.base, "--text", text, "--output", result)
	client := inspectFixture(t, result)
	actual, err := client.notes.GetNote(context.Background(), &notes.NoteId{Nid: fixture.noteID})
	if err != nil || len(actual.GetFields()[0]) != 17<<20 {
		t.Fatalf("large field truncated: %v", err)
	}
	if native := os.Getenv("NATIVE_ANKI_CLI"); native != "" {
		binary := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"), native)
		verified := inspectFixtureWithBinary(t, result, binary)
		assertFixtureMedia(t, verified)
		note, err := verified.notes.GetNote(t.Context(), &notes.NoteId{Nid: fixture.noteID})
		if err != nil || len(note.GetFields()[0]) != 17<<20 {
			t.Fatalf("native Anki rejected WASM package payload/identity: %v", err)
		}
		for _, id := range fixture.cardIDs {
			request := &cards.CardId{Cid: id}
			before, err := client.cards.GetCard(t.Context(), request)
			if err != nil {
				t.Fatalf("read WASM card %d: %v", id, err)
			}
			after, err := verified.cards.GetCard(t.Context(), request)
			if err != nil || !proto.Equal(before, after) {
				t.Fatalf("native reopen changed card %d scheduling or identity: %v", id, err)
			}
			history, err := client.stats.GetReviewLogs(t.Context(), request)
			if err != nil {
				t.Fatalf("read WASM review history for card %d: %v", id, err)
			}
			reopened, err := verified.stats.GetReviewLogs(t.Context(), request)
			if err != nil || !proto.Equal(history, reopened) {
				t.Fatalf("native reopen changed card %d review history: %v", id, err)
			}
		}
	}
	saveFixtureArtifact(t, result, "large-field-roundtrip.colpkg")
	output := runFixtureCLI(t, "col", "plan", "--collection", result, "--text", text)
	plan := &ankiascodepb.CollectionChangePlan{}
	if err := codec.Decode(output, codec.JSON, plan); err != nil {
		t.Fatal(err)
	}
	if len(plan.Changes) != 0 {
		t.Fatalf("large-field CLI plan did not converge: %v", plan.Changes)
	}
	if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
		if err := os.WriteFile(filepath.Join(artifacts, "large-api-plan.json"), output, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestLargeCollectionWorkflow verifies native package media beyond 512 MiB.
func TestLargeCollectionWorkflow(t *testing.T) {
	fixture := newFixture(t)
	root := t.TempDir()
	largeMedia := filepath.Join(filepath.Dir(fixture.database), "collection.media", "large.bin")
	file, err := os.Create(largeMedia)
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.Truncate(largeMedia, 513<<20); err != nil {
		t.Fatal(err)
	}
	text := filepath.Join(root, "text")
	// Prepare bindings on the database so setup does not repeatedly publish large media.
	if err := prepareFixtureConfigs(t, context.Background(), fixture.database, text); err != nil {
		t.Fatal(err)
	}
	base := filepath.Join(root, "large.colpkg")
	packageFixture(t, fixture.database, base)
	result := filepath.Join(root, "large-rebuilt.colpkg")
	if err := buildForFixture(context.Background(), base, text, result, ""); err != nil {
		t.Fatal(err)
	}
	client := inspectFixture(t, result)
	conn, err := client.worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	// Borrow the restored connection so observing convergence cannot repackage the archive.
	plan, err := fixtureApplication().ColPlan(t.Context(), ankiascodeapi.CollectionOptions{Connection: conn}, text)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Changes) != 0 {
		t.Fatalf("large-media build did not converge: %v", plan.Changes)
	}
	path, err := media.NewMediaServiceClient(conn).GetAbsoluteMediaPath(context.Background(), &generic.String{Val: "large.bin"})
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path.Val)
	if err != nil || info.Size() != 513<<20 {
		t.Fatalf("large native media: %v %v", info, err)
	}
	saveFixtureArtifact(t, result, "large-api-roundtrip.colpkg")
}

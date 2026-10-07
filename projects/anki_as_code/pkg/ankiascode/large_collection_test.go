package ankiascode

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"github.com/ankitects/anki/proto/anki/generic"
	"github.com/ankitects/anki/proto/anki/media"
	"github.com/ankitects/anki/proto/anki/notes"
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
	note := &collectionpb.NoteConfig{}
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
	assertConverged(t, result, text)
	output := runFixtureCLI(t, "col", "plan", "--collection", result, "--text", text)
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
	base := filepath.Join(root, "large.colpkg")
	packageFixture(t, fixture.database, base)
	text := filepath.Join(root, "text")
	if err := prepareFixtureConfigs(t, context.Background(), base, text); err != nil {
		t.Fatal(err)
	}
	result := filepath.Join(root, "large-rebuilt.colpkg")
	if err := buildForFixture(context.Background(), base, text, result, ""); err != nil {
		t.Fatal(err)
	}
	assertConverged(t, result, text)
	client := inspectFixture(t, result)
	conn, err := client.worker.Connection()
	if err != nil {
		t.Fatal(err)
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

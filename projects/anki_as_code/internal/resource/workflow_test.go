package resource

import (
	"cmp"
	"context"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"github.com/ankitects/anki/proto/anki/cards"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/ankitects/anki/proto/anki/import_export"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"github.com/ankitects/anki/proto/anki/search"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// TestResourceWorkflow verifies shared native access without policy or lifecycle ownership.
func TestResourceWorkflow(t *testing.T) {
	ctx := context.Background()
	binary, err := runfiles.Rlocation("_main/projects/anki_grpc/cmd/anki_grpc/anki_grpc")
	if err != nil {
		t.Fatal(err)
	}
	worker := ankiclient.NewSTD(ankiclient.ClientSTDOptions{Binary: binary, Arguments: []string{"serve"}})
	collection := ankiclient.NewCollection(ankiclient.ClientCollectionOptions{Backend: worker, Path: filepath.Join(t.TempDir(), "collection.anki2")})
	t.Cleanup(func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := collection.Stop(cleanup); err != nil {
			t.Error(err)
		}
		if err := worker.Stop(cleanup); err != nil {
			t.Error(err)
		}
	})
	if err := worker.Start(ctx); err != nil {
		t.Fatal(err)
	}
	if err := collection.Start(ctx); err != nil {
		t.Fatal(err)
	}
	connection, err := collection.Connection()
	if err != nil {
		t.Fatal(err)
	}
	store := New(connection)
	created, err := store.NoteTypes.AddNotetype(ctx, &notetypes.Notetype{
		Name: "Shared", Config: &notetypes.Notetype_Config{Css: ".card { color: blue; }"},
		Fields:    []*notetypes.Notetype_Field{{Name: "Front", Config: &notetypes.Notetype_Field_Config{}}, {Name: "Back", Config: &notetypes.Notetype_Field_Config{}}},
		Templates: []*notetypes.Notetype_Template{{Name: "Card", Config: &notetypes.Notetype_Template_Config{QFormat: "{{Front}}", AFormat: "{{Back}}"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	added, err := store.Notes.AddNote(ctx, &notes.AddNoteRequest{DeckId: 1, Note: &notes.Note{NotetypeId: created.Id, Fields: []string{"<b>世界</b>", "line one\nline two"}}})
	if err != nil {
		t.Fatal(err)
	}
	filtered, err := store.Decks.AddDeck(ctx, &decks.Deck{Name: "Filtered", Kind: &decks.Deck_Filtered_{Filtered: &decks.Deck_Filtered{Reschedule: true}}})
	if err != nil {
		t.Fatal(err)
	}
	ids, err := store.Search.SearchCards(ctx, &search.SearchRequest{})
	if err != nil || len(ids.GetIds()) != 1 {
		t.Fatalf("native card generation: %v", err)
	}
	card, err := store.Cards.GetCard(ctx, &cards.CardId{Cid: ids.Ids[0]})
	if err != nil {
		t.Fatal(err)
	}
	card.OriginalDeckId, card.DeckId = 1, filtered.Id
	if _, err := store.Cards.UpdateCards(ctx, &cards.UpdateCardsRequest{Cards: []*cards.Card{card}, SkipUndoEntry: true}); err != nil {
		t.Fatal(err)
	}
	card, err = store.Cards.GetCard(ctx, &cards.CardId{Cid: card.Id})
	if err != nil {
		t.Fatal(err)
	}
	native, err := store.Notes.GetNote(ctx, &notes.NoteId{Nid: added.NoteId})
	if err != nil {
		t.Fatal(err)
	}
	first, err := store.Snapshot(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Notes) != 1 || !proto.Equal(first.Notes[0], native) || len(first.Cards) != 1 || !proto.Equal(first.Cards[0], card) || first.Presets[1] != "Default" {
		t.Fatal("snapshot changed or omitted native values, filtered-card references or presets")
	}
	if !slices.IsSortedFunc(first.Decks, func(a, b *decks.Deck) int { return cmp.Compare(a.Name, b.Name) }) || !slices.IsSortedFunc(first.NoteTypes, func(a, b *notetypes.Notetype) int { return cmp.Compare(a.Id, b.Id) }) {
		t.Fatal("snapshot ordering is unstable")
	}
	native.Fields[0] = "Updated"
	if _, err := store.Notes.UpdateNotes(ctx, &notes.UpdateNotesRequest{Notes: []*notes.Note{native}}); err != nil {
		t.Fatal(err)
	}
	second, err := store.Snapshot(ctx)
	if err != nil || second.Notes[0].Fields[0] != "Updated" || first.Notes[0].Fields[0] != "<b>世界</b>" {
		t.Fatalf("snapshot was cached or shared between reads: %v", err)
	}
	first.Notes[0].Fields[0] = "Caller mutation"
	unchanged, err := store.Notes.GetNote(ctx, &notes.NoteId{Nid: added.NoteId})
	if err != nil || !proto.Equal(unchanged, second.Notes[0]) {
		t.Fatalf("snapshot mutation reached native storage: %v", err)
	}
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if value, err := store.Snapshot(cancelled); err == nil || value != nil {
		t.Fatal("cancelled snapshot succeeded or exposed partial data")
	}
	if _, err := store.Cards.GetCard(ctx, &cards.CardId{Cid: card.Id}); err != nil {
		t.Fatalf("snapshot closed borrowed collection: %v", err)
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		data, err := (protojson.MarshalOptions{Indent: "  "}).Marshal(second.Notes[0])
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "snapshot-note.json"), data, 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := import_export.NewBackendImportExportServiceClient(connection).ExportCollectionPackage(ctx, &import_export.ExportCollectionPackageRequest{OutPath: filepath.Join(output, "resources.colpkg")}); err != nil {
			t.Fatal(err)
		}
	}
}

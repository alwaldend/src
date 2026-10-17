package ankiascode

import (
	"context"
	"path/filepath"
	"testing"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"github.com/ankitects/anki/proto/anki/cards"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"google.golang.org/protobuf/proto"
)

// TestNoteTypeStructureWorkflow verifies field/template migrations and retained reviewed cards through Anki.
func TestNoteTypeStructureWorkflow(t *testing.T) {
	fixture := newFixture(t)
	root := t.TempDir()
	text := filepath.Join(root, "text")
	if err := prepareFixtureConfigs(t, context.Background(), fixture.base, text); err != nil {
		t.Fatal(err)
	}
	client := inspectFixture(t, fixture.base)
	previous, err := client.cards.GetCard(context.Background(), &cards.CardId{Cid: fixture.cardIDs[0]})
	if err != nil {
		t.Fatal(err)
	}
	history, err := client.stats.GetReviewLogs(context.Background(), &cards.CardId{Cid: previous.Id})
	if err != nil {
		t.Fatal(err)
	}
	if err := client.close(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(text, "note_types", "Chinese.notetype.anki.toml")
	settings := &collectionpb.NoteTypeConfig{}
	readFixtureTOML(t, path, settings)
	settings.Fields[0].Name = "Character"
	settings.Config.Css = ".card { color: blue; }"
	settings.Templates[0].Config.QFormat = "{{Character}}"
	settings.Templates[1].Config.AFormat = "{{Character}}"
	settings.Fields[0], settings.Fields[2] = settings.Fields[2], settings.Fields[0]
	settings.Templates[0], settings.Templates[1] = settings.Templates[1], settings.Templates[0]
	writeFixtureTOML(t, filepath.Join(text, "collection.anki.toml"), &collectionpb.CollectionConfig{Version: 1, Moved: []*collectionpb.ResourceMove{{From: "note_types/Chinese.notetype.anki.toml#fields/hanzi", To: "note_types/Chinese.notetype.anki.toml#fields/Character"}}})
	writeFixtureTOML(t, path, settings)
	notePath := fixtureNotePath(text)
	note := &collectionpb.NoteConfig{}
	readFixtureTOML(t, notePath, note)
	note.Fields["Character"] = note.Fields["hanzi"]
	delete(note.Fields, "hanzi")
	for _, card := range note.Cards {
		card.Ordinal = 1 - card.Ordinal
	}
	writeFixtureTOML(t, notePath, note)
	plan := fixturePlan(t, fixture.base, text)
	if plan.NoteTypesChanged != 1 {
		t.Fatalf("schema change absent: %v", plan)
	}
	result := filepath.Join(root, "structured.colpkg")
	if err := buildForFixture(context.Background(), fixture.base, text, result, ""); err != nil {
		t.Fatal(err)
	}
	assertConverged(t, result, text)
	actual := inspectFixture(t, result)
	definition, err := actual.types.GetNotetype(context.Background(), &notetypes.NotetypeId{Ntid: fixture.typeID})
	if err != nil || definition.GetFields()[0].Name != "english" || definition.GetFields()[2].Name != "Character" {
		t.Fatalf("field migration: %v", err)
	}
	changed, err := actual.notes.GetNote(context.Background(), &notes.NoteId{Nid: fixture.noteID})
	if err != nil || changed.GetFields()[2] != "<b>你好</b>" {
		t.Fatalf("field content migration: %v", err)
	}
	retained, err := actual.cards.GetCard(context.Background(), &cards.CardId{Cid: previous.Id})
	if err != nil || !proto.Equal(fixtureCardSchedule(previous), fixtureCardSchedule(retained)) {
		t.Fatalf("retained scheduling changed: %v", err)
	}
	logs, err := actual.stats.GetReviewLogs(context.Background(), &cards.CardId{Cid: previous.Id})
	if err != nil || !proto.Equal(history, logs) {
		t.Fatalf("history changed: %v", err)
	}
	saveFixtureArtifact(t, result, "api-structure.colpkg")
	saveFixtureArtifact(t, path, "api-structure.notetype.anki.toml")
}

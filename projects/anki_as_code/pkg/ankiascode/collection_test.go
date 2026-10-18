package ankiascode

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"github.com/ankitects/anki/proto/anki/cards"
	"github.com/ankitects/anki/proto/anki/deck_config"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/ankitects/anki/proto/anki/generic"
	"github.com/ankitects/anki/proto/anki/notes"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// TestReadOnlyPlanWorkflow verifies a changing plan reads the supplied database without mutation or copies.
func TestReadOnlyPlanWorkflow(t *testing.T) {
	fixture := newFixture(t)
	ctx := context.Background()
	text := filepath.Join(t.TempDir(), "configs")
	if err := prepareFixtureConfigs(t, ctx, fixture.database, text); err != nil {
		t.Fatal(err)
	}
	path := fixtureNotePath(text)
	note := &ankiascodepb.NoteConfig{}
	readFixtureTOML(t, path, note)
	note.Fields["hanzi"] = "planned content"
	writeFixtureTOML(t, path, note)
	writeFixtureTOML(t, filepath.Join(filepath.Dir(path), "new.note.anki.toml"), note)
	before, err := os.ReadFile(fixture.database)
	if err != nil {
		t.Fatal(err)
	}
	bindings := stateFixtureDump(t, fixture.database)
	scratch := t.TempDir()
	t.Setenv("TMPDIR", scratch)
	first := fixturePlan(t, fixture.database, text)
	second := fixturePlan(t, fixture.database, text)
	if first.NotesAdded != 1 || first.NotesChanged != 1 || !first.CardGenerationDeferred || !proto.Equal(first, second) {
		t.Fatalf("read-only plan did not describe the pending changes: %v", first)
	}
	after, err := os.ReadFile(fixture.database)
	if err != nil || string(before) != string(after) || !proto.Equal(bindings, stateFixtureDump(t, fixture.database)) {
		t.Fatal("planning changed the database or stored bindings")
	}
	files, err := os.ReadDir(scratch)
	if err != nil || len(files) != 0 {
		t.Fatalf("planning created a temporary collection: %v %v", files, err)
	}
	client := connectFixture(t, fixture.database)
	actual, err := client.notes.GetNote(ctx, &notes.NoteId{Nid: fixture.noteID})
	if err != nil || actual.GetFields()[0] != "<b>你好</b>" {
		t.Fatal("planning applied the requested field edit")
	}
	if err := client.close(); err != nil {
		t.Fatal(err)
	}
	if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
		data, err := (protojson.MarshalOptions{Indent: "  "}).Marshal(first)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(artifacts, "read-only-plan.json"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := applyForFixture(ctx, fixture.database, text, ""); err != nil {
		t.Fatal(err)
	}
	assertConverged(t, fixture.database, text)
	saveFixtureArtifact(t, fixture.database, "planned-and-applied.anki2")
}

// TestCollectionWorkflow verifies CLI edits, reviewed-card retention, media, and convergence.
func TestCollectionWorkflow(t *testing.T) {
	fixture := newFixture(t)
	root := t.TempDir()
	text := filepath.Join(root, "text")
	if err := prepareFixtureConfigs(t, context.Background(), fixture.base, text); err != nil {
		t.Fatal(err)
	}
	before := inspectFixture(t, fixture.base)
	card, err := before.cards.GetCard(context.Background(), &cards.CardId{Cid: fixture.cardIDs[0]})
	if err != nil {
		t.Fatal(err)
	}
	history, err := before.stats.GetReviewLogs(context.Background(), &cards.CardId{Cid: card.Id})
	if err != nil {
		t.Fatal(err)
	}
	preset, err := before.presets.GetDeckConfig(context.Background(), &deck_config.DeckConfigId{Dcid: 1})
	if err != nil {
		t.Fatal(err)
	}
	if err := before.close(); err != nil {
		t.Fatal(err)
	}
	notePath := fixtureNotePath(text)
	note := &ankiascodepb.NoteConfig{}
	readFixtureTOML(t, notePath, note)
	note.Fields["hanzi"] = "<strong>你好，世界</strong>\nSecond\x00\x1f line"
	note.Tags = []string{" edited ", "edited"}
	writeFixtureTOML(t, notePath, note)
	marker := filepath.Join(text, "decks", "lesson", "deck.anki.toml")
	writeFixtureTOML(t, marker, &ankiascodepb.DeckConfig{Normal: &decks.Deck_Normal{}})
	target := filepath.Join(filepath.Dir(marker), "notes", "greeting.note.anki.toml")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(notePath, target); err != nil {
		t.Fatal(err)
	}
	note.Cards = nil
	writeFixtureTOML(t, target, note)
	writeFixtureTOML(t, filepath.Join(text, "collection.anki.toml"), &ankiascodepb.CollectionConfig{Version: 1, Moved: []*ankiascodepb.ResourceMove{{From: "decks/grammar/notes/你好.note.anki.toml", To: "decks/lesson/notes/greeting.note.anki.toml"}}})
	plan := fixturePlan(t, fixture.base, text)
	if plan.NotesChanged != 1 || plan.CardsChanged != 2 || plan.DecksAdded != 1 {
		t.Fatalf("unexpected edits: %v", plan)
	}
	result := filepath.Join(root, "edited.colpkg")
	runFixtureCLI(t, "col", "build", "--base", fixture.base, "--text", text, "--output", result)
	if pending := fixturePlan(t, result, text); pending.NotesChanged != 1 || !pending.CardGenerationDeferred {
		t.Fatal("plan hid requested values normalized by Anki during apply")
	}
	actual := inspectFixture(t, result)
	changed, err := actual.notes.GetNote(context.Background(), &notes.NoteId{Nid: fixture.noteID})
	if err != nil || changed.GetFields()[0] != "<strong>你好，世界</strong>\nSecond line" || !slices.Equal(changed.Tags, []string{"edited"}) {
		t.Fatalf("note edit: %v", err)
	}
	recorded := &ankiascodepb.NoteConfig{}
	readFixtureTOML(t, target, recorded)
	if recorded.Fields["hanzi"] != note.Fields["hanzi"] || !slices.Equal(recorded.Tags, note.Tags) {
		t.Fatal("reconciliation rewrote source note content")
	}
	retained, err := actual.cards.GetCard(context.Background(), &cards.CardId{Cid: card.Id})
	if err != nil || !proto.Equal(fixtureCardSchedule(card), fixtureCardSchedule(retained)) {
		t.Fatalf("reviewed card state changed: %v", err)
	}
	logs, err := actual.stats.GetReviewLogs(context.Background(), &cards.CardId{Cid: card.Id})
	if err != nil || !proto.Equal(history, logs) {
		t.Fatalf("review history changed: %v", err)
	}
	currentPreset, err := actual.presets.GetDeckConfig(context.Background(), &deck_config.DeckConfigId{Dcid: 1})
	if err != nil || !proto.Equal(preset, currentPreset) {
		t.Fatalf("scheduling preset changed: %v", err)
	}
	if err := actual.close(); err != nil {
		t.Fatal(err)
	}
	runFixtureCLI(t, "col", "apply", "--collection", result, "--text", text)
	if pending := fixturePlan(t, result, text); pending.NotesChanged != 1 || !pending.CardGenerationDeferred {
		t.Fatal("plan hid requested values normalized by Anki during apply")
	}
	mediaClient := inspectFixture(t, result)
	assertFixtureMedia(t, mediaClient)
	if err := mediaClient.close(); err != nil {
		t.Fatal(err)
	}
	saveFixtureArtifact(t, result, "api-reconciled.colpkg")
	saveFixtureArtifact(t, target, "api-edited.note.anki.toml")
}

// TestRejectedWorkflow verifies declaration and native-template errors leave inputs intact.
func TestRejectedWorkflow(t *testing.T) {
	for _, scenario := range []string{"malformed", "missing-type", "template", "missing-field", "native-guid"} {
		t.Run(scenario, func(t *testing.T) {
			fixture := newFixture(t)
			text := filepath.Join(t.TempDir(), "text")
			if err := prepareFixtureConfigs(t, context.Background(), fixture.base, text); err != nil {
				t.Fatal(err)
			}
			path := fixtureNotePath(text)
			switch scenario {
			case "malformed":
				if err := os.WriteFile(path, []byte("invalid = ["), 0o600); err != nil {
					t.Fatal(err)
				}
			case "template":
				settingsPath := filepath.Join(text, "note_types", "Chinese.notetype.anki.toml")
				settings := &ankiascodepb.NoteTypeConfig{}
				readFixtureTOML(t, settingsPath, settings)
				settings.Templates[0].Config.QFormat = "{{UnknownField}}"
				writeFixtureTOML(t, settingsPath, settings)
			case "missing-type":
				note := &ankiascodepb.NoteConfig{}
				readFixtureTOML(t, path, note)
				note.NoteType = "note_types/missing.notetype.anki.toml"
				writeFixtureTOML(t, path, note)
			case "native-guid":
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, append([]byte("guid = 'old'\n"), data...), 0o600); err != nil {
					t.Fatal(err)
				}

			case "missing-field":
				note := &ankiascodepb.NoteConfig{}
				readFixtureTOML(t, path, note)
				delete(note.Fields, "hanzi")
				writeFixtureTOML(t, path, note)
			}
			before, err := os.ReadFile(fixture.base)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := planForFixture(context.Background(), fixture.base, text, ""); (err == nil) != (scenario == "template") {
				t.Fatalf("unexpected declaration validation result: %v", err)
			}
			if err := applyForFixture(context.Background(), fixture.base, text, ""); err == nil {
				t.Fatal("invalid apply accepted")
			}
			if err := buildForFixture(context.Background(), fixture.base, text, filepath.Join(text, "rejected.colpkg"), ""); err == nil {
				t.Fatal("invalid build accepted")
			}
			after, err := os.ReadFile(fixture.base)
			if err != nil || string(before) != string(after) {
				t.Fatal("rejected operation changed the archive")
			}
		})
	}
}

// TestRetainedNoteDeckRemovalWorkflow rejects unsafe deletion before writes and accepts explicit destinations.
func TestRetainedNoteDeckRemovalWorkflow(t *testing.T) {
	for _, format := range []string{"database", "package"} {
		t.Run(format, func(t *testing.T) {
			fixture := newFixture(t)
			ctx := t.Context()
			client := connectFixture(t, fixture.database)
			temporary, err := client.decks.NewDeck(ctx, &generic.Empty{})
			if err != nil {
				t.Fatal(err)
			}
			temporary.Name = "Temporary"
			created, err := client.decks.AddDeck(ctx, temporary)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := client.cards.SetDeck(ctx, &cards.SetDeckRequest{CardIds: fixture.cardIDs, DeckId: created.Id}); err != nil {
				t.Fatal(err)
			}
			history, err := client.stats.GetReviewLogs(ctx, &cards.CardId{Cid: fixture.cardIDs[0]})
			if err != nil {
				t.Fatal(err)
			}
			if err := client.close(); err != nil {
				t.Fatal(err)
			}
			packageFixture(t, fixture.database, fixture.base)
			input := fixture.database
			if format == "package" {
				input = fixture.base
			}
			text := filepath.Join(t.TempDir(), "configs")
			if err := prepareFixtureConfigs(t, ctx, input, text); err != nil {
				t.Fatal(err)
			}
			source := filepath.Join(text, "decks", "Temporary", "notes", "你好.note.anki.toml")
			note := &ankiascodepb.NoteConfig{}
			readFixtureTOML(t, source, note)
			note.Cards = nil
			destination := filepath.Join(text, "notes", "retained.note.anki.toml")
			writeFixtureTOML(t, destination, note)
			if err := os.RemoveAll(filepath.Join(text, "decks", "Temporary")); err != nil {
				t.Fatal(err)
			}
			writeFixtureTOML(t, filepath.Join(text, "collection.anki.toml"), &ankiascodepb.CollectionConfig{Version: 1, Moved: []*ankiascodepb.ResourceMove{
				{From: "decks/Temporary/notes/你好.note.anki.toml", To: "notes/retained.note.anki.toml"},
			}})
			before, err := os.ReadFile(input)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := planForFixture(ctx, input, text, ""); err == nil {
				t.Error("plan accepted removing a retained card's deck without a destination")
			}
			if err := applyForFixture(ctx, input, text, ""); err == nil {
				t.Error("apply accepted removing a retained card's deck without a destination")
			}
			after, err := os.ReadFile(input)
			if err != nil || string(before) != string(after) {
				t.Fatal("unsafe deck deletion changed collection bytes")
			}
			note.Cards = []*ankiascodepb.CardConfig{
				{Ordinal: 0, Deck: proto.String("decks/grammar/deck.anki.toml")},
				{Ordinal: 1, Deck: proto.String("decks/grammar/deck.anki.toml")},
			}
			writeFixtureTOML(t, destination, note)
			if err := applyForFixture(ctx, input, text, ""); err != nil {
				t.Fatal(err)
			}
			assertConverged(t, input, text)
			if format == "package" {
				client = inspectFixture(t, input)
			} else {
				client = connectFixture(t, input)
			}
			retained, err := client.notes.GetNote(ctx, &notes.NoteId{Nid: fixture.noteID})
			if err != nil || retained.GetGuid() != "fixture-note" {
				t.Fatalf("retained note disappeared: %v", err)
			}
			for _, id := range fixture.cardIDs {
				card, err := client.cards.GetCard(ctx, &cards.CardId{Cid: id})
				if err != nil || card.GetDeckId() != 1 {
					t.Fatalf("retained card did not move to the surviving deck: %v", err)
				}
			}
			retainedHistory, err := client.stats.GetReviewLogs(ctx, &cards.CardId{Cid: fixture.cardIDs[0]})
			if err != nil || !proto.Equal(history, retainedHistory) {
				t.Fatalf("retained review history changed: %v", err)
			}
			if err := client.close(); err != nil {
				t.Fatal(err)
			}
			saveFixtureArtifact(t, input, "retained-note-deck-removal-"+format+filepath.Ext(input))
			saveFixtureArtifact(t, destination, "retained-note-"+format+".note.anki.toml")
		})
	}
}

// TestConvergedDatabaseApplyWorkflow verifies repeated empty-plan applies do not write database bytes.
func TestConvergedDatabaseApplyWorkflow(t *testing.T) {
	fixture := newFixture(t)
	text := filepath.Join(t.TempDir(), "configs")
	if err := prepareFixtureConfigs(t, t.Context(), fixture.database, text); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(fixture.database)
	if err != nil {
		t.Fatal(err)
	}
	for range 2 {
		assertConverged(t, fixture.database, text)
		if err := applyForFixture(t.Context(), fixture.database, text, ""); err != nil {
			t.Fatal(err)
		}
		after, err := os.ReadFile(fixture.database)
		if err != nil || string(before) != string(after) {
			t.Fatal("converged apply changed database bytes")
		}
	}
	saveFixtureArtifact(t, fixture.database, "converged-apply.anki2")
}

// TestCreatedNoteBatchWorkflow verifies each new declaration receives its own native identity and deck assignment.
func TestCreatedNoteBatchWorkflow(t *testing.T) {
	fixture := newFixture(t)
	ctx := t.Context()
	text := filepath.Join(t.TempDir(), "configs")
	if err := prepareFixtureConfigs(t, ctx, fixture.base, text); err != nil {
		t.Fatal(err)
	}
	original := &ankiascodepb.NoteConfig{}
	readFixtureTOML(t, fixtureNotePath(text), original)
	writeFixtureTOML(t, filepath.Join(text, "decks", "batch", "deck.anki.toml"), &ankiascodepb.DeckConfig{Normal: &decks.Deck_Normal{}})
	labels, homes := map[string]string{}, map[string]string{}
	for index := range 32 {
		home := "grammar"
		if index%2 != 0 {
			home = "batch"
		}
		address := fmt.Sprintf("decks/%s/notes/created-%02d.note.anki.toml", home, index)
		note := proto.Clone(original).(*ankiascodepb.NoteConfig)
		label := fmt.Sprintf("Batch note %02d", index)
		note.Fields["hanzi"] = label
		note.Cards = nil
		writeFixtureTOML(t, filepath.Join(text, filepath.FromSlash(address)), note)
		labels[address] = label
		homes[address] = "decks/" + home + "/deck.anki.toml"
	}
	sources := fixtureSources(t, text)
	plan := fixturePlan(t, fixture.base, text)
	if plan.NotesAdded != int64(len(labels)) || !plan.CardGenerationDeferred {
		t.Fatalf("incorrect batch creation plan: %v", plan)
	}
	output := filepath.Join(t.TempDir(), "created-notes.colpkg")
	runFixtureCLI(t, "col", "build", "--base", fixture.base, "--text", text, "--output", output)
	assertFixtureSources(t, text, sources)
	assertConverged(t, output, text)
	bindings := stateFixtureDump(t, output)
	client := inspectFixture(t, output)
	seen := map[int64]bool{fixture.noteID: true}
	for address, label := range labels {
		id := bindings.GetResources().GetNotes()[address].GetId()
		if id == 0 || seen[id] {
			t.Fatalf("batch path %q received duplicate or missing identity %d", address, id)
		}
		seen[id] = true
		note, err := client.notes.GetNote(ctx, &notes.NoteId{Nid: id})
		if err != nil || len(note.GetFields()) == 0 || note.Fields[0] != label {
			t.Fatalf("batch identity mapped to incorrect content at %q: %v", address, err)
		}
		generated, err := client.notes.CardsOfNote(ctx, &notes.NoteId{Nid: id})
		if err != nil || len(generated.GetCids()) != 2 {
			t.Fatalf("batch note %q generated incorrect cards: %v", address, err)
		}
		home := bindings.GetResources().GetDecks()[homes[address]].GetId()
		for _, cardID := range generated.Cids {
			card, err := client.cards.GetCard(ctx, &cards.CardId{Cid: cardID})
			if err != nil || card.GetDeckId() != home {
				t.Fatalf("batch note %q has incorrect card destination: %v", address, err)
			}
		}
	}
	if err := client.close(); err != nil {
		t.Fatal(err)
	}
	saveFixtureArtifact(t, output, "created-note-batch.colpkg")
}

// TestDeckKindConversionWorkflow protects occupied decks and permits empty conversions without losing reviews.
func TestDeckKindConversionWorkflow(t *testing.T) {
	for _, format := range []string{"database", "package"} {
		t.Run(format, func(t *testing.T) {
			fixture := newFixture(t)
			ctx := t.Context()
			client := connectFixture(t, fixture.database)
			created, err := client.decks.AddDeck(ctx, &decks.Deck{Name: "Convertible", Kind: &decks.Deck_Normal_{Normal: &decks.Deck_Normal{ConfigId: 1}}})
			if err != nil {
				t.Fatal(err)
			}
			if _, err := client.cards.SetDeck(ctx, &cards.SetDeckRequest{CardIds: fixture.cardIDs, DeckId: created.Id}); err != nil {
				t.Fatal(err)
			}
			history, err := client.stats.GetReviewLogs(ctx, &cards.CardId{Cid: fixture.cardIDs[0]})
			if err != nil {
				t.Fatal(err)
			}
			if err := client.close(); err != nil {
				t.Fatal(err)
			}
			packageFixture(t, fixture.database, fixture.base)
			input := fixture.database
			if format == "package" {
				input = fixture.base
			}
			text := filepath.Join(t.TempDir(), "configs")
			if err := prepareFixtureConfigs(t, ctx, input, text); err != nil {
				t.Fatal(err)
			}
			source := filepath.Join(text, "decks", "Convertible", "notes", "你好.note.anki.toml")
			note := &ankiascodepb.NoteConfig{}
			readFixtureTOML(t, source, note)
			note.Cards = nil
			destination := filepath.Join(text, "notes", "retained.note.anki.toml")
			writeFixtureTOML(t, destination, note)
			if err := os.Remove(source); err != nil {
				t.Fatal(err)
			}
			writeFixtureTOML(t, filepath.Join(text, "collection.anki.toml"), &ankiascodepb.CollectionConfig{Version: 1, Moved: []*ankiascodepb.ResourceMove{
				{From: "decks/Convertible/notes/你好.note.anki.toml", To: "notes/retained.note.anki.toml"},
			}})
			marker := filepath.Join(text, "decks", "Convertible", "deck.anki.toml")
			writeFixtureTOML(t, marker, &ankiascodepb.DeckConfig{Filtered: &decks.Deck_Filtered{}})
			before, err := os.ReadFile(input)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := planForFixture(ctx, input, text, ""); err == nil {
				t.Error("plan accepted converting an occupied normal deck to filtered")
			}
			if err := applyForFixture(ctx, input, text, ""); err == nil {
				t.Error("apply accepted converting an occupied normal deck to filtered")
			}
			after, err := os.ReadFile(input)
			if err != nil || string(before) != string(after) {
				t.Fatal("unsafe conversion changed collection bytes")
			}
			writeFixtureTOML(t, marker, &ankiascodepb.DeckConfig{Normal: &decks.Deck_Normal{}})
			note.Cards = []*ankiascodepb.CardConfig{
				{Ordinal: 0, Deck: proto.String("decks/grammar/deck.anki.toml")},
				{Ordinal: 1, Deck: proto.String("decks/grammar/deck.anki.toml")},
			}
			writeFixtureTOML(t, destination, note)
			if err := applyForFixture(ctx, input, text, ""); err != nil {
				t.Fatal(err)
			}
			writeFixtureTOML(t, marker, &ankiascodepb.DeckConfig{Filtered: &decks.Deck_Filtered{}})
			if err := applyForFixture(ctx, input, text, ""); err != nil {
				t.Fatal(err)
			}
			assertConverged(t, input, text)
			if format == "package" {
				client = inspectFixture(t, input)
			} else {
				client = connectFixture(t, input)
			}
			converted, err := client.decks.GetDeck(ctx, &decks.DeckId{Did: created.Id})
			if err != nil || converted.GetFiltered() == nil {
				t.Fatalf("empty normal deck was not converted: %v", err)
			}
			for _, id := range fixture.cardIDs {
				card, err := client.cards.GetCard(ctx, &cards.CardId{Cid: id})
				if err != nil || card.GetDeckId() != 1 || card.GetOriginalDeckId() != 0 {
					t.Fatalf("conversion changed retained normal cards: %v", err)
				}
			}
			retainedHistory, err := client.stats.GetReviewLogs(ctx, &cards.CardId{Cid: fixture.cardIDs[0]})
			if err != nil || !proto.Equal(history, retainedHistory) {
				t.Fatalf("conversion changed review history: %v", err)
			}
			if err := client.close(); err != nil {
				t.Fatal(err)
			}
			writeFixtureTOML(t, marker, &ankiascodepb.DeckConfig{Normal: &decks.Deck_Normal{}})
			addition := proto.Clone(note).(*ankiascodepb.NoteConfig)
			addition.Cards = []*ankiascodepb.CardConfig{
				{Ordinal: 0, Deck: proto.String("decks/grammar/deck.anki.toml")},
				{Ordinal: 1},
			}
			address := "decks/Convertible/notes/created.note.anki.toml"
			writeFixtureTOML(t, filepath.Join(text, filepath.FromSlash(address)), addition)
			if plan := fixturePlan(t, input, text); plan.NotesAdded != 1 || plan.DecksChanged != 1 {
				t.Fatalf("conversion and note creation missing from plan: %v", plan)
			}
			if err := applyForFixture(ctx, input, text, ""); err != nil {
				t.Fatal(err)
			}
			assertConverged(t, input, text)
			id := stateFixtureDump(t, input).GetResources().GetNotes()[address].GetId()
			if id == 0 {
				t.Fatal("created note has no stored binding")
			}
			if format == "package" {
				client = inspectFixture(t, input)
			} else {
				client = connectFixture(t, input)
			}
			generated, err := client.notes.CardsOfNote(ctx, &notes.NoteId{Nid: id})
			if err != nil || len(generated.GetCids()) != 2 {
				t.Fatalf("created note is missing generated cards: %v", err)
			}
			for _, id := range generated.Cids {
				card, err := client.cards.GetCard(ctx, &cards.CardId{Cid: id})
				if err != nil {
					t.Fatal(err)
				}
				destination := int64(1)
				if card.TemplateIdx == 1 {
					destination = created.Id
				}
				if card.DeckId != destination || card.OriginalDeckId != 0 {
					t.Fatalf("converted deck ignored generated-card assignment: %v", card)
				}
			}
			if err := client.close(); err != nil {
				t.Fatal(err)
			}
			saveFixtureArtifact(t, input, "deck-kind-conversion-"+format+filepath.Ext(input))
			saveFixtureArtifact(t, marker, "converted-"+format+".deck.anki.toml")
		})
	}
}

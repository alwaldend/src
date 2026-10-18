package ankiascode

import (
	"context"
	"path/filepath"
	"testing"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"github.com/ankitects/anki/proto/anki/cards"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"google.golang.org/protobuf/proto"
)

// TestFilteredCardProtection verifies filtered cards are excluded and conflicts never publish.
func TestFilteredCardProtection(t *testing.T) {
	fixture := newFixture(t)
	client := connectFixture(t, fixture.database)
	ctx := context.Background()
	definition, err := client.types.GetNotetype(ctx, &notetypes.NotetypeId{Ntid: fixture.typeID})
	if err != nil {
		t.Fatal(err)
	}
	for index, template := range definition.Templates {
		template.Config.Id = proto.Int64(int64(index + 101))
	}
	if _, err := client.types.UpdateNotetype(ctx, definition); err != nil {
		t.Fatal(err)
	}
	created, err := client.decks.AddDeck(ctx, &decks.Deck{Name: "Filtered", Kind: &decks.Deck_Filtered_{Filtered: &decks.Deck_Filtered{}}})
	if err != nil {
		t.Fatal(err)
	}
	homeDeck, err := client.decks.AddDeck(ctx, &decks.Deck{Name: "Home", Kind: &decks.Deck_Normal_{Normal: &decks.Deck_Normal{ConfigId: 1}}})
	if err != nil {
		t.Fatal(err)
	}
	resident, err := client.cards.GetCard(ctx, &cards.CardId{Cid: fixture.cardIDs[0]})
	if err != nil {
		t.Fatal(err)
	}
	resident.OriginalDeckId, resident.DeckId, resident.OriginalDue = homeDeck.Id, created.Id, resident.Due
	if _, err := client.cards.UpdateCards(ctx, &cards.UpdateCardsRequest{Cards: []*cards.Card{resident}, SkipUndoEntry: true}); err != nil {
		t.Fatal(err)
	}
	resident, err = client.cards.GetCard(ctx, &cards.CardId{Cid: resident.Id})
	if err != nil {
		t.Fatal(err)
	}
	if err := client.close(); err != nil {
		t.Fatal(err)
	}
	packageFixture(t, fixture.database, fixture.base)
	root := t.TempDir()
	text := filepath.Join(root, "text")
	if err := prepareFixtureConfigs(t, ctx, fixture.base, text); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(text, "decks", "Filtered", "deck.anki.toml")
	declaration := &ankiascodepb.DeckConfig{}
	readFixtureTOML(t, marker, declaration)
	declaration.Filtered.Reschedule = true
	writeFixtureTOML(t, marker, declaration)
	result := filepath.Join(root, "filtered.colpkg")
	if err := buildForFixture(ctx, fixture.base, text, result, ""); err != nil {
		t.Fatal(err)
	}
	assertConverged(t, result, text)
	actual := inspectFixture(t, result)
	deck, err := actual.decks.GetDeck(ctx, &decks.DeckId{Did: created.Id})
	if err != nil || !deck.GetFiltered().GetReschedule() {
		t.Fatalf("filtered API settings: %v", err)
	}
	retained, err := actual.cards.GetCard(ctx, &cards.CardId{Cid: resident.Id})
	if err != nil || !proto.Equal(resident, retained) {
		t.Fatalf("resident filtered card changed: %v", err)
	}
	if err := actual.close(); err != nil {
		t.Fatal(err)
	}
	// A filtered card's original deck is occupied even without resident normal cards.
	homeMarker := filepath.Join(text, "decks", "Home", "deck.anki.toml")
	originalHome := &ankiascodepb.DeckConfig{}
	readFixtureTOML(t, homeMarker, originalHome)
	writeFixtureTOML(t, homeMarker, &ankiascodepb.DeckConfig{Filtered: &decks.Deck_Filtered{}})
	originalPackage, err := fixtureContent(t, result)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := planForFixture(ctx, result, text, ""); err == nil {
		t.Fatal("filtered card home-deck conversion accepted")
	}
	if err := applyForFixture(ctx, result, text, ""); err == nil {
		t.Fatal("filtered card home-deck conversion applied")
	}
	retainedPackage, err := fixtureContent(t, result)
	if err != nil || string(originalPackage) != string(retainedPackage) {
		t.Fatal("rejected home-deck conversion changed archive")
	}
	writeFixtureTOML(t, homeMarker, originalHome)

	notePath := fixtureNotePath(text)
	note := &ankiascodepb.NoteConfig{}
	readFixtureTOML(t, notePath, note)
	if len(note.Cards) != 1 || note.Cards[0].Ordinal == int64(resident.TemplateIdx) {
		t.Fatal("filtered card exported as managed")
	}
	settingsPath := filepath.Join(text, "note_types", "Chinese.notetype.anki.toml")
	settings := &ankiascodepb.NoteTypeConfig{}
	readFixtureTOML(t, settingsPath, settings)
	settings.Templates[0], settings.Templates[1] = settings.Templates[1], settings.Templates[0]
	writeFixtureTOML(t, settingsPath, settings)
	original, err := fixtureContent(t, result)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := planForFixture(ctx, result, text, ""); err == nil {
		t.Fatal("template identity change affecting a filtered card accepted")
	}
	if err := applyForFixture(ctx, result, text, ""); err == nil {
		t.Fatal("template identity change affecting a filtered card applied")
	}
	retainedArchive, err := fixtureContent(t, result)
	if err != nil || string(original) != string(retainedArchive) {
		t.Fatal("rejected template identity change replaced the archive")
	}
	settings.Templates[0], settings.Templates[1] = settings.Templates[1], settings.Templates[0]
	writeFixtureTOML(t, settingsPath, settings)
	note.Cards = append(note.Cards, &ankiascodepb.CardConfig{Ordinal: int64(resident.TemplateIdx), Deck: proto.String("decks/grammar/deck.anki.toml")})
	writeFixtureTOML(t, notePath, note)
	before, err := fixtureContent(t, result)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := planForFixture(ctx, result, text, ""); err == nil {
		t.Fatal("filtered card conflict accepted")
	}
	if err := applyForFixture(ctx, result, text, ""); err == nil {
		t.Fatal("filtered card conflict applied")
	}
	after, err := fixtureContent(t, result)
	if err != nil || string(before) != string(after) {
		t.Fatal("rejected filtered edit changed archive")
	}
	saveFixtureArtifact(t, marker, "filtered-api.deck.anki.toml")
}

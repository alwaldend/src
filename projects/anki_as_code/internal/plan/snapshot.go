package plan

import (
	"fmt"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/resource"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"google.golang.org/protobuf/proto"
)

// snapshot indexes caller-supplied native data and its stored ownership.
type snapshot struct {
	// bindings contains the expected stored state before declared moves.
	bindings *collectionpb.State
	// presets indexes existing unmanaged scheduling presets by native ID.
	presets map[int64]string
	// notes indexes native notes by ID.
	notes map[int64]runtimeNote
	// noteTypes indexes native type definitions by ID.
	noteTypes map[int64]*runtimeNoteType
	// decks contains captured normal and filtered decks.
	decks []runtimeDeck
	// cardDeckIDs indexes normal-card destinations by note ID and ordinal.
	cardDeckIDs map[cardKey]int64
	// occupiedDecks records resident and filtered-home deck references.
	occupiedDecks map[int64]bool
	// filteredCards retains references needed to protect unmanaged cards in filtered decks.
	filteredCards []filteredCard
}

// filteredCard contains only the references needed to protect an unmanaged card.
type filteredCard struct {
	// id identifies the unmanaged card.
	id int64
	// noteID identifies the native note referenced by the card.
	noteID int64
	// ordinal identifies the card template within its note type.
	ordinal int
	// deckID identifies the filtered deck currently holding the card.
	deckID int64
	// homeDeckID identifies the original deck retained by Anki for restoration.
	homeDeckID int64
}

// readSnapshot indexes supplied native values and bindings without mutating either input.
func readSnapshot(native *resource.Snapshot, bindings *collectionpb.State) (snapshot, error) {
	current := snapshot{
		notes:         map[int64]runtimeNote{},
		noteTypes:     map[int64]*runtimeNoteType{},
		cardDeckIDs:   map[cardKey]int64{},
		occupiedDecks: map[int64]bool{},
	}
	current.bindings = bindings
	deckPaths, notePaths, typePaths := map[int64]string{}, map[int64]string{}, map[int64]string{}
	for address, deck := range bindings.GetResources().GetDecks() {
		deckPaths[deck.Id] = address
	}
	for address, note := range bindings.GetResources().GetNotes() {
		notePaths[note.Id] = address
	}
	for address, definition := range bindings.GetResources().GetNoteTypes() {
		typePaths[definition.Id] = address
	}
	for _, source := range native.NoteTypes {
		definition := proto.Clone(source).(*notetypes.Notetype)
		settings := &collectionpb.NoteTypeConfig{Name: definition.Name, Config: definition.Config, Fields: definition.Fields}
		for _, field := range settings.Fields {
			field.Ord = nil
		}
		for _, template := range definition.Templates {
			settings.Templates = append(settings.Templates, &collectionpb.TemplateConfig{Name: template.Name, Config: template.Config})
		}
		current.noteTypes[definition.Id] = &runtimeNoteType{id: definition.Id, path: typePaths[definition.Id], definition: settings, native: source}
	}
	normalDecks := map[int64]bool{}
	for _, value := range native.Decks {
		declaration := &collectionpb.DeckConfig{Normal: value.GetNormal(), Filtered: value.GetFiltered()}
		current.decks = append(current.decks, runtimeDeck{id: value.Id, path: deckPaths[value.Id], config: declaration, title: value.Name})
		normalDecks[value.Id] = declaration.Normal != nil
	}
	current.presets = native.Presets
	for _, value := range native.Notes {
		definition := current.noteTypes[value.NotetypeId]
		if definition == nil || len(value.Fields) != len(definition.definition.Fields) {
			return current, fmt.Errorf("note %d has an unavailable field layout", value.Id)
		}
		declaration := &collectionpb.NoteConfig{NoteType: definition.path, Tags: value.Tags, Fields: map[string]string{}}
		for index, field := range definition.definition.Fields {
			declaration.Fields[field.Name] = value.Fields[index]
		}
		current.notes[value.Id] = runtimeNote{id: value.Id, typeID: value.NotetypeId, path: notePaths[value.Id], config: declaration, bodies: value.Fields}
	}
	for _, value := range native.Cards {
		id := value.Id
		note, exists := current.notes[value.NoteId]
		if !exists {
			return current, fmt.Errorf("card %d references unavailable note %d", id, value.NoteId)
		}
		_, exists = normalDecks[value.DeckId]
		if !exists {
			return current, fmt.Errorf("card %d references unavailable deck %d", id, value.DeckId)
		}
		current.occupiedDecks[value.DeckId] = true
		if value.OriginalDeckId != 0 {
			current.occupiedDecks[value.OriginalDeckId] = true
		}
		if value.OriginalDeckId != 0 || !normalDecks[value.DeckId] {
			current.filteredCards = append(current.filteredCards, filteredCard{id: id, noteID: value.NoteId, ordinal: int(value.TemplateIdx), deckID: value.DeckId, homeDeckID: value.OriginalDeckId})
			continue
		}
		note.config.Cards = append(note.config.Cards, &collectionpb.CardConfig{Ordinal: int64(value.TemplateIdx), Deck: proto.String(deckPaths[value.DeckId])})
		current.notes[value.NoteId] = note
		key := cardKey{noteID: note.id, ordinal: int64(value.TemplateIdx)}
		current.cardDeckIDs[key] = value.DeckId
	}
	return current, nil
}

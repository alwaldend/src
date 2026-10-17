package plan

import (
	"cmp"
	"fmt"
	"maps"
	"slices"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"github.com/ankitects/anki/proto/anki/decks"
	"google.golang.org/protobuf/proto"
)

// compile prepares each resource delta once, producing native operations and report entries together.
func compile(desired desiredState, before snapshot) *collectionpb.CollectionPlan {
	result := &collectionpb.CollectionPlan{Report: &collectionpb.CollectionChangePlan{}, ExpectedState: proto.CloneOf(before.bindings), ProposedBindings: bindingInputs(desired)}
	compileDecks(result, desired, before)
	compileNoteTypes(result, desired, before)
	compileNotes(result, desired, before)
	compileCards(result, desired, before)
	result.Report.CardGenerationDeferred = result.Report.NotesAdded > 0 || result.Report.NotesChanged > 0 || len(result.NoteTypes) > 0
	return result
}

// temporaryName reserves an unused root name for native rename swaps and replacements.
func temporaryName(kind string, id int64, reserved map[string]bool) string {
	prefix := fmt.Sprintf("_anki_as_code_pending_%s_%d", kind, id)
	name := prefix
	for suffix := 1; reserved[folded(name)]; suffix++ {
		name = fmt.Sprintf("%s_%d", prefix, suffix)
	}
	reserved[folded(name)] = true
	return name
}

// compileDecks compares prepared decks and emits collision-free name releases.
func compileDecks(result *collectionpb.CollectionPlan, desired desiredState, before snapshot) {
	reserved := map[string]bool{}
	for name := range before.deckNames {
		reserved[name] = true
	}
	for _, deck := range desired.decksByPath {
		reserved[folded(deck.Value.Name)] = true
	}
	for id, deck := range before.decks {
		path := before.deckPaths[id]
		if path == "" {
			continue
		}
		next := desired.decksByID[id]
		if next != nil && next.Value.Name == deck.Name {
			continue
		}
		result.DeckRenames = append(result.DeckRenames, &decks.RenameDeckRequest{DeckId: id, NewName: temporaryName("deck", id, reserved)})
		if next == nil {
			result.RemovedDecks = append(result.RemovedDecks, id)
			result.Report.DecksDeleted++
			addChange(result.Report, "deck", path, "delete")
		}
	}
	for _, deck := range desired.decksByPath {
		previous := before.decks[deck.Value.Id]
		changed := previous == nil || previous.Name != deck.Value.Name || !proto.Equal(previous.GetNormal(), deck.Value.GetNormal()) || !proto.Equal(previous.GetFiltered(), deck.Value.GetFiltered())
		if previous == nil {
			result.Report.DecksAdded++
			addChange(result.Report, "deck", deck.Path, "add")
		} else {
			addMove(result.Report, "deck", before.deckPaths[previous.Id], deck.Path)
			if changed {
				result.Report.DecksChanged++
				addChange(result.Report, "deck", deck.Path, "update")
			}
		}
		// Parent name releases may move retained descendants, so restore their final names too.
		if changed || len(result.DeckRenames) != 0 {
			result.Decks = append(result.Decks, deck)
		}
	}
	slices.SortFunc(result.Decks, func(a, b *collectionpb.PlannedDeck) int {
		return cmp.Compare(folded(a.Value.Name), folded(b.Value.Name))
	})
	slices.Sort(result.RemovedDecks)
	slices.SortFunc(result.DeckRenames, func(a, b *decks.RenameDeckRequest) int { return cmp.Compare(a.DeckId, b.DeckId) })
}

// compileNoteTypes compares prepared definitions and releases tracked native names.
func compileNoteTypes(result *collectionpb.CollectionPlan, desired desiredState, before snapshot) {
	reserved := map[string]bool{}
	for name := range before.typeNames {
		reserved[name] = true
	}
	for _, definition := range desired.noteTypes {
		reserved[folded(definition.Value.Name)] = true
	}
	for _, id := range slices.Sorted(maps.Keys(before.noteTypes)) {
		old := before.noteTypes[id]
		path := before.typePaths[id]
		if path == "" {
			continue
		}
		next := desired.typesByID[id]
		if next != nil && next.Value.Name == old.Name {
			continue
		}
		value := proto.CloneOf(old)
		value.Name = temporaryName("type", id, reserved)
		result.NoteTypeRenames = append(result.NoteTypeRenames, value)
		if next == nil {
			result.RemovedNoteTypes = append(result.RemovedNoteTypes, id)
			result.Report.NoteTypesDeleted++
			addChange(result.Report, "note_type", path, "delete")
		}
	}
	for _, operation := range desired.noteTypes {
		old := before.noteTypes[operation.Value.Id]
		if old == nil {
			result.Report.NoteTypesAdded++
			addChange(result.Report, "note_type", operation.Path, "add")
		} else {
			addMove(result.Report, "note_type", before.typePaths[old.Id], operation.Path)
			if unchangedNoteType(old, operation) {
				continue
			}
			result.Report.NoteTypesChanged++
			addChange(result.Report, "note_type", operation.Path, "update")
		}
		for _, template := range operation.Value.Templates {
			if operation.TemplateDecks[template.Name] != "" {
				template.Config.TargetDeckId = 0
			}
		}
		result.NoteTypes = append(result.NoteTypes, operation)
	}
	slices.SortFunc(result.NoteTypes, func(a, b *collectionpb.PlannedNoteType) int { return cmp.Compare(a.Path, b.Path) })
}

// compileNotes compares prepared content and emits writes required after native migrations.
func compileNotes(result *collectionpb.CollectionPlan, desired desiredState, before snapshot) {
	for _, note := range desired.notes {
		previous := before.notes[note.Value.Id]
		changed := previous == nil || !slices.Equal(previous.Fields, note.Value.Fields) || !slices.Equal(previous.Tags, note.Value.Tags)
		if previous == nil {
			result.Report.NotesAdded++
			addChange(result.Report, "note", note.Path, "add")
		} else {
			addMove(result.Report, "note", before.notePaths[previous.Id], note.Path)
			if changed {
				result.Report.NotesChanged++
				addChange(result.Report, "note", note.Path, "update")
			}
		}
		structural := previous != nil && !sameCardStructure(before.noteTypes[previous.NotetypeId], desired.noteTypes[note.NoteType].Value)
		if changed || structural {
			result.Notes = append(result.Notes, note)
		}
	}
	for id := range before.notes {
		if path := before.notePaths[id]; desired.notesByID[id] == nil && path != "" {
			result.RemovedNotes = append(result.RemovedNotes, id)
			result.Report.NotesDeleted++
			addChange(result.Report, "note", path, "delete")
		}
	}
	slices.Sort(result.RemovedNotes)
	slices.SortFunc(result.Notes, func(a, b *collectionpb.PlannedNote) int { return cmp.Compare(a.Path, b.Path) })
}

// compileCards reports known destinations and schedules only assignments that may need native work.
func compileCards(result *collectionpb.CollectionPlan, desired desiredState, before snapshot) {
	needed, changedTypes := map[string]bool{}, map[string]bool{}
	for _, operation := range result.Notes {
		needed[operation.Path] = true
	}
	for _, operation := range result.NoteTypes {
		changedTypes[operation.Path] = true
	}
	for key, deckID := range before.cardDeckIDs {
		previous := before.notePaths[key.noteID]
		if previous == "" {
			continue
		}
		note := desired.notesByID[key.noteID]
		if note == nil {
			result.Report.CardsDeleted++
			addChange(result.Report, "card", cardReference(previous, key.ordinal), "delete")
			continue
		}
		destination := cardDestination(desired.cards[note.Path], key.ordinal)
		if destination != "" && desired.decksByPath[destination].Value.Id != deckID {
			result.Report.CardsChanged++
			addChange(result.Report, "card", cardReference(note.Path, key.ordinal), "update")
			needed[note.Path] = true
		}
	}
	for path, assignments := range desired.cards {
		if assignments.Deck == "" && len(assignments.Cards) == 0 {
			continue
		}
		for _, card := range assignments.Cards {
			if _, exists := before.cardDeckIDs[cardKey{noteID: desired.notes[path].Value.Id, ordinal: card.Ordinal}]; !exists {
				needed[path] = true
			}
		}
		if needed[path] || changedTypes[desired.notes[path].NoteType] {
			result.CardAssignments = append(result.CardAssignments, assignments)
		}
	}
	slices.SortFunc(result.CardAssignments, func(a, b *collectionpb.PlannedCardAssignments) int { return cmp.Compare(a.Path, b.Path) })
	if len(result.CardAssignments) > 0 {
		for id, deck := range before.decks {
			if wanted := desired.decksByID[id]; wanted != nil {
				deck = wanted.Value
			}
			if deck.GetFiltered() != nil {
				result.FilteredDecks = append(result.FilteredDecks, id)
			}
		}
		slices.Sort(result.FilteredDecks)
	}
}

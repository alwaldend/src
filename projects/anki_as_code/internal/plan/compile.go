package plan

import (
	"cmp"
	"fmt"
	"slices"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"google.golang.org/protobuf/proto"
)

// compile prepares each resource delta once, producing native operations and report entries together.
func compile(desired desiredState, before snapshot) *collectionpb.CollectionPlan {
	result := &collectionpb.CollectionPlan{Report: &collectionpb.CollectionChangePlan{}, ExpectedState: proto.Clone(before.bindings).(*collectionpb.State), ProposedBindings: bindingInputs(desired)}
	compileDecks(result, desired, before)
	compileNoteTypes(result, desired, before)
	compileNotes(result, desired, before)
	compileCards(result, desired, before)
	if len(result.CardAssignments) == 0 {
		result.FilteredDecks = nil
	}
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

// compileDecks compares deck values and prepares ordered upserts and collision-free name releases.
func compileDecks(result *collectionpb.CollectionPlan, desired desiredState, before snapshot) {
	current, wanted := map[int64]runtimeDeck{}, map[int64]runtimeDeck{}
	reserved, filtered := map[string]bool{}, map[int64]bool{}
	for _, deck := range before.decks {
		current[deck.id] = deck
		reserved[folded(deck.title)] = true
		filtered[deck.id] = deck.config.Filtered != nil
	}
	for _, deck := range desired.decks {
		if deck.id != 0 {
			wanted[deck.id] = deck
			filtered[deck.id] = deck.config.Filtered != nil
		}
		reserved[folded(deck.title)] = true
	}
	for _, deck := range before.decks {
		if deck.path == "" {
			continue
		}
		next, retained := wanted[deck.id]
		if retained && next.title == deck.title {
			continue
		}
		result.DeckRenames = append(result.DeckRenames, &decks.RenameDeckRequest{DeckId: deck.id, NewName: temporaryName("deck", deck.id, reserved)})
		if !retained {
			result.RemovedDecks = append(result.RemovedDecks, deck.id)
			result.Report.DecksDeleted++
			addChange(result.Report, "deck", deck.path, "delete")
		}
	}
	for _, deck := range desired.decks {
		previous, exists := current[deck.id]
		changed := !exists || previous.title != deck.title || !proto.Equal(previous.config.Normal, deck.config.Normal) || !proto.Equal(previous.config.Filtered, deck.config.Filtered)
		if !exists {
			result.Report.DecksAdded++
			addChange(result.Report, "deck", deck.path, "add")
		} else {
			addMove(result.Report, "deck", previous.path, deck.path)
			if changed {
				result.Report.DecksChanged++
				addChange(result.Report, "deck", deck.path, "update")
			}
		}
		// Parent name releases may move retained descendants, so restore their final names too.
		if exists && !changed && len(result.DeckRenames) == 0 {
			continue
		}
		value := &decks.Deck{Id: deck.id, Name: deck.title}
		if deck.config.Normal != nil {
			value.Kind = &decks.Deck_Normal_{Normal: proto.Clone(deck.config.Normal).(*decks.Deck_Normal)}
		} else {
			value.Kind = &decks.Deck_Filtered_{Filtered: proto.Clone(deck.config.Filtered).(*decks.Deck_Filtered)}
		}
		result.Decks = append(result.Decks, &collectionpb.PlannedDeck{Path: deck.path, Value: value})
	}
	for id, enabled := range filtered {
		if enabled {
			result.FilteredDecks = append(result.FilteredDecks, id)
		}
	}
	slices.Sort(result.FilteredDecks)
	slices.Sort(result.RemovedDecks)
	slices.SortFunc(result.DeckRenames, func(a, b *decks.RenameDeckRequest) int { return cmp.Compare(a.DeckId, b.DeckId) })
}

// compileNoteTypes compares prepared definitions and emits native payloads with path references.
func compileNoteTypes(result *collectionpb.CollectionPlan, desired desiredState, before snapshot) {
	wanted := map[int64]*runtimeNoteType{}
	reserved := map[string]bool{}
	for _, definition := range before.noteTypes {
		reserved[folded(definition.definition.Name)] = true
	}
	for _, definition := range desired.noteTypes {
		if definition.id != 0 {
			wanted[definition.id] = definition
		}
		reserved[folded(definition.definition.Name)] = true
	}
	ids := make([]int64, 0, len(before.noteTypes))
	for id := range before.noteTypes {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	for _, id := range ids {
		old := before.noteTypes[id]
		if old.path == "" {
			continue
		}
		next := wanted[id]
		if next != nil && next.definition.Name == old.definition.Name {
			continue
		}
		value := proto.Clone(old.native).(*notetypes.Notetype)
		value.Name = temporaryName("type", id, reserved)
		result.NoteTypeRenames = append(result.NoteTypeRenames, value)
		if next == nil {
			result.RemovedNoteTypes = append(result.RemovedNoteTypes, id)
			result.Report.NoteTypesDeleted++
			addChange(result.Report, "note_type", old.path, "delete")
		}
	}
	for _, entry := range desired.noteTypes {
		old := before.noteTypes[entry.id]
		if old == nil {
			result.Report.NoteTypesAdded++
			addChange(result.Report, "note_type", entry.path, "add")
		} else {
			addMove(result.Report, "note_type", old.path, entry.path)
			if unchangedNoteType(old.definition, entry) {
				continue
			}
			result.Report.NoteTypesChanged++
			addChange(result.Report, "note_type", entry.path, "update")
		}
		operation := &collectionpb.PlannedNoteType{Path: entry.path, Value: proto.Clone(entry.native).(*notetypes.Notetype), SortFieldIndex: entry.definition.Config.SortFieldIdx}
		for index, source := range entry.definition.Templates {
			if source.TargetDeck != "" {
				if operation.TemplateDecks == nil {
					operation.TemplateDecks = map[string]string{}
				}
				operation.TemplateDecks[source.Name] = source.TargetDeck
				operation.Value.Templates[index].Config.TargetDeckId = 0
			}
		}
		result.NoteTypes = append(result.NoteTypes, operation)
	}
	slices.SortFunc(result.NoteTypes, func(a, b *collectionpb.PlannedNoteType) int { return cmp.Compare(a.Path, b.Path) })
}

// compileNotes compares content and prepares writes required after native field/template migrations.
func compileNotes(result *collectionpb.CollectionPlan, desired desiredState, before snapshot) {
	retained := map[int64]bool{}
	for _, note := range desired.notes {
		previous, exists := before.notes[note.id]
		changed := !exists || !slices.Equal(previous.bodies, note.bodies) || !slices.Equal(previous.config.Tags, note.config.Tags)
		if !exists {
			result.Report.NotesAdded++
			addChange(result.Report, "note", note.path, "add")
		} else {
			retained[note.id] = true
			addMove(result.Report, "note", previous.path, note.path)
			if changed {
				result.Report.NotesChanged++
				addChange(result.Report, "note", note.path, "update")
			}
		}
		structural := exists && !sameCardStructure(before.noteTypes[previous.typeID].definition, desired.noteTypes[note.config.NoteType])
		if changed || structural {
			result.Notes = append(result.Notes, &collectionpb.PlannedNote{Path: note.path, NoteType: note.config.NoteType, Deck: note.homeDeck, Value: &notes.Note{Id: note.id, Fields: slices.Clone(note.bodies), Tags: slices.Clone(note.config.Tags)}})
		}
	}
	for id, note := range before.notes {
		if !retained[id] && note.path != "" {
			result.RemovedNotes = append(result.RemovedNotes, note.id)
			result.Report.NotesDeleted++
			addChange(result.Report, "note", note.path, "delete")
		}
	}
	slices.Sort(result.RemovedNotes)
	slices.SortFunc(result.Notes, func(a, b *collectionpb.PlannedNote) int { return cmp.Compare(a.Path, b.Path) })
}

// compileCards reports known destinations and schedules only assignments that may need native work.
func compileCards(result *collectionpb.CollectionPlan, desired desiredState, before snapshot) {
	notesByID := map[int64]runtimeNote{}
	decksByPath := map[string]runtimeDeck{}
	needed, changedTypes := map[string]bool{}, map[string]bool{}
	for _, note := range desired.notes {
		if note.id != 0 {
			notesByID[note.id] = note
		}
	}
	for _, deck := range desired.decks {
		decksByPath[deck.path] = deck
	}
	for _, operation := range result.Notes {
		needed[operation.Path] = true
	}
	for _, operation := range result.NoteTypes {
		changedTypes[operation.Path] = true
	}
	for key, deckID := range before.cardDeckIDs {
		previous := before.notes[key.noteID]
		if previous.path == "" {
			continue
		}
		note, retained := notesByID[key.noteID]
		if !retained {
			result.Report.CardsDeleted++
			addChange(result.Report, "card", cardReference(previous.path, key.ordinal), "delete")
			continue
		}
		destination := cardDestination(note, key.ordinal)
		if destination != "" && decksByPath[destination].id != deckID {
			result.Report.CardsChanged++
			addChange(result.Report, "card", cardReference(note.path, key.ordinal), "update")
			needed[note.path] = true
		}
	}
	for _, note := range desired.notes {
		if note.homeDeck == "" && len(note.config.Cards) == 0 {
			continue
		}
		for _, card := range note.config.Cards {
			if _, exists := before.cardDeckIDs[cardKey{noteID: note.id, ordinal: card.Ordinal}]; !exists {
				needed[note.path] = true
			}
		}
		if !needed[note.path] && !changedTypes[note.config.NoteType] {
			continue
		}
		assignments := &collectionpb.PlannedCardAssignments{Path: note.path, Deck: note.homeDeck}
		for _, card := range note.config.Cards {
			assignments.Cards = append(assignments.Cards, proto.Clone(card).(*collectionpb.CardConfig))
		}
		result.CardAssignments = append(result.CardAssignments, assignments)
	}
	slices.SortFunc(result.CardAssignments, func(a, b *collectionpb.PlannedCardAssignments) int { return cmp.Compare(a.Path, b.Path) })
}

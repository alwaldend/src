package reconcile

import (
	"context"
	"fmt"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"github.com/ankitects/anki/proto/anki/cards"
	ankicol "github.com/ankitects/anki/proto/anki/collection"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"google.golang.org/protobuf/proto"
)

// nativeExecution holds only execution results for one detached prepared plan.
type nativeExecution struct {
	// engine supplies borrowed native clients and state persistence.
	engine *Engine
	// plan is a private copy updated with actual Anki identities.
	plan *collectionpb.CollectionPlan
	// changed retains Anki's native mutation flags.
	changed bool
}

// execute follows the prepared dependency phases without recomputing their intent.
func (run *nativeExecution) execute(ctx context.Context) error {
	if err := run.decks(ctx); err != nil {
		return fmt.Errorf("execute decks: %w", err)
	}
	if err := run.noteTypes(ctx); err != nil {
		return fmt.Errorf("execute note types: %w", err)
	}
	if err := run.notes(ctx); err != nil {
		return fmt.Errorf("execute notes: %w", err)
	}
	if err := run.removeNotesAndTypes(ctx); err != nil {
		return fmt.Errorf("execute note/type removals: %w", err)
	}
	if err := run.assignDecks(ctx); err != nil {
		return fmt.Errorf("execute card destinations: %w", err)
	}
	if len(run.plan.RemovedDecks) > 0 {
		response, err := run.engine.resources.Decks.RemoveDecks(ctx, &decks.DeckIds{Dids: run.plan.RemovedDecks})
		if err != nil {
			return fmt.Errorf("remove planned decks: %w", err)
		}
		run.recordChanges(response.GetChanges())
	}
	return nil
}

// decks applies prepared name releases and native values while preserving common metadata.
func (run *nativeExecution) decks(ctx context.Context) error {
	for _, rename := range run.plan.DeckRenames {
		response, err := run.engine.resources.Decks.RenameDeck(ctx, rename)
		if err != nil {
			return fmt.Errorf("release deck %d: %w", rename.DeckId, err)
		}
		run.recordChanges(response)
	}
	for _, operation := range run.plan.Decks {
		value := operation.Value
		if value.Id != 0 {
			current, err := run.engine.resources.Decks.GetDeck(ctx, &decks.DeckId{Did: value.Id})
			if err != nil {
				return fmt.Errorf("read deck metadata %q: %w", operation.Path, err)
			}
			current.Name, current.Kind = value.Name, value.Kind
			response, err := run.engine.resources.Decks.UpdateDeck(ctx, current)
			if err != nil {
				return fmt.Errorf("update deck %q: %w", operation.Path, err)
			}
			run.recordChanges(response)
		} else {
			response, err := run.engine.resources.Decks.AddDeck(ctx, value)
			if err != nil {
				return fmt.Errorf("add deck %q: %w", operation.Path, err)
			}
			run.recordChanges(response.GetChanges())
			binding := run.plan.ProposedBindings.Decks[operation.Path]
			binding.Id = response.Id
		}
	}
	return nil
}

// noteTypes executes prepared member mappings and restores declared presentation after migration.
func (run *nativeExecution) noteTypes(ctx context.Context) error {
	for _, value := range run.plan.NoteTypeRenames {
		response, err := run.engine.resources.NoteTypes.UpdateNotetype(ctx, value)
		if err != nil {
			return fmt.Errorf("release note type %d: %w", value.Id, err)
		}
		run.recordChanges(response)
	}
	for _, operation := range run.plan.NoteTypes {
		value := operation.Value
		for _, template := range value.Templates {
			if path := operation.TemplateDecks[template.Name]; path != "" {
				template.Config.TargetDeckId = run.deckID(path)
			}
		}
		if value.Id != 0 {
			current, err := run.engine.resources.NoteTypes.GetNotetype(ctx, &notetypes.NotetypeId{Ntid: value.Id})
			if err != nil {
				return fmt.Errorf("read type metadata %q: %w", operation.Path, err)
			}
			value.Config = typeConfig(value.Config, current.Config)
			response, err := run.engine.resources.NoteTypes.UpdateNotetype(ctx, value)
			if err != nil {
				return fmt.Errorf("update type %q: %w", operation.Path, err)
			}
			run.recordChanges(response)
		} else {
			response, err := run.engine.resources.NoteTypes.AddNotetype(ctx, value)
			if err != nil {
				return fmt.Errorf("add type %q: %w", operation.Path, err)
			}
			run.recordChanges(response.GetChanges())
			binding := run.plan.ProposedBindings.NoteTypes[operation.Path]
			binding.Id, value.Id = response.Id, response.Id
		}
		if err := run.applyPresentation(ctx, operation); err != nil {
			return fmt.Errorf("restore type presentation %q: %w", operation.Path, err)
		}
		current, err := run.engine.resources.NoteTypes.GetNotetype(ctx, &notetypes.NotetypeId{Ntid: value.Id})
		if err != nil {
			return fmt.Errorf("read resulting type members %q: %w", operation.Path, err)
		}
		binding := run.plan.ProposedBindings.NoteTypes[operation.Path]
		if len(current.Fields) != len(value.Fields) || len(current.Templates) != len(value.Templates) {
			return fmt.Errorf("Anki changed the planned member layout %q", operation.Path)
		}
		for index, field := range current.Fields {
			binding.Fields[value.Fields[index].Name].Id = field.Config.GetId()
		}
		for index, template := range current.Templates {
			binding.Templates[value.Templates[index].Name].Id = template.Config.GetId()
		}
	}
	return nil
}

// typeConfig preserves Anki-owned provenance and derived data when applying editable values.
func typeConfig(wanted, previous *notetypes.Notetype_Config) *notetypes.Notetype_Config {
	value := proto.Clone(wanted).(*notetypes.Notetype_Config)
	value.OriginalId, value.TargetDeckIdUnused = previous.OriginalId, previous.TargetDeckIdUnused
	value.Reqs, value.Other = previous.Reqs, previous.Other
	return value
}

// applyPresentation restores declared values after Anki rewrites structurally migrated templates.
func (run *nativeExecution) applyPresentation(ctx context.Context, operation *collectionpb.PlannedNoteType) error {
	wanted := operation.Value
	current, err := run.engine.resources.NoteTypes.GetNotetype(ctx, &notetypes.NotetypeId{Ntid: wanted.Id})
	if err != nil {
		return fmt.Errorf("read migrated type: %w", err)
	}
	if len(current.Templates) != len(wanted.Templates) {
		return fmt.Errorf("Anki changed the planned template count")
	}
	current.Config = typeConfig(wanted.Config, current.Config)
	current.Config.SortFieldIdx = operation.SortFieldIndex
	for index, template := range wanted.Templates {
		previous := current.Templates[index].Config
		current.Templates[index].Config = proto.Clone(template.Config).(*notetypes.Notetype_Template_Config)
		current.Templates[index].Config.Id, current.Templates[index].Config.Other = previous.Id, previous.Other
	}
	response, err := run.engine.resources.NoteTypes.UpdateNotetype(ctx, current)
	if err != nil {
		return fmt.Errorf("apply migrated presentation: %w", err)
	}
	run.recordChanges(response)
	return nil
}

// notes applies prepared ordered content and binds newly allocated note identities.
func (run *nativeExecution) notes(ctx context.Context) error {
	updates, additions := []*notes.Note{}, []*notes.AddNoteRequest{}
	paths := []string{}
	for _, operation := range run.plan.Notes {
		value := operation.Value
		value.NotetypeId = run.plan.ProposedBindings.NoteTypes[operation.NoteType].GetId()
		if value.Id != 0 {
			current, err := run.engine.resources.Notes.GetNote(ctx, &notes.NoteId{Nid: value.Id})
			if err != nil {
				return fmt.Errorf("read note metadata %q: %w", operation.Path, err)
			}
			current.Fields, current.Tags = value.Fields, value.Tags
			updates = append(updates, current)
		} else {
			generated, err := run.engine.resources.Notes.NewNote(ctx, &notetypes.NotetypeId{Ntid: value.NotetypeId})
			if err != nil {
				return fmt.Errorf("initialize note %q: %w", operation.Path, err)
			}
			generated.Fields, generated.Tags = value.Fields, value.Tags
			additions = append(additions, &notes.AddNoteRequest{Note: generated, DeckId: run.deckID(operation.Deck)})
			paths = append(paths, operation.Path)
		}
	}
	if len(additions) > 0 {
		response, err := run.engine.resources.Notes.AddNotes(ctx, &notes.AddNotesRequest{Requests: additions})
		if err != nil {
			return fmt.Errorf("add planned notes: %w", err)
		}
		run.recordChanges(response.GetChanges())
		if len(response.Nids) != len(paths) {
			return fmt.Errorf("Anki returned %d IDs for %d planned notes", len(response.Nids), len(paths))
		}
		for index, path := range paths {
			run.plan.ProposedBindings.Notes[path].Id = response.Nids[index]
		}
	}
	if len(updates) > 0 {
		response, err := run.engine.resources.Notes.UpdateNotes(ctx, &notes.UpdateNotesRequest{Notes: updates, SkipUndoEntry: true})
		if err != nil {
			return fmt.Errorf("update planned notes: %w", err)
		}
		run.recordChanges(response)
	}
	return nil
}

// removeNotesAndTypes performs explicit tracked deletions before card destinations and decks.
func (run *nativeExecution) removeNotesAndTypes(ctx context.Context) error {
	if len(run.plan.RemovedNotes) > 0 {
		response, err := run.engine.resources.Notes.RemoveNotes(ctx, &notes.RemoveNotesRequest{NoteIds: run.plan.RemovedNotes})
		if err != nil {
			return fmt.Errorf("remove planned notes: %w", err)
		}
		run.recordChanges(response.GetChanges())
	}
	for _, id := range run.plan.RemovedNoteTypes {
		response, err := run.engine.resources.NoteTypes.RemoveNotetype(ctx, &notetypes.NotetypeId{Ntid: id})
		if err != nil {
			return fmt.Errorf("remove planned type %d: %w", id, err)
		}
		run.recordChanges(response)
	}
	return nil
}

// assignDecks applies supplied destinations to Anki's actual generated normal cards.
func (run *nativeExecution) assignDecks(ctx context.Context) error {
	filtered := map[int64]bool{}
	for _, id := range run.plan.FilteredDecks {
		filtered[id] = true
	}
	byDeck := map[int64][]int64{}
	for _, assignment := range run.plan.CardAssignments {
		identities, err := run.engine.resources.Notes.CardsOfNote(ctx, &notes.NoteId{Nid: run.plan.ProposedBindings.Notes[assignment.Path].Id})
		if err != nil {
			return fmt.Errorf("read generated cards %q: %w", assignment.Path, err)
		}
		generated := map[int64]bool{}
		for _, id := range identities.Cids {
			card, err := run.engine.resources.Cards.GetCard(ctx, &cards.CardId{Cid: id})
			if err != nil {
				return fmt.Errorf("read generated card %d: %w", id, err)
			}
			if card.OriginalDeckId != 0 || filtered[card.DeckId] {
				continue
			}
			ordinal := int64(card.TemplateIdx)
			generated[ordinal] = true
			destination := assignment.Deck
			for _, override := range assignment.Cards {
				if override.Ordinal == ordinal {
					destination = override.GetDeck()
					break
				}
			}
			if destination != "" && card.DeckId != run.deckID(destination) {
				byDeck[run.deckID(destination)] = append(byDeck[run.deckID(destination)], card.Id)
			}
		}
		for _, override := range assignment.Cards {
			if !generated[override.Ordinal] {
				return fmt.Errorf("note %q has no generated card at planned ordinal %d", assignment.Path, override.Ordinal)
			}
		}
	}
	for id, cardIDs := range byDeck {
		response, err := run.engine.resources.Cards.SetDeck(ctx, &cards.SetDeckRequest{CardIds: cardIDs, DeckId: id})
		if err != nil {
			return fmt.Errorf("assign planned deck %d: %w", id, err)
		}
		run.recordChanges(response.GetChanges())
	}
	return nil
}

// deckID resolves a prepared path binding after any native deck creation.
func (run *nativeExecution) deckID(path string) int64 {
	return run.plan.ProposedBindings.Decks[path].GetId()
}

// recordChanges retains native entity flags, including non-undoable note writes.
func (run *nativeExecution) recordChanges(changes *ankicol.OpChanges) {
	run.changed = run.changed || changes.GetMtime() || changes.GetCard() || changes.GetNote() || changes.GetDeck() || changes.GetTag() || changes.GetNotetype() || changes.GetConfig() || changes.GetDeckConfig()
}

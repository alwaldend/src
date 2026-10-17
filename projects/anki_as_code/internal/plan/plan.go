package plan

import (
	"fmt"
	"slices"
	"strings"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/resource"
	"google.golang.org/protobuf/proto"
)

// Build prepares executable native operations from caller-supplied configs, resources and state.
// It performs no I/O and changes no inputs. Callers own freshness and exclusive
// access to the same collection and state namespace until execution. Anki owns
// native validation, normalization and generated-card migrations during apply.
func Build(configs *collectionpb.CollectionConfigSet, native *resource.Snapshot, stored *collectionpb.State) (*collectionpb.CollectionPlan, error) {
	desired, before, err := resolve(configs, native, stored)
	if err != nil {
		return nil, fmt.Errorf("resolve collection plan: %w", err)
	}
	result := compile(desired, before)
	result.Report.StateChanged = !proto.Equal(&collectionpb.State{Version: stored.Version, Resources: result.ProposedBindings}, stored)
	if result.Report.StateChanged && len(result.Report.Changes) == 0 {
		addChange(result.Report, "state", "collection.anki.toml", "update")
	}
	slices.SortFunc(result.Report.Changes, func(a, b *collectionpb.CollectionChange) int {
		if order := strings.Compare(a.Resource, b.Resource); order != 0 {
			return order
		}
		if order := strings.Compare(a.Path, b.Path); order != 0 {
			return order
		}
		if order := strings.Compare(a.Action, b.Action); order != 0 {
			return order
		}
		return strings.Compare(a.From, b.From)
	})
	return result, nil
}

// cardDestination resolves an explicit ordinal override or the note's declaring deck.
func cardDestination(note runtimeNote, ordinal int64) string {
	for _, assignment := range note.config.Cards {
		if assignment.Ordinal == ordinal {
			return assignment.GetDeck()
		}
	}
	return note.homeDeck
}

// addChange appends one path-addressed resource action.
func addChange(plan *collectionpb.CollectionChangePlan, resource, path, action string) {
	plan.Changes = append(plan.Changes, &collectionpb.CollectionChange{Resource: resource, Path: path, Action: action})
}

// addMove reports a retained resource whose bound configuration path changed.
func addMove(plan *collectionpb.CollectionChangePlan, resource, before, after string) {
	if before != after {
		plan.Changes = append(plan.Changes, &collectionpb.CollectionChange{Resource: resource, Path: after, From: before, Action: "move"})
	}
}

// equalNoteTypes compares editable values while ignoring native metadata and member hints.
func equalNoteTypes(left, right *collectionpb.NoteTypeConfig) bool {
	if left == nil || right == nil {
		return left == right
	}
	return proto.Equal(editableType(left), editableType(right))
}

// editableType clones resolved configuration without native provenance, caches or identity hints.
func editableType(definition *collectionpb.NoteTypeConfig) *collectionpb.NoteTypeConfig {
	value := proto.Clone(definition).(*collectionpb.NoteTypeConfig)
	value.SortField = ""
	value.Config.OriginalId = nil
	value.Config.TargetDeckIdUnused = 0
	value.Config.Reqs = nil
	value.Config.Other = nil
	for _, field := range value.Fields {
		field.Ord = nil
		field.Config.Id = nil
		field.Config.Other = nil
	}
	for _, template := range value.Templates {
		template.TargetDeck = ""
		template.Config.Id = nil
		template.Config.Other = nil
	}
	return value
}

// cardReference identifies a generated card by its note path and template or cloze ordinal.
func cardReference(address string, ordinal int64) string {
	return fmt.Sprintf("%s#cards/%d", address, ordinal)
}

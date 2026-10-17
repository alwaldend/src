package plan

import (
	"cmp"
	"fmt"
	"slices"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/resource"
	"github.com/ankitects/anki/proto/anki/notetypes"
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
	previous := stored.GetResources()
	if previous == nil {
		previous = &collectionpb.Resources{}
	}
	result.Report.StateChanged = !proto.Equal(result.ProposedBindings, previous)
	if result.Report.StateChanged && len(result.Report.Changes) == 0 {
		addChange(result.Report, "state", "collection.anki.toml", "update")
	}
	slices.SortFunc(result.Report.Changes, func(a, b *collectionpb.CollectionChange) int {
		return cmp.Or(
			cmp.Compare(a.Resource, b.Resource),
			cmp.Compare(a.Path, b.Path),
			cmp.Compare(a.Action, b.Action),
			cmp.Compare(a.From, b.From),
		)
	})
	return result, nil
}

// cardDestination resolves an explicit ordinal override or the note's declaring deck.
func cardDestination(assignments *collectionpb.PlannedCardAssignments, ordinal int64) string {
	for _, assignment := range assignments.Cards {
		if assignment.Ordinal == ordinal {
			return assignment.GetDeck()
		}
	}
	return assignments.Deck
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

// editableType clones native editable values without provenance, caches or identity hints.
func editableType(definition *notetypes.Notetype) *notetypes.Notetype {
	value := proto.CloneOf(&notetypes.Notetype{
		Name:      definition.Name,
		Config:    definition.Config,
		Fields:    definition.Fields,
		Templates: definition.Templates,
	})
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
		template.Ord = nil
		template.Config.Id = nil
		template.Config.Other = nil
	}
	return value
}

// cardReference identifies a generated card by its note path and template or cloze ordinal.
func cardReference(address string, ordinal int64) string {
	return fmt.Sprintf("%s#cards/%d", address, ordinal)
}

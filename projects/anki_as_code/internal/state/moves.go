package state

import (
	"fmt"
	"strings"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
)

// validMoveAddress validates the filesystem portion of a resource or member address.
func validMoveAddress(address string) error {
	if err := validAddress(resourceRoot(address)); err != nil {
		return fmt.Errorf("move address %q: %w", address, err)
	}
	return nil
}

// addressDescendant tests exact resources, directory descendants, and nested type resources.
func addressDescendant(address, prefix string) bool {
	if resourceRoot(prefix) != prefix {
		return address == prefix
	}
	file := resourceRoot(address)
	return file == prefix || strings.HasPrefix(file, prefix+"/")
}

// moveState applies independent explicit moves to a clone; retained moves are idempotent.
func moveState(state *collectionpb.State, moves []*collectionpb.ResourceMove) (*collectionpb.State, error) {
	result := cloneState(state)
	for index, move := range moves {
		for _, other := range moves[:index] {
			endpoints := []string{other.From, other.To}
			for _, endpoint := range endpoints {
				if addressDescendant(move.From, endpoint) || addressDescendant(endpoint, move.From) || addressDescendant(move.To, endpoint) || addressDescendant(endpoint, move.To) {
					return nil, fmt.Errorf("move declarations overlap; use independent source and destination addresses")
				}
			}
		}
		if err := validMoveAddress(move.From); err != nil {
			return nil, fmt.Errorf("move source: %w", err)
		}
		if err := validMoveAddress(move.To); err != nil {
			return nil, fmt.Errorf("move destination: %w", err)
		}
		if move.From == move.To || addressDescendant(move.To, move.From) || addressDescendant(move.From, move.To) {
			return nil, fmt.Errorf("move paths overlap: %q -> %q", move.From, move.To)
		}
		if resourceRoot(move.From) != move.From || resourceRoot(move.To) != move.To {
			if err := moveMember(result, move); err != nil {
				return nil, fmt.Errorf("move note-type member: %w", err)
			}
			continue
		}
		occupied := map[string]bool{}
		for address := range result.Resources.Decks {
			occupied[address] = true
		}
		for address := range result.Resources.Notes {
			occupied[address] = true
		}
		for address := range result.Resources.NoteTypes {
			occupied[address] = true
		}
		if err := moveBindings(result.Resources.Decks, move, occupied); err != nil {
			return nil, fmt.Errorf("move deck state: %w", err)
		}
		if err := moveBindings(result.Resources.Notes, move, occupied); err != nil {
			return nil, fmt.Errorf("move note state: %w", err)
		}
		if err := moveBindings(result.Resources.NoteTypes, move, occupied); err != nil {
			return nil, fmt.Errorf("move note-type state: %w", err)
		}
	}
	return result, nil
}

// memberAddress separates a declared member move into its file, category, and exact name.
func memberAddress(address string) (string, string, string) {
	for _, category := range []string{"fields", "templates"} {
		file, name, found := strings.Cut(address, ".notetype.anki.toml#"+category+"/")
		if found {
			return file + ".notetype.anki.toml", category, name
		}
	}
	return address, "", ""
}

// moveBindings moves file or directory keys without changing the typed values.
func moveBindings[T any](bindings map[string]T, move *collectionpb.ResourceMove, occupied map[string]bool) error {
	replacements := map[string]T{}
	for address, binding := range bindings {
		if addressDescendant(address, move.From) {
			destination := move.To + strings.TrimPrefix(address, move.From)
			if occupied[destination] {
				return fmt.Errorf("move destination %q is already bound", destination)
			}
			replacements[destination] = binding
		}
	}
	for address := range bindings {
		if addressDescendant(address, move.From) {
			delete(bindings, address)
		}
	}
	for address, binding := range replacements {
		bindings[address] = binding
	}
	return nil
}

// moveMember renames a field or template while retaining its native identity and ordinal.
func moveMember(value *collectionpb.State, move *collectionpb.ResourceMove) error {
	fromFile, fromCategory, fromName := memberAddress(move.From)
	toFile, toCategory, toName := memberAddress(move.To)
	if fromCategory == "" || fromCategory != toCategory {
		return fmt.Errorf("member moves require matching field or template endpoints")
	}
	from := value.Resources.NoteTypes[fromFile]
	if from == nil {
		return nil
	}
	to := value.Resources.NoteTypes[toFile]
	if fromCategory == "fields" {
		if from.Fields[fromName] == nil {
			return nil
		}
		if to == nil {
			return fmt.Errorf("member destination %q must be a bound note type", toFile)
		}
		if to.Fields == nil {
			to.Fields = map[string]*collectionpb.FieldState{}
		}
		if err := moveNamedMember(from.Fields, to.Fields, fromName, toName); err != nil {
			return fmt.Errorf("move field: %w", err)
		}
	} else {
		if from.Templates[fromName] == nil {
			return nil
		}
		if to == nil {
			return fmt.Errorf("member destination %q must be a bound note type", toFile)
		}
		if to.Templates == nil {
			to.Templates = map[string]*collectionpb.TemplateState{}
		}
		if err := moveNamedMember(from.Templates, to.Templates, fromName, toName); err != nil {
			return fmt.Errorf("move template: %w", err)
		}
	}
	return nil
}

// moveNamedMember transfers one typed member between note-type maps without overwriting a destination.
func moveNamedMember[T any](source, destination map[string]T, from, to string) error {
	if _, exists := destination[to]; exists {
		return fmt.Errorf("member destination %q is already bound", to)
	}
	destination[to] = source[from]
	delete(source, from)
	return nil
}

// resourceRoot separates virtual type members while allowing # in literal filesystem names.
func resourceRoot(address string) string {
	file, _, _ := memberAddress(address)
	return file
}

// validateMoveSources prevents retained moves from claiming a path that is also newly declared.
func validateMoveSources(configs *collectionpb.CollectionConfigSet) error {
	addresses := []string{}
	for _, deck := range configs.Decks {
		addresses = append(addresses, deck.Path)
		for _, note := range deck.Notes {
			addresses = append(addresses, note.Path)
		}
	}
	for _, note := range configs.Notes {
		addresses = append(addresses, note.Path)
	}
	for _, definition := range configs.NoteTypes {
		addresses = append(addresses, definition.Path)
		for _, field := range definition.Config.Fields {
			addresses = append(addresses, definition.Path+"#fields/"+field.Name)
		}
		for _, template := range definition.Config.Templates {
			addresses = append(addresses, definition.Path+"#templates/"+template.Name)
		}
	}
	for _, move := range configs.GetConfig().GetMoved() {
		for _, address := range addresses {
			if addressDescendant(address, move.From) {
				return fmt.Errorf("move source %q is still declared; remove its declaration or the completed move", move.From)
			}
		}
	}
	return nil
}

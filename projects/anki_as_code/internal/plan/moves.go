package plan

import (
	"fmt"
	"io/fs"
	"strings"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"google.golang.org/protobuf/proto"
)

// validMoveAddress validates the filesystem portion of a resource or member address.
func validMoveAddress(address string) error {
	file := resourceRoot(address)
	if file == "." || !fs.ValidPath(file) || strings.ContainsAny(file, "\\\x00") {
		return fmt.Errorf("invalid collection-relative move address %q", address)
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
	result := proto.Clone(state).(*collectionpb.State)
	if result.Resources == nil {
		result.Resources = &collectionpb.Resources{}
	}
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

// moveMember renames a field or template within its original native note type.
func moveMember(value *collectionpb.State, move *collectionpb.ResourceMove) error {
	fromFile, fromCategory, fromName := memberAddress(move.From)
	toFile, toCategory, toName := memberAddress(move.To)
	if fromFile != toFile || fromCategory == "" || fromCategory != toCategory {
		return fmt.Errorf("member moves require matching field or template endpoints within the same note type")
	}
	binding := value.Resources.NoteTypes[fromFile]
	if binding == nil {
		return nil
	}
	if fromCategory == "fields" {
		if err := renameMember(binding.Fields, fromName, toName); err != nil {
			return fmt.Errorf("move field: %w", err)
		}
	} else {
		if err := renameMember(binding.Templates, fromName, toName); err != nil {
			return fmt.Errorf("move template: %w", err)
		}
	}
	return nil
}

// renameMember changes one map key while preserving its typed identity and ordinal.
// An absent source makes a retained move idempotent; existing destinations cannot be overwritten.
func renameMember[T any](bindings map[string]T, from, to string) error {
	value, exists := bindings[from]
	if !exists {
		return nil
	}
	if _, exists := bindings[to]; exists {
		return fmt.Errorf("member destination %q is already bound", to)
	}
	bindings[to] = value
	delete(bindings, from)
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

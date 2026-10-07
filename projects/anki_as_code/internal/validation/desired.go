package validation

import (
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/ankiformat"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/model"
	"github.com/ankitects/anki/proto/anki/notetypes"
)

// clozePattern matches positive cloze deletion numbers in note field content.
var clozePattern = regexp.MustCompile(`\{\{c([1-9][0-9]*)::`)

// clozes returns the distinct cloze deletion numbers found in note fields.
func clozes(bodies []string) []string {
	set := map[string]bool{}
	for _, body := range bodies {
		for _, m := range clozePattern.FindAllStringSubmatch(body, -1) {
			set[m[1]] = true
		}
	}
	result := []string{}
	for ord := range set {
		result = append(result, ord)
	}
	sort.Strings(result)
	return result
}

// ValidateDesired validates desired resources against the current collection snapshot.
func ValidateDesired(d *model.CollectionDesiredState, s model.CollectionState) error {
	if err := validateAppearance(d, s); err != nil {
		return fmt.Errorf("validate card appearance: %w", err)
	}
	deckNames, err := validateDecks(d, s)
	if err != nil {
		return fmt.Errorf("validate desired decks: %w", err)
	}
	normalDecks := make(map[string]bool, len(d.Decks))
	for _, deck := range d.Decks {
		normalDecks[ankiformat.Folded(deck.Title)] = deck.Normal != nil
	}
	currentCards := map[int64]int64{}
	for id, n := range s.Notes {
		for _, c := range n.Cards {
			currentCards[c.Id] = id
		}
	}
	guids := map[string]bool{}
	cards := map[int64]bool{}
	for i := range d.Notes {
		n := &d.Notes[i]
		nt := s.Types[n.NoteTypeId]
		if nt == nil || nt.Name != n.NoteType {
			return fmt.Errorf("note %d has unknown/mismatched note type", n.NoteId)
		}
		if n.Guid == "" || guids[n.Guid] {
			return fmt.Errorf("note %d has empty/duplicate GUID", n.NoteId)
		}
		guids[n.Guid] = true
		if len(n.Fields) != len(nt.Fields) {
			return fmt.Errorf("note %d field count differs from note type", n.NoteId)
		}
		orderedFields := make([]string, len(nt.Fields))
		orderedBodies := make([]string, len(nt.Fields))
		for j, field := range nt.Fields {
			name := field.Name
			body, ok := n.Fields[name]
			if !ok {
				return fmt.Errorf("note %d is missing field %q", n.NoteId, name)
			}
			orderedFields[j] = name
			orderedBodies[j] = body
		}
		n.OrderedFields = orderedFields
		n.Bodies = orderedBodies
		if err := ValidateBodies(*n); err != nil {
			return fmt.Errorf("validate note: %w", err)
		}
		noteClozes := clozes(n.Bodies)
		if old, ok := s.Notes[n.NoteId]; ok {
			if old.Guid != n.Guid || old.NoteTypeId != n.NoteTypeId {
				return fmt.Errorf("note %d identity/type cannot change", n.NoteId)
			}
			if !slices.Equal(clozes(old.Bodies), noteClozes) {
				return fmt.Errorf("note %d changes cloze card ordinals; declare a new note instead", n.NoteId)
			}
		}
		ordinals := map[int64]bool{}
		for _, c := range n.Cards {
			if c.Id <= 0 || cards[c.Id] || c.Ordinal < 0 || c.Ordinal > int64(^uint(0)>>1) || ordinals[c.Ordinal] {
				return fmt.Errorf("note %d has invalid/duplicate card identity or ordinal", n.NoteId)
			}
			cards[c.Id] = true
			ordinals[c.Ordinal] = true
			if owner, ok := currentCards[c.Id]; ok {
				if owner != n.NoteId {
					return fmt.Errorf("card %d cannot change owning note", c.Id)
				}
				for _, old := range s.Notes[owner].Cards {
					if old.Id == c.Id && old.Ordinal != c.Ordinal {
						return fmt.Errorf("card %d cannot change template ordinal; declare a new card id", c.Id)
					}
				}
			}
			if nt.Config.Kind == notetypes.Notetype_Config_KIND_CLOZE {
				found := false
				for _, ord := range noteClozes {
					v, err := strconv.Atoi(ord)
					if err != nil {
						return fmt.Errorf("parse cloze ordinal: %w", err)
					}
					if int64(v)-1 == c.Ordinal {
						found = true
					}
				}
				if !found {
					return fmt.Errorf("card %d has no matching cloze ordinal", c.Id)
				}
			}
			if !(nt.Config.Kind == notetypes.Notetype_Config_KIND_CLOZE) && !nt.Ordinals[int(c.Ordinal)] {
				return fmt.Errorf("card %d has no template at ordinal %d", c.Id, c.Ordinal)
			}
			if !deckNames[ankiformat.Folded(c.GetDeck())] {
				return fmt.Errorf("card %d references undeclared deck %q", c.Id, c.GetDeck())
			}
			if !normalDecks[ankiformat.Folded(c.GetDeck())] {
				return fmt.Errorf("card %d requires a normal home deck; %q is filtered", c.Id, c.GetDeck())
			}
		}
		for _, tag := range n.Tags {
			if tag == "" || strings.ContainsAny(tag, " \t\r\n\x00") {
				return fmt.Errorf("note %d has invalid tag", n.NoteId)
			}
		}
		if int(nt.Config.SortFieldIdx) >= len(n.Bodies) {
			return fmt.Errorf("note type %d has invalid sort field", n.NoteTypeId)
		}
	}

	return nil
}

// validateDecks validates deck identities, titles, parents, and deck-kind transitions.
func validateDecks(d *model.CollectionDesiredState, s model.CollectionState) (map[string]bool, error) {
	deckNames := map[string]bool{}
	deckIDs := map[int64]bool{}
	existingDecks := map[int64]model.Deck{}
	filteredDecks := map[int64]bool{}
	for _, deck := range s.Decks {
		if err := ankiformat.DecodeDeck(&deck); err != nil {
			return nil, fmt.Errorf("decode current deck %q: %w", deck.Title, err)
		}
		existingDecks[deck.Id] = deck
		filteredDecks[deck.Id] = deck.Normal == nil
	}
	maxID := int64(0)
	for id := range existingDecks {
		if id > maxID {
			maxID = id
		}
	}
	for _, deck := range d.Decks {
		if deck.Id > maxID {
			maxID = deck.Id
		}
	}
	sort.Slice(d.Decks, func(i, j int) bool {
		return d.Decks[i].Title < d.Decks[j].Title
	})
	for i := range d.Decks {
		deck := &d.Decks[i]
		if err := ValidDeck(deck.Title); err != nil {
			return nil, fmt.Errorf("validate deck: %w", err)
		}
		if deckNames[ankiformat.Folded(deck.Title)] {
			return nil, fmt.Errorf("duplicate deck name %q", deck.Title)
		}
		deckNames[ankiformat.Folded(deck.Title)] = true
		if _, err := ankiformat.DecodeDeckKind(deck.Kind); err != nil {
			return nil, fmt.Errorf("validate deck %q kind: %w", deck.Title, err)
		}
		if deck.Id == 0 {
			for _, old := range s.Decks {
				if ankiformat.Folded(old.Title) == ankiformat.Folded(deck.Title) {
					deck.Id = old.Id
				}
			}
			if deck.Id == 0 {
				maxID++
				deck.Id = maxID
			}
		}
		if deck.Id <= 0 || deckIDs[deck.Id] {
			return nil, fmt.Errorf("invalid or duplicate deck id %d", deck.Id)
		}
		if deck.Normal != nil && filteredDecks[deck.Id] && s.OccupiedDecks[deck.Id] {
			return nil, fmt.Errorf("cannot convert occupied filtered deck %q to normal; empty it in Anki first", deck.Title)
		}
		if current, exists := existingDecks[deck.Id]; exists {
			if err := ankiformat.PreserveDeckKind(current, deck); err != nil {
				return nil, fmt.Errorf("preserve deck %q settings: %w", deck.Title, err)
			}
		}
		deckIDs[deck.Id] = true
	}
	for _, deck := range d.Decks {
		parts := strings.Split(deck.Title, "::")
		for i := 1; i < len(parts); i++ {
			if !deckNames[ankiformat.Folded(strings.Join(parts[:i], "::"))] {
				return nil, fmt.Errorf("deck %q lacks declared parent", deck.Title)
			}
		}
	}
	return deckNames, nil
}

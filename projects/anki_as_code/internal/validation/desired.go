package validation

import (
	"fmt"
	"math"
	"slices"
	"sort"
	"strings"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/ankiformat"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/model"
	"github.com/ankitects/anki/proto/anki/notetypes"
)

// ValidateDesired resolves declarations and validates them against current identities.
// Cards are explicit resources; note HTML and cloze content are not interpreted.
func ValidateDesired(desired *model.CollectionDesiredState, current model.CollectionState) error {
	if err := validateDecks(desired, current); err != nil {
		return fmt.Errorf("validate desired decks: %w", err)
	}
	if err := validateNoteTypes(desired, current); err != nil {
		return fmt.Errorf("validate note types: %w", err)
	}
	if err := validateNotes(desired, current); err != nil {
		return fmt.Errorf("validate notes: %w", err)
	}
	return nil
}

// noteValidator keeps cross-note identity indexes local to a validation pass.
type noteValidator struct {
	// current contains existing note identities and managed card ownership.
	current model.CollectionState
	// types indexes the desired note-type declarations.
	types map[int64]*collectionpb.NoteTypeSettings
	// decks indexes desired deck declarations by normalized title.
	decks map[string]model.Deck
	// owners indexes existing managed cards by their owning note.
	owners map[int64]int64
	// notes records note identities already encountered.
	notes map[int64]bool
	// cards records card identities already encountered.
	cards map[int64]bool
	// guids records note GUIDs already encountered.
	guids map[string]bool
}

// validateNotes checks note declarations using the resolved desired schema and decks.
func validateNotes(desired *model.CollectionDesiredState, current model.CollectionState) error {
	validator := noteValidator{current: current, types: map[int64]*collectionpb.NoteTypeSettings{}, decks: map[string]model.Deck{}, owners: map[int64]int64{}, notes: map[int64]bool{}, cards: map[int64]bool{}, guids: map[string]bool{}}
	for _, settings := range desired.NoteTypes {
		validator.types[settings.Id] = settings
	}
	for _, deck := range desired.Decks {
		validator.decks[ankiformat.Folded(deck.Title)] = deck
	}
	for id, note := range current.Notes {
		for _, card := range note.Cards {
			validator.owners[card.Id] = id
		}
	}
	for i := range desired.Notes {
		if err := validator.validateNote(&desired.Notes[i]); err != nil {
			return fmt.Errorf("validate note declaration: %w", err)
		}
	}
	return nil
}

// validateNote checks identity, field layout, tags, and explicit managed cards.
func (validator *noteValidator) validateNote(note *model.Note) error {
	settings, err := validator.validateNoteIdentity(note)
	if err != nil {
		return fmt.Errorf("note identity: %w", err)
	}
	if err := validateNoteFields(note, settings); err != nil {
		return fmt.Errorf("validate note %d fields: %w", note.NoteId, err)
	}
	ordinals := map[int64]bool{}
	for _, card := range note.Cards {
		if err := validator.validateCard(note.NoteId, card, settings, ordinals); err != nil {
			return fmt.Errorf("validate note %d card: %w", note.NoteId, err)
		}
	}
	for _, tag := range note.Tags {
		if tag == "" || strings.ContainsAny(tag, " \t\r\n\x00") {
			return fmt.Errorf("note %d has invalid tag", note.NoteId)
		}
	}
	return nil
}

// validateNoteFields derives SQLite field order from the desired note-type definition.
func validateNoteFields(note *model.Note, settings *collectionpb.NoteTypeSettings) error {
	if len(note.Fields) != len(settings.Fields) {
		return fmt.Errorf("field count differs from desired note type")
	}
	note.Bodies = make([]string, len(settings.Fields))
	for i, field := range settings.Fields {
		body, exists := note.Fields[field.Name]
		if !exists {
			return fmt.Errorf("missing field %q", field.Name)
		}
		note.Bodies[i] = body
	}
	if err := validateBodies(*note); err != nil {
		return fmt.Errorf("field content: %w", err)
	}
	return nil
}

// validateCard checks uniqueness, stable ownership, template ordinal, and normal deck assignment.
func (validator *noteValidator) validateCard(noteID int64, card *collectionpb.Card, settings *collectionpb.NoteTypeSettings, ordinals map[int64]bool) error {
	if card == nil || card.Id <= 0 || validator.cards[card.Id] || card.Ordinal < 0 || card.Ordinal > math.MaxInt32 || ordinals[card.Ordinal] {
		return fmt.Errorf("invalid/duplicate card identity or ordinal")
	}
	validator.cards[card.Id], ordinals[card.Ordinal] = true, true
	if owner, exists := validator.owners[card.Id]; exists && owner != noteID {
		return fmt.Errorf("card %d cannot change owning note", card.Id)
	}
	if settings.Config.Kind != notetypes.Notetype_Config_KIND_CLOZE && !slices.ContainsFunc(settings.Templates, func(template *collectionpb.TemplateSettings) bool { return template.Ordinal == card.Ordinal }) {
		return fmt.Errorf("card %d has no template at ordinal %d", card.Id, card.Ordinal)
	}
	deck, exists := validator.decks[ankiformat.Folded(card.GetDeck())]
	if !exists {
		return fmt.Errorf("card %d references undeclared deck %q", card.Id, card.GetDeck())
	}
	if deck.Normal == nil {
		return fmt.Errorf("card %d requires a normal home deck; %q is filtered", card.Id, card.GetDeck())
	}
	return nil
}

// validateDecks resolves deck identities and validates configuration and hierarchy.
func validateDecks(desired *model.CollectionDesiredState, current model.CollectionState) error {
	existing, nextID := map[int64]model.Deck{}, int64(0)
	for _, deck := range current.Decks {
		existing[deck.Id] = deck
		nextID = max(nextID, deck.Id)
	}
	for _, deck := range desired.Decks {
		if deck.Deck == nil {
			return fmt.Errorf("nil deck declaration")
		}
		nextID = max(nextID, deck.Id)
	}
	names, ids := map[string]bool{}, map[int64]bool{}
	sort.Slice(desired.Decks, func(i, j int) bool { return desired.Decks[i].Title < desired.Decks[j].Title })
	for i := range desired.Decks {
		deck := &desired.Decks[i]
		if err := resolveDeckID(deck, current.Decks, &nextID); err != nil {
			return fmt.Errorf("resolve deck %q: %w", deck.Title, err)
		}
		if err := validateDeckDeclaration(deck, names, ids); err != nil {
			return fmt.Errorf("validate deck %q: %w", deck.Title, err)
		}
		if old, exists := existing[deck.Id]; exists {
			if err := validateDeckTransition(old, deck, current.OccupiedDecks[deck.Id]); err != nil {
				return fmt.Errorf("update deck %q: %w", deck.Title, err)
			}
		}
	}
	if err := validateDeckParents(desired.Decks, names); err != nil {
		return fmt.Errorf("validate deck hierarchy: %w", err)
	}
	return nil
}

// resolveDeckID matches omitted identities by title or allocates after all reserved identities.
func resolveDeckID(deck *model.Deck, current []model.Deck, nextID *int64) error {
	if deck.Id != 0 {
		return nil
	}
	for _, old := range current {
		if ankiformat.Folded(old.Title) == ankiformat.Folded(deck.Title) {
			deck.Id = old.Id
			return nil
		}
	}
	if *nextID == math.MaxInt64 {
		return fmt.Errorf("deck identity range exhausted")
	}
	*nextID += 1
	deck.Id = *nextID
	return nil
}

// validateDeckDeclaration checks unique identity/title and a valid encoded deck kind.
func validateDeckDeclaration(deck *model.Deck, names map[string]bool, ids map[int64]bool) error {
	if err := ValidDeck(deck.Title); err != nil {
		return fmt.Errorf("deck title: %w", err)
	}
	name := ankiformat.Folded(deck.Title)
	if names[name] {
		return fmt.Errorf("duplicate deck name")
	}
	if deck.Id <= 0 || ids[deck.Id] {
		return fmt.Errorf("invalid or duplicate deck id %d", deck.Id)
	}
	if _, err := ankiformat.DecodeDeckKind(deck.Kind); err != nil {
		return fmt.Errorf("deck kind: %w", err)
	}
	names[name], ids[deck.Id] = true, true
	return nil
}

// validateDeckTransition protects resident filtered cards and preserves opaque configuration.
func validateDeckTransition(current model.Deck, desired *model.Deck, occupied bool) error {
	kind, err := ankiformat.DecodeDeckKind(current.Kind)
	if err != nil {
		return fmt.Errorf("read current kind: %w", err)
	}
	if kind.GetFiltered() != nil && desired.Normal != nil && occupied {
		return fmt.Errorf("cannot convert occupied filtered deck to normal; empty it in Anki first")
	}
	if err := ankiformat.PreserveDeckKind(current, desired); err != nil {
		return fmt.Errorf("preserve settings: %w", err)
	}
	return nil
}

// validateDeckParents requires an explicit declaration for every hierarchy ancestor.
func validateDeckParents(decks []model.Deck, names map[string]bool) error {
	for _, deck := range decks {
		parts := strings.Split(deck.Title, "::")
		for i := 1; i < len(parts); i++ {
			if !names[ankiformat.Folded(strings.Join(parts[:i], "::"))] {
				return fmt.Errorf("deck %q lacks declared parent", deck.Title)
			}
		}
	}
	return nil
}

// validateNoteIdentity checks uniqueness and stable identity, and resolves the desired note type.
func (validator *noteValidator) validateNoteIdentity(note *model.Note) (*collectionpb.NoteTypeSettings, error) {
	if note.Note == nil || note.NoteId <= 0 || validator.notes[note.NoteId] {
		return nil, fmt.Errorf("invalid or duplicate note identity")
	}
	validator.notes[note.NoteId] = true
	if note.Guid == "" || validator.guids[note.Guid] {
		return nil, fmt.Errorf("note %d has empty/duplicate GUID", note.NoteId)
	}
	validator.guids[note.Guid] = true
	if note.NoteTypeId == 0 {
		for id, settings := range validator.types {
			if settings.Name == note.NoteType {
				note.NoteTypeId = id
			}
		}
	}
	settings := validator.types[note.NoteTypeId]
	if settings == nil || settings.Name != note.NoteType {
		return nil, fmt.Errorf("note %d has unknown/mismatched desired note type", note.NoteId)
	}
	if old, exists := validator.current.Notes[note.NoteId]; exists && (old.Guid != note.Guid || old.NoteTypeId != note.NoteTypeId) {
		return nil, fmt.Errorf("note %d identity/type cannot change", note.NoteId)
	}
	return settings, nil
}

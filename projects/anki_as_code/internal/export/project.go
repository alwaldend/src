package export

import (
	"cmp"
	"context"
	"fmt"
	"path"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/resource"
	"github.com/ankitects/anki/proto/anki/card_rendering"
	"github.com/ankitects/anki/proto/anki/cards"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"google.golang.org/protobuf/proto"
)

// maxNoteNameBytes leaves room for collision suffixes and the configuration extension.
const maxNoteNameBytes = 120

// invalidPathComponentCharacters contains separators and NUL excluded from exported filenames.
const invalidPathComponentCharacters = "/\\\x00"

// configProjection holds per-export inputs, outputs and cross-resource references.
type configProjection struct {
	// current contains the native snapshot projected by this export.
	current *resource.Snapshot
	// stored supplies stable filenames and scheduling preset aliases.
	stored *ankiascodepb.State
	// configs collects the editable configuration files.
	configs *ankiascodepb.CollectionConfigSet
	// identities records native IDs at the actual projected config paths.
	identities map[string]int64
	// decks resolves native deck IDs to their projected destinations.
	decks map[int64]*ankiascodepb.DeckConfigSet
	// typePaths resolves native note-type IDs to their projected filenames.
	typePaths map[int64]string
}

// exportConfigs projects native protobufs into editable configs using stored names.
func exportConfigs(
	ctx context.Context,
	current *resource.Snapshot,
	stored *ankiascodepb.State,
	rendering card_rendering.CardRenderingServiceClient,
) (*ankiascodepb.CollectionConfigSet, map[string]int64, error) {
	projection := configProjection{
		current:    current,
		stored:     stored,
		configs:    &ankiascodepb.CollectionConfigSet{Config: &ankiascodepb.CollectionConfig{Version: config.Version}},
		decks:      map[int64]*ankiascodepb.DeckConfigSet{},
		typePaths:  map[int64]string{},
		identities: map[string]int64{},
	}
	if err := projection.projectDecks(); err != nil {
		return nil, nil, fmt.Errorf("project decks: %w", err)
	}
	if err := projection.projectNoteTypes(); err != nil {
		return nil, nil, fmt.Errorf("project note types: %w", err)
	}
	if err := projection.projectNotes(ctx, rendering); err != nil {
		return nil, nil, fmt.Errorf("project notes: %w", err)
	}
	return projection.configs, projection.identities, nil
}

// projectDecks assigns literal deck paths and resolves editable preset references.
func (projection *configProjection) projectDecks() error {
	deckPaths := map[string]bool{}
	presetIDs := map[string]int64{}
	presetAliases := map[int64]string{}
	for alias, preset := range projection.stored.GetResources().GetPresets() {
		id := preset.GetId()
		if projection.current.Presets[id] == "" {
			continue
		}
		presetIDs[alias] = id
		if previous := presetAliases[id]; previous == "" || alias < previous {
			presetAliases[id] = alias
		}
	}
	for _, deck := range projection.current.Decks {
		directory := ""
		for _, part := range strings.Split(deck.Name, "::") {
			if part == "" || part == "." || part == ".." || strings.ContainsAny(part, invalidPathComponentCharacters) {
				return fmt.Errorf("deck %q cannot be represented as literal directory names", deck.Name)
			}
			directory = path.Join(directory, "decks", part)
		}
		value := proto.Clone(&ankiascodepb.DeckConfig{Normal: deck.GetNormal(), Filtered: deck.GetFiltered()}).(*ankiascodepb.DeckConfig)
		if value.Normal != nil {
			presetID := value.Normal.ConfigId
			value.Preset = presetAliases[presetID]
			if value.Preset == "" {
				value.Preset = projection.current.Presets[presetID]
			}
			if value.Preset == "" {
				return fmt.Errorf("deck %q has an unavailable scheduling preset", deck.Name)
			}
			if previous := presetIDs[value.Preset]; previous != 0 && previous != presetID {
				return fmt.Errorf("scheduling preset name %q is ambiguous", value.Preset)
			}
			presetIDs[value.Preset] = presetID
			value.Normal.ConfigId = 0
		}
		group := &ankiascodepb.DeckConfigSet{Title: deck.Name, Path: path.Join(directory, "deck.anki.toml"), Deck: value}
		if deckPaths[group.Path] {
			return fmt.Errorf("export deck path %q represents multiple objects", group.Path)
		}
		deckPaths[group.Path] = true
		projection.configs.Decks = append(projection.configs.Decks, group)
		projection.decks[deck.Id] = group
		projection.identities[group.Path] = deck.Id
	}
	return nil
}

// projectNoteTypes reserves filenames and projects editable fields and templates.
func (projection *configProjection) projectNoteTypes() error {
	typePathsUsed := map[string]bool{}
	typeNames := map[string]bool{}
	for address, binding := range projection.stored.GetResources().GetNoteTypes() {
		projection.typePaths[binding.Id] = address
	}
	for _, definition := range projection.current.NoteTypes {
		if address := projection.typePaths[definition.Id]; address != "" {
			typeNames[strings.ToLower(path.Base(address))] = true
		}
	}
	for _, definition := range projection.current.NoteTypes {
		address := projection.typePaths[definition.Id]
		filename := path.Base(address)
		if address == "" {
			name := definition.Name
			if strings.ContainsAny(name, invalidPathComponentCharacters) {
				name = fmt.Sprintf("type-%d", len(projection.configs.NoteTypes)+1)
			}
			for suffix := 1; ; suffix++ {
				filename = name + ".notetype.anki.toml"
				if suffix > 1 {
					filename = fmt.Sprintf("%s-%d.notetype.anki.toml", name, suffix)
				}
				key := strings.ToLower(filename)
				if !typeNames[key] {
					typeNames[key] = true
					break
				}
			}
		}
		address = path.Join("note_types", filename)
		projection.typePaths[definition.Id] = address
		if typePathsUsed[address] {
			return fmt.Errorf("export note-type path %q represents multiple objects", address)
		}
		typePathsUsed[address] = true
		projection.identities[address] = definition.Id
		copy := proto.Clone(definition).(*notetypes.Notetype)
		value := &ankiascodepb.NoteTypeConfig{Name: copy.Name, Config: copy.Config, Fields: copy.Fields}
		if int(value.Config.SortFieldIdx) >= len(value.Fields) {
			return fmt.Errorf("invalid sort field for type %q", value.Name)
		}
		if value.Config.SortFieldIdx != 0 {
			value.SortField = value.Fields[value.Config.SortFieldIdx].Name
		}
		value.Config.OriginalId = nil
		value.Config.TargetDeckIdUnused = 0
		value.Config.SortFieldIdx = 0
		value.Config.Reqs = nil
		value.Config.Other = nil
		for _, field := range value.Fields {
			field.Ord = nil
			field.Config.Id = nil
			field.Config.Other = nil
		}
		for _, template := range copy.Templates {
			entry := &ankiascodepb.TemplateConfig{Name: template.Name, Config: template.Config}
			if id := entry.Config.TargetDeckId; id != 0 {
				destination := projection.decks[id]
				if destination == nil || destination.Deck.Normal == nil {
					return fmt.Errorf("template %q has unavailable normal target deck", entry.Name)
				}
				entry.TargetDeck = destination.Path
			}
			entry.Config.Id = nil
			entry.Config.TargetDeckId = 0
			entry.Config.Other = nil
			value.Templates = append(value.Templates, entry)
		}
		projection.configs.NoteTypes = append(projection.configs.NoteTypes, &ankiascodepb.NoteTypeConfigFile{Path: address, Config: value})
	}
	slices.SortFunc(projection.configs.NoteTypes, func(a, b *ankiascodepb.NoteTypeConfigFile) int { return strings.Compare(a.Path, b.Path) })
	return nil
}

// noteCards groups managed regular cards by note in template order.
func (projection *configProjection) noteCards() (map[int64][]*cards.Card, error) {
	noteCards := map[int64][]*cards.Card{}
	for _, card := range projection.current.Cards {
		destination := projection.decks[card.DeckId]
		if destination == nil {
			return nil, fmt.Errorf("card %d has an unavailable deck", card.Id)
		}
		if card.OriginalDeckId != 0 || destination.Deck.Normal == nil {
			continue
		}
		noteCards[card.NoteId] = append(noteCards[card.NoteId], card)
	}
	for _, entries := range noteCards {
		slices.SortFunc(entries, func(a, b *cards.Card) int { return cmp.Compare(a.TemplateIdx, b.TemplateIdx) })
	}
	return noteCards, nil
}

// projectNotes preserves bound filenames and projects note fields and deck assignments.
func (projection *configProjection) projectNotes(ctx context.Context, rendering card_rendering.CardRenderingServiceClient) error {
	noteCards, err := projection.noteCards()
	if err != nil {
		return fmt.Errorf("resolve note cards: %w", err)
	}
	notePaths := map[int64]string{}
	for address, binding := range projection.stored.GetResources().GetNotes() {
		notePaths[binding.Id] = address
	}
	nativeTypes := map[int64]*notetypes.Notetype{}
	for _, definition := range projection.current.NoteTypes {
		nativeTypes[definition.Id] = definition
	}
	addresses := map[int64]string{}
	reserved := map[string]bool{}
	nextSuffix := map[string]int{}
	for _, note := range projection.current.Notes {
		if bound := notePaths[note.Id]; bound != "" {
			address := path.Join(exportNoteDirectory(exportNoteHome(noteCards[note.Id], projection.decks)), path.Base(bound))
			key := strings.ToLower(address)
			if !reserved[key] {
				addresses[note.Id] = address
				reserved[key] = true
			}
		}
	}
	for _, note := range projection.current.Notes {
		definition := nativeTypes[note.NotetypeId]
		if definition == nil || len(definition.Fields) != len(note.Fields) {
			return fmt.Errorf("note %d has an unavailable field layout", note.Id)
		}
		value := &ankiascodepb.NoteConfig{NoteType: projection.typePaths[note.NotetypeId], Tags: note.Tags, Fields: map[string]string{}}
		for index, field := range definition.Fields {
			value.Fields[field.Name] = note.Fields[index]
		}
		home := exportNoteHome(noteCards[note.Id], projection.decks)
		address := addresses[note.Id]
		if address == "" {
			var text string
			if bound := notePaths[note.Id]; bound != "" {
				text = strings.TrimSuffix(path.Base(bound), ".note.anki.toml")
			} else {
				plain, err := rendering.HtmlToTextLine(ctx, &card_rendering.HtmlToTextLineRequest{
					Text: note.Fields[definition.Config.SortFieldIdx],
				})
				if err != nil {
					return fmt.Errorf("convert sort field for note %d to text: %w", note.Id, err)
				}
				text = plain.Val
			}
			name := exportNoteName(text)
			directory := exportNoteDirectory(home)
			nameKey := strings.ToLower(path.Join(directory, name))
			for suffix := max(1, nextSuffix[nameKey]); ; suffix++ {
				filename := name
				if suffix > 1 {
					filename = fmt.Sprintf("%s-%d", name, suffix)
				}
				address = path.Join(directory, filename+".note.anki.toml")
				key := strings.ToLower(address)
				if !reserved[key] {
					reserved[key] = true
					nextSuffix[nameKey] = suffix + 1
					break
				}
			}
		}
		for _, card := range noteCards[note.Id] {
			assignment := &ankiascodepb.CardConfig{Ordinal: int64(card.TemplateIdx)}
			destination := projection.decks[card.DeckId]
			if home == nil || destination.Path != home.Path {
				assignment.Deck = proto.String(destination.Path)
			}
			value.Cards = append(value.Cards, assignment)
		}
		file := &ankiascodepb.NoteConfigFile{Path: address, Config: value}
		projection.identities[address] = note.Id
		if home == nil {
			projection.configs.Notes = append(projection.configs.Notes, file)
		} else {
			home.Notes = append(home.Notes, file)
		}
	}
	return nil
}

// exportNoteHome selects the first regular card's destination from ordinal-sorted cards.
func exportNoteHome(cards []*cards.Card, decks map[int64]*ankiascodepb.DeckConfigSet) *ankiascodepb.DeckConfigSet {
	if len(cards) == 0 {
		return nil
	}
	return decks[cards[0].DeckId]
}

// exportNoteDirectory places cardless notes at the root and others in their first card's deck.
func exportNoteDirectory(home *ankiascodepb.DeckConfigSet) string {
	if home == nil {
		return "notes"
	}
	return path.Join(path.Dir(home.Path), "notes")
}

// exportNoteName makes plain text a bounded Unicode filename without modifying note content.
func exportNoteName(text string) string {
	var name strings.Builder
	separator := false
	for _, character := range text {
		if !unicode.IsLetter(character) && !unicode.IsNumber(character) && !unicode.IsMark(character) {
			separator = name.Len() > 0
			continue
		}
		length := utf8.RuneLen(character)
		if separator {
			length++
		}
		if name.Len()+length > maxNoteNameBytes {
			break
		}
		if separator {
			name.WriteByte('-')
		}
		name.WriteRune(character)
		separator = false
	}
	label := name.String()
	if label == "" {
		return "note"
	}
	upper := strings.ToUpper(label)
	if upper == "CON" || upper == "PRN" || upper == "AUX" || upper == "NUL" ||
		(len(upper) == 4 && (strings.HasPrefix(upper, "COM") || strings.HasPrefix(upper, "LPT")) && upper[3] >= '1' && upper[3] <= '9') {
		return "note-" + label
	}
	return label
}

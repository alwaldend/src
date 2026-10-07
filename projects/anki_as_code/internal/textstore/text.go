package textstore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/ankiformat"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/model"
	"github.com/BurntSushi/toml"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// readFile reads declaration bytes and reports read and close failures.
func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open text %q: %w", path, err)
	}
	b, readErr := io.ReadAll(f)
	closeErr := f.Close()
	if err := errors.Join(readErr, closeErr); err != nil {
		return nil, fmt.Errorf("read text %q: %w", path, err)
	}
	return b, nil
}

// DecodeTOML decodes TOML values into the supplied protobuf message.
func DecodeTOML(b []byte, message proto.Message) error {
	if message == nil || !message.ProtoReflect().IsValid() {
		return fmt.Errorf("nil protobuf input message")
	}
	var values map[string]any
	if _, err := toml.Decode(string(b), &values); err != nil {
		return fmt.Errorf("decode TOML: %w", err)
	}
	// Older root configs carry a checksum that does not constrain the archive.
	if _, config := message.(*collectionpb.CollectionConfig); config {
		if _, checksum := values["base_sha256"].(string); checksum {
			delete(values, "base_sha256")
		}
	}
	data, err := json.Marshal(values)
	if err != nil {
		return fmt.Errorf("encode TOML values as JSON: %w", err)
	}
	if err := protojson.Unmarshal(data, message); err != nil {
		return fmt.Errorf("parse TOML with protobuf schema: %w", err)
	}
	return nil
}

// ReadTOML reads a TOML declaration and validates it through its protobuf schema.
func ReadTOML(path string, message proto.Message) error {
	b, err := readFile(path)
	if err != nil {
		return fmt.Errorf("read TOML input: %w", err)
	}
	if err := DecodeTOML(b, message); err != nil {
		return fmt.Errorf("parse %q: %w", path, err)
	}
	return nil
}

// encodeTOML encodes a value tree with the configured TOML encoder.
func encodeTOML(v any) ([]byte, error) {
	var b bytes.Buffer
	encoder := toml.NewEncoder(&b)
	encoder.Indent = ""
	if err := encoder.Encode(v); err != nil {
		return nil, fmt.Errorf("encode TOML: %w", err)
	}
	return b.Bytes(), nil
}

// WriteTOML serializes a protobuf message to a TOML declaration file.
func WriteTOML(path string, message proto.Message) error {
	values, err := tomlValues(message)
	if err != nil {
		return fmt.Errorf("project config %q: %w", path, err)
	}
	b, err := encodeTOML(values)
	if err != nil {
		return fmt.Errorf("serialize config %q: %w", path, err)
	}
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return fmt.Errorf("write config %q: %w", path, err)
	}
	return nil
}

// tomlValues uses ProtoJSON for field names and presence, retaining numeric TOML integers.
func tomlValues(message proto.Message) (map[string]any, error) {
	if message == nil || !message.ProtoReflect().IsValid() {
		return nil, fmt.Errorf("nil protobuf output message")
	}
	data, err := (protojson.MarshalOptions{UseProtoNames: true, EmitDefaultValues: true}).Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("encode protobuf JSON: %w", err)
	}
	var values map[string]any
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&values); err != nil {
		return nil, fmt.Errorf("decode protobuf JSON values: %w", err)
	}
	if err := tomlIntegers(values, message.ProtoReflect().Descriptor()); err != nil {
		return nil, fmt.Errorf("convert protobuf integers for TOML: %w", err)
	}
	return values, nil
}

// tomlIntegers converts integer fields in a ProtoJSON object without changing text content.
func tomlIntegers(values map[string]any, schema protoreflect.MessageDescriptor) error {
	for name, value := range values {
		field := schema.Fields().ByName(protoreflect.Name(name))

		switch {
		case field.IsMap():
			for key, entry := range value.(map[string]any) {
				converted, err := tomlScalar(entry, field.MapValue())
				if err != nil {
					return fmt.Errorf("convert map field %s key %q: %w", field.FullName(), key, err)
				}
				value.(map[string]any)[key] = converted
			}
		case field.IsList():
			for index, entry := range value.([]any) {
				converted, err := tomlScalar(entry, field)
				if err != nil {
					return fmt.Errorf("convert list field %s index %d: %w", field.FullName(), index, err)
				}
				value.([]any)[index] = converted
			}
		default:
			converted, err := tomlScalar(value, field)
			if err != nil {
				return fmt.Errorf("convert field %s: %w", field.FullName(), err)
			}
			values[name] = converted
		}
	}
	return nil
}

// DeckPath maps an Anki deck hierarchy to an escaped export directory path.
func DeckPath(dir, deck string) string {
	parts := strings.Split(deck, "::")
	for i, p := range parts {
		parts[i] = escapeDeckPart(p)
	}
	return filepath.Join(append([]string{dir, "decks"}, parts...)...)
}

// escapeDeckPart escapes characters that cannot safely identify a deck directory.
func escapeDeckPart(part string) string {
	if part == "notes" || part == DeckMarker {
		return fmt.Sprintf("%%%02X%s", part[0], part[1:])
	}
	var b strings.Builder
	for i, r := range part {
		if strings.ContainsRune(`/\%<>:"|?*`, r) || r < 32 || r == 127 || ((r == '.' || r == ' ') && i == len(part)-1) || part == "." || part == ".." {
			for _, v := range []byte(string(r)) {
				fmt.Fprintf(&b, "%%%02X", v)
			}
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// WriteDeck writes a deck marker and its configured note location.
func WriteDeck(dir string, deck model.Deck) error {
	deck.Deck = proto.Clone(deck.Deck).(*collectionpb.Deck)
	if err := ankiformat.DecodeDeck(&deck); err != nil {
		return fmt.Errorf("decode deck settings: %w", err)
	}
	path := DeckPath(dir, deck.Title)
	if err := os.MkdirAll(filepath.Join(path, "notes"), 0o755); err != nil {
		return fmt.Errorf("create deck directory: %w", err)
	}
	deck.NotesPath = "notes"
	if err := WriteTOML(filepath.Join(path, DeckMarker), deck.Deck); err != nil {
		return fmt.Errorf("write deck %q: %w", deck.Title, err)
	}
	return nil
}

// WriteNote writes a note declaration with named fields and explicit cards.
func WriteNote(path string, n model.Note) error {
	n.Note = proto.Clone(n.Note).(*collectionpb.Note)
	for i := range n.Cards {
		if n.Cards[i].GetDeck() == n.HomeDeck {
			n.Cards[i].Deck = nil
		}
	}
	fields := n.Fields
	if fields == nil {
		fields = map[string]string{}
		for i, field := range n.OrderedFields {
			fields[field] = n.Bodies[i]
		}
	}
	n.Fields = fields
	if err := WriteTOML(path, n.Note); err != nil {
		return fmt.Errorf("write note %d: %w", n.NoteId, err)
	}
	return nil
}

// ReadNote reads a note declaration and records its source path.
func ReadNote(path string) (model.Note, error) {
	n := model.Note{Note: &collectionpb.Note{}}
	if err := ReadTOML(path, n.Note); err != nil {
		return n, fmt.Errorf("read note %q: %w", path, err)
	}
	return n, nil
}

// Load discovers the complete desired collection state from configured declarations.
func Load(input string) (model.CollectionDesiredState, error) {
	d := model.CollectionDesiredState{Config: &collectionpb.CollectionConfig{}}
	config, err := collectionConfig(input)
	if err != nil {
		return d, fmt.Errorf("resolve collection config: %w", err)
	}
	if err := ReadTOML(config, d.Config); err != nil {
		return d, fmt.Errorf("load manifest: %w", err)
	}
	if d.Config.FormatVersion != 2 {
		return d, fmt.Errorf("unsupported text format %d", d.Config.FormatVersion)
	}
	noteTypesPath, err := referencePath(config, d.Config.NoteTypesPath)
	if err != nil {
		return d, fmt.Errorf("resolve note types: %w", err)
	}
	appearance, err := readNoteTypes(noteTypesPath)
	if err != nil {
		return d, fmt.Errorf("load note types: %w", err)
	}
	d.Appearance = appearance
	decks, err := discoverDecks(config, d.Config)
	if err != nil {
		return d, fmt.Errorf("discover decks: %w", err)
	}
	d.Decks = decks
	seen := map[int64]bool{}
	owners := map[string]string{}
	loadNotes := func(path, home string) error {
		files, err := discoverFiles(path, noteSuffix, true)
		if err != nil {
			return fmt.Errorf("discover notes: %w", err)
		}
		for _, file := range files {
			if owner, exists := owners[file]; exists {
				return fmt.Errorf("note %q belongs to both %q and %q", file, owner, home)
			}
			n, err := ReadNote(file)
			if err != nil {
				return fmt.Errorf("load note text: %w", err)
			}
			if n.NoteId <= 0 || seen[n.NoteId] {
				return fmt.Errorf("invalid or duplicate note id %d", n.NoteId)
			}
			seen[n.NoteId], owners[file] = true, home
			n.HomeDeck, n.Path = home, file
			for i := range n.Cards {
				if n.Cards[i].GetDeck() == "" {
					n.Cards[i].Deck = proto.String(home)
				}
			}
			d.Notes = append(d.Notes, n)
		}
		return nil
	}
	for _, deck := range decks {
		path, err := referencePath(deck.Path, deck.NotesPath)
		if err != nil {
			return d, fmt.Errorf("resolve deck notes: %w", err)
		}
		if err := loadNotes(path, deck.Title); err != nil {
			return d, fmt.Errorf("load deck %q notes: %w", deck.Title, err)
		}
	}
	if d.Config.GetNotesPath() != "" {
		path, err := referencePath(config, d.Config.GetNotesPath())
		if err != nil {
			return d, fmt.Errorf("resolve cardless notes: %w", err)
		}
		if err := loadNotes(path, ""); err != nil {
			return d, fmt.Errorf("load cardless notes: %w", err)
		}
	}
	sort.Slice(d.Notes, func(i, j int) bool { return d.Notes[i].NoteId < d.Notes[j].NoteId })
	return d, nil
}

// PersistDeckIDs writes allocated deck identities back to their marker files.
func PersistDeckIDs(decks []model.Deck) error {
	for _, deck := range decks {
		current := model.Deck{Deck: &collectionpb.Deck{}}
		if err := ReadTOML(deck.Path, current.Deck); err != nil {
			return fmt.Errorf("read generated deck identity: %w", err)
		}
		if current.Id != 0 {
			continue
		}
		current.Id = deck.Id
		if err := WriteTOML(deck.Path, current.Deck); err != nil {
			return fmt.Errorf("record generated deck id %d: %w", deck.Id, err)
		}
	}
	return nil
}

// NotePath returns the export filename for a note within its deck directory.
func NotePath(dir, deck string, id int64) string {
	if deck == "" {
		return filepath.Join(dir, "notes", strconv.FormatInt(id, 10)+noteSuffix)
	}
	return filepath.Join(DeckPath(dir, deck), "notes", strconv.FormatInt(id, 10)+noteSuffix)
}

// tomlScalar retains scalar values and converts integer or nested message fields for TOML.
func tomlScalar(value any, scalar protoreflect.FieldDescriptor) (any, error) {
	switch scalar.Kind() {
	case protoreflect.MessageKind:
		if err := tomlIntegers(value.(map[string]any), scalar.Message()); err != nil {
			return nil, fmt.Errorf("convert nested field %s: %w", scalar.FullName(), err)
		}
	case protoreflect.Int32Kind, protoreflect.Int64Kind, protoreflect.Sint32Kind, protoreflect.Sint64Kind,
		protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind, protoreflect.Uint32Kind, protoreflect.Uint64Kind,
		protoreflect.Fixed32Kind, protoreflect.Fixed64Kind:
		integer, err := strconv.ParseInt(fmt.Sprint(value), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("convert integer field %s: %w", scalar.FullName(), err)
		}
		return integer, nil
	}
	return value, nil
}

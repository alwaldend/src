package ankiformat

import (
	"encoding/base64"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/model"
	"github.com/ankitects/anki/proto/anki/decks"
	"google.golang.org/protobuf/proto"
)

// MarshalConfig serializes an upstream protobuf configuration for database storage.
func MarshalConfig(message proto.Message) ([]byte, error) {
	data, err := (proto.MarshalOptions{Deterministic: true}).Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("encode Anki configuration: %w", err)
	}
	return data, nil
}

// DecodeDeckKind decodes a base64 upstream deck-kind message.
func DecodeDeckKind(encoded string) (*decks.Deck_KindContainer, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("decode deck kind: %w", err)
	}
	kind := &decks.Deck_KindContainer{}
	if err := proto.Unmarshal(data, kind); err != nil {
		return nil, fmt.Errorf("read Anki deck kind: %w", err)
	}
	if kind.Kind == nil {
		return nil, fmt.Errorf("deck needs a normal/filtered kind container")
	}
	return kind, nil
}

// PreserveDeckKind carries unedited upstream deck-kind data into the desired declaration.
func PreserveDeckKind(current model.Deck, desired *model.Deck) error {
	old, err := DecodeDeckKind(current.Kind)
	if err != nil {
		return fmt.Errorf("read existing deck kind: %w", err)
	}
	kind, err := DecodeDeckKind(desired.Kind)
	if err != nil {
		return fmt.Errorf("read desired deck kind: %w", err)
	}
	kind.ProtoReflect().SetUnknown(old.ProtoReflect().GetUnknown())
	if proto.Equal(old, kind) {
		// Keep untouched blobs verbatim, including explicit defaults and ordering.
		desired.Kind = current.Kind
		return nil
	}
	data, err := MarshalConfig(kind)
	if err != nil {
		return fmt.Errorf("preserve deck container settings: %w", err)
	}
	desired.Kind = base64.StdEncoding.EncodeToString(data)
	return nil
}

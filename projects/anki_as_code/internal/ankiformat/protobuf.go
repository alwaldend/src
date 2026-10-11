package ankiformat

import (
	"encoding/base64"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/model"
	"github.com/ankitects/anki/proto/anki/decks"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
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
	if err := PreserveUnknown(old, kind); err != nil {
		return fmt.Errorf("preserve deck unknown configuration: %w", err)
	}
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

// PreserveUnknown retains opaque protobuf data recursively in corresponding messages.
// Removing a nested message containing unknown data fails instead of losing it.
// Repeated messages correspond by position; maps correspond by key.
func PreserveUnknown(current, desired proto.Message) error {
	if current == nil || desired == nil || !current.ProtoReflect().IsValid() || !desired.ProtoReflect().IsValid() {
		return fmt.Errorf("configuration messages must be present")
	}
	if current.ProtoReflect().Descriptor().FullName() != desired.ProtoReflect().Descriptor().FullName() {
		return fmt.Errorf("configuration message types must match")
	}
	if err := preserveUnknown(current.ProtoReflect(), desired.ProtoReflect()); err != nil {
		return fmt.Errorf("preserve unknown configuration: %w", err)
	}
	return nil
}

// preserveUnknown copies opaque data and checks nested message removal.
func preserveUnknown(current, desired protoreflect.Message) error {
	if desired.IsValid() {
		desired.SetUnknown(append([]byte(nil), current.GetUnknown()...))
	} else if len(current.GetUnknown()) != 0 {
		return fmt.Errorf("cannot remove %s containing unknown configuration", current.Descriptor().FullName())
	}
	var failure error
	current.Range(func(field protoreflect.FieldDescriptor, value protoreflect.Value) bool {
		if field.Message() == nil {
			return true
		}
		var target protoreflect.Value
		if desired.IsValid() {
			target = desired.Get(field)
		}
		switch {
		case field.IsMap():
			if field.MapValue().Message() == nil {
				return true
			}
			value.Map().Range(func(key protoreflect.MapKey, entry protoreflect.Value) bool {
				var next protoreflect.Message
				if target.IsValid() && target.Map().Has(key) {
					next = target.Map().Get(key).Message()
				}
				failure = preserveUnknown(entry.Message(), validMessage(next, entry.Message()))
				return failure == nil
			})
		case field.IsList():
			for i := 0; i < value.List().Len(); i++ {
				var next protoreflect.Message
				if target.IsValid() && i < target.List().Len() {
					next = target.List().Get(i).Message()
				}
				previous := value.List().Get(i).Message()
				if failure = preserveUnknown(previous, validMessage(next, previous)); failure != nil {
					break
				}
			}
		default:
			var next protoreflect.Message
			if target.IsValid() {
				next = target.Message()
			}
			failure = preserveUnknown(value.Message(), validMessage(next, value.Message()))
		}
		return failure == nil
	})
	if failure != nil {
		return fmt.Errorf("preserve %s: %w", current.Descriptor().FullName(), failure)
	}
	return nil
}

// validMessage supplies an invalid message of the same type for an absent target.
func validMessage(message, current protoreflect.Message) protoreflect.Message {
	if message == nil {
		return current.Type().Zero()
	}
	return message
}

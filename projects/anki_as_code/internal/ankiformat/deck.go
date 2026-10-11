package ankiformat

import (
	"encoding/base64"
	"fmt"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"github.com/ankitects/anki/proto/anki/decks"
)

// DecodeDeck returns editable configuration from an encoded upstream deck kind.
// Identity and persistence metadata remain the caller's responsibility.
func DecodeDeck(encoded string) (*collectionpb.Deck, error) {
	kind, err := DecodeDeckKind(encoded)
	if err != nil {
		return nil, fmt.Errorf("read deck settings: %w", err)
	}
	if filtered := kind.GetFiltered(); filtered != nil {
		return &collectionpb.Deck{Filtered: filtered}, nil
	}
	config := kind.GetNormal()
	if config.ConfigId < 0 {
		return nil, fmt.Errorf("deck configuration ID must be nonnegative")
	}
	return &collectionpb.Deck{Normal: config}, nil
}

// EncodeDeck returns the upstream deck-kind payload for an editable declaration.
// It does not change the declaration or persistence metadata.
func EncodeDeck(editable *collectionpb.Deck) (string, error) {
	if editable.Normal != nil && editable.Filtered != nil {
		return "", fmt.Errorf("deck cannot be both normal and filtered")
	}
	kind := &decks.Deck_KindContainer{}
	if editable.Normal != nil {
		if editable.Normal.ConfigId < 0 {
			return "", fmt.Errorf("deck configuration ID must be nonnegative")
		}
		kind.Kind = &decks.Deck_KindContainer_Normal{Normal: editable.Normal}
	} else if editable.Filtered != nil {
		kind.Kind = &decks.Deck_KindContainer_Filtered{Filtered: editable.Filtered}
	} else {
		return "", fmt.Errorf("deck needs normal or filtered configuration")
	}
	data, err := MarshalConfig(kind)
	if err != nil {
		return "", fmt.Errorf("encode deck kind: %w", err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

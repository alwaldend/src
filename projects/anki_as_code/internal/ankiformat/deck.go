package ankiformat

import (
	"encoding/base64"
	"fmt"
	"math"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/model"
	"github.com/ankitects/anki/proto/anki/decks"
	"google.golang.org/protobuf/proto"
)

// DecodeDeck projects upstream deck configuration into its editable declaration.
func DecodeDeck(d *model.Deck) error {
	kind, err := DecodeDeckKind(d.Kind)
	if err != nil {
		return fmt.Errorf("read deck settings: %w", err)
	}
	if filtered := kind.GetFiltered(); filtered != nil {
		data, err := MarshalConfig(filtered)
		if err != nil {
			return fmt.Errorf("encode filtered settings: %w", err)
		}
		d.FilteredBase64 = proto.String(base64.StdEncoding.EncodeToString(data))
		d.Normal = nil
		return nil
	}
	normal := kind.GetNormal()
	if normal.ConfigId < 0 {
		return fmt.Errorf("deck configuration ID must be nonnegative")
	}
	d.Normal = &collectionpb.DeckConfig{
		ConfigId:            uint64(normal.ConfigId),
		ExtendNew:           uint64(normal.ExtendNew),
		ExtendReview:        uint64(normal.ExtendReview),
		Description:         normal.Description,
		MarkdownDescription: normal.MarkdownDescription,
	}
	d.FilteredBase64 = nil
	// The extra payload includes unmanaged members known to the pinned schema,
	// as well as unknown fields retained by the protobuf runtime.
	normal.ConfigId = 0
	normal.ExtendNew = 0
	normal.ExtendReview = 0
	normal.Description = ""
	normal.MarkdownDescription = false
	extra, err := MarshalConfig(normal)
	if err != nil {
		return fmt.Errorf("encode extra deck settings: %w", err)
	}
	d.Normal.ExtraFieldsBase64 = nil
	if len(extra) > 0 {
		d.Normal.ExtraFieldsBase64 = proto.String(base64.StdEncoding.EncodeToString(extra))
	}
	return nil
}

// EncodeDeck updates upstream deck configuration from its editable declaration.
func EncodeDeck(d *model.Deck) error {
	if d.Normal != nil && d.GetFilteredBase64() != "" {
		return fmt.Errorf("deck cannot be both normal and filtered")
	}
	kind := &decks.Deck_KindContainer{}
	if d.Normal != nil {
		n := d.Normal
		if n.ConfigId > math.MaxInt64 || n.ExtendNew > math.MaxUint32 || n.ExtendReview > math.MaxUint32 {
			return fmt.Errorf("deck settings exceed Anki's schema limits")
		}
		extra, err := base64.StdEncoding.DecodeString(n.GetExtraFieldsBase64())
		if err != nil {
			return fmt.Errorf("decode extra deck settings: %w", err)
		}
		normal := &decks.Deck_Normal{}
		if err := proto.Unmarshal(extra, normal); err != nil {
			return fmt.Errorf("read extra deck settings: %w", err)
		}
		if normal.ConfigId != 0 || normal.ExtendNew != 0 || normal.ExtendReview != 0 || normal.Description != "" || normal.MarkdownDescription {
			return fmt.Errorf("extra settings contain editable deck members")
		}
		normal.ConfigId = int64(n.ConfigId)
		normal.ExtendNew = uint32(n.ExtendNew)
		normal.ExtendReview = uint32(n.ExtendReview)
		normal.Description = n.Description
		normal.MarkdownDescription = n.MarkdownDescription
		kind.Kind = &decks.Deck_KindContainer_Normal{Normal: normal}
	} else {
		data, err := base64.StdEncoding.DecodeString(d.GetFilteredBase64())
		if err != nil {
			return fmt.Errorf("decode filtered settings: %w", err)
		}
		filtered := &decks.Deck_Filtered{}
		if err := proto.Unmarshal(data, filtered); err != nil {
			return fmt.Errorf("read filtered settings: %w", err)
		}
		kind.Kind = &decks.Deck_KindContainer_Filtered{Filtered: filtered}
	}
	data, err := MarshalConfig(kind)
	if err != nil {
		return fmt.Errorf("encode deck kind: %w", err)
	}
	d.Kind = base64.StdEncoding.EncodeToString(data)
	return nil
}

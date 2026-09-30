// Package draftjs models the subset of the DraftJS document format that the X
// Articles draft endpoint accepts as `content_state`.
//
// Entity and inline ranges are character offsets into the enclosing block's
// final text, counted in UTF-16 code units because the document is interpreted
// as a JavaScript string.
package draftjs

import "strconv"

// Block types the Articles API accepts.
const (
	BlockUnstyled          = "unstyled"
	BlockHeaderOne         = "header-one"
	BlockHeaderTwo         = "header-two"
	BlockHeaderThree       = "header-three"
	BlockUnorderedListItem = "unordered-list-item"
	BlockOrderedListItem   = "ordered-list-item"
	BlockBlockquote        = "blockquote"
	BlockAtomic            = "atomic"
)

// Inline styles the Articles API accepts.
const (
	StyleBold          = "bold"
	StyleItalic        = "italic"
	StyleStrikethrough = "strikethrough"
)

// Entity types the Articles API accepts.
const (
	EntityLink     = "link"
	EntityImage    = "image"
	EntityMarkdown = "markdown"
	EntityDivider  = "divider"
)

// Entity mutability values.
const (
	MutabilityImmutable = "immutable"
	MutabilityMutable   = "mutable"
	MutabilitySegmented = "segmented"
)

// MediaCategoryImage is the media category the upload endpoint returns for an
// uploaded image.
const MediaCategoryImage = "tweet_image"

// InlineStyleRange selects a styled span of a block's text.
type InlineStyleRange struct {
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Style  string `json:"style"`
}

// EntityRange selects the span of a block's text an entity applies to.
type EntityRange struct {
	Offset int `json:"offset"`
	Length int `json:"length"`
	Key    int `json:"key"`
}

// EntityValue is the payload an entity range selects.
type EntityValue struct {
	Type       string         `json:"type"`
	Mutability string         `json:"mutability"`
	Data       map[string]any `json:"data"`
}

// Entity is a document-level record referenced by an entity range. The Articles
// endpoint models the entity set as a mapping: each entry carries a string key
// that an entity range names, and the entity itself under `value`. This is not
// the DraftJS `entityMap` shape, whose keys are numeric and whose payload is
// inline.
type Entity struct {
	Key   string      `json:"key"`
	Value EntityValue `json:"value"`
}

// Block is one paragraph-like unit of the document. The Articles endpoint names
// the range fields with underscores and sets `additionalProperties: false`, so
// a block carries no `depth`: list nesting is not part of the accepted schema.
type Block struct {
	Key               string             `json:"key"`
	Type              string             `json:"type"`
	Text              string             `json:"text"`
	InlineStyleRanges []InlineStyleRange `json:"inline_style_ranges"`
	EntityRanges      []EntityRange      `json:"entity_ranges"`
	Data              map[string]any     `json:"data"`
}

// Document is the `content_state` payload.
type Document struct {
	Blocks   []Block  `json:"blocks"`
	Entities []Entity `json:"entities"`
}

// New returns an empty document with non-nil slices, so the emitted JSON uses
// arrays rather than nulls.
func New() *Document {
	return &Document{
		Blocks:   []Block{},
		Entities: []Entity{},
	}
}

// AddEntity appends an entity and returns its numeric key, which is also the
// string key an entity range records.
func (d *Document) AddEntity(entityType, mutability string, data map[string]any) int {
	key := len(d.Entities)
	if data == nil {
		data = map[string]any{}
	}
	d.Entities = append(d.Entities, Entity{
		Key: strconv.Itoa(key),
		Value: EntityValue{
			Type:       entityType,
			Mutability: mutability,
			Data:       data,
		},
	})
	return key
}

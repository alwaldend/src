// Package textstore reads and writes editable Anki collection declarations.
//
// It discovers configured deck markers, note files, and note-type appearance
// files; encodes TOML through generated protobuf contracts; and stages exports
// and identity updates. Collection validation is shared with validation, while
// database access and archive resources remain in their own adapters.
package textstore

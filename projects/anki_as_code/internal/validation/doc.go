// Package validation checks desired collection resources against current state.
//
// It validates identities, deck destinations, note fields, cloze ordinals, and
// editable appearance structure without reading files or mutating SQLite.
// Persistence-specific guards for unmanaged resources belong to repository.
package validation

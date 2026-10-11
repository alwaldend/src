// Package validation checks desired collection resources against current state.
//
// It validates identities, deck destinations, note fields, explicit card ordinals, and
// editable note-type structure and configuration without reading files or mutating SQLite.
// Persistence-specific guards for unmanaged resources belong to repository.
package validation

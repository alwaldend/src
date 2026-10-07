// Package repository owns SQLite persistence for Anki collections.
//
// It reads collection snapshots and reconciles managed resources in database
// transactions. Database handles, preservation bytes, and scheduling cursors
// remain private; review history, presets, and cards in filtered decks remain
// outside the managed reconciliation surface.
package repository

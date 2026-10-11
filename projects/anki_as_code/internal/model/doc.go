// Package model defines shared runtime models for Anki collections.
//
// Notes and decks wrap generated declarations with runtime metadata.
// Note types use their generated declarations directly, including ordered fields.
// Collection snapshots carry state exchanged between reconciliation packages;
// validation, planning, and persistence remain outside this package.
package model

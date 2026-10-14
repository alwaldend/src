// Package model defines shared runtime models for Anki collections.
//
// Notes and decks wrap generated declarations with runtime metadata.
// Note types wrap generated declarations with their native identity and source path.
// Collection snapshots carry state exchanged between reconciliation packages;
// validation, planning, and persistence remain outside this package.
package model

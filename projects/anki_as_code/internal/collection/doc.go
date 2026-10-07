// Package collection orchestrates offline collection export and reconciliation.
//
// It coordinates archive resources, text storage, SQLite persistence, and
// change planning for export, plan, apply, build, and note identity generation.
// Resource lifetimes and convergence checks are owned by this package; storage
// formats and reconciliation SQL remain in their respective adapters.
package collection

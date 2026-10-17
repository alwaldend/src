// Package app coordinates complete configuration and collection state workflows.
//
// Typed ColStateDump, ColStateImport, ColStateUpdate and ColStateReplace methods
// load required configs, borrow prepared connections from collection.Store and request
// archive publication for successful native changes. The CLI owns flags, state
// serialization and output. App exposes no raw connection or collection callback.
// Backend managers own native lifecycle; Anki owns storage semantics.
// Config dump needs no native backend. Collection export and reconciliation
// remain command placeholders.
package app

// Package app wires state commands to a shared native collection workflow.
//
// WithCollection composes anki_grpc's process and collection backends, lends a
// prepared connection and releases both backends. Callers construct stores from
// that connection. Anki owns database access and package restore/export.
// Config dump needs no native backend. Collection export and reconciliation
// remain command placeholders.
package app

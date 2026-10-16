// Package app wires configuration and state operations consumed by the command tree.
//
// Pure configuration operations need no Anki access. Start prepares one managed
// worker and collection; WithState borrows their connection across operations;
// Stop closes them in order. Collection input belongs to Options, while state
// options select only a namespace. Reconciliation remains a separate change.
package app

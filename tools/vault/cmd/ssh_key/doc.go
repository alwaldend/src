// Package main creates or reuses a dedicated Ed25519 SSH key pair in Vault.
// Callers own serialization and deployment; this command never rotates keys or
// writes private material to files, arguments, Terraform state, or output.
package main

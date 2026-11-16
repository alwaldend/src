// Package main bootstraps repository-scoped Forgejo runner credentials into
// Vault during an authorized Terraform apply. It reuses a cached record,
// refuses unmanaged name collisions, and never prints secret values.
package main

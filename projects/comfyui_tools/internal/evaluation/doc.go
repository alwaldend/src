// Package evaluation runs bounded image-generation workloads and owns their
// ephemeral rental lifecycles. Deployment supplies profiles, credentials and
// limits; this package never reads personal rental state or writes Vault.
package evaluation

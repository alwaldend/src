// Package evaluation runs bounded image-generation workloads and owns their
// rental lifecycles, including explicit reuse of retained installations.
// Deployment supplies profiles, credentials and limits; this package never reads
// personal rental state or writes Vault. Retained rentals require explicit release.
package evaluation

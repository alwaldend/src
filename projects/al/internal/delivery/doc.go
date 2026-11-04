// Package delivery owns guarded repository preparation, publication, and review.
// Typed receipts bind exact candidates, transport selection, remote leases, and
// observed forge state. Isolated Git transport uses SSH by default or explicitly
// selected GitHub CLI credentials over HTTPS without importing Git configuration.
// Callers run their own build and test checks through Bazel before publication.
package delivery

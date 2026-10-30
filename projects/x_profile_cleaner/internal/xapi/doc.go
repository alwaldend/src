// Package xapi adapts the shared X SDK to lazy inventory and explicit mutations.
// It caches timeline pages, normalizes requested owned kinds, and paces endpoint
// quotas, signing admitted requests with the shared OAuth 1.0a signer.
// Collection and quota progress use the command-owned logger.
// It ignores foreign and article content, refuses redirects, and never
// automatically replays failed requests or accepts unconfirmed mutations.
package xapi

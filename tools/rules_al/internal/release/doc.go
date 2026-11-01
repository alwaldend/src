// Package release generates project bundles and coordinates immutable artifact publication.
// It owns manifest identity, staged verification, website activation, and deletion.
// The AL command adapter uses these operations without starting credential plugins.
package release

// Package release generates project bundles and coordinates immutable artifact publication.
// It owns manifest identity, staged verification, website activation, and deletion.
// Offline plans share publication selection; read-only inspection reports local
// payload verification and bounded destination observations without remote locks.
// The AL command adapter uses these operations without starting credential plugins.
package release

// Package web adapts the cleanup application to an authenticated loopback interface.
// It owns embedded assets, temporary sessions, serialized jobs, review selections,
// resumable execution progress, and bounded log delivery. The application retains
// planning and mutation semantics;
// only subsets of server-owned plans may cross the browser execution boundary.
package web

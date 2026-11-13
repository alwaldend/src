// Package cleanup supplies the cleanup application library for archive scripts.
// Each library instance adapts Lua cleanup declarations and callbacks to the Go
// planner and validates exact operation/report agreement and archive provenance.
// Cleanup owns its graph projection of canonical archivepb records. Generic
// evaluation options and Lua/protobuf conversion belong to script; source parsing
// belongs to archive. Credentials and mutation remain outside them.
package cleanup

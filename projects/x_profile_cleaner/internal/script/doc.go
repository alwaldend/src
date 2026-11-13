// Package script evaluates bounded Lua scripts over read-only archive datasets.
// Schema-derived plan methods collect native protobuf operations per evaluation;
// scripts return no values. The consuming application owns authorization and
// execution. The interpreter has no filesystem, process, package or network API.
// Evaluator options configure deadlines and optional per-call application libraries.
// The caller owns archive lifetime and operating-system memory limits.
package script

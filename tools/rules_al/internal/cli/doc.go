// Package cli assembles AL's command tree and adapts command input and output.
// Release, delivery, versioning, and Bazel packages retain their domain behavior.
// Configuration and plugin startup are selected by config and run commands;
// ordinary repository inspection does not start credential plugins.
package cli

// Package main checks the public AL rules from an independent generated Bazel
// workspace and retains a verified manifest and payload result. Run it outside
// build actions so its nested Bazel client does not contend with its own build.
package main

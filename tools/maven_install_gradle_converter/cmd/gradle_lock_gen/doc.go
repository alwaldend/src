// Package main converts a rules_jvm_external maven_install.json
// lock file into one Gradle build input file per invocation, so both build
// systems can share the same pinned dependency set.
package main

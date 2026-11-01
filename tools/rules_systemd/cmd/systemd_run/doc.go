// Package main adapts command-line or BUILD-generated configuration and executes an existing command
// through a transient user systemd scope. It leaves builds and application
// configuration to callers and replaces itself to preserve process behavior.
package main

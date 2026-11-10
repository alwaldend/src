// Package main runs the packaged ComfyUI CLI through a temporary authenticated
// HTTPS bridge. It owns the listener and CLI lifetime; the caller supplies
// credentials, trust and endpoint discovery without granting live provisioning.
package main

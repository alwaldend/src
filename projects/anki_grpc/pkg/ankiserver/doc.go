// Package ankiserver runs Anki's native CLI through one linked entry point.
// The serve subcommand owns the calling process's protocol streams or listener until it exits.
// Offline applications should host it in their managed worker process.
package ankiserver

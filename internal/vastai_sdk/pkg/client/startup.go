package client

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
)

// maximumStartupBytes is Vast's inline startup command limit, including encoding.
const maximumStartupBytes = 16384

// StartupCommand prepares bounded API input without including script contents in errors.
func StartupCommand(script string, compress bool) (string, error) {
	command := script
	if compress && script != "" {
		var compressed bytes.Buffer
		writer, err := gzip.NewWriterLevel(&compressed, gzip.BestCompression)
		if err != nil {
			return "", fmt.Errorf("initialize startup compression: %w", err)
		}
		if _, err := writer.Write([]byte(script)); err != nil {
			return "", fmt.Errorf("compress startup: %w", err)
		}
		if err := writer.Close(); err != nil {
			return "", fmt.Errorf("finish startup compression: %w", err)
		}
		encoded := base64.StdEncoding.EncodeToString(compressed.Bytes())
		wrapped := "# Compressed Vast.ai startup\nset -eu\numask 077\nstartup_file=$(mktemp)\ntrap 'rm -f \"$startup_file\"' EXIT\nprintf '%s' '" + encoded + "' | base64 -d | gzip -dc > \"$startup_file\"\n/bin/bash \"$startup_file\"\n"
		if len(wrapped) < len(script) {
			command = wrapped
		}
	}
	if len(command) > maximumStartupBytes {
		return "", fmt.Errorf("startup command exceeds %d bytes after the selected encoding; reduce the payload or enable compression for a compatible image", maximumStartupBytes)
	}
	return command, nil
}

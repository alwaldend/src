package web

import (
	"embed"
	"fmt"
)

// assets contains the UI, generated shared theme, and maintained example policy.
//
//go:embed index.html app.js app.css theme.css policy.lua
var assets embed.FS

// DefaultPolicy returns the embedded maintained example for an unconfigured editor.
func DefaultPolicy() (string, error) {
	content, err := assets.ReadFile("policy.lua")
	if err != nil {
		return "", fmt.Errorf("read embedded policy: %w", err)
	}
	return string(content), nil
}

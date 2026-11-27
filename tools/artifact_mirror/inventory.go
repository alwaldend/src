package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// inventory contains the pinned artifacts from one upstream package.
type inventory struct {
	// Charts lists the upstream Helm archives.
	Charts []chart `json:"charts"`
	// Images lists the upstream OCI images.
	Images []image `json:"images"`
}

// loadInventory reads the pin file supplied by the package.
func loadInventory(filename string) (pins inventory, result error) {
	file, err := os.Open(filename)
	if err != nil {
		return inventory{}, fmt.Errorf("read artifact pins %q: %w", filename, err)
	}
	defer closeResource(file, filename, &result)
	// A pin file contains release metadata, never the artifacts themselves.
	const limit = 4 << 20
	data, err := io.ReadAll(io.LimitReader(file, limit+1))
	if err != nil {
		return inventory{}, fmt.Errorf("read artifact pins %q: %w", filename, err)
	}
	if len(data) > limit {
		return inventory{}, fmt.Errorf("artifact pins %q exceed %d bytes", filename, limit)
	}
	if err := json.Unmarshal(data, &pins); err != nil {
		return inventory{}, fmt.Errorf("decode artifact pins %q: %w", filename, err)
	}
	if len(pins.Charts) == 0 || len(pins.Images) == 0 {
		return inventory{}, fmt.Errorf("artifact pins must include charts and images")
	}
	for _, artifact := range pins.Charts {
		name := artifact.filename()
		if !filepath.IsLocal(name) || filepath.Base(name) != name || strings.ContainsAny(name, `\?#`) {
			return inventory{}, fmt.Errorf("chart name and version must form a local archive filename")
		}
	}
	for _, artifact := range pins.Images {
		name := path.Base(artifact.Repository)
		if name == "." || name == ".." || name == "/" || artifact.Tag == "" || strings.ContainsAny(artifact.Tag, `/\:`) {
			return inventory{}, fmt.Errorf("image repository and tag must identify a local layout tag")
		}
	}
	return pins, nil
}

// chart identifies an upstream Helm archive.
type chart struct {
	// Name is the chart name published in its metadata.
	Name string `json:"name"`
	// Version is the Helm chart version.
	Version string `json:"version"`
	// Source is the chart archive or OCI repository URL.
	Source string `json:"source"`
	// Digest is the SHA-256 checksum of the chart archive.
	Digest string `json:"digest"`
}

// filename is the archive name used by Helm and Forgejo.
func (c chart) filename() string {
	return c.Name + "-" + c.Version + ".tgz"
}

// image identifies an OCI manifest or multi-platform index.
type image struct {
	// Repository is the upstream registry and repository name.
	Repository string `json:"repository"`
	// Tag is the upstream release tag to use in the local layout and in Forgejo.
	Tag string `json:"tag"`
	// Digest pins the upstream manifest or index.
	Digest string `json:"digest"`
}

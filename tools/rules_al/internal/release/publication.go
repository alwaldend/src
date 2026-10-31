package release

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"git.alwaldend.com/alwaldend/src/tools/rules_al/api/contracts"
	"git.alwaldend.com/alwaldend/src/tools/rules_al/internal/versioning"
)

// publicationManifest is the transport-independent identity stored at a destination.
type publicationManifest struct {
	// Schema identifies the canonical publication format.
	Schema string `json:"schema"`
	// Namespace distinguishes first-party projects from external artifacts.
	Namespace string `json:"namespace"`
	// Project names the independently published project.
	Project string `json:"project"`
	// Version is the shared source version or the external dependency version.
	Version string `json:"version"`
	// Commit binds first-party bytes to repository source.
	Commit string `json:"commit,omitempty"`
	// Channel selects the repository version channel.
	Channel string `json:"channel,omitempty"`
	// Files lists the entire selected payload set in filename order.
	Files []publicationFile `json:"files"`
}

// publicationFile records the expected bytes of one ordinary downloadable file.
type publicationFile struct {
	// Name is a safe, unique basename.
	Name string `json:"name"`
	// Size is the expected byte count.
	Size int64 `json:"size"`
	// SHA256 is the lowercase digest of the complete payload.
	SHA256 string `json:"sha256"`
}

// validateReleaseSource rejects unstamped or modified first-party publications.
func validateReleaseSource(release *contracts.Release) error {
	if release.Channel != "release" && release.Channel != "nightly" {
		return fmt.Errorf("first-party publication requires a nightly or release channel")
	}
	tag, err := versioning.ParseVersionTag("v" + release.Name)
	if err != nil {
		return fmt.Errorf("release version: %w", err)
	}
	if (tag.Nightly != "") != (release.Channel == "nightly") {
		return fmt.Errorf("release version and channel disagree")
	}
	if release.TreeState != "clean" {
		return fmt.Errorf("first-party publication requires a clean source tree")
	}
	decoded, err := hex.DecodeString(release.SourceCommit)
	if err != nil || (len(decoded) != 20 && len(decoded) != 32) || strings.ToLower(release.SourceCommit) != release.SourceCommit {
		return fmt.Errorf("first-party publication requires a full source commit")
	}
	return nil
}

// makePublication canonicalizes the selected identity without deployment metadata.
func makePublication(release *contracts.Release, target *contracts.ReleaseDeploymentSsh, items []*contracts.ReleaseItem) (publicationManifest, string, error) {
	manifest := publicationManifest{Schema: "release/publication/v1", Namespace: target.Namespace, Project: target.Project, Version: release.Name, Commit: release.SourceCommit, Channel: release.Channel}
	if target.Namespace == "projects" {
		if err := validateReleaseSource(release); err != nil {
			return manifest, "", fmt.Errorf("validate source: %w", err)
		}
	}
	if len(items) == 0 {
		return manifest, "", fmt.Errorf("publication requires at least one file")
	}
	seen := map[string]bool{}
	for _, item := range items {
		if item.File == nil {
			return manifest, "", fmt.Errorf("release item has no file")
		}
		file := item.File
		if err := pathComponent(file.Name); err != nil {
			return manifest, "", fmt.Errorf("filename: %w", err)
		}
		if file.Name == "release.json" || seen[file.Name] || file.Size < 0 {
			return manifest, "", fmt.Errorf("invalid or duplicate payload %q", file.Name)
		}
		seen[file.Name] = true
		digest := ""
		for _, hash := range file.Hashes {
			if hash.Algo == "SHA-256" {
				if digest != "" {
					return manifest, "", fmt.Errorf("duplicate SHA-256 for %q", file.Name)
				}
				digest = hash.Content
			}
		}
		decoded, err := hex.DecodeString(digest)
		if err != nil || len(decoded) != sha256.Size || strings.ToLower(digest) != digest {
			return manifest, "", fmt.Errorf("missing or invalid SHA-256 for %q", file.Name)
		}
		manifest.Files = append(manifest.Files, publicationFile{Name: file.Name, Size: file.Size, SHA256: digest})
	}
	sort.Slice(manifest.Files, func(i, j int) bool { return manifest.Files[i].Name < manifest.Files[j].Name })
	data, err := json.Marshal(manifest)
	if err != nil {
		return manifest, "", fmt.Errorf("encode publication: %w", err)
	}
	return manifest, string(data) + "\n", nil
}

// verifyPublicationFile follows local runfiles symlinks to regular payloads and
// checks their bytes. Publication verifies the independent staging copy again.
func verifyPublicationFile(filename string, file publicationFile) error {
	info, err := os.Stat(filename)
	if err != nil {
		return fmt.Errorf("inspect payload %q: %w", filename, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("payload %q is not a regular file", filename)
	}
	input, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open payload %q: %w", filename, err)
	}
	hash := sha256.New()
	size, copyErr := io.Copy(hash, input)
	closeErr := input.Close()
	if copyErr != nil {
		return fmt.Errorf("hash payload %q: %w", filename, copyErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close payload %q: %w", filename, closeErr)
	}
	if size != file.Size || hex.EncodeToString(hash.Sum(nil)) != file.SHA256 {
		return fmt.Errorf("payload %q does not match its manifest", filename)
	}
	return nil
}

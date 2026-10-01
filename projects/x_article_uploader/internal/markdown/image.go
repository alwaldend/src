package markdown

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// acceptedImageTypes are the image media types the X media upload endpoints
// accept. image/svg+xml is deliberately absent: the endpoints reject it, so the
// converter reports it instead of emitting an entity publication cannot resolve.
var acceptedImageTypes = map[string]bool{
	"image/jpeg":  true,
	"image/gif":   true,
	"image/bmp":   true,
	"image/png":   true,
	"image/webp":  true,
	"image/pjpeg": true,
	"image/tiff":  true,
}

// ValidateImageType verifies accepted bytes against both the source extension
// and the artifact's declared media type before an upload can begin.
func ValidateImageType(reference string, content []byte, recordedType string) error {
	actualType := mediaTypeFor(reference, content)
	if actualType == "" {
		return fmt.Errorf("image %q has bytes that do not match its declared image extension", reference)
	}
	if !acceptedImageTypes[actualType] {
		return fmt.Errorf("image %q has media type %q, which the media upload endpoints do not accept", reference, actualType)
	}
	if actualType != recordedType {
		return fmt.Errorf("image %q media type %q does not match the artifact's %q", reference, actualType, recordedType)
	}
	return nil
}

// mediaTypeFor returns the media type an image is treated as, or "" when a
// recognized extension's declared type does not match the bytes. The declared
// extension is authoritative for the formats the endpoints accept, and the
// file's own signature confirms the declared type matches its bytes. Returning
// "" for a mismatch of two individually accepted formats is what lets the
// caller report it: falling back to the sniffed type would silently accept
// bytes renamed from one accepted format to another.
func mediaTypeFor(path string, content []byte) string {
	ext := strings.ToLower(filepath.Ext(path))
	byExtension := map[string]string{
		".jpg":   "image/jpeg",
		".jpeg":  "image/jpeg",
		".pjpeg": "image/pjpeg",
		".gif":   "image/gif",
		".bmp":   "image/bmp",
		".png":   "image/png",
		".webp":  "image/webp",
		".tif":   "image/tiff",
		".tiff":  "image/tiff",
		".svg":   "image/svg+xml",
	}
	if declared, ok := byExtension[ext]; ok {
		// The SVG type is already unaccepted, so it needs no byte check.
		if declared != "image/svg+xml" && !matchesImageFormat(declared, content) {
			return ""
		}
		return declared
	}
	return http.DetectContentType(content)
}

// matchesImageFormat reports whether content carries the leading signature of
// the given media type. Go's content sniffer recognizes the common raster
// formats except TIFF, so TIFF's signature is checked here too.
func matchesImageFormat(mediaType string, content []byte) bool {
	switch mediaType {
	case "image/png":
		return len(content) >= 8 && string(content[:8]) == "\x89PNG\r\n\x1a\n"
	case "image/jpeg", "image/pjpeg":
		return len(content) >= 3 && content[0] == 0xFF && content[1] == 0xD8 && content[2] == 0xFF
	case "image/gif":
		return len(content) >= 6 && (string(content[:6]) == "GIF87a" || string(content[:6]) == "GIF89a")
	case "image/bmp":
		return len(content) >= 2 && string(content[:2]) == "BM"
	case "image/webp":
		return len(content) >= 12 && string(content[:4]) == "RIFF" && string(content[8:12]) == "WEBP"
	case "image/tiff":
		return len(content) >= 4 &&
			(string(content[:4]) == "II*\x00" || string(content[:4]) == "MM\x00*")
	default:
		return true
	}
}

// readImage reads an image referenced by a post and returns its bytes, media
// type, and digest. It refuses a path that escapes the post's directory, and
// enforces containment with a descriptor-anchored handle rather than by a
// pathname check followed by a read: an apparently safe relative path could be
// a symlink whose target lies outside the post directory, and reading it would
// both inspect bytes outside the promised boundary and record a locator
// publication would later reject.
//
// When root is a handle the caller already opened and validated, the read goes
// through it; otherwise the post directory is opened here. Either way the
// resolution is descriptor-anchored, so the object whose containment is checked
// is the object read and a writer that swaps the target for a symlink after a
// check cannot redirect the read.
func readImage(root *os.Root, postDir, reference string) ([]byte, string, string, error) {
	cleaned := filepath.Clean(reference)
	if filepath.IsAbs(cleaned) || escapes(cleaned) {
		return nil, "", "", fmt.Errorf("image reference %q escapes the post directory", reference)
	}
	// Classify by the reference the author wrote, not the canonical target:
	// the declared extension is what the post claims the image is, and a
	// symlink to a differently named file must not let mismatched bytes pass the
	// declared-extension check. The bytes still come from the validated target.
	declaredPath := filepath.Join(postDir, cleaned)
	if root == nil {
		opened, err := os.OpenRoot(postDir)
		if err != nil {
			return nil, "", "", fmt.Errorf("open post directory %q: %w", postDir, err)
		}
		defer opened.Close()
		root = opened
	}
	content, err := root.ReadFile(cleaned)
	if err != nil {
		return nil, "", "", fmt.Errorf("read image %q: %w", reference, err)
	}
	sum := sha256.Sum256(content)
	return content, mediaTypeFor(declaredPath, content), hex.EncodeToString(sum[:]), nil
}

// escapes reports whether a cleaned relative path leaves its directory.
func escapes(cleaned string) bool {
	return cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator))
}

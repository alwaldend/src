package markdown

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// acceptedImageTypes follows X's documented simple-image upload formats.
// Animated images need a different upload flow and are refused separately.
var acceptedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/gif":  true,
	"image/png":  true,
	"image/webp": true,
}

const (
	// X documents a 5 MB maximum for tweet_image uploads.
	maxImageBytes = 5_000_000
	// This local decoder guard bounds allocation for highly compressed images;
	// it is not an X API limit. RGBA64 uses at most eight bytes per pixel.
	maxImagePixels = 32_000_000
)

// ValidateImageType checks the source extension, recorded type, byte limit,
// static-image contract, and content before any upload. JPEG, PNG, and GIF are
// decoded; WebP receives bounded container/header validation without a decoder.
func ValidateImageType(reference string, content []byte, recordedType string) error {
	if len(content) > maxImageBytes {
		return fmt.Errorf("image %q is %d bytes; static image uploads allow at most %d bytes", reference, len(content), maxImageBytes)
	}
	actualType := mediaTypeFor(reference, content)
	if actualType == "" {
		return fmt.Errorf("image %q has bytes that do not match its declared image extension", reference)
	}
	if !acceptedImageTypes[actualType] {
		return fmt.Errorf("image %q has unsupported media type %q; simple uploads support static JPEG, PNG, GIF, and WebP", reference, actualType)
	}
	if actualType != recordedType {
		return fmt.Errorf("image %q media type %q does not match the artifact's %q", reference, actualType, recordedType)
	}
	if err := validateStaticImage(content, actualType); err != nil {
		return fmt.Errorf("validate image %q: %w", reference, err)
	}
	return nil
}

func validateStaticImage(content []byte, mediaType string) error {
	if mediaType == "image/webp" {
		if err := validateStaticWebP(content); err != nil {
			return fmt.Errorf("validate WebP container: %w", err)
		}
		return nil
	}
	config, _, err := image.DecodeConfig(bytes.NewReader(content))
	if err != nil {
		return fmt.Errorf("decode image dimensions: %w", err)
	}
	if err := validateImageDimensions(config.Width, config.Height); err != nil {
		return fmt.Errorf("validate image dimensions: %w", err)
	}
	switch mediaType {
	case "image/png":
		if err := validateStaticPNG(content); err != nil {
			return fmt.Errorf("validate PNG chunks: %w", err)
		}
		if _, err := png.Decode(bytes.NewReader(content)); err != nil {
			return fmt.Errorf("decode PNG: %w", err)
		}
	case "image/jpeg":
		if _, err := jpeg.Decode(bytes.NewReader(content)); err != nil {
			return fmt.Errorf("decode JPEG: %w", err)
		}
	case "image/gif":
		if err := validateStaticGIF(content); err != nil {
			return fmt.Errorf("validate GIF frames: %w", err)
		}
		if _, err := gif.DecodeAll(bytes.NewReader(content)); err != nil {
			return fmt.Errorf("decode GIF: %w", err)
		}
	}
	return nil
}

func validateImageDimensions(width, height int) error {
	if width <= 0 || height <= 0 || uint64(width)*uint64(height) > maxImagePixels {
		return fmt.Errorf("image dimensions %dx%d exceed the local decoding limit of %d pixels", width, height, maxImagePixels)
	}
	return nil
}

// validateStaticPNG inspects actual chunk boundaries so text metadata containing
// an animation marker cannot be mistaken for APNG. The decoder checks CRCs.
func validateStaticPNG(content []byte) error {
	for offset := 8; offset < len(content); {
		if len(content)-offset < 12 {
			return fmt.Errorf("truncated chunk header")
		}
		size := uint64(binary.BigEndian.Uint32(content[offset : offset+4]))
		if size+12 > uint64(len(content)-offset) {
			return fmt.Errorf("truncated chunk payload")
		}
		switch string(content[offset+4 : offset+8]) {
		case "acTL", "fcTL", "fdAT":
			return fmt.Errorf("animated PNG is unsupported; use a static image")
		}
		offset += int(size) + 12
	}
	return nil
}

// validateStaticGIF counts frames without decompressing them, avoiding unbounded
// allocations for a small file that encodes many large frames. DecodeAll then
// verifies the complete single-frame image rather than just its signature.
func validateStaticGIF(content []byte) error {
	if len(content) < 13 {
		return fmt.Errorf("truncated logical screen descriptor")
	}
	offset := 13
	if content[10]&0x80 != 0 {
		offset += 3 << ((content[10] & 7) + 1)
	}
	frames := 0
	for offset < len(content) {
		block := content[offset]
		offset++
		switch block {
		case 0x3b:
			if frames != 1 || offset != len(content) {
				return fmt.Errorf("invalid single-frame GIF trailer")
			}
			return nil
		case 0x21:
			offset++ // Extension label precedes its data sub-blocks.
		case 0x2c:
			frames++
			if frames > 1 {
				return fmt.Errorf("animated GIF is unsupported; use a static image")
			}
			if len(content)-offset < 9 {
				return fmt.Errorf("truncated image descriptor")
			}
			flags := content[offset+8]
			offset += 9
			if flags&0x80 != 0 {
				offset += 3 << ((flags & 7) + 1)
			}
			offset++ // LZW minimum code size precedes the image sub-blocks.
		default:
			return fmt.Errorf("invalid GIF block")
		}
		for {
			if offset >= len(content) {
				return fmt.Errorf("truncated GIF data sub-block")
			}
			size := int(content[offset])
			offset++
			if size > len(content)-offset {
				return fmt.Errorf("truncated GIF data sub-block payload")
			}
			offset += size
			if size == 0 {
				break
			}
		}
	}
	return fmt.Errorf("missing GIF trailer")
}

// validateStaticWebP follows the RIFF container and frame-header contract:
// https://developers.google.com/speed/webp/docs/riff_container.
// It checks framing, dimensions, and animation, not compressed pixel decoding.
func validateStaticWebP(content []byte) error {
	if len(content) < 12 || string(content[:4]) != "RIFF" || string(content[8:12]) != "WEBP" {
		return fmt.Errorf("invalid RIFF signature")
	}
	if uint64(binary.LittleEndian.Uint32(content[4:8]))+8 != uint64(len(content)) {
		return fmt.Errorf("RIFF length does not match the file")
	}
	width, height, canvasWidth, canvasHeight, frames := 0, 0, 0, 0, 0
	for offset := 12; offset < len(content); {
		if len(content)-offset < 8 {
			return fmt.Errorf("truncated RIFF chunk header")
		}
		kind := string(content[offset : offset+4])
		size := uint64(binary.LittleEndian.Uint32(content[offset+4 : offset+8]))
		padded := size + size%2
		if padded > uint64(len(content)-offset-8) {
			return fmt.Errorf("truncated %s chunk", kind)
		}
		chunk := content[offset+8 : offset+8+int(size)]
		if size%2 != 0 && content[offset+8+int(size)] != 0 {
			return fmt.Errorf("nonzero %s chunk padding", kind)
		}
		switch kind {
		case "ANIM", "ANMF":
			return fmt.Errorf("animated WebP is unsupported; use a static image")
		case "VP8X":
			if offset != 12 || len(chunk) != 10 {
				return fmt.Errorf("invalid extended header")
			}
			if chunk[0]&2 != 0 {
				return fmt.Errorf("animated WebP is unsupported; use a static image")
			}
			canvasWidth = 1 + int(chunk[4]) + int(chunk[5])<<8 + int(chunk[6])<<16
			canvasHeight = 1 + int(chunk[7]) + int(chunk[8])<<8 + int(chunk[9])<<16
		case "VP8 ":
			if len(chunk) <= 10 || chunk[0]&1 != 0 || !bytes.Equal(chunk[3:6], []byte{0x9d, 0x01, 0x2a}) {
				return fmt.Errorf("invalid VP8 key-frame header")
			}
			width = int(binary.LittleEndian.Uint16(chunk[6:8]) & 0x3fff)
			height = int(binary.LittleEndian.Uint16(chunk[8:10]) & 0x3fff)
			frames++
		case "VP8L":
			if len(chunk) <= 5 || chunk[0] != 0x2f || chunk[4]&0xe0 != 0 {
				return fmt.Errorf("invalid VP8L header")
			}
			bits := binary.LittleEndian.Uint32(chunk[1:5])
			width, height = 1+int(bits&0x3fff), 1+int((bits>>14)&0x3fff)
			frames++
		}
		offset += 8 + int(padded)
	}
	if frames != 1 {
		return fmt.Errorf("expected one image bitstream, found %d", frames)
	}
	if canvasWidth != 0 && (canvasWidth != width || canvasHeight != height) {
		return fmt.Errorf("canvas dimensions do not match the image bitstream")
	}
	if err := validateImageDimensions(width, height); err != nil {
		return fmt.Errorf("validate WebP dimensions: %w", err)
	}
	return nil
}

// mediaTypeFor returns the media type an image is treated as, or "" when a
// recognized extension's declared type does not match the bytes. The declared
// extension is authoritative for recognized raster formats, and the
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

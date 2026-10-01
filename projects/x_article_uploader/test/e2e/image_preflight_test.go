package e2e

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"path/filepath"
	"strings"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/draft"
	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/markdown"
)

// Image preflight failure modes, specified before implementation:
//   - unsupported formats or files larger than the simple-image limit reach X;
//   - a recognizable signature hides corrupt or truncated image contents;
//   - animated GIF, PNG, or WebP is silently uploaded as a static tweet_image;
//   - banner validation is stronger than retained-body validation;
//   - a forged artifact bypasses conversion checks and spends network calls;
//   - valid static formats or the inclusive size boundary are rejected.
// The generated artifact and local HTTP exchanges remain repeatable evidence.

type preflightImage struct {
	name, filename, mediaType string
	content                   []byte
}

func decodeImageFixture(t *testing.T, encoded string) []byte {
	t.Helper()
	content, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("decode image fixture: %v", err)
	}
	return content
}

func staticImageFixtures(t *testing.T) []preflightImage {
	t.Helper()
	pixel := image.NewPaletted(image.Rect(0, 0, 2, 2), color.Palette{color.Black, color.White})
	pixel.SetColorIndex(1, 1, 1)
	var jpegBytes, gifBytes bytes.Buffer
	if err := jpeg.Encode(&jpegBytes, pixel, nil); err != nil {
		t.Fatalf("encode JPEG fixture: %v", err)
	}
	if err := gif.Encode(&gifBytes, pixel, nil); err != nil {
		t.Fatalf("encode GIF fixture: %v", err)
	}
	return []preflightImage{
		{"png", "image.png", "image/png", tinyPNG},
		{"jpeg", "image.jpeg", "image/jpeg", jpegBytes.Bytes()},
		{"gif", "image.gif", "image/gif", gifBytes.Bytes()},
		// The lossless WebP and animation fixtures were encoded from two 2x2
		// colored frames with Pillow 12.3.0; they are real decodable images.
		{"webp", "image.webp", "image/webp", decodeImageFixture(t, "UklGRhwAAABXRUJQVlA4TA8AAAAvAUAAAAcQ/Y/+ByKi/wEA")},
	}
}

// paddedPNG adds a legal ancillary text chunk before IEND, keeping the image
// decodable while exercising the inclusive byte-size boundary.
func paddedPNG(t *testing.T, size int) []byte {
	t.Helper()
	payload := make([]byte, size-len(tinyPNG)-12)
	copy(payload, []byte("comment\x00"))
	for index := len("comment\x00"); index < len(payload); index++ {
		payload[index] = 'x'
	}
	var chunk bytes.Buffer
	if err := binary.Write(&chunk, binary.BigEndian, uint32(len(payload))); err != nil {
		t.Fatalf("encode PNG chunk length: %v", err)
	}
	chunk.WriteString("tEXt")
	chunk.Write(payload)
	checksum := crc32.ChecksumIEEE(chunk.Bytes()[4:])
	if err := binary.Write(&chunk, binary.BigEndian, checksum); err != nil {
		t.Fatalf("encode PNG checksum: %v", err)
	}
	result := append([]byte(nil), tinyPNG[:len(tinyPNG)-12]...)
	result = append(result, chunk.Bytes()...)
	return append(result, tinyPNG[len(tinyPNG)-12:]...)
}

func rejectedImageFixtures(t *testing.T) []preflightImage {
	t.Helper()
	static := staticImageFixtures(t)
	pixel := image.NewPaletted(image.Rect(0, 0, 2, 2), color.Palette{color.Black, color.White})
	var animation bytes.Buffer
	if err := gif.EncodeAll(&animation, &gif.GIF{Image: []*image.Paletted{pixel, pixel}, Delay: []int{10, 10}}); err != nil {
		t.Fatalf("encode animated GIF fixture: %v", err)
	}
	corruptPNG := append([]byte(nil), tinyPNG...)
	corruptPNG[45] ^= 0xff
	largeCanvasPNG := append([]byte(nil), tinyPNG...)
	binary.BigEndian.PutUint32(largeCanvasPNG[16:20], 100_000)
	binary.BigEndian.PutUint32(largeCanvasPNG[20:24], 100_000)
	binary.BigEndian.PutUint32(largeCanvasPNG[29:33], crc32.ChecksumIEEE(largeCanvasPNG[12:29]))
	truncatedWebP := append([]byte(nil), static[3].content[:len(static[3].content)-4]...)
	return []preflightImage{
		{"oversized", "image.png", "image/png", paddedPNG(t, 5_000_001)},
		{"bmp", "image.bmp", "image/bmp", []byte("BMunsupported bitmap")},
		{"tiff", "image.tiff", "image/tiff", []byte("II*\x00unsupported TIFF")},
		{"pjpeg", "image.pjpeg", "image/pjpeg", static[1].content},
		{"png_signature_only", "image.png", "image/png", tinyPNG[:8]},
		{"png_corrupt_pixels", "image.png", "image/png", corruptPNG},
		{"excessive_decoded_pixels", "image.png", "image/png", largeCanvasPNG},
		{"jpeg_truncated", "image.jpeg", "image/jpeg", static[1].content[:100]},
		{"gif_truncated", "image.gif", "image/gif", static[2].content[:13]},
		{"webp_signature_only", "image.webp", "image/webp", []byte("RIFF\x04\x00\x00\x00WEBP")},
		{"webp_truncated", "image.webp", "image/webp", truncatedWebP},
		{"animated_gif", "image.gif", "image/gif", animation.Bytes()},
		{"animated_png", "image.png", "image/png", decodeImageFixture(t, "iVBORw0KGgoAAAANSUhEUgAAAAIAAAACCAYAAABytg0kAAAACGFjVEwAAAACAAAAAPONk3AAAAAaZmNUTAAAAAAAAAACAAAAAgAAAAAAAAAAAAEACgAA6FTcAAAAABRJREFUeJxj/M/A8J+BgYGBiQEKAB8XAgJPlM6+AAAAGmZjVEwAAAABAAAAAgAAAAIAAAAAAAAAAAABAAoAAHMnNtQAAAAYZmRBVAAAAAJ4nGNkYPj/n4GBgYGJAQoAHRkCAunm7jEAAAAASUVORK5CYII=")},
		{"animated_webp", "image.webp", "image/webp", decodeImageFixture(t, "UklGRoQAAABXRUJQVlA4WAoAAAACAAAAAQAAAQAAQU5JTQYAAAAAAAAAAABBTk1GKAAAAAAAAAAAAAEAAAEAAGQAAAJWUDhMDwAAAC8BQAAABxD9j/4HIqL/AQBBTk1GKAAAAAAAAAAAAAEAAAEAAGQAAABWUDhMDwAAAC8BQAAABxDR//4HIqL/AQA=")},
	}
}

func TestImagePreflightConversionRefusesInvalidSources(t *testing.T) {
	for _, fixture := range rejectedImageFixtures(t) {
		for _, placement := range []string{"banner", "body"} {
			t.Run(fixture.name+"/"+placement, func(t *testing.T) {
				metadata, body := "", "Article text.\n\n![image]("+fixture.filename+")"
				if placement == "banner" {
					metadata, body = "images: ["+fixture.filename+"]\n", "Article text."
				}
				_, postDir, source := bannerPost(t, metadata, body)
				writeFile(t, filepath.Join(postDir, fixture.filename), fixture.content)
				article, err := markdown.New(postDir, "post").Convert(source)
				if err == nil {
					t.Fatalf("conversion accepted %s %s", fixture.name, placement)
				}
				if placement == "banner" {
					if !strings.Contains(err.Error(), "front matter images[0] as banner") || !strings.Contains(err.Error(), fixture.filename) {
						t.Fatalf("invalid banner lacks source context: %v", err)
					}
					return
				}
				var diagnosticError *markdown.DiagnosticError
				if !errors.As(err, &diagnosticError) || article == nil || !hasDiagnostic(diagnosticCodes(article), "image-media-type-rejected") {
					t.Fatalf("invalid image lacks a conversion diagnostic: %v", err)
				}
			})
		}
	}
}

func TestImagePreflightForgedArtifactMakesNoRequests(t *testing.T) {
	for _, fixture := range rejectedImageFixtures(t) {
		for _, placement := range []string{"banner", "body"} {
			t.Run(fixture.name+"/"+placement, func(t *testing.T) {
				root, postDir, source := bannerPost(t, "images: [banner.png]\n", "![body](body.png)")
				path, artifactJSON := bannerArtifact(t, root, postDir, source)
				artifact, err := draft.ReadArtifact(path)
				if err != nil {
					t.Fatalf("read baseline artifact: %v", err)
				}
				writeFile(t, filepath.Join(postDir, fixture.filename), fixture.content)
				locator := markdown.ImageSource{
					Path: fixture.filename, PostPackage: "post", MediaType: fixture.mediaType,
					Digest: fmt.Sprintf("%x", sha256.Sum256(fixture.content)),
				}
				if placement == "banner" {
					artifact.Banner = &locator
				} else {
					artifact.Locators[0].ImageSource = locator
				}
				service, client := newBannerService(t)
				_, err = draft.New(client, root, "").CreateDraft(artifact)
				artifactJSON["rejected_source"] = locator
				artifactJSON["placement"] = placement
				if err != nil {
					artifactJSON["error"] = err.Error()
				}
				exchanges := service.recorded()
				writeBannerEvidence(t, artifactJSON, exchanges)
				if err == nil || len(exchanges) != 0 {
					t.Fatalf("invalid %s reached publication: error=%v, requests=%d", placement, err, len(exchanges))
				}
			})
		}
	}
}

func TestImagePreflightStaticFormatsHTTPPipeline(t *testing.T) {
	fixtures := append(staticImageFixtures(t), preflightImage{"size_boundary", "image.png", "image/png", paddedPNG(t, 5_000_000)})
	for _, fixture := range fixtures {
		t.Run(fixture.name, func(t *testing.T) {
			root, postDir, source := bannerPost(t, "images: ["+fixture.filename+"]\n", "Article text.\n\n![image]("+fixture.filename+")")
			writeFile(t, filepath.Join(postDir, fixture.filename), fixture.content)
			path, artifactJSON := bannerArtifact(t, root, postDir, source)
			artifact, err := draft.ReadArtifact(path)
			if err != nil {
				t.Fatalf("read image artifact: %v", err)
			}
			service, client := newBannerService(t)
			if _, err := draft.New(client, root, "").CreateDraft(artifact); err != nil {
				t.Fatalf("create draft with static image: %v", err)
			}
			exchanges := service.recorded()
			writeBannerEvidence(t, artifactJSON, exchanges)
			if len(exchanges) != 2 || exchanges[0].Path != "/media/upload" || exchanges[1].Path != "/articles/draft" {
				t.Fatalf("static banner/body should share one upload: %#v", exchanges)
			}
			if exchanges[0].Digest != fmt.Sprintf("%x", sha256.Sum256(fixture.content)) {
				t.Error("uploaded bytes differ from the validated image")
			}
		})
	}
}

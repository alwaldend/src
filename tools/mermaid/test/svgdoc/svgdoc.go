// Package svgdoc inspects rendered Mermaid artifacts: the paint layers of an
// SVG, the canvas an SVG draws, and the dimensions of the raster projection a
// service-upload consumer renders.
//
// Mermaid writes one group per layer under the root group, and SVG paints in
// document order, so an element's offset in the document is its paint order.
// The checks that need these scans live in more than one workspace, so the
// scans live here rather than in each.
package svgdoc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ClusterTitleMarker identifies a container title group. Mermaid may write the
// class with a trailing space that a renderer trims, so the marker stops before
// the closing quote and matches either spelling.
const ClusterTitleMarker = `class="cluster-label`

// NodesMarker identifies the node layer, which is the last layer Mermaid writes
// before any lifted container title.
const NodesMarker = `<g class="nodes">`

// GroupEnd returns the offset just past the group that opens at the first
// occurrence of marker, and whether that group closes.
//
// Counting group tags is exact for a rendered document: the renderer writes no
// markup as an attribute value, so no group tag can hide inside one. A
// self-closed group, such as the empty edge-label layer, adds no depth.
func GroupEnd(svg, marker string) (int, bool) {
	start := strings.Index(svg, marker)
	if start < 0 {
		return 0, false
	}
	depth := 0
	for index := start; index < len(svg); {
		switch {
		case strings.HasPrefix(svg[index:], "<g "), strings.HasPrefix(svg[index:], "<g>"):
			tag := strings.Index(svg[index:], ">")
			if tag < 0 {
				return 0, false
			}
			if !strings.HasSuffix(svg[index:index+tag], "/") {
				depth++
			}
			index += tag + 1
		case strings.HasPrefix(svg[index:], "</g>"):
			depth--
			index += len("</g>")
			if depth == 0 {
				return index, true
			}
		default:
			index++
		}
	}
	return 0, false
}

// ClusterTitlesPaintLast reports whether every container title is written after
// the node layer. Mermaid emits a title inside the clusters layer ahead of the
// edges, so an edge crossing a container border would paint over the title; the
// renderer lifts each title to the end of the root group to fix that.
func ClusterTitlesPaintLast(svg string) bool {
	if !strings.Contains(svg, ClusterTitleMarker) {
		return false
	}
	nodesEnd, ok := GroupEnd(svg, NodesMarker)
	if !ok {
		return false
	}
	return strings.Index(svg, ClusterTitleMarker) > nodesEnd
}

// CanvasSize returns the size a rendered SVG draws, read from its `viewBox`.
//
// A Mermaid render writes the drawing's own size into `viewBox`; a raster
// projection of that document must reproduce it, so the value is read here
// rather than from the document's `width`/`height`, which a consumer may have
// rewritten for its own layout.
func CanvasSize(svg string) (float64, float64, error) {
	index := strings.Index(svg, `viewBox="`)
	if index < 0 {
		return 0, 0, errors.New("rendered SVG carries no viewBox")
	}
	rest := svg[index+len(`viewBox="`):]
	end := strings.Index(rest, `"`)
	if end < 0 {
		return 0, 0, errors.New("rendered SVG has an unterminated viewBox")
	}
	fields := strings.Fields(rest[:end])
	if len(fields) != 4 {
		return 0, 0, fmt.Errorf("viewBox %q does not carry four numbers", rest[:end])
	}
	width, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("viewBox %q is not numeric", rest[:end])
	}
	height, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("viewBox %q is not numeric", rest[:end])
	}
	return width, height, nil
}

// WebPDimensions returns a WebP image's canvas size from its RIFF header.
//
// The pinned browser emits a `VP8X` chunk when it writes an ICC profile and a
// `VP8`/`VP8L` frame otherwise, so all three layouts are handled rather than
// depending on which encoder path Chrome chose or adding an image dependency.
func WebPDimensions(data []byte) (int, int, error) {
	if len(data) < 16 || string(data[0:4]) != "RIFF" || string(data[8:12]) != "WEBP" {
		return 0, 0, fmt.Errorf("not a WebP image")
	}
	chunk := data[20:]
	switch string(data[12:16]) {
	case "VP8X":
		if len(chunk) < 10 {
			return 0, 0, errors.New("VP8X chunk is truncated")
		}
		width := int(chunk[4]) | int(chunk[5])<<8 | int(chunk[6])<<16
		height := int(chunk[7]) | int(chunk[8])<<8 | int(chunk[9])<<16
		return width + 1, height + 1, nil
	case "VP8 ":
		if len(chunk) < 10 || !strings.HasPrefix(string(chunk[3:6]), "\x9d\x01\x2a") {
			return 0, 0, errors.New("VP8 frame start code is missing")
		}
		width := int(chunk[6]) | int(chunk[7])<<8
		height := int(chunk[8]) | int(chunk[9])<<8
		return width & 0x3fff, height & 0x3fff, nil
	case "VP8L":
		if len(chunk) < 5 || chunk[0] != 0x2f {
			return 0, 0, errors.New("VP8L frame header is missing")
		}
		bits := uint32(chunk[1]) | uint32(chunk[2])<<8 | uint32(chunk[3])<<16 | uint32(chunk[4])<<24
		return int(bits&0x3fff) + 1, int((bits>>14)&0x3fff) + 1, nil
	default:
		return 0, 0, fmt.Errorf("unrecognized WebP chunk %q", string(data[12:16]))
	}
}

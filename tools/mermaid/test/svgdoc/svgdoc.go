// Package svgdoc locates paint layers in a rendered Mermaid SVG.
//
// Mermaid writes one group per layer under the root group, and SVG paints in
// document order, so an element's offset in the document is its paint order.
// Two test packages need the same scan, so it lives here rather than in each.
package svgdoc

import "strings"

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

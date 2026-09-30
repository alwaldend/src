package markdown

import (
	"bytes"

	gast "github.com/yuin/goldmark/ast"
)

// lineStart returns the offset of the first byte of the line containing offset.
func lineStart(body []byte, offset int) int {
	if offset > len(body) {
		offset = len(body)
	}
	i := bytes.LastIndexByte(body[:offset], '\n')
	return i + 1
}

// lineEnd returns the offset just past the line containing offset, excluding
// the trailing newline.
func lineEnd(body []byte, offset int) int {
	if offset > len(body) {
		offset = len(body)
	}
	i := bytes.IndexByte(body[offset:], '\n')
	if i < 0 {
		return len(body)
	}
	return offset + i
}

// fencedSource recovers a fenced code block's original source, including its
// opening fence, info string, and closing fence. Goldmark records the block's
// inner lines, not its fences, so the bounds come from the block's position and
// its last inner line rather than from the inner lines' own start.
func fencedSource(body []byte, node gast.Node) []byte {
	return body[fencedStart(body, node):fencedEnd(body, node)]
}

// indentedSource recovers an indented code block's original source, which is
// its content lines. An indented block has no fence, so no delimiter line is
// included.
//
// The block's last line segment stops just past its own newline, so the slice
// ends at that offset directly rather than scanning for the next line end.
// Scanning would swallow the first line of a following paragraph into the
// markdown entity, and the block walker would then emit that paragraph again.
func indentedSource(body []byte, firstInner, lastInner int) []byte {
	start := lineStart(body, firstInner)
	if lastInner > len(body) {
		lastInner = len(body)
	}
	if start > lastInner {
		start = lastInner
	}
	return body[start:lastInner]
}

// tableSource recovers a table's original Markdown source from the line span
// between its first header cell and its last row cell, which also includes the
// delimiter line between them.
func tableSource(body []byte, firstCell, lastCell int) []byte {
	start := lineStart(body, firstCell)
	end := lineEnd(body, lastCell)
	return body[start:end]
}

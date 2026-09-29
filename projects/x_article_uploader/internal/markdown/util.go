package markdown

import (
	"strconv"
	"strings"

	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"

	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/draftjs"
)

func itoa(v int) string { return strconv.Itoa(v) }

func trimSpace(s string) string { return strings.TrimSpace(s) }

func blockKey(index int) string { return "block-" + strconv.Itoa(index) }

func nonNilStyles(v []draftjs.InlineStyleRange) []draftjs.InlineStyleRange {
	if v == nil {
		return []draftjs.InlineStyleRange{}
	}
	return v
}

func nonNilRanges(v []draftjs.EntityRange) []draftjs.EntityRange {
	if v == nil {
		return []draftjs.EntityRange{}
	}
	return v
}

// blockOffset returns the byte offset of a block node's first line, falling
// back to the start of the body when the node has no lines.
func blockOffset(body []byte, node gast.Node) int {
	if lines := node.Lines(); lines != nil && lines.Len() > 0 {
		return lines.At(0).Start
	}
	return 0
}

// withinBlockquote reports whether a node sits inside a block quote, however
// deeply nested. A quoted heading can be nested beneath a list item, so the
// current output context alone cannot tell whether the quote still encloses the
// node.
func withinBlockquote(node gast.Node) bool {
	for parent := node.Parent(); parent != nil; parent = parent.Parent() {
		if _, ok := parent.(*gast.Blockquote); ok {
			return true
		}
	}
	return false
}

// inlineOffset returns the byte offset of the first text within an inline
// subtree, used to locate a construct that spans inline nodes. A titled link or
// image with an empty label or alt text has no descendant Text node, so the
// construct's own Pos is used instead: Goldmark sets it to the delimiter that
// opens the construct, which is where the diagnostic belongs.
func inlineOffset(node gast.Node) int {
	offset := -1
	_ = gast.Walk(node, func(n gast.Node, entering bool) (gast.WalkStatus, error) {
		if !entering {
			return gast.WalkContinue, nil
		}
		if textNode, ok := n.(*gast.Text); ok && offset < 0 {
			offset = textNode.Segment.Start
		}
		return gast.WalkContinue, nil
	})
	if offset < 0 {
		if pos := node.Pos(); pos >= 0 {
			return pos
		}
		return 0
	}
	return offset
}

// imageOffset locates an image within its container. An image with alt text is
// located at that text; an image with empty alt text has no child to locate, so
// the image's own Pos is used, which is the delimiter that opens it. A
// block-level image has no text line of its own, so the container's first line
// is the final fallback.
func imageOffset(body []byte, image *gast.Image, container gast.Node) int {
	if child := image.FirstChild(); child != nil {
		if textNode, ok := child.(*gast.Text); ok {
			return textNode.Segment.Start
		}
	}
	if pos := image.Pos(); pos >= 0 {
		return pos
	}
	if lines := container.Lines(); lines != nil && lines.Len() > 0 {
		return lines.At(0).Start
	}
	return 0
}

// innerStart returns the first line's start of an indented code block, or the
// line start of the block when it has no lines.
func innerStart(node gast.Node) int {
	if lines := node.Lines(); lines != nil && lines.Len() > 0 {
		return lines.At(0).Start
	}
	return 0
}

// innerEnd returns the last line's end of an indented code block.
func innerEnd(node gast.Node) int {
	if lines := node.Lines(); lines != nil && lines.Len() > 0 {
		return lines.At(lines.Len() - 1).Stop
	}
	return 0
}

// fencedStart returns the offset of a fenced code block's opening fence line.
// Pos() is set to the opening fence, so the preserved source keeps the fence
// and info string; an inner line's start would drop both.
func fencedStart(body []byte, node gast.Node) int {
	return lineStart(body, node.Pos())
}

// fencedEnd returns the offset just past a fenced code block's closing fence
// line, or the offset just past its last line when no closing fence exists.
//
// A block with inner lines closes on the line after its last one. An empty
// block has no inner lines, so its closing fence is the line after the opening
// fence; a fence opened and closed at the end of input has no following line,
// and the source then ends at the opening fence.
func fencedEnd(body []byte, node gast.Node) int {
	lines := node.Lines()
	if lines != nil && lines.Len() > 0 {
		return lineEnd(body, lines.At(lines.Len()-1).Stop)
	}
	openingEnd := lineEnd(body, lineStart(body, node.Pos()))
	if openingEnd >= len(body) {
		return openingEnd
	}
	return lineEnd(body, openingEnd+1)
}

// firstCellOffset returns the first table cell's text start.
func firstCellOffset(body []byte, table gast.Node) int {
	first := -1
	_ = gast.Walk(table, func(n gast.Node, entering bool) (gast.WalkStatus, error) {
		if !entering {
			return gast.WalkContinue, nil
		}
		if n.Type() != gast.TypeBlock {
			return gast.WalkContinue, nil
		}
		if cell, ok := n.(interface{ Lines() *text.Segments }); ok && first < 0 {
			if lines := cell.Lines(); lines != nil && lines.Len() > 0 {
				first = lines.At(0).Start
			}
		}
		return gast.WalkContinue, nil
	})
	if first < 0 {
		return 0
	}
	return first
}

// lastCellOffset returns the last table cell's text end.
func lastCellOffset(body []byte, table gast.Node) int {
	last := -1
	_ = gast.Walk(table, func(n gast.Node, entering bool) (gast.WalkStatus, error) {
		if !entering {
			return gast.WalkContinue, nil
		}
		if n.Type() != gast.TypeBlock {
			return gast.WalkContinue, nil
		}
		if cell, ok := n.(interface{ Lines() *text.Segments }); ok {
			if lines := cell.Lines(); lines != nil && lines.Len() > 0 {
				last = lines.At(lines.Len() - 1).Stop
			}
		}
		return gast.WalkContinue, nil
	})
	if last < 0 {
		return len(body)
	}
	return last
}

// tableSourceFromNode recovers a table's source Markdown from its cell span.
func tableSourceFromNode(body []byte, table gast.Node) []byte {
	return tableSource(body, firstCellOffset(body, table), lastCellOffset(body, table))
}

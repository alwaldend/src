package markdown

import (
	"strings"

	gast "github.com/yuin/goldmark/ast"
	extensionast "github.com/yuin/goldmark/extension/ast"
	gmutil "github.com/yuin/goldmark/util"

	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/draftjs"
)

// inlineBuilder accumulates the text, inline styles, and entity ranges for one
// block. Offsets are measured in UTF-16 code units against the text built so
// far, matching how the document is interpreted as a JavaScript string.
type inlineBuilder struct {
	text   []rune
	styles []draftjs.InlineStyleRange
	ranges []draftjs.EntityRange
	doc    *draftjs.Document
	// open holds the inline styles and entities whose text is still being
	// written. A range is recorded when the context closes; the portion before
	// an inline image split is recorded when the split flushes, so a styled or
	// linked span that an image interrupts keeps its range on both sides.
	open []openContext
}

// openContext is an inline style or entity whose text is still being written.
type openContext struct {
	start      int
	style      string
	entityType string
	mutability string
	data       map[string]any
	entity     bool
}

// pushStyle starts a style context at the current text length.
func (b *inlineBuilder) pushStyle(style string) {
	b.open = append(b.open, openContext{start: b.length(), style: style})
}

// pushEntity starts an entity context at the current text length.
func (b *inlineBuilder) pushEntity(entityType, mutability string, data map[string]any) {
	b.open = append(b.open, openContext{
		start: b.length(), entityType: entityType, mutability: mutability, data: data, entity: true,
	})
}

// popContext closes the innermost context and records the range it covers.
func (b *inlineBuilder) popContext() {
	if len(b.open) == 0 {
		return
	}
	ctx := b.open[len(b.open)-1]
	b.open = b.open[:len(b.open)-1]
	b.emitContext(ctx, b.length())
}

// emitContext records the range a context covers over text ending at end. An
// empty span records nothing, so a context that selects no text adds no range.
func (b *inlineBuilder) emitContext(ctx openContext, end int) {
	if end <= ctx.start {
		return
	}
	if ctx.entity {
		key := b.doc.AddEntity(ctx.entityType, ctx.mutability, ctx.data)
		b.ranges = append(b.ranges, draftjs.EntityRange{Offset: ctx.start, Length: end - ctx.start, Key: key})
		return
	}
	b.styles = append(b.styles, draftjs.InlineStyleRange{Offset: ctx.start, Length: end - ctx.start, Style: ctx.style})
}

func (b *inlineBuilder) write(s string) {
	b.text = append(b.text, []rune(s)...)
}

// length returns the current text length in UTF-16 code units.
func (b *inlineBuilder) length() int {
	return utf16Len(b.text)
}

// utf16Len returns the number of UTF-16 code units needed for runes.
func utf16Len(runes []rune) int {
	n := 0
	for _, r := range runes {
		if r > 0xFFFF {
			n += 2
			continue
		}
		n++
	}
	return n
}

// writeInline appends node's inline content to the builder.
func (c *conversion) writeInline(node gast.Node, b *inlineBuilder, blockType string) {
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		switch n := child.(type) {
		case *gast.Text:
			b.write(resolveText(n.Segment.Value(c.body), n.IsRaw()))
			// Both break flags end the text node on a new line: a soft break
			// is the source's own newline, and a hard break is a newline the
			// author asked for. Neither may be dropped, which would join two
			// lines into one word.
			if n.SoftLineBreak() || n.HardLineBreak() {
				b.write("\n")
			}
		case *gast.String:
			b.write(resolveText(n.Value, n.IsRaw()))
		case *gast.Emphasis:
			style := draftjs.StyleItalic
			if n.Level >= 2 {
				style = draftjs.StyleBold
			}
			b.pushStyle(style)
			c.writeInline(n, b, blockType)
			b.popContext()
		case *extensionast.Strikethrough:
			b.pushStyle(draftjs.StyleStrikethrough)
			c.writeInline(n, b, blockType)
			b.popContext()
		case *gast.CodeSpan:
			start := b.length()
			b.write(string(n.Text(c.body)))
			c.report(Continuing, "inline-code-style-lost",
				"an inline-code span keeps its text but loses its monospace styling", inlineOffset(n))
			_ = start
		case *gast.Link:
			c.reportTitleLoss(n.Title, "link-title-lost", "a link keeps its text and destination but loses its title", n)
			b.pushEntity(draftjs.EntityLink, draftjs.MutabilityMutable, map[string]any{"url": resolveDestination(n.Destination)})
			c.writeInline(n, b, blockType)
			b.popContext()
		case *gast.AutoLink:
			// An email autolink's URL is the bare address, so it needs the
			// mailto scheme Goldmark's own renderer adds; without it the link
			// would be a relative URL rather than an email address.
			url := string(n.URL(c.body))
			if n.AutoLinkType == gast.AutoLinkEmail && !strings.HasPrefix(strings.ToLower(url), "mailto:") {
				url = "mailto:" + url
			}
			b.pushEntity(draftjs.EntityLink, draftjs.MutabilityMutable, map[string]any{"url": url})
			b.write(resolveText(n.Label(c.body), false))
			b.popContext()
		case *extensionast.FootnoteLink:
			b.write("[" + itoa(n.Index) + "]")
		case *extensionast.FootnoteBacklink:
			// Definitions are emitted separately, without their backlinks.
		case *extensionast.TaskCheckBox:
			// Goldmark represents a task-list checkbox as a childless node, so
			// the default recursion below would drop its state and leave only
			// the item text. No DraftJS construct carries it, so preserve it as
			// bracketed text the way a footnote reference is preserved.
			if n.IsChecked {
				b.write("[x] ")
			} else {
				b.write("[ ] ")
			}
		case *gast.Image:
			// An image carries no text range, so a link around it cannot become
			// a DraftJS link over the image: the image is always emitted
			// unlinked. When the link has no text of its own, the destination
			// would otherwise vanish entirely; when it does have sibling text,
			// that text keeps the destination but the image still loses its own
			// link. Both are reported so the author is told which part of the
			// hyperlink will not appear, rather than the outcome depending on
			// whether unrelated sibling text happens to exist.
			if link := enclosingLink(node); link != nil {
				message := "an image wrapped in a link keeps the image but loses its destination "
				if !linkTextAbsent(link, c.body) {
					message = "an image inside a link keeps the image, and the link's text keeps the destination, " +
						"but the image itself is emitted without its link to "
				}
				c.report(Continuing, "linked-image-destination-lost",
					message+resolveDestination(link.Destination), imageOffset(c.body, n, node))
			}
			// An image in the middle of a block splits it: the text so far is
			// flushed, the image becomes its own atomic block, and collection
			// continues in a fresh builder. Any context open across the split
			// records its range over the text before the image here, and
			// restarts at the new text so the text after the image also carries
			// its style or entity.
			end := b.length()
			for _, ctx := range b.open {
				b.emitContext(ctx, end)
			}
			c.flushBlock(b, blockType)
			c.emitImage(n, node)
			b.text = b.text[:0]
			b.styles = nil
			b.ranges = nil
			for i := range b.open {
				b.open[i].start = 0
			}
		case *gast.RawHTML:
			c.report(Failing, "unrepresentable-html",
				"an inline HTML fragment cannot be represented in content_state", inlineOffset(n))
		default:
			c.writeInline(n, b, blockType)
		}
	}
}

// reportTitleLoss reports a link or image title the target entity cannot carry.
// The API's entity `data` object accepts only the documented properties, so a
// title has nowhere to go; losing it silently would break the no-silent-loss
// contract, so it is reported as a continuing diagnostic the way other lost
// formatting is.
func (c *conversion) reportTitleLoss(title []byte, code, message string, node gast.Node) {
	if len(strings.TrimSpace(string(title))) == 0 {
		return
	}
	c.report(Continuing, code, message+": "+string(title), inlineOffset(node))
}

// resolveText resolves the Markdown escaping and character references a text
// node's source carries, matching how Goldmark's own renderers present prose.
// Goldmark leaves the source encoding in the segment for the renderer to
// resolve, so a paragraph containing `\\*literal\\*` or `AT&amp;T` would
// otherwise show the author's backslashes or the raw entity. Raw text — a code
// span or a raw-HTML string — is emitted verbatim, because its source encoding
// is its content.
func resolveText(source []byte, raw bool) string {
	if raw {
		return string(source)
	}
	return resolveDestination(source)
}

// resolveInlineText resolves the inline content of a node into a caption or
// label, applying reference resolution to each text node and leaving a raw
// descendant (a code span or raw-HTML string) verbatim. Resolving the flattened
// Node.Text instead would run the whole caption through the non-raw path and
// strip a code span's literal backslashes, so the raw-code guarantee would not
// hold inside a caption.
func resolveInlineText(node gast.Node, body []byte) string {
	var b strings.Builder
	_ = gast.Walk(node, func(n gast.Node, entering bool) (gast.WalkStatus, error) {
		if !entering {
			return gast.WalkContinue, nil
		}
		switch v := n.(type) {
		case *gast.Text:
			b.WriteString(resolveText(v.Segment.Value(body), v.IsRaw()))
			if v.SoftLineBreak() || v.HardLineBreak() {
				b.WriteString("\n")
			}
		case *gast.String:
			b.WriteString(resolveText(v.Value, v.IsRaw()))
		}
		return gast.WalkContinue, nil
	})
	return b.String()
}

// resolveDestination resolves the Markdown escaping and character references a
// destination may carry, matching how Goldmark's own renderers normalize a
// destination before use. A destination such as `a\(1\)?x=1&amp;y=2` is stored
// as `a(1)?x=1&y=2`, so the recorded link or image target names the resource the
// author wrote rather than one containing literal backslashes or an unresolved
// entity.
func resolveDestination(destination []byte) string {
	resolved := gmutil.UnescapePunctuations(destination)
	resolved = gmutil.ResolveNumericReferences(resolved)
	resolved = gmutil.ResolveEntityNames(resolved)
	return string(resolved)
}

// enclosingLink returns the nearest ancestor link of node, or nil. An image
// wrapped in a link may sit inside inline wrappers such as emphasis, so the
// link is an ancestor rather than the image's immediate parent.
func enclosingLink(node gast.Node) *gast.Link {
	for ancestor := node; ancestor != nil; ancestor = ancestor.Parent() {
		if link, ok := ancestor.(*gast.Link); ok {
			return link
		}
	}
	return nil
}

// linkTextAbsent reports whether a link carries no non-whitespace text of its
// own. A link whose only content is an image has none: the image's alt text is
// a child of the image, not link text, so it does not carry the destination.
// The link's whole subtree is walked, so an image nested inside a style wrapper
// is still recognized as the only content.
func linkTextAbsent(link *gast.Link, body []byte) bool {
	absent := true
	_ = gast.Walk(link, func(n gast.Node, entering bool) (gast.WalkStatus, error) {
		if !entering {
			return gast.WalkContinue, nil
		}
		if _, ok := n.(*gast.Image); ok {
			// The alt text belongs to the image, not to the link.
			return gast.WalkSkipChildren, nil
		}
		// Every other leaf that produces text carries the link: plain text,
		// inline code, and a string leaf all keep a text span the destination
		// can attach to.
		switch leaf := n.(type) {
		case *gast.Text:
			if trimSpace(string(leaf.Segment.Value(body))) != "" {
				absent = false
				return gast.WalkStop, nil
			}
		case *gast.CodeSpan:
			if trimSpace(string(leaf.Text(body))) != "" {
				absent = false
				return gast.WalkStop, nil
			}
		case *gast.String:
			if trimSpace(string(leaf.Value)) != "" {
				absent = false
				return gast.WalkStop, nil
			}
		}
		return gast.WalkContinue, nil
	})
	return absent
}

// emitInlineBlock emits one block for a block node's inline content.
func (c *conversion) emitInlineBlock(node gast.Node, blockType string) {
	b := &inlineBuilder{doc: c.doc}
	c.writeInline(node, b, blockType)
	c.flushBlock(b, blockType)
}

// flushBlock appends the builder's accumulated block when it carries text.
func (c *conversion) flushBlock(b *inlineBuilder, blockType string) {
	if len(b.text) == 0 {
		return
	}
	c.doc.Blocks = append(c.doc.Blocks, draftjs.Block{
		Key:               blockKey(len(c.doc.Blocks)),
		Type:              blockType,
		Text:              string(b.text),
		InlineStyleRanges: nonNilStyles(b.styles),
		EntityRanges:      nonNilRanges(b.ranges),
		Data:              map[string]any{},
	})
}

// emitAtomic appends an atomic block backed by an entity and returns the
// entity's key.
func (c *conversion) emitAtomic(entityType, mutability string, data map[string]any) int {
	key := c.doc.AddEntity(entityType, mutability, data)
	c.doc.Blocks = append(c.doc.Blocks, draftjs.Block{
		Key:               blockKey(len(c.doc.Blocks)),
		Type:              draftjs.BlockAtomic,
		Text:              " ",
		InlineStyleRanges: []draftjs.InlineStyleRange{},
		EntityRanges:      []draftjs.EntityRange{{Offset: 0, Length: 1, Key: key}},
		Data:              map[string]any{},
	})
	return key
}

// emitDivider appends a thematic break as an atomic block with a divider
// entity.
func (c *conversion) emitDivider() {
	c.emitAtomic(draftjs.EntityDivider, draftjs.MutabilityImmutable, map[string]any{})
}

// emitMarkdownEntityAt appends source-preserving Markdown as an atomic block.
func (c *conversion) emitMarkdownEntityAt(source []byte) {
	payload := string(source)
	c.markdownWeightEstimate += markdownWeightedLengthEstimate(payload)
	c.emitAtomic(draftjs.EntityMarkdown, draftjs.MutabilityMutable, map[string]any{"markdown": payload})
}

// emitImage appends an atomic block with an unresolved image entity and records
// the locator publication resolves it from. An image whose media type the
// upload endpoints reject, or whose bytes cannot be read, is reported instead
// of being emitted.
func (c *conversion) emitImage(image *gast.Image, container gast.Node) {
	reference := resolveDestination(image.Destination)
	caption := resolveInlineText(image, c.body)
	c.reportTitleLoss(image.Title, "image-title-lost",
		"an image keeps its media but loses its title", image)
	content, mediaType, digest, err := readImage(c.converter.PackageRoot, c.converter.PostDir, reference)
	offset := imageOffset(c.body, image, container)
	if err != nil {
		c.report(Failing, "image-not-readable", err.Error(), offset)
		return
	}
	if err := ValidateImageType(reference, content, mediaType); err != nil {
		c.report(Failing, "image-media-type-rejected", err.Error(), offset)
		return
	}
	key := c.emitAtomic(draftjs.EntityImage, draftjs.MutabilityImmutable, map[string]any{"caption": caption})
	c.locators = append(c.locators, ImageLocator{
		EntityKey: key,
		ImageSource: ImageSource{
			Path:        reference,
			PostPackage: c.converter.PostPackage,
			Digest:      digest,
			MediaType:   mediaType,
		},
		Caption: caption,
	})
}

// soleImage returns the image a node contains when it is the node's only
// content, ignoring surrounding whitespace.
func soleImage(node gast.Node, body []byte) (*gast.Image, bool) {
	var found *gast.Image
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if textNode, ok := child.(*gast.Text); ok {
			if trimSpace(string(textNode.Segment.Value(body))) == "" {
				continue
			}
			return nil, false
		}
		if image, ok := child.(*gast.Image); ok {
			if found != nil {
				return nil, false
			}
			found = image
			continue
		}
		return nil, false
	}
	return found, found != nil
}

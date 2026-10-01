package markdown

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	extensionast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"

	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/draftjs"
)

// MarkdownPayloadBudget bounds the total `markdown` entity payload for one
// article. The API documents a 10,000 weighted-length limit; this fails below
// that so an over-budget article is refused locally rather than by the API.
const MarkdownPayloadBudget = 9500

// Article is the emitted draft artifact: the content_state document, the title
// draft creation needs, the locators publication resolves images from, and the
// diagnostics conversion produced.
type Article struct {
	Title       string            `json:"title"`
	Document    *draftjs.Document `json:"content_state"`
	Locators    []ImageLocator    `json:"image_locators"`
	Banner      *ImageSource      `json:"banner_locator,omitempty"`
	Diagnostics []Diagnostic      `json:"diagnostics"`
}

// ImageSource records where an image's bytes come from, resolved
// against the post package so the artifact stays usable after it is moved.
type ImageSource struct {
	Path        string `json:"path"`
	PostPackage string `json:"post_package"`
	Digest      string `json:"digest"`
	MediaType   string `json:"media_type"`
}

// ImageLocator associates an unresolved body image with its document entity.
// Embedding ImageSource preserves the existing locator's JSON field names.
type ImageLocator struct {
	EntityKey int `json:"entity_key"`
	ImageSource
	Caption string `json:"caption"`
}

// Converter converts parsed posts into draft artifacts. PostDir is the
// directory a post's relative image references resolve against, and
// PostPackage is the Bazel package recorded in each locator. PostSource is the
// path diagnostics name; it defaults to PostDir.
//
// PackageRoot, when set, is a descriptor-anchored handle on the post directory
// the caller validated. Image reads then resolve through it instead of
// re-opening PostDir, so a symlink retargeted after validation cannot redirect
// a read while the artifact still records the validated package.
type Converter struct {
	PostDir     string
	PostPackage string
	PostSource  string
	PackageRoot *os.Root
}

// New returns a converter for a post directory and its owning package.
func New(postDir, postPackage string) *Converter {
	return &Converter{PostDir: postDir, PostPackage: postPackage, PostSource: postDir}
}

// Convert parses a source post and returns its draft artifact. A failing
// diagnostic makes conversion fail with a DiagnosticError instead of emitting a
// document that would lose content.
func (c *Converter) Convert(raw []byte) (*Article, error) {
	post, err := ParsePost(c.PostDir, raw)
	if err != nil {
		return nil, fmt.Errorf("parse post: %w", err)
	}
	var banner *ImageSource
	if len(post.Images) > 0 {
		reference := post.Images[0]
		content, mediaType, digest, err := readImage(c.PackageRoot, c.PostDir, reference)
		if err != nil {
			return nil, fmt.Errorf("%s: resolve front matter images[0] as banner: %w", c.PostSource, err)
		}
		if err := ValidateImageType(reference, content, mediaType); err != nil {
			return nil, fmt.Errorf("%s: validate front matter images[0] as banner: %w", c.PostSource, err)
		}
		banner = &ImageSource{Path: reference, PostPackage: c.PostPackage, Digest: digest, MediaType: mediaType}
	}
	conv := &conversion{converter: c, post: post, body: post.Body, doc: draftjs.New()}

	md := goldmark.New(goldmark.WithExtensions(extension.GFM, extension.Footnote))
	doc := md.Parser().Parse(text.NewReader(post.Body))

	if err := conv.walkBlocks(doc, blockContext{}); err != nil {
		return nil, fmt.Errorf("convert post blocks: %w", err)
	}
	conv.appendFootnotes(doc)

	if conv.markdownLen > MarkdownPayloadBudget {
		conv.report(Failing, "markdown-payload-over-budget", fmt.Sprintf(
			"the article's markdown payload is %d bytes, over the %d budget",
			conv.markdownLen, MarkdownPayloadBudget), 0)
	}

	article := &Article{
		Title:       post.Title,
		Document:    conv.doc,
		Locators:    conv.locators,
		Banner:      banner,
		Diagnostics: conv.diagnostics,
	}
	if err := conv.failure(); err != nil {
		return article, fmt.Errorf("convert post: %w", err)
	}
	return article, nil
}

// JSON marshals the artifact with a stable field order.
func (a *Article) JSON() ([]byte, error) {
	return json.MarshalIndent(a, "", "  ")
}

// blockContext carries the block-mapping state for the blocks being emitted.
type blockContext struct {
	// blockType is the DraftJS type paragraphs in this context map to. A block
	// quote maps its paragraphs onto blockquote blocks; a list item supplies
	// its own type.
	blockType string
}

type conversion struct {
	converter   *Converter
	post        *Post
	body        []byte
	doc         *draftjs.Document
	locators    []ImageLocator
	diagnostics []Diagnostic
	markdownLen int
}

func (c *conversion) report(severity Severity, code, message string, offset int) {
	c.diagnostics = append(c.diagnostics, Diagnostic{
		Severity:   severity,
		Code:       code,
		Message:    message,
		Source:     c.converter.PostSource,
		Position:   positionAt(c.body, c.post.FileLine, offset),
		ByteOffset: offset,
	})
}

// failure returns a DiagnosticError when a failing diagnostic was reported.
func (c *conversion) failure() error {
	var failing []Diagnostic
	for _, d := range c.diagnostics {
		if d.Severity == Failing {
			failing = append(failing, d)
		}
	}
	if len(failing) == 0 {
		return nil
	}
	return &DiagnosticError{Source: c.post.Source, Diagnostics: failing}
}

// walkBlocks maps the block-level children of node.
func (c *conversion) walkBlocks(node gast.Node, ctx blockContext) error {
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if err := c.walkBlock(child, ctx); err != nil {
			return err
		}
	}
	return nil
}

func (c *conversion) walkBlock(node gast.Node, ctx blockContext) error {
	switch n := node.(type) {
	case *gast.Heading:
		// A heading inside a block quote has no DraftJS block type: the
		// blockquote type already carries the quote, and X exposes no
		// quoted-heading form. Emitting header-one here would silently move the
		// line out of the quote, so keep it quoted and report the lost heading
		// level. Quote ancestry is checked rather than the current output block
		// type, because a quote may enclose the heading through a list, where
		// the current type is the item's; the block is emitted as a blockquote
		// so the quote is preserved rather than replaced by the item type.
		if withinBlockquote(node) {
			c.report(Continuing, "blockquote-heading-unrepresentable",
				"a heading inside a block quote keeps its text as quoted content but loses its heading level", blockOffset(c.body, node))
			c.emitInlineBlock(node, draftjs.BlockBlockquote)
			break
		}
		c.emitInlineBlock(node, headingType(n.Level))
	case *gast.Paragraph:
		if image, ok := soleImage(n, c.body); ok {
			c.emitImage(image, n)
			return nil
		}
		c.emitInlineBlock(node, paragraphType(ctx))
	case *gast.TextBlock:
		if image, ok := soleImage(n, c.body); ok {
			c.emitImage(image, n)
			return nil
		}
		c.emitInlineBlock(node, paragraphType(ctx))
	case *gast.Blockquote:
		return c.walkBlocks(node, blockContext{blockType: draftjs.BlockBlockquote})
	case *gast.List:
		return c.walkList(n, ctx)
	case *gast.FencedCodeBlock:
		c.emitMarkdownEntityAt(fencedSource(c.body, n))
	case *gast.CodeBlock:
		c.emitMarkdownEntityAt(indentedSource(c.body, innerStart(n), innerEnd(n)))
	case *gast.ThematicBreak:
		c.emitDivider()
	case *extensionast.Table:
		c.emitMarkdownEntityAt(tableSourceFromNode(c.body, n))
	case *gast.HTMLBlock:
		c.report(Failing, "unrepresentable-html",
			"an HTML block cannot be represented in content_state", blockOffset(c.body, n))
	case *extensionast.FootnoteList:
		// Footnote definitions are appended as a trailing section.
		return nil
	default:
		if !node.HasChildren() {
			return nil
		}
		return c.walkBlocks(node, ctx)
	}
	return nil
}

// walkList maps a list's items, recursing into nested lists so their items are
// emitted as their own blocks rather than merged into the parent item.
func (c *conversion) walkList(list *gast.List, ctx blockContext) error {
	itemType := draftjs.BlockUnorderedListItem
	if list.IsOrdered() {
		itemType = draftjs.BlockOrderedListItem
	}
	for item := list.FirstChild(); item != nil; item = item.NextSibling() {
		listItem, ok := item.(*gast.ListItem)
		if !ok {
			continue
		}
		// One list item is one block. A loose item exposes each of its
		// paragraphs as a separate child, so accumulate them into one builder
		// rather than emitting one block per paragraph, which DraftJS would
		// read as several items and which would renumber an ordered list. A
		// blank line between the item's own paragraphs is preserved.
		b := &inlineBuilder{doc: c.doc}
		// A block construct inside the item — a nested list, a fenced or
		// indented code block, a table, a thematic break, or a single image —
		// is not inline content. Such constructs are collected in source order
		// and emitted after the item's text, so the item stays one block while
		// the constructs keep the order they appear in; emitting one as soon as
		// it is met would place it before a construct collected earlier, and
		// flushing the item text at each one would split the item.
		var deferred []gast.Node
		for child := listItem.FirstChild(); child != nil; child = child.NextSibling() {
			switch n := child.(type) {
			case *gast.List:
				deferred = append(deferred, n)
			case *gast.Paragraph:
				if _, ok := soleImage(n, c.body); ok {
					deferred = append(deferred, n)
					continue
				}
				if len(b.text) > 0 {
					b.write("\n")
				}
				c.writeInline(n, b, itemType)
			case *gast.TextBlock:
				if _, ok := soleImage(n, c.body); ok {
					deferred = append(deferred, n)
					continue
				}
				if len(b.text) > 0 {
					b.write("\n")
				}
				c.writeInline(n, b, itemType)
			default:
				deferred = append(deferred, child)
			}
		}
		c.flushBlock(b, itemType)
		for _, node := range deferred {
			// A nested list's items carry their own type rather than inheriting
			// the parent item's.
			if nested, ok := node.(*gast.List); ok {
				if err := c.walkList(nested, blockContext{blockType: ctx.blockType}); err != nil {
					return err
				}
				continue
			}
			if err := c.walkBlock(node, blockContext{blockType: itemType}); err != nil {
				return err
			}
		}
	}
	return nil
}

// paragraphType returns the block type a paragraph maps onto in this context.
func paragraphType(ctx blockContext) string {
	if ctx.blockType != "" {
		return ctx.blockType
	}
	return draftjs.BlockUnstyled
}

// headingType clamps a heading level onto the three levels X exposes.
func headingType(level int) string {
	switch {
	case level <= 2:
		return draftjs.BlockHeaderOne
	case level == 3:
		return draftjs.BlockHeaderTwo
	default:
		return draftjs.BlockHeaderThree
	}
}

// appendFootnotes emits a trailing `Footnotes` section when the document has
// definitions. It is a no-op for a post with no footnotes.
func (c *conversion) appendFootnotes(doc gast.Node) {
	var definitions []*extensionast.Footnote
	_ = gast.Walk(doc, func(n gast.Node, entering bool) (gast.WalkStatus, error) {
		if !entering {
			return gast.WalkContinue, nil
		}
		if list, ok := n.(*extensionast.FootnoteList); ok {
			for item := list.FirstChild(); item != nil; item = item.NextSibling() {
				if definition, ok := item.(*extensionast.Footnote); ok {
					definitions = append(definitions, definition)
				}
			}
			return gast.WalkSkipChildren, nil
		}
		return gast.WalkContinue, nil
	})
	if len(definitions) == 0 {
		return
	}
	c.flushBlock(&inlineBuilder{doc: c.doc, text: []rune("Footnotes")}, draftjs.BlockHeaderOne)
	for _, definition := range definitions {
		// Emit each of the definition's block children through the normal
		// inline mapper rather than flattening the definition with Node.Text:
		// flattening joins multiple paragraphs without a separator and drops
		// inline entities such as links, so the mapper is what keeps a
		// multi-paragraph or formatted definition intact.
		for child := definition.FirstChild(); child != nil; child = child.NextSibling() {
			switch n := child.(type) {
			case *gast.Paragraph:
				c.emitInlineBlock(n, draftjs.BlockOrderedListItem)
			case *gast.TextBlock:
				c.emitInlineBlock(n, draftjs.BlockOrderedListItem)
			case *gast.List:
				if err := c.walkList(n, blockContext{blockType: draftjs.BlockOrderedListItem}); err != nil {
					return
				}
			default:
				// A structured definition can nest a code block, block quote,
				// table, or thematic break. Dispatch it through the block
				// mapper under the definition's list context so its content is
				// preserved rather than silently dropped.
				if err := c.walkBlock(child, blockContext{blockType: draftjs.BlockOrderedListItem}); err != nil {
					return
				}
			}
		}
	}
}

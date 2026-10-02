package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/draftjs"
	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/markdown"
)

// convertSource converts an in-memory post whose images resolve against dir.
func convertSource(t *testing.T, dir, source string) (*markdown.Article, error) {
	t.Helper()
	converter := markdown.New(dir, "content/blog/example")
	converter.PostSource = "example/index.md"
	return converter.Convert([]byte(source))
}

// findBlockByType returns the first block of a type, or fails.
func findBlockByType(t *testing.T, article *markdown.Article, blockType string) draftjs.Block {
	t.Helper()
	return findBlockByTypeN(t, article, blockType, 1)
}

// findBlockByTypeN returns the n-th (one-based) block of a type, or fails.
func findBlockByTypeN(t *testing.T, article *markdown.Article, blockType string, n int) draftjs.Block {
	t.Helper()
	seen := 0
	for _, block := range article.Document.Blocks {
		if block.Type != blockType {
			continue
		}
		seen++
		if seen == n {
			return block
		}
	}
	t.Fatalf("no %s block at position %d (blocks: %v)", blockType, n, blockTypes(article))
	return draftjs.Block{}
}

func blockTypes(article *markdown.Article) []string {
	types := make([]string, 0, len(article.Document.Blocks))
	for _, block := range article.Document.Blocks {
		types = append(types, block.Type)
	}
	return types
}

// TestConversionFailureCases enumerates the ways conversion can fail and asserts
// each produces the named diagnostic at a source position.
func TestConversionFailureCases(t *testing.T) {
	cases := []struct {
		name   string
		source string
		code   string
		// write, when set, creates a file beside the post first.
		write func(t *testing.T, dir string)
	}{
		{
			name:   "no title",
			source: "---\ndescription: no title here\n---\n\nBody.\n",
			code:   "",
		},
		{
			name:   "unacceptable image media type",
			source: "---\ntitle: Example\n---\n\n![diagram](./diagram.svg)\n",
			code:   "image-media-type-rejected",
			write: func(t *testing.T, dir string) {
				writeFile(t, filepath.Join(dir, "diagram.svg"), []byte("<svg></svg>"))
			},
		},
		{
			name:   "escaping artifact locator",
			source: "---\ntitle: Example\n---\n\n![outside](../../secret.png)\n",
			code:   "image-not-readable",
		},
		{
			name:   "symlink whose target extension differs from the reference",
			source: "---\ntitle: Example\n---\n\n![mismatch](./image.png)\n",
			code:   "image-media-type-rejected",
			write: func(t *testing.T, dir string) {
				// The reference declares PNG, but the in-directory symlink
				// points at WebP bytes. Containment passes, so only classifying
				// by the reference extension catches the mismatch.
				writeFile(t, filepath.Join(dir, "actual.webp"), tinyWebP)
				if err := os.Symlink("actual.webp", filepath.Join(dir, "image.png")); err != nil {
					t.Fatalf("create symlink: %v", err)
				}
			},
		},
		{
			name:   "symlinked image outside the post directory",
			source: "---\ntitle: Example\n---\n\n![linked](./linked.png)\n",
			code:   "image-not-readable",
			write: func(t *testing.T, dir string) {
				// An apparently safe relative path whose canonical target lies
				// outside the post directory must be refused at conversion, not
				// read and then rejected only at publication.
				outside := filepath.Join(t.TempDir(), "outside.png")
				writeFile(t, outside, tinyPNG)
				if err := os.Symlink(outside, filepath.Join(dir, "linked.png")); err != nil {
					t.Fatalf("create symlink: %v", err)
				}
			},
		},
		{
			name:   "unrepresentable construct",
			source: "---\ntitle: Example\n---\n\n<div>raw html</div>\n",
			code:   "unrepresentable-html",
		},
		{
			name:   "over-budget markdown payload",
			source: "---\ntitle: Example\n---\n\n```text\n" + strings.Repeat("x", markdown.MarkdownPayloadBudget+1) + "\n```\n",
			code:   "markdown-payload-over-budget",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if tc.write != nil {
				tc.write(t, dir)
			}
			article, err := convertSource(t, dir, tc.source)
			if err == nil {
				t.Fatalf("expected conversion to fail, got success (blocks: %v)", blockTypes(article))
			}
			codes := diagnosticCodes(article)
			if tc.code == "" {
				// A missing title is refused before any diagnostic exists.
				if !strings.Contains(err.Error(), "title") {
					t.Errorf("expected a title diagnostic, got: %v", err)
				}
				return
			}
			if !hasDiagnostic(codes, tc.code) {
				t.Fatalf("expected diagnostic %s, got %s (%v)", tc.code, joinCodes(codes), err)
			}
			assertDiagnosticHasPosition(t, article, tc.code)
		})
	}
}

// assertDiagnosticHasPosition checks a diagnostic names a real source line.
func assertDiagnosticHasPosition(t *testing.T, article *markdown.Article, code string) {
	t.Helper()
	for _, d := range article.Diagnostics {
		if d.Code != code {
			continue
		}
		if d.Position.Line <= 0 || d.Position.Column <= 0 {
			t.Errorf("diagnostic %s reports position %+v", code, d.Position)
		}
		if d.Source == "" {
			t.Errorf("diagnostic %s names no source", code)
		}
		return
	}
}

// TestInlineStyleOffsets asserts bold, italic, and strikethrough produce ranges
// that select exactly the styled text, measured in UTF-16 code units.
func TestInlineStyleOffsets(t *testing.T) {
	source := "---\ntitle: Example\n---\n\nPlain **bold** and *italic* and ~~struck~~ end.\n"
	dir := t.TempDir()
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	block := findBlockByType(t, article, draftjs.BlockUnstyled)
	want := map[string]string{
		draftjs.StyleBold:          "bold",
		draftjs.StyleItalic:        "italic",
		draftjs.StyleStrikethrough: "struck",
	}
	for style, text := range want {
		selected := selectStyled(t, block, style)
		if selected != text {
			t.Errorf("style %s selects %q, expected %q", style, selected, text)
		}
	}
}

// TestInlineStyleOffsetAfterSupplementaryCharacter asserts a range after a
// supplementary character still selects exactly the styled span, which is the
// difference between UTF-16 and rune offsets.
func TestInlineStyleOffsetAfterSupplementaryCharacter(t *testing.T) {
	source := "---\ntitle: Example\n---\n\nEmoji \U0001F600 then **bold** end.\n"
	dir := t.TempDir()
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	block := findBlockByType(t, article, draftjs.BlockUnstyled)
	if selected := selectStyled(t, block, draftjs.StyleBold); selected != "bold" {
		t.Fatalf("style bold selects %q after a supplementary character, expected %q", selected, "bold")
	}
}

// TestLinkEntityRange asserts a link becomes an entity whose range selects the
// link text and whose payload carries the destination.
func TestLinkEntityRange(t *testing.T) {
	source := "---\ntitle: Example\n---\n\nSee [the docs](https://example.com/docs) now.\n"
	dir := t.TempDir()
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	block := findBlockByType(t, article, draftjs.BlockUnstyled)
	if len(block.EntityRanges) != 1 {
		t.Fatalf("expected one entity range, got %d", len(block.EntityRanges))
	}
	entity := article.Document.Entities[block.EntityRanges[0].Key]
	if entity.Value.Type != draftjs.EntityLink {
		t.Errorf("expected a link entity, got %s", entity.Value.Type)
	}
	if entity.Value.Data["url"] != "https://example.com/docs" {
		t.Errorf("link entity carries url %v", entity.Value.Data["url"])
	}
	if selected := selectRange(block, block.EntityRanges[0].Offset, block.EntityRanges[0].Length); selected != "the docs" {
		t.Errorf("link range selects %q, expected %q", selected, "the docs")
	}
}

// TestHeadingClamping covers all Markdown heading depths without emitting the
// schema-enumerated header-three type that the live draft endpoint rejects.
func TestHeadingClamping(t *testing.T) {
	source := "---\ntitle: Example\n---\n\n" +
		"# One\n\n## Two\n\n### Three\n\n#### Four\n\n##### Five\n\n###### Six\n"
	dir := t.TempDir()
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	want := []string{
		draftjs.BlockHeaderOne,
		draftjs.BlockHeaderOne,
		draftjs.BlockHeaderTwo,
		draftjs.BlockHeaderTwo,
		draftjs.BlockHeaderTwo,
		draftjs.BlockHeaderTwo,
	}
	got := blockTypes(article)
	if len(got) != len(want) {
		t.Fatalf("expected %d heading blocks, got %v", len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("heading %d mapped to %s, expected %s", i+1, got[i], want[i])
		}
	}
	for i, text := range []string{"One", "Two", "Three", "Four", "Five", "Six"} {
		if article.Document.Blocks[i].Text != text {
			t.Errorf("heading %d text = %q, want %q", i+1, article.Document.Blocks[i].Text, text)
		}
	}
}

// TestBlockAndNestedListMapping asserts paragraphs, both list kinds, nested
// list items, and block quotes map to their expected block types.
func TestBlockAndNestedListMapping(t *testing.T) {
	source := "---\ntitle: Example\n---\n\nParagraph.\n\n- alpha\n  - nested\n\n1. first\n\n> quoted\n"
	dir := t.TempDir()
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	types := blockTypes(article)
	for _, want := range []string{
		draftjs.BlockUnstyled,
		draftjs.BlockUnorderedListItem,
		draftjs.BlockOrderedListItem,
		draftjs.BlockBlockquote,
	} {
		if !contains(types, want) {
			t.Errorf("expected a %s block, got %v", want, types)
		}
	}
	// The Articles endpoint accepts no block depth, so a nested item is its own
	// list-item block after its parent rather than a deeper one.
	nested := findBlockByTypeN(t, article, draftjs.BlockUnorderedListItem, 2)
	if nested.Text != "nested" {
		t.Errorf("second unordered list item text %q, expected the nested item's own text", nested.Text)
	}
	if strings.Contains(nested.Text, "- ") {
		t.Errorf("nested list item text %q contains raw Markdown", nested.Text)
	}
}

// TestLooseListItemStaysOneBlock asserts a list item with a continuation
// paragraph — a loose item — becomes one list-item block whose text carries the
// item's paragraphs, rather than one block per paragraph, which DraftJS reads
// as two separate items and which changes the article's structure.
func TestLooseListItemStaysOneBlock(t *testing.T) {
	source := "---\ntitle: Example\n---\n\n1. first\n\n   continued\n\n2. second\n"
	dir := t.TempDir()
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	var items []draftjs.Block
	for _, block := range article.Document.Blocks {
		if block.Type == draftjs.BlockOrderedListItem {
			items = append(items, block)
		}
	}
	if len(items) != 2 {
		t.Fatalf("expected two ordered-list-item blocks, got %d: %v", len(items), blockTypes(article))
	}
	if !strings.Contains(items[0].Text, "first") || !strings.Contains(items[0].Text, "continued") {
		t.Errorf("first item %q does not carry both of the item's paragraphs", items[0].Text)
	}
	if strings.Contains(items[1].Text, "continued") {
		t.Errorf("the continuation paragraph leaked into the next item: %q", items[1].Text)
	}
}

// TestBlockquoteHeadingNestedInListIsReported asserts quote ancestry is tracked
// through a list, so a heading that is both quoted and listed — as in
// `> - # Warning`, where Goldmark nests the heading beneath a list item — still
// cannot escape the quote as an ordinary heading.
func TestBlockquoteHeadingNestedInListIsReported(t *testing.T) {
	for _, source := range []string{
		"---\ntitle: Example\n---\n\n> - # Warning\n",
		"---\ntitle: Example\n---\n\n- > # Warning\n",
	} {
		dir := t.TempDir()
		article, err := convertSource(t, dir, source)
		if err != nil {
			t.Fatalf("convert %q: %v", source, err)
		}
		if !hasDiagnostic(diagnosticCodes(article), "blockquote-heading-unrepresentable") {
			t.Errorf("source %q: expected a blockquote-heading diagnostic, got %v", source, joinCodes(diagnosticCodes(article)))
		}
		var carriers []draftjs.Block
		for _, block := range article.Document.Blocks {
			if strings.Contains(block.Text, "Warning") {
				carriers = append(carriers, block)
			}
		}
		if len(carriers) != 1 {
			t.Errorf("source %q: expected one block carrying the quoted heading, got %d", source, len(carriers))
			continue
		}
		if carriers[0].Type != draftjs.BlockBlockquote {
			t.Errorf("source %q: quoted heading lost its quote, emitted as %s: %q", source, carriers[0].Type, carriers[0].Text)
		}
	}
}

// TestListItemResumesAfterNestedList asserts an item that continues with
// another paragraph after a nested list — as in `1. first` / nested / continued
// — stays one parent item: its own paragraphs are one ordered-list-item block
// and the nested list follows, rather than the continuation becoming a second
// ordered-list item that DraftJS would read as a new entry and renumber. The
// Articles endpoint exposes no block depth, so the parent items are told apart
// from the nested items by their block type and order.
func TestListItemResumesAfterNestedList(t *testing.T) {
	source := "---\ntitle: Example\n---\n\n1. first\n\n   - nested\n\n   continued\n\n2. second\n"
	dir := t.TempDir()
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	var parents []draftjs.Block
	for _, block := range article.Document.Blocks {
		if block.Type == draftjs.BlockOrderedListItem {
			parents = append(parents, block)
		}
	}
	if len(parents) != 2 {
		t.Fatalf("expected two ordered-list-item blocks, got %d: %v", len(parents), blockTypes(article))
	}
	if !strings.Contains(parents[0].Text, "first") || !strings.Contains(parents[0].Text, "continued") {
		t.Errorf("first item %q does not carry both of the item's paragraphs", parents[0].Text)
	}
	nested := findBlockByType(t, article, draftjs.BlockUnorderedListItem)
	if nested.Text != "nested" {
		t.Errorf("nested list item text %q, expected the nested item's own text", nested.Text)
	}
}

// TestListItemNestedListPrecedesLaterConstruct asserts source order is kept when
// an item has text, then a nested list, then another block construct: the
// nested list must still precede the later construct, not be deferred past it.
func TestListItemNestedListPrecedesLaterConstruct(t *testing.T) {
	source := "---\ntitle: Example\n---\n\n1. first\n\n   - nested\n\n   ```go\n   code\n   ```\n"
	dir := t.TempDir()
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	nested := -1
	markdown := -1
	for i, block := range article.Document.Blocks {
		switch block.Type {
		case draftjs.BlockUnorderedListItem:
			if nested < 0 {
				nested = i
			}
		case draftjs.BlockAtomic:
			if markdown < 0 {
				markdown = i
			}
		}
	}
	if nested < 0 || markdown < 0 {
		t.Fatalf("expected a nested list item and a markdown block, got %v", blockTypes(article))
	}
	if nested > markdown {
		t.Errorf("nested list at %d follows the later construct at %d; source order was reversed: %v", nested, markdown, blockTypes(article))
	}
}

// TestBlockquoteHeadingIsReported asserts a heading inside a block quote — a
// construct X has no block type for — is reported rather than silently emitted
// as an ordinary heading outside the quote's context.
func TestBlockquoteHeadingIsReported(t *testing.T) {
	source := "---\ntitle: Example\n---\n\n> # Warning\n> quoted body\n"
	dir := t.TempDir()
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	if !hasDiagnostic(diagnosticCodes(article), "blockquote-heading-unrepresentable") {
		t.Errorf("expected a blockquote-heading diagnostic, got %v", joinCodes(diagnosticCodes(article)))
	}
	for _, block := range article.Document.Blocks {
		if block.Type == draftjs.BlockHeaderOne {
			t.Errorf("quoted heading was emitted as an ordinary header-one block: %q", block.Text)
		}
	}
}

// TestThematicBreakBecomesDivider asserts a thematic break is an atomic block
// backed by a divider entity rather than literal paragraph characters.
func TestThematicBreakBecomesDivider(t *testing.T) {
	source := "---\ntitle: Example\n---\n\nBefore.\n\n---\n\nAfter.\n"
	dir := t.TempDir()
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	divider := findBlockByType(t, article, draftjs.BlockAtomic)
	entity := article.Document.Entities[divider.EntityRanges[0].Key]
	if entity.Value.Type != draftjs.EntityDivider {
		t.Errorf("atomic block backs entity %s, expected a divider", entity.Value.Type)
	}
	for _, block := range article.Document.Blocks {
		if strings.TrimSpace(block.Text) == "---" {
			t.Errorf("thematic break leaked into paragraph text: %q", block.Text)
		}
	}
}

// TestInlineCodeKeepsTextAndReportsLostStyle asserts an inline-code span's text
// survives and the lost styling is a continuing diagnostic.
func TestInlineCodeKeepsTextAndReportsLostStyle(t *testing.T) {
	source := "---\ntitle: Example\n---\n\nUse `envsubst` here.\n"
	dir := t.TempDir()
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	block := findBlockByType(t, article, draftjs.BlockUnstyled)
	if !strings.Contains(block.Text, "envsubst") {
		t.Errorf("inline-code text does not survive: %q", block.Text)
	}
	if !hasDiagnostic(diagnosticCodes(article), "inline-code-style-lost") {
		t.Errorf("expected the lost code styling to be reported, got %v", diagnosticCodes(article))
	}
}

// TestFootnotesBecomeTrailingSection asserts references become bracketed text
// and definitions become an ordered list beneath a Footnotes heading, and that
// no heading is added when there are none.
func TestFootnotesBecomeTrailingSection(t *testing.T) {
	source := "---\ntitle: Example\n---\n\nStatement[^note] here.\n\n[^note]: The definition.\n"
	dir := t.TempDir()
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	bodyText := findBlockByType(t, article, draftjs.BlockUnstyled).Text
	if !strings.Contains(bodyText, "[1]") {
		t.Errorf("footnote reference is not bracketed text: %q", bodyText)
	}
	heading := findBlockByType(t, article, draftjs.BlockHeaderOne)
	if heading.Text != "Footnotes" {
		t.Errorf("expected a Footnotes heading, got %q", heading.Text)
	}
	definition := findBlockByType(t, article, draftjs.BlockOrderedListItem)
	if !strings.Contains(definition.Text, "The definition.") {
		t.Errorf("footnote definition is missing: %q", definition.Text)
	}

	without, err := convertSource(t, dir, "---\ntitle: Example\n---\n\nNo footnotes here.\n")
	if err != nil {
		t.Fatalf("convert without footnotes: %v", err)
	}
	for _, block := range without.Document.Blocks {
		if block.Text == "Footnotes" {
			t.Errorf("a Footnotes heading was added to a post with no footnotes")
		}
	}
}

// TestImageLocatorPerImage asserts each image becomes an atomic block with an
// unresolved image entity and its own locator, and that the entity's caption
// preserves the alt text.
func TestImageLocatorPerImage(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "one.png"), tinyPNG)
	writeFile(t, filepath.Join(dir, "two.png"), tinyPNG)
	source := "---\ntitle: Example\n---\n\n![first](./one.png)\n\n![second](./two.png)\n"
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	if len(article.Locators) != 2 {
		t.Fatalf("expected two locators, got %d", len(article.Locators))
	}
	captions := map[string]bool{}
	for _, locator := range article.Locators {
		if locator.Digest == "" {
			t.Errorf("locator %q records no digest", locator.Path)
		}
		entity := article.Document.Entities[locator.EntityKey]
		if entity.Value.Type != draftjs.EntityImage {
			t.Errorf("locator names entity %d, which is a %s entity", locator.EntityKey, entity.Value.Type)
		}
		captions[entity.Value.Data["caption"].(string)] = true
	}
	if !captions["first"] || !captions["second"] {
		t.Errorf("image captions did not preserve the alt text: %v", captions)
	}
}

// TestParserConstructs asserts the parser exposes the CommonMark and GFM
// constructs the posts rely on — tables, footnotes, and strikethrough — rather
// than assuming a default.
func TestParserConstructs(t *testing.T) {
	source := "---\ntitle: Example\n---\n\n" +
		"| a | b |\n| - | - |\n| 1 | 2 |\n\n" +
		"Struck ~~text~~ and a footnote[^x].\n\n[^x]: Definition.\n"
	dir := t.TempDir()
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	// A table preserves its rows and cell contents inside a markdown entity.
	table := findMarkdownEntity(t, article, "| a | b |")
	if !strings.Contains(table, "| 1 | 2 |") {
		t.Errorf("table entity lost its rows: %q", table)
	}
	// Strikethrough produces a style range.
	paragraph := findBlockByType(t, article, draftjs.BlockUnstyled)
	if len(paragraph.InlineStyleRanges) == 0 {
		t.Errorf("strikethrough produced no inline style range")
	}
	// Footnotes produce the trailing section.
	findBlockByType(t, article, draftjs.BlockHeaderOne)
}

// findMarkdownEntity returns the payload of a markdown entity containing want.
func findMarkdownEntity(t *testing.T, article *markdown.Article, want string) string {
	t.Helper()
	for _, entity := range article.Document.Entities {
		if entity.Value.Type != draftjs.EntityMarkdown {
			continue
		}
		payload, _ := entity.Value.Data["markdown"].(string)
		if strings.Contains(payload, want) {
			return payload
		}
	}
	t.Fatalf("no markdown entity contains %q", want)
	return ""
}

// selectStyled returns the text a style range selects.
func selectStyled(t *testing.T, block draftjs.Block, style string) string {
	t.Helper()
	for _, r := range block.InlineStyleRanges {
		if r.Style == style {
			return selectRange(block, r.Offset, r.Length)
		}
	}
	t.Fatalf("block %q carries no %s range", block.Text, style)
	return ""
}

// selectRange returns the substring a UTF-16 offset and length select.
func selectRange(block draftjs.Block, offset, length int) string {
	units := []rune(block.Text)
	start := 0
	seen := 0
	for i, r := range units {
		if seen == offset {
			start = i
			break
		}
		if r > 0xFFFF {
			seen += 2
		} else {
			seen++
		}
	}
	end := len(units)
	seen = 0
	for i, r := range units {
		if seen == offset+length {
			end = i
			break
		}
		if r > 0xFFFF {
			seen += 2
		} else {
			seen++
		}
	}
	return string(units[start:end])
}

func contains(values []string, want string) bool {
	for _, v := range values {
		if v == want {
			return true
		}
	}
	return false
}

func writeFile(t *testing.T, path string, content []byte) {
	t.Helper()
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// TestHardLineBreakBecomesNewline asserts both Markdown hard-break spellings —
// two trailing spaces and a trailing backslash — become a newline in the block
// text rather than merging the two lines into one word.
func TestHardLineBreakBecomesNewline(t *testing.T) {
	for _, tc := range []struct {
		name   string
		source string
	}{
		{name: "trailing spaces", source: "first  \nsecond"},
		{name: "trailing backslash", source: "first\\\nsecond"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			source := "---\ntitle: Example\n---\n\n" + tc.source + "\n"
			dir := t.TempDir()
			article, err := convertSource(t, dir, source)
			if err != nil {
				t.Fatalf("convert: %v", err)
			}
			block := findBlockByType(t, article, draftjs.BlockUnstyled)
			if block.Text != "first\nsecond" {
				t.Errorf("hard line break produced %q, expected %q", block.Text, "first\nsecond")
			}
		})
	}
}

// TestIndentedCodeStopsAtItsOwnLines asserts an indented code block followed
// immediately by a paragraph keeps its entity payload to the code lines and
// emits the paragraph once, rather than swallowing it into the entity.
func TestIndentedCodeStopsAtItsOwnLines(t *testing.T) {
	source := "---\ntitle: Example\n---\n\n    code line\nFollowing paragraph.\n"
	dir := t.TempDir()
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	payload := findMarkdownEntity(t, article, "code line")
	if payload != "    code line\n" {
		t.Errorf("indented code entity carries %q, expected %q", payload, "    code line\n")
	}
	paragraphs := 0
	for _, block := range article.Document.Blocks {
		if block.Type != draftjs.BlockUnstyled {
			continue
		}
		paragraphs++
		if block.Text != "Following paragraph." {
			t.Errorf("unexpected paragraph text %q", block.Text)
		}
	}
	if paragraphs != 1 {
		t.Errorf("following paragraph emitted %d time(s), expected exactly once", paragraphs)
	}
}

// TestListNestedConstructsArePreserved asserts a block construct nested in a
// list item — a fenced code block, here — is preserved as its own block rather
// than silently dropped because only inline children are walked.
func TestListNestedConstructsArePreserved(t *testing.T) {
	source := "---\ntitle: Example\n---\n\n- before\n\n  ```go\n  fmt.Println(\"nested\")\n  ```\n\n- after\n"
	dir := t.TempDir()
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	entity := findMarkdownEntity(t, article, "fmt.Println")
	if !strings.Contains(entity, "```go") {
		t.Errorf("nested fenced code lost its fence: %q", entity)
	}
	if !strings.Contains(entity, "nested") {
		t.Errorf("nested fenced code lost its body: %q", entity)
	}
}

// TestImageExtensionMustMatchBytes asserts a file whose bytes do not match its
// declared image extension is reported rather than accepted as that media type,
// so publication never sends an invalid payload to the upload endpoint.
func TestImageExtensionMustMatchBytes(t *testing.T) {
	dir := t.TempDir()
	// A real PNG renamed to .svg and a real SVG renamed to .png: neither
	// declaration matches the bytes.
	if err := os.WriteFile(filepath.Join(dir, "renamed.svg"), tinyPNG, 0o644); err != nil {
		t.Fatalf("write renamed svg: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "renamed.png"), []byte(`<svg xmlns="http://www.w3.org/2000/svg"></svg>`), 0o644); err != nil {
		t.Fatalf("write renamed png: %v", err)
	}
	for _, name := range []string{"renamed.svg", "renamed.png"} {
		source := "---\ntitle: Example\n---\n\n![diagram](./" + name + ")\n"
		article, err := convertSource(t, dir, source)
		if err == nil {
			t.Errorf("%s: a file whose bytes do not match its extension was accepted", name)
			continue
		}
		if !hasDiagnostic(diagnosticCodes(article), "image-media-type-rejected") {
			t.Errorf("%s: expected image-media-type-rejected, got %v", name, diagnosticCodes(article))
		}
	}
}

// TestFootnoteDefinitionKeepsBlocksAndFormatting asserts a footnote definition
// with multiple paragraphs keeps them separated and preserves an inline link
// entity rather than flattening everything to plain text.
func TestFootnoteDefinitionKeepsBlocksAndFormatting(t *testing.T) {
	source := "---\ntitle: Example\n---\n\nStatement[^note] here.\n\n[^note]: First paragraph with a [link](https://example.com).\n\n    Second paragraph.\n"
	dir := t.TempDir()
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	var definitions []draftjs.Block
	for _, block := range article.Document.Blocks {
		if block.Type == draftjs.BlockOrderedListItem {
			definitions = append(definitions, block)
		}
	}
	if len(definitions) == 0 {
		t.Fatalf("no footnote definition blocks were emitted (blocks: %v)", blockTypes(article))
	}
	joined := ""
	for _, block := range definitions {
		joined += block.Text
	}
	if !strings.Contains(joined, "First paragraph") || !strings.Contains(joined, "Second paragraph") {
		t.Errorf("footnote definition lost a paragraph: %q", joined)
	}
	if strings.Contains(joined, "First paragraphSecond paragraph") {
		t.Errorf("footnote definition paragraphs were joined without a separator: %q", joined)
	}
	hasLink := false
	for _, entity := range article.Document.Entities {
		if entity.Value.Type == draftjs.EntityLink {
			hasLink = true
		}
	}
	if !hasLink {
		t.Errorf("footnote definition dropped its link entity (blocks: %v)", blockTypes(article))
	}
}

// TestRasterExtensionMismatchIsReported asserts valid bytes for one accepted
// raster format renamed to another accepted raster extension are reported, not
// silently accepted under the sniffed type.
func TestRasterExtensionMismatchIsReported(t *testing.T) {
	dir := t.TempDir()
	// PNG bytes declared as JPEG; both types are individually accepted, so a
	// fallback to the sniffed type would let the mismatch through.
	if err := os.WriteFile(filepath.Join(dir, "photo.jpg"), tinyPNG, 0o644); err != nil {
		t.Fatalf("write mismatched image: %v", err)
	}
	source := "---\ntitle: Example\n---\n\n![photo](./photo.jpg)\n"
	article, err := convertSource(t, dir, source)
	if err == nil {
		t.Fatalf("PNG bytes declared as .jpg were accepted")
	}
	if !hasDiagnostic(diagnosticCodes(article), "image-media-type-rejected") {
		t.Errorf("expected image-media-type-rejected, got %v", diagnosticCodes(article))
	}
}

// TestFootnoteDefinitionKeepsNestedConstructs asserts a block construct nested
// in a footnote definition is preserved rather than silently dropped.
func TestFootnoteDefinitionKeepsNestedConstructs(t *testing.T) {
	source := "---\ntitle: Example\n---\n\nStatement[^note] here.\n\n[^note]: First paragraph.\n\n    - nested item\n\n    ```go\n    x := 1\n    ```\n"
	dir := t.TempDir()
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	nested := findMarkdownEntity(t, article, "x := 1")
	if !strings.Contains(nested, "x := 1") {
		t.Errorf("nested fenced code in a footnote was lost: %q", nested)
	}
	foundItem := false
	for _, block := range article.Document.Blocks {
		if strings.Contains(block.Text, "nested item") {
			foundItem = true
		}
	}
	if !foundItem {
		t.Errorf("nested list item in a footnote was lost (blocks: %v)", blockTypes(article))
	}
}

// TestInlineImageKeepsSurroundingStyle asserts an inline image splitting a
// styled span does not drop the style from the text before it.
func TestInlineImageKeepsSurroundingStyle(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "image.png"), tinyPNG)
	source := "---\ntitle: Example\n---\n\n**before ![alt](./image.png) after**\n"
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	var styled []string
	for _, block := range article.Document.Blocks {
		for _, r := range block.InlineStyleRanges {
			if r.Style == draftjs.StyleBold {
				styled = append(styled, selectRange(block, r.Offset, r.Length))
			}
		}
	}
	joined := strings.Join(styled, "|")
	if !strings.Contains(joined, "before") {
		t.Errorf("bold text before the inline image lost its style: %v", styled)
	}
}

// TestTaskListCheckboxIsPreserved asserts a GFM task-list item keeps its
// checkbox state rather than silently dropping it.
func TestTaskListCheckboxIsPreserved(t *testing.T) {
	source := "---\ntitle: Example\n---\n\n- [x] done\n- [ ] todo\n"
	dir := t.TempDir()
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	checked := findBlockByTypeN(t, article, draftjs.BlockUnorderedListItem, 1)
	unchecked := findBlockByTypeN(t, article, draftjs.BlockUnorderedListItem, 2)
	if !strings.Contains(checked.Text, "done") || !strings.Contains(checked.Text, "[x]") {
		t.Errorf("checked task item lost its state: %q", checked.Text)
	}
	if !strings.Contains(unchecked.Text, "todo") || !strings.Contains(unchecked.Text, "[ ]") {
		t.Errorf("unchecked task item lost its state: %q", unchecked.Text)
	}
}

// TestLinkedImageReportsLostDestination asserts a link whose only content is an
// image reports the lost destination as a continuing diagnostic: the image is
// still emitted, but the operator is told the hyperlink will not appear on X.
func TestLinkedImageReportsLostDestination(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "image.png"), tinyPNG)
	source := "---\ntitle: Example\n---\n\n[![alt](./image.png)](https://example.com/target)\n"
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	codes := diagnosticCodes(article)
	if !hasDiagnostic(codes, "linked-image-destination-lost") {
		t.Fatalf("expected a linked-image diagnostic, got %v", joinCodes(codes))
	}
	// The image itself survives the lost link.
	block := findBlockByType(t, article, draftjs.BlockAtomic)
	if len(block.EntityRanges) != 1 {
		t.Fatalf("expected the image entity range, got %v", block.EntityRanges)
	}
	var reported string
	for _, d := range article.Diagnostics {
		if d.Code == "linked-image-destination-lost" {
			reported = d.Message
		}
	}
	if !strings.Contains(reported, "https://example.com/target") {
		t.Errorf("diagnostic did not name the lost destination: %q", reported)
	}
	assertDiagnosticHasPosition(t, article, "linked-image-destination-lost")

	// The image may sit inside an inline wrapper such as emphasis, so the
	// enclosing link is an ancestor rather than the image's immediate parent.
	styled, err := convertSource(t, dir, "---\ntitle: Example\n---\n\n[**![alt](./image.png)**](https://example.com/target)\n")
	if err != nil {
		t.Fatalf("convert styled linked image: %v", err)
	}
	if !hasDiagnostic(diagnosticCodes(styled), "linked-image-destination-lost") {
		t.Errorf("a styled image inside a link must still report the lost destination, got %v",
			joinCodes(diagnosticCodes(styled)))
	}

	// A link that carries its own text must not be reported: the destination
	// survives as a link entity over that text.
	texted, err := convertSource(t, dir, "---\ntitle: Example\n---\n\n[take a look](./image.png)\n")
	if err != nil {
		t.Fatalf("convert texted link: %v", err)
	}
	if hasDiagnostic(diagnosticCodes(texted), "linked-image-destination-lost") {
		t.Error("a link carrying text must not be reported as a lost destination")
	}

	// When a link mixes text with an image, the text keeps the destination but
	// the image is still emitted unlinked. That partial loss must be reported
	// too, or the same image-in-link construct reports in one form and not the
	// other depending only on whether sibling text happens to exist.
	mixed, err := convertSource(t, dir, "---\ntitle: Example\n---\n\n[label ![alt](./image.png)](https://example.com/target)\n")
	if err != nil {
		t.Fatalf("convert mixed link and image: %v", err)
	}
	mixedCodes := diagnosticCodes(mixed)
	if !hasDiagnostic(mixedCodes, "linked-image-destination-lost") {
		t.Errorf("a link mixing text with an image must report the image's lost link, got %v", joinCodes(mixedCodes))
	}
	var mixedMessage string
	for _, d := range mixed.Diagnostics {
		if d.Code == "linked-image-destination-lost" {
			mixedMessage = d.Message
		}
	}
	if !strings.Contains(mixedMessage, "https://example.com/target") {
		t.Errorf("mixed-link diagnostic did not name the lost destination: %q", mixedMessage)
	}
	assertDiagnosticHasPosition(t, mixed, "linked-image-destination-lost")

	// The link's text still keeps the destination, so the destination must not
	// be dropped from the document merely because the image lost it.
	var urls []string
	for _, entity := range mixed.Document.Entities {
		if entity.Value.Type == draftjs.EntityLink {
			urls = append(urls, entity.Value.Data["url"].(string))
		}
	}
	if !contains(urls, "https://example.com/target") {
		t.Errorf("mixed-link destinations %v do not contain the surviving link target", urls)
	}
}

// TestLinkDestinationResolvesMarkdownEscapes asserts a link destination that a
// valid Markdown source escapes or entity-encodes is stored as the resolved URL
// rather than with the source encoding intact.
func TestLinkDestinationResolvesMarkdownEscapes(t *testing.T) {
	source := "---\ntitle: Example\n---\n\nSee [the docs](https://example.test/a\\(1\\)?x=1&amp;y=2) and [n](https://example.test/caf\u00e9).\n"
	dir := t.TempDir()
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	var urls []string
	for _, entity := range article.Document.Entities {
		if entity.Value.Type == draftjs.EntityLink {
			urls = append(urls, entity.Value.Data["url"].(string))
		}
	}
	want := "https://example.test/a(1)?x=1&y=2"
	if !contains(urls, want) {
		t.Errorf("link destinations %v do not contain the resolved %q", urls, want)
	}
}

// TestProseResolvesMarkdownEscapes asserts ordinary prose is emitted with the
// Markdown escaping and character references resolved, matching how a Markdown
// renderer presents it, while a code span keeps its literal source.
func TestProseResolvesMarkdownEscapes(t *testing.T) {
	source := "---\ntitle: Example\n---\n\nKeep \\*literal\\* and AT&amp;T and &#65; and &#x42;, but `co\\*de` stays raw.\n"
	dir := t.TempDir()
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	var text string
	for _, block := range article.Document.Blocks {
		text += block.Text
	}
	for _, want := range []string{"*literal*", "AT&T", "A", "B", "co\\*de"} {
		if !strings.Contains(text, want) {
			t.Errorf("emitted text %q does not contain the expected %q", text, want)
		}
	}
	for _, unwanted := range []string{"\\*literal\\*", "&amp;", "&#65;", "&#x42;"} {
		if strings.Contains(text, unwanted) {
			t.Errorf("emitted text %q still contains the unresolved source encoding %q", text, unwanted)
		}
	}
}

// TestImageCaptionResolvesProseAndKeepsRawCode asserts an image caption is
// resolved the way prose is — escaped punctuation and character references
// presented as a renderer would — while a code span inside the caption keeps
// its literal source, matching the guarantee ordinary prose already has. The
// caption comes from the flattened alt text, so resolving the whole caption as
// non-raw would strip the code span's backslash.
func TestImageCaptionResolvesProseAndKeepsRawCode(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "diagram.png"), tinyPNG)
	source := "---\ntitle: Example\n---\n\n![a \\*star\\* AT&amp;T and `co\\*de`](./diagram.png)\n"
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	if len(article.Locators) != 1 {
		t.Fatalf("expected one locator, got %d", len(article.Locators))
	}
	entity := article.Document.Entities[article.Locators[0].EntityKey]
	caption, _ := entity.Value.Data["caption"].(string)
	// A code span keeps its inner source without the delimiting backticks, so
	// the raw-code guarantee here is the literal backslash.
	for _, want := range []string{"a *star* AT&T", "co\\*de"} {
		if !strings.Contains(caption, want) {
			t.Errorf("caption %q does not contain the expected %q", caption, want)
		}
	}
	for _, unwanted := range []string{"\\*star\\*", "AT&amp;T", "co*de"} {
		if strings.Contains(caption, unwanted) {
			t.Errorf("caption %q still contains the unresolved or over-resolved form %q", caption, unwanted)
		}
	}
}

// TestLinkTitleLossIsReported asserts a link or image title, which the target
// entity cannot carry, is reported as a continuing diagnostic rather than
// dropped silently.
func TestLinkTitleLossIsReported(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "image.png"), tinyPNG)
	titled, err := convertSource(t, dir, "---\ntitle: Example\n---\n\nRead [the docs](https://example.test \"Reference\") now.\n")
	if err != nil {
		t.Fatalf("convert titled link: %v", err)
	}
	if !hasDiagnostic(diagnosticCodes(titled), "link-title-lost") {
		t.Errorf("expected a link-title diagnostic, got %v", joinCodes(diagnosticCodes(titled)))
	}
	var reported string
	for _, d := range titled.Diagnostics {
		if d.Code == "link-title-lost" {
			reported = d.Message
		}
	}
	if !strings.Contains(reported, "Reference") {
		t.Errorf("diagnostic did not name the lost title: %q", reported)
	}

	untitled, err := convertSource(t, dir, "---\ntitle: Example\n---\n\nRead [the docs](https://example.test) now.\n")
	if err != nil {
		t.Fatalf("convert untitled link: %v", err)
	}
	if hasDiagnostic(diagnosticCodes(untitled), "link-title-lost") {
		t.Error("a link without a title must not be reported as losing one")
	}

	titledImage, err := convertSource(t, dir, "---\ntitle: Example\n---\n\n![alt](./image.png \"Reference\")\n")
	if err != nil {
		t.Fatalf("convert titled image: %v", err)
	}
	if !hasDiagnostic(diagnosticCodes(titledImage), "image-title-lost") {
		t.Errorf("expected an image-title diagnostic, got %v", joinCodes(diagnosticCodes(titledImage)))
	}

	// The position must name the titled construct, not the paragraph's first
	// text: the link starts after "prefix " on the same line.
	positioned, err := convertSource(t, dir, "---\ntitle: Example\n---\n\nprefix [docs](https://example.test \"Reference\") now.\n")
	if err != nil {
		t.Fatalf("convert positioned link: %v", err)
	}
	// The body starts at column 1 and "prefix " is seven characters, so the
	// link text begins at column 9; the image's alt text begins at column 10.
	assertTitleColumn(t, positioned, "link-title-lost", 9)
	positionedImage, err := convertSource(t, dir, "---\ntitle: Example\n---\n\nprefix ![alt](./image.png \"Reference\") now.\n")
	if err != nil {
		t.Fatalf("convert positioned image: %v", err)
	}
	assertTitleColumn(t, positionedImage, "image-title-lost", 10)

	// A titled construct with an empty label or alt text has no text node of
	// its own. The position must then come from the construct's delimiter, not
	// the paragraph's first byte, so the diagnostic still names what produced
	// it rather than the paragraph start.
	emptyLabelLink, err := convertSource(t, dir, "---\ntitle: Example\n---\n\nprefix [](https://example.test \"Reference\") now.\n")
	if err != nil {
		t.Fatalf("convert empty-label titled link: %v", err)
	}
	assertTitleColumn(t, emptyLabelLink, "link-title-lost", 8)
	assertTitleOffset(t, emptyLabelLink, "link-title-lost", 8)

	emptyAltImage, err := convertSource(t, dir, "---\ntitle: Example\n---\n\nprefix ![](./image.png \"Reference\") now.\n")
	if err != nil {
		t.Fatalf("convert empty-alt titled image: %v", err)
	}
	assertTitleColumn(t, emptyAltImage, "image-title-lost", 8)
	assertTitleOffset(t, emptyAltImage, "image-title-lost", 8)
}

// assertTitleColumn asserts a title diagnostic names the construct on the line
// rather than the paragraph's first text.
func assertTitleColumn(t *testing.T, article *markdown.Article, code string, wantColumn int) {
	t.Helper()
	for _, d := range article.Diagnostics {
		if d.Code != code {
			continue
		}
		if d.Position.Column != wantColumn {
			t.Errorf("%s diagnostic column is %d, expected %d (the titled construct)",
				code, d.Position.Column, wantColumn)
		}
		return
	}
	t.Errorf("expected a %s diagnostic", code)
}

// assertTitleOffset asserts a title diagnostic reports the construct's own byte
// offset. For an empty label or alt text no text node exists, so that offset
// must be the construct's delimiter rather than the paragraph start.
func assertTitleOffset(t *testing.T, article *markdown.Article, code string, wantOffset int) {
	t.Helper()
	for _, d := range article.Diagnostics {
		if d.Code != code {
			continue
		}
		if d.ByteOffset != wantOffset {
			t.Errorf("%s diagnostic byte offset is %d, expected %d (the construct delimiter)",
				code, d.ByteOffset, wantOffset)
		}
		return
	}
	t.Errorf("expected a %s diagnostic", code)
}

// TestEmailAutolinkGetsMailto asserts an email autolink is recorded as a
// mailto URL rather than a bare address, which would be a relative link.
func TestEmailAutolinkGetsMailto(t *testing.T) {
	source := "---\ntitle: Example\n---\n\nWrite to <user@example.com> today.\n"
	dir := t.TempDir()
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	var urls []string
	for _, entity := range article.Document.Entities {
		if entity.Value.Type == draftjs.EntityLink {
			urls = append(urls, entity.Value.Data["url"].(string))
		}
	}
	want := "mailto:user@example.com"
	if !contains(urls, want) {
		t.Errorf("autolink destinations %v do not contain %q", urls, want)
	}
}

// TestImageReferenceResolvesMarkdownEscapes asserts a local image reference
// whose filename requires Markdown escaping is read from the unescaped path, so
// the file it names is found rather than reported unreadable.
func TestImageReferenceResolvesMarkdownEscapes(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "plot(final).png"), tinyPNG)
	source := "---\ntitle: Example\n---\n\n![plot](plot\\(final\\).png)\n"
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	if len(article.Locators) != 1 {
		t.Fatalf("expected one image locator, got %d (diagnostics: %s)",
			len(article.Locators), joinCodes(diagnosticCodes(article)))
	}
	if got := article.Locators[0].Path; got != "plot(final).png" {
		t.Errorf("locator records path %q, expected the resolved %q", got, "plot(final).png")
	}
}

// tinyWebP is the leading signature of a WebP file, enough for the converter's
// format check; the bytes are never decoded.
var tinyWebP = []byte{
	'R', 'I', 'F', 'F', 0x00, 0x00, 0x00, 0x00,
	'W', 'E', 'B', 'P',
}

// TestDiagnosticColumnCountsCharacters asserts a diagnostic's column is the
// character position in the source line, not the byte offset: text after a
// multi-byte character must not shift the reported column.
func TestDiagnosticColumnCountsCharacters(t *testing.T) {
	columns := map[string]int{}
	for _, tc := range []struct {
		name   string
		source string
	}{
		{"non-ascii prefix", "---\ntitle: Example\n---\n\ncaf\u00e9 `code` end.\n"},
		{"ascii prefix", "---\ntitle: Example\n---\n\ncafe `code` end.\n"},
	} {
		dir := t.TempDir()
		article, err := convertSource(t, dir, tc.source)
		if err != nil {
			t.Fatalf("%s: convert: %v", tc.name, err)
		}
		found := false
		for _, d := range article.Diagnostics {
			if d.Code == "inline-code-style-lost" {
				columns[tc.name] = d.Position.Column
				found = true
			}
		}
		if !found {
			t.Fatalf("%s: expected an inline-code diagnostic, got %v", tc.name, joinCodes(diagnosticCodes(article)))
		}
	}
	if columns["non-ascii prefix"] != columns["ascii prefix"] {
		t.Errorf("column after a multi-byte character is %d, expected %d (the character position)",
			columns["non-ascii prefix"], columns["ascii prefix"])
	}
}

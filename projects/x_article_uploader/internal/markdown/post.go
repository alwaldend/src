package markdown

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
)

// Post is a source blog post split into its YAML front matter metadata and its
// Markdown body. Front matter is metadata, never article text: draft creation
// requires a title that the content_state document cannot carry.
type Post struct {
	// Source is the path the post was read from, for diagnostics.
	Source string
	// Title is the parsed, non-empty front matter title.
	Title string
	// Body is the Markdown body with the front matter block removed.
	Body []byte
	// FileLine is the one-based line in the source file where Body begins, so a
	// diagnostic measured against Body can report the file's own line.
	FileLine int
}

// frontMatterDelimiter is the line that opens and closes the YAML block.
const frontMatterDelimiter = "---"

// frontMatter is the subset of the blog's front matter this project reads.
type frontMatter struct {
	Title       string `yaml:"title"`
	LinkTitle   string `yaml:"linkTitle"`
	Description string `yaml:"description"`
}

// ParsePost splits a source file into front matter and body, and returns the
// parsed title. A source without a non-empty title is refused here rather than
// producing a document that cannot be uploaded.
func ParsePost(source string, raw []byte) (*Post, error) {
	meta, body, fileLine, err := splitFrontMatter(raw)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", source, err)
	}
	if strings.TrimSpace(meta.Title) == "" {
		return nil, fmt.Errorf("%s: front matter has no non-empty title", source)
	}
	return &Post{Source: source, Title: meta.Title, Body: body, FileLine: fileLine}, nil
}

// splitFrontMatter returns the parsed front matter and the Markdown body. The
// body is the original bytes after the closing delimiter, so every byte offset
// goldmark reports stays valid against it.
func splitFrontMatter(raw []byte) (frontMatter, []byte, int, error) {
	var meta frontMatter
	if !bytes.HasPrefix(raw, []byte(frontMatterDelimiter+"\n")) {
		return meta, nil, 0, fmt.Errorf("source does not begin with a front matter block")
	}
	rest := raw[len(frontMatterDelimiter)+1:]
	closing := []byte("\n" + frontMatterDelimiter + "\n")
	end := bytes.Index(rest, closing)
	if end < 0 {
		return meta, nil, 0, fmt.Errorf("front matter block is not closed")
	}
	block := rest[:end]
	body := rest[end+len(closing):]
	if err := yaml.Unmarshal(block, &meta); err != nil {
		return meta, nil, 0, fmt.Errorf("parse front matter: %w", err)
	}
	// The body's first line is one past the closing delimiter, which is the
	// number of newlines the prefix consumed.
	fileLine := bytes.Count(raw[:len(raw)-len(body)], []byte("\n")) + 1
	return meta, body, fileLine, nil
}

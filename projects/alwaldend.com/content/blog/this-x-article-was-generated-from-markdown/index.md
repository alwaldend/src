---
title: This X Article was generated from markdown
linkTitle: This X Article was generated from markdown
date: 2026-09-30
description: A tour of the elements the X Articles endpoint accepts, exercised by the uploader's own conversion path
images:
  - pipeline.webp
tags:
  - x
  - markdown
  - repo
draft: true
---

This draft was written as Markdown and converted by `x_article_uploader`, which compiles a post into the DraftJS `content_state` the X Articles draft endpoint accepts. The post is also a test fixture: it contains every element the endpoint accepts, so one conversion exercises the whole vocabulary.

![The uploader pipeline](pipeline.webp)

## Text and inline styles

A paragraph is a plain text block, and X has no block nesting beyond the three heading levels below, so Markdown headings deeper than `header-three` clamp onto `header-three`. Inline styling is carried as ranges over the block text: **bold**, _italic_, and ~~strikethrough~~ all survive the conversion, and a span may combine them, as in **bold with _nested italic_ inside**.

### Heading levels

A third-level heading becomes `header-two`; a fourth-level heading and deeper become `header-three`.

#### A clamped fourth-level heading

## Lists

Unordered lists and ordered lists map onto their own block types:

- Draft creation resolves each referenced image.
- It uploads the bytes through the media endpoint.
- It attaches the returned `media_id` to the entity its locator names.

1. Convert the post offline.
2. Create the draft with the injected credential.
3. Review the draft before it becomes public.

A nested item is emitted as its own list item, because the endpoint's schema has no block depth:

- Parent item.
  - Nested item.
- Another parent item.

## Quotes

> X exposes a blockquote type, so quoted prose keeps its quote. A heading inside a quote stays quoted text and reports the lost heading level, because there is no quoted-heading form.

## Code

A fenced code block is preserved as a `markdown` entity, so the source travels intact:

```go
func (c Credentials) Sign(method, rawURL string, params url.Values, now time.Time) (string, error) {
    // HMAC-SHA1 over the normalized parameter string, per RFC 5849.
}
```

An inline-code span keeps its text but loses its monospace styling, and the conversion reports that loss.

## Tables

A GFM table has no native block, so it becomes a `markdown` entity that preserves the source:

| Element         | Entity type | Mutability  |
| --------------- | ----------- | ----------- |
| Link            | `link`      | `mutable`   |
| Image           | `image`     | `immutable` |
| Divider         | `divider`   | `immutable` |
| Code and tables | `markdown`  | `immutable` |

## Links and autolinks

A link carries its destination in a `link` entity: [the repository skill for Mermaid diagrams](https://example.com/mermaid). An autolink such as <https://example.com> becomes a link to the same destination.

## Divider

---

## Footnotes

A footnote reference[^draftjs] becomes a trailing section, so its text is not lost.

> **Note on images and captions.** An image becomes an `atomic` block backed by an `image` entity carrying its alt text as a caption, and the media is resolved at draft time.

## Emphasis edge cases

An underscore inside a word like snake_case stays literal, a backslash escape such as \*not italic\* stays literal, and an `entity` such as AT&amp;T is resolved to the character the author wrote.

[^draftjs]: DraftJS measures inline and entity ranges in UTF-16 code units of the block's final text, which is why the converter counts them that way.

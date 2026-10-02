---
title: This X Article was generated from Markdown
linkTitle: This X Article was generated from Markdown
date: 2026-09-30
description: Building a Markdown-to-X Article uploader, followed by a conversion showcase.
images:
  - pipeline.webp
resources:
  - src: pipeline.webp
    params:
      alt: The uploader pipeline
tags:
  - x
  - markdown
  - repo
---

I've had the displeasure of using X's Article editor and it has left me unsatisfied - it's utter dogshit. Naturally, I've decided to avoid it altogether. There probably exists a ready-made solution, but I've decided to just write it myself. After all, how hard can it be?

Not that hard, to be honest. X exposes an API to upload articles in [DraftJS](https://draftjs.org/) format, so that leaves two tasks:

- Build DraftJS
- Upload it

## Build DraftJS

Here is an example of a DraftJS payload that X expects:

```json
{
  "title": "My article",
  "content_state": {
    "blocks": [
      {
        "key": "title",
        "text": "My article",
        "type": "header-one",
        "inline_style_ranges": [],
        "entity_ranges": [],
        "data": {}
      }
    ],
    "entities": []
  }
}
```

Seems simple enough. Since I use Hugo to build my blog posts, I wanted to just render the payload with it - that did not work.

Hugo does not expose Markdown AST to templates, only raw text, so if I wanted to write a template for DraftJS, then I would have to parse raw Markdown with Go templating tools and then construct DraftJS from it which is theoretically possible, but insane and brittle.

Another option that seemed promising was [Render hooks](https://gohugo.io/render-hooks/introduction/) which override conversion from Markdown to HTML. It was not suitable because while they allow you to override some elements, it's still not full AST parsing.

Because Hugo was unsuitable, I've moved on to the boring option of just parsing Markdown in Go. I chose [Goldmark](https://github.com/yuin/goldmark/) since it's the parser Hugo uses, then wrote a tool that walks through the Markdown AST and emits DraftJS.

X does not allow you to inline media in the payload, so I had to add a couple of custom fields for the uploader to handle media correctly.

Here is an example of an artifact that is built from Markdown:

```json
{
  "payload": {
    "title": "Example",
    "content_state": {
      "blocks": [
        {
          "key": "block-0",
          "type": "unstyled",
          "text": "Hello world.",
          "inline_style_ranges": [
            { "offset": 6, "length": 5, "style": "bold" }
          ],
          "entity_ranges": [],
          "data": {}
        },
        {
          "key": "block-1",
          "type": "atomic",
          "text": " ",
          "inline_style_ranges": [],
          "entity_ranges": [{ "offset": 0, "length": 1, "key": 0 }],
          "data": {}
        }
      ],
      "entities": [
        {
          "key": "0",
          "value": {
            "type": "image",
            "mutability": "immutable",
            "data": { "caption": "Diagram" }
          }
        }
      ]
    }
  },
  "image_locators": [
    {
      "entity_key": 0,
      "path": "diagram.webp",
      "post_package": "out/article-payload/example",
      "digest": "422ba47ab0602f2875e0cbab8c5b2479a9140f1790122f6899885e870230de89",
      "media_type": "image/webp",
      "caption": "Diagram"
    }
  ],
  "banner_locator": {
    "path": "diagram.webp",
    "post_package": "out/article-payload/example",
    "digest": "422ba47ab0602f2875e0cbab8c5b2479a9140f1790122f6899885e870230de89",
    "media_type": "image/webp"
  }
}
```

The uploader uploads the media from locator fields and then patches the payload with media IDs.

## Upload it

The upload seemed like the most straightforward part to me before I started - just generate a token, put it in Vault, and hit the API, easy. None of it was easy.

Figuring out which token I need to generate to upload articles and how to do it through the dashboard took some effort. Either the dashboard is bad, or it was simply a skill issue, maybe both. I'm not sure.

Putting a token in Vault was harder than I expected because there are two options for Article upload: OAuth 1.0a and OAuth 2.0. OAuth 2.0 seems like an obvious choice - it's version 2, after all, but to utilize it I would need to use a short-lived token and then refresh it using a long-lived refresh token that is also rotated after a refresh. This seems like a very complex setup and I would need to update the secret pretty often, so I've settled on OAuth 1.0a. It worked alright, although it requires 4 different secret values and some request signing magic.

At this point I have my tokens, I have my DraftJS - I can just upload the media, and then the article, right? Wrong. X has a rolling limit of 10 article uploads per day which I immediately hit. After an investigation, I've figured out that X's API kept returning me `503 Service Unavailable` errors in response to my requests and I obviously kept retrying them, but those failed requests still count toward the limit. Very baffling behavior - there were no useful error messages, nothing, I just kept getting `503` errors out of nowhere. After some more research, the reason was the payload - I had a heading with depth 3, but X does not support it. After reducing heading depth, the problem was solved and I got a full Markdown to X upload.

## In conclusion

I don't have to touch the editor anymore, which is a win in my book. The blog post ends here. Everything else is just some Markdown to show conversion.

---

## Text and inline styles

A paragraph is a plain text block. The converter uses the two heading styles accepted by X's draft endpoint, mapping third-level and deeper Markdown headings onto `header-two`. Inline styling is carried as ranges over the block text: **bold**, _italic_, and ~~strikethrough~~ all survive the conversion, and a span may combine them, as in **bold with _nested italic_ inside**.

### Heading levels

First- and second-level headings become `header-one`; third-level and deeper headings become `header-two`.

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
| Code and tables | `markdown`  | `mutable`   |

## Links and autolinks

A link carries its destination in a `link` entity: [the repository skill for Mermaid diagrams](https://github.com/alwaldend/src/blob/master/tools/mermaid/skills/mermaid-diagrams/SKILL.md). An autolink such as <https://example.com> becomes a link to the same destination.

## Divider

---

## Footnotes

A footnote reference[^draftjs] becomes a trailing section, so its text is not lost.

> **Note on images and captions.** An image becomes an `atomic` block backed by an `image` entity carrying its alt text as a caption, and the media is resolved at draft time.

## Emphasis edge cases

An underscore inside a word like snake_case stays literal, a backslash escape such as \*not italic\* stays literal, and the HTML character reference `&amp;` in AT&amp;T is resolved to an ampersand.

[^draftjs]: DraftJS measures inline and entity ranges in UTF-16 code units of the block's final text, which is why the converter counts them that way.

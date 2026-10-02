## Why

Front matter `images` currently supplies social metadata only. Authors repeat
the image in Markdown to display it, and X conversion then removes selected
opening images. The user requested that Hugo render the metadata images above
the post instead, with explicit body duplicates removed from the source.

## What Changes

- Render the declared images before blog content, including the print edition.
- Remove the two existing matching opening image paragraphs and preserve their
  alternative text in resource metadata.
- Keep later body illustrations, social cards, post metadata, and prose intact.
- Update blog authoring guidance to describe the new rendering behavior.

## Capabilities

### Modified Capabilities

- `blog-section`: render front matter images as the post's opening illustrations.

## Impact

The shared Hugo layouts, two post sources, browser checks, and blog skill are
affected. No new dependency is needed. The companion uploader change is
`projects/x_article_uploader/openspec/changes/preserve-explicit-body-images`;
the uploader continues to use the first metadata image as an X banner.

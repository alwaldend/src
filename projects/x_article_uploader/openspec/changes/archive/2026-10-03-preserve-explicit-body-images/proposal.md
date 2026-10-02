## Why

The blog's opening images belong to its front matter rendering, while the
uploader uses the first metadata image as an X Article banner. Removing matching
Markdown images during conversion can erase body content the author explicitly
included.

## What Changes

- Preserve every explicit Markdown image, including opening images, repeated
  references, and files whose bytes match the selected banner.
- Keep metadata images separate from the body: the first entry supplies the
  banner without inserting any DraftJS image blocks.
- Retain image validation, body order, and shared media upload caching.
- Update existing HTTP E2E cases before removing suppression.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `markdown-to-draftjs`: Preserve explicit body images independently of banner
  metadata and remove the opening-image omission contract.

## Impact

The Markdown converter, its image helper, uploader README, and existing E2E
fixtures change. Artifact schemas, banner upload handling, and dependencies
remain compatible. The related Hugo owner change is
`projects/alwaldend.com/openspec/changes/render-blog-frontmatter-images`; it owns
rendering metadata images above blog bodies and removing redundant source images.
This uploader change makes no live X requests.

The existing `Banner image metadata` requirement is explicitly replaced by
`Independent banner and body image conversion`, retiring its opening-image
suppression scenarios while preserving banner and body validation guarantees.

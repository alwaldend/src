## Why

An opening blog image selected as the X Article banner currently appears twice
in the draft. Conversion should retain the cover while omitting its opening
body repetition, without removing later illustrations or changing the blog.

## What Changes

- Omit the consecutive opening body images whose validated bytes match the
  selected banner, including adjacent images in one image-only paragraph.
- Stop at the first nonmatching image or other content; preserve later repeats,
  mixed text, links, and images inside lists or quotations.
- Keep selecting only the first front matter `images` entry as the banner.
- Verify the converted artifact and draft HTTP payload with offline E2E cases.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `markdown-to-draftjs`: Omit opening repetitions of the selected banner.

## Impact

The Markdown converter, uploader README, and existing E2E fixtures change.
The artifact and HTTP schemas, media cache, site source, and dependencies stay
compatible. This work makes no live X API requests.

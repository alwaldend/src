## Why

Code files need portable PNG captures matching the repository website.

## What Changes

- Add a Bazel command that captures literal UTF-8 source using the shared site stylesheet.
- Support light and dark modes, scale, and a fixed width with line numbers and wrapping.

## Capabilities

### New Capabilities

- `code-capture`: Offline PNG capture of literal code using site theme tokens.

### Modified Capabilities

None.

## Impact

Adds tools/code_screenshot and exposes shared application theme sources to this tool. Reuses pinned Puppeteer, Chrome, Sass, and the existing Highlight.js version; promotes Highlight.js to a direct tools dependency without adding a package version.

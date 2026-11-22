## Why

The pinned Danbooru vocabulary provides prompt tags. A ComfyUI utility makes
them searchable while composing prompts.

## What Changes

- Add a CPU-only tag search node with a saved graph and native App.
- Package the node with the pinned vocabulary as an opt-in installation.
- Show ranked, copyable tag names for a partial or exact query.

## Capabilities

### New Capabilities

- `danbooru-tag-search` provides locally searchable prompt tags in ComfyUI.

### Modified Capabilities

## Impact

The ComfyUI tools catalog gains a utility workflow and App. The installation
target consumes the existing Danbooru dataset pin from `master`.

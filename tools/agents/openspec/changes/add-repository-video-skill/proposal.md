## Why

Repository video work needs a reusable procedure for selecting available media
tools, preserving requested audio, and verifying the resulting file. A recent
still-image export exposed an unavailable H.264 encoder and a dropped Opus
packet despite a successful encode.

## What Changes

- Add `repo-video` with tool discovery, input inspection, container and codec
  selection, stream mapping, verification, and artifact delivery guidance.
- Include a still-image and audio recipe that drains the complete audio source.
- Register canonical skill packaging, generated discovery, repository routing,
  and offline evaluation configuration.

## Capabilities

### New Capabilities

- `repository-video`: Guide video operations and verify their resulting artifacts.

### Modified Capabilities

None.

## Impact

The generic `tools/agents` owner holds the skill. `.agents/BUILD.bazel` owns its
discovery and `AGENTS.md` owns routing. Existing media tools and dependency pins
remain unchanged; user media and generated videos stay in ignored task output.

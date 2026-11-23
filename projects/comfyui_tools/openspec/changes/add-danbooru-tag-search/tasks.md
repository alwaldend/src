## 1. Search utility

- [x] 1.1 Implement local CSV lookup and package the node with the pinned vocabulary.
- [x] 1.2 Generate the graph and App and document installation and use.

## 2. Verification

- [x] 2.1 Exercise searches over the packaged vocabulary and check the generated workflow and text output.

The packaged CSV matched the pinned checksum. ComfyUI v0.43.0 ran the graph
on CPU, saved the seven requested results for `BLUE HAIR`, and served the
result through its native output endpoint. Local searches also checked ranking,
normalization, deprecated tags, result limits and empty or missing data.

## 3. Prompt formatting and sort order

- [x] 3.1 Add optional prompt spacing and popularity sorting to the node and expose both in the App.
- [x] 3.2 Verify both formats, both sort modes, and execution with the legacy inputs.

The bundled node passed a local fixture covering relevance and popularity
ordering, raw and spaced output, deprecation filtering, and the legacy
three-input invocation. A search over the pinned vocabulary returned five
spaced results using popularity sorting. The generated App exposes both options.

## 4. Build file generation

- [x] 4.1 Regenerate the node's BUILD file and verify the Gazelle workspace drift check.

## 5. Terminal access

- [x] 5.1 Expose the pinned vocabulary and shared ranking as a local CLI.
- [x] 5.2 Verify prompt spacing, raw output and limits through the packaged CLI.

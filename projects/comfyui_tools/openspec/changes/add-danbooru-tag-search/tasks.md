## 1. Search utility

- [x] 1.1 Implement local CSV lookup and package the node with the pinned vocabulary.
- [x] 1.2 Generate the graph and App and document installation and use.

## 2. Verification

- [x] 2.1 Exercise searches over the packaged vocabulary and check the generated workflow and text output.

The packaged CSV matched the pinned checksum. ComfyUI v0.43.0 ran the graph
on CPU, saved the seven requested results for `BLUE HAIR`, and served the
result through its native output endpoint. Local searches also checked ranking,
normalization, deprecated tags, result limits and empty or missing data.

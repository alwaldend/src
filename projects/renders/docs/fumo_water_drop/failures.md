# Fumo water-drop failures and lessons

[Back to goal](README.md)

## 2026-08-30 — Blender MCP unavailable

The Blender MCP execution endpoint could not connect to Blender at
`localhost:9876`. No Blender code ran and no scene state was inspected or
changed. This blocks an honest `.blend` scaffold, not the interface preflight.
Do not bypass the required Blender path with a different editor and then imply
that MCP inspected the result.

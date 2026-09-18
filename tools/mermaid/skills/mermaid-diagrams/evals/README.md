---
title: Mermaid diagram skill evaluations
---

The offline target validates that Promptfoo can load the skill and cases. The
cases cover the vertical flowchart default, keeping appearance in the shared
theme, direct documentation consumption, generated blog publication assets,
absolute paths and matching post-processing for the interactive CLI, container
containment, and hermetic rendering.

A live target is omitted because realistic evaluation requires rendering
diagrams through Bazel and comparing rendered geometry. Offline validation does
not claim that those renders succeed; the theme contract is covered by
`//tools/mermaid/test/renderer`.

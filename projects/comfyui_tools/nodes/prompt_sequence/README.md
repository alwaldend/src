---
title: ComfyUI workflow nodes
---

This package owns the dependency-free `Prompt Sequence` implementation, its
native multiline text box and offline checks. In `line` mode, each nonempty line
is one panel; blank and whitespace-only lines are ignored. Terraform stages its Python source through the existing native
ComfyUI provisioner; the repository Python library is only for offline checking.
One optional multiline negative prompt is shared across every panel without
splitting or adding the positive context. It is appended to the existing inputs;
the third output is scalar text, while prompt and filename outputs retain their
original list slots. Existing four-input workflows remain accepted with an empty
negative output.
See the [workflow instructions](../../README.md#prompt-sequences)
for panel input and the boundary between shared references and panel continuity.

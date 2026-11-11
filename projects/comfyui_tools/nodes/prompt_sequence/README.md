---
title: ComfyUI workflow nodes
---

This package owns the `Prompt Sequence` implementation, its
native multiline text box and offline checks. In `line` mode, each nonempty line
is one panel; blank and whitespace-only lines are ignored. Terraform stages its Python source through the existing native
ComfyUI provisioner; the repository Python library is only for offline checking.
One optional multiline negative prompt is shared across every panel without
splitting or adding the positive context. It is appended to the existing inputs;
the third output is scalar text, while prompt and filename outputs retain their
original list slots. Existing four-input workflows remain accepted with an empty
negative output.

One Run stays one native ComfyUI job, so its outputs stay grouped in media.
Every run gets a fresh UUID prefix, for example
`comic_5ba8d9cae1374e90aeed4482cf57841c_001`; the remaining panels use the same
UUID and increasing suffixes. Unchanged text queued again creates a new run
identity rather than reusing cached output filenames.

Current sequence workflows pass their model through the node's optional MODEL
input/output. A wrapper on a cloned native model patcher displays `Panel 2 of 3`
as native sampler progress text while the sampler's step counter continues to
work. Sampling and conditioning are unchanged. Older graphs retain their first
three output slots and still work; import the current workflow to add the model
connection required for panel progress. No browser extension is installed.
See the [workflow instructions](../../README.md#prompt-sequences)
for panel input and the boundary between shared references and panel continuity.

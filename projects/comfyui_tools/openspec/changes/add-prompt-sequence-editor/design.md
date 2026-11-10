## Context

The node currently accepts line or blank-line text. Its first-party Python
module is already staged through the native provisioner. ComfyUI supports
custom DOM widgets and extension lifecycle hooks; native Autogrow creates
connection sockets, not editable text boxes.

## Goals / Non-Goals

Provide the operator's automatically growing inline text fields. No modal,
reorder controls, third-party editor, new runtime package or live operation.

## Decisions

Replace the frontend prompts widget with one DOM widget at the same serialized
position. Keep one trailing blank textarea and collect nonempty fields in order.
Use nodeCreated and loadedGraphNode hooks to support both fresh and restored
nodes without monkey-patching framework prototypes.

Store a JSON string array in the existing prompts field under a new explicit
JSON separator mode. This preserves paragraphs and existing workflow field
positions. Backend validation accepts only arrays of strings; old separators
continue to work. A missing frontend extension leaves the normal text widget
and JSON mode available as a fallback.

Stage one optional web/index.js per approved first-party node module. Validate
that script keys refer only to those modules. Export WEB_DIRECTORY in the node.
No arbitrary file paths, Git-pack changes or additional runtime services.

## Risks / Trade-offs

Frontend API changes can affect rendering; use a browser fixture with the native
DOM-widget value contract, preserve source output and report runtime coverage
separately. Test failures before implementation: focus loss on growth, empty
panels, multiline splitting, changed ordering, stale reload state, duplicate
widgets, malformed JSON, injected markup and frontend assets missing from the
Terraform payload.

## Migration Plan

Publish source and packaging only. No rental exists and none is authorized.
A future deployment must load the packaged script and validate the target
ComfyUI frontend before claiming runtime acceptance.

## Duplication correction

Native node cloning restores widget values synchronously without invoking the
full graph-loading hook. Restoring prompts before the separator previously
rendered a saved JSON list using the default line separator. Defer rendering
to a microtask so the complete widget state is available; keep the existing
node lifecycle hooks and avoid patching framework prototypes. Verify the
cloning path in the existing browser fixture before changing implementation.

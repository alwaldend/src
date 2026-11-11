## Context

The first-party editor replaced ComfyUI's native text widget with an automatically
growing DOM widget. The operator rejected its usability and rendering after using
it, and requested the previous text box with blank lines for visual spacing.

## Decision

Use the existing STRING widget with multiline enabled. Remove WEB_DIRECTORY,
the browser script, its Bazel target and the rental's browser-asset input.
The existing line parser already drops blank and whitespace-only lines while
preserving ordered prompts and paired filenames; no parser change is required.
Keep the blank-line mode for paragraph prompts and JSON mode for previously
saved lists. Do not silently convert stored JSON or change panel boundaries.

The former browser fixture tests a removed editor rather than native ComfyUI,
so remove it. Existing backend checks cover blank omission, order, saved JSON,
paragraph prompts and rejection. Verify packaged Terraform inputs and generated
workflows without accessing the live rental.

## Deployment boundary

The operator initially prohibited rental inspection while using it. Source work
respected that boundary. The operator subsequently requested destruction,
replacement and tests of all workflows on the rented GPU. The owning
infra/image_gen change records those reviewed operations and neutral-fixture
acceptance; this component owns only the editor source and packaging.

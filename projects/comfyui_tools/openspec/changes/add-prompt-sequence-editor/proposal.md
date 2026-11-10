## Why

A single delimited text field makes panel prompts hard to edit and reorder.
Provide an automatically growing list of text boxes while retaining existing workflow and API inputs.

## What Changes

- Display one multiline text box per panel directly in the node.
- Add an empty box when the last box is filled; ignore empty boxes when running.
- Store edited lists in the existing prompts field as a JSON string array, with
  an explicit JSON separator mode that preserves blank lines within panels.
- Package the first-party script through the existing native provisioner.

## Capabilities

### New Capabilities

- `prompt-sequence-editor`: Editable ordered panels with lossless workflow
  storage, existing-list compatibility and paired output order.

### Modified Capabilities

None.

## Impact

The node, its browser extension and Terraform asset packaging change. No new
dependency or live deployment is required. Existing line and blank-line modes
remain supported, and users can continue editing raw lists without the extension.

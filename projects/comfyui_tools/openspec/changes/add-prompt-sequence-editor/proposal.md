## Why

The operator found the automatically growing custom editor awkward and poorly
rendered. Restore ComfyUI's standard multiline text box and allow visual spacing
without adding empty panels.

## What Changes

- Remove the custom browser widget and its deployment inputs.
- Keep one prompt per nonempty line by default, ignoring blank and whitespace-only lines.
- Preserve blank-line and JSON separator modes for existing saved workflows.
- Remove the browser fixture for the retired editor; retain backend behavior checks.

## Capabilities

### New Capabilities

- `prompt-sequence-editor`: Native multiline panel entry and saved-list compatibility.

### Modified Capabilities

None.

## Impact

Node declarations, asset packaging and documentation change. No new dependency
is added. The running rental remains untouched while the operator uses it.

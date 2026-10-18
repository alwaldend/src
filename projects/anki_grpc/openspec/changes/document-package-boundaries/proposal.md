## Why

Package documentation must identify native hosting, client lifetime and test ownership accurately so consumers preserve those boundaries.

## What Changes

- Review ankiclient, ankiserver and transport package responsibilities against source and consumers.
- Clarify lifecycle ownership, native hosting, excluded application policy and fixture responsibilities in doc.go.

## Capabilities

No behavior changes; skip_specs is enabled.

## Impact

Go package comments only. Preserve Rust execution, transports, lifecycle, write-back and existing dependencies.

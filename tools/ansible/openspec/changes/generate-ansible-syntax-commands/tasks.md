## 1. Packaged syntax commands

- [x] 1.1 Generate syntax targets and verify their argument and plugin wiring.
- [x] 1.2 Remove standalone checks replaced by generated targets and document
      the shared naming convention.

## 2. Validation

- [x] 2.1 Run every generated syntax target and validate the OpenSpec change.

All 31 syntax tests across 20 packages passed. The first run identified a
Hermes role import that needed the collection-qualified name.

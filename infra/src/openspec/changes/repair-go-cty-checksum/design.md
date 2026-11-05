## Repair

Use `GOWORK=off bazel run //tools/go:go -- mod download github.com/zclconf/go-cty@v1.13.1`.
The pinned Go command verifies the module through the configured checksum
database and writes its archive hash. Inspect the diff to ensure there are no
version changes or unrelated lockfile changes.

Workspace mode already finds this dependency through another workspace module
and leaves the root checksum file unchanged. Disabling workspace mode for this
one command writes the checksum to the root module consumed by Gazelle.

## Acceptance

The failing broad lint command is the existing behavioral reproducer; no new
unit test is needed. The standalone repair must pass repository quality and
semantic lint of the shared tools and owner record, and build the existing HCL
library that consumes go-cty. Retain logs under out/checksum-extraction and
publish the independent trunk-based candidate through AL delivery.

This repair is extracted from PR 211. That feature retains its copy until the
standalone fix merges, so extraction does not break the feature's lint gate.

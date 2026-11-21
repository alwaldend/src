## Context

The ruleset produces repositories through a Bzlmod extension. The monorepo
currently consumes it with a local module override. Its OpenSpec and
documentation targets also cross that module boundary.

## Decisions

Use the extension's root package label for monorepo consumers. Resolve the
flash rule's canonical label at the rule definition so generated repositories
load the same source regardless of where the extension is installed.

Build the archive from the same rule files. Preserve its existing public
package layout and publish a separate registry module declaration.

## Verification

Check root packages and representative generated repository targets. Build
the archive from a separate consumer using a temporary local registry to
exercise the published dependency graph.

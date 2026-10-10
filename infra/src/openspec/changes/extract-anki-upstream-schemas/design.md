## Context

The published project candidate is fde7f10d04987bfda9bb92424837e8e10e2f64ab. Its Anki dependency pins release 26.09.3 at commit 29bb700b951e3f0c0cb69b77c0180fc1fe33e6ba. Master is 2e42480c2e78dc2c17b0e7e95ba1d15abbb11461.

## Goals / Non-Goals

Extract the dependency, retaining upstream provenance and the schema import closure. Review fixes use upstream Go import paths, a tagged release download with integrity verification, and visibility limited to the local wrapper package. CLI code, TOML format, and schema version upgrades are outside this split.

## Decisions

Create the dependency branch from fresh master and copy only its four package files. Regenerate MODULE.bazel rather than copying the project branch's root configuration. Build both generated consumer aliases to verify compilation without requiring the unmerged CLI.

Both PRs target master, following the user's earlier explicit split preference. #128 temporarily retains byte-identical prerequisite files so its required checks remain valid. The new PR owns the dependency; after it merges, rebasing #128 removes the overlap.

## Risks / Trade-offs

Temporary overlapping files require merge ordering. Link the prerequisite in #128 and preserve identical dependency trees. Delivery receipts under ignored out/anki-schema-split own exact validation and publication state.

## Review fixes

Use `github.com/ankitects/anki/proto/anki/<schema>` for generated imports. Fetch the 26.09.3 tag archive and compute its own integrity; retain the resolved commit in provenance. External schema targets grant access only to `@//third_party/com_github_ankitects_anki:__pkg__`; sibling targets retain same-package access. Update #128 with identical dependency files and upstream imports so both candidates remain buildable.

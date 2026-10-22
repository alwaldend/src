## 1. Export command integration

- [x] 1.1 Adapt existing config-only and filesystem export E2E cases before production changes; verify their setup uses native fixtures and explicit state commands without reconciliation helpers.
- [x] 1.2 Extract App.Export and its dependency wiring; verify focused app/CLI/config E2E targets while plan/apply/build remain placeholders.

## 2. Documentation and delivery

- [x] 2.1 Document the current export command and verify strict OpenSpec; select integrated rendered documentation checks for delivery.
- [x] 2.2 Inspect the master-based diff and representative TOML/native outputs; select quality, affected-target lint and exact-candidate E2E/site gates. Delivery receipts own validation and publication results.

## 3. Review correction

- [x] 3.1 Extend the existing native E2E with source/destination directory aliases, reproduce bypassed containment, and resolve symlinks through the output's nearest existing ancestor; verify the focused app E2E passes.
- [x] 3.2 Carry the correction into #128, inspect representative outputs, and select fresh exact-candidate gates and the reasoned review reply. Delivery receipts own publication and resolution evidence.

## 4. Application boundary

- [x] 4.1 Define isolated command failure cases before implementation: application errors, flag/syntax rejection, context forwarding, state codecs and batched imports without native startup.
- [x] 4.2 Make commands consume a CLI-owned application interface; move state store construction, config loading and collection callbacks into App operations.
- [x] 4.3 Carry the boundary into #128, document ownership, and validate isolated commands alongside native E2E and representative artifacts before publication.

## Test ownership

- [x] Move config parsing/layout/export behavior into internal/config tests and collection input failures into internal/app tests. Keep syntax, flags, dispatch, codecs and streams in CLI tests; update build targets and references.

## Operation arguments

- [x] First adapt isolated CLI cases to require an interface-only root, forwarded explicit namespace arguments, default keys, and repeated commands without retained configuration.
- [x] Remove StateOptions and its mutable App field. Pass the namespace as a string argument through collection methods, keeping root flag parsing local to CLI and state interpretation in App. Update native callers and documentation on both branches.
- [x] Verify command isolation, native workflows and representative artifacts, then select exact-candidate delivery gates. Delivery receipts own publication and review evidence.

## Reusable collection scopes

- [x] First carry native failure cases for shared state/export operations, duplicate aliases, expired references, overlapping access and failed archive publication.
- [x] Add App-owned collection references and registry reuse while keeping plan, apply and build as placeholders; adapt CLI file selectors and document scoped lifetime.
- [x] Verify native and isolated command workflows and representative artifacts; select delivery gates with matching shared changes in #128. Delivery receipts own publication and review evidence.

## Collection command group and aliases

- [x] Adapt isolated dispatch, syntax and native CLI expectations before production edits: collection operations live under col, collection is an alias of col, config remains at root with alias conf, namespace flags belong to col, and removed root spellings fail before workflows.
- [x] Group collection commands and state under col using Cobra aliases; update help and examples, preserve command contexts/streams/required flags, and carry the same CLI into #128.
- [x] Verify command and native workflows and representative CLI output, then select fresh exact-candidate delivery gates. Delivery receipts own publication status.

## Application handlers

- [x] First adapt isolated CLI and native App callers to Config and Collection handler accessors; keep workflow dispatch, errors, context, state and file behavior expectations intact.
- [x] Reduce App and the CLI root interface to Config() and Collection(). Define handler contracts and move config/collection implementation and command factories to their owning files, retaining zero-value construction and scoped native ownership.
- [x] Document handler ownership and verify grouped CLI, native workflows and staged builds where implemented. Delivery receipts own exact-head checks and publication.

## Separate state workflows

- [x] Move state orchestration into its own handler, sharing the existing collection registry, prepared connections, namespaces and collection lifecycle through Collection().State().

## Public application contract

- [x] Adapt existing isolated CLI and native scope cases before implementation to use pkg/ankiascodeapi contracts, including caller-defined paths, nil/empty references, foreign ownership and expired references.
- [x] Move App and handler interfaces to pkg/ankiascodeapi and the implementation/native E2E to pkg/ankiascode with only protobuf and standard-library dependencies. Make CollectionRef an interface, retain private managed-reference state and validation, and update callers and Bazel dependencies without aliases or adapters.
- [x] Update documentation and specifications, verify native artifacts and the API dependency boundary, and select exact-candidate delivery gates. Receipts under out/application-api own validation and publication evidence.

## Flat namespaced application API

- [x] Adapt existing isolated CLI and native workflow callers before implementation to one App interface with Config, Col and ColState method prefixes; preserve command behavior, scoped references, namespace forwarding and failure expectations.
- [x] Remove handler interfaces, accessors and the separate API package; keep concrete App and collection reference contracts in pkg/ankiascode and define the application interface privately in the CLI. Implement namespaced operations directly on App and retain shared collection lifecycle and service boundaries. Update current documentation and specifications on both application branches.
- [x] Verify isolated CLI dispatch, native artifacts and the shared public package; select strict OpenSpec and exact-candidate delivery gates. Receipts under out/flat-application-api own final validation and publication evidence.

## Direct collection services and isolated CLI package

- [x] Adapt existing CLI and native workflow cases before implementation to plain collection paths; reuse the collection.Store E2E cases for direct prepared services, lifetime, aliases, namespaces and source publication.
- [x] Remove collection references, lookup IDs and App lifecycle wrappers. Reuse collection.Store and its injected Collection directly, keeping ordinary App operations path-based and existing service business rules intact.
- [x] Move Cobra commands and their private application interface to internal/cli; keep cmd/anki_as_code as composition and verify the CLI package does not depend on pkg/ankiascode or native backend packages.
- [x] Update docs and OpenSpec, inspect representative native artifacts, and select exact-candidate checks. Delivery receipts under out/direct-collection-cli own final validation and publication evidence.

## Configuration-set export command

- [x] Define command failure cases and a subprocess JSON/file/stdin/TOML round-trip artifact before implementation.
- [x] Add ConfigExport and config export through the existing codec and config.Store.Export; keep native lifecycle and semantic validation outside this operation.
- [x] Document the command, verify focused behavior and representative output, and select exact-candidate delivery gates on both application branches. Delivery receipts under out/config-export-cli own final validation and publication evidence.

## Prepared connection callbacks

- [x] Adapt existing native cases before implementation to locally reused services on the callback connection; retain persistence, cleanup, source-alias and publication coverage.
- [x] Carry the simplified collection.Store boundary from #165, remove Collection/service caches and construct application services on the injected grpc.ClientConnInterface.
- [x] Update current documentation and select exact-candidate native workflow, artifact, quality, lint and build gates. Receipts under out/simplify-collection own check execution and publication.

## Synchronize after collection-store merge

- [x] Rebase onto merged #165, preserve its typed App state-workflow E2E in the moved package and retain the feature's CLI/application boundaries.
- [x] Resolve visibility/relocation conflicts, update current ownership documentation and select exact-candidate quality, lint, native/API/CLI and site-build gates. Review and delivery records under out/rebase-after-165 own their results.

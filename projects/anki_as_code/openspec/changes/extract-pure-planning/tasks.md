## 1. Pure planning

- [x] 1.1 Write the native-backed planning failure and artifact cases before implementation, and retain persisted-state E2E coverage while moving preview logic into reconciliation.
- [x] 1.2 Extract a pure Plan function with private moves and bindings, simplify direct resource comparison, and verify planning after backend shutdown with unchanged inputs. Include native-validation deferral and filtered-card home-deck protection.
- [x] 1.3 Document the narrow API and deferred native effects, and verify the change with strict OpenSpec validation and native E2E artifacts.

Mandatory quality, semantic lint, consumer builds, exact-candidate review and publication are tracked by repository delivery receipts.

## 2. Simplification

- [x] 2.1 Split resolution and comparison into focused resource helpers and index notes by native identity without requiring GUID uniqueness.
- [x] 2.2 Verify native fixtures before implementing protection for unmanaged template deck references, then validate immutable plans and repeatable output.

- [x] Align plan counter documentation with deferred native card generation and structural migrations; known deletion counts cover deleted notes only.

- [x] Reproduce cross-note-type field/template transfers before restricting member moves to same-type renames; simplify the typed map helper and verify native identity retention.

## Executable planning boundary

- [x] Write native E2E cases for detached protobuf plans, execution without resource discovery, dynamic identity binding, stale-state rejection, input immutability and convergence before implementation.
- [x] Move pure planning to internal/plan and return a protobuf execution plan with prepared resource operations and a change report. Keep runtime resolution models private.
- [x] Make reconciliation execute supplied plans without configuration resolution or diffing. Let App compose capture, build and execution in one collection scope.
- [x] Update package docs and OpenSpec; select native artifact, quality, lint, consumer and documentation gates. Exact-candidate receipts own execution and publication.

## Review simplification

- [x] Write native-backed failure cases for member drift, empty converged operations, destination-only edits and path references before implementation.
- [x] Resolve stable member IDs or original stored names without stale ordinal fallbacks; reject ambiguous mappings.
- [x] Key desired objects by path, remove provisional IDs and prepare typed operations and reports together.
- [x] Emit only necessary card assignments; make expected state and proposed bindings explicit in the protobuf contract.
- [x] Verify native planning artifacts and the main PR's consumer E2E; select exact-candidate quality, lint, build and review gates.

## Filtered cloze ownership review

- [x] Capture native filtered cloze cards at multiple ordinals and reproduce unsafe template replacement before implementation.
- [x] Protect the shared cloze template's native identity, permit explicit identity-preserving template moves and verify native artifacts; select consumer and exact-candidate gates.

- [x] Add and inspect a maintained planning diagram showing App borrowing a collection connection, config/resource/state capture and pure plan construction; publish themed SVGs and a dark WebP preview.

## Empty namespace convergence

- [x] Capture absent and explicitly empty namespaces through the native state API, reproduce synthetic state changes before implementation and preserve empty-state convergence without input mutation.

## Tracked note-type replacement

- [x] Reproduce same-name replacement at a new config path through native fixtures before removing the redundant collision scan; retain untracked-name protection and verify native execution in #128.

## Native planning data

- [x] Extend the native E2E fixture to verify returned-payload ownership before implementation and capture the existing planning behavior.
- [x] Retain native snapshot messages, resolve desired configs directly into prepared protobufs, and reuse path/ID indexes across ownership checks and compilation.
- [x] Verify detached planning and execution, input independence, moves, member identities, filtered-card protection and convergence; select exact-candidate quality, lint, build and review gates for #168 and #128.

## Case-insensitive hierarchy ordering

- [x] Reproduce accepted mixed-case parent/child ordering before implementation and verify the prepared plan through Anki.
- [x] Order deck payloads using the same folded-name semantics as parent validation and verify native artifacts and final delivery gates.

## Tracked deck replacement and further simplification

- [x] Reproduce tracked deck name reuse through detached planning and native execution before removing its redundant name guard.
- [x] Cover independent directory moves, occupied destinations and typed preset binding ownership before removing duplicated maps and move passes.
- [x] Reuse typed preset bindings, check occupied destinations directly in typed resource maps and rename disjoint keys directly; verify native artifacts and select exact-candidate delivery gates.

## Direct planning projections

- [x] Verify declaration-order independence and filtered-deck kind projection in the native-backed fixture before changing implementation.
- [x] Remove the duplicate ordered deck collection and filtered-kind map; retain only resource bindings for moved lookup data and use the pinned typed protobuf clone API.
- [x] Verify unchanged reports, ordered operations, filtered-card protections, input ownership and native execution; select exact-candidate delivery gates.

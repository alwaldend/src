## 1. Behavior coverage

- [x] 1.1 Add borrowed-connection workflow and invalid-source coverage before production changes; verify initial compilation fails because the new API is absent.

## 2. Shared API and source resolution

- [x] 2.1 Add `ankiascodeapi` source contracts and update App/CLI composition; verify borrowed workflow, cancellation, source validation, path cleanup and CLI behavior with native E2E outputs.
- [x] 2.2 Rename schema directory, Go package and protobuf namespace; verify generated Go/Python/Java targets and all application consumers build.
- [x] 2.3 Update current documentation and architecture references; inspect source and rendered output for ownership accuracy.

## 3. Validation and delivery

- [x] 3.1 Run scoped native/CLI tests, semantic lint, repository quality and OpenSpec validation; inspect the produced archive/state artifacts and declared CLI dependency boundary.
- [x] 3.2 Review the final diff and prepare scoped delivery inputs; verify source-only task paths and a validation plan covering each affected package. Publication state belongs to the delivery receipt.

## 4. Review findings

- [x] 4.1 Reject typed-nil sources and connections before getter or RPC access; verify the expanded invalid-source workflow reproduces the panic before the fix and returns validation errors afterward.

## 5. Public application contract and synchronization

- [x] 5.1 Rebase onto current master, preserving explicit worker configuration alongside shared sources and renamed schemas.
- [x] 5.2 Move the complete application interface into `ankiascodeapi`, consume it in CLI and check concrete App compatibility in cmd; update API documentation and architecture.
- [x] 5.3 Verify existing native and command workflows, generated contracts, dependency boundaries, semantic lint and repository quality on the rebased candidate; prepare scoped delivery inputs. Publication state belongs to the delivery receipt.
- [x] 5.4 Name the public interface `ankiascodeapi.App`, update its consumers and documentation, and validate the renamed contract through existing workflows.
- [x] 5.5 Match the API directory, import path and Bazel target to `ankiascodeapi`; update consumers and verify the relocated package.

## 6. Reconcile current application boundaries

- [x] 6.1 Rebase onto current master and preserve its server command delegation and package documentation; expose the server workflows through `ankiascodeapi.App` while retaining per-call sources and schema renames.
- [x] 6.2 Verify server dispatch, native hosting, borrowed/path workflows, dependency boundaries and rendered documentation; prepare scoped delivery inputs for the reconciled candidate.

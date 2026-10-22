## 1. Behavioral coverage

- [x] 1.1 Before implementation, add memory-backed App configuration workflows covering complete nested round trips, independent instances, relative paths, fresh reads, malformed/missing inputs, read-only output and repeatable artifacts.
- [x] 1.2 Before implementation, add memory-backed build success/failure cleanup coverage and a borrowed native collection workflow exercising config export, plan, apply and state import from injected storage.

## 2. Implementation

- [x] 2.1 Add the shared filesystem adapter and explicit application/store injection; verify configuration workflows and existing OS layout/symlink cases.
- [x] 2.2 Route collection source checks, export/state output and rebuild staging/publication through the supplied filesystem; verify memory build and borrowed-native workflows plus a source scan for direct OS filesystem access.
- [x] 2.3 Expose the existing pinned Afero module, update generated module metadata and document filesystem paths/backend boundaries; verify builds and documentation examples.

## 3. Validation and delivery

- [x] 3.1 Record focused integration and strict OpenSpec evidence and select the required exact-candidate delivery checks. Delivery receipts own final gate and publication status.

## 4. Automated review follow-up

- [x] 4.1 Add a failing memory-backed source/state-output alias workflow before the fix, compare resolved paths independently of OS metadata, and verify rejection before backend startup without changing source or output.
- [x] 4.2 Record the review finding and focused validation evidence; final candidate checks, publication and review-thread response remain owned by delivery receipts.

## 5. Rooted filesystem aliases

- [x] 5.1 Before implementation, reproduce BasePathFs source/destination aliases through archive-build and collection-export workflows, including nested roots, relative names, dangling outputs and rejected external targets.
- [x] 5.2 Preserve the virtual namespace when resolving rooted symlink targets, verify the new workflows and existing alias protections, and record exact-candidate delivery and review evidence.

## 6. Shared filesystem boundary and complete review

- [x] 6.1 Before implementation, reproduce parent traversal during configuration export and state publication through directly nested BasePathFs; verify valid configuration paths and cleanup in the same workflows.
- [x] 6.2 Make internal/filesystem the shared application I/O boundary with private Afero storage, guarded paths and temporary filenames in the caller's namespace; pass that boundary through every internal owner.
- [x] 6.3 Review the complete PR, address confirmed in-scope defects, and record focused behavior/artifact validation and the final delivery plan. Delivery receipts own exact-candidate quality, lint, WASI and publication results.

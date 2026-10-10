- [x] Copy the shared implementation files without changes and verify exact extraction scope. Update skill task records to refer to their new delivery owner.
- [x] Validate the extracted candidate, existing consumers, and OpenSpec records. The final delivery receipt records exact-candidate revalidation.
- [x] Select the standalone prerequisite PR's scoped delivery against master. The delivery receipt owns publication and final verification state.
- [x] After the prerequisite merges, synchronize and validate PR #128's reduced diff.

The prerequisite's initial candidate f8599d527f204bd7a9d59182c27f59c09f6de396
passed all three validation groups: quality and tests, semantic lint, and existing
Go consumer builds. The implementation files match PR #128's 33cc18d candidate
exactly. The final receipt is under ignored out/anki-split and owns validation,
publication, and review observations. PR #128 was subsequently synchronized
after this prerequisite merged. Its project-owned collection-text-workflow
record and delivery receipts own subsequent validation and publication state.

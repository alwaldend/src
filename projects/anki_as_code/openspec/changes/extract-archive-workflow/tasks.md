## Native archive resources

- [x] 1. Record native lifecycle, ownership, and behavioral failure cases.
- [x] 2. Write native E2E before rewriting production archive handling.
- [x] 3. Delegate package operations to Anki and reduce Source to private collection lifetime
     and guarded filesystem publication; remove ZIP/Zstandard handling.
- [x] 4. Update README and build wiring, review the full diff, simplify, and fix findings.
- [x] 5. Format and verify the native implementation; select exact-candidate
     gates for #152 against master. Delivery receipts own final publication state.

Earlier manual extraction evidence is superseded by this native revision.
Current evidence and delivery receipts live under out/archive-native-api/.
Pending #128 files remain unpublished and must match their recorded hashes.

## Review corrections

- [x] 6. Accept a borrowed existing connection, remove worker lifecycle management,
     and verify success/failure/cancellation leave the worker reusable.
- [x] 7. Make the server target visible throughout the repository and select native
     workflow/build delivery gates. Receipts track publication and review replies.

## 1. Independent roots

- [x] 1.1 Extract NAS, development and AL roots, preserve behavior and update core identity lookups.
- [x] 1.2 Validate providers, packaging, backend selection and source boundaries.

## 2. State recovery

- [x] 2.1 Split the immutable recovery copy using Terraform state mv and verify resource/attribute conservation and smaller core state.
- [x] 2.2 Guard backend versions, push each state without force and verify ownership and no-change resource plans; keep Vault running.
- [x] 2.3 Remove protected migration copies after verified persistence and record sanitized results.

## 3. Delivery

- [x] 3.1 Resume and verify the authorized NAS/runner deployment; validate and review source plus sanitized migration/deployment receipts for publication.

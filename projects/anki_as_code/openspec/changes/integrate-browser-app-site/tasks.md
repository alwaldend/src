## 1. Site integration

- [x] 1.1 Simplify the themed page to one TOML editor, a file selector and Open JSON, Check and Download JSON actions; verify navigation and shared presentation.
- [x] 1.2 Automatically load the example and retain edits across file switches, checks and downloads; verify direct JSON opening and error recovery through the UI.
- [x] 1.3 Package runtime assets beside the Hugo page through the project content target; build the site and verify its generated HTML and resource URLs.

## 2. Verification

- [x] 2.1 Update the rendered-site browser E2E for the simplified workflow and produce desktop light/dark/mobile screenshots and a structured report; run it alongside the existing SDK E2E.
- [x] 2.2 Document the simplified page; run site/DNS checks, OpenSpec validation and formatting before exact-candidate delivery.

## 3. Collection and saved workspace integration

- [x] 3.1 Extend the rendered-site E2E for collection import, plan, apply, download and save/reload before implementation; retain the existing JSON and editor failure cases.
- [x] 3.2 Add themed collection controls using the reusable App API and stage replacement workspaces so failed imports preserve current edits.
- [x] 3.3 Add explicit local saving and restoration at startup, clear storage status and recoverable failure/conflict feedback while preserving the current JSON workflow.
- [x] 3.4 Verify the rendered workflow, light/dark/mobile presentation and emitted collection/screenshots/report; update documentation and delivery evidence. Final delivery gates remain owned by the delivery receipt.

## 4. Review fixes and simplification

- [x] 4.1 Extend rendered-site E2E to import JSON and collections with unavailable IndexedDB, retain downloads after save failures, and preserve conflicts after workspace replacement.
- [x] 4.2 Remove storage reads from local imports, retain the caller's save revision and unify staged workspace creation and disposal.
- [x] 4.3 Verify the site workflow and simplify current documentation before updating the PR.

## 5. Configuration import isolation

- [x] 5.1 Extend the rendered-site E2E before implementation to reject escaped configuration files with and without an open collection, preserving editor contents and exact package bytes through save/reload.
- [x] 5.2 Validate imported persistent files before copying the retained collection into the candidate workspace; reject files outside the editor directory without adopting the candidate.
- [x] 5.3 Verify the rendered workflow and record the emitted artifacts before exact-candidate delivery.

## 6. Reject imports into temporary storage

- [x] 6.1 Extend the rendered-site E2E before implementation with imports into reserved `/tmp` storage, preserving editor contents and package bytes through save/reload.
- [x] 6.2 Validate every exported file path before adopting imported configuration, including temporary paths, and verify the rendered workflow.

## 7. Large workspace adoption

- [x] 7.1 Extend the rendered-site workflow before implementation to restore a workspace whose file count exceeds the JavaScript argument limit, verify its editor and file selection, and confirm replaced workers are disposed.
- [x] 7.2 Prepare the file selector without argument spreading before adopting a replacement workspace and verify the workflow. Exact-candidate validation and publication are owned by the delivery receipt.

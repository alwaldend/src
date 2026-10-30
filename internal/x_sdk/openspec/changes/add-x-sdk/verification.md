## Verified behavior

- Live authenticated schema fetch succeeded on 2026-10-30 through the uploader injection. X API schema version 2.173, 200 paths and 871 component schemas. No secret values were written.
- OpenAPI Generator CLI reports 7.26.0; the downloaded JAR matched the publisher SHA-1 and is pinned with SHA-256 SRI.
- The user approved validator.v2 v2.0.1, requested generated clients remain build artifacts, and selected internal/x_sdk for repository-only use.
- The complete generated Go client compiles from a Go-only Bazel tree artifact. Client E2E verifies user lookup, bearer authentication, decoding, rejected authentication and base64 media upload against local HTTP fixtures. Command E2E verifies atomic schema updates and preservation for HTTP, redirect, size and JSON errors. Existing uploader E2E also passes with the shared signer.
- Generated client files were removed from source before any commit or push. Only the schema URL/checksum pin, generator configuration and handwritten source are retained.

## Final validation

All six integration test targets pass, including SDK and uploader E2E and all three affected OpenSpec workspaces. The main site builds and its rendered internal SDK reference documents the new commands. All 26 repository quality checks pass. Delivery receipts and exact-candidate check logs remain under ignored out/x_sdk.

## Session ergonomics

- XSDK-RUNFILES-001 (fixture-tested): corrected command E2E binary lookup to the existing rules_go runfiles resolver.
- XSDK-ARTIFACT-OWNERSHIP-001 (fixture-tested): initially staged generated client source; removed it at user direction and wired the generated Go-only tree directly into rules_go.
- XSDK-MEDIA-UNION-001 (fixture-tested): pinned upstream generator emits invalid pointer identifiers for three binary/base64 unions; narrow schema mappings plus raw-JSON aliases retain JSON upload values without editing generated code.
- XSDK-DISCOVERY-001 (live): early oversized skill reads obscured outputs; later reads used bounded artifact fields. No shared agent-policy changes are proposed.

## Mirrored schema revision

At the user request, the 44,259-line schema body was removed from source after
publishing its exact bytes to both download hosts. The versioned third-party
http_file pin carries SHA-256 SRI; SDK generation and E2E pass against this
anonymous mirror. Root now returns a public namespace index without redirects.
Deployment and checksum receipts are under ignored out/x_schema_mirror.

## Review correction acceptance

XSDK-PROBLEM-DECODING-001 is fixture-tested: the new local HTTP partial-success
case first failed with the multiple-oneOf-match error, then passed after enabling
discriminator lookup. Its partial-users.json artifact retains both the user and
the typed resource-not-found problem. Failure and success logs remain under ignored
out/x_sdk_review_fixes; generated code remains solely in Bazel outputs.

The correction reuses the existing task worktree and component change. An
oversized worktree inventory was avoidable; subsequent inspection used exact
paths. No shared policy change is proposed.

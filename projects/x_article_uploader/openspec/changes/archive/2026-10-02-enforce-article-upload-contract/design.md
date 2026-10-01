## Context

The user requested fixes for every documented review finding and exactly one
new upload attempt. The previous attempt at 2026-10-02T20:22:15Z returned 503
with nine daily slots remaining. Its static WebP and request shape matched
the documented contract; no finding established that response's internal cause.

## Goals / Non-Goals

Correct the documented contract gaps without adding dependencies or replaying
mutations. Preserve the existing Vault injection and unpublished-draft default.
Do not change credentials, query account identity merely to name a cache, or
retry draft creation automatically. Source delivery is authorized separately
from the one live attempt's observed outcome.

## Decisions

- Keep the supported media flow static-only. Reject animations rather than
  silently submitting them with `tweet_image`; adding animated or video upload
  is outside this change. Standard libraries decode supported formats where
  available; WebP receives bounded container/header checks without claiming a
  full entropy decode.
- Scope cache reuse to a versioned one-way fingerprint of the endpoint,
  category, and OAuth credential context. This conservatively invalidates reuse
  on token rotation as well as account changes and needs no identity request.
  Cache files contain no credential values. Legacy unscoped entries are misses.
- Retain only media with known future lifetimes across runs. Account for upload
  and processing time, reserve time for draft submission, and limit unknown
  lifetimes to the resolution in which the upload just succeeded.
- Poll read-only status only for reported pending processing, respecting the
  server's delay within a bounded deadline. Any failed, unknown, inconsistent,
  or incomplete processing result stops before draft creation. No upload or
  draft POST is replayed.
- Follow the recommended mutable Markdown entity mode. The Articles schema
  specifies a weighted limit but does not define its precise algorithm; the
  local budget must state its conservative assumptions rather than claim
  equivalence to an undisclosed backend validator.

## Decision review

Verdict: proceed with the bounded changes. An identity lookup could bind caches
to a stable user ID, but it adds a credentialed request and token permissions;
credential-context binding prevents cross-account reuse with a conservative
miss on rotation. Doing nothing preserves a possible invalid media reference.
Rejecting animations narrows acceptance but preserves content by failing
explicitly; pretending an animated file is a static image is unreliable.

## Verification and execution state

Write behavioral E2E failure cases before implementation. The root agent owns
serialized Bazel runs, formatting, validation, source delivery, and the single
live attempt. Raw credentialed output stays restricted and temporary; durable
receipts retain only operation/status, allowed response headers, and artifact
identity. Task-local evidence is under `out/x-article-upload-contract-fixes/`.

Verification passed: full uploader HTTP E2E and conversion command tests,
site output test, buildifier, affected Go semantic lint, and strict change
validation. Behavioral red runs preceded implementation for image validation,
cache scope/expiry, asynchronous processing, and Markdown mutability/weighting.
Independent review found an expired duplicate-image reference; the added
failure case reproduced it, and resolution now stops without another upload.
Legacy signature-only PNG fixtures were replaced with valid images; banner
assertions preserve the established contextual-error contract.

The final offline artifact has 47 blocks, seven entities, two mutable Markdown
entities preserving 468 source characters, and one 29,930-byte WebP banner with
no duplicate body image. UTF-16 ranges and the banner digest were verified.
Artifact SHA256:
`285a8a737c535c337f24da6513f6da7e18134ba599ba48fd3d5e4e100b4ecc69`.
Its ten continuing diagnostics describe the documented loss of inline-code
styling; there are no fatal diagnostics.

Exactly one authorized invocation ran at 2026-10-02T20:49:39Z. Media resolution
completed, then `/2/articles/draft` returned HTTP 503 Service Unavailable at
20:49:41Z, without a draft ID. The response reported 40,000/40,000 endpoint
requests remaining and 8/10 daily requests remaining. The daily reset was
2026-10-03T20:22:15Z (23:22:15 Moscow time); that is the quota reset, not a
prediction of service recovery. No Retry-After or internal failure explanation
was supplied. That authorization was used once.

The user then explicitly authorized one additional attempt without proxy
variables. At 2026-10-02T20:51:33Z, all upper- and lowercase HTTP, HTTPS, FTP,
ALL, and NO proxy variables were removed from the command environment and
explicitly unset at the uploader process boundary. It reused the ready scoped
banner cache and reached draft submission, then failed at 20:51:45Z with
`net/http: TLS handshake timeout` connecting to `api.x.com`. It returned no
HTTP response, quota headers, or draft ID. This was one additional invocation,
not an automatic retry. Both authorized attempts are complete; no article is
confirmed created or published. Further attempts need new authorization.

Task-local receipts: `source-review.json`, `artifact-verification.json`,
`pre-upload-checks.json`, `pre-upload-lint.json`, and
`single-upload-attempt.json` and `direct-upload-attempt.json` under the scratch
directory above. The exclusive
attempt marker prevents an accidental repeat during continuation. Source
publication requires the delivery adapter's final candidate validation.

NFC normalization reuses the repository's existing pinned `golang.org/x/text`
module; no new dependency or version was introduced. The weighted estimate uses
published v3 character weights, separately counts emoji components, and assigns
at least 23 units to plausible dotted domain fragments. It is deliberately
identified as an estimate, not a strict upper bound or exact Articles validator.

---
name: x-article-upload
description: >-
  Upload a blog post to X as an Article draft through the repository's
  x_article_uploader project, and understand the endpoint's rate and usage
  limits. Use when asked to upload, draft, or publish an article to X, or when
  a draft request is rejected with 429; do not use to edit the post's content
  or to make an article publicly visible.
---

# Upload an article to X as a draft

`projects/x_article_uploader` is the only supported path. It converts a post's
Markdown to the `content_state` the Articles endpoint accepts, uploads the post's
images, and creates an Article draft. It never publishes: a draft becomes public
only through a separate review step, so no invocation of this skill can make an
article visible. Read that project's `README.md` for the conversion contract and
the accepted image media types.

Do the conversion before the draft, and never export `X_*` values by hand.

1. Convert the post offline. This needs no credentials and no network, so it is
   the right place to catch an unrepresentable construct or a rejected image
   type:

   ```sh
   bazel_agent bazel run //projects/x_article_uploader/cmd/convert -- \
     --source projects/alwaldend.com/content/blog/<slug>/index.md \
     --out out/x-article-uploader/<slug>.json \
     --post-package projects/alwaldend.com/content/blog/<slug>
   ```

   A failing conversion prints every diagnostic and exits non-zero; a
   successful one prints the block and diagnostic counts, including any
   continuing `inline-code-style-lost` losses.

2. Create the draft through the Vault-injected wrapper, which is the only step
   that needs credentials and the only one that touches the network:

   ```sh
   bazel_agent bazel run //projects/x_article_uploader:draft -- \
     --artifact out/x-article-uploader/<slug>.json
   ```

   Pass workspace-absolute paths for `--artifact`, `--out`, and `--source` when
   running through Bazel, which may execute the tool from its runfiles tree. On
   success it prints `created draft <id> (not published)`.

3. Report the draft identifier and that the article is not published. A reviewer
   makes it public in the X composer.

## Rate and usage limits

The draft endpoint enforces two independent limits, and `429` does not say
which one was hit, so read the response headers before choosing how long to
wait. The uploader includes both sets of limit headers and `Retry-After`, when
present, in errors, and decodes valid reset timestamps into UTC. The standard
`X-Rate-Limit-*` headers describe the endpoint window, and
`X-User-Limit-24hour-*` describes the account cap.

- **24-hour per-user cap** (the usual cause): `X-User-Limit-24hour-Limit`,
  `-Remaining`, and `-Reset`. The observed allowance is 10 per 24 hours. When
  this counter is exhausted, wait until its reported reset (a Unix seconds
  timestamp); an earlier endpoint reset does not replenish this allowance.
- **Endpoint rate limit**: `X-Rate-Limit-Limit`, `-Remaining`, and `-Reset`.
  This window is minutes long and resets on the normal schedule.

Do not loop or retry a `429` automatically: the uploader already reports the
status and does not auto-retry, and the draft command must not be re-run to
probe the limit. A failed request, including a `503`, does not establish that
quota was preserved. Use the latest captured counters; an earlier positive
remaining value does not prove that a slot is still available.

If both counters are exhausted, use the later valid reset and respect any
`Retry-After` hint. Report the timestamp in the user's timezone, distinguishing
a reported waiting boundary from guaranteed upload success. Missing or invalid
headers leave the current waiting boundary unknown; label any older reset as
historical evidence. Wait for the reported reset, then run the draft step once
within the user's retry authorization.

A `429` here is not an authentication problem. Credentials, the AppRole, and the
media upload can all succeed while the draft call is capped: the media upload
endpoint has its own budget and is not part of the 24-hour Article limit, so a
cached `media_id` from a failed attempt stays reusable and only the draft call
needs to be repeated.

## Diagnosing a non-429 failure

A failure before any request is sent is a credential-path problem, not an X
problem:

- `could not create secret id for the approle: Vault request failed (HTTP 403)`
  means the host's Vault bootstrap token is invalid or expired, so the AppRole
  can never be logged into. Re-authenticate the host with `//tools/vault/login`;
  the uploader and its `al.lua` are not at fault, and no retry of the draft
  command will help until the token is renewed.
- `missing credential reference(s): X_...` means a field is absent from the
  injected environment; do not work around it by exporting `X_*` values, which
  would bypass the AppRole boundary.

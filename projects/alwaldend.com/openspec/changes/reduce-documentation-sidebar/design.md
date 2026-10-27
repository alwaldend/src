## Decisions

Enable Docsy sidebar_root_enabled and apply sidebar_root_for=self to section
pages matching /docs/_/_ through one Hugo cascade. This groups navigation by
repository owner without maintaining a project list. Docsy supplies parent
navigation and preserves foldable branches within each owner.

Compact and foldable settings only change visibility in pinned Docsy 0.17;
they do not prune generated HTML. Its shared chrome mode avoids duplication
but is experimental and documented for previews rather than production.
Keep full chrome. Upgrade the existing theme pin to Docsy 0.18.0 under the
user's subsequent request, using its canonical docsy/docsy repository. The
release declares Hugo 0.160.1 as its minimum and retains the existing pinned
Bootstrap and Font Awesome versions; verify compatibility with Hugo 0.165.0.

Render sidebar state per page by setting a high native cache threshold: cached
trees derive active state from href, but a scoped root links to its parent.
Keep the complete global taxonomy cloud and per-page tags so visitors can
navigate all tags. An A/B build measured the reduced cloud at 18,248,842 bytes
compressed and 208,661,474 bytes extracted, versus 34,574,345 and 271,245,906
with the full cloud. Restore the full cloud because the user prioritizes tag
navigation over that saving. Enable Hugo minification for release builds.
The user separately approved FontTools 4.67.0 to subset embedded diagram fonts;
that renderer change belongs to tools/mermaid and preserves standalone SVGs.

## Acceptance

Compare compressed archive size and uncompressed file totals with the measured
112,692,474-byte archive. Inspect project, tool, and infrastructure pages for
their own tree, parent links, current selection, and content. Verify the global
Docs entry, mobile navigation, theme behavior, and a nested page in a browser.
The setting must preserve all content and file downloads.

## First-interaction investigation (2026-10-26)

The user reports delayed first hover after every Firefox reload, and delayed
first finger tap in Vanadium. It affects both the sidebar and landing page.
The stable investigation identifier is `first-interaction-latency`.

The homepage background and sidebar pseudo-element are driven by CSS hover
selectors. The custom navigation script reveals the active row on initial
execution, DOMContentLoaded, and shown.bs.collapse; it has no hover, pointer,
or touch listener. This source observation does not rule out browser rendering
costs or other scripts in the user's environment. No speculative behavior
change was made.

The initial 36-case headless fixture did not reproduce the delay. That result
was invalid as a test of eager indexing: the fixture served JSON with an
application/octet-stream MIME type, so jQuery did not parse the document array.
It also measured from pointerover, excluding time spent waiting to dispatch
input. Corrected tests serve application/json and measure native X11 input
dispatch through the pointer event in headed Firefox private browsing.

On 2026-10-27, Firefox 158.0 reproduced first-hover delays of 667–712 ms
on the deployed 2026.44.0 homepage and documentation, with repeated hovers
at 34–62 ms. Instrumenting Lunr measured synchronous index construction at
756–776 ms. Input delivered during indexing was dispatched immediately after
index construction ended. Intercepting only the search document response as
an empty array reduced first-hover latency to 2–18 ms. The production JSON
contains 2,885 documents and occupies 3,965,455 bytes. The custom hover CSS
and navigation script are not the source of this blocking work.

Docsy 0.18.0 still constructs the Lunr index synchronously at page load;
its JavaScript conversion from jQuery to native APIs does not defer indexing.
There is no configuration option for prebuilt indexes or background indexing.
Search-result limits and excerpt lengths do not reduce searchable body text.
Keep search and the complete tag catalogue intact during the theme upgrade;
measure the new release before claiming any improvement. Physical Vanadium
with JIT disabled remains untested; slower JavaScript execution could amplify
this cost, but that device-specific explanation is not established.

The corrected harness and causal controls are retained under
out/hover-analysis/input-delay. The stable source references are Docsy's
v0.18.0 theme/assets/js/offline-search.js and upstream issue #1024; Lunr's
index_prebuilding guide describes precomputation as a separate future fix.

## Docsy 0.18.0 verification (2026-10-27)

The pinned archive SHA-256 SRI is recorded in hugo_lock.json. The site builds
with the unchanged Hugo 0.165.0 toolchain and Bootstrap/Font Awesome pins.
The site and Mermaid browser tests pass, covering theme selection, mobile
layout, diagram assets and print output. The browser fixture now serves JSON
with application/json to exercise local search correctly.

Headed Firefox 158.0 private browsing reproduced the delay on the rebuilt
homepage, Docs overview and owner-scoped documentation: 721–780 ms from native
pointer movement to the first pointerover, versus 10–12 ms on repeated hovers.
Lunr construction occupied 859–919 ms. These runs establish that the upgrade
does not fix the issue; different content and host load prevent interpreting
the old/new timing difference as a regression. All search content and tags
remain enabled. The upgraded site has not been deployed to either server.

Session review: the initial fixture MIME error and event-only timing obscured
the blocking workload. Native input dispatch and correctly served JSON are
required inputs for future reproduction. Firefox denied preload injection on
about:privatebrowsing; navigating the private tab to an ordinary local page
before instrumentation avoided needing privileged remote access. The durable
change record supersedes the earlier negative result; raw traces stay ignored.

## Independent website branch and responsive search

The user requested a new worktree and website-only pull request, including
all prior site fixes. Branch t3/update-docsy-site starts at master
e4ea1f28be358efe2a1a1f01ac782d9ff4086a0a. Website changes and their approved
FontTools dependency are ported from 736512db6301f8fc5efc5192ee2114606735c6fc;
infrastructure runtime configuration and deployment changes remain outside it.
The source-port manifest and fresh validation live under out/docsy-site.

Decision: proceed with build-time indexing plus a dedicated browser worker.
A fresh Firefox 158.0 private-window benchmark of 2,883 documents measured
751 ms to build, 55 ms to parse and 156 ms to load a serialized index.
Disabling baseline/optimizing JIT in the isolated browser profile increased
construction to 2,345 ms; parse/load still cost 68/217 ms. The serialized
index occupies 14,959,847 bytes. Prebuilding alone reduces CPU cost but leaves
visible main-thread stalls. Loading on search focus alone only moves the stall;
a worker also keeps query processing off the UI thread. The user approved
Lunr 2.3.9 as a build dependency on 2026-10-27.

Use Hugo segmented rendering to export Docsy's existing searchable document
template without restricting the available page graph. A Bazel action compiles
it with the existing fields and boosts, retains only result metadata alongside
the index, and packages the same pinned Lunr runtime for the worker. The final
Hugo pass fingerprints both generated index and worker. Preview, release and
staging use the same flags in their corpus and final rendering so draft policy
cannot diverge. Never rewrite a fingerprinted asset after Hugo renders it.

The site owns its search-input partial and search client; Docsy's eager offline
search is disabled. Initialize the worker on focus and retain it for the page.
Show loading/failure state, correlate asynchronous queries, and ignore stale
responses. Preserve literal-text result rendering and Docsy's search ranking,
result limits, excerpts, page exclusions and complete tag navigation. The
worker and library are served locally; no search CDN or server process remains.

Acceptance: full site browser tests and real search results; search remains
responsive while loading, errors are visible and retryable, stale results
cannot replace a newer query, and markup in queries is text. First hover/tap
and navigation must remain responsive before and during search initialization.
Compare native Firefox input timing, including a private profile with JIT
turned off, against the reproduced delay. Inspect release archives for project
binaries, generated release records, and raw intermediate search documents.

## Worker verification (2026-10-27)

The real Chromium search E2E check passes: no initial search/CDN request,
full-corpus results including download documentation, delayed loading with
query clearing, literal-text markup queries, retry after HTTP 503, all 127
global tags, parent navigation, and first mobile tap navigation. The existing
Mermaid/theme browser check also passes. Both checks now need no external
network; JSON uses its actual MIME type. Search screenshots wait for the
Bootstrap fade to complete.

Headed Firefox 158.0 private windows measured native first-hover dispatch at
6/13/13 ms for homepage/Docs/owner pages, compared with the previously reproduced
roughly 700 ms stall. With baseline and optimizing JIT disabled in an isolated
profile, the same pages measured 5/14/16 ms. Focusing search before hovering
while its worker loaded measured 11 ms with JIT disabled, with a maximum 24 ms
20-ms timer gap; the real query then returned the expected results. These are
controller measurements, not verification on the user's physical Vanadium
device. Traces and the native-input harness are under out/docsy-site/input-delay.

Hugo 0.165.0 coarse segment exclusions skip output formats before advancing
the page-output index. This made an excluded HTML format select the wrong
output while exporting search. Use include-only segment filters, including
the negated searchdocuments output for final rendering. The explicit output
filters preserve one shared page graph and avoid upgrading Hugo for this fix.

Session review: first-interaction-latency was reproduced using native input
after correcting the fixture MIME and event timing. Subsequent failed browser
cases exposed a real Enter retry issue (now explicitly handled) and a fixture
selection assumption (use Select All rather than triple-click). The build
export failure was localized to Hugo's coarse filter; include-only filters
produce the document corpus. Generated Python mappings were refreshed through
their owner. A direct release build refused missing versioning bootstrap; use
the documented bootstrap for stamped validation. No host configuration changed.

The independent release archive contains 3,499 HTML files and
2,883 searchable documents. It measures 35,230,302 bytes
compressed and 282,167,316 bytes extracted with the complete tag
catalogue. The prebuilt search asset occupies 15,367,051 bytes;
it loads only on search focus. All indexed references resolve to packaged pages,
release and staging reference lists match, and the archive contains the Lunr
license but no raw document-export JSON, generated project release pages, or
project download payloads. Explicitly mount generated static assets so the
worker library's license is distributed.

Review identified `search-result-blur-resubmission`: submitting with Enter did
not commit the input's native change event, so clicking a result submitted the
same query again and disposed the clicked link. The new first-click navigation
case reproduces the failure before the fix. Track each input's submitted value
and skip unchanged blur submissions; input editing clears that value, and
explicit Enter remains available for retries. The E2E check now follows a real
result link, rather than only inspecting its destination.

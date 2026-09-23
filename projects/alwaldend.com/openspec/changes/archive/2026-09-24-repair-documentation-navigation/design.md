## Context

The published owned-dns specification has an empty Hugo title and a second
body H1. Its parent directories have no section pages, so navigation skips
the OpenSpec hierarchy. The packaged source contains 56 Markdown documents
with a leading heading and no front matter.

## Decision

Proceed with a site-owned archive preparation step. It adds metadata only
to documentation whose first nonblank line is an H1 and creates missing
section indexes along those documents' paths. Existing metadata and files
remain authoritative. The heading render hook preserves the original anchor
without repeating the derived title.

Hugo content adapters cannot overwrite existing content deterministically;
the [upstream documentation](https://gohugo.io/content-management/content-adapters/#page-collisions)
explicitly describes page collisions. Template
fallbacks would require overriding every title consumer, including search
and feeds. Editing canonical OpenSpec files would duplicate their headings
as presentation metadata. Archive preparation keeps one derived title for
all Hugo consumers, at the cost of a small build step.

## Failure modes to verify

- Existing titles, draft state, metadata, source paths, and resources change.
- Generated section indexes collide with existing content or turn a leaf
  bundle into a branch bundle.
- A title is repeated, blank, incorrectly escaped, or taken from a code block.
- Original heading links break, including combined print output.
- Expanded branches overflow, overlap chevrons, or clip wrapped labels.
- Collapsed descendants remain focusable or retain visible highlights.
- Clicking a chevron navigates, or clicking a row expands the wrong branch.
- Mobile navigation obscures content, cannot close, or causes horizontal overflow.
- TOC anchors land beneath the header or their active state becomes stale.
- Taxonomy, release, print, and error layouts lose theme or spacing consistency.

## Verification

Use the complete packaged site and real Firefox interaction. Preserve the
repeatable browser script, screenshots, DOM measurements, and archive
comparison under `out/monochrome-theme/page-audit/`. No unit or
change-detector tests are added. Confirm actual viewport dimensions rather
than reporting requested dimensions: Firefox initially clamps narrow windows
to 500 CSS pixels.

## UI review findings

- Deep tree indentation squeezed specification names into several short lines;
  smaller steps retain hierarchy and readable labels.
- Tiny pointer-only chevrons now have a larger target and native keyboard toggling.
- Expanded mobile navigation now scrolls within a bounded region; desktop trees
  leave a bottom gutter, and the current page is revealed inside the tree.
- The tag index now wraps entries instead of placing one tag on each line.
- The 404 page now has content, recovery links, and space below the fixed header.
- Expanded release panels and print notices use the active neutral palette.
- Book tables retain covers and wrap long text within narrow screens.

## Reproduction

Build and serve `//projects/alwaldend.com:site_serve` with `bazel_agent`, then run
`python3 out/monochrome-theme/page-audit/verify_pages.py`. The task-owned script
writes unique screenshot directories and a JSON report with screenshot hashes,
actual viewport sizes, page measurements, and interaction results. Select one
color mode per process with `THEME_VERIFY_MODES=dark` or `light` for bounded
browser memory. `THEME_VERIFY_BASE_URL` selects the deployed site for a live
verification pass. Evidence includes homepage, Infra, project landing, and a
real blog article, plus OpenSpec, taxonomy, release, print, 404, book-table,
and generated-reference layouts.

## Acceptance evidence

The complete site build passed. Archive comparison preserved 26,917 original
files byte-for-byte, adding metadata to 56 plain Markdown documents while
preserving all their body bytes, and generating 83 missing section indexes.
A scan of all rendered HTML found no empty H1 headings.

Firefox verified 19 representative routes at 1440px and 390px in both light
and dark modes: 76 page views, plus nine interaction groups covering native
keyboard and pointer disclosure, hidden descendants, full-row hit areas,
visible active rows, TOC anchor navigation, release panels, and indexed titles.
`verification-combined.json` and `browser-input-hashes.json` under the task
output bind the screenshots to the rendered inputs. The browser harness has
bounded history; one earlier long session exhausted memory after 56 captures,
so the completed dark pass was retained and light mode verified separately.

A follow-up at 768px exposed an initially clipped active row: Docsy hydrates
its cached navigation on DOMContentLoaded. The shared navigation script now
reveals the active row after hydration, and tablet indentation leaves room for
its label. The 768px browser interaction pass verifies the row is fully visible,
keyboard and pointer toggling work, and the full-row hit target is preserved.

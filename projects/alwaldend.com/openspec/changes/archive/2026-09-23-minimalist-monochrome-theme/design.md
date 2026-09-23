## Context

See proposal.md. Docsy exposes pre-Bootstrap variable and post-Bootstrap style
hooks. All pages already share these hooks.

## Goals / Non-Goals

Use the existing shell for a coherent neutral palette and flat surfaces.
Preserve navigation, search, theme selection, and source content. This change
does not replace Docsy or deploy the public site.

## Decisions

Set neutral tokens before Bootstrap compilation so derived component colors
stay consistent. Runtime overrides retain the existing light/dark selectors.
Use a charcoal dark canvas, white light canvas, and understated gray surfaces.
Keep the existing system font stack, modest corner radii, and text hierarchy.
Present home links in a centered panel, with titles on the left and visible
destinations on the right. Indent project rows in the same flat list and keep
them permanently visible. Add space below the fixed header, and stack each
row’s title and destination on narrow screens. A separate theme or copied vendor templates would add ownership
and upgrade costs without improving this scoped presentation change.

## Risks / Trade-offs

Removing accent color weakens link recognition: underline prose links and
retain visible focus outlines. Docsy component specificity can hide palette
changes: inspect computed styles and representative pages in both modes.

## Verification

Observed 2026-09-23 from the task worktree based on `9b98eedb`:

- The site build passed with the existing Hugo toolchain and warning policy.
- Firefox screenshots confirmed the desktop projects page in light mode and
  the documentation page at 390 pixels in dark mode. Browser computed styles
  confirmed white/charcoal canvases and neutral text in both modes.
- At 390 pixels, the homepage, blog, project index, Bazel agent landing,
  documentation, and Hugo taxonomy fit without document-level horizontal
  overflow. Long documentation links initially overflowed; wrapping fixed it.
- The existing theme menu changed modes, offline search returned Hugo results,
  and a keyboard-focused home link displayed a 2px outline.
- Source review confirmed that all sections consume the same theme hooks;
  no vendor copies, dependency changes, or deployment operations were needed.

Raw build logs and Firefox screenshots are under `out/monochrome-theme/`.
Exact prepared-candidate build, lint, quality, and publication results belong
to the delivery receipts in that directory.

The first remote review identified `dark-outline-button-contrast`: Bootstrap's
compiled primary outline color matched the dark canvas. The shared theme now
binds the normal, hover, active, and disabled outline-button colors to the
runtime palette, including release attribute dropdowns. The revised candidate
requires fresh build, visual, lint, and quality validation in its receipts.

User screenshot feedback added thin, muted scrollbars and understated rounded
taxonomy labels. Standard scrollbar properties supply Firefox support, with
rounded thumb styling for WebKit. The taxonomy style removes Docsy's arrow
clip and nested count backgrounds; the project index uses the same class.
Updated screenshots and exact-candidate checks supersede the earlier visuals.

The homepage panel reuses configured GitHub and GitLab links and main-menu
Blog, Docs, and Projects destinations. Project rows come from the existing
Hugo section in title order. This keeps one project registry and requires no
custom JavaScript. The final layout removes the earlier disclosure control
and duplicate Projects heading, adds visible destinations, widens the panel,
and includes the fixed header's height in its desktop top spacing.

Rendered output verification covers the five requested links in order and
17 individual project links whose destinations exist in the built site.
Firefox screenshots at 1440x1280 and 390x1100 cover the revised panel layout.
The refreshed Infra screenshot shows thin muted scrollbar thumbs and rounded
taxonomy labels. The Autoscroll release screenshot confirms visible attribute
dropdown text and borders in dark mode.

To reproduce these visual artifacts, build the candidate and run
`bazel_agent bazel run //projects/alwaldend.com:site_serve`, select dark mode,
and capture `/`, `/docs/infra/`, `/projects/bazel_agent/`, `/projects/`, and
`/docs/projects/autoscroll/releases/head/`. Use the desktop dimensions above
for the homepage and 1440x1100 for the other pages; repeat `/` at the mobile
dimensions. Artifacts are `root-dark.png`, `root-mobile.png`, `infra-dark.png`,
`project-landing-dark.png`, `projects-refined.png`, and `release-refined.png`
under the task's ignored output directory. Delivery receipts identify the
validated revision.

## Documentation cohesion refinement

A subsequent Infra review identified remaining Docsy defaults: large bold
headings, prominent utility icons, heavy sidebar dividers, and dense tag
outlines. Shared overrides now use restrained headings, plain action rows
with subtle hover backgrounds, quiet active navigation rows, aligned chevrons,
and compact wrapping sidebar taxonomy items with inline counts.
Sidebar group titles use uppercase, spaced lettering to distinguish them
from their items, including the table-of-contents heading. Article and project
labels retain their rounded outlines. Search,
code, tables, quotes, and callouts share the same neutral surfaces and radii.
Project landings display their previously omitted title.

The first browser pass also exposed inline syntax colors that bypassed theme
CSS. Hugo now emits highlighting classes so the shared code palette controls
both modes. Browser verification and refreshed screenshots in the delivery
artifacts cover Infra, the homepage, project index, project landing, blog,
code samples, and mobile layouts. The isolated Firefox browser check passed
20 page captures across light/dark and desktop/mobile views, with no observed
page overflow or chromatic UI colors. Both theme-menu choices, mobile
navigation to Architecture, offline-search navigation, and a visible 2px
keyboard-focus outline passed. Its repeatable runner is
`python3 out/monochrome-theme/verify_theme.py` with the local server running;
`browser-verification.json` records the candidate, source dirtiness, viewports,
interaction results, and screenshot hashes. Final receipts bind the rerun to
the prepared candidate.

## Sidebar hierarchy decision review

Verdict: revise. The user rejected underlined actions and found the initial
heading changes indistinct. Uppercase, spaced group labels establish a clear
hierarchy while actions remain plain hoverable rows. A single-column taxonomy
variant clarified alignment but consumed too much sidebar height; the user
identified that trade-off as undesirable. The final taxonomy layout uses
compact wrapping items beneath the distinct group labels. Every destination
remains available without adding disclosure controls.

The capture runner saves every run to a unique directory to make revisions
unambiguous; identical image paths can obscure comparison even when the file
changes. Every render set includes a real blog post at
`/blog/dns-management-in-a-monorepo/`, as well as the homepage, Infra page,
and project landing. Screenshot metadata identifies the revision and paths.

## Session review

The agent-ergonomics review identified avoidable discovery output from a full
worktree inventory and a broad README search. Querying the supplied checkout
and site owner directly is sufficient. One formatter invocation used relative
paths in a Bazel runfiles directory; absolute source paths corrected it.
Preview snapshot and later automation failures were tooling errors, not site
failures; a task-local headless Firefox profile provided visual evidence.
These are execution corrections; no shared policy or host changes are proposed.

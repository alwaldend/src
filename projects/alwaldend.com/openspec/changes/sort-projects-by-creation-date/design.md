## Context

The homepage and project index currently use `.Pages.ByTitle`. Landing pages
are packaged through archives and do not declare dates.

## Goals / Non-Goals

Use explicit, manually maintained creation dates in both lists. Keep hosting,
page membership, and documentation navigation unchanged.

## Decisions

Each landing owns a mandatory ISO creation `date` in front matter. Initial
dates are manually assigned using project history and the original GitLab backup
repositories as references;
they are not inferred during builds. Corrections belong in the landing page.
The templates use `.Pages.ByDate.Reverse`, preserving existing presentation.
Manual weights are unnecessary. Filesystem timestamps would reflect archive
packaging rather than project creation.

Update the Hugo skill's authoring contract and existing evaluation cases.
Validate the built homepage and project index against the declared dates,
including unchanged project membership and representative documentation links.

## Risks / Trade-offs

Historical dates may approximate project creation when early work predates its
README or current directory. Maintained dates remain easy to correct.
New landing pages must provide the mandatory field; the skill records that
requirement. No automatic date-discovery mechanism is added.

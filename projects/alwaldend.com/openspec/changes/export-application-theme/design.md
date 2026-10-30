## Context

The site owns neutral light/dark colors, Bootstrap variables, and shared palette
selectors. A local tool needs those same sources without the full Hugo archive.

## Goals / Non-Goals

Provide a generated, offline application stylesheet with a single theme owner.
Do not alter public-site behavior, duplicate color values, or deploy the site.

## Decisions

Extract the existing palette mixins/selectors into a shared partial without
changing their output. Compose that partial with pinned Bootstrap and the shared
accent variables for applications. Move the small Sass action into a generic
site-owned stylesheet rule; retain Mermaid's public alias and existing outputs.
The cleaner embeds the generated CSS and owns only its workflow layout styles.

## Risks / Trade-offs

Validate the application stylesheet in a real browser in light and dark modes,
and build the site to catch changes to its Sass import context. Reusing Bootstrap
adds CSS size but avoids a second component/palette implementation.

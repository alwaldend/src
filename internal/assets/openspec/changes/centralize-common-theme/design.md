## Context

Hugo and Mermaid generate styles from the identity palette. Forgejo's existing
managed custom header applies mobile diff fixes and loads after its built-in
stylesheet.

## Decisions

- Keep the palette, common web styles, Mermaid configuration, and Forgejo theme
  in one internal package tree. Retain each consumer's loading and packaging
  adapters with their respective owners.
- Generate the custom Forgejo header from the shared palette. Apply the light or
  dark mapping based on Forgejo's selected mode and use the browser preference
  for Auto. Continue using Forgejo's built-in status and diff colors.
- Package the generated header at the existing Ansible controller path to reuse
  the role's staging, recovery, and checksum checks.
- Package the canonical monogram once and install it at Forgejo’s SVG favicon
  and logo override paths.
- Share the site’s navigation weight, prose underline offset, and keyboard
  focus dimensions through the common theme tokens.

## Risks / Trade-offs

Forgejo owns the CSS variable names and where the custom header loads. Confirm
the mappings and the rendered result against the repository's Forgejo release.

## Context

The favicon already consumes the canonical S8 SVG. See the proposal for the
missing header integration.

## Decisions

Enable Docsy's navigation logo and use its existing logo hook to consume the
same packaged artwork, rather than maintaining another SVG. Preserve the
shared header and its typography. Inspect both themes in real screenshots.

## Risks / Trade-offs

- Enabling the flag alone may select Docsy's default artwork. Inspect the
  emitted header and connect its logo hook to the canonical resource if needed.

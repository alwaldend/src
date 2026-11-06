## Context

The pinned OpenAPI schema already supplies template search and an instance
creation template hash. Reuse those generated operations and the bounded client.

## Decisions

Search by exact hash, select only necessary public metadata, and require one
matching result with an image and Docker environment options. Classify options
conservatively: only literal `-e`/`--env` assignments qualify as environment-only.
Support quoted whitespace and escaped characters without executing or expanding
shell expressions. Unknown flags and malformed quoting fail closed.

## Risks / Trade-offs

Future Docker option forms can be rejected until supported. This is intentional
for security consumers. Do not render template environment values in errors or
persist them in Terraform state. The template hash is passed through creation
only; existing instance refresh behavior remains authoritative.

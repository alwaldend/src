## Context

The Terraform reference already owns stage selection, packaging, injection,
and operation authority. Extend that reference rather than adding a parallel
procedure or a new infrastructure tool framework.

## Decisions

Document terraform_data with local-exec as a bounded fallback when an existing
provider lacks an operation and a small packaged tool is sufficient. Explain
that the tool owns retries and recovery, stdout is not a resource attribute,
and deletion and drift handling require an explicit design. Link the Forgejo
implementation as an example while leaving its service facts with its owner.

## Risks / Trade-offs

Provisioners can conceal side effects from plans. Keep invocation source and
non-secret lifecycle triggers reviewable, preserve normal interactive apply
review, and require isolated behavioral checks without live bootstrap calls.

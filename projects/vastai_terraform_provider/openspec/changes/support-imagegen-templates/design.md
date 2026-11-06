## Context

The existing resource requires an image and supports only image-based rentals.
The image-generation module needs to inspect and select existing templates.

## Decisions

Add `vastai_template` with a required hash and computed image and
`environment_only` safety classification. Add optional creation-only
`template_hash_id` to instances; retain image as a required input so ordinary
creation, refresh, and import remain compatible. The template lookup resolves
its image for callers. Startup overrides remain explicit instance settings.

## Risks / Trade-offs

Templates can contain untrusted application code. `environment_only` describes
Docker options, not application trust. The module uses the AI-Dock startup
contract; no GPU or image-generation result is claimed by offline validation.

Before implementation, exercise missing, duplicate, mismatched and malformed
template results, public Docker flags, quoted values, explicit template
selection at creation, refresh, and a replacement when the hash changes.

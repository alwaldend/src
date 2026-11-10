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

## Startup transport

The provider owns the 16,384-byte startup request limit. Add explicit
`compress_onstart` opt-in for images providing bash, base64, gzip and mktemp;
leave other consumers' execution behavior unchanged. Compress the complete
plain startup script with Go's standard library and use a private temporary
file with cleanup. Preserve the configured script in state when the API returns
the corresponding encoded command, while detecting externally changed commands.
Reject commands still exceeding the byte limit before the rental API call,
including when startup is unknown until apply. Diagnostics contain no script.

Before implementation, exercise literal script recovery and actual shell
execution, a no-change refresh, external command drift, and oversized plain or
incompressible input. Reuse the real Terraform/provider loopback fixture and
retain only non-sensitive sizes and outcomes as acceptance evidence.

Final suite validation exposed TEAM_FIXTURE_OWNER_KEY_COLLISION: the second
team's synthetic owner key used ID 100, also used by the first team's generated
management key. Revocation verification failed whenever the second team still
existed. Assign synthetic owner IDs above the fixture's generated-key range;
production key handling is unchanged.

The complete configured image-generation profile must also run through the
loopback Terraform/provider fixture. Consume only the public model configuration,
requirements, first-party node and generated workflow assets; exercise the real
module without Vault or rental credentials. Verify the final encoded byte count,
every selected workflow and model destination before any paid creation. This
closes the gap where small synthetic manifests passed but the deployed bundle
exceeded the API limit.

## Decisions

Expose one public API in pkg/client. Keep generation configuration at the SDK
root and generated protocol code in internal/generated. The public client owns
bounded authenticated operations, strict response checks, and safe diagnostics;
the generated implementation owns protocol serialization. The provider owns
Terraform behavior. No second public vastapi package remains.

New accepts an Options value with required Endpoint and APIKey fields. Generated
request rewriting preserves both decoded and escaped URL paths so a role name
containing a slash remains one path segment.

The user explicitly selected shared internal placement, including while the
Terraform provider is the initial runtime consumer. No independent product or
release surface is introduced. API and credential behavior remain unchanged.

The user requested a separate trunk-based PR including related tools, third_party,
and reusable tf_modules changes. Provider and infrastructure integration stay in
PR 249. Publish this prerequisite first; the dependent provider is updated after
the prerequisite lands, without a stacked pull request or implicit merge authority.

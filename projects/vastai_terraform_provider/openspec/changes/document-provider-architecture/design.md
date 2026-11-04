## Decisions

Use the left-to-right arrangement reviewed by the user. Terraform and Vast.ai / REST API
remain external nodes. The provider group contains Provider server,
Provider, Resources, and Data resources (Terraform data sources). Provider
configures the client and registers the two resource branches, both of which
use the SDK client. This diagram omits the internal providerdata package.
The Provider -> Resources and Provider -> Data resources edges denote registration.
Label Resources -> Vast.ai client as Read/Write and Data resources -> Vast.ai
client as Read, identifying the components that perform those operations.
The SDK group contains Vast.ai client and Generated client / OpenAPI.
The SDK client connects through the
generated client and directly to the external API for endpoints omitted from
the schema. This is a functional resource-request flow, not a package graph.
OpenAPI labels the generated client's origin rather than a runtime service.
Use native Mermaid flowchart layout and verify the rendered routing before
publishing.
Set diagramPadding to 48 for additional space on all four canvas sides.

The authoritative source is docs/architecture.mmd. Use the existing native
site Mermaid rules for light/dark SVGs and a dark WebP preview. Build outputs
are generated assets, not maintained source or separately copied illustrations.

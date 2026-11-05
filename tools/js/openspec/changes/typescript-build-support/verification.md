## Acceptance

The E2E driver was written before the fixture programs and build implementation.
The fixtures compile with the native compiler. A Node program imports an emitted
model module and writes the expected magnitude; a separate browser consumer
compiles against the model declarations and Three.js types, then draws visible
coordinates in pinned Chrome without external requests. The compiler rejects an
invalid assignment with TS2322. The browser canvas was visually inspected.

The fixtures use ordinary `.ts` source and explicit package module metadata,
as requested by the user. Acceptance additionally verifies that the manifest
is available at runtime. Source remains ordinary Git text without changing
the repository's media policy.

The test retains the Node result, browser PNG, type-error diagnostic, and a
passed JSON manifest with emitted-module and declaration hashes. Existing
Mermaid and website consumer tests passed. The OpenSpec record validates.

The pnpm workflow preserves every existing package version and adds the approved
compiler/type packages and their upstream dependencies. Adding Node declarations
also resolves existing optional peer references in the generated lock. The
TypeScript rules select an existing-compatible JavaScript rule version; their
telemetry dependency advances by one patch release through normal resolution.

Final delivery includes repository quality, the TypeScript E2E, existing native
Mermaid and website checks, and semantic lint of affected packages. AL task
evidence owns exact candidate identities, gate results, publication, and review.

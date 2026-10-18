## Context

The existing site render targets already produce light and dark SVGs and a WebP
preview from `docs/architecture.mmd`. See proposal.md for the motivation.

## Goals / Non-Goals

Show executable, linked-library, protocol, and storage boundaries with the
agreed labels. Keep package responsibilities available as a separate diagram.
App's built-in external-server selection remains work for another PR.
The base's shared API already accepts caller-prepared connections.

## Decisions

- Nest the linked `anki_grpc` graph and its Server and Anki API nodes inside
  `anki_as_code`, beside `ankiascode.App`. Add the `CLI` entry point inside the
  executable, with an arrow to App. Server delegates to Anki API, which
  accesses storage. Run (C ABI) targets the library; gRPC targets the Server. The
  subprocess edge targets the executable boundary to represent self-launch.
  This overview describes composition and entry points, not separate process
  instances, which would require duplicate App/library nodes.
  An invisible routing anchor outside the executable keeps the self-launch
  arrow visible; a direct edge to App's ancestor collapses in Mermaid. Route
  from App through the anchor below the executable and back to its boundary.
  Give the two loop segments named edges with `bumpX` curves so their horizontal
  tangents form a rounded turn at the anchor. Increase node spacing to 64 for
  a broader return. Other edges retain the site curve.
- Keep Configs outside the executable and storage graphs. App's Read/Write
  edge denotes operations on the text configuration files. Use a single
  arrow from App to Configs. The label describes the operations; the arrow
  identifies the component performing them. Allow more ranks for this edge and
  order Configs before the external server with an invisible link so both
  connections remain clear of nodes and the self-launch return.
- Label the linked call Run (C ABI), expressing the interface boundary on
  the edge without adding another component node.
- Use a solid gRPC edge for App's planned external-server selection and
  explain its planned status in nearby prose. Both servers point to the same
  conceptual Anki storage graph with Resources and State; this does not imply concurrent
  access to one collection.
- Preserve the old source as `docs/package-dependencies.mmd` and publish a
  separate render. Keep the main architecture asset URL stable.
- Configure the existing architecture WebP preview for dark appearance and
  natural proportions at scale 2, matching nearby documentation previews.
  Generate artifacts through Bazel instead of checking in image binaries.

## Risks / Trade-offs

- Nested/self-launch edges can cause poor Mermaid geometry. Inspect the actual
  render for readable titles, contained labels, and unobstructed edge captions.
- The overview combines composition and runtime paths. Explain self-launch and
  planned support next to the image so readers do not infer in-process RPCs or
  implemented external-server selection. Preserve the updated shared-API package
  graph from base `b31e3f3bca640c433c3bea9523a0411f75a0e83a` when relocating it.

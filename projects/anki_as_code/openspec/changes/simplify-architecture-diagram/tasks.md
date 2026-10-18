## 1. Diagrams and documentation

- [x] 1.1 Add the runtime overview and preserve the package diagram; verify all
      agreed nodes, boundaries, and edge labels in the rendered SVGs.
- [x] 1.2 Publish both diagrams in the README and explain self-launch and planned
      external support; verify the packaged documentation asset paths.
- [x] 1.3 Configure and render the dark architecture WebP; inspect it and show
      the generated image to the user.
- [x] 1.4 Add Anki API beside Server in the linked library; verify the rendered
      Server-to-API-to-storage path and show the updated dark preview.
- [x] 1.5 Add an external Configs node and App's Read/Write edge; verify both in
      the generated SVGs and updated dark preview.
- [x] 1.6 Round the self-launch loop without changing its endpoints and use a
      solid external-server gRPC edge; inspect both in the dark preview.

- [x] 1.7 Point the Read/Write arrow from App to Configs; verify the rendered
      edge has an arrowhead only at Configs.

- [x] 1.8 Add `CLI` inside the executable pointing to App; inspect the rendered
      entry point and update the dark WebP.

## 2. Validation

- [x] 2.1 Validate OpenSpec, formatting, and the documentation asset consumer;
      inspect the generated output and record actual results.

## Evidence

- Strict OpenSpec validation passed. Repository formatting and targeted
  formatting of the new artifacts completed without additional source changes.
- `//projects/anki_as_code:docs` and `//projects/anki_as_code/docs:architecture_preview`
  built successfully. Both SVG variants contain the expected labels and ten
  visible nondegenerate routed edge segments, including the self-launch return
  arrow, plus one invisible layout link.
  The loop segments meet with horizontal tangents, and the external-server
  gRPC edge is solid. The wider spacing leaves the loop caption unobstructed.
  App's Read/Write edge has its only arrowhead at Configs. CLI points to App
  inside the executable boundary.
- The dark WebP was visually inspected and shown in the task conversation.
  The retained package diagram matches the updated base source byte for byte.
- Delivery receipts own exact-candidate quality/lint gates, commit, push, and
  pull-request state; those records are kept under ignored task scratch.

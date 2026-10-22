## 1. Extract diagrams

- [x] 1.1 Convert all 13 remaining Drawio pages and reuse the 3x-ui source from #141; preserve topology, protocols, and archive classification.
- [x] 1.2 Publish generated site SVGs with each owning project's documentation and place ownerless archived diagrams in infra/archive.
- [x] 1.3 Remove infra/arch, obsolete validation registration, and stale operational references.

## 2. Verify and deliver

- [x] 2.1 Compare source inventories and inspect rendered outputs, including nested containers.
- [x] 2.2 Format and validate the exact candidate, semantic lint, repository quality, OpenSpec, and documentation site packaging.
- [x] 2.3 Prepare guarded publication and register the resulting pull request.

## 3. Separate 3x-ui views

- [x] 3.1 Create an overview with simple control/data nodes in datacenters and external services.
- [x] 3.2 Create a detailed control-plane view with one simple data-plane node, and a detailed data-plane view with one simple control-plane node.
- [x] 3.3 Package all three diagrams in project documentation, inspect their WebP renders, and select the exact-candidate delivery checks.

## 4. Simplify overview connections

- [x] 4.1 Group data-plane providers under one boundary, connect overview links to that group, use proper provider names without "datacenter", put node names below "Node", and inspect the rendered result.

## 5. Refine overview hierarchy

- [x] 5.1 Put User above Admin, rename the control node Master, nest PVE under Control Plane, split the PVE caption across two lines, show full node hostnames, and inspect the result.

## 6. Connect the control-plane boundary

- [x] 6.1 Attach incoming administration and outgoing management/ACME links to Control Plane rather than Master and inspect the compact render.

## 7. Refine control-plane captions

- [x] 7.1 Split both detailed control-plane container captions into two lines, name the VM "3x-ui VM", and inspect the render.

## 8. Publish one complete diagram

- [x] 8.1 Combine the overview, control-plane detail, and data-plane detail into one Mermaid source and render with separate readable sections; replace the three documentation assets and inspect the complete WebP.

## 9. Select dark appearance

- [x] 9.1 Select the site's dark palette for the full WebP and inspect its labels and 5:2 canvas.

## 10. Keep raster previews outside repository consumers

- [x] 10.1 Remove the preview-only WebP build target and README reference; retain build-time SVG documentation assets and the user's WebP under ignored task scratch.

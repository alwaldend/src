## 1. Dependencies and acceptance

- [x] 1.1 Upgrade and pin Mermaid and its CLI, reconcile the existing patch, and verify the installed versions and existing rendering checks.
- [x] 1.2 Write browser acceptance before implementation, covering source coverage, spatial geometry, routes, theme, controls, and screenshot/result artifacts.

## 2. Renderer and gallery

- [x] 2.1 Implement the version-pinned model adapters and spatial rendering pipeline; verify all six sample definitions produce valid scenes.
- [x] 2.2 Add hermetic HTML packaging, theme inputs, and interactive controls; verify offline rendering, both appearances, and camera/selection behavior.
- [x] 2.3 Wire the existing gallery sources to spatial outputs and document use; build the collection and inspect representative screenshots.

## 3. Delivery

- [x] 3.1 Run affected rendering tests and site integration, inspect outputs, and record acceptance evidence.
- [x] 3.2 Run final candidate quality and semantic lint gates, publish, and verify through repository delivery. Delivery receipts own exact validation and publication state.

## 4. Visual readability

- [x] 4.1 Extend browser acceptance for readable labels, projected label collisions, resize/framing, and phone screenshots before implementation.
- [x] 4.2 Improve side ports and spatial routes, anchored label placement, class hierarchy, and group contrast while retaining the shared theme.
- [x] 4.3 Provide readable scrollable framing and an explicit fitted overview; inspect all six samples in both themes and on phones.
- [x] 4.4 Record visual acceptance and the delivery validation scope; delivery receipts own final gates and exact publication state for the existing pull request.

## 5. Architecture documentation

- [x] 5.1 Add one authoritative architecture diagram, render it with the maintained theme, and embed it in the README and project landing with an interactive 3D link.
- [x] 5.2 Inspect the diagram and built pages, verify asset links and rendering, and record the delivery validation scope.

## 6. Curved connections

- [x] 6.1 Extend browser acceptance for smooth bends and curve clearance before implementation.
- [x] 6.2 Curve graph connections and round routed bends and sequence self-calls, preserving ports, markers, dash spacing, and the shared theme.
- [x] 6.3 Inspect the six samples and architecture view and record the delivery validation scope; delivery receipts own final gates and publication.

## 7. Review corrections

- [x] 7.1 Add browser coverage for selection scope, arrowless messages, notes in nested and note-only blocks, and numeric participant declaration order before implementation.
- [x] 7.2 Preserve selection scope and sequence semantics, including note membership and complete group volumes.
- [x] 7.3 Inspect affected renders and record validation evidence; delivery receipts own final gates and publication.

## 8. Rendering refinements

- [x] 8.1 Add browser acceptance for multiline sequence cards, offline payload size, rendering cost, and compact mobile overview before implementation.
- [x] 8.2 Share text layout between sizing and drawing, embed the Three.js core once, and give diamond nodes a flat front face.
- [x] 8.3 Unify path geometry and batch dashes, reduce rear group-frame emphasis, and size overview panels to their content.
- [x] 8.4 Inspect light and dark outputs and record the delivery validation scope; delivery receipts own final gates and publication.

## 9. Adapter correctness

- [x] 9.1 Extend browser acceptance before implementation for activation timing and nesting, diamond text containment, decoded labels and generic class titles, invisible links, and local subgraph direction.
- [x] 9.2 Normalize display text at the parser boundary, preserve layout-only constraints and local directions, and fit decision text inside diamond faces.
- [x] 9.3 Anchor activations to message events and preserve visible nested intervals.
- [x] 9.4 Inspect corrected fixtures and maintained examples and record the delivery validation scope; delivery receipts own final gates and publication.

## 10. Element and layout boundaries

- [x] 10.1 Extend browser coverage before implementation for shared node/connection identifiers, self-loop projection, and long multiline group titles on phones.
- [x] 10.2 Preserve element kind in selection and picking, route self-loops outside their node faces, and include measured group titles in layout and framing.
- [x] 10.3 Inspect corrected views and record the validation scope; delivery receipts own rebasing, final gates, and publication.

## 11. Implementation structure

- [x] 11.1 Separate command orchestration, internal model/routing code, and browser scene/label/controller modules; retain one offline packaging pipeline.
- [x] 11.2 Make rendering explicit after viewer updates and derive browser target wiring from one fixture and input list.
- [x] 11.3 Verify the existing browser behavior and compare rendered outputs with the prior candidate; update ownership documentation and delivery scope.

## 12. Static snapshots

- [x] 12.1 Retain the site configuration, Sass palette, and diagram font inputs independently of the shared-theme extraction.
- [x] 12.2 Add hermetic PNG/WebP snapshots from generated 3D HTML with selectable appearance, framing, and camera view; verify behavior end to end before implementation.
- [x] 12.3 Build sample snapshots and use a representative snapshot in the project README and landing, with the interactive view still available.
- [x] 12.4 Inspect snapshots and both themes, update documentation, and record delivery scope; receipts own final checks and publication.

## 13. Independent delivery

- [x] 13.1 Rebase onto current master, preserve generated registries and application theme consumers, and extract shared identity theming into its own trunk-based PR.
- [x] 13.2 Define renderer, snapshot, website integration, quality, and semantic lint checks for the split candidates; receipts own actual check results and publication state.

## 14. Interactive camera depth

- [x] 14.1 Add a production-pipeline 50-node flowchart fixture and demonstrate the blank interactive canvas before changing the viewer.
- [x] 14.2 Fit clipping planes to scene geometry and labels for interactive camera changes, sharing the calculation with snapshots.
- [x] 14.3 Verify visible canvas pixels after opening, orbiting, resetting, changing view and appearance, and resizing; inspect artifacts and record delivery scope. Receipts own final gates and publication.

## Context

The approved prototype is specified in the prior gallery change's
`assessment.md`. The site preset owns appearance and the Sass palette owns
colors. Native SVGs embed their fonts and are delivered without rewriting.

## Goals / Non-Goals

Apply the approved combination across current site consumers. Preserve shape,
marker, and dash semantics. Broader visual alternatives remain proposals and
must not change the production preset during exploration.

## Decisions

- Set normal graph strokes to 1.75px with type-specific selectors scoped to
  normal edges. Preserve thick and invisible edges and existing marker sizes.
  Explicitly exclude thick and invisible classes: Mermaid also marks thick
  paths with the normal class.
- Change the existing caption border token in place; retain its margins and
  fill. Append only the remaining approved CSS to the preset.
- Add an actor-line palette parameter sourced from existing gray-600/gray-500
  tokens. Sequence SVGs do not have the flowchart's `.root` variable scope.
- Use explicit site colors for mindmap nodes and branches. Uniform branch
  weight is the approved trade-off for readable dark output.
- Extend browser-based render verification before implementation. Consume the
  gallery render targets and a semantic-edge fixture; save screenshots and
  computed geometry/contrast reports. Reuse existing browser dependencies.
- Render design alternatives from the updated preset under task scratch,
  varying color, typography, geometry, and line treatment. Use only pinned
  fonts and the native rendering pipeline, never hand-edited SVGs.

## Failure cases to cover

- Normal graph strokes remain subpixel at the gallery's 728px display width.
- A broad selector flattens thick/invisible weights or changes dash semantics.
- Sequence guides disappear or compete with message arrows.
- Mindmap branch or root-label colors fall back to low-contrast derived values.
- Light/dark variants lose geometry, embedded fonts, or semantic markers.
- A newly styled shared image breaks existing site/mobile/print behavior.

## Risks / Trade-offs

- Larger path weights do not enlarge fixed markers: inspect class/ER symbols.
- Very wide diagrams still need native-size viewing: do not promise that a
  theme alone makes dense labels legible on a phone.
- Appearance prototypes can require diagram-family-specific rules: render
  every proposal across all six grammars before recommending adoption.

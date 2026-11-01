## 1. Typography and palettes

- [x] 1.1 Extend existing browser font checks before implementation, then switch
      the site font and native fixtures to the pinned monospace face.
- [x] 1.2 Promote the four palettes to a named catalog, preserve their values,
      and retain the active neutral palette.
- [x] 1.3 Render all sample families in both modes; preserve article-width
      readability, semantic edges, and caption containment.

## 2. Documentation and delivery

- [x] 2.1 Update current consumer documentation and record render evidence.
- [x] 2.2 Format, validate OpenSpec, and prepare scoped delivery inputs;
      receipts own final candidate checks and publication state.

## 3. Generated article asset follow-up

- [x] 3.1 Regenerate the Markdown-to-X article pipeline WebP through its owning updater after the typography change; inspect the image and pass its freshness test.

## 4. Historical consumer wiring

- [x] 4.1 Bind the diagram-migration blog's existing T3 Code image to the explicit historical render target when Host Bot's normal targets switch to the current site style; verify the existing image freshness test without changing publication assets.

The consumer now selects Host Bot's `rendered_historical` target. All six
existing image freshness tests passed on 2026-11-01, including the unchanged
T3 Code publication SVG. The source diagram and blog image remain unchanged.
Evidence and exact-candidate delivery receipts are under ignored
`out/comfyui-diagram-style/` in the repository worktree.

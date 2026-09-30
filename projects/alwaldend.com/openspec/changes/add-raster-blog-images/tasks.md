## 1. Raster policy

- [x] 1.1 State in the blog authoring guidance that a post intended for syndication references only image media types a syndication target accepts, naming the accepted types
- [x] 1.2 State that a Mermaid diagram in such a post is referenced through the raster render beside `index.md`, produced by the repository's `mermaid_webp` rule rather than hand-committed, with the `.mmd` source remaining authoritative
- [x] 1.3 Confirm the guidance keeps the SVG render for a post that is not intended for syndication and requires no existing post to change

## 2. Validation

- [x] 2.1 Run `bazel_agent bazel run //tools/openspec -- validate --all --strict --no-interactive` with `OPENSPEC_PROJECT=projects/alwaldend.com` and confirm this change validates strictly
- [x] 2.2 Confirm the guidance names the render target and the accepted media types, and that the maintained Mermaid skill documents the same rule

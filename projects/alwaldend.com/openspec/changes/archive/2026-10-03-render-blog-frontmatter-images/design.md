## Context

The site currently emits social image metadata without displaying those images
in post content. The user's clarification explicitly assigns automatic image
rendering to Hugo, while the uploader continues to attach the banner.

## Decisions

Use Docsy's blog content hook and one shared image partial, reusing the existing
image renderer. Include the same images in the existing print layout. Keep
image declarations in `images`, and preserve supplied alternative text through
resource parameters. Remove only the matching opening references from the X
Article story and the diagrams post; later body illustrations remain.

The initial browser test failed before its new assertions because a CDN script
did not finish loading. An independent check of the built HTML confirmed that
the expected metadata image group is absent. Browser transport must be repaired
or configured without skipping its actual theme interactions before acceptance.

## Verification

Evidence lives under `out/hugo-blog-images/`. Test cases were updated before
template implementation. Validate desktop/mobile layout, print output, preserved
alternative text and later images, and the absence of an empty image group on
posts without `images`. Inspect the actual generated HTML and screenshot as
well as command success. Source delivery follows the existing PR workflow.

The full browser test passed after inheriting the configured HTTPS proxy for
its exact SRI-pinned CDN scripts; local site requests retain Chromium's loopback
bypass. No browser assertions were weakened and no new dependency was added.
The browser report confirms one pipeline header on its post, the diagrams header
plus five later body illustrations, both headers in print output, correct social
URLs, and a mobile pipeline width within its content column. The mobile
screenshot was inspected. Independent review found no blockers. Uploader,
conversion command, site, and skill configuration tests also passed. Strict
change validation passed. Final candidate quality/lint and publication are owned
by the delivery receipts; the earlier site redeploy authorization is retained.

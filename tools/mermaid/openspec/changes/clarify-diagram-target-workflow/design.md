## Decisions

The owner-local `.mmd` and BUILD declarations select the normal rendering path.
Reuse or add a small render target instead of reconstructing private renderer
arguments, fonts or runfiles. Match the existing consumer's preset; select dark
WebP through the rule attributes rather than a comparison theme.

Keep context diagrams at the requested actor/process/service/storage boundary.
Preserve explicit grouping, arrows and labels. Identify planned relationships
without claiming that drawing them implements runtime behavior.

Local previews use the same build path, respecting withheld publication. Verify
the actual output and provide the requested asset. Existing offline Promptfoo
validation checks staging and case loading; it is not live behavioral evaluation.

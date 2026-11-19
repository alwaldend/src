# Native numeric-transform failure receipt

The attempt opened a task-owned A157 copy whose SHA-256 remained identical to
the protected source:
`433d08ad36be488bb16e4221a85f831d4390660c258a43ea0b08775811574b73`.
No candidate geometry was saved or changed.

In multi-object Edit Mode, deterministic source-space selection produced the
planned counts:

- left red tail: 132 of 650 outer-band vertices;
- left white ruffle: 67 of 332 outer-band vertices.

Connected proportional editing, smooth falloff, and radius `0.045` were active.
`bpy.ops.transform.translate` returned `RUNNING_MODAL`, and the mapped window
reported `TRANSFORM_OT_translate` as the live modal operator. Timer-spaced
native key events sent `-0.0267` followed by Return. The modal then exited with
no error.

Post-operation edit-mesh inspection found zero changed vertices in both
objects. Selected-band X displacement, protected-root displacement, Y change,
and Z change were all exactly `0.0`. Return therefore confirmed the transform
at its default zero value: Blender's modal numeric parser did not accept the
simulated number-key stream in this context.

This is an input-parsing failure, not a failure of the deterministic selection
or proportional solver. Replaying different spellings of the same numeric key
stream is prohibited. A corrective plan may supply the measured vector through
the native transform operator's explicit `value` property, which removes this
failed parsing mechanism while preserving Blender's own proportional transform
implementation.

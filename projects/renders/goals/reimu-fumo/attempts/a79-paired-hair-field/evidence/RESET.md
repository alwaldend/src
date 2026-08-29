# A79 variant A — RESET

Stage: pre-Blender, pre-render. Blender was never run for this variant.

The frozen handoff is rejected and must not be consumed. The later scratch
interface revision can legalize the four-band quad annulus, bridge-only
vertices, and target-facing outer roots, but it cannot repair these independent
representation and validation defects:

1. Whole-pocket contour convergence reaches `0.774 px` in the front/crown
   three-quarter view, above the `0.5 px` gate. The published proof measures
   only the outer perimeter.
2. Bridge geometry is absent from the quality/projection proof. It creates a
   `+2.720 px` leaf side extremum and a `+1.090 px` rear extremum.
3. `aligned_inner_control()` rewrites every inner control point by as much as
   `3.593 mm` and rotates authored displacement by as much as `71.09°`. That is
   not the separately authored inner field claimed by the specification.
4. `A79_ROOT_RECEIVER_FACE_LOOP` is open; its endpoints are `121.39 mm` apart,
   so it cannot register to the closed 196-anchor face-opening contract.
5. Rear receiver-independence/camber and leaf-width gates are described but
   are not enforced by the frozen-path preflight.

Fixing these together requires a new inner-skin, bridge, face-root, and
whole-pocket validation architecture. Under the parent timebox, that is a
strategy reset rather than an in-place correction.

Current rejected artifact hashes:

- `geometry_preflight.py`:
  `b8b471b535e87fa8e725c32fb63747bc02dbd55fdf8c49b8566b41df772ab95e`
- `builder.py`:
  `30bb97dee54dcc73bcdd75e31733fddcd1b44254c8c85af59ce8f27fc8a4a73d`
- `control_nets.json`:
  `a660ad4dfd022eb36fe6f83e8a3c7e0e36a5ce656d6a449c250e0cb8e33d39f2`
- rejected `frozen_geometry.json`:
  `ee9708c076bd443d59b2828c727b4ae85784c8043a87aebfdb8df896c0e44812`
- `geometry_preflight_report.json` (`REVISE`):
  `1a76e455b41ba4f9a0c18412762fa417c213bef34a50d3de9a2aa7309c38aa74`

The current report still names old-interface incompatibilities that the
pending scratch contract revision may remove. Those are not the reset basis;
the five independent defects above are.

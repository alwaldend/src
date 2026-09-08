# Head034c stopped before object creation

Helper: `head_034c_draft.py`, entrypoint `build_head_034c()`, exact targets
the old fringe and hood. Evaluated SHA256:
`68a035190a68d602813e93a501ad25208b68680483bdeb4fed02e3dff20c234d`.
Source remains032 SHA256
`6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8`.

One causal derivative-only test ran. Cuts, 432-knot correspondence,
derivative magnitudes, thickness endpoints, materials and captured source
surfaces were unchanged. No replacement object was created and neither
target was hidden. No model save, render or second trial occurred.

## Measured result

All 1,728 co-normal endpoint-direction records passed: minimum dot with the
required incident-edge side is0.0221195; maximum residual longitudinal dot
is7.45e-8. The 864 scaled derivative pairs are logged before creation.

The fixed strip nevertheless failed its precreation Jacobian gate:

- 456 invalid cells of6,896:226 outer and230 inner.
- 214 have mixed corner Jacobian signs;242 are wholly reversed relative to
  their expected skin orientation. No failure came from the numerical area
  or normalized-sine lower bounds.
- 343 invalid cells are wholly above145 mm, including170 outer cells.
- Six precreation junction angles exceed90 degrees.

All862 junctions: median1.585 degrees, p95 6.844, maximum167.181.
The complete distributions and raw edges are in the JSON. Above145 mm,
fringe/bridge junctions have median1.808, p95 4.563, max75.148 degrees;
hood/bridge junctions have median1.452, p95 6.212, max167.181 degrees.
These use the first oriented triangle of each precreation polygon, not a
claim that a final Blender object was tessellated and audited.

Concrete mixed-Jacobian witness, outer array face74697:

```text
Vertices: 8599, 40806, 42046, 42045
XYZ, mm:
(-55.982504, -33.874404, 145.570755)
(-55.982504, -33.862129, 145.584062)
(-56.193467, -33.201665, 145.360664)
(-56.162704, -33.005167, 145.775855)
Corner orientation cosines:
+0.999534, +0.999815, -0.999861, -0.999853
```

This sign reversal within one cell cannot be dismissed as a benign
scanline crossing. The corrected endpoint signs do not establish a usable
interior mapping across the fixed stepped cuts.

## Evidence and stop

`head_034c_diagnostic.py/.json` preserve vectors, signs, all cell points and
Jacobians, precreation junction coordinates/normals, the assertion and
source/helper hashes. A literal trailing `\\n` serialization terminator in
the report was normalized after execution; data were not changed, and the
diagnostic source still records the formatting defect. There was no model
replay for that report correction.

Actual strip-intersection and evaluated hit audits did not run: the guard
correctly stopped before an object existed. Their absence is explicit, not
a negative collision result. No geometric or visual suitability is claimed.

Root rejects this fixed ruled-strip family. The next plan-only question is
whether the exact arcs bound a simple domain in a surface chart, permitting
constrained triangulation and a smooth single-valued depth field without
forced knot-to-knot cross-strip correspondence. No dose, width, cut-value
or resolution tweak is proposed for034c.

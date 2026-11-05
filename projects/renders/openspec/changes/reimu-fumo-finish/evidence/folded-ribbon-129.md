# Folded ribbon fixture 129

Study 128 is rejected: it changed the garment envelope but did not establish
sewn construction or gathered fabric. A separate ribbon fixture now tests
literal material fullness before another full garment replacement. This is
not a model candidate or a cloth-solver calibration.

## Bounded decision

Build one 9 mm-wide strip with three explicit omega-shaped return folds.
Use approximately 16 mm projected pitch, 30–34 mm material length per cup,
8 mm cup depth and 48 mm total attachment span. The material path doubles
back at the fold mouth. Extruding the path across the strip width retains a
developable surface; it does not rely on a sinusoidal free-edge border.

Measure actual edge length, compression ratio, area, strain, return spans
and clearance. Render front, side and elevated views at the real 25 cm asset
scale. Reject a corrugated fence, intersections or closed cups with no
readable interiors. One fixed construction is allowed, with no shape sweep.
Sewing its attachment edge and orienting it around a curved hem remain
separate work; neither is implied by fixture success. The seated-cloth
solver still requires the calibration recorded in study 118.

Root remains the sole native writer on the same pinned four-thread Blender
route. Save a new standalone fixture under ignored `assembly_129/`, then
clean-reopen its bytes for inspection. Do not modify any model source or
integrate the fixture before its pixel review. No final criterion advances.

## Result

The fixture saved as
`c72bc2e531a0bbceaa9666396ae4eb019a0e78914f8a14c1e325055676c1d229`.
Protected source 126 remains unchanged. The first wrapper invocation stopped
before loading a model because its copied output path still named 127. The
corrected wrapper derives its directory from its own file; both the failed
wrapper and log remain preserved. This setup correction changed no geometry.

The single construction contains 96.918 mm of material over a 46.639 mm
projected attachment span, a 2.078:1 ratio. The rigid-panel fold construction
preserves edge/diagonal lengths within 1.46e-13 relative error. Native checks
found zero nonadjacent midsurface intersections and 1.638 mm minimum sampled
ruling clearance after the thickness allowance. Those witnesses do not prove
continuous-surface clearance, sewing, equilibrium or visual fidelity.

Root inspected all four clean-reopened views. Independent pixel-only reviewer
`/root/ribbon129_blind`, given the reference and fixture images without source,
metrics or prior models, **rejects the fixture**: construction 5/10, fabric
read 3/10 and contact/self-overlap 4/10. The nearly uniform walls and matching
upper/lower contours look extruded; the front becomes a band with rectangular
tabs, and the return folds read as flanges. Open interiors alone are not soft
gathered cups. No geometry is transferred to the model.

This fixture family is closed. Increasing the crease skew or fold amplitude
would repeat its failure. The next ruffle method must establish a narrow sewn
attachment, a flared free rim and compressed three-dimensional fabric folds
together. If it uses simulation, first satisfy study 118's small-scale cloth
calibration requirement; this static fixture does not meet that prerequisite.

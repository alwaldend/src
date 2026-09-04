# 099 arms saved after first build, visual review pending

arm_099_candidate.blend925350ae56377cea29f22fbf55045f26e6df1b5040ca38e8efd4d4526ed1290a.
Source0982bee72c54ba571662e36fead903c9516ecb0f13b04d9138fbe897e0c1ffb6193
unchanged; only two arm meshes changed. No body/cloth/hair/material edits.
First native execution exited0, no data repair used, one remains available
only for causal API/data issues (not post-render shape changes).

Both arms1106v/2208tri, closed consistent positive one-component, zero
nonadjacent self/duplicate/degenerate triangles. Volumes9.0240483/9.0241187e-6m3.
All actual visible-mesh receiver crossings0 except intended body attachment.
All positive sampled section centers passed two actual closed-wall winding
loops; centerline chords had0sleevewallhits.117vertices fitted per arm,
minimum radial scale0.87065192, max12distance steps, minimum sampled gap
0.80020458/0.80020929mm. This is not a global0.8mm triangle-gap certificate.

New body contact pairs350left/324right;106crossing armtriangles per side form
one band. Removing it leaves727buried and1375exterior triangles, strict3ray
component witnesses unambiguous. Original proximal polesinside, useful095
distal targetsoutside. Contact locus changed as explicitly permitted.

Builder26d152a5112c730600c9ab8cd794a5b017d4368f3878e66d97a8bd062912bdad,
arrays151523770feaafef4c7ff1650180b8fc5ea9e7a6b51ece4c671f38529c336806,
preflightd1fe957f8936e8dea59ad18610c5c916b354866872e084b11dd2b428c1f2bce4.
Clean-reopen fullfive/presentation now rendering. No retention or stage pass
until root and implementation-blind review. Retained098 remains fallback.

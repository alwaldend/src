# Shingled-lock reset evidence

Blender MCP 5.1.1 created four independently closed curved crown meshes and
two closed side-lock meshes from protected A157, whose digest remained
`433d08ad36be488bb16e4221a85f831d4390660c258a43ea0b08775811574b73`.
Candidate digest:
`3847634874fecdccb50795d7519c532449338c711c762c2f0350a580e1ba713f`.

Pinned Blender 5.2.1 clean-reopened that exact candidate and rendered the
frozen fast trio. Candidate hashes before and after rendering match. Outputs:

- front: `cb34341d40911b4ae952d010d9b6aeeb076b314c8c29e7c191b62e01a7973874`;
- side: `767cd286871a8e342981922182cc760767a62b21f1594558eb9a2b5e5386a76b`;
- three-quarter:
  `77eacb884a3b50803ce0151377a2098eecb5182bc384d7c8d8f003a7442bb23f`.

All three views fail the predefined stop gate. The front exposes most of the
beige crown. Side and three-quarter show the new side lock as a large detached
oval pad, with sharp crossings at its root. Only a thin strip of the crown
cushions is visible because their analytic surface approximation sits inside
the A157 head after its lattice deformation.

This is a source-to-evaluated-geometry placement failure, not evidence that the
closed-pillow data structure worked. Measurements and full views are therefore
unauthorized and uninformative for this candidate.

# A100 topology and edit receipt

- Parent and working copy initially matched SHA-256
  `c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`.
- Blender: repository-pinned `5.2.1 LTS`, interactive MCP host.
- Exact owner: `Garment42 rear pooled dress panel`.
- Owner mesh: 325 vertices, 612 edges, 288 polygons, identity object
  transform, single-user mesh data.
- Topology is a 13-by-25 quad grid. Pinned waist boundary is `0..24` at
  world `Y=0.021 m`, `Z=0.067 m`.
- Free rear boundary is explicit indices `300..324`, initially all at world
  `Y=0.049 m`; first diagnostic moved only those key coordinates to
  `Y=0.0395 m`, preserving their Z coordinates.
- One shape key was added: `A100_RearHemSupport` with value `1.0`; the Basis
  remains unchanged.
- Attached `Garment42 rear pooled ruffle` received object translation
  `Y=-0.0102 m`; attached rear stitches `00..08` received object translation
  `Y=-0.0095 m`.
- The unauthorized diagnostic correction additionally moved only adjacent
  loop indices `275..299` by `-0.005 m` in the same shape key.
- No topology, modifier, material, seat, foot, leg, front/side panel, or
  upper-body change was made.


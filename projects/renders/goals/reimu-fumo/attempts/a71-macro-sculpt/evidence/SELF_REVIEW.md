# A71 P0 sculpt-control review

## Verdict

Reset at P0. Do not proceed to front-silhouette sculpting. The fresh one-remesh
volume is valid, but the synthetic native brush interface does not supply a
reliably broad macro deformation in this live Blender configuration.

## Evidence

- Fresh `A71_MacroClay`: 3,863 vertices / 3,861 polygons after exactly one
  4 mm voxel remesh; dimensions about 130.9 x 101.3 x 129.1 mm.
- Raw P0 blend:
  `sha256:0004a3f0bc4987a250f7028b8697a7f740c70866ad8b60f7181e2b4eafa96400`.
- Ordinary Grab, 50 mm world radius and 0.40 strength, changed only 11
  vertices (`0.285%`) in a roughly 10 mm support patch; maximum displacement
  was 1.80 mm.
- Elastic Grab, 50 mm world radius and 0.30 strength, produces numeric changes
  on all vertices, but only 40 vertices move more than 0.1 mm and only six move
  more than 0.5 mm. Its maximum displacement is 0.814 mm and mean displacement
  only 0.0045 mm.
- Both operators poll and finish successfully. Setup and stroke were separated
  across MCP calls so the VIEW_3D context could settle. The failure is support,
  not an unexecuted operator.
- Protected rung 003 remains byte-exact at
  `sha256:c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`.

## Decision

A68 already showed that reported screen radius is not evidence of useful
support. A71 repeats that failure on a fresh uniformly remeshed object, so it
cannot be blamed on the old cap topology. Increasing strength would deepen the
same local dent; repeated strokes would violate the P0 gate.

The next macro-form attempt should use a deterministic low-frequency control
cage or analytic rest-surface deformation whose affected vertex set and
world-space displacement are explicit before rendering. Native sculpt brushes
may return later for local organic breakup only after the macro silhouette is
approved.

## Process audit

The Flatpak MCP loop remains fast and useful, but interactive state must settle
in one MCP call before a native stroke is sent in the next. That is now captured
in reusable scripts. The process correctly spent no time on artistic strokes,
materials, or retopology after the P0 support veto. Future support gates must
count materially moved vertices (for example greater than 0.1 mm), not every
floating-point delta from an elastic falloff.

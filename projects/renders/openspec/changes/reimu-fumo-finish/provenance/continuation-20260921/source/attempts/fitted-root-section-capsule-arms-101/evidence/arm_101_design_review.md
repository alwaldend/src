# Arm 101 fitted root-section review

Observed 2026-09-06 at approximately 04:39 UTC. Recommended verdict:
**proceed with this bounded root-section reconstruction**. It directly
addresses the measured interpolation failure without changing the useful
100 hand/cuff shape. The coordinator owns the verdict and implementation;
one added section is not guaranteed to pass finite-surface clearance.

I read the complete 100 collision diagnosis, failed-result, plan, and image
direction reports, and inspected its front, side, and mirrored oblique
diagnostics. The direction review supports preserving the rounded high hand
and front containment. The geometry diagnosis isolates all 62 sleeve pairs
per side to row 0, where 14 connecting edges cross between clear fitted
spokes on opposite sides of A=0. The failure is therefore neither the cuff
nor a defective radial-path check.

Adding fitted samples at the actual root is a causal change to the skin
interpolation. Merely splitting the existing world triangles would retain
their intersections. The proposed new points are instead evaluated through
the 100 field and fitted at the root, so they intentionally need not lie on
the old straight world-space edges. This locally changes the connecting
surface while preserving the 1,106 existing points and distal cap.

The following details make the construction well-defined:

- Use the exact hash-bound 100 triangle connectivity and source vertex
  identities. Leave all uncut triangles unchanged. Each child inherits its
  original triangle's material and winding; do not silently retriangulate
  unaffected source quads by a different rule.
- Distinguish the global root parameter from each edge's clipping fraction.
  Compute the root's normalized longitudinal coordinate using the global
  proximal/distal axial bounds. On a crossing edge, interpolate to that
  coordinate using that edge's endpoint parameters. Substituting edge bounds
  into the global normalization formula would produce a different parameter.
- Key each shared intersection by the unordered original edge IDs and reuse
  that one vertex in both incident triangles. Interpolate angular XY before
  normalizing its direction; do not linearly interpolate wrapped angles
  across the angular seam. Use the exact nominal root as the mapped center,
  then apply the unchanged capsule profile and conservative full-sleeve fit.
- Preserve each clipped polygon's cyclic order when triangulating its
  triangle or quad. The refitted world polygon may be nonplanar, so the
  chosen finite triangles, not the parameter polygon alone, need the final
  winding, degeneracy, self-intersection, and receiver checks. Shared-edge
  reuse prevents cracks only if every incident crossing triangle is split.
- Treat 96 added vertices and 192 added triangles as a consistency
  expectation derived from the actual crossed source band, not counts to
  force. Assert actual shared-edge incidence, closedness, original-point
  equality, and unchanged uncut connectivity/materials. Report discrepancies.

The proposed exact preservation check must compare actual original 100
fitted world points, not merely unchanged input parameters. The existing
root frame, capsule endpoints/profile, and all nontarget geometry remain
fixed by the proposal. Adding the critical section can still produce new
root-loop edges or joining faces that cross between safe samples; rerun all
finite topology, self, body-band/containment, and receiver gates without a
root-contact exception. Changed triangle counts alone are not collision
severity or a reason to waive a geometric crossing.

No extra radius, endpoint, or sampling grid is justified by this review.
If the new finite surface still fails, the exact residual crossing geometry
must determine the next decision. If it passes, fixed-view review must
confirm the retained 100 shape direction; preserving original vertices does
not alone certify identical shading or local attachment appearance.

## Evidence and scope

| Read input | SHA-256 |
| --- | --- |
| `arm_100_collision_diagnosis.md` | `ba8b7b3676b3ac1a7f718e3d520e3cc416c6c901b0fc350fd3635ac912f69de8` |
| `arm_100_pixel_direction.md` | `6ec1e124ab75ff50fe3e36ff26a14e85d0138bb59bb93097b1847bcfec37a845` |
| `arm_100_failed.md` | `26b1d523c6df5542ea1f853d4b111e1d0b7657b80c3ced2baf596196441c08a5` |
| `arm_100_plan.md` | `1689815a8c1e5e13573334e10c9964d9d7ae2209fb481e00fecc49877d114790` |

Paths are relative to this ignored report; the reports bind the inspected
images and failed arrays. The detailed 101 proposal came from the
coordinator's task message. No native run, scalar geometry computation,
model/goal edit, or Git mutation was performed. Only this report was written.
The linked worktree remains on `t3code/continue-fumo-desktop-use`, HEAD
`c7601f0fc80e0a94b585910459639ffccbdbdbd4`. Decision-review and reference-fidelity
informed the distinction between this supported reconstruction hypothesis,
its mandatory finite-surface checks, and visual acceptance.

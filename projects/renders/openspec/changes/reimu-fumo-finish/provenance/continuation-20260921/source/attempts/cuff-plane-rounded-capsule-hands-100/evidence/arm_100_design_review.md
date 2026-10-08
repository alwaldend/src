# Arm 100 capsule-at-cuff design review

Observed 2026-09-06 04:20 UTC. Recommended verdict: **proceed with this
bounded capsule hypothesis**, with the implementation distinctions below
explicit before freeze. The coordinator owns the verdict and geometry.
This report proposes no alternate dimensions or position/radius grid.

## Causal assessment

The 099 native result establishes a feasible continuous arm route, clean
cloth/receiver contacts, and a meaningful buried body connection. Its pixel
review nevertheless rejects the placement: the hand is a small crescent at
the upper rim in side and largely disappears in oblique views. I read that
review, the 099 plan, and saved-result report in full, and inspected the
099 front, side, and mirrored three-quarter images. I agree that the empty
opening is the dominant local failure. This is an informed design critique,
not an implementation-blind visual review or acceptance of 099.

The proposed 100 change addresses that failure directly. A tapered sphere
ending before the actual lip leaves little terminal volume visible. A
rounded capsule cap ending at the measured mean outgoing-cuff plane places
substantial volume immediately behind the opening. Using the actual cuff
as the terminal anchor supplies a construction reason for the extension;
retaining the sleeve/root and verified fitting method avoids reopening the
earlier attachment and transition failures without cause.

The supplied mean plane is a summary of a curved lip, not a physical flat
boundary or a reference-measured fingertip target. The reported lip axial
range 39.066–46.263 mm means different parts of the rim will expose or
occlude the same cap differently. Ending at the mean is a credible single
hypothesis to test, not a guarantee that the hand will be visible enough.

Nor does preserving the 095 C/Y offsets guarantee front containment.
Increasing the A coordinate moves the endpoint both outward and downward
in world space. The thicker capsule also has a lower/outboard silhouette
away from its pole. Thus 100 can fix the hidden hand and still recreate a
low frontal protrusion. The existing front gate must reject that result,
even if the capsule clears every cloth triangle. Side and both oblique
views must simultaneously show a substantial rounded hand high within the
opening, rather than only a longer peg or a cap pressed against the rim.

## Geometry details to make explicit

1. **Construct a true capsule profile.** The constant middle transverse
   radii and the specified half-ellipsoid end caps have matching value and
   zero axial radius derivative at their junctions, so their intended raw
   surface can be C1. Apply that profile to the normalized source local-XY
   direction; do not multiply by the sphere's latitude radius again, which
   would restore the taper this change is meant to remove. At the two source
   poles, local XY has zero length: assign the exact centerline endpoint
   directly rather than normalizing it. Reuse the existing oriented topology
   and right-handed basis. A C1 raw profile does not certify the final fitted
   mesh after independent radial contractions.
2. **Keep fitting conservative.** The inherited distance-bound march still
   queries the actual entire sleeve, including both sides of the root and
   the cap/opening region. The minimum fitted scale 0.6 remains a rejection
   threshold. Never expand an unsafe fit up to that floor. Preserve the
   conservative termination rules and reject unresolved centers or fits.
   The 0.8 mm margin continues to describe sampled radial paths, with no
   new claim of a global finite-triangle gap.
3. **Resolve the cuff cutoff's exact index set.** Compute the proposed
   opening threshold from both evaluated wall halves of outgoing row 40,
   return rows 41–44, and the **distal cuff rim**. Exclude the proximal/root
   rim from a generic end-rim collection; including it would move the minimum
   near the root and incorrectly exempt most of the sleeve from the section
   test. The stated 0.1 mm separation remains the coordinator's parameter,
   not a value selected by this review.
4. **Keep the closed-section and opening claims distinct.** Before that
   measured cutoff, retain the two actual nested contour/winding test for
   sampled centers. Do not reinterpret a failed contour there as an opening
   merely because a looser gate would pass. In the finite opening region,
   the wavy hem can make two closed nested loops geometrically inapplicable;
   labeling that region as an opening is appropriate. It does not prove
   fully enclosed lumen membership there. Continue checking every spine
   chord and the entire final arm against every actual cloth/receiver
   surface, including faces spanning the cutoff. No old sleeve-contact
   exception or untested hand cap follows from the opening label.
5. **Recheck the full attachment and finite mesh.** A constant-girth middle
   changes the proximal envelope as well as the hand. The inherited
   connected body-crossing band, one buried/one exterior component, strict
   containment witnesses, pole tests, closedness, positive consistent
   volume, self/duplicate/degenerate checks, and full receiver contacts must
   apply to the new mesh. The prior 099 pass cannot certify the changed
   capsule. The prescribed no-post-fit-smoothing rule prevents a later
   operation from invalidating those checks.

These points do not require another dimension choice. They clarify the
proposed profile, the exact surface set defining the opening zone, and the
limits of the available geometric certificates.

## Alternatives and decision boundary

Keeping 099 preserves a clean geometric result but retains its rejected
empty-sleeve appearance. Translating the same thin tapered endpoint alone
would be less directly responsive to the lack of visible rounded volume.
Moving the sleeve or its root would disturb the reference-facing garment
without evidence that its placement causes this newly isolated failure.
The proposed capsule therefore has the strongest causal basis among these
bounded options.

Its main risks are a low/front protrusion, an overlong or rigid-looking hand,
or fitting contractions that erase the intended cap volume. Treat those as
possible outcomes to judge in the fixed pixels, not as reasons to adjust
dimensions repeatedly before an evaluated candidate exists. If the one
candidate passes geometry but fails these pixels, record which part of its
actual visible envelope causes the failure before choosing another method.

The expected visual result is a contained frontal silhouette together with
a visibly rounded hand occupying the upper opening from side and both
obliques. Passing native checks alone is insufficient, as 099 demonstrates.
Neither this design review nor a local arm pass clears the whole asset.

## Evidence and scope

The 100 proposal, cuff axial measurements, and gate changes were supplied by
the coordinator in the task message. They were not independently recomputed
here. Source remains 098 per that proposal; 099 is failed comparison evidence,
not a retained input. The earlier 099 plan binds the source model and exact
receiver authority.

| Inspected input | SHA-256 |
| --- | --- |
| `arm_099_pixel_review.md` | `326f25724b46691b6614ab86dde49fcfa7036587d89b52b158411c7c02c388c9` |
| `arm_099_plan.md` | `55b9078a00d134e7c6baf5ce095c76325c786bb9b0632cd97283cf65838ea785` |
| `arm_099_saved.md` | `b0eea2306fa031ff76705865ac8e997e8252d807fbb2aadb7d001b45e9ff4c49` |
| `arm_099_all_review/front.png` | `8e93a5214056bf1dfcd5a31c64cc6dd69a53f70c950d9dbf8fe3f3304132c594` |
| `arm_099_all_review/side.png` | `db69d9f6f240bb2a5f242a4a2902013682cc2b8f30dca4b55252fd341aacd69b` |
| `arm_099_all_review/three_quarter_mirror.png` | `f558848dee78969602a7fab2a9cf8755e9cc4ba31cac90967a85af66688dd936` |

Paths are relative to this ignored report. No native/Blender process,
geometry computation, parameter sweep, model edit, goal write, or Git mutation
was performed. Only `arm_100_design_review.md` was written. The previously
verified linked worktree remains on `t3code/continue-fumo-desktop-use` at
`c7601f0fc80e0a94b585910459639ffccbdbdbd4`. Decision-review motivated the
causal comparison; reference-fidelity keeps geometric success separate from
the mandatory front/side/oblique visual result.

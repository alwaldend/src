# Sleeve 095 correction to assembly review

Observed 2026-09-06 02:30 UTC. This new report corrects the retained-root
attribution and resulting recommendation in
[`sleeve_094_assembly_review.md`](sleeve_094_assembly_review.md). The earlier
report remains unchanged as historical evidence. The correction is based on
the completed [`sleeve_095_actual_assembly.md`](sleeve_095_actual_assembly.md)
audit, read in full, and the reference/candidate images inspected for 094.
This reviewer performed no Blender run, geometry edit, new calculation, or
candidate implementation.

## Factual correction

The 094 review incorrectly called the |X|=38 mm, Z=56.829788 mm throat
"current" and treated that lower/outboard root as part of retained 087/094.
It belonged to rejected 088–092 sleeve constructions and was never retained.
The resulting claim that retaining the current sleeve root is a practical
dead end is withdrawn. The advice to relocate that root does not follow for
the actual baseline.

The 095 audit opened retained `hair_094_candidate.blend` and verified that
both sleeve base meshes match the 067 loft and 068 cuff return with zero
coordinate residual. Their actual cloth-center-surface roots are centered
nominally at **(±24, −3, 70) mm**, with bilateral bounds |X|=19.423–28.577 mm,
Y=−13–7 mm, and Z=63.439–76.561 mm. Those are source cloth-center-surface
measurements, not evaluated Solidify wall rings. They replace the rejected
38 mm root as the baseline for this decision.

The actual outgoing cuff row 40 has center approximately
**(±59.231, −2.391, 45.347) mm** and bounds |X|=47.662–73.248 mm,
Y=−20.702–16 mm, Z=28.148–61.734 mm. Rows 41–44 form the existing cuff return.
The retained arm's actual distal principal-axis endpoints are
**(±68.978, −13.526, 39.693) mm**. The audit confirms that both evaluated
arms are unchanged from the source 087 evidence.

The 094 scalar numbers—38.927% section height and 11.226184 mm extension—
remain conditional findings for the **provisional analytic corrected cuff**
used in that report. They are not measurements against the actual retained
067/068 cuff. Do not relabel them as retained-sleeve clearance or protrusion
distances. Their useful conclusion is that the old arm can remain low and
distally extended even after an assumed cuff admits its planar section.
The visible 087 hand protrusions remain direct image evidence; their exact
relationship to the retained sleeve needs that sleeve's actual surfaces.

## Revised bounded recommendation

Recommended verdict: **proceed with one arm-only construction hypothesis
against the retained sleeve**, subject to the existing visual and geometry
gates. The coordinator owns the decision. This is a credible smaller change
than relocating or rebuilding the sleeve, not an acceptance finding.

The corrected root is already in the shoulder region implicated by the
reference's visible attachment witnesses. Those witnesses are occluded and
unregistered, so they cannot prove exact root fidelity, but they no longer
support the specific lower/outboard accusation in 094. The retained sleeve
also already extends substantially below the current distal arm. These facts
make a shorter, higher stuffed arm inside the existing cloth a coherent
first hypothesis without requiring a new root or lower cuff.

Keep the retained sleeve and its attachment fixed for that test. Derive the
shorter/higher arm from the actual retained opening and evaluated interior,
while keeping its proximal connection continuous with the body. Retraction
and elevation should reduce the conspicuous frontal hand exposure and leave
the reference's empty lower sleeve pocket visible from the side. Do not
transfer the rejected root coordinates, the provisional analytic cuff plane,
or an invented hand-height fraction into the new geometry.

The main tradeoff is that an arm-only change cannot alter the retained
aperture's silhouette or rim construction. If matched front, side, and
three-quarter pixels still show the wrong sleeve profile after the hand is
correctly placed, that would support a separate cloth correction. The 095
audit does not establish that such a correction is already necessary. It also
does not certify current sleeve collision clearance, so the arm must be
checked against the evaluated walls, cuff return, body, and nearby garment;
an outer bounding box or a plane-section test alone is insufficient.

Thus the durable distinction is: rejected sleeve trials had the problematic
38 mm root; retained 087/094 has a verified 24 mm root. The unchanged arm is
still a credible source of the exposed-hand mismatch, while retained-sleeve
relocation is currently unsupported. Camera registration remains relevant to
precise image comparisons and does not change this provenance correction.

## Input identity and scope

- `sleeve_095_actual_assembly.md` SHA-256:
  `d1d6bcc7f85bd02512f5e7563fa910c7f842c7c482050604210893f81a75a520`.
- Unmodified `sleeve_094_assembly_review.md` SHA-256:
  `e3e6729f2094eaf7ec9a964823716ba98b9cca739d84005b1f248883a6edabf1`.
- Audited retained model hash, as recorded by 095:
  `af5a61921ae309a69cfe98b0da092d206b2d0ff29c9123ac2a4886de5a0f9add`.
- Audited geometry JSON hash, as recorded by 095:
  `4a23e0161fa80d992aedce2598ac7943b887bf83c2d5f6b1f3f10cd74783c3bf`.

The audit report is the source of the geometry identities and measurements;
this reviewer did not reread the model or the 48 MB JSON. The prior 094 report
records the inspected image hashes and reviewer-context limitation. The
branch remains `t3code/continue-fumo-desktop-use` at
`c7601f0fc80e0a94b585910459639ffccbdbdbd4` in the previously verified linked
worktree. Only this new ignored scratch report was written. No acceptance,
clearance exception, goal mutation, or model publication is authorized by
this correction. Decision-review's requirement to reopen a conclusion when
its premise fails directly caused this revision.

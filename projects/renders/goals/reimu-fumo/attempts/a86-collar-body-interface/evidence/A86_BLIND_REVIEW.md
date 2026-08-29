# A86 collar/body interface implementation-blind review

## Review isolation

I inspected only the S00, S01, and S02 fixed front/profile beauty pixels,
their front ID pixels, and the supplied physical front, side, and canonical
turntable references. I did not inspect Blend geometry, topology, scripts,
object names, manifests, or the implementation intent.

This is a representation test: can a collar alone hide the known S12
underside while still reading as the compact constructed collar in the
physical plush?

## Absolute answer

**No. Collar-only coverage is not a viable representation.** None of the
three states hides the S12 underside in profile, and increasing the same
horizontal collar band enough to do so would make the already prominent yoke
wider/deeper and less faithful to the references.

The physical references do not show a continuous structural ring carrying the
entire head width or depth. From the front, only a compact central collar and
tie region is visible beneath the face, while hair locks and the dress/sleeve
mass control the lateral transitions. From the side and rear, layered hair and
garment parts cover most of the head-to-body junction. A single collar band is
therefore being asked to own pixels that belong to several constructed parts.

## State-by-state pixel verdict

### S00

- The collar is an obvious rounded horizontal spacer between two disconnected
  blocks.
- Its side lobes create a dumbbell or neck-pillow silhouette in front.
- In profile, the dark underside cavity remains fully visible above it.
- The ID image confirms that the apparent light region is a single continuous
  collar owner, not several reference-faithful overlapping panels.

**Verdict: reject.** It neither covers the defect nor matches collar logic.

### S01

- Raising or enlarging the collar makes the front junction more continuous,
  but the white band becomes a harder shelf across nearly the torso width.
- The profile still shows the dark head-underside trough as a separate line
  above the collar.
- The ID image again exposes a continuous broad yoke, not the compact split
  collar visible in the physical plush.

**Verdict: reject.** More reach produces more yoke without achieving coverage.

### S02

- S02 is the **best diagnostic state**: the collar sits most tightly under the
  head, and frontal overlap suppresses more of the disconnected-spacer read.
- It is still not a viable candidate. The profile cavity remains plainly
  visible, and the front ID pass shows that the clean beauty overlap is still
  supplied by one broad horizontal collar band.
- Any further width/depth growth of this owner would approach a rigid head
  support or shoulder yoke, while the reference collar remains compact and
  partly hidden by hair, tie, and dress construction.

**Verdict: reject as an asset state; retain only as the strongest negative
representation test.**

## Best-state decision

| State | Relative result | Absolute result |
| --- | --- | --- |
| S00 | Weakest; disconnected dumbbell spacer | Reject |
| S01 | Tighter, but harder and broader yoke | Reject |
| S02 | Best coverage probe | Reject |

No state is safe to keep as the provisional collar/body interface. S02 is the
best evidence because it demonstrates that the collar-only approach reaches
its visual limit before hiding the underside.

## Exact representation verdict

**`COLLAR_ONLY_BAND: REJECT`**

Replace the single-owner coverage hypothesis with a **composite constructed
interface**:

1. a compact split front collar/neck opening owns only the small central white
   collar pixels visible beneath the face;
2. front and side hair panels own the lateral head-to-shoulder coverage;
3. rear hair panels and the upper dress/torso own the rear junction; and
4. any underside pixels still exposed between those owners require a bounded
   head-interface repair rather than further collar growth.

This verdict does not require splitting the reusable head file or rebuilding
the whole model. It changes only ownership of the interface pixels. The next
coupon should test the compact collar together with minimal proxy hair-side
and rear panels; a collar-only fourth size variant would repeat a representation
that the three reviewed states have already falsified.

## Acceptance condition for the replacement

The composite interface is viable only if beauty and ID views jointly show:

- no exposed dark S12 underside in front, profile, or three-quarter views;
- no continuous white yoke spanning the head or torso width;
- a compact central collar consistent with the physical front reference;
- layered lateral/rear coverage assigned to hair or garment panels; and
- no clipping, floating gap, or accidental tangent between the owners.

Until that test passes, S12 remains only a provisional receiver and the collar
interface remains unaccepted.

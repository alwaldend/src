# A69 C1 v4 diagnostic review

## Absolute verdict

**Reject C1 v4 for promotion and do not advance this candidate to C2.** The
bounded deformation produces a real but modest reduction of the rear-depth
silhouette in the side and both three-quarter cameras. It does not solve the
dominant reference mismatch: the head still reads as one deep, rigid
helmet/block rather than a shallow stuffed receiver plus a separate thin rear
hair layer.

No obvious new bald patch, skin breakthrough, local dimple, hard daylight gap,
or front face-aperture regression is visible in the 320 px diagnostic renders.
That does not justify overriding `promotion_contact_gate=false`. Exact
clean-reopen measurements show that the three rear-lock attachment bands retain
only 82.58%, 72.44%, and 78.12% of their baseline within-1-mm contact samples.
The visual gain is too small to accept that construction regression, and a
larger C2 amplitude would predictably increase the detachment risk.

A69 may be retried only from its last safe C0 with a revised
contact-preserving deformation interface—for example, root bands pinned or
carried surface-relatively with the receiver and displacement tapered into the
free lock lengths. That would be a replacement C1, not continuation of this v4
branch. The failed promotion gate must pass before another amplitude increase.

## Evidence reviewed

### Immutable candidate and renderer

- Candidate:
  `c1_live_hybrid/reimu_fumo_a69_c1_live_v4.blend`
- Candidate SHA-256 before and after rendering:
  `12f2a09d376d754c838fc2dddb9c528e59709e60119882bbf5563572903bcfe9`
- Diagnostic spec SHA-256:
  `ae81f6cb3f0478a5100bfc6ce25d6b425b76ed0a148d9933b649d93f24ac8ae2`
- Renderer: repository target
  `//projects/renders/cmd/fumo_review:render_packet`
- Blender: `5.2.1 LTS`, build hash `9e2066aef7ef`
- Render result: exit status 0; four nonblank 320 by 320 RGBA PNGs
- Manifest SHA-256:
  `ad8a0ead852cde1059ab1a3293ff92bea578ffe405901183d0bdbdc87bea15dc`

The candidate was rendered read-only with:

```sh
bazel_agent run //projects/renders/cmd/fumo_review:render_packet -- \
  --blend-file \
  out/reimu_fumo_attempt_069_head_cap_interface/c1_live_hybrid/reimu_fumo_a69_c1_live_v4.blend \
  --spec \
  out/reimu_fumo_attempt_069_head_cap_interface/c1_live_hybrid/diagnostic_render_spec.json \
  --output-dir \
  out/reimu_fumo_attempt_069_head_cap_interface/c1_live_hybrid/diagnostic_packet
```

### Controlling references

- Canonical front 25 cm, SHA-256
  `864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c`,
  controls front outline and identity.
- Canonical turn frame 07, SHA-256
  `28ef5155434ca05970b605047b4b0db223f7041fbe0db786d36535c841cbf9a4`,
  controls the near three-quarter receiver and lock overlap.
- Canonical turn frame 11, SHA-256
  `13d70b2ed0c790dc1938d66f8c250196419a3a4b215af2a113641924984f250b`,
  controls the first side silhouette.
- Canonical turn frame 25, SHA-256
  `349f56a207cf8dc230aa8f4de584be3a97f1e3450271faae02fba78809003a13`,
  controls the opposite side and asymmetry.
- Physical side, SHA-256
  `cbb39e70f95fa464f6dc94862e0300d15771f3ff4c046d005849891aca55a19d`,
  controls thin fabric layering and seated roots rather than metric depth.

The same-spec A68 C0 diagnostic control is under
`c1_live_hybrid/a68_c0_diagnostic_control/`. Its render isolates the C1 change
without a resolution, camera, lighting, or renderer difference. The original
640 px A68 C0 packet remains under
`out/reimu_fumo_attempt_068_sculpt_coupon/candidate_c0_pre_sculpt/render_packet/packet/`.

## Exact pixel comparison to A68 C0

All comparisons decode the PNGs to row-major 8-bit RGB before measuring.

| View | A68 SHA-256 | C1 SHA-256 | Changed pixels | Fraction | Mean absolute channel delta | Max channel delta |
| --- | --- | --- | ---: | ---: | ---: | ---: |
| Front | `cd96cee2cfc99b0f7df8a9a90b720504821eef8e3f724a876787dda595391520` | `f245b78d296952225957755b3a818a07682fa11971ae159db79c2923a19db222` | 1003 | 0.9795% | 0.0050 | 6 |
| Side | `63213a1ff52fdcaa2525a17a48a9ed899700591b78b074704707472fc9819580` | `774313e49ae59ae52b0b94feeee4473ae1dedf5025efbd012e29943580ad5506` | 6539 | 6.3857% | 0.6255 | 179 |
| Three-quarter | `97d0d8e7b0ab0ee7104f20399cc762cedc9f20fdc5472841bb9a46b72c26d231` | `6ba9a4ff86fead130a974dfd8b3fe650b323c7e645567bf6cbbb057287d9c46f` | 3895 | 3.8037% | 0.4815 | 189 |
| Mirrored three-quarter | `1ee647080f99e57a5b8751706f58eb44822bcc35d10ea85012f606eb56edbec4` | `d3288748bd2b146035b58ccdd9384a170aab6993dd521eb343622b3329d789f3` | 4318 | 4.2168% | 0.4937 | 185 |

The side and bilateral three-quarter differences are concentrated in the
rear/crown/lock region, while the front delta is sparse and visually
negligible. Difference witnesses are:

- `side_difference.png`
- `three_quarter_difference.png`
- `three_quarter_mirror_difference.png`

## View-by-view visual review

### Front

The canonical frontal outline, fringe, face opening, eyes, cheek locks, bow,
and body remain visually stable. There is no meaningful front silhouette or
identity regression. This view passes the local regression gate, although the
whole model remains far below final reference quality.

### Side

The rear boundary moves forward slightly and the broad depth reduction is
visible relative to A68 C0. This is directionally correct. The result still
has a long nearly horizontal crown and a near-vertical back wall, so it remains
much deeper and more monolithic than canonical frame 11. The reference reads
as a shallow receiver plus a long, thin, separately draping rear layer; C1
still reads as one solid cranium with attached pieces.

No obvious sharp dent or exposed beige patch appears at this resolution. The
lower root/receiver transition remains dark and compressed, however, and the
numeric contact loss prevents treating the lack of obvious daylight as proof
of sound attachment.

### Three-quarter

The cap reserve behind the near lock is slightly reduced without collapsing
the face opening. The same oversized rounded rear mass remains dominant, and
the transition into the rear locks is still abrupt rather than a soft,
surface-seated fabric overlap like canonical frame 07 and the physical side
photo. No new bald crown or gross local dimple is visible.

### Mirrored three-quarter

The intended silhouette movement survives on the opposite side, so the change
is not a single-camera accident. It remains modest and does not reach the
reference construction. The numeric loss of root contact also occurs in all
three locks, so bilateral pixel improvement cannot override the attachment
failure.

## Contact and continuation gate

The clean-reopen report records:

| Rear lock | Baseline samples within 1 mm | C1 samples within 1 mm | Retention |
| --- | ---: | ---: | ---: |
| Left asymmetric | 534 | 441 | 82.58% |
| Off-center main | 450 | 326 | 72.44% |
| Short right | 425 | 332 | 78.12% |

`diagnostic_contact_gate_pass=true` only establishes that the candidate is
useful for diagnosis. `promotion_contact_gate_pass=false` is the controlling
construction result. Proceeding to C2 would violate the attempt's explicit
stop conditions and trade attachment integrity for a silhouette that still
misses the source by a large margin.

## Absolute reference-fidelity gate

- Unlabeled same-subject recognition: **yes**, as Reimu, but not as a faithful
  match to the controlling plush variant.
- Overall reference likeness: **4/10**
- Macro silhouette and proportions: **3/10**
- Constructed-plush logic: **3/10**
- Identity-defining head/hair features: **5/10**
- Contact, attachment, and occlusion: **3/10**
- Intended plush-medium read: **3/10**
- Diagnostic presentation readability: **8/10**
- Major visible failure present: **yes**—the monolithic helmet/block rear head
  construction remains.
- Absolute decision: **reject**.

This is an internal diagnostic rejection, not a claim that A69 made no useful
progress. It proves that coupled deformation can move the correct projected
region while preserving the front, but its current root-contact interface is
not promotion-safe.

# A85 constructed head-cage cycle

## Outcome at this checkpoint

A85 remains open. `A85_head_cage_s12.blend` is the provisional head
foundation for collar and hair interface work, not an approved sculpt and not
a replacement for the tracked reusable asset. It is a new head-only scene;
the old C1b visible-form meshes were never imported.

The main process correction was to stop sculpting a rounded cube and directly
author a sparse quad cage whose rows own the crown, cheeks, chin, front/rear
planes, and side depth. Fixed front/profile renders were clean-reopened by the
repository-pinned Blender 5.2.1 after every kept or rejected delta.

## Candidate lineage

| State | Decision | Evidence-backed result |
| --- | --- | --- |
| raw / S00 | reject | New 152-vertex cage and attached references existed, but the first face-selection transform had zero geometric effect. |
| S01 | reject as final; valid parent | Corrected row transforms established plausible gross bounds. Blind review rejected the constant-depth rounded extrusion. |
| S02 | reject | Combined X/Y underside deformation created a pinched dark lip. |
| S03 | keep locally | Lower-cheek widening improved the front without changing profile, but the cheek band remained slightly narrow. |
| S04 | reject | Full-height Y taper became an egg/biconvex capsule in profile. |
| S05 | reject; retain front-arc lesson | Front arc worked, but bottom-row Y draw-in produced a folded-lip profile cue. |
| S06 | reject | Added the correct cheek width on top of S05 but inherited its profile defect. |
| S07 | keep locally | Isolated the successful front arc and preserved S03 depth. |
| S08 | strongest reviewed survivor | Added the remaining cheek width. Pixel gate passed the complete `.18-.43 Wh` cheek band at `.969-1.000 Wh`, lower width `.825 Wh`, height `.992 Wh`, max depth `.628 Wh`, and crown shoulder `.738 Wh`. |
| S09 | reject | Scalar arc reduction weakened lower-width accuracy without materially reducing the cavity read. |
| S10 | reject | Cosine redistribution worsened the underside shelf and introduced highlight speckling. |
| S11 | reject | An added low perimeter support ring did not materially remove the underside band; it was unnecessary topology. |
| S12 | provisional foundation | Returned to S08 and compressed only the two crown rows. All measurable foundation silhouette gates pass; blind review keeps it for integration. The underside remains an explicit collar-interface risk, not an accepted visible surface. |

## Bound artifacts

- S01 source SHA-256: `72df91b9ad3ccb1c65c88f63386408946135cbf301e68733133a91c5f8d28e63`.
- S08 source SHA-256: `dbd8ede7b58d369c9edb5163ccc51afd4b955c6bdf2c99a80838f9a85eaad2c2`.
- S12 source SHA-256: `982da6404ea6edcbb4432903e67dad4ee5c130a203a5a5727a374b773fc9ad8a`.
- S12 front SHA-256: `2d4a9f4c81780a71a4575a857ae144b250abb0e5043697f07c37af36ecf30d93`.
- S12 profile SHA-256: `84943201026e9ce9e14752a06d30d9140a490104fca6f662b7c0190b6f8d5d7a`.
- S12 manifest SHA-256: `e10a19bd6b54d83b25d041f6c4d68fd1b0c1be0fb5845b31e637686240dc3d93`.
- S08 measurement gate SHA-256: `717b964ae467faa1669094acb53c562284ad4d75ce53728a309a3fca34572ca8`.
- S08 blind review SHA-256: `80152d4edaef00b38bf86b989f0f24d4481a6377bc8edf08228557000d13814b`.
- S12 measurement gate SHA-256: `2df87c2d39a06019007398bc48964a370860c88b87d6d0f033a8a6b4b2952cb2`.
- S12 blind review SHA-256: `2b775e89abda220e6c2f358aab58f4d6681f8fc4da26e0ace58fe068a8032a0e`.

## Honest acceptance state and next move

The full Fumo is not accepted. The bare head still cannot independently pass
the constructed-medium gate because its exposed underside shades as a recessed
band. Repeated local underside tuning did not solve that read and began to
degrade passing silhouettes. Because the real plush collar and hair obscure
this interface, the result-oriented next work unit is a separate collar/body
context coupon around S12. Keep S12 only if the whole-model views hide the
underside cleanly with plausible fabric contact and no overlap; otherwise the
receiver interface, not the entire head, must be rebuilt.

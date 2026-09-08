# Correct inherited mirrored hand winding

First029 candidate93b2c163bff3e0a6e0d478b7b2312dd8578610d618b3570140cdf76ca400f386
visually improves cuff filling, but independent audit finds its right shell
consistently inward: signed volume-5.4361371e-6m3. Independent028 source
inspection confirms inherited inward orientation, volume-2.4726993e-6m3 and
front-hit normalY+.99055 versus left's outwardY-.99118.

029b is a winding-only correction from the same frozen028 input and exact
distal geometry field. Compute signed closed-mesh volume and reverse face
order only when negative; preserve every vertex, material region and weight.
The root coordinate hash is explicitly source-indexed, not a false claim that
the corrected right face winding is unchanged. Predicted result: outward
positive right volume, identical geometry positions and cuff-contact locations.
029 and its exact scripts stay immutable; source028 remains protected.

Cuff crossings are reported separately: new hand/trim intersections occur
inside the cuff's axial band, none beyond the checked distal lip. Some occur
on the first visible surfaces, so they are not all buried proximal contact.
Pixels show no obvious major clipping. Treat these as measured attachment-
band overlaps with limits, not an exhaustive collision-free pass or an
automatic demand to reshape all sleeve trim. Skirt edge checks have no
crossings and sampled gaps above4.2mm. Final technical evidence will qualify
the exact corrected candidate.

Observed029b: positions, material regions and weights are exact; right
volume is positive. The contact-location prediction needed qualification:
reversed quads select different internal diagonals, so right-hand contact
was rescanned rather than inferred. Matched cloth-edge hits move at most
.051432mm; all remain in the same cuff band, with no new exposed-distal
crossing region. Skirt separation stays4.331653mm sampled. See the029b
technical addendum for exact counts and sampled trim-depth changes.

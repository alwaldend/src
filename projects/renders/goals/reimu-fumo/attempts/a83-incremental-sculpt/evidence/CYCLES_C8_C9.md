# A83 cycles C8--C9 evidence

## C8: bounded cap-and-receiver tangent-plane edit

C8 restarted from the retained C1b checkpoint and changed only localized
camera-side neighborhoods on the continuous hair cap and head receiver through
reversible shape keys. The edit selected `1,244` cap vertices and `1,287`
receiver vertices, but the fixed front and three-quarter renders retained the
same rolling helmet highlight and lacked a readable fabric-panel break. The
implementation-independent review found no visible preference over C1b.
**UNDO**.

- Candidate blend:
  `sha256:6d394807cd33062a3067ed4de632c3f16dfb4c4394ecb0c7a4649796f2983b70`.
- Blind verdict:
  `sha256:71bb9cd5c77f200f5f59fb0f62397570c1c938f8bf3249eaad98b8c1a0eb8bf9`.

## C9: paired native Blender sculpt stroke

C9 again restarted from exact C1b and exercised Blender's native Grab sculpt
operator on the same paired cap and receiver neighborhoods in a realized
Xvfb viewport. The receiver stroke moved `195 / 11,184` vertices with a
`0.9804 mm` maximum displacement; the cap stroke moved `196 / 9,449` vertices
with a `0.9919 mm` maximum displacement. Both operators returned `FINISHED`,
and the candidate saved with the intended object-local changes.

The fixed 512-pixel front and three-quarter renders nevertheless remained
visually indistinguishable from C1b in the named defect. No plane, lap,
gusset, panel thickness, or interruption of the rolling side highlight became
legible. Numerical mesh and pixel activity did not produce reference-fidelity
progress. **UNDO**.

- Candidate blend:
  `sha256:2c46187fd1933b7afff7c092d841883ec0e87d5096d403aef00fc9288ee83219`.
- Front render:
  `sha256:e326c9ed7b246dc9e7245f3bb6c9d84046fcc3964d80a307866bb81aeab59ab5`.
- Three-quarter render:
  `sha256:556e11dc61f84e436a6ce1f353502816fe87e314abbdbb5c84566fb75c459699`.
- Blind verdict:
  `sha256:8783051a69a3ca4658011f6f446f7306e625c15e293e6d2ebb1c295fb4565542`.

Protected rung 003 and the tracked reusable model retained their expected
hashes throughout.

## Process and structural conclusion

C8 established that a scripted coordinate deformation can be mechanically
localized yet visually inert. C9 then removed uncertainty about access to the
native sculpt operator: the operator works, but a roughly one-millimetre
continuous-shell stroke is still the wrong discriminating edit for this
defect. The process now judges absolute reference fidelity before fixed A/B,
requires a visible preference before keeping a candidate, and does not count
geometry or pixel churn as progress.

The next bounded test must change the local owner without rebuilding the head:
keep C1b's cap, head, front landmarks, bow seat, and bounds fixed, then add one
shallow, broad, root-seated side-hair panel coupon with a lapped edge and
explicit felt thickness. It passes only if the three-quarter view reads as a
soft padded plane rather than a helmet, card, armor plate, or floating strip,
while the front remains stable.

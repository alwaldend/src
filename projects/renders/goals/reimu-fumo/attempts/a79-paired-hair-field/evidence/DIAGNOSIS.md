# A79 frozen-geometry diagnostic

Status: **NON-CANDIDATE — reject this representation**

The five-view sheet is
`a79_non_candidate_five_view.png` (SHA-256
`7ca7566127546b645f54a80f8c104496da09bd4bdc153e4777b4d1beece1b277`).
It renders the byte-exact frozen geometry with SHA-256
`2483e49d684836cc4af349e381fd35d2ab91fcf65e7be8a688c980b81988931d`.

## Pixel verdict

The representation fails before material detail or topology refinement:

1. The front and three-quarter views expose a very large beige bald crown.
   The proposed crown field covers only a narrow perimeter and the fringe; it
   does not form the photographed brown hair cap.
2. The rear base reads as a long, rigid rectangular curtain with nearly
   vertical sides and an abrupt bottom. It does not match the compact,
   rounded, layered rear hair mass in the turntable.
3. The broad asymmetric leaf reads as a thick planar hanging board. In side
   and three-quarter views it is too long, too isolated, and insufficiently
   curved around the head.
4. Crown, rear base, and leaf do not visually join into one manufactured hair
   assembly. Their roots produce gaps, hard tangencies, and abrupt changes in
   depth.
5. Smooth shading and the copied brown material cannot rescue these macro
   silhouette and coverage failures.

The useful conclusion is categorical: mesh validity and frozen-surface
metrics were not predictive of reference likeness. The next attempt must
constrain visible brown coverage and the complete multi-view hair silhouette
before spending time on topology density or construction details.

## Safety and reproducibility

- Blender: repository-pinned 5.2.1 LTS through `bazel_agent`.
- Protected rung003 remained SHA-256
  `c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`.
- Frozen geometry remained unchanged.
- The scratch blend clean-reopened successfully; see
  `reopen_validation.json`.
- This directory is disposable evidence and is not authorized for promotion.

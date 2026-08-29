# A74 P0 self-review

## Verdict

Reset A74 at P0. The candidate proves that a source-traced open saddle can be
built deterministically without a closed receiver, but the resulting whole
plush is substantially worse than the protected rung and fails every visual
gate. Do not add rear pockets or tune this representation.

## Exact tested subject

- Candidate blend:
  `sha256:d25cc2b5940c16ab694f72b39ed5f409baab8caf01df38bd0b0759723447ec44`
- Protected parent:
  `sha256:c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`
- Source-only geometry contract:
  `sha256:76a00ad04f51a52f7252f476f122405ee3e0bf182b37b14994e2e8e319de40b4`
- Render manifest:
  `sha256:54f872ba5ba2d99d52e198ead809cecb3b5c645d4126e8ebc86ac31da4cd6dbf`
- Five-view contact sheet:
  `sha256:b15a8ac5a9f58075b9f8f4d3e4f2da01fae9a62d3e25b69f68480debaddb56c1`
- Support-mask sheet:
  `sha256:b244f2979e6fd809ad920881ee21193852ad7840302bcf3f58d0089f8fb926cd`

Pinned Blender 5.2.1 built and reopened the candidate. The protected parent
hash was verified before append and after save. Exactly 15 audited legacy
objects were hidden; all 22 frozen face/bow witnesses retained their geometry,
transforms, materials, and visibility.

## Mechanical result

- The saddle has 480 base vertices, 449 quad faces, an authored paired skin,
  and the expected open rear-root boundary. It has no receiver, rear wall,
  Solidify modifier, pole, or hidden closure.
- The left and right locks are independent meshes with 70 vertices and 68
  quad faces each. They do not use mirroring or shared mesh data.
- The support is an open one-sided 81-vertex, 64-quad front patch.
- All four candidate objects have identity transforms and new materials.

These mechanics validate the experiment but do not satisfy the goal.

## Decisive visual failures

1. The face support is a flat rectangular card. Its square corners and nearly
   white render dominate the front and both three-quarter views.
2. The brown saddle reads as a smooth human haircut or hard canopy, not as the
   layered, softly stuffed brown hair of the references.
3. Both locks are pointed detached blades rather than blunt flexible plush
   locks. Their reported widths are also about `.20 Wh`, outside the intended
   `.14-.18 Wh` construction band.
4. The side view is card-thin with a large unsupported gap.
5. The rear exposes the face and eyes because P0 hid every old rear owner while
   explicitly forbidding replacement rear geometry.
6. Component masks show support pixels in side, both three-quarter views, and
   rear, violating the zero-leak gate.

Independent implementation-blind review scores macro identity `2/10`, plush
construction `1/10`, and contact/integration `1/10`, all below the required
P0 score of 6 and far below final acceptance.

## Process and interface audit

The source-contract and fixed-camera loop worked: the candidate reached a
reliable reject packet quickly, and no protected asset was touched. The
failure is the attempt interface itself. A "front-only" module cannot be
judged honestly in whole-plush context after its boundary hides every rear
hair owner; the resulting bald rear is guaranteed before geometry begins.
Likewise, preserving facial appliques while replacing their stuffed carrier
with a one-sided rectangular diagnostic sheet guarantees a card face.

The next attempt must not repeat another full head rebuild. It should first
construct a polished, rounded stuffed face cushion as a small module while
retaining enough protected hair context to judge integration. Only after that
module is preferred should a later atomic whole-hair candidate replace all
silhouette-owning hair objects together. This separates a well-defined support
problem from the still uncertain full hair construction without accepting a
partial bald model.


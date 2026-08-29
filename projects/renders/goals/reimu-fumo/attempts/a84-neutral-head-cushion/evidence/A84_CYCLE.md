# A84 neutral head-cushion cycle

## Absolute verdict

**Reset the rounded-cube scaffold.** The attempt proved that a broad elastic
stroke can taper the lower front silhouette without the hard folds produced
by ordinary Grab, but the same topology did not provide controllable crown
ownership. The bare cushion still reads as a rounded rectangular primitive,
not a stuffed gusseted head cushion.

The canonical front and canonical turn remained the controlling references.
The old C1b visible-form meshes were never imported.

## Candidate ledger

- `S00`, SHA-256
  `7dc1bf434a2c8e2300d4d1b2480d622f4465e0828226aaf6c57f6233cfcc9246`:
  measured rounded-cube baseline. It satisfies the coarse width, height, and
  depth envelope but has a flat top, broad lower corners, and a uniform box
  read.
- `S01`, SHA-256
  `3edc8674c3667a3178f68a65789c0cbac303ed2c1f4e409a5159bf59ebc67568`:
  ordinary Grab at Multires level 1. The pinned packet showed only a tiny
  silhouette delta. A stronger unsaved probe produced categorical hard
  corner creases and was restored from S00.
- `S02`, SHA-256
  `af314e5ca27538fe72bcd8908684526056ad77bf4cbf1f9ae1701fa843867de8`:
  two mirrored Elastic Grab strokes at the coarsest Multires level. Keep only
  as a local process result: the lower front taper improved without a profile
  regression or pinched folds. It is not an acceptable head owner.
- `S03`, SHA-256
  `a838e979ceb16ceb3b736438bd2103d70cf7bb654f08f11737e82b5f93cbc778`:
  crown stroke started on the silhouette boundary. Its pinned front and
  owner-ID images are byte-identical to S02. Undo as zero effect.
- `S04`, SHA-256
  `797af9bf0293195dad58cfaad1b71f73ce41c6b4f79b6c0cdfa31e47167f00c6`:
  crown stroke moved inside the front plane. It produced only a negligible
  crown change and a localized beauty-pass shading singularity. Undo.

## Process evaluation

The interactive X-display Flatpak host supplied a valid `982 x 851` sculpt
viewport, and repository-pinned Blender 5.2.1 clean-reopened every frozen
candidate. This authoring/render split works. The main latency remains Bazel
startup plus approximately 11 seconds of packet rendering; one snapshot
should continue to feed one pinned packet without duplicate live renders.

The failed variable is no longer brush choice. The rounded cube has no
explicit front silhouette rows or crown/lower-chin control loops, so broad
strokes either under-act or concentrate deformation at unsuitable corners.
Further brush tuning would optimize the wrong scaffold.

The next attempt should construct a low-resolution quad cage with explicit
crown, cheek, lower-chin, front-plane, rear-plane, and side-gusset loops in
Edit Mode. It should establish the measured silhouette and shallow depth
before subdivision. Direct sculpting should then add only broad stuffing
softness and restrained asymmetry; it should not be responsible for inventing
the primary silhouette.

## Review artifacts

- `A84_front_progression.png`, SHA-256
  `478eaeb07dea8374d8ac56f8a0fea29c7620945ed4146b91880767b1a0d6f251`.
- `A84_profile_progression.png`, SHA-256
  `a2fcc8302ed359fb03f2e4bb36ca6207a754f17bccdf5824a8067881937fbb58`.
- The exact S00, S02, and S04 packets and manifests are under
  `packets/head_s00_pair`, `packets/head_s02_pair`, and
  `packets/head_s04_pair`.


# A91 first-pixel gate

## Exact candidate and evidence

- Candidate Blend SHA-256:
  `6007f96f6fc11b662bff5ddad99e91749c4a3b8f2013f4a7af57d08815558fe7`.
- Front render SHA-256:
  `645858bfbc8a67497afb74a8db40d8c8a3542469c963eb0a1fa8c8869cc55772`.
- Worst three-quarter render SHA-256:
  `b13d915ecdb29bb0631c7e558cc0bb10892116a6e83d2135003ce3993f67bbe4`.
- Auxiliary three-quarter render SHA-256:
  `638cfdc15c2b3422cc362618475981190e3c50c40d8fab0b6c1ddcb554d19e0f`.
- Candidate manifest SHA-256:
  `4600ec2eb81b06e0b42f51da1b8eff3df3745a1e82bb4a56e22873fe8354393b`.
- Historical rung003 comparator SHA-256:
  `c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`.

The candidate and pixels remain in the ignored diagnostic directory
`out/reimu_fumo_attempt_091_whole_plush_bakeoff/candidate_fresh/`. They are not
accepted or promoted artifacts. Their exact observations and identities are
retained here; the deterministic diagnostic script remains beside them for
the current workspace only.

## Technical evidence

Pinned Blender 5.2.1 produced the final batch in 18.238 seconds; front pixels
appeared at 5.712 seconds. Clean reopen found 49 objects, 37 meshes, 11
materials, all required macro components, and bounds
`0.242874 × 0.135 × 0.249337 m`. The tracked Blend and controlling reference
hashes remained unchanged.

The initial renders were overexposed and could not be judged. A same-geometry
lighting correction restored readable pixels; this was an evidence repair,
not a geometry correction. No geometry correction was attempted after the
readable pixels exposed a representation-level miss.

## Independent correction review

One fresh `gpt-5.6-sol` Ultra reviewer ran with `fork_turns=none`, no mutable
access, and no delegation. The reviewer verified the final pixel digests and
returned `RESET`.

Scores for A91:

| Category | Score |
| --- | ---: |
| Overall likeness | 3.5/10 |
| Macro silhouette/proportions | 3.0/10 |
| Construction | 1.5/10 |
| Identity features | 4.0/10 |
| Contact/occlusion | 2.5/10 |
| Intended-medium read | 1.5/10 |
| Presentation | 6.0/10 |

Decisive failures were the tall non-seated profile, deep rectangular
hair/head shell, obscured canonical eyes and face opening, cone dress, block
sleeves, tube-like feet, and rigid card bow. Rung003 remains about 2 points
better on whole-subject likeness. The candidate does not clear any required
6/10 parent-selection category and contains automatic-fail wrong-category
construction.

## Coordinator disposition

`follow RESET`. The visible miss agrees with the author's hard stop and the
coordinator's implementation-blind inspection. A bounded local correction
cannot turn the card/cone/block representation into a seated sewn plush. A91
is rejected without promotion, detailed panels, six-view work, materials,
retopology, or rigging.

The next discriminator is the planned A93 direct-sculpt input coupon: three
broad independently directed deformations must move meaningful low-frequency
silhouette regions without dimples, zero-effect output, or camera-only
illusion. Only a passing coupon authorizes another complete low-resolution
manual sculpt.


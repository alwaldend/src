# A73 P0 profile-cage review

## Verdict

Reset at P0. Do not make the allowed row correction and do not build the
visible-layer coupon from this cage. The explicit profile and cap mechanics
all work, but the evaluated pixels remain a rounded box/mattress in front,
rear, top, and bottom and an egg with a long face-side wall in profile.

## Mechanical evidence

- 514 vertices, 1,024 edges, and 512 all-quad faces.
- Every edge has exactly two incident faces; there is no axis pole, triangle
  fan, n-gon, or high-valence cap center.
- Static gates pass: `Wh=132 mm`, evaluated depth `105.6 mm`, evaluated height
  `116.16 mm`, rear maximum at `t=.69`, lower-rear undercut `.065 Wh`, and
  face-zone reach variation `.005 Wh`.
- Pinned Blender 5.2.1 built, reopened, and rendered the exact candidate in
  eight fixed views.

## Pixel evidence

- The front retains a long, nearly uniform inflated wall and broad horizontal
  plateau. It is the prohibited rounded mattress, not a plausible hidden
  stuffed support.
- Both profiles remain egg-like. The high rear maximum is numerically present
  but visually too weak; the lower-rear turn-in does not read decisively.
- Rear, top, and bottom remain rounded rectangles, proving that the failure is
  the loft family rather than one favored camera.
- Independent implementation-blind review scores receiver macro silhouette
  4/10, below the 6/10 P0 gate, and calls it not directionally viable for the
  mandatory hair coupon.
- Reference scanline audit independently shows that same-crown registration
  would protrude by `.052--.178 Wh` through the canonical visible crown rows;
  no vertical offset alone fits all controlling rows.

## Exact identities

- Candidate blend:
  `sha256:1fe46cdd11c329b2511426eaee2c099603bbf6f5a2a04cbfa2d5186c2e5aec33`
- Coordinate set:
  `sha256:b728903ce7f4ba870c86fe7bbf76922ea0fae3a31dc5c4e86f7d81e7f09f6e1d`
- Build report:
  `sha256:0c6649fd5ad73b5bcb48945a4d177c45d787e50d3186d2762218fdd6de0c2b2f`
- Profile contract:
  `sha256:9718ff7d830abeca6ada2fa5b51ed8fbaa790fd78d9064ee4809d7918e93ed35`
- Render manifest:
  `sha256:dd060e2adee944d86e2cb43e798be36d5b15dd6276bc7ee020175d4ba89a09d3`
- Eight-view contact sheet:
  `sha256:609aaa15353caf4a723a00258c6fac4ca3986c689c5378e83f59f7315cc67e6f`
- Protected parent:
  `sha256:c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`

## Whole-process conclusion

A73 fixes A72's false preflight, hidden formula coupling, pole closures, and
evidence incompleteness. Those process gains made the failure cheaper and more
conclusive, but they did not improve the plush. The deeper error is asking an
unobserved hidden receiver to carry visible identity. The next representation
must start with the measured visible brown head/hair envelope and its distinct
compact field/free-leaf ownership. A hidden support may follow that visible
construction; it may not define it.

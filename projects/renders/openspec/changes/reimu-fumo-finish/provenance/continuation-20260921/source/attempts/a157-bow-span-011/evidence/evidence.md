# Bow-span failure receipt

## Protected input

The recovered A157 source and its task-owned working copy both remained at
SHA-256
`433d08ad36be488bb16e4221a85f831d4390660c258a43ea0b08775811574b73`.
No source or checked-in blend was modified.

## Native strokes

The mapped XWayland Blender 5.1.1 host used the essential `Grab` brush with
scene-space diameter `0.11`, strength `1.0`, symmetry disabled, and nine
timer-spaced native events per planned stroke.

The left-tail modal operator completed and moved 306 of 650 vertices. Maximum
local displacement was `0.0368124097`, mean changed displacement was
`0.0132040498`, and the protected root-region maximum was exactly `0.0`. The
front-view outer extremum moved 265 live-view pixels as requested.

The identically configured right-tail modal operator opened, consumed all
nine events, and exited normally, but moved zero vertices. The right mesh had
no mask or hidden-vertex attribute; front-face restriction and every
automasking option were off; its target normal faced the camera comparably to
the successful left target. The settled failure is therefore the assembled
viewport surface-picking path at that projected outer extremum. Occlusion by
overlapping bow parts is the leading hypothesis, not a proven fact. A later
trial must isolate the target surface or replace the selection mechanism; it
must not resend this stroke.

## Automatic visual rejection

The partial failure state was saved only as evidence:

- blend:
  `out/reimu_fumo_finish/attempt_011_a157_bow_span/a157_bow_span_011_left_only_failure.blend`
- SHA-256:
  `01de615a44baa6b6ca20ead49d0dd2d7f79eefea689b0fa33ec68c13e5289f05`

Repository-pinned Blender 5.2.1 clean-opened the exact bytes and rendered the
fixed 512 by 512 front camera. The render SHA-256 is
`2ae0afbcb9c5e01a83da1f8e40a38539447c95e7443e1855d73df970c204ff5f`;
its manifest SHA-256 is
`4cc59662c0bcbba0c7dce4f77e190745060a35f60363b01c2dc891bd8ec69485`.

The render shows the left tail stretched into a long pointed fin and cropped
by the frame while the right tail remains unchanged. This violates the
plan's fin, framing, symmetry, and complete-subsystem gates before independent
review. The partial geometry is rejected and must not become a new baseline.

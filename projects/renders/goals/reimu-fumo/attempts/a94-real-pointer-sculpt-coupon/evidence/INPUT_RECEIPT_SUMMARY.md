# A94 exact input-receipt summary

This durable summary binds the raw task-local JSON and pixel artifacts without
promoting temporary binaries or images into tracked evidence.

## Verified pre-input state

- Pinned Blender: `5.2.1 LTS`, foreground, Sculpt mode.
- Exact session working-copy SHA-256:
  `0004a3f0bc4987a250f7028b8697a7f740c70866ad8b60f7181e2b4eafa96400`.
- Native Essentials brush: `Grab`, scene radius `0.050 m`, strength `0.40`.
- Orthographic view; requested world direction `(0, 1, 0)`.
- Projected world `+Y` distance `0.005 m`: `26.0417485676 px`.
- Projected world `+X` depth-axis distance: `0.0 px`.
- Planned Blender window-local bottom-left path:
  `(1160.4375, 790.6341)` to `(1186.4792, 790.6341)`.
- Raw gesture-plan SHA-256:
  `b5c7727de26704deb36a25f4c5781235533d41994f232d1ff1a661f282b8bb82`.
- Raw window-state SHA-256:
  `11340d0ad9bc760cb8d795f074c1ee09089531eb52371c54f9cd612827fe184b`.

The disposable target was normalized to an identity object transform by
baking only its copied mesh. Pre/post world-coordinate SHA-256 was identical:
`8893cfd6d8b5d96742fcef2f6ea8fc13750eb1418e40f9e9c6ff4197881a334b`,
with measured maximum error `0.0 m`. The source stayed untouched.

## Sole XTest delivery

- Mapped X11 window: `2097154`, origin `(0, 0)`, size `2560 x 1440`.
- Quantized window-local bottom-left path: `(1160, 791)` to `(1186, 791)`.
- Delivered screen top-left path: `(1160, 648)` to `(1186, 648)`.
- Eight held motions over `320 ms`, button 1, injector exit `0`.
- API: `XTestFakeMotionEvent` / `XTestFakeButtonEvent`.
- Raw injector-receipt SHA-256:
  `fcd8e964828380b0383c9157ae55dd08828a0a2c591ad834e0dc2ab3017e328b`.

Blender produced no raw event receipt, post snapshot, or `DONE` marker after
more than 30 seconds. Therefore XTest delivery cannot be reconciled to an
in-region Blender press/motion/release sequence and fails provenance.

## Independent pixel invalidity

The nominal live `before_view3d.png` exists and hashes to
`f0e354bbb43ea78d799f8506cccc9ec658b0aed2d85d937dfda8106bd5bf4484`,
but post-stop inspection found every pixel black (`2100 x 1389`, minimum,
maximum, and mean all zero). It is not valid visual evidence and should itself
have stopped input. The fixed-camera baseline renders were nonblank but cannot
replace the required live viewport pair.

No second gesture, `-X` branch, `-Z` branch, post geometry claim, visual claim,
or parameter adjustment occurred.

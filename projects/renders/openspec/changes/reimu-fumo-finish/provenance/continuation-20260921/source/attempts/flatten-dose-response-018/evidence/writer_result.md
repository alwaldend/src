# Attempt 018 Blender-writer result

Status: failed at the 27-stroke plane cap. Root and tip tests did not run.

## Exact input and output identities

- Baseline input:
  `out/reimu_fumo_finish/attempt_017_localized_sculpt_coupon/localized_sculpt_fixture_baseline_frozen.blend`
- Baseline SHA-256:
  `6f49bd4e0a8af6b45870d9d4224a520c1398e52e6e9c42f8fb5bee7b8c17118e`
- Baseline coordinate digest:
  `41ee23670f67335ac070d95bd782436f53405034f1e24efdebdad709f7d47df2`
- Partial input:
  `out/reimu_fumo_finish/attempt_017_localized_sculpt_coupon/localized_sculpt_coupon_017_partial_plane_fail.blend`
- Partial SHA-256:
  `2428de0a0b65e572de9437a8d3ef35f1ee21c18bd9dbf27ff01de1816418c0bd`
- Partial coordinate digest:
  `f2fdcbdaea90335b9de47861e8ce59b64896ec81853d19ee2ea675725d9fb16e`
- Terminal failure state:
  `out/reimu_fumo_finish/attempt_018_flatten_dose_response/flatten_dose_response_018_final_fail.blend`
- Terminal file SHA-256:
  `b9ad59c15901b0fe22cd96208e60792e2dc35e8bc0f3e73d0e1a8181697dd6fc`
- Terminal coordinate digest:
  `fc3d57571800ecade593e64d2750687f883c7da9ef6177fa9eec622a8346c2e0`

The baseline and partial input hashes were unchanged after the attempt.
Saving the terminal state did not change its pre-save coordinate digest or
metric.

## Frozen metric reconciliation

The authoritative attempt-017 measurement window is:

`front layer; abs(current vertex x) <= 0.65; abs(grid u) <= 0.65`

Using embedded baseline X instead of current-state X incorrectly included 16
vertices that Flatten had moved outside the window. That alternate 837-vertex
selection produced variance `0.013351381197580402` and was rejected before
adding any stroke. The frozen current-coordinate definition selects 821
vertices and exactly reproduces the attempt-017 partial variance
`0.01352948526060966`. Baseline variance is
`0.014215173048885853`.

## Held variables and event receipt

- Brush: essential `Flatten/Contrast`, type `PLANE`.
- Strength: `0.6499999761581421`.
- Scene-locked size: `0.5799999833106995`.
- View: front orthographic, distance `3.9000000953674316`, quaternion
  `(0.7071067690849304, 0.7071067690849304, 0, 0)`.
- Mask: 1,548 plane vertices unmasked and 7,958 vertices masked.
- Every added block replayed the exact saved 33-event pattern: three hover
  moves, three presses, 24 drag moves, and three releases.
- Eight blocks delivered 264 events and 24 added native strokes without an
  event error.
- The requested timer interval remained `0.12` seconds. The full delivered
  timestamps, coordinates, event phases, block numbers, and interval deltas
  are embedded in `a018_event_log_json` in the terminal blend.
- The only varied input was cumulative identical Flatten dose.

## Complete block metrics

Every checkpoint used an 821-vertex current-coordinate window. Every block
had zero non-plane displacement and zero `CONTROL_plane` displacement.

### Cumulative 6 strokes

- Variance: `0.012768618841366354`.
- Reduction from baseline: `0.10176128018595432`.
- Absolute block drop: `0.0007608664192433057`.
- Maximum cumulative displacement: `0.05335054547476773`.
- Changed vertices from partial: 964, all in the plane region.
- Coordinate digest:
  `da81e585891e0a73618a917e4d155653d400576b3bf6e694e2d13669a349d6b9`.
- Low-response block: false.

### Cumulative 9 strokes

- Variance: `0.012129297108030106`.
- Reduction from baseline: `0.14673588099718782`.
- Absolute block drop: `0.0006393217333362476`.
- Maximum cumulative displacement: `0.06822419574281582`.
- Changed vertices from partial: 967, all in the plane region.
- Coordinate digest:
  `fe13223658fef0e8820443a8651761f6383f4288a26c533dede77e9236db71f4`.
- Low-response block: false.

### Cumulative 12 strokes

- Variance: `0.01157981745100081`.
- Reduction from baseline: `0.18539032826558477`.
- Absolute block drop: `0.0005494796570292958`.
- Maximum cumulative displacement: `0.0798227757766021`.
- Changed vertices from partial: 969, all in the plane region.
- Coordinate digest:
  `fa715e1c840785770b2aab0e3b9cdbb6dbc66bef3b2b76a2a992ec1e8e9be3ae`.
- Low-response block: false.

### Cumulative 15 strokes

- Variance: `0.011100507927945294`.
- Reduction from baseline: `0.21910849134437216`.
- Absolute block drop: `0.00047930952305551625`.
- Maximum cumulative displacement: `0.089495743765136`.
- Changed vertices from partial: 974, all in the plane region.
- Coordinate digest:
  `65a8cca8a36e691fcfffe5deab71231c11f31c0dfccc47b773521c06722cbe6a`.
- Low-response block: false.

### Cumulative 18 strokes

- Variance: `0.01066979347286459`.
- Reduction from baseline: `0.24940811932635176`.
- Absolute block drop: `0.00043071445508070447`.
- Maximum cumulative displacement: `0.09797136617882572`.
- Changed vertices from partial: 974, all in the plane region.
- Coordinate digest:
  `61ab6eaabcbfc8477e296c60db2c7e15ad3b4f6f0ac8d12fb2f74966ddbb20e3`.
- Low-response block: false.

### Cumulative 21 strokes

- Variance: `0.010283509693619804`.
- Reduction from baseline: `0.27658216623498666`.
- Absolute block drop: `0.00038628377924478617`.
- Maximum cumulative displacement: `0.10544778471843798`.
- Changed vertices from partial: 974, all in the plane region.
- Coordinate digest:
  `905823e7d1442d1f5e51aa741e78da412ffad38eb13cf24144216397de88cbfb`.
- Low-response block: false.

### Cumulative 24 strokes

- Variance: `0.00992794400483365`.
- Reduction from baseline: `0.3015952763507318`.
- Absolute block drop: `0.0003555656887861535`.
- Maximum cumulative displacement: `0.112262968919496`.
- Changed vertices from partial: 975, all in the plane region.
- Coordinate digest:
  `0c6112281c5d29217cfd4f3afb406811190ac810a4341d9ad067a1856e5b381f`.
- Low-response block: false.

### Cumulative 27 strokes

- Variance: `0.009603326433540808`.
- Reduction from baseline: `0.3244312678772848`.
- Absolute block drop: `0.00032461757129284216`.
- Maximum cumulative displacement: `0.11849009312857924`.
- Changed vertices from partial: 975, all in the plane region.
- Coordinate digest:
  `fc3d57571800ecade593e64d2750687f883c7da9ef6177fa9eec622a8346c2e0`.
- Low-response block: true; this was the first consecutive low-response
  block.
- Stop reason: `27_strokes_without_plane_pass`.

The fixed low-response floor was `0.0003428438941380965`. Variance decreased
at every checkpoint, maximum displacement stayed below `0.20`, and isolation
remained exact. The terminal reduction was `32.44312678772848%`, below the
required `35%`. The 27-stroke cap therefore stopped the coupon before the
root and tip tests.

## Native undo and live final state

- Native undo delivered 24 simulated Ctrl-Z press/release pairs: 48 events.
- Undo error: none.
- Live post-undo coordinate digest:
  `f2fdcbdaea90335b9de47861e8ce59b64896ec81853d19ee2ea675725d9fb16e`.
- Exact partial-digest match: true.
- Changed vertices from partial after undo: 0.
- Maximum delta from partial after undo: `0.0`.
- Live host: Blender 5.1.1, foreground, Sculpt mode.
- Live file-path label remains the terminal failure path, but in-memory
  geometry is the exact partial input after undo; the live scene is clean and
  was not saved again.

No Reimu blend was opened or modified. No root or tip operation ran. Pinned
clean-open, audit, and rendering remain coordinator work. The Blender writer
lock is relinquished after this receipt.


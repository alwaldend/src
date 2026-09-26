# Camera and body occlusion investigation 121

Outcome remains open and execution active. Studies 118–120 are closed
rejected. Before another model change, investigate the level frontal camera
and the unseen body's projection. This is a diagnostic investigation, not a
new accepted review camera or a modeling candidate.

## Evidence and decision

The canonical front shows background between the feet below the white hem.
All recent drafts show a broad central red mass down to the floor. Study 120
reaches the required hem width and depth without sampled support penetration,
yet the discrepancy persists. A level orthographic camera projects height
independently of depth, making it impossible for a higher foreground hem to
occlude the low rear seat. The contract's `fixed_views` explicitly specifies
a level front camera; that binding remains unchanged. Previous front overlays
align scale and image position but do not independently establish that this
pitch matches the reference. An initial reading of only the contract's
`camera` section missed its explicit transforms; root corrected that reading
before any publication or change to the acceptance packet.

For an elevated orthographic view, projected height is
`z * cos(elevation) + y * sin(elevation)`. Using the present model's hem and
seat only, approximately 21.5 degrees could align a seat-bottom point with
the foreground hem. This is a diagnostic calculation, **not** an estimate
of the reference's camera angle. Choosing a camera solely to hide bad
geometry would not establish reference fidelity.

**Proceed with bounded diagnostics.** Independently inspect reference
turntable motion and multiple depth/occlusion cues. Keep canonical-photo and
turntable camera evidence distinct. Do not infer the hidden pelvis or recolor
it from an unverified projection. The physical side photograph identifies
pale fabric near leg roots but does not prove the entire pelvis shape/color.

## Bounded plan

Protected input is saved study120, SHA-256
`fe7f3f638ac1b59f59c744af0533e3ffb88f12d26bd3b1784bc5f17285355892`.
Root uses pinned background Blender with four threads to render two additional
front diagnostics at 12 and 22 degrees elevation. Keep the original fixed
packet unchanged, preserve orthographic scale and lights, and register the
front ground point to the existing image height. Do not edit or save model
geometry. Record all camera transforms and any model hash changes; the input
must remain unchanged. These views may explain a failure but cannot replace
the frozen acceptance packet without independent calibration evidence.

Reference analysis may extract frames into local `assembly_121/reference_frames/`
and estimate a feature's vertical/horizontal motion ratio when a rigid turn
and identifiable landmarks permit it. Report uncertainty or inability to
constrain elevation rather than selecting a convenient angle. Stop after
these two diagnostic renders and the bounded reference analysis, then record
what is established and what remains unknown.

## Result

Both diagnostics completed in the pinned runtime and the protected model hash
remains unchanged. Local `assembly_121/camera_diagnostics.json` binds each
image to its camera transform. At 12 degrees the visible seat is reduced but
still substantial. At 22 degrees it is mostly occluded, while the face and
crown projections shift and the bow clips the diagnostic frame. Neither is a
valid replacement presentation or acceptance view. Camera elevation affects
the failure, but selecting 22 degrees to hide it would not validate the model.

Independent reference analysis tracked two black toe-cap region centers in
the canonical turn GIF. The 30 frames each last 100 ms and show a full turn:
approximately frame 4 front, 12 side, 19 rear and 25 opposite side. Equal
angular increments remain unverified. Under a uniform-turn assumption, the
two sinusoidal fits give nominal elevations 10.14 and 8.86 degrees, with
position residuals about 0.5–2.3 pixels. Visible-region centroids are not
material landmarks; occlusion, shape and alignment cause systematic error.
The suggested 7–13 degree range is therefore a provisional diagnostic range,
not a calibrated interval. Raw marks and method are preserved in local
`assembly_121/reference_frames/trajectory_diagnostic.json`.

The canonical photograph is a separate capture. No sufficiently reliable cue
was established to transfer the GIF estimate or determine the photograph's
elevation independently. The reference variant authority remains unchanged;
this investigation does not establish a different variant. Keep the exact
review contract and its level fixed cameras unchanged.

**Investigation complete; no model or camera accepted.** The observed missing
negative space is real, but its division between projection, hidden body
contour and garment layers remains unresolved. Next, trace the visible seated
body/leg contours and under-hem negative space across the controlling images,
and independently test projection consistency before committing to a new
support volume. Do not assume the current broad flat-bottom red seat is a
verified input; equally, do not recolor or shrink it merely to erase a bad
silhouette. Any model trial must bind a new measured body/contact hypothesis.
Cloth fixture calibration remains necessary before using numerical drape
again. Execution stays active with no external blocker.

## Session ergonomics

Four new full candidates received clean-open four-view packets and four
implementation-blind pixel reviews; two additional images tested camera
sensitivity without changing model bytes. Two isolated module drafts shared
the sole native writer. The physical run took 34.42 seconds; its unsettled
state was recorded instead of treated as cloth acceptance.

`FUMO-REVIEW-SAVE-GATE`: one review launched after a failed B build, then a
second encountered its already-created output directory. Both failed before
rendering. Subsequent launches check successful save status and use a fresh
packet path. This is a local orchestration correction, not a shared tool change.

`FUMO-CAMERA-CALIBRATION`: matching image scale and front landmarks had been
treated as sufficient camera evidence. The missing under-hem space exposed
the independent projection assumption. Future records must distinguish fixed
contract transforms from evidence that those transforms match a capture.
The partial contract read was corrected before publication; no acceptance
constraint was silently changed.

`FUMO-BOUNDED-READS`: a combined catalog/source response repeated structured
content and caused avoidable output truncation. Selected fields and bounded
file reads restored useful context. Existing delivery launcher reuse avoided
the previous session's cold rebuild. These findings stay task-local.

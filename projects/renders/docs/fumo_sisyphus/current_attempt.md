# Current attempt

[Back to goal](README.md)

## Attempt 01 — reference hierarchy reset

### Hypothesis

The scene can make measurable progress without touching the changing Fumo if
its macro composition is rebuilt around the source landmarks: a dominant
upper-left boulder, a fractured diagonal rock face with readable midtones, a
pale amber sky, and directional cast shadows.

### Boundary

- In scope: slope, boulder, background, camera, light, reference-at-side, and
  a replaceable neutral placeholder envelope.
- Out of scope: Reimu geometry, rig, materials, pose, integration, typography,
  and promotion into either tracked blend.

### Falsification gates

Reject the candidate if any condition holds in the first review frame:

- boulder width is outside `0.40–0.55` of frame width;
- the lower/right rock reads as a flat black wedge or leaves no readable
  surface relief;
- the slope direction is not clearly lower-left to upper-right;
- no visible contact/cast shadow anchors the boulder to the incline;
- the sky becomes dark brown or loses the source's pale negative space;
- a Reimu/Fumo object from the source remains render-visible.

### Result and process review

V01 and V02 were rendered and reviewed. V01 has rounded, textured masses and a
dominant upper-left boulder, but fills the frame with brown rock and background
instead of preserving the source's pale sky and readable diagonal face. V02
replaces the right side with conspicuous low-poly slabs while retaining the
same dark negative space. Both fail the macro hierarchy gate and are rejected.

The useful retained evidence is that boulder scale can be reached without a
black wedge. The next attempt must block the sky/slope/right-face composition
before procedural texture or additional rocks.

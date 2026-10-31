## Context

See `proposal.md`. Available host FFmpeg builds differ in their usable encoders.
The repository already owns a Video2X toolchain, diagram rendering, and Blender
procedures, but no component owns general video processing.

## Goals / Non-Goals

Document reliable media operations and verification without introducing a new
media implementation, dependency, or host configuration change.

## Decisions

- Use `tools/agents/skills/repo-video` for the cross-product procedure rather
  than assigning generic video work to TruffleHog or the Video2X dependency.
- Keep common decisions in `SKILL.md` and a tested still-image recipe in a
  conditional reference. Reuse existing tools instead of adding a wrapper.
- Prefer copying compatible audio. Encode the finite still-image video first,
  then mux complete finite inputs without an audio duration cap. Verify decoded
  audio for stream-copy exports rather than relying on metadata duration.
- Prefer a broadly compatible H.264/AAC MP4 for ordinary video delivery. When
  host H.264 is a stub, use the already available pinned browser encoder, as
  established in the earlier video thread, rather than silently choosing WebM.
  Document browser codec probing, timestamping, queue management, and decoding.
- Validate skill packaging and evaluation configuration offline. Exercise the
  recipe on generated audio and an existing diagram; omit a live model target
  because the offline harness cannot provide media tools and encoder state.

## Risks / Trade-offs

- A container fallback can violate a user's requirements; preserve explicit
  format constraints and describe any chosen fallback.
- Lossy conversions cannot retain identical decoded hashes; verify coverage and
  expected conversion separately from exact-copy checks.
- Codec and source timing vary; use observed capabilities and content evidence
  instead of making this session's encoder failure a universal rule.

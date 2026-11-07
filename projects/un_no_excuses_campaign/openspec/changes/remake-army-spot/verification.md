# Army spot verification

## Visual comparison

The revised scene is a stylized reconstruction, not an exact replica of the
source. Two implementation-blind reviews of intermediate candidates found
remaining differences in facial anatomy, uniforms, gestures, and staging.
Subsequent work adjusted the salute, the departing recruits, camera framing,
colors, shadows, garment proportions, and character visibility. Those later
changes were checked against representative source frames, but no independent
review has established an exact likeness pass.

Remaining visible limitations include simplified facial and hand anatomy,
rigid clothing and body motion, approximate lip synchronization, and simplified
background vehicles and buildings. Hidden rear anatomy is inferred. The source
video is not used as a visual layer in the scene.

## Scene verification

Pinned Blender clean-reopened the editable scene with automatic scripts
disabled. The scene has 452 objects, 11 cameras, 12 camera-cut markers,
299 animation actions, one packed sound, no linked libraries, no embedded
scripts, and no frame-change handlers. The timeline spans 898 frames at
30000/1001 fps. The packed sound bytes match the extracted source soundtrack.
The source attachment's SHA-256 remained unchanged.

The camera markers are at source frames 1, 79, 120, 198, 262, 350, 429, 648,
688, 717, 778, and 842. Scene inspection confirms the two woodland recruits
move away in the last shot. A visibility audit found and corrected stray facial
objects in officer-only shots, and the pair shots were isolated at their cuts.

## Media verification

All 898 PNG frames were present, newer than the final scene, and decoded at
960×540. The rendered video was decoded with pinned Blender at frames 1, 120,
301, 569, 648, 842, 898, and 902; each differed from its source PNG by a mean
of at most 1.1 RGB levels. The encoder held the last frame for four extra
frames to cover the full soundtrack. The MP4 has its metadata before media
data for fast start.

The decoded audio PCM SHA-256 is
`821d3b681271fe1195545dc8b13d89aea73c23ef1d763cc03d18f29ec2851d17`,
matching the extracted source soundtrack. A clean render of frame 301 from the
portable saved scene matched its full-sequence PNG exactly.

## Repository verification

The campaign targets, project documentation, and website built successfully;
the rendered landing and project index contain the new project and its stated
limitations. A 320×180, 24 fps encoder sample preserved the full resized image,
produced three frames from two input frames and one hold, and decoded through
pinned Blender. Repository validation passed 35 tests, semantic lint, and
strict OpenSpec validation. Local review identified no remaining code or
project-layout findings; the visual limitations above remain open.

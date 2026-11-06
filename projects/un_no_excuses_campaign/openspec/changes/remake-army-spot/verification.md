# Army spot verification

## Visual status

The replacement remains unfinished and has no user acceptance. The full MP4
and poster still show the rejected draft. Only selected stills and short motion
sequences have been rendered since the user stopped the full render.

One independent, read-only reviewer assessed the preceding character studies.
The findings identified disconnected-looking facial features, square uniform
forms, and overly frontal attention poses. The later candidate shades the
actual mouth cavities, revises facial deformation, rebuilds the main jacket
and collar, reshapes the officer's crown, and adjusts the confrontation and
salute. Buttons and ribbons are seated against the officer's jacket.

The implementing agent inspected the current close-ups, opening and wide-shot
contact sheets, and the dialogue and salute sequence. The independent review
does not certify this later candidate. Facial likeness, smaller recruits,
helmet placement, cloth shape, acting, word-level speech synchronization, and
scenery still differ materially from the reference. The source video is not
used as a background, image texture, or other visual layer.

## Native scene verification

Pinned Blender 5.2.1 LTS (`9e2066aef7ef`) clean-reopened the exact scene with
scripts disabled. It contains 437 objects, 13 cameras, 12 camera cuts, 413
animation actions, eight armatures, and one packed sound. There are no linked
libraries, image textures, or embedded scripts. The packed sound's bytes match
the extracted original soundtrack. The reference attachment's SHA-256 was
rechecked and remains unchanged.

The timeline has 898 frames at 30000/1001 fps. Cuts occur at frames 1, 79,
120, 198, 262, 350, 429, 648, 688, 717, 778, and 842. The maintained scene uses
960×540 CPU Cycles rendering with eight samples and denoising.

A full-timeline check followed by nonsequential frame access evaluated 7,240
leg chains. Maximum knee and ankle gaps were 0.000000308 and 0.000000716 scene
units, respectively, with no failures. Evaluated foot surfaces in representative
walking and standing poses remained above the ground plane. Opening frames
1, 23, 48, and 78 were also inspected visually.

All 100 salute frames keep the hand origin within 0.000000492 scene units of
the deformed cuff endpoint. The reaction hand remains within 0.000000388 scene
units over its 101 visible frames. Reaction-hand visibility passes checks on
both sides of the relevant cuts. These measurements establish attachment and
continuity; they do not establish visual likeness or natural motion.

The opening leg defect came from stale parent transforms and assumed segment
lengths. Its bake uses native bone lengths, explicit desired parent poses, and
independent shin scale inheritance. The revised wide stance also uses explicit
joint transforms. New head and collar channels have explicit values outside
their changed intervals, preventing extrapolation into other shots.

## Short-preview verification

The three maintained previews were rendered from the exact current scene at
480×270, six samples, and 30000/1001 fps. Each includes its corresponding source
audio excerpt; the character preview joins two excerpts at its internal cut.

| File                    | Source frames         | Encoded frames | Decoded samples | Maximum RGB error |
| ----------------------- | --------------------- | -------------: | --------------: | ----------------: |
| `opening_preview.mp4`   | 1–78                  |             78 |               6 |             1.385 |
| `wide_shot_preview.mp4` | 350–428               |             79 |               6 |             1.455 |
| `character_preview.mp4` | 289–318, then 548–607 |             90 |               7 |             1.579 |

Blender's bundled decoder verified dimensions, frame count, and frame rate.
Decoded samples were compared with the exact PNG inputs, including both sides
of the character preview's internal cut. Error is mean absolute RGB difference
on a 0–255 scale. All three MP4 files have fast-start metadata. The manifest owns
the exact scene, preview, and previous-draft file identities.

A separate 640×180 reference comparison joins the same dialogue and salute
excerpts. Seven decoded frames match its PNG inputs within 1.897 RGB levels.
It is review evidence, not the complete replacement video. The previous
full-render checks apply only to the old MP4 and poster, whose scene mismatch
remains explicit in the manifest.

## Repository validation

Publication requires the root `//:repo_quality_test`, the campaign OpenSpec
validator, semantic lint for the campaign packages and validation consumer,
and local review of the prepared aggregate. These gates are distinct from the
visual acceptance still needed for the replacement. Final full-frame coverage,
video, soundtrack, and exact-file verification remain pending in tasks.md.

## Implementation lessons

`LEG-PARENT-TRANSFORM`: calculate the child bone basis against the explicit
parent pose; setting a parent matrix does not guarantee its getter is current.

`MOUTH-CAVITY-SHADING`: blue interior walls produced the apparent raised lip.
Inspect cavity topology and material assignment before repeatedly changing
surrounding facial depth.

`GARMENT-SHOULDER-BINDING`: broadening a chest requires revisiting arm weights.
An X-only weight rule pulled the chest into the raised arm, and volume-preserving
skinning amplified the bulge. Linear skinning and separate chest influences
retain a controlled shoulder bend.

`NEW-CHANNEL-EXTRAPOLATION`: newly keyed channels need explicit original values
outside their changed intervals. Representative stills alone miss these leaks;
check contacts over every affected frame and inspect a short motion sequence.

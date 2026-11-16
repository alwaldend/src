---
title: Plush petting sketch comparison
---

`reference.jpg` retains the original operator-supplied photograph.
[`lift-sketch.png`](lift-sketch.png) is a prepared structural input with one
continuous head outline and a visible gap beneath the entire raised hand.
The sketch is not a Qwen output or a successful animation frame. It simplifies
some details and is not a pixel-aligned mask for the original photograph.

Use these inputs in order:

1. Upload `lift-sketch.png` and open `qwen-sketch-reconstruction`. Run the default
   full-strength reconstruction and inspect the raw PNG. Reject contact, a second
   hand, a second crown, disconnected anatomy or a different scene layout.
2. Once that baseline succeeds, upload `reference.jpg` and open
   `qwen-structure-guide-edit`. Image 1 remains the same sketch; image 2 adds the
   original appearance. Keep prompt, seed, negative text and sampling settings
   identical to the baseline. Compare `QwenStructureGuideRaw` with the baseline
   before looking at the optional composite.
3. If only the second run fails, investigate reference competition. If the
   baseline also fails, investigate sketch reconstruction before producing a
   sequence. A successful file or graph execution is not visual acceptance.

The checked-in graphs share the default prompt and sampler configuration.
Edits made separately after import do not synchronize automatically. The
baseline requires only the sketch; both retain native semantic and VAE image
conditioning on positive and negative branches. The appearance workflow also
restores original dimensions and produces a feathered composite using the
existing repair rectangle. Adjust that rectangle to cover the complete movement;
compositing cannot repair pose drift, differing outlines or changed framing.

## Sketch provenance and limitations

The built-in image generation editor prepared this input on 2026-11-16 using
the exact [asset preparation prompt](sketch-prompt.txt). Its input was the
white-background scene sketch from the rejected native-reference trial,
SHA-256 `81ce3fbee43e508b5caa510de445fe5c893eaf8355ec797ffb6c5567bef1d36f`.
That source combined native Canny edges from the public photograph with moved
hand contours and an added crown curve. Inspection found a conflicting retained
outer hair contour; the repaired asset removes that ambiguity. The trial record
retains its lineage and failed Qwen outputs.

The repaired asset is an opaque RGB PNG, 1297 by 1213 pixels, SHA-256
`c19aa520a4cc2ee4c809921bd639a3987bcbbb93927958ad71e6fa037610176c`.

Visual inspection confirms a single connected forearm/hand silhouette, one
rounded crown and a continuous gap. The [native Qwen trial](../../openspec/changes/add-parallel-evaluations/design.md#retained-gpu-trial)
records the subsequent GPU comparison: the default sketch-only reconstruction
retains the gap and coherent anatomy, while adding the original restores its
appearance but introduces unwanted crown folds. Preserve the two variants and
compare them at identical settings. Neither command success nor the improved
still establishes a finished animation.

The same trial also tests this targeted negative text with the original positive
prompt unchanged in both variants:

```text
animal ears, raised ear-like folds on the plush's head, horns, hats, cushions above the hair
```

It retains a clean sketch-only reconstruction but does not remove the folds
with the original appearance reference. It is a recorded experiment, not the
default or a verified correction. The saved workflows retain editable negative
prompt fields and the original blank default.

## Other poses and animation

The default prompt describes the supplied lifted-hand sketch. For a contact or
pressing guide, replace its hovering and air-gap instructions with the contact,
finger curvature and head compression shown in that guide. Keep image 1 as the
target geometry and image 2 as original appearance. An unchanged example prompt
can contradict a new sketch.

Approve the sketch motion before reconstructing its distinct poses. Inspect the
raw photographs before optional compositing or interpolation: the contact-pose
trial preserves hand placement but changes finger arrangement and portal
texture between frames. The existing rectangular composite can introduce seams.
Fixed seeds and shared appearance references do not guarantee temporal stability.

Photographic pose guides are another input to `qwen-pose-guide-edit`: prepare a
clean background and a separate hand layer, place the hand at each target pose,
then refine each guide with the original appearance reference. This preparation
is separate from Qwen sampling. In the [visible-motion trial](../../openspec/changes/add-parallel-evaluations/design.md#photographic-pose-guides-and-visible-motion),
the image editor prepares the layers and FFmpeg positions them. The earlier
drawn sketch is not an input to that trial.

The trial keeps seed 42, 40 Euler/simple steps, CFG 4 and blank negative text.
Start by comparing the lifted and pressed raw outputs at denoise 0.7; the
instruction must preserve either a gap or contact according to the current
guide. Lower denoise can retain more of the guide, including its defects: 0.45
removes one extra digit but fails to fix a separate skin-texture artifact.
Inspect all poses and omit rejected outputs before RIFE assembly. Larger
movement exposes the hand-layer warp's unnatural forearm shape, which Qwen can
preserve. The resulting loop demonstrates visible motion, not stable natural
anatomy.

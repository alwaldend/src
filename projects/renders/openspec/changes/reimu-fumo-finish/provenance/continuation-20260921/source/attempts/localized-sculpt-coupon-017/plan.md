# Localized sculpt coupon with repaired fixture setup

This is a technique-only coupon on disposable geometry. It does not open,
edit, or accept Reimu model bytes. One named writer owns Blender until the
coupon state is frozen; the coordinator independently audits and renders it.

## Repaired setup

Attempt 016 stopped before creating a mesh because its empty-scene helper used
a Blender 5.2-only render enum and assumed that `read_homefile(use_empty=True)`
creates a World datablock. This attempt removes both irrelevant dependencies:
the live authoring setup does not configure a render engine or access a World.
Pinned rendering configures its own scene only after candidate freeze.

This predicts that the same empty Blender 5.1.1 host reaches fixture creation
and a saved baseline. A setup correction is permitted only after naming a new
causal mechanism and predicted observable result. Repeating any settled setup
failure stops; changing unrelated parameters does not count as progress.

## Causal sculpt test

Attempt 015 proved broad Grab and proportional Edit transforms but could not
flatten pillow inflation, seat stem-like roots, or taper free edges. Build one
connected dense padded-panel fixture, separate from all Reimu files, with a
raised pillow center, narrow root beside a locked receiver, blunt tip, and
named plane, root, tip, mid-panel, and untouched-control regions. Deterministic
fixture construction is allowed because it is a measuring instrument. Freeze
and hash its exact bytes and coordinate baseline before sculpting.

Every claimed correction must use Blender's native interactive sculpt path,
not direct coordinate assignment, scripted replacement, or object-only
transforms.

1. Activate and read back essential Grab, Smooth, Scrape or Flatten, and
   Inflate or Clay brush types.
2. Set and read back scene-locked radius and unified strength at the planned
   values, including two strengths with a predicted magnitude response.
3. Establish mask or face-set isolation; non-target control vertices may move
   no more than `1e-6` scene units for each operation.
4. From the fixed front view, use native press, at least six timer-spaced move
   events, and release to reduce pillow height or plane variance by at least
   35 percent with Scrape or Flatten.
5. From the fixed three-quarter view, use controlled volume and Smooth strokes
   to increase root contact width by at least 25 percent and reduce root
   curvature discontinuity by at least 30 percent.
6. From the three-quarter view, use Grab and Smooth to taper the free tip to no
   more than 0.55 of mid-panel width and reduce tip roughness by at least
   25 percent.
7. Save the first complete coupon state to a new path, then undo the complete
   geometry sequence natively and require the exact baseline coordinate
   digest in the live scene.
8. Clean-open exact baseline and candidate bytes with pinned Blender 5.2.1,
   recompute the metrics and isolation, and render fixed front and
   three-quarter before/after views into a fresh immutable packet.

## Decision and stop rule

Pass only when all brush classes, settings, view paths, isolation checks, shape
metrics, exact undo, pinned clean-open, audit, and renders pass together.
Partial success fails. A passing result accepts only this authoring capability
and permits exactly one new whole-head attempt. A failing sculpt result blocks
autonomous Reimu variants pending a skilled Blender artist. Generic desktop
control is not a remedy unless a specific UI-only transport failure is first
demonstrated.

# Physical garment drape study 117

Outcome remains open and execution active under the user's request to
continue. Study 116 is closed rejected. Root owns this plan and is the sole
native writer using the unchanged pinned Blender mesh-data route.

## Decision and bounded plan

**Revise.** Two causal repairs to the explicit garment panels changed their
outline but did not produce cloth construction. An independent reviewer found
the final grounded result less convincing, not a successful stage. Replacing
its broad folds with more formula adjustments would repeat the same failed
mechanism. Test a rectangular cloth pattern gathered at the waist and draped
under gravity against the seated body and feet instead.

Protected input: assembly116 C, SHA-256
`012713486a2b2c9ce15da049005187b89fa178b36687dcc7179946e73ff2d829`.
Keep its head, bow, face and sleeve context unchanged for this isolated
garment mechanism test. Reuse the calibrated scale and fixed cameras. No
final material, fiber, rigging or acceptance work is authorized by this study.

The hypothesis is that a cloth solver responding to actual seat/foot contacts
will form broad nonuniform gathers and a compressed skirt, without a rigid
lap shelf or suspended rear edge. Pin and gather only the waist, preserve the
free lower perimeter, and simulate a bounded moderate-resolution mesh on CPU
with four threads. Retain the evaluated result in a new candidate rather than
depending on a mutable simulation cache for review. Add fabric thickness and
the white hem only after draping; do not hide contact errors with trim.

One physical drape candidate is permitted, with at most two causal setup
repairs for unsupported API, solver binding or invalid collision setup. Record
each before execution; no tuning sweep. Stop on explosive motion, penetrating
or missing cloth, or exhaustion of a 60-frame simulation. Protected files must
remain unchanged. Clean-reopen and inspect the saved whole-figure front and
side plus three-quarter pixels. Cloth appearance and reference likeness still
require independent review; technical simulation success alone passes nothing.

## Execution

The script used the official [ClothSettings API](https://docs.blender.org/api/5.2/bpy.types.ClothSettings.html)
and [ClothCollisionSettings API](https://docs.blender.org/api/5.2/bpy.types.ClothCollisionSettings.html),
then inspected the pinned runtime RNA before opening C. Runtime build remained
Blender 5.2.1 LTS `9e2066aef7ef`. The preflight passed; no native setup repair
was needed. Cloth contact and self-contact margins were 1 mm. The material
settings are an explicit uncalibrated hypothesis, not measured cotton values.

The 96-by-20 cylindrical pattern retained its Basis rest shape while a shape
key moved only the waist pins. Sixty frames completed in 7.84 seconds without
nonfinite coordinates or escape from the scene bounds. The result remained
in motion: maximum per-frame displacement during the last ten frames was
0.982 mm. Do not describe it as settled cloth. The source candidate remained
unchanged. Simulation and collision modifiers were removed after copying the
evaluated coordinates; the saved result has no cloth-cache dependency.

Candidate `assembly_117_drape.blend` saved and clean-reopened with SHA-256
`e03a85f42090ca5f3c48b81db7be1c2940e84747c334d37e27deecb5175cc074`.
The simulation script is
`454031e5546057307d57b6027f3734476a6984080bdc6af34f4aafaff2cd08db`.
Local evidence is in `out/reimu_fumo_finish/assembly_117/`: preflight, rest
pattern, frame records, frozen coordinates, build receipt, measurements,
aligned front comparison and overlay, and the four-view `review_a/` packet.
These are local working artifacts, not a published accepted asset.

## Result and next action

**Reject.** Root's front/side/rear/three-quarter inspection shows broad folds,
but a short gathered skirt perched over the exposed seat. Evaluated cloth
width is 105.76 mm versus the reference hem target of 140.95 mm; depth is
123.07 mm versus the 131.02–152.20 mm supporting band. The lowest cloth is
12.80 mm above the floor. These bounds do not replace a visible-silhouette
or penetration audit, and they already disprove the intended silhouette.

Implementation-blind reviewer `/root/cloth_117_blind`, supplied only the
references and frozen pixels, independently rejects it: likeness 5.5,
proportions 6, construction 4, identity 7, contacts 5 and cloth read 4 out of 10. Its largest failures are the exposed red seat, short lumpy skirt shelf,
narrow white piping instead of a gathered hem, excessive exposed thighs, and
rigid head/hair. Root agrees. No module, whole-character stage or final asset
criterion passes. The bounded study is closed rejected; execution remains
active with no external blocker.

The next construction decision must reconcile the reference's garment pattern,
seated-body volume and lower hem before another simulation. The previous
hand-shaped surface and this uncalibrated cylindrical pattern are not accepted
starting solutions. Preserve the fast physical solver as a tested tool, but do
not extend this run or sweep its settings. Inspect the supporting body's
actual reference silhouette and choose a new measurable pattern/contact
hypothesis first. Hair, bow and sleeves remain unresolved and receive no
inherited acceptance from the garment experiment.

## Session ergonomics

Two bounded module writers, one independent construction critique and four
frozen-pixel reviews supported root's sole native writer. All four new full
assemblies were clean-reopened and judged; rejection was not relabeled as
progress in overall likeness. Large batched source reads caused two avoidable
truncations; later inspection used selected records. A stale delivery launcher
failed before operating on Git; rebuilding its missing runfiles took 174
seconds. This is recorded setup cost, not a reason to bypass delivery.

The cloth solver itself took under eight seconds. For the next study, establish
pattern and support-volume evidence before authoring a generalized setup: the
long preparation did not compensate for the incorrect garment envelope.
Keep this finding task-local; no shared policy or infrastructure was changed.

# NoobAI sequence acceptance

## Initial still-image acceptance

The native graph/App pair has six nodes, eight connections and seventeen exposed
App controls. The App differs only in `extra.linearMode`; its `.app.json` suffix
matches frontend v1.57.6 bundled with ComfyUI v0.44.0. Upstream `workflowSchema.ts`,
`appModeStore.ts`, `comfyWorkflow.ts`, `workflowService.ts` and `uploadImage.ts`
at that frontend tag establish the serialized mode, selected widgets, output nodes
and native media-upload contract. No new browser extension is installed.

The offline native-call fixture exercises two scenes with four actions. Observed
seeds are 42–45; denoise is 1, 0.45, 1, 0.45. Edits consume the preceding decoded
image, while new scenes consume empty latents. Range continuation keeps the
absolute seed and uses the supplied preceding image. All actions share the patched
reference model. Invalid scenes/ranges/controls fail before sampling. Cancellation
preserves completed files, and portrait reference edges survive square preparation.
The fixture emits real PNGs and `image-sequence-result.json`; GPU boundaries are
deterministic substitutes, so this does not establish model execution or fidelity.

The configured-profile Terraform CLI fixture packages every model profile using
loopback resources only. NoobAI retains all ten generated graph files, with its
native App selected for startup. Its encoded startup is 15,802 bytes against the
unchanged 16,384-byte limit. Exact first-party source, staged workflow content and
model destinations are verified. `NOOBAI-APP-BASENAME` and `NOOBAI-STARTUP-SIZE` are
resolved by the bounded native suffix and existing startup selection respectively.

Repeatable checks:

```sh
bazel test //projects/comfyui_tools/nodes/image_sequence:image_sequence_test
bazel test //projects/comfyui_tools/nodes/prompt_sequence:prompt_sequence_test
bazel test //projects/comfyui_tools/workflows:format_test
bazel test //projects/vastai_terraform_provider/test/acceptance:acceptance_test
```

Task-local artifacts under `out/noobai-image-sequence/` include the generated graph,
App, source/import ZIP and `artifact-verification.json`. That report checks reciprocal
links, selected widget names against the actual source contract, seed-widget layout
and graph/App parity. The selected infrastructure profile remains MOSS-Image; no
live rental, inference job, deployment or user image was accessed for this change.

## Session ergonomics

`NOOBAI-CONTEXT-BOUND`: broad skill/catalog and query reads produced unnecessary
output during discovery. Subsequent checks use exact owner files, named targets,
bounded diagnostics and task-local logs. `NOOBAI-BAZEL-OUTPUT`: tool invocations can
change the convenience `bazel-bin` symlink; artifact inspection uses the actual
generated configuration path and copies deliverables into task scratch. No shared
instruction, cache or host configuration was changed. Delivery receipts own exact
candidate and publication state; the local review record owns the final verdict.

The prepared implementation passed all 38 executed repository-quality and affected
behavior/packaging checks, the corrected aggregate test command, strict OpenSpec
validation and affected-target semantic lint. Local review found no actionable
defects. Final task-state edits are documentation only and remain subject to the
same preparation and publication gates.

## GIF and default-demonstration extension

The GIF-only extension had eighteen controls, including image hold time for looping GIFs.
Its default has six images split into three snowy-park poses and three library
poses. The native-call fixture writes and decodes actual scene/combined GIFs,
checking frame order, dimensions, 500-millisecond holds, looping, single-image
partial runs, repeated-image playback time and preservation of a completed scene
on later interruption. All six behavior cases pass. The fixture still substitutes
GPU model boundaries and does not establish artistic quality.

`NOOBAI-GIF-STARTUP-SIZE` is resolved through the existing post-boot workflow API
staging pattern. The GIF-enabled node is installed declaratively; generated workflow
documents are saved afterward through the owning authenticated command. NoobAI's
encoded startup is 15,230 bytes with no inline workflow files and the complete
ten-document importable catalog preserved. All configured startup profiles pass.

The operator explicitly authorizes the existing 96 GB profile and requests a live
default demonstration, retained rental and endpoint. Deployment and live output
observations will be recorded under the infrastructure owner's change after the
reviewed apply and generation complete.

## Canny, interpolation and live acceptance

The final graph/App has seven nodes, nine connections and 21 exposed controls.
It retains the six-image, two-scene example in concise NoobAI tags. The App uses
change strength 0.6, Canny strength 0.2 until 0.75, reference influence 0.05 and
4x RIFE interpolation. Generic node defaults remain compatible with older graphs.
All nine offline behavior cases pass, including native conditioning selection,
scene resets, partial continuation and decoded interpolation timing/scene cuts.

The updated startup fixture verifies the approved ControlNet destination and
deferred source omission with retained whitelist membership. NoobAI's final
startup is 8,402 bytes; the owning Ansible workflow installs the expanded node,
App and pinned interpolation engine after boot. The packaged portable AppImage
is an ELF file matching the existing lock's SHA-256 integrity.

The retained 96 GB deployment completed the default sequence with native NoobAI,
IP-Adapter, Canny and RIFE. Final owned job
`2af51190-554b-45c4-bb90-503eaedf03eb` completed in 54.7 seconds, saving six
1024-square PNGs, two nine-frame/1.5-second GIFs and an eighteen-frame/three-second
combined GIF. Decoding verifies loop, dimensions, timing, absolute seeds 42–47,
scene-reset denoise and effective Canny settings. The exact App was uploaded.
Original media, metadata, hashes and native CLI receipts remain in task scratch.

Visual acceptance is limited: layouts and broad appearance remain stable within
a scene, while the requested hand and book poses change weakly. Fresh scenes
can alter character details. Canny does not supply pose editing or identity
locking. Strong generic SDXL IP-Adapter weights degraded NoobAI's colors and
anatomy, so those results were rejected and the App starts at low influence.
The default demonstrates the workflow and its stability/change tradeoff rather
than establishing reliable instruction-following or exact cross-scene identity.

Bounded refinement lineage: the original full run
`faae95b6-0527-41bc-a924-f971c4534d0a` was rejected for reference overweight.
Three owned one-image comparisons tested 0, 0.15 and 0.3 influence. The next
full run `3bb73673-4345-43f8-a067-9b970c3a8cfc` retained clean scenes but strongly
restricted pose changes; the final revision loosened continuation and Canny.
No further generation is scheduled. Only submitted native IDs were inspected;
the installer uses aggregate idle counts without reading other jobs.

## Interpolation cancellation correction

Automated review identified a blocking interpolation wait that delayed native
cancellation. The corrected node polls interruption during bounded waits and
terminates its owned process group on interruption, timeout or failure. A real
parent/descendant fixture reproduces cancellation while a descendant ignores
termination and closes its output pipes. The failure-first check exposed that
case; the corrected check verifies both processes stop, the two completed PNGs
remain and run-owned temporary media disappears. Successful decoded playback
coverage remains in the same suite.

The in-place updater installed only the corrected node with workflow copies
skipped, using aggregate idle counts. The remote source hash matches the local
source and the updated native interface responds. Later live App edits and the
separate uploaded example are preserved. No additional GPU generation was run
for this correction; the existing inspected demonstration remains the visual
acceptance evidence.

## Separate reference and scene-composer Apps

The operator rejected the earlier robot demonstration as unacceptable. Its
preserved outputs are negative visual evidence, not acceptance of character
identity or action fidelity. The supplied Scenes Composer instead shares fixed
prompt components and a sampling seed across independently generated images,
with a second pass over each image's latent. It has no image-reference path.

The revised implementation exposes both strategies as separate graph/App pairs.
The reference version retains uploaded references and previous-image feedback;
the composer exposes character, clothing, style, scene settings, actions and an
optional locally installed LoRA. Its default is 25 first-pass steps and 20
refinement steps, with a fixed seed. Both versions now publish PNGs and exactly
one combined GIF; interpolation remains scene-local with hard cuts.

All thirteen offline behavior cases pass. The composer fixture executes eight
native sampler-boundary calls for four completed PNGs, decodes its four-frame
two-second GIF and records effective prompts, modes, seeds and refinement
settings in a repeatable JSON artifact. The reference fixture verifies four
calls, carry-forward pixels, scene resets, partial continuation and the same
single-GIF contract. Additional cases reject malformed settings/LoRAs and stop
between refinement passes while retaining completed PNGs. The real-process
RIFE cancellation case also passes after correcting a procfs exit/read race in
its fixture. These fixtures substitute GPU model operations; they establish
execution behavior, not NoobAI art quality.

Generated graph/App parity, reciprocal links and all widget bindings pass for
both variants: seven nodes/nine links/21 controls for the reference version and
four nodes/seven links/22 controls for the composer. Each has a six-image,
two-scene default. The composer uses only native operations and first-party
nodes; no external node pack, new model, or LoRA was installed. The import ZIP
contains both graph/App pairs and the exact native module.

Live upload is not completed. On 2026-11-15 at 08:43 UTC the old HTTPS endpoint
returned 502, SSH was refused, and the infrastructure owner's specific provider
lookup returned no instance. Its historical Terraform output still names the
old rental. No live node/workflow update or new generation was attempted after
those observations. The revised variants therefore have offline behavioral and
packaging evidence only; their NoobAI visual quality remains unverified.

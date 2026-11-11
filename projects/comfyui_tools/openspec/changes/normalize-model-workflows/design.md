## Context

Juggernaut XL v9 is an all-in-one SDXL checkpoint with bundled encoders and VAE.
The approved RES4LYF revision includes an SDXL style example and ReSDPatcher;
the archived Chroma graph uses ReChromaPatcher and cannot be reused unchanged.

## Decisions

Use <profile>-<purpose>.json everywhere: file, profile manifest and saved workflow.
Use full text-to-image and image-to-image purposes; keep single-reference,
character-style, prompt-sequence, style-reference and style-prompt-sequence
where those accurately describe the graph. The folders retain model ownership.
Preserve existing Flux, Qwen and Chroma JSON bytes during renaming.

Native Juggernaut graphs use CheckpointLoaderSimple, CLIPTextEncode, KSampler,
VAE decode and SaveImage. Image-to-image scales and center-crops the reference
to 1024 square and encodes it as the starting latent. Text-to-image uses fresh
noise. Prompt sequences use the existing first-party node and shared text.

For style-reference graphs, encode the reference with the checkpoint VAE and
route it through ClownGuide_Style_Beta to ClownsharKSampler_Beta. Patch the SDXL
model with ReSDPatcher in float32. Start from an empty SDXL latent rather than
the reference's composition. Expose guide weight, use a neutral illustrated
prompt, and keep a native negative prompt. The sequence variant reuses the same
guide for every panel. This is style guidance, not guaranteed style isolation or
character identity preservation.

Decision review: proceed with the previously approved pack and exact dependency
lock. Plain img2img carries composition and does not meet the reference-style
request. A new adapter and vision-model dependency are unnecessary for this
source option. Native-only graphs remain available if the pack fails on the
next rental. Do not claim inference or quality acceptance from offline checks.

## Validation

Before acceptance, check malformed/nonreciprocal links, mismatched slot types,
wrong checkpoint/VAE paths, a missing model patcher or guide edge, disconnected
sequence outputs, missing widget metadata, unpinned dependencies, loss of old
profile bytes and startup-payload constraints. Existing browser and Terraform
CLI fixtures cover editor duplication and provisioner behavior. Retain a graph
and dependency verification receipt under ignored task scratch. No rental or
private content is used by these checks.

## Jsonnet authoring

The operator approved go-jsonnet v0.22.0 as a separate build dependency. Shared
tooling is recorded by infra/src/add-jsonnet-tooling. Source definitions use one
graph library for node defaults and link bookkeeping, model-specific builders
for their conditioning/sampling contracts and shared sequence composition.
Generate native workflow JSON as Bazel outputs at existing model/name paths;
remove independently maintained JSON copies. Preserve every parsed field against
the public 14-document baseline at ec49a504, including editor positions, settings,
IDs and widget links. Libraries take model filenames as parameters; the generation
entrypoint may read the existing public rental configuration as its single owner.
Archived Chroma filenames remain local example defaults. No model/rental changes.

Failure cases: reordered widget values, missing input metadata, changed reference
conditioning, unstable IDs or link IDs, stale configuration filenames, lost layout,
undeclared imports, missing generated files in Terraform/Hugo packaging, changes
to prompts or sampler defaults, and requiring Jsonnet on the rented host.
Use actual compiler generation plus exact parsed-output comparison and existing
consumer fixtures. No new GPU/private-content acceptance is claimed.

## SDXL WCT precision correction

The operator reported `double != float` at RES4LYF StyleWCT.get and prohibited
rental inspection while using the instance. At the approved revision, SDXL style
sampling promotes the generated branch to float64 while Retrojector.embed returns
the caller's dtype; the configured float32 reference branch therefore reaches the
WCT matrix multiply with a different dtype. Set ReSDPatcher.style_dtype to the
upstream float64 default for both composed Juggernaut style graphs. This changes
style-projection math, not the checkpoint weight precision. No dependency pin or
model is changed by this source correction. Subsequent explicitly authorized
replacement-GPU tests completed both style paths without the dtype failure;
the owning infra/image_gen record holds deployment acceptance.

## Approved several-reference conditioning

The operator requested several references containing both characters and style,
and explicitly approved ComfyUI_IPAdapter_plus at
`a0f451a5113cf9becb0847b92884cb10cbdec0ef`, the general SDXL Plus adapter,
and its CLIP-ViT-H encoder. Model filenames and fixed publisher URLs belong to
infra/image_gen/config.json. The adapter file at revision
9bf28b38530e55ffa91c6d82e5161a982c22f284 has SHA256
3f5062b8400c94b7159665b21ba5c62acdcd7682262743d7f2aefedef00e6581;
the encoder at 0859e809306db97aa2338e370587ab284e8a754f has SHA256
6ca9667da1ca9e0b0f75e46bb030f7e011f44f86cbfb8d5a36590fcd7507b030.
They add approximately 3.4 GB. FaceID/InsightFace is outside the approved path.

Compose three resized reference images into one native image batch and apply
IPAdapterAdvanced with concat embeddings, linear weight 0.8 and the existing
native SDXL sampler. A sequence variant shares that conditioned model for each
nonempty line. This guides appearance/style rather than guaranteeing identity
or subject isolation. No additional editor abstraction is introduced.

Pinned backend node schemas match the generated input/output and widget
contracts. The complete seven-workflow profile passed the real offline
Terraform/provider lifecycle and encoded to 12,310 bytes, below 16,384. The
fixture compares every staged native graph with its generated source artifact.
This validates packaging and contracts, not GPU image quality. The upstream
plugin is maintenance-only; future ComfyUI upgrades need compatibility review.

## Shared sequence negative prompt

The operator requested all prompt text in Prompt Sequence and removal of the
separate negative-prompt editor. Append an optional multiline negative_prompt
input and a third scalar STRING output. Preserve the existing widget positions,
prompt output slot 0 and filename output slot 1 so older saved workflows remain
usable. Negative text is shared across panels and must not inherit the positive
context or be split into a panel list. Omitted input emits an empty string.

The shared Libsonnet sequence transformation moves the base workflow's existing
negative text into the sequence widget and connects its output to the native
negative encoder's converted text input. The connected input replaces the
standalone editable text box; retain the encoder because samplers consume
conditioning, including Qwen's reference-conditioned negative branch. Preserve
all non-sequence graphs, model-specific conditioning, guidance and sampler defaults.

Before implementation, check nonempty and multiline shared negatives, empty and
omitted input, old widget/output positions, blank panel lines and paired filenames.
Verify all five generated sequence variants, reciprocal typed links and native
CLI conversion without any rental connection. The operator is using the rental;
this source-only change does not authorize inspection, jobs or deployment.

The operator subsequently requested destruction and deployment of the identified
FLUX.2 dev 32B model, selecting the prepared 96 GB hardware option. Add native
text-to-image, three-reference and shared-reference sequence graphs under
workflows/flux2-dev. Use BasicGuider, FluxGuidance 4, Euler, Flux2Scheduler and
fresh 1024-square latents; start at 28 steps. Reference VAE latents chain through
native ReferenceLatent nodes. This guidance-distilled path has no negative
branch, so the shared sequence transformation must also accept no negative
encoder rather than adding an ineffective conditioning node. Existing Base,
Qwen and Juggernaut sequence graphs retain their shared negative connections.

The five negative-conditioned sequences passed native CLI UI-to-API conversion
using public node schemas and the repository-pinned Python/CLI dependency
closure, with three ordered panels each. Qwen uses its actual prompt input name
for both converted encoders; other families use text. The new dev sequence
retains only positive/prefix connections. All six relevant package checks passed,
including the complete selected profile through the Terraform/provider fixture.
The operator began using the replacement before GPU acceptance and prohibited
further access, so no inference acceptance is claimed for these new graphs.

The operator identified the missing single-reference path after starting use.
Compose a single-reference variant by retaining the first native reference chain,
with a matching sequence that shares that image. Add generated JSON to the
declared future provisioning bundle and provide the operator importable artifacts.
Do not install or inspect anything on the active rental.

## Full Chroma1-HD profile

The operator requests the full, previously approved Chroma1-HD checkpoint on
the earlier RTX PRO 5000 48 GB rental and explicitly authorizes generation tests.
Reuse the Chroma builder with model filenames and sampling parameters; retain
the three archived Flash graphs unchanged. Add five chroma-hd outputs: native
text-to-image, img2img, RES4LYF style-reference, prompt-sequence and
style-prompt-sequence. Use the shared sequence transformation so the positive
prompts, paired filenames and scalar negative prompt remain connected.

The publisher's full-model example uses CFG 3 and 40 steps. Its native workflow
uses Euler and T5TokenizerOptions [0, 0]. Use those contracts for native HD
graphs with beta scheduling; retain RES 2M/beta57 for the approved style path.
Use the full 16-bit checkpoint and FP16 T5 with the original VAE, not quantized
substitutes from the publisher's optional workflow example. Native img2img
carries composition; RES4LYF guides style without promising character identity.
Check all five workflows on the separately authorized rental with neutral
fixtures, including blank-spaced panel lines and shared negatives.

## Chained frames and native queue visibility

Use the existing Prompt Sequence parser as the sole owner of panel separation.
The operator prioritizes the existing media grouping over separate queue jobs.
Keep the native list mapping within one job and group. Append an optional model
input/output to Prompt Sequence, clone its patcher and wrap OUTER_SAMPLE to emit
native progress text such as Panel 2 of 3, leaving sampler results unchanged.
Use the shared graph transformation to place that connection before each native
sampler or guider. Do not add a browser extension or an independent queue.
Allocate one UUID per actual node execution and include it in every panel prefix.
Mark the node non-cacheable so identical queued runs receive new identities and
fresh counters; retain unchanged model weights and conditioning caches upstream.

The chained img2img GIF is one output node using the installed native sampler,
VAE, Torch and Pillow. It consumes the ordered conditioning list, samples from
the previous frame's latent, decodes each generated frame and writes one GIF.
It does not create independent queue entries because frames depend on each
other. Native sampler progress remains available within that job. Start at
denoise 0.45 and 6 fps; prompt lines are discrete frames, not interpolated motion.

Before implementation, cover: blank/JSON/paragraph parsing mismatches, unsafe
prefixes, reused run identities, wrong panel filenames, counter state leaking
between jobs, mutating the original model, changing sampler return values,
missing model connections and splitting a media group. For the GIF cover:
resetting the latent instead of feeding back, prompt/negative ordering, invalid
batch/list lengths, inconsistent dimensions, unsafe paths, incorrect GIF frame
count/timing/looping, interruption leaving partial output and filename collisions.
Use deterministic native-sampler boundaries for isolated behavior checks and
real neutral GPU workflows for acceptance. No new dependency is introduced.

The controlled full Chroma style comparison favored guide weight 0.5 with
end_step 15 over weight 1 through all steps. Parameterize those HD defaults while
preserving archived Flash settings; this improves the tested sharpness but does
not establish general style fidelity.

The follow-up temporary GPU acceptance completed all six HD workflows, a
repeated sequence and a zero-denoise identical-frame GIF control. Native
progress messages count all panels, each sequence remains one job with three
outputs, and repeated runs use distinct UUID prefixes. Both GIFs retain three
1024-square frames at 170 ms, including identical frames. The owning
infra/image_gen change records timings, visual limitations and immediate
teardown. No paid rental remains after acceptance.

---
title: Image-generation evaluations
description: Bounded parallel ComfyUI evaluation configuration and cleanup
---

![Evaluation lifecycle](../assets/evaluations.dark.webp)

Default cleanup is shown; `--cleanup=false` retains rentals and skips the final sweep.

This subproject instantiates the reusable evaluator in projects/comfyui_tools.
It owns the run matrix, credential references and deployment targets. All
executable code, cases and references stay with that project. Personal rental
resources and Terraform state remain independent.

The matrix in config.json compares Juggernaut, NoobAI, Illustrious, Chroma and
HiDream on the approved thorn-landscape case, with text-only and reference-guided
workflows. Chroma, HiDream and Qwen also run the plush-petting GIF case. Two fixed seeds
are configured for still images. Hardware and publisher pins come from the parent
config.json. Configuration is not evidence of GPU quality.

Compatible matrix groups cross hardware_profiles, model_profiles, workflow
purposes and case_ids, then apply shared seeds and optional sampling controls.
A purpose such as text-to-image resolves to <model-profile>-text-to-image.json.
Separate groups express compatibility without a second workflow catalog. Unknown
files, unsupported bindings, duplicate axis choices/seeds and expansions above
128 workloads fail locally before allocation. The current matrix expands to
eighteen workflow/case/settings combinations and twenty-eight seeded queue entries, pooled
onto six hardware/model rentals. Add another hardware profile to a group to
compare that rental configuration too.

All workflows, cases, settings and seeds sharing a hardware/model pair reuse one
rental, including across matrix groups. Each worker allocates only when a rental
concurrency slot is available, waits for its own ComfyUI endpoint and immediately
starts submitting native queue entries. A bounded lookahead overlaps submission
with completion waits and downloads, so later jobs can render while earlier
outputs are retrieved. After all its assigned outputs are downloaded, the worker
destroys that rental before local media inspection. If another eval worker
holds the available capacity, an unallocated worker waits and can reuse that
offer after confirmed release. Failed or uncertain reservations stay excluded. A failed workload cancels
further submissions and releases the rental. Confirmed provider exit, removal or
ownership change also cancels active work promptly; transient control-plane
observation errors retain the existing workload deadline. In failure cases, already downloaded artifacts remain.

Maximum parallelism is three rentals. Each has a 45-minute lifetime for its
entire assigned workload, a five-minute readiness deadline and a separate cleanup deadline. The standard hardware profile selects any GPU with at least 32 GiB VRAM,
100 GB disk and a $1/hour reported offer-price ceiling. Offers are ordered by
lowest reported hourly price among offers satisfying both transfer-rate caps;
no GPU model is pinned. The canonical shared offer declaration caps incoming
and outgoing bandwidth at USD 0.025/GB each. This rejects the expensive-transfer
hosts used in the first experiment without disabling any host requirements. Datacenter and verification are
mandatory in both the query declaration and returned tier metadata. The explicit
minimum_reliability value of zero removes only the inherited 99-percent floor;
architecture, CUDA and port requirements remain. The hardware profile caps the
reported hourly offer price. Network fees and delayed
provider operations mean this is not a guaranteed account spend cap. No timer or
persistent service is deployed.

With automatic cleanup enabled (the default), the final step and standalone cleanup use the same rule: destroy every instance
whose label starts with `imagegen-eval:image_gen:`, regardless of run or age.
They preserve normal rental labels and verify that no matching instances remain
by enumerating all pages of the current Vast v1 instances API.
Run one controller per namespace. A shared per-OS-user cache lock excludes local
controllers with the same namespace even when their output directories differ;
coordinate this operationally across machines or OS accounts. Manual cleanup deliberately terminates matching test rentals too; stop an active
controller before manual cleanup so it cannot allocate more queued jobs. Forced termination or host failure can skip the final
step; invoke the standalone cleanup target to recover.

The existing Vault-owned team key and web login are injected through the parent
component. Only the web login reaches ComfyUI access. The runner never writes
Vault or reads the personal Terraform backend. Startup uses the existing
hardening template, compiled by the pinned Terraform expression evaluator in an
empty directory. Terraform manages no evaluation resources or state.

Results, journals, effective inputs and credential-masked CLI failure diagnostics
are written under ignored
out/imagegen/evaluations/results. Automatic checks cover lifecycle, decodable
media, frame count, timing, looping metadata, distinct frame colors and boundary
pixel difference. Pixel change does not establish a correct petting gesture. Subject adherence, style, reference
following, natural motion and artistic quality have explicit review criteria
and remain unscored until reviewed. Live runs require an explicit operator
invocation; offline checks never rent capacity.

## Commands

```sh
bazel run //infra/image_gen/evals:check
bazel run //infra/image_gen/evals:run
bazel run //infra/image_gen/evals:gif
bazel run //infra/image_gen/evals:cleanup
```

For a diagnosed capacity failure, `run -- --group GROUP --exclude-offer ID`
selects only the requested group and skips that offer for this invocation.
For example, append those arguments to `bazel run //infra/image_gen/evals:run`.
The IDs must be unique and positive; mandatory host filters remain enforced.
The final report records exclusions. Cleanup still destroys the full namespace.

Only check is offline. Run rents the configured matrix; gif selects only the
three plush groups (Chroma, HiDream and Qwen), producing eight experimental
clips on three reused rentals. Cleanup destroys all
matching evaluation instances. Both live targets require an explicit operator
invocation. Run supplies its final cleanup automatically, including failures
and ordinary interruption. Append `-- --cleanup=false` to the run target to
retain all allocated rentals after success or failure instead. Both per-rental
destruction and the final sweep are skipped, and run.json records
cleanup_skipped=true and cleanup_confirmed=false. The journal identifies each
retained rental; invoke the standalone cleanup target when debugging is finished. Source validation builds the live targets without
executing them.

Edit config.json to select matrix groups, profiles, cases, workflow purposes,
seeds and limits. Optional per-group sampling controls include steps, cfg, denoise,
adapter_weight, adapter_end and style_weight. Omitted controls retain the
selected workflow defaults; overrides must apply to every workflow in their
group. Effective workflows and expanded rental/workload identities are saved
with the results. queue_and_generation_seconds includes time spent in the
native queue and CLI polling overhead; retrieval_seconds measures downloads.
completion_observed_at is captured by an independent bounded CLI observer for each
submitted trial, so preceding downloads do not inflate later completion observations.
It is not an authoritative GPU execution timestamp.
Each trial has a forty-minute absolute budget, bounded by the rental's 45-minute lifetime. The trial retains that deadline across upload, submission, native queue
residence, generation and download. Frame lists use the node's JSON mode so
multiline entries remain one frame; empty entries fail before uploading.
GIF controls `frame_source` and `frame_seed_mode` select original-reference
sampling or feedback and fixed noise or incrementing seeds. The selected
comparisons use fixed noise with feedback at denoise 0.45, original-reference
sampling at denoise 0.65, and a masked version of that original-reference path.
Qwen compares full-image and masked edits at its native denoise 1. The case uses
static hand-pose instructions with identical first/last prompts. Its edit_region
owns the hand rectangle and soft boundary. These are experiments, not a motion
quality guarantee; the previous two-clip baseline is recorded below.
The masked and Qwen GIF graphs are generated repository assets submitted by the
evaluator. They are not added to the personal rental's startup workflow bundle,
which must stay inside Vast's compressed startup-size limit. Load the generated
JSON manually when using these experimental graphs outside the evaluator.
The first matrix compares model/workflow combinations with their model-specific
starting settings, not isolated model architectures under an identical sampler.

```sh
bazel test //infra/image_gen/evals:configuration_test \
  //infra/image_gen/evals:workload_e2e_test \
  //projects/comfyui_tools/test/evaluation:evaluation_test
```

The first check validates the actual packaged matrix and canonical startup
template using public dummy login values and no API client. The lifecycle
fixture uses a loopback Vast server and verifies matrix pooling, offer reuse,
namespace locking, early rejection and cleanup; the workload fixture uses the real pinned
CLI through an isolated HTTPS ComfyUI API, including queued image/GIF retrieval
and cancellation of a stalled CLI descendant. These checks cannot evaluate image
quality or establish new GPU/model compatibility.

## GIF acceptance

The authorized 2026-11-11 run completed both configured GIF workloads on verified
datacenter CMP 170HX 40 GiB rentals. Both outputs decode as eight 1024-by-960
frames with 170 ms delays and infinite looping. Chroma's native queue and
generation took 278 seconds; HiDream's took 410 seconds. These timings include
model loading after submission. Each rental was destroyed after its own
download, and the final paginated sweep confirmed the eval namespace empty.

Media validity passed; convincing petting motion did not. Chroma preserves the
scene with little readable hand motion. HiDream changes more but drifts in
the plush's appearance and the background. Results remain in ignored local
output, with the detailed evidence in the owning OpenSpec change. The evaluator
recorded successful execution and cleanup; the wrapping Vault plugin separately
reported a token-revocation error during shutdown, so the overall command exited
nonzero. No successful whole-matrix or still-image quality claim follows.

## GIF experiment outcome

The subsequent eight-variant comparison retrieved six Chroma/HiDream clips.
Fixed-noise feedback still accumulates severe noise or appearance drift.
Original-reference sampling prevents that accumulation and repeats the initial
frame at the loop boundary, but Chroma barely moves the hand. The masked HiDream
variant produces more localized movement while retaining malformed fingers.
None passes the case's natural-petting criterion.

Qwen produced no clip: its first rental missed readiness, and a different
verified datacenter offer reached readiness but returned a generation execution
error. The retained CLI summary does not identify the model exception; no cause
is asserted. These Qwen graphs remain experimental, without live GPU acceptance.
All four rentals across the initial comparison and isolated retry were destroyed;
final account enumeration found no instances. Billing and downloaded outputs
remain local, with failed attempts included in the operator's cost report.

Qwen's generated GIF workflows now encode each previous decoded image through
the native image-edit encoder for both conditioning branches, and encode those
same pixels for sampling. Their new node must be loaded before these graphs
are used. On the retained rental, an explicitly authorized ComfyUI restart
loads the updated first-party source; saving workflow JSON alone does not.

For an explicitly authorized retained-rental update, the packaged SSH client is
`bazel run //infra/image_gen/evals:ssh -- -p PORT root@HOST COMMAND`. Supply the
current instance's provider-reported endpoint and existing account SSH identity.
The target keeps batch authentication and strict host-key checking enabled;
bootstrap a task-owned host-key pin before first use. It never selects or
discovers user jobs.

The later native-decoder retry retrieved both Qwen clips. After the operator
authorized a service update, both corrected decoded-image feedback variants
also completed, producing eight-frame GIFs. Full-image feedback still drifts
and deforms fingers; the masked variant remains almost static. Neither passes
the petting criterion. The operator subsequently requested destruction, and
account-wide enumeration confirmed zero rentals. The selected NoobAI personal
configuration is source-only and has not been rented.

## 1. Implementation

- [x] 1.1 Implement the declared behavior and owning documentation.
- [x] 1.2 Verify offline behavior and preserve repeatable evidence.

## 2. Delivery

- [x] 2.1 Complete formatting, semantic lint, repository quality, and local review.
- [x] 2.2 Prepare the checked candidate for guarded publication; the delivery receipt owns publication state.

## 3. Dedicated team credential

- [x] 3.1 Add an independent alwaldend API-key creation and regular-user Vault-storage root, with offline checks; remove the shared-key policy grant.
- [x] 3.2 Create and store the credential through an interactive reviewed apply after the Vault session is renewed; verify a no-change plan without renting a GPU.

## 4. Component ownership

- [x] 4.1 Consume the image resolver from tools/vastai and the dedicated key from infra/vastai; retain only rental and credential-read configuration here.

## 5. Default checkpoint

- [x] 5.1 Configure the explicitly approved Juggernaut XL v9 publisher checkpoint in the rental's Terraform defaults and document overrides.
- [x] 5.2 Verify the fixed-revision public URL, existing download-manifest behavior, and offline rental checks without renting capacity.

## 6. Direct image tag and authorized rental

- [x] 6.1 Replace dynamic image discovery with a direct published tag and remove unused provider/tool wiring.
- [x] 6.2 Review security/privacy and validate the deployment candidate before creating capacity.
- [x] 6.3 Rent one matching instance through the owning interactive apply, verify its observed HTTPS endpoint and authentication, and report its URL and Vault login reference.

## 7. Chroma replacement

- [x] 7.1 Configure Chroma1-HD and its publisher-recommended T5 encoder and VAE, selecting them in Forge.
- [x] 7.2 Validate and review the candidate, then replace the rental through the owning interactive apply. The user explicitly confirmed that no local files need preservation.
- [x] 7.3 Verify HTTPS authentication and model availability on the replacement; report the new endpoint and retained Vault login reference. The user requested immediate deployment and will evaluate generation interactively.

## 8. ComfyUI and Chroma Flash source update

- [x] 8.1 Replace the Forge consumer with the ComfyUI module and full Chroma1-HD-Flash checkpoint, retaining T5, VAE and Vault ownership.
- [x] 8.2 Publish the concise dark architecture diagram with correct actor, Vault retrieval, Vast instance and native model provisioner boundaries.
- [x] 8.3 Validate and prepare the source for guarded PR publication without creating another rental; the delivery receipt owns publication state.

## 9. Ready workflows and authorized ComfyUI deployment

- [x] 9.1 Supply native text-to-image and reference-image img2img workflows with the declared Chroma Flash, T5 and VAE filenames.
- [x] 9.2 Review the fresh interactive plan and rent the explicitly requested instance, retaining login and host constraints.
- [x] 9.3 Verify trusted HTTPS, authentication, native model/workflow availability and a representative inference result.
- [x] 9.4 Rebase after deployment, complete source validation/review and prepare guarded publication of the workflow follow-up.

## 10. Reference-style behavior correction

- [x] 10.1 Reproduce the ignored-reference failure with fixed prompt/seed and guide-on/off comparisons; inspect the approved pack's required model patch.
- [x] 10.2 Correct the native workflow using the existing approved node pack and update its deployed user document without recreating the rental or disturbing user files.
- [x] 10.3 Verify that reference changes affect generation, document the limits, and deliver the corrected PR candidate.

## 11. Native Qwen reference generation

- [x] 11.1 Select the requested Qwen-Image-Edit-2511 bundle and native one-/two-reference graphs, with neutral prompts and no additional custom nodes.

The planned Qwen live replacement and inference acceptance were superseded by
the user's teardown and fresh-IaC request. They were not performed; the native
workflow defaults remain available for a later authorized deployment.

## 12. Infrastructure ownership and full teardown

- [x] 12.1 Destroy the legacy rental, generated login, dedicated team API key and obsolete delivery permission through owning interactive applies.
- [x] 12.2 Move the component and maintained change to infra/image_gen with an independent AppRole, KV namespace and fresh backend.
- [x] 12.3 Wire team issuance into the component's exact key entry and remove personal credential/state references from active configuration.
- [x] 12.4 Validate source, representative diagram and native workflow structure without deploying replacement resources.
- [x] 12.5 Review the new desired state and prepare it for guarded publication; delivery receipts own publication and remote-review state.

## 13. Credential-only deployment

- [x] 13.1 Apply the merged AppRole/policies and dedicated key through owning scoped interactive workflows.
- [x] 13.2 Apply only password and web-login resources, retaining the password after correcting metadata access; create no instance.
- [x] 13.3 Verify credential access and absence of a rental, then record the sanitized receipt.

## 14. Authorized rental and Qwen acceptance

- [x] 14.1 Review and apply one rental with unchanged credentials and existing Secure Cloud/VRAM/price constraints.
- [x] 14.2 Verify TLS/authentication, native model/node availability, and both staged workflows.
- [x] 14.3 Execute both workflows with neutral character/style references; inspect prompt, character and style behavior and correct in-scope defects.
- [x] 14.4 Refresh connection outputs through the owning workflow and deliver sanitized acceptance evidence.

## 15. Requested rental release

- [x] 15.1 Destroy only the current rental through its owning reviewed apply and clear connection outputs through refresh-only apply.
- [x] 15.2 Verify the instance is absent and retain the component identity, credentials and reusable IaC.
- [x] 15.3 Update current deployment status and prepare the sanitized release record for delivery.

## 16. Selectable profiles and authenticated CLI

- [x] 16.1 Add the approved packaged CLI target with Vault login and Terraform endpoint discovery.
- [x] 16.2 Preserve Qwen as a profile and add FLUX.2 klein 9B with matching native reference workflows.
- [x] 16.3 Validate source and download access before renting one Flux instance through the reviewed owning apply.
- [x] 16.4 Verify authenticated CLI execution and neutral reference generation, and deliver sanitized evidence.
- [x] 16.5 Supply a first-party prompt sequence node and workflows for both profiles; verify prompt splitting and paired filename ordering offline.
- [x] 16.6 Execute a neutral three-prompt workflow on the GPU and verify one numbered image per prompt with shared references.

## 17. Base checkpoint replacement

- [x] 17.1 Destroy the current rental through its owning reviewed apply and clear connection outputs, retaining component credentials.
- [x] 17.2 Pin the approved FLUX.2 klein Base 9B FP8 checkpoint and update all three Flux workflows to native Base conditioning, 50 steps and CFG 4.
- [x] 17.3 Validate the source and verify gated model download access before creating paid capacity.
- [x] 17.4 Create the requested Base rental through its reviewed apply; stop workflow acceptance when the operator subsequently requests destruction.
- [x] 17.5 Prepare the checked source update and sanitized teardown status for guarded publication; delivery receipts own publication state.
- [x] 17.6 Wire an opt-in 96 GB hardware preset and owning targets without selecting it for deployment.

- [x] 17.7 Destroy the brief Base rental and clear its connection outputs, retaining credentials and source; no inference jobs were submitted.

## 18. ComfyUI tools relocation

- [x] 18.1 Consume custom nodes and workflows from projects/comfyui_tools; keep the CLI in tools/comfyui and verify Terraform packaging offline.

- [x] 18.2 Include the Prompt Sequence browser script through the native provisioner and verify the staged asset without renting.

## 19. Externally supplied model token

- [x] 19.1 Relinquish the redundant root's Terraform ownership through a reviewed removed block with destroy=false; verify the token entry is retained.
- [x] 19.2 Remove the completed transition, verify a no-change plan, then delete the token-writing root and its unused authentication/backend wiring.
- [x] 19.3 Document manual or external-IaC token ownership and verify the rental still reads the existing entry through offline checks; prepare the reviewed cleanup for publication.

## 20. Single JSON configuration

- [x] 20.1 Move image, model definitions, hardware presets and selected profiles into config.json; remove root variables, the tfvars preset and duplicate high-memory targets. Verify the existing mocked-provider checks read the packaged JSON.
- [x] 20.2 Update operator instructions and review the selected defaults and all profile contents against the previous candidate before publication; do not rent.

## 21. Restore Juggernaut

- [x] 21.1 Restore the previously approved Juggernaut XL v9 checkpoint at its existing revision as the default JSON profile, retaining Flux/Qwen options and hardware constraints.
- [x] 21.2 Supply neutral native SDXL text-to-image, image-to-image and prompt-sequence workflows; verify model paths, node types and graph links without renting.
- [x] 21.3 Document style-prompting and img2img limits, complete source validation and review, and prepare the follow-up for publication. No deployment is authorized.

- [x] 21.4 Restore the previously approved RES4LYF revision and unchanged requirements only for Juggernaut; stage matching style-reference and style-sequence graphs without a rental.

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

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

## 1. Role and integration

- [x] 1.1 Write a local packaged-role acceptance fixture before implementation.
- [x] 1.2 Implement validated asset installation, idle checks and restart/readiness handling.
- [x] 1.3 Package the role and use it from the profile-driven deployment.
- [x] 1.4 Verify installation, repeatability and failures; document exact limitations.

## 2. Complete the profile interface (superseded by section 3)

- [x] 2.1 Extend the deployment E2E fixture to install a profile by invoking the role directly with configuration/runtime inputs, before implementation.
- [x] 2.2 Move profile expansion, asset path mapping and readiness selection into the role while preserving individual-asset inputs and workflow tags.
- [x] 2.3 Verify the standalone profile path, all generated workflow catalogs and the existing role fixture; keep the infra playbook to configuration loading and one role invocation.

## 3. Plain role variables

- [x] 3.1 Adapt the deployment E2E fixture to supply the JSON registry to the playbook and verify a non-default profile before implementation.
- [x] 3.2 Remove registry inputs, conditional profile defaults and temporary profile facts from the role; translate JSON once in playbook variables.
- [x] 3.3 Verify all workflow catalogs, full local NoobAI deployment and generic role installation; retain AL-only Vault access.

## 4. Packaged node bundles

- [x] 4.1 Install selected repository-built node archives before updating the node whitelist and restarting ComfyUI.
- [x] 4.2 Verify the archive, NoobAI App and Tag Search through the packaged deployment fixture.

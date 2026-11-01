## Source acceptance on 2026-11-01

Upstream [3.7.17](https://github.com/traefik/traefik/releases/tag/v3.7.17)
was the latest stable release, published 2026-10-30. The publisher's checksum
manifest and GitHub asset digest agreed on Linux AMD64 archive SHA-256
`a564abd5fc8ee97d838270f188bb94cb3172cc113f4b51aa5e7dd6117eedb393`.
The lock carries its SRI equivalent. Bazel fetched and verified that archive.

The built binary reports `Version: 3.7.17` and `OS/Arch: linux/amd64`.
The shared Traefik role and both `infra/local_inference/ansible:ansible_bin`
and `infra/download/ansible:ansible_bin` packaged successfully. The module
update and lock reconciliation workflows produced no generated-file changes.

The [migration guide](https://doc.traefik.io/traefik/migrate/v3/#v3717)
changes Ingress-NGINX resource names. The prepared T3 host uses the file
provider. Its network-isolated fixture confirmed native pairing, real T3
routing through the new proxy, trusted-client success, rejection of missing
and untrusted client certificates, rejection of wrong SNI, and HTTP redirect.
See `infra/local_inference/openspec/changes/add-t3code-mtls/validation.md`.

Mandatory repository quality, affected-target lint, and owner OpenSpec checks
are selected in guarded delivery; private receipts own their exact-candidate
results. No service was upgraded or deployed, and live ACME was not exercised.

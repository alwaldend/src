## 1. Ingress routing

- [x] 1.1 Add the ComfyUI service map entry; verify rendered router, HTTPS target, SNI, and inherited mTLS.

## 2. Validation

- [x] 2.1 Build the ingress Ansible package and run repository checks; inspect the full rendered route alongside T3 Code.

## Evidence

Validated on 2026-11-01: ingress Ansible binary build, strict OpenSpec
validation, repository quality checks, and affected-package semantic lint
passed. Rendered TOML confirms the ComfyUI HTTPS backend, service-hostname
SNI, inherited `RequireAndVerifyClientCert`, and unchanged existing routes.
The coordinated Host Bot change records DNS and loopback verification.
Local artifacts and exact-candidate delivery receipts are under ignored
`out/comfyui-ingress-traefik/`. No live deployment was performed.

## 3. Authorized deployment

- [x] 3.1 Deploy the ComfyUI Traefik route through the owning Ansible playbook with `--tags traefik`; verify the active route.

On 2026-11-01, the authorized deployment ran the owning playbook with
`--tags traefik --limit host1.ingress.alwaldend.com`. Strict SSH host
verification used the existing server CA from Vault's TLS-verified public-key
endpoint in task-private known-hosts storage. The recap was 10 successful
tasks, 3 changes, and zero failures or unreachable hosts. The public ComfyUI
SNI endpoint presented a valid TLS certificate and required a client
certificate. Sanitized deployment evidence is under ignored
`out/comfyui-routing-deploy/`; authenticated public application access was
not exercised without a client certificate.

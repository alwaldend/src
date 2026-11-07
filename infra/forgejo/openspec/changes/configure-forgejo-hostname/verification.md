## Offline hostname verification

Pinned Ansible rendered the actual Forgejo configuration and both download
inventory views with placeholder credentials. The pinned Traefik binary then
passed 14 response checks against an isolated loopback upstream:

- Local Git root, repository and escaped-path/query requests redirect to Forgejo.
- Canonical and inventory hostnames reach the upstream; unrelated hosts return 404.
- Public Git requests redirect to GitHub; Go discovery reaches its upstream route.
- The local download view has no Git router and continues serving downloads.

Rendered Forgejo HTTP, SSH and root URL values use `forgejo.alwaldend.com`.
TLS configuration retains the Git alias and inventory hostname. The generated
shared AL configuration uses the canonical URL. Forgejo and runner packages
build, and Vault's offline Terraform formatting test passes.

The fixture uses proxy-generated local TLS and a synthetic upstream. It does
not establish live TLS issuance, real OIDC login or production service health.
No DNS, Vault or host configuration was deployed. Sanitized response evidence
and the repeatable helper remain in `out/forgejo-upgrade/hostname-fixture/result.json`
and `out/forgejo-upgrade/verify_hostname.py` in the task worktree.

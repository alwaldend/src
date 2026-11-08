# Repository infrastructure flow

Infrastructure components live under `infra/<component>` and, for personal
hosts, `users/<user>/<host>`. A component may carry any subset of these
directories, and each is a separate stage:

| Stage            | Owns                                                       |
| ---------------- | ---------------------------------------------------------- |
| `tf_setup`       | Provisioning: VMs, resource sets, bootstrap objects        |
| `tf`             | Service configuration: API-side resources, policies, roles |
| `ansible/`       | Host deployment: packages, users, files, services, ingress |
| `dnsconfig.json` | The component's record declarations                        |
| `BUILD.bazel`    | Vault authentication, backend, and secret injection        |

`README.md` owns the component's purpose and boundaries; `BUILD.bazel` owns
target structure and configuration wiring. Read them before changing any stage.

Prefer a component-owned `dns/` Terraform root with its own state and DNS-only
credentials, consuming the shared `dns_records` module. Keep canonical record
declarations in that component's `dnsconfig.json`. Existing combined roots can
remain until a scoped migration is requested; source movement alone does not
transfer state ownership. The
[DNS workflow](dns.md) describes runtime source discovery, ownership checks,
and the prerequisites for using those roots.

## Authentication and injection

Each component authenticates with its own Vault AppRole and receives the
minimum policy for its work. An `al_vault_auth` target extends `//tools/al:auth`
and selects the component AppRole. Compose small targets with `al_config` and
keep short declarations inline. Use local `defs.bzl` only when extraction
provides significant reuse or substantially simplifies complex graphs/templates.
Sibling BUILD files show the wiring. Components then call plugins:

- `injector` reads KV values, files, and operations and exposes them as
  environment variables or files.
- `tf_backend` supplies the Vault-backed Terraform HTTP state backend.
- Provider logins (`xo_login`, `pve_login`, `forgejo_login`, and so on) obtain
  provider credentials for a specific stage.

A plugin call only takes effect when its label matches the wrapper's
`plugin_labels` dictionary. The wrapper's `configs` dependencies
provide the registered plugin executables and their runtime files transitively;
do not repeat plugin binaries in the wrapper's `data`.

## Why one component can have several modules

A Terraform Vault AppRole is created and owned by `infra/vault/tf`. Components
consume it; they do not declare it. When a component writes Vault resources of
its own, it uses its own provider configuration and the injected token.

## Boundaries to preserve

- Keep one implementation of each behavior. Extend the owning module or role
  instead of copying it into a component.
- Keep environment-specific values in component configuration and reusable
  behavior in the shared module or role tree.
- Preserve target visibility and publication boundaries declared by the owning
  `README.md` and `BUILD.bazel`.

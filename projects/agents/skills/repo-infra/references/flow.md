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
| `al.lua`         | Vault authentication, backend, and secret injection        |

`README.md` owns the component's purpose and boundaries; `BUILD.bazel` owns
target structure; `al.lua` owns wiring. Read them before changing any stage.

## Authentication and injection

Each component authenticates with its own Vault AppRole and receives the
minimum policy for its work. `lib.vault_auth({name = "default", approle = {
name = "<component>"}})` selects it, and sibling `al.lua` files show the exact
DSL shape. Components then call plugins:

- `injector` reads KV values, files, and operations and exposes them as
  environment variables or files.
- `tf_backend` supplies the Vault-backed Terraform HTTP state backend.
- Provider logins (`xo_login`, `pve_login`, `forgejo_login`, and so on) obtain
  provider credentials for a specific stage.

A plugin call only takes effect when three things agree: the call's label in
`al.lua`, the `--plugin_label` arguments in the target's `run_args`, and the
plugin binary in the target's `data`. Keep them aligned.

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

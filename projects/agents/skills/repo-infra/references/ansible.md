# Ansible in this repository

## Follow the packaged execution model

1. Read the component's `ansible/BUILD.bazel`, parent `al.lua`, inventory,
   `group_vars`, and the closest sibling deployment.
2. Keep runtime inputs in the `pkg_files`/`pkg_filegroup` source set. Update its
   globs or explicit entries when adding a new file type or location.
3. Use `al_ansible_binary` with `//infra:ansible_cfg` and the repository Ansible
   collection. Execute through `al_binary_run`/`al_binary_run_map`, not an
   ad-hoc local command that bypasses packaging and injection.
4. Keep Vault injector and SSH/auth plugin labels aligned between `al.lua`,
   Bazel `data`, and `run_args`.

The root `infra/ansible.cfg` establishes Python, roles, collections, and YAML
result behavior. Do not add component-local overrides without a demonstrated
need. Reuse the host's package backend: `dev_vm/tasks/packages.yml` owns
ordinary versus rpm-ostree installation where that role manages dependencies.
A service role can expose caller-managed package installation instead of
bypassing that owner. On an OSTree host, staged dependencies may require a
reboot before installation or startup can use them; honor that prerequisite
and the host's writable systemd-unit path, commonly `/etc/systemd/system`.

## Roles

Reusable roles live in `projects/ansible_collection/roles/<role>`. A role
declares `defaults/main.yaml`, `tasks/main.yaml`, optional `templates/` and
`files/`, and an `al_ansible_role` target. `al_ansible_role` requires a defaults
file and auto-generates its documentation page, so keep defaults documented and
accurate.

Role defaults use `{{ undefined() }}` to mark values a deployment must supply.
An inventory `group_vars` file then provides them, often from an injected
environment variable, for example
`"{{ lookup('ansible.builtin.env', 'OPENHANDS_SESSION_API_KEY',
default=undefined()) }}"`.

## Write reliable automation

- Use fully qualified collection names where surrounding code does.
- Prefer idempotent modules over `shell` or `command`. When those modules are
  necessary, define `changed_when`, `creates`, or another explicit idempotency
  condition and quote untrusted input.
- Use handlers for restarts and notify them only when configuration changes.
- Preserve inventory and group naming, and put defaults at the narrowest
  correct scope.
- Check merged variable precedence: a child group must not redefine a variable
  as a template referring to itself. Inherit the parent value or use a distinct
  input name; otherwise it hides the injected value and recurses at rendering.
- Never commit secret values. Receive them through the component's Vault-backed
  injection flow and mark tasks that could reveal them `no_log: true`.
- Be cautious with destructive tasks, quorum changes, host reboots, and rolling
  operations; encode serial and health behavior where appropriate.

## Validate without deploying

Inspect generated labels, then build the packaged binary and use any existing
non-mutating test or syntax target:

```sh
bazel_agent bazel query '//path/to/ansible:*'
bazel_agent bazel build //path/to/ansible:ansible_bin
```

Do not run the normal `:ansible` or named playbook targets as a syntax check:
they can connect to inventory hosts and mutate infrastructure. If no repository
syntax-check target exists, validate YAML and packaging locally and report that
live execution was intentionally not performed. A successful package build
does not prove merged variables or templates render; use non-secret fixture
values for the affected rendering. Inspect the pinned launcher and rendered
service command when flags control exposure: an upstream frontend-only mode
can still spawn additional listeners. Verify the actual command and child
listener behavior relevant to the change, using the
[diagnostics reference](diagnostics.md) for live read-only observations.
Never use `--check` as proof of safety unless every involved role and module is
known to support check mode.

## Shared host roles

Component playbooks commonly import the shared `host`, `traefik`, and
`install_ca` roles before the component role. Ingress is owned by the Traefik
role: component roles install and run their service on loopback, and the
deployment declares the dynamic config that exposes it. Keep that split, and
do not duplicate a reverse proxy inside a component role.

For chained Traefik proxies, trace both TLS SNI and the HTTP Host header.
The backend URL normally supplies SNI, while
[`passHostHeader`](https://doc.traefik.io/traefik/reference/routing-configuration/http/load-balancing/service/)
defaults to retaining the original request Host; downstream routers must match
that host. Set explicit priorities for overlapping API and catch-all routers:
Traefik's [default priority](https://doc.traefik.io/traefik/reference/routing-configuration/http/routing/rules-and-priority/#priority-calculation)
is the rule length. Verify representative requests against the rendered rules.

When a backend router also matches a forwarded hostname owned by another
proxy, prevent certificate-domain inference from requesting that foreign
name. Configure explicit router `tls.domains` for the names that this host
must serve; with HTTP-01, their challenge requests must reach this certificate
issuer. Traefik's [ACME domain selection](https://doc.traefik.io/traefik/reference/install-configuration/tls/certificate-resolvers/acme/)
otherwise uses router Host rules when explicit domains are absent.

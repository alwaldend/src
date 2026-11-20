# Ansible in this repository

## Follow the packaged execution model

1. Read the component's `ansible/BUILD.bazel`, parent `BUILD.bazel`, inventory,
   `group_vars`, and the closest sibling deployment.
2. Keep runtime inputs in the `pkg_files`/`pkg_filegroup` source set. Update its
   globs or explicit entries when adding a new file type or location.
3. Use `al_ansible_binary` with `//infra:ansible_cfg` and the repository Ansible
   collection. Execute through `al_binary_run`/`al_binary_run_map`, not an
   ad-hoc local command that bypasses packaging and injection.
4. Keep Vault injector and SSH/auth plugin labels aligned between `BUILD.bazel`,
   Bazel `data`, and `plugin_labels`.

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
- Prefer Ansible's Jinja templates for Ansible-specific files and data
  transformations. Use Bazel to package the templates and supply input data
  through extra-vars files.
- Preserve inventory and group naming. Prefer `group_vars/<group>.yaml` for
  shared deployment settings, including SSH users and role inputs. Keep
  host and group membership in the inventory and prefer `host_vars`
  for host-specific overrides.
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
bazel query '//path/to/ansible:*'
bazel build //path/to/ansible:ansible_bin
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

The default web-service deployment runs `alwaldend.main.host`, then
`alwaldend.main.traefik` as the reverse proxy, then the component service
roles. The host role installs the CA and configures the firewall. Supply
firewall and Vault SSH host-signing inputs through group variables.

Keep the service listener private and put the reverse proxy's routing
configuration in the deployment. For a service without HTTP ingress, use the
host role followed by the service roles. Use a service-only play when the
operator explicitly targets an already configured host.

## Traefik configuration

Keep listeners, providers, and ACME resolvers in the deployment's static
config. Keep HTTP routers, middleware, certificate files, and upstreams in
the dynamic config.
For browser access, configure HTTPS and redirect HTTP on the networks where
clients can reach the proxy. Apply the deployment's firewall allowlist there,
and keep service and metrics listeners on loopback when only local processes
need them.

Use Let's Encrypt ACME for public Traefik instances, including shared ingress.
Use Vault ACME for local Traefik instances with a service-scoped Vault PKI role
and the shared role's external account binding enrollment. For HTTP-01, verify
that each issuer resolves the service hostname to the Traefik instance that
answers its challenge on port 80.

For direct LAN access, follow the [DNS record pattern](dns.md#scopes-and-landing-sites)
so both paths use the same service hostname. Public DNS sends Let's Encrypt
to shared ingress. DC1 DNS sends Vault to the guest's ACME listener. The
guest router, certificate, and application authentication serve both paths.

For services behind shared ingress, point ingress at the separate private
backend DNS name. The shared ingress transport verifies backend TLS against
the public service hostname. Issue the guest Traefik certificate for that
public hostname and match it in the guest's HTTP router. If Vault cannot
reach the guest's challenge listener, use a Vault-signed certificate from the
component signer before starting guest Traefik.

When a service needs a shared Basic Auth login, generate the username and
password in its Terraform root and store them in Vault. Use at least 128
characters for generated passwords. If a tool supports fewer, use its maximum
supported length. For bcrypt, use 72 ASCII characters to stay within its
72-byte input limit. Grant only the intended administrator group read access
to that login. Store a separate
bcrypt verifier for the deployment AppRole. Deliver that verifier to the
reverse proxy in a restricted file, keep it out of task logs, and strip the
authorization header before proxying. Preserve any existing ingress
authentication unless the change explicitly replaces it. Keep generated
credentials in protected Terraform state when managed generation is needed
to reconcile the login and verifier together.

Trace both TLS SNI and the HTTP Host header for chained proxies.
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

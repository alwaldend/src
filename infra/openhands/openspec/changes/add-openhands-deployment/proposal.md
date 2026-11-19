## Why

OpenHands provides an agent runtime as three cooperating services: Agent Canvas
(the browser client), an agent server (conversations, tools, and workspace
operations), and an automation server (schedules and event dispatch). The
repository has no way to deploy or reach them, and no Vault identity, DNS
records, or host automation for any of the three.

## What Changes

Add the `infra/openhands` component and its supporting Vault identity:

- A `src_infra_openhands` Vault AppRole declared as a module under
  `infra/vault/tf/approle_src_infra_openhands/`, with its own entity, group, SSH
  host-key signing role, and server PKI role for the component hostnames.
- `infra/openhands/dnsconfig.json` declaring the site-local component names, the
  `dc1`-prefixed ingress targets, and the public canvas CNAME.
- `infra/openhands/tf_setup` creating the canvas, secured agent server, and
  automation server VMs in the `src_infra_openhands` Xen Orchestra resource set.
- `infra/openhands/ansible` deploying each component natively (no containers)
  through the shared host and Traefik roles, with per-component dynamic config.
- An `openhands_server`, `openhands_canvas`, and `openhands_automation` role in
  the shared Ansible collection.
- Canvas in `infra/ingress` so it is publicly reachable while the existing
  client-certificate requirement is preserved.
- The unsecured agent server in `users/simeonwarren/host_bot`, with that host's
  deployment and secret injection owning its loopback-only service.

## Capabilities

### New Capabilities

- `infra-openhands`: Deploy and operate the OpenHands canvas, agent server, and
  automation server.

### Modified Capabilities

None.

## Impact

- `infra/vault/tf`: new AppRole module and shared group membership.
- `infra/xcp_ng/tf`: named template, storage, and network assignment for the
  OpenHands AppRole under its existing resource-set contract.
- `infra/openhands`: new component, provisioning, deployment, and DNS.
- `infra/ingress`: new canvas service entry under its existing routing and
  client-authentication requirements.
- `users/simeonwarren/host_bot`: native unsecured agent server and secret
  injection.
- `projects/ansible_collection/roles`: three new component roles.
- `third_party/com_openhands_openhands_canvas_ui_tool`: pinned upstream
  compatibility module for the agent server's tool registration.

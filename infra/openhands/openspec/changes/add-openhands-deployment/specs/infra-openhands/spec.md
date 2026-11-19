## ADDED Requirements

### Requirement: Separate components by responsibility and trust boundary

The deployment SHALL run Agent Canvas, the agent server, and the automation
server as distinct components, because each owns a different responsibility.
The agent server SHALL run natively with its service account's host
permissions. The deployment documentation SHALL identify those permissions as
the execution boundary and the workspace as its working directory.

#### Scenario: Inspect the component set

- **WHEN** a reader inspects the deployment
- **THEN** the canvas serves the browser client, the agent server runs
  conversations and workspace operations, and the automation server stores
  schedules and dispatches conversations
- **AND** the documentation states that the agent server acts with its
  service account's permissions on its host

### Requirement: Scope Vault access to the component AppRole

`infra/vault/tf` SHALL declare the `src_infra_openhands` AppRole through a
module under `infra/vault/tf/approle_src_infra_openhands/`, including its
entity, group, SSH host-key signing role for `openhands.alwaldend.com`, and
server PKI role for the component hostnames. The AppRole SHALL receive only the
secrets, SSH, and certificate policies it needs, and the component SHALL
authenticate with it rather than a shared or administrative identity.

#### Scenario: Inspect the AppRole scope

- **WHEN** the Vault configuration is evaluated
- **THEN** the AppRole, SSH role, and server PKI role are declared in the
  AppRole module and referenced by the shared identity groups
- **AND** the component's `al.lua` selects that AppRole by name

### Requirement: Provision hosts separately from deployment

`tf_setup` SHALL create the canvas, secured agent server, and automation server
VMs in the `src_infra_openhands` Xen Orchestra resource set. Each VM SHALL
derive its hostname, address, and MAC from one record in `dnsconfig.json`.
Ansible SHALL deploy the components natively, without containers.

#### Scenario: Add a component host

- **WHEN** a new component is added
- **THEN** one record in `dnsconfig.json` owns its name and static address
- **AND** the provisioning `vms` map references that record prefix

### Requirement: Own ingress outside the component roles

Component roles SHALL install and run their services on loopback and SHALL NOT
declare their own ingress. The deployment SHALL expose them through the shared
Traefik role with per-component dynamic configuration.

#### Scenario: Expose a component through Traefik

- **WHEN** a component host is deployed
- **THEN** the shared `traefik` role writes its dynamic configuration
- **AND** the component role itself declares no router, entry point, or
  certificate resolver

#### Scenario: Start the canvas frontend

- **WHEN** the canvas service starts
- **THEN** the package's static server serves the packaged browser assets on
  the declared loopback address
- **AND** it starts no agent server, automation server, or additional proxy

### Requirement: Keep the browser on one origin

The deployment SHALL document configuring the browser backend with the canvas
origin and shared session key. The canvas origin SHALL proxy the backend API
paths to the agent server and automation server ahead of the frontend route.
Backend targets SHALL use site-local hostnames, and downstream routers SHALL
accept the forwarded canvas Host header.

#### Scenario: Serve the browser client

- **WHEN** the browser uses the canvas origin as its configured backend and
  requests an automation or agent server API path
- **THEN** Traefik forwards it to the matching backend service
- **AND** the browser issues no cross-origin request
- **AND** the inter-service hop selects the downstream router without
  traversing the public ingress

### Requirement: Preserve the public ingress client-authentication policy

The component SHALL declare the public `canvas.openhands.alwaldend.com` CNAME
and the site-local canvas backend address consumed by its ingress service
entry. Canvas integration SHALL satisfy the existing
[ingress routing and client-authentication requirements](../../../../../../ingress/openspec/specs/infra-ingress/spec.md).

#### Scenario: Inspect public canvas routing

- **WHEN** the OpenHands DNS and ingress service declarations are inspected
- **THEN** the public canvas CNAME selects the shared ingress and the service
  entry selects the component's site-local backend address
- **AND** the integration uses the ingress owner's existing TLS policy

### Requirement: Inject secrets through the component's Vault flow

Secret values SHALL be stored in Vault and referenced only by path, environment
variable name, or plugin label. The component `al.lua` SHALL inject them through
the repository's `injector` plugin, and tasks that could reveal them SHALL use
`no_log`. The session API key SHALL be shared by the agent server, the
automation server, and the browser because all three authenticate the same
`X-Session-API-Key` value.

#### Scenario: Deploy with injected secrets

- **WHEN** the deployment reads the session key, secret key, or automation KV
  secret
- **THEN** it receives them from the injected environment rather than committed
  files
- **AND** the environment file that receives them is written with `no_log`

### Requirement: Deploy the unsecured agent server within its host's owner

The unsecured agent server SHALL be owned and deployed by the host project that
owns its host, not by the OpenHands VM package. It SHALL bind to loopback,
receive no session API key, and add no ingress or firewall exposure of its own.

#### Scenario: Review the unsecured server deployment

- **WHEN** a reader inspects the unsecured agent server deployment
- **THEN** the host project's inventory, playbook, group variables, and secret
  injection declare it
- **AND** the OpenHands Ansible inventory contains only the three VM components

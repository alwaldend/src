## Context

The repository provisions Fedora VMs through Xen Orchestra and configures
guest services through packaged Ansible playbooks.

## Goals / Non-Goals

**Goals:** One guest, persistent local storage, direct scraping, and VMUI
access through authenticated ingress or SSH.

**Non-goals:** Grafana, high availability, and off-VM backup
provisioning.

## Decisions

- Keep the AppRole root under the component's `approle` directory. It uses
  the shared Vault infrastructure credentials and dedicated Vault HTTP state.
- Use the guest address in the component-owned DNS declaration.
- Use the Xen Orchestra endpoint and session token injected by the XO login
  plugin through the provider's native environment variables.
- Set the Ansible SSH user in the VictoriaMetrics inventory group variables.
- Keep the database on the persistent guest disk. Use a dedicated systemd
  account and reserve disk space for the OS.
- Fetch and verify the Linux AMD64 release with Bazel. Package the extracted
  executable with the Ansible playbook. Render scrape jobs and active ingress
  WireGuard addresses in an Ansible YAML template using the packaged ingress
  inventory. Scrape VictoriaMetrics itself and other Traefik hosts over the
  private network.
- Run the shared host and Traefik roles before VictoriaMetrics after SSH and
  cloud-init are ready. Issue host certificates for the inventory and service
  names. Limit the Ansible user's forwarding to loopback VMUI. Supply the
  binary path through group variables.
- Keep the VictoriaMetrics self-scrape on IPv4 loopback. Guest Traefik serves
  shared ingress peers and the DC1 LANs over HTTPS. It obtains and renews its
  service-hostname certificate through Vault ACME. Vault reaches guest port 80
  using DC1 DNS. Shared ingress uses Let's Encrypt, verifies the guest's
  certificate against the service hostname, and requires a trusted client
  certificate at the public endpoint.

## Risks / Trade-offs

- One disk and guest constitute a single failure domain. Restart or host
  failure interrupts scraping and queries.
- The disk allocation and retention are starting values for an unmeasured
  workload. VictoriaMetrics stops ingesting at the configured disk floor.

## Web login

Terraform generates a shared login for the client-certificate-protected route.
Vault grants the administrators read access to the login and stores the
bcrypt verifier separately for the deployment AppRole. Traefik reads it from
a private file and removes the authorization header before proxying.
The managed random password remains in restricted Terraform state so both
Vault entries can be reconciled without splitting the credential pair.

Direct LAN clients use the same service hostname as public clients. DC1 DNS
resolves it directly to the guest, and public DNS points to shared ingress.
The `dc1.` upstream hostname resolves to the guest's private address through
public DNS. Both access paths use the same guest router, certificate, and
Basic Auth login. Shared ingress requires a client certificate.

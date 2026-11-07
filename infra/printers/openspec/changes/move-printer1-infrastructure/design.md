## Context

See [proposal.md](proposal.md) for the requested ownership change.

## Goals / Non-Goals

Implement and verify checked-in definitions, then execute the separately
authorized provisioning and state adoption while preserving live printer
objects. Printer device configuration and firmware remain outside this change.

## Decisions

**Verdict: proceed.** Move DNS to dns/ and the firewall root to tf/. Retain Terraform resource addresses, record keys, and firewall comment to minimize live changes. Give each root a new owner-local backend. Freeze former roots and transfer ownership using targeted state removal and exact-ID imports, without deleting live objects; do not save state or plan artifacts. The old firewall backend is retired only after adoption. Deployment and state operations received separate explicit user authorization; the completed receipt records their scope. Linked credential changes live in infra/mikrotik/openspec/changes/provision-printer1-credentials and Vault changes in infra/vault/openspec/changes/authorize-printers-infrastructure.

Use the shared infra.mikrotik helper to inject MIKROTIK_HOST, MIKROTIK_USER
and MIKROTIK_PASSWORD directly. The default provider needs no Terraform
credential variables; DNS aliases keep their separate explicit inputs.

The network diagram groups the LAN actor and printer inside LAN, and the Wi-Fi
actor inside Wi-Fi. Both nodes are labelled Actor; the router remains outside
the network groups, with Wi-Fi to its left and LAN to its right. The printer
connection needs no redundant wired-network label. The Wi-Fi edge identifies the printer-only
exception to default forwarding restrictions. The printer connects to the
router with a bidirectional network arrow. DNS labels identify the devices;
the printer's name still resolves to its own address.

The authorized live migration used a temporary state-only MikroTik target
with the original firewall backend and owner identity. It needed no ingress
credential or RouterOS login to list and remove state bindings. It declared no
router resources and exposed no apply or destroy operation. The target
and its backend selector were removed after transferring the firewall binding.

## Risks / Trade-offs

Moving files does not move state. Follow the printer project's migration guide
before running operational roots. Offline checks cannot establish provider
permissions, deployed adoption or Wi-Fi connectivity.

## Placement adoption decision

Proceed with creation-only insertion metadata. The pinned provider marks
place_before as ForceNew but does not read it back from RouterOS. Import
therefore cannot retain that attribute. Ignore its updates after creation,
preserving guarded initial insertion and the existing rule ID; verify actual
order through the read-only API. Reject replacement of the live rule and avoid
editing or exporting raw state merely to restore metadata.

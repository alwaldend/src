# Terraform extraction

Extracted provisioning source from the task-owned download candidate
`bfdcdb6f3f035d60cfe562b0b5db99eb4b042298`, on trunk
`636f0baef81f1051f6a796b391b6b85fec80f4d1`. The three Terraform roots, their
existing mock-provider tests, provider assignments, DNS declarations, and XO
helper retain their source bytes. Component AL retains only provisioning
authentication; Ansible injection remains in the host-configuration PR.

The original candidate passed backend-disabled initialization, validation,
and mocked plans for all three roots. The split is checked again against its
own packaged inputs, including DNS ownership, AL provider labels, Terraform
formatting, semantic lint, OpenSpec validation, and repository quality.
Delivery receipts and sanitized command results live in ignored `out/split/`.

Generated DNS snapshots remain unchanged, following the manual regeneration
workflow merged in #107. No live provider plan, apply, VM provisioning, or DNS
cutover is part of this extraction. The local address, provider inventory,
and real host/disk behavior still require authorized operational acceptance.

## Hostname bootstrap review

The host consumer now connects directly by its inventory FQDN. Following the
old local host CNAME to the apex would reach the previous deployment before
cutover. The revised source gives that unique host the canonical A record,
uses it for the VM address, and makes the local download name an alias.
The component DNS root can establish host reachability before the separate
apex cutover. Ansible carries no IP address. The apex owner retains its separately declared
A record so local mail/TXT records keep their existing behavior.

Verdict: revise the host DNS dependency direction. Keeping the old alias would
require premature apex cutover; retaining an Ansible IP override contradicts
the requested hostname-only connection. The shared RouterOS module already
supports these A/CNAME mappings. Offline graph and DNS checks cover the
revised declarations; no live DNS change is performed.

## Apex cutover review

The apex A/AAAA instances and replacement CNAME have independent Terraform
keys. The component DNS runbook now requires two reviewed, committed source
revisions and separate saved plans: first delete only the old public apex
address records with no new CNAME present, verify deletion, then create the
CNAME from the final revision. It records the maintenance interval, ownership
adoption, unchanged saved-plan inputs, and recovery preconditions. This is a
required future operating procedure; neither phase has been executed.

## XO identity bootstrap review

The deployment sequence now links the owning XO bootstrap procedure and
names `XO_BOOTSTRAP_APPROLE=src_infra_download` explicitly before resource-set
assignment. That first-login operation needs its own authorization and recorded
result because identity discovery does not create XO users. The subsequent
owning-root apply reconciles bindings in IaC; the AppRole must leave
`approles_pending_oidc_login` before planning the local VM. No bootstrap was
executed by this source change.

## Endpoint ownership review

The final endpoint model removes repeated provisioning values without
extending the shared DNS schema. The public apex follows the stable download
service name; only the component records contain the Yandex host endpoint.
The local apex A record is the sole local address, read by the VM root and
followed by the unique host CNAME. Ansible still uses its inventory FQDN with
no override or DNS-file lookup.

Verdict: use existing DNS aliases and an explicit local preparation phase.
Before Ansible, an authorized, committed intermediate revision changes only
the dc1 apex while retaining public Pages records. The runbook documents the
local maintenance interval and checks named-host resolution before continuing.
This supersedes the earlier independent-host-A layout above. The public
A/AAAA-to-CNAME cutover remains a later two-phase operation. Source validation
covers the final graph; none of these DNS phases has run.

## Authentication configuration review

PR #108 review of `abbe37992e71914eee67771e38b22d1eccba46b5` identified a
hardcoded Vault OIDC provider in the shared XO helper and an unnecessary
download-specific Yandex role exception. The helper now requires the caller's
`oidc_provider`; download and OpenHands explicitly select the existing
provider. Missing or empty names fail configuration rendering.

The Yandex folder owner and shared module return to their trunk definitions:
every AppRole receives the standard folder-scoped `admin` service account.
The download-only `editor` override and its unused module parameter are
removed. This changes source configuration only; no IAM operation was run.

## Decisions

Use the existing dns_access module for src_infra_dev with dc1 and global views. The standard AppRole namespace policy already covers its own DNS backend; do not grant access to either legacy backend. Remove the unused runner-config grant because the runner role is inactive. The user retired OpenHands during rollout; remove its service-secret policy and read reference.

Decision review: proceed. A separate DNS identity would preserve duplicate authentication for one owner. The shared role already administers both hosts; the added policy only reads the two existing DNS provider credentials. Both legacy backends are now empty; retire their identities and grants after the authorized DNS handoff and secure1 recreation.

## Validation and rollout

Validate Terraform formatting, semantic lint and packaged AL configuration offline. Applying Vault policy and transferring DNS state are separate operations; retain old state access through that handoff.

The user chose secure1 recreation: add the shared development entity to the standard provisioning group for its scoped XO allocation. Remove the old runner identity only after its VM is destroyed and its backend has no managed resources.

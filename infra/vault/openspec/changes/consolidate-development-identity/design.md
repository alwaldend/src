## Decisions

Use the existing dns_access module for src_infra_dev with dc1 and global views. The standard AppRole namespace policy already covers its own DNS backend; do not grant access to either legacy backend or add generic infrastructure-group membership. Remove the unused runner-config grant because the runner role is inactive. Keep the OpenHands reference at its existing canonical secret location.

Decision review: proceed. A separate DNS identity would preserve duplicate authentication for one owner. The shared role already administers both hosts; the added policy only reads the two existing DNS provider credentials. Legacy identities remain transition dependencies until an authorized state handoff proves they can be retired.

## Validation and rollout

Validate Terraform formatting, semantic lint and packaged AL configuration offline. Applying Vault policy and transferring DNS state are separate operations; retain old state access through that handoff.

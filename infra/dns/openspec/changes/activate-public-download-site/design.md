## Context

See [proposal.md](proposal.md). The adopted shared DNS state currently owns
eight public Pages apex records. Local routing and public mail share this
declaration, so complete saved plans must be inspected for unrelated changes.

## Goals / Non-Goals

Complete only the approved apex cutover. Keep credentials in the existing
Vault injection flow and do not regenerate manual DNS snapshots or change
other owners' records.

## Decisions

Use the canonical download CNAME instead of copying the Yandex IP. Remove the
eight A/AAAA records in one saved Terraform plan, then create the CNAME through
a second saved plan; a combined plan cannot guarantee Cloudflare's conflicting
record types are removed first. Preserve each temporary declaration and plan
in private task evidence, with the final desired declaration checked in.

Enable and verify the public Nginx site through its owning Ansible tag after
content sync. Change DNS, then deploy the Traefik tag to request apex/www
HTTP-01 certificates against the correct public address. This avoids issuing
certificate requests while the names still point to Pages.

## Risks / Trade-offs

- A short DNS transition and certificate bootstrap interval is unavoidable.
  Inspect exact plans first and execute the prepared stages consecutively.
- Cached Pages answers may persist for the prior 300-second TTL. Verify both
  authoritative answers and direct public-host HTTPS, then public recursion.
- If acceptance fails, preserve diagnostics and restore the prior Pages
  declaration through the same owner workflow, removing the CNAME first.

## Migration Plan

1. Verify the local release, authoritative Ansible mirror, and public-host content.
2. Review/apply exactly eight public apex deletions, then exactly one CNAME creation.
3. Activate public Traefik routing and verify trusted HTTPS, assets, redirects,
   mail/local preservation, and a no-change Terraform plan.
4. Record acceptance and publish the final source through repository delivery.

## Context

The user explicitly requested reverting the public root-site DNS after an
authorized cutover. The Yandex VM is still blocked by HDD quota. The current
source contains the download CNAME; the prior GitHub Pages addresses are
available in Git and the pre-cutover provider inventory.

## Goals / Non-Goals

**Goals:** Restore the eight public addresses and preserve local routing and
all unrelated records. Make the restored state the maintained declaration.

**Non-Goals:** Change VM sizing, quota, host services, download subdomains,
local DNS, or generated DNS documentation snapshots.

## Decisions

Proceed with the requested rollback through `infra/dns/tf`. Restore the prior
record keys and values from Git, changing their destination to `global` so
they cannot replace the current local apex. Remove only `download_global`
from the apex owner's declaration.

Use two reviewed saved plans against this declaration: target the existing
apex CNAME for deletion first, then create the eight missing A/AAAA records
with a fresh plan. Inspect each plan and authoritative provider inventory
between phases. This preserves the runbook's required ordering without
introducing an intermediate declaration that could be deployed accidentally.

The alternative of applying the full plan at once can race record creation
against CNAME removal; direct provider writes would bypass desired state.

## Risks / Trade-offs

- DNS has an address gap between phases; complete the second phase promptly.
- Resolver caches may retain the former CNAME for its 300-second TTL.
- Restoring historical `dsp: all` would change the local site; restrict the
  restored records to the global view and compare the remaining declaration.
- A plan may include unrelated drift; reject actions outside the single
  CNAME deletion and eight address additions. Verify final convergence.

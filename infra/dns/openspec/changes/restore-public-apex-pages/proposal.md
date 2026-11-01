## Why

The public apex was switched to the download endpoint before the Yandex VM
could be provisioned because of HDD quota. The user requested restoring the
root site to GitHub Pages while the public download server remains undeployed.

## What Changes

- Restore the four public apex A records and four AAAA records for GitHub Pages.
- Remove the public apex CNAME to download from the maintained declaration.
- Preserve local apex routing, download subdomains, mail, and unrelated records.
- Apply the rollback through the existing DNS Terraform root in two saved plans.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `infra-dns`: Keep the public root site on GitHub Pages until a later authorized cutover.

## Impact

Only the public apex records in `infra/dns/dnsconfig.json` change. No new
provider, package, credential, VM, or host configuration is required. The user
explicitly authorized the rollback, including removal of the apex CNAME.

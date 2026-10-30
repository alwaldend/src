## Why

The user authorized publishing the website locally, mirroring it to Yandex,
verifying the public copy, and then moving the root site from GitHub Pages.
The earlier restoration remains historical evidence of the deferred rollout.

## What Changes

- Replace only the eight public apex Pages A/AAAA records with an apex CNAME
  to `download.alwaldend.com`, through the existing owner Terraform state.
- Retain local apex, mail, staging, and the existing `www` alias.
- Verify authoritative DNS, trusted apex/www HTTPS, and existing public services.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `infra-dns`: Route the activated public root site through the download host.

## Impact

`infra/dns/dnsconfig.json`, its Terraform root and runbook, and the coordinated
download-host rollout. No new provider, dependency, or credential is required.

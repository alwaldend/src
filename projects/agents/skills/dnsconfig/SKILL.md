---
name: dnsconfig
description: Work with DNS record definitions and safe DNS validation in this repository.
---

# Manage DNS declarations

DNS configuration is declarative and split by owner. Project-owned records
live in `projects/<project>/dnsconfig.json`; infrastructure records live under
the matching `infra/` owner. `infra/dns/dnsconfig.js` is the bridge that turns
JSON records into DNSControl domains and modifiers.

## Required workflow

1. Read the nearest existing `dnsconfig.json` before changing its schema.
2. Add project records to the owning project, not to a global list.
3. Use the existing record shape:
   - `CNAME.name` is the subdomain.
   - `CNAME.target` is the upstream hostname.
   - Cloudflare proxy flags are not used.
   - `dsp` selects `global`, `dc1`, or another supported scope.
4. Wire a new owner JSON file into `infra/dns/dnsconfig.js` and its owning
   `BUILD.bazel` `dnsconfig` filegroup.
5. Validate with:

   ```sh
   bazel_agent bazel build //infra/dns:all
   bazel_agent bazel test //:repo_quality_test
   ```

Never run `//infra/dns` or `//infra/dns:dns.deploy` as validation. Those
targets contact DNS providers and can change live records. Use preview only
when the user explicitly authorizes that exact live read operation.

## Landing subdomains

Project landing sites use `<project-slug>.alwaldend.com`, CNAME to
`pages.alwaldend.com`, and are unproxied so GitHub Pages can serve them
directly and issue their certificates. Keep the `pages` A/AAAA records with
the landing records so GitHub Pages traffic follows one address set. The apex
remains separately configured; never replace the apex with a project wildcard
or CNAME.

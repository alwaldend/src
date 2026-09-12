# DNS declarations in this repository

DNS configuration is declarative and split by owner. Project-owned records live
in `projects/<project>/dnsconfig.json`; infrastructure records live under the
matching `infra/` owner. `infra/dns/dnsconfig.js` is the bridge that turns JSON
records into DNSControl domains and modifiers.

## Required workflow

1. Read the nearest existing `dnsconfig.json` before changing its schema.
2. Add project records to the owning project, not to a global list.
3. Use the existing record shape:
   - One record per key, holding exactly one `A`, `AAAA`, `CNAME`, or `NS`.
   - `CNAME.name` is the subdomain; `CNAME.target` is the upstream hostname.
   - Cloudflare proxy flags are not used.
   - `dsp` selects `global`, `dc1`, or `all`.
4. Expose the file through the owning `BUILD.bazel` `dnsconfig` filegroup, and
   add the owner to the record lists in `infra/BUILD.bazel` and
   `infra/dns/BUILD.bazel`. For project workspaces, this is driven by the
   `PROJECTS` registry instead.
5. Validate with:

   ```sh
   bazel_agent bazel test //infra/dns:config_test
   ```

Never run `//infra/dns` or `//infra/dns:dns.deploy` as validation. Those targets
contact DNS providers and can change live records. Use preview only when the
user explicitly authorizes that exact live read operation.

`infra/dns/zones/*.zone` are generated BIND snapshots used to diff against
`infra/mikrotik/router1.rsc`; never hand-edit them.

## Scopes

- `dc1`: site-local names resolved by the dc1 router.
- `global`: public Cloudflare records.
- `all`: emitted to every scope.

A service that must be reachable both site-locally and publicly typically
declares a `dc1` record for the real host and a `global` `CNAME` to
`ingress.alwaldend.com.`, following the `infra/ingress` convention of an
explicitly prefixed site-local target. Public reachability is not the same as
public authorization: an ingress service that keeps a TLS client-auth policy
still requires a client certificate.

## Landing subdomains

Project landing sites use `<project-slug>.alwaldend.com`, CNAME to
`pages.alwaldend.com`, and are unproxied so GitHub Pages can serve them directly
and issue their certificates. Keep the `pages` A/AAAA records with the landing
records so GitHub Pages traffic follows one address set. The apex remains
separately configured; never replace the apex with a project wildcard or CNAME.

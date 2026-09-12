# DNS declarations in this repository

Canonical DNS declarations remain in each owner's `dnsconfig.json`. Read the
[DNS README](../../../../../infra/dns/README.md) for linting, owner-local
Terraform roots, and the cutover procedure. The shared
[`dns_records` Terraform module](../../../../../projects/tf_modules/dns_records/README.md)
owns the schema, normalization, and provider resource mapping.

## Declare and package records

1. Read the owner's existing declarations and Terraform stage before editing.
2. Keep records in the owning source file. Each logical record key may contain
   multiple supported type members (`A`, `AAAA`, `CNAME`, `NS`, `MX`, `TXT`)
   and a `dsp` list selecting `global`, `dc1`, or `all`. Follow the shared
   module's schema; Cloudflare records remain unproxied.
3. Use the shared `dns_records` module in the owner's `tf_setup` root when it
   exists, otherwise `tf`. Consume the JSON instead of copying record values
   into Terraform. Include the source filegroup and module sources in the
   root's Bazel data and preserve its AppRole, injection, and backend flow.
4. Run the offline checks:

   ```sh
   bazel_agent bazel run //infra/dns:lint
   bazel_agent bazel test //infra/dns:config_test
   ```

The linter discovers raw `dnsconfig.json` files in the current workspace at
runtime, including nested project workspaces, and prints the declarations as
a table. There is no central source registry to update. `config_test` runs
the offline linter unit suite and checks every current checkout declaration;
run `:lint` to inspect the declaration table.

Each canonical domain name must be managed by one source file across all
record types and views. Multiple values or types for that name may share that
file; separate names in the same DNS zone may have different owners. Duplicate
JSON object keys are errors. Resolve a duplicate at its owning source rather
than excluding that source from discovery.

## Operational prerequisites and migration

Before invoking an operational root, establish its deployed AppRole access and
the referenced Vault credentials: Cloudflare token, plus RouterOS username and
password for `dc1`. An explicit Cloudflare zone ID is optional; enabled global
records resolve the zone by name when it is absent or empty, requiring zone-list
and zone-read access. DNS and ingress share the RouterOS endpoint in the infra
AL configuration. These defaults require no additional metadata fields in
existing credential entries. Follow the owning `al.lua` and cutover procedure
for exact references and independent endpoint recovery requirements. Disabled
DNS management does not bypass credential injection or guarantee that provider
configuration is unnecessary.

Source implementation does not establish live adoption. Follow the DNS
README's cutover procedure to freeze old central deployment jobs and source
revisions, import existing records into the owning state, and review adoption
before enabling writes. The same procedure owns rollback and recovery.
Terraform plans contact live systems; imports, applies, and state operations
can mutate them. Keep those operations within the user's exact authorization.

## Scopes and landing sites

`dc1` selects site-local router records, `global` public Cloudflare records,
and `all` both views. Follow the owning component's ingress convention when a
service needs both local and public names. Public DNS reachability does not
remove an ingress service's authentication requirements.

For project landing names and targets, follow the DNS README and the project's
existing declarations. Keep apex and shared hosting records with their
declared owner instead of copying them into project files.

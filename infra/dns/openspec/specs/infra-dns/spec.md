# Infrastructure DNS

## Purpose

Describe the repository-owned DNSControl configuration for the global and dc1
views of `alwaldend.com`, including project-owned landing records. This baseline
describes checked-in configuration and wrappers; deployed records were not
observed.

Baseline source revision: `550d7e79b1f5fdbc2b6017b75178471d6914082f`.
Observation date: 2026-09-08. Sources are linked in full; no excerpts are used.

Sources: [component documentation](../../../README.md),
[target definitions](../../../BUILD.bazel),
[DNSControl configuration](../../../dnsconfig.js),
[provider configuration](../../../providers.json), and
[credential injection](../../../al.lua).

## Requirements

### Requirement: Aggregate records from their owning components

The DNS package SHALL assemble its record input manifest from the declared
infrastructure sources and the project registry, including nested `rules_*`
modules through their external Bazel repositories. Project landing CNAMEs SHALL
remain owned by each project's `dnsconfig.json`.

#### Scenario: A registered project contributes landing records

- **WHEN** the DNS configuration is assembled for a registered project
- **THEN** its owning DNS target is included in the record manifest
- **AND** the landing CNAME points directly to `alwaldend.github.io.` as required
  by the component documentation

### Requirement: Preserve distinct global and site-local views

The configuration SHALL combine common apex and mail records with records
assigned to the `global` or `dc1` view. It SHALL reject a record that names an
unknown destination view.

#### Scenario: A record is assigned only to dc1

- **WHEN** a record declares `dsp: ["dc1"]`
- **THEN** its modifier is added to the dc1 domain configuration
- **AND** it is not added to the global domain configuration by that declaration

#### Scenario: A record names an unsupported destination

- **WHEN** a record references a destination missing from the modifiers map
- **THEN** DNSControl configuration evaluation fails with `invalid dsp`

### Requirement: Validate configuration without provider access

The `//infra/dns:config_test` target SHALL invoke the pinned DNSControl `check`
command with the configuration and generated manifest, without requiring
provider credentials or contacting DNS providers.

#### Scenario: Validate a proposed record change locally

- **WHEN** the configuration test is run
- **THEN** DNSControl checks the assembled configuration without a deployment

### Requirement: Inject credentials into operational wrappers

The preview and push wrappers SHALL pass `providers.json` explicitly and obtain
its environment references through the AL Vault injector. The provider mapping
SHALL use Cloudflare for the global view, MikroTik for dc1, and BIND files for
both views.

#### Scenario: An authorized operator previews a DNS change

- **WHEN** `//infra/dns:dns.preview` is invoked with its injected environment
- **THEN** DNSControl receives the declared configuration and provider mapping
- **AND** credentials are read from environment references rather than a rendered
  credentials file

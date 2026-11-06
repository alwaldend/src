---
title: Vast.ai Terraform provider
description: Terraform management of Vast.ai compute and team resources
statuses:
  - in_progress
languages:
  - go
tags:
  - terraform
---

![Terraform provider request flow](assets/architecture.svg)

This project implements a Terraform Plugin Framework provider for the
[Vast.ai REST API](https://docs.vast.ai/api-reference/hello-world).
The provider type is `vastai`; the source address is
`registry.terraform.io/alwaldend/vastai`. This is a development provider;
no Terraform Registry release is published by this project setup.

The shared [API client](../../internal/vastai_sdk/README.md) is generated offline from
the official OpenAPI 3.1 schema, pinned by commit and SHA-256. The repository's
pinned OpenAPI Generator produces Go source directly into Bazel outputs.
The adapter retains endpoint validation and a few operations omitted from the schema.

## Coverage

- `vastai_offers` searches GPU offers using a JSON object query and returns
  `offers` in API order: `id`, `gpu_name`, `num_gpus`, `gpu_ram_mb`, and
  `hourly_price`. An empty search produces an empty list.
- `vastai_instance` rents an offer with `offer_id`, Docker `image`, and
  positive `disk_gb`. Optional creation settings are `runtype` (default
  `ssh`) and sensitive `onstart` (default empty). Supported launch modes are
  `ssh`, `ssh_proxy`, `ssh_direct`, `jupyter`, `jupyter_proxy`, and
  `jupyter_direct`. Changing creation settings replaces the instance.
  Refresh normalizes compound API modes and retains a compatible configured
  value; imports select a supported mode. Missing or unsupported modes produce
  a read error and retain state.
- `label` (default empty) and `state` (`running` or `stopped`, default
  `running`) update in place. Refresh reads drift, observed `actual_status`,
  `ssh_host`, `ssh_port`, and `public_ip`. Destroy terminates the contract;
  HTTP 404 during refresh or deletion, or an explicit `instances: null` in
  a successful refresh response, means it is already absent. Missing instance
  objects, invalid image/disk/launch metadata, mismatched IDs, and API failures
  remain errors and preserve state.

Apply acknowledges the rental or update request. It does not wait for boot,
power-state convergence, or SSH readiness. Connection fields are initially
null; use `terraform apply -refresh-only` to observe them after provisioning.
Offers can disappear between plan and apply. Stopped instances continue to
incur storage charges. Spot bids, templates, volumes, and
membership management are outside the supported surface.

## Non-compute resources

- `vastai_ssh_key`: required `public_key`, computed numeric-string `id`.
  Updates change public material in place; import uses the key ID. The
  provider never generates or stores private keys. Vast.ai may propagate
  registered keys to existing instances in the account.
  SSH keys require personal account context; Vast.ai rejects them in teams.
- `vastai_environment_variable`: required uppercase `name` and sensitive
  `value`; `id` equals the name. Values update in place; changing the name
  requires replacement. Import uses the uppercase name. Values remain in
  state, so use a secure backend.
  Masked API reads preserve the stored value and cannot detect external value
  changes. After importing a masked variable, configure its value for the next apply.
- `vastai_team_role`: required `name` and requested JSON-object `permissions`, computed
  `id` (role name), numeric `role_id`, and `effective_permissions` including
  Vast.ai's baseline grants. Permissions update in place;
  changing the name requires replacement. Import uses the role name. The
  resource does not assign or invite members.
  Names containing a forward slash are retained as one encoded API path segment.
  Refresh compares effective permissions with their previous observed value to
  detect drift while preserving configured permission categories. Imported roles
  may require an apply to establish that configured representation.

`vastai_team` creates a separate team with zero credit transfer and generates
an independent management key in sensitive state. `vastai_api_key` generates
a project key using its configured provider, `name`, and JSON `permissions`;
its computed `key` is sensitive. Key name, permission, and team-binding changes
require replacement. Optional `team_api_key` is reserved for a management key
generated during the same bootstrap apply. Ordinary resources omit it. Removing
that override verifies key visibility through the provider and updates only
state metadata, preserving the generated key ID and value. A new explicit
bootstrap issuer retains replacement behavior because it can select another team.
An empty permission string omits the permission field and requests default
full access; `{}` is not equivalent. Optional `team_id` resolves an existing
owned team's owner context, retains the new key in state, and attempts binding
only on that generated key. It requires empty permissions for explicit full
access and rejects scoped input before key creation. Computed bindings may require
read-only team refresh before validation. The issuer is never switched. Binding failure
fails apply before downstream Vault storage and leaves the key recoverable by
destruction. Offline fixtures verify this behavior; the live alwaldend attempt
created a key but received HTTP 400 from context selection, so that binding
has not passed live acceptance. The failed key was revoked through its root.
Their generated credentials cannot be recovered through import. API-key
imports remain unsupported. Use a secure backend and pass generated keys
directly to your secret manager rather than Terraform outputs.

An existing owned team can be imported by numeric account ID. Import and
refresh read the personal account's owner-context metadata without generating
keys or changing credential context. Imported teams have null `api_key` and
`management_key_id`; project keys can instead use the configured provider and
the imported team ID.
Their destruction is refused while they remain in state; use `prevent_destroy`
in their configuration and an explicit state-removal workflow to detach them.
Confirmed external absence removes an imported identity without key cleanup.
Import failures and refresh errors preserve account resources.

```sh
terraform import vastai_team.existing 846045
```

Team deletion requires a private record of confirmed creation and uses the
captured account ID. Failed creation never adopts an existing team with the
same name. If a team is deleted externally, refresh retains its management key
and plans replacement; apply revokes that key before creating another team.
Failed key revocation retains state for an explicit retry.

The [shared team module](../tf_modules/vastai_team/README.md) composes
new or imported teams and keys and returns sensitive values for caller-owned
Vault storage. The [bootstrap root](../../infra/vastai/tf/README.md) uses the personal token
from infra/vastai to create the `test_` team and project key, then writes the key
directly into the project's Vault namespace. It retains the personal token's
existing context.

The [team test root](test/team/README.md) exercises API keys, an environment marker,
and an unassigned role
using the generated key and the project's own AppRole and injector. Its object
names begin with `test_`; its environment name uses uppercase
`TEST_` because that API uppercases identifiers. It uses isolated Vault-backed
state and contains no compute resources or team-ID guard.

## Build and develop

```sh
bazel build //projects/vastai_terraform_provider
```

Use a private Terraform CLI configuration with a development override pointing
at the directory containing the built `terraform-provider-vastai` executable:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/alwaldend/vastai" = "/absolute/path/to/provider/directory"
  }
  direct {}
}
```

Set `TF_CLI_CONFIG_FILE` to that private file. With a development override,
run `terraform validate`, `plan`, and `apply` directly; `terraform init` would
attempt Registry resolution for this unreleased provider. In configurations
with other providers, initialize those separately as Terraform requires.
The provider uses protocol 6 and requires Terraform 1.0 or newer.

Inject the API credential as `VAST_API_KEY` from your secret manager. The
optional sensitive provider attribute `api_key` overrides it; avoid literal
credentials in configuration. `endpoint` defaults to
`https://console.vast.ai/api/v0`. Custom endpoints require HTTPS, except
literal loopback HTTP for fixture testing. Requests have a 60-second timeout
and an 8 MiB response limit, reject redirects, and do not retry mutations.
Inspect the account after an ambiguous creation failure before retrying.

Sensitive attributes are hidden from normal CLI display but remain in state.
Secure Terraform state and saved plans. API response bodies and returned
instance credentials are never copied into diagnostics or persisted state.

```hcl
terraform {
  required_providers {
    vastai = { source = "alwaldend/vastai" }
  }
}

provider "vastai" {}

resource "vastai_instance" "worker" {
  offer_id = 123456 # Select a currently available on-demand offer.
  image    = "ubuntu:24.04"
  disk_gb  = 20
  label    = "terraform-worker"
}
```

The [checked-in example](examples/gpu_rental/main.tf) separates offer discovery
from explicit offer selection so refreshing volatile search results does not
replace an existing rental.

## Import

```sh
terraform import vastai_instance.worker 12345678
```

Import accepts a canonical positive decimal contract ID. The API does not
return the original offer ID. Before applying imported resources, supply a
placeholder positive `offer_id` and add
`lifecycle { ignore_changes = [offer_id] }`. Reconcile `image`, `disk_gb`,
`runtype`, `onstart`, `label`, and `state` with the imported state. Review the
plan carefully; creation-only differences require replacement. When replacing
an imported instance deliberately, remove the ignore rule and select a valid
current offer first.

## Validation and ownership

```sh
bazel test //projects/vastai_terraform_provider/test/acceptance:acceptance_test
```

The acceptance suite runs the repository-pinned Terraform CLI and built
provider against a loopback fixture. It exercises the full lifecycle,
replacement plans, import, drift, remote disappearance, failure diagnostics,
and credential-safe redirect handling. `lifecycle.json` in test undeclared
outputs records the observed non-sensitive API operations. It requires no
Vast.ai account, network access outside loopback, or paid instances; live
compute compatibility remains unverified.

On 2026-11-04, the authorized live deployment applied the project AppRole,
team/token bootstrap, and the three-resource team test root. Both Vast.ai
Terraform roots then produced no-change plans using their own credentials.
No compute was rented or credit transferred; the personal bootstrap context
remained unchanged.

`cmd/terraform-provider-vastai` owns plugin startup. `internal/provider` owns
provider configuration and registration; `internal/resources` owns resource
schemas, plans, state, and imports; `internal/datasources` owns offer queries.
`internal/providerdata` defines the shared `Data` payload and checks it during
resource and data-source configuration. The shared `internal/vastai_sdk` owns generated endpoint operations,
bounded transport, and API response contracts.
`openspec` contains the maintained implementation record.

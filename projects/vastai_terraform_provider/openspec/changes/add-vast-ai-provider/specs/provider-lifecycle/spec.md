## Purpose

Allow Terraform users to discover available Vast.ai GPU offers and manage
on-demand instances through a reproducible declarative lifecycle.

## ADDED Requirements

### Requirement: Authenticated bounded requests

The provider SHALL accept a sensitive API key or VAST_API_KEY, require a
nonempty credential, use bearer authentication, bound HTTP operations, refuse
redirects, and avoid disclosing API response bodies or credentials in errors.

#### Scenario: Unauthorized or failed API call

- **WHEN** the API returns an authentication failure or success=false
- **THEN** Terraform reports a diagnostic without treating it as a missing instance.

### Requirement: Offer discovery

The provider SHALL expose a data source accepting a JSON object query and
return typed offer IDs, GPU details, and hourly prices in API order.

#### Scenario: Search and empty results

- **WHEN** a valid query returns zero or more offers
- **THEN** the data source returns the corresponding typed list without renting.

### Requirement: Instance lifecycle

The provider SHALL create an instance from an offer, image, and disk size,
retain its contract ID, refresh observable status and connection fields,
update label and desired state, replace creation-only settings, import a
positive decimal instance ID, and destroy managed instances.

#### Scenario: Full Terraform lifecycle

- **WHEN** an operator applies, updates, refreshes, imports, and destroys an instance
- **THEN** Terraform state and the corresponding API operations agree.

#### Scenario: Instance no longer exists

- **WHEN** refresh or destroy receives HTTP 404, or refresh receives a
  successful response with an explicitly null `instances` field
- **THEN** the resource is removed from Terraform state.

#### Scenario: Incomplete or mismatched successful response

- **WHEN** refresh receives a response with an omitted `instances` field or
  an object whose ID differs from the managed contract
- **THEN** refresh reports an error and retains resource ownership in state.

#### Scenario: Incomplete creation metadata

- **WHEN** refresh returns a matching instance ID but missing or invalid image or disk metadata
- **THEN** it reports a read error and retains the prior image and disk values rather than planning replacement.

#### Scenario: Configured creation-only change

- **WHEN** the image or disk allocation changes
- **THEN** Terraform plans replacement instead of an unsupported update.

### Requirement: Offline verification and usage documentation

The project SHALL provide Bazel build targets, examples, a landing page,
and an offline acceptance suite using the real Terraform CLI with repeatable
evidence. Documentation SHALL state the initial API coverage and asynchronous
readiness semantics.

#### Scenario: Isolated validation

- **WHEN** the acceptance suite runs
- **THEN** it exercises the provider with a local fixture and produces a report
  without contacting Vast.ai or provisioning live resources.

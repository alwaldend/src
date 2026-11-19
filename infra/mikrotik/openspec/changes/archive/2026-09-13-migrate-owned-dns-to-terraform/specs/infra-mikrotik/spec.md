## MODIFIED Requirements

### Requirement: Preserve export snapshots as the package's documented interface

The package SHALL retain the router export files as documentation inputs and
describe the RouterOS `/export` collection workflow. Documentation SHALL
distinguish these snapshots from automatically applied desired state. The
Terraform root SHALL manage DNS records only and SHALL NOT apply router exports.

#### Scenario: Inspect the router package

- **WHEN** a reader opens the component documentation
- **THEN** the reader can inspect both router export snapshots
- **AND** the Terraform package exposes owner-local DNS management while the
  snapshots remain documentary inputs

## ADDED Requirements

### Requirement: Standard Terraform validation test

The Terraform test map SHALL expose an optional standard validation test for
backend-free configuration inputs.

#### Scenario: Request the standard checks

- **WHEN** a caller enables validation with the default test map
- **THEN** the map declares the format and validation tests under the caller's prefix.

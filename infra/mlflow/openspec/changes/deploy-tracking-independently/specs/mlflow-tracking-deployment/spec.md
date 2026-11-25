## ADDED Requirements

### Requirement: Deploy tracking before worker enrollment

The default MLflow deployment SHALL use the Terraform-managed Vault OIDC credentials and start tracking without requiring Vast.ai or workload inputs.

#### Scenario: Tracking without a coordinator

- **WHEN** An operator runs the default MLflow guest deployment before registering a rental
- **THEN** The tracking service starts with Vault OIDC and the coordinator remains disabled

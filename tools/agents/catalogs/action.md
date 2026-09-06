# Action catalog

> Generated deterministic projection. The JSON document at `tools/agents/catalogs/action.json` is authoritative.

- ID: `agent-system.action`
- Schema: `agents.alwaldend.com/catalog/v1alpha1/action-catalog`
- Derivation: `1.0.0`
- Source revision: `c96e9118d7beed918a968be0da0be2d8`
- Completeness: `complete`
- JSON digest: `sha256:f9fb833e98ca04ffb9e49adf38f5bb97a3e044cd08ecd0273d22351ba287170f`

## Providers

- `cordis.runtime` (owned by `projects.mcp-cordis`) — projects/mcp_cordis/internal/mcp.mjs
- `github.forge` (owned by `tools.repo-delivery`) — tools/repo_delivery/cmd/repo_delivery/command.go
- `goal.local-store` (owned by `projects.goal`) — projects/goal/cmd/goal/command.go
- `terraform.runner` (owned by `tools.terraform`) — tools/terraform/defs.bzl

## Actions

- `cordis.define` (cordis.runtime.cordis_define): classified
- `cordis.promote` (cordis.runtime.cordis_promote): classified
- `cordis.remove` (cordis.runtime.cordis_remove): classified
- `goal.checkpoint` (goal.local-store.checkpoint): classified
- `goal.migrate` (goal.local-store.migrate): classified
- `goal.promote` (goal.local-store.promote): classified
- `repo-delivery.continue` (github.forge.continue): classified
- `repo-delivery.prepare` (github.forge.prepare): classified
- `repo-delivery.publish` (github.forge.publish): classified
- `repo-delivery.review-mutate` (github.forge.review): classified
- `repo-delivery.validate` (github.forge.validate): classified
- `terraform.apply` (terraform.runner.apply): classified
- `terraform.deploy` (terraform.runner.deploy): classified
- `terraform.deploy-y` (terraform.runner.deploy_y): classified
- `terraform.destroy` (terraform.runner.destroy): classified
- `terraform.direct` (terraform.runner.direct): requires_migration
- `terraform.fmt` (terraform.runner.fmt): classified
- `terraform.fmt-check` (terraform.runner.fmt_check): classified
- `terraform.force-unlock` (terraform.runner.force_unlock): classified
- `terraform.import` (terraform.runner.import): classified
- `terraform.init` (terraform.runner.init): classified
- `terraform.migrate` (terraform.runner.migrate): classified
- `terraform.output` (terraform.runner.output): classified
- `terraform.plan` (terraform.runner.plan): classified
- `terraform.show` (terraform.runner.show): classified
- `terraform.state` (terraform.runner.state): requires_migration

## Aliases

None.

## Limitations

None.

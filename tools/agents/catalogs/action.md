# Action catalog

> Generated deterministic projection. The JSON document at `tools/agents/catalogs/action.json` is authoritative.

- ID: `agent-system.action`
- Schema: `agents.alwaldend.com/catalog/v1alpha1/action-catalog`
- Derivation: `1.0.0`
- Source revision: `7aa44af2a240a52814308e26898ea9ff`
- Completeness: `complete`
- JSON digest: `sha256:1ec39e7ce0f671e2b22ffec83ed8238ccaa0fc3b2efa39a7adcd26543e05dba5`

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
- `repo-delivery.prepare` (github.forge.prepare): classified
- `repo-delivery.publish` (github.forge.publish): classified
- `repo-delivery.review-mutate` (github.forge.review): classified
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

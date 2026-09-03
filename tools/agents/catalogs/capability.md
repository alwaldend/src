# Capability catalog

> Generated deterministic projection. The JSON document at `tools/agents/catalogs/capability.json` is authoritative.

- ID: `agent-system.capability`
- Schema: `agents.alwaldend.com/catalog/v1alpha1/capability-catalog`
- Derivation: `1.0.0`
- Producer: `repository.capability-compiler`
- Source revision: `6105344746529abebfe130dda440ea5f`
- Completeness: `complete`
- JSON digest: `sha256:e5f4733ddba410341d97799b182e7c9182296a97297204157162248092a00ff0`

## Limitations

None.

## Providers

- `agent-system` (direct_binary, owned by `tools/agents`) — tools/agents/cmd/agent_system/main.go
- `bazel-agent` (direct_binary, owned by `projects/bazel_agent`) — projects/bazel_agent/cmd/bazel_agent/main.go
- `control-status` (direct_binary, owned by `tools/agents`) — tools/agents/cmd/control_status/main.go
- `cordis.runtime` (operation_provider, owned by `projects.mcp-cordis`) — projects/mcp_cordis/internal/mcp.mjs
- `cordis_define` (runtime_tool, owned by `projects/mcp_cordis`) — projects/mcp_cordis/internal/mcp.mjs
- `cordis_inspect` (runtime_tool, owned by `projects/mcp_cordis`) — projects/mcp_cordis/internal/mcp.mjs
- `cordis_invoke` (runtime_tool, owned by `projects/mcp_cordis`) — projects/mcp_cordis/internal/mcp.mjs
- `cordis_list` (runtime_tool, owned by `projects/mcp_cordis`) — projects/mcp_cordis/internal/mcp.mjs
- `cordis_list_tools` (runtime_tool, owned by `projects/mcp_cordis`) — projects/mcp_cordis/internal/mcp.mjs
- `cordis_promote` (runtime_tool, owned by `projects/mcp_cordis`) — projects/mcp_cordis/internal/mcp.mjs
- `cordis_reload` (runtime_tool, owned by `projects/mcp_cordis`) — projects/mcp_cordis/internal/mcp.mjs
- `cordis_remove` (runtime_tool, owned by `projects/mcp_cordis`) — projects/mcp_cordis/internal/mcp.mjs
- `cordis_run` (runtime_tool, owned by `projects/mcp_cordis`) — projects/mcp_cordis/internal/mcp.mjs
- `cordis_stop` (runtime_tool, owned by `projects/mcp_cordis`) — projects/mcp_cordis/internal/mcp.mjs
- `github.forge` (operation_provider, owned by `tools.repo-delivery`) — tools/repo_delivery/cmd/repo_delivery/command.go
- `goal` (direct_binary, owned by `projects/goal`) — projects/goal/cmd/goal/main.go
- `goal.local-store` (operation_provider, owned by `projects.goal`) — projects/goal/cmd/goal/command.go
- `mcp-cordis` (direct_binary, owned by `projects/mcp_cordis`) — projects/mcp_cordis/cmd/mcp_cordis/main.mjs
- `repo-delivery` (direct_binary, owned by `tools/repo_delivery`) — tools/repo_delivery/cmd/repo_delivery/main.go
- `terraform.runner` (operation_provider, owned by `tools.terraform`) — tools/terraform/defs.bzl

## Skills

- `agent-ergonomics-review` (owned by `projects/agents`): layer `review`, activation `substantial task close or recurring agent friction`, cost `medium`
  - exclusions: routine single-turn responses
  - capabilities: source.read
- `answer-question` (owned by `projects/agents`): layer `procedure`, activation `substantive user questions`, cost `medium`
  - exclusions: inert quoted questions
  - capabilities: source.read
- `bazel-agent` (owned by `projects/agents`): layer `execution`, activation `agent-executed Bazel commands`, cost `small`
  - exclusions: non-Bazel host commands
  - capabilities: code.execute
- `ast-grep` (owned by `projects/agents`): layer `procedure`, activation `structure-aware code search or rewrite`, cost `medium`
  - exclusions: plain-text searches expressible with rg
  - capabilities: code.execute
  - dependencies: bazel-agent
- `blender-reference-fidelity` (owned by `projects/agents`): layer `procedure`, activation `reference-controlled Blender modeling`, cost `large`
  - exclusions: freeform designs
  - capabilities: code.execute, source.write
- `codex-migration` (owned by `projects/agents`): layer `migration`, activation `Codex model, provider, authentication, or shared configuration migration`, cost `large`
  - exclusions: generic application configuration changes
  - capabilities: credential.consume, network.read
- `decision-review` (owned by `projects/agents`): layer `review`, activation `material decisions`, cost `medium`
  - exclusions: routine reversible choices
  - capabilities: source.read
- `full-repo-check` (owned by `projects/agents`): layer `validation`, activation `complete repository health checks`, cost `large`
  - exclusions: ordinary package checks
  - capabilities: code.execute
  - dependencies: bazel-agent
- `git-rebase-remote` (owned by `projects/agents`): layer `delivery`, activation `task-owned feature rebase`, cost `medium`
  - exclusions: shared or human-owned history
  - capabilities: history.write, network.read, remote.write
- `goal` (owned by `projects/goal`): layer `coordination`, activation `durable multi-step work`, cost `large`
  - exclusions: simple one-step tasks
  - capabilities: task_state.write
  - dependencies: bazel-agent
- `host-bot-diagnostics` (owned by `projects/agents`): layer `diagnostics`, activation `host bot audits`, cost `medium`
  - exclusions: infrastructure implementation
  - capabilities: network.read, source.read
- `project-layout` (owned by `projects/agents`): layer `policy`, activation `source placement and moves`, cost `small`
  - exclusions: unchanged existing layout
  - capabilities: source.read
- `repo-bazel` (owned by `projects/agents`): layer `procedure`, activation `Bazel graph changes and validation`, cost `medium`
  - exclusions: non-Bazel source-only edits
  - capabilities: code.execute, source.write
  - dependencies: bazel-agent
- `bazel-nested-module` (owned by `projects/agents`): layer `procedure`, activation `standalone nested Bzlmod modules`, cost `medium`
  - exclusions: ordinary Bazel packages
  - capabilities: source.write
  - dependencies: bazel-agent, repo-bazel
- `bazel-rules-skill` (owned by `projects/agents`): layer `procedure`, activation `repository skill packaging`, cost `medium`
  - exclusions: personal untracked skills
  - capabilities: source.write
  - dependencies: bazel-agent, repo-bazel
- `repo-ansible` (owned by `projects/agents`): layer `procedure`, activation `repository Ansible changes`, cost `medium`
  - exclusions: non-Ansible infrastructure
  - capabilities: code.execute, source.write
  - dependencies: bazel-agent, repo-bazel
- `repo-blender` (owned by `projects/agents`): layer `procedure`, activation `repository Blender assets and renders`, cost `large`
  - exclusions: non-Blender image edits
  - capabilities: code.execute, source.write
  - dependencies: bazel-agent, repo-bazel
- `repo-delivery` (owned by `tools/repo_delivery`): layer `delivery`, activation `feature publication and PR review`, cost `large`
  - exclusions: read-only reviews
  - capabilities: history.write, network.read, remote.write
  - dependencies: bazel-agent
- `repo-external-dependency` (owned by `projects/agents`): layer `procedure`, activation `external dependency pins`, cost `medium`
  - exclusions: first-party versioning
  - capabilities: network.read, source.write
  - dependencies: bazel-agent, repo-bazel
- `repo-gazelle-plugin` (owned by `projects/agents`): layer `procedure`, activation `Gazelle language plugins`, cost `medium`
  - exclusions: ordinary generated BUILD updates
  - capabilities: code.execute, source.write
  - dependencies: bazel-agent, repo-bazel
- `repo-secrets` (owned by `projects/agents`): layer `security`, activation `credentials and secret routing`, cost `large`
  - exclusions: ordinary public configuration
  - capabilities: credential.consume, source.write
- `repo-terraform` (owned by `projects/agents`): layer `procedure`, activation `repository Terraform changes`, cost `large`
  - exclusions: non-Terraform infrastructure
  - capabilities: code.execute, source.write
  - dependencies: bazel-agent, repo-bazel
- `spellcheck` (owned by `projects/agents`): layer `procedure`, activation `prose proofreading and rewriting`, cost `small`
  - exclusions: code linting, translation-only requests
  - capabilities: source.read, source.write
- `versioning` (owned by `tools/versioning`): layer `release`, activation `first-party repository releases`, cost `medium`
  - exclusions: third-party dependency versions
  - capabilities: history.write, source.read
  - dependencies: bazel-agent

# Agent system index

> Generated deterministic projection. The JSON document at `tools/agents/catalogs/index.json` is authoritative.

- ID: `agent-system.index`
- Schema: `agents.alwaldend.com/catalog/v1alpha1/agent-system-index`
- Derivation: `1.0.0`
- Source revision: `76b4b430b1369514f21497dd3b838472`
- Completeness: `complete`
- JSON digest: `sha256:9cff493fdf12a4980542a2b6179d7354aa78e1afb8f4f377a616281e1efca608`

## Catalogs

- `agent-system.action` (action-catalog): complete, digest `sha256:1ec39e7ce0f671e2b22ffec83ed8238ccaa0fc3b2efa39a7adcd26543e05dba5`
  - route: `compile:action-catalog`
  - route: `check:action-catalog`
- `agent-system.capability` (capability-catalog): partial, digest `sha256:8849f080a74eef0cc214e7d828718e6dc6e013b3a7732fc8ac5fb0835120f3f2`
  - route: `compile:capability-catalog`
  - route: `check:capability-catalog`
- `agent-system.goal` (goal-catalog): complete, digest `sha256:a66567353b89e0dd49324055f2628ba3cccbc2fbb192b6a1e962a3a8ecd2608c`
  - route: `compile:goal-catalog`
  - route: `check:goal-catalog`
- `agent-system.policy` (policy-catalog): complete, digest `sha256:0b64de26a8fd013756b0ea5dfb768ff05492d4f752491cc57a49170ea7a3dd0e`
  - route: `compile:policy-catalog`
  - route: `check:policy-catalog`
- `agent-system.topology` (topology-catalog): complete, digest `sha256:24cfd78b69ac9e4821cc1ac27a4a3ad36be99a583b65ed2c090842256cd2c3b3`
  - route: `compile:topology-catalog`
  - route: `check:topology-catalog`
- `agent-system.workspace-check` (workspace-check-catalog): complete, digest `sha256:40d46b3b913dc607a7c55a9b0b499cd4d500fdc99a779a5b29843ef4bcb696c0`
  - route: `compile:workspace-check-catalog`
  - route: `check:workspace-check-catalog`

## Conflicts

- `agent-system.policy.policy.agents` (policy_axis_conflict): AGENTS.md

## Limitations

None.

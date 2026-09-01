# Agent system index

> Generated deterministic projection. The JSON document at `tools/agents/catalogs/index.json` is authoritative.

- ID: `agent-system.index`
- Schema: `agents.alwaldend.com/catalog/v1alpha1/agent-system-index`
- Derivation: `1.0.0`
- Source revision: `38c4c174d3b28e53bcd9a0e026f3b8ec`
- Completeness: `complete`
- JSON digest: `sha256:51515b92138c81929a888269154e856370cd906bcc45afe1b6f55cec41758fa7`

## Catalogs

- `agent-system.action` (action-catalog): complete, digest `sha256:bea5d306989cf40af723160b48ccdf306f66037e087052ed2080f297890990d6`
  - route: `compile:action-catalog`
  - route: `check:action-catalog`
- `agent-system.capability` (capability-catalog): partial, digest `sha256:a376d14ba0a6960106bce790e44c06cc04b9d4b7bb21d60b0bf419ac09629076`
  - route: `compile:capability-catalog`
  - route: `check:capability-catalog`
- `agent-system.goal` (goal-catalog): complete, digest `sha256:635221e76fc2001cad2ee39228475da6e95f296349d3ad8e54b94ae40fe7acb2`
  - route: `compile:goal-catalog`
  - route: `check:goal-catalog`
- `agent-system.policy` (policy-catalog): complete, digest `sha256:60de6ee935594334ad954d638174b42bf25b0578e85099f696e34d8b967ea840`
  - route: `compile:policy-catalog`
  - route: `check:policy-catalog`
- `agent-system.topology` (topology-catalog): complete, digest `sha256:d288102cbdd5a5a84f6bfb7a4ced8ed89d4f5b69e51b7c0661abca42e930c7af`
  - route: `compile:topology-catalog`
  - route: `check:topology-catalog`
- `agent-system.workspace-check` (workspace-check-catalog): complete, digest `sha256:2b5fc640a81d1fb57659c2c3b601f3dc284da83de44c72cd1bc616cab39fbaa8`
  - route: `compile:workspace-check-catalog`
  - route: `check:workspace-check-catalog`

## Conflicts

- `agent-system.policy.policy.agents` (policy_axis_conflict): AGENTS.md

## Limitations

None.

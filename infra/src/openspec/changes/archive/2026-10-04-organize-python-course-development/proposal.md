## Why

PR #124 adds an MLOps course placeholder and prepares shared Python dependencies.
The user requested that dependency ownership move to tools/py while shared tool
settings remain at root. This record was created after the initial implementation
at revision 0fffe5b489f887020c1b11a4e815f00f14a961ff; it does not claim earlier planning.

## What Changes

- Add projects/mlops_2026 as an MLOps course placeholder with the supplied Yandex
  Practicum link and a documentation target, excluded from release deployment.
- Add explicitly requested FastAPI and uvicorn[standard] dependencies.
- Keep only tool settings in root pyproject.toml; move project metadata and
  dependencies to tools/py/pyproject.toml and the lock to tools/py/requirements.txt.
- Update Bazel consumers and regenerate dependency and Gazelle mappings.
- Require an owner-local OpenSpec change for every repository change in AGENTS.md,
  including trivial changes, while keeping artifact size proportional.
- Retain gazelle_python.yaml at root for default discovery; document the directive
  that permits relocation. No manifest relocation was requested.

## Capabilities

### New Capabilities

None. The course is a documentation placeholder without executable behavior.

### Modified Capabilities

- repository: Require OpenSpec records for every repository change. Python
  dependency preparation and ownership refactoring preserve existing checks.

## Impact

Affected owners are projects/mlops_2026, tools/py, tools/gazelle, and root shared
configuration. The linked agent-skill change is
../../../../../../tools/agents/openspec/changes/archive/2026-10-04-add-python-and-change-review-skills/.
Acceptance requires correct source ownership, consistent generated files, existing
checker behavior, repository quality, and semantic lint on the delivered candidate.

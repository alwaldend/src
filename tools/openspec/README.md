# OpenSpec

`openspec` is the pinned Bazel wrapper for the OpenSpec CLI. Its generated
agent skills are packaged from the pinned release archive declared in
`third_party/org_fissionai_openspec/` and written into `.agents/skills/` by
`//.agents:write_skills`.

The repository-specific [al-openspec skill](skills/al-openspec/SKILL.md) owns
repository integration: owner-local workspaces, the pinned CLI invocation, and
routing to the upstream operation skills. Its
discovery link is generated alongside the upstream skills by
`//.agents:write_skills`.

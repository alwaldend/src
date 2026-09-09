# OpenSpec

The pinned OpenSpec CLI wrapper lives at `//tools/openspec`. Its generated
`openspec-*` agent skills are materialized from the pinned upstream release
archive declared in `third_party/org_fissionai_openspec/`; the resulting
`.agents/skills/openspec-*` entries are written as regular files by
`//.agents:write_skills` and verified by `//.agents:write_skills_test`.

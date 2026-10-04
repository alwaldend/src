---
title: OpenAI Codex review agent
description: Pinned upstream read-only code review skill
---

The [review-agent skill](https://github.com/openai/codex/blob/afb436df8b70bb5bc57b86d9a3e829968988cd21/codex-rs/skills/src/assets/samples/review-agent/SKILL.md)
is consumed unchanged from OpenAI Codex revision `afb436df8b70bb5bc57b86d9a3e829968988cd21`.
`include.MODULE.bazel` pins the source archive and its SHA-256 integrity.
Only the review-agent directory is packaged; the Codex executable is not used.

The source is licensed under [Apache 2.0](LICENSE). The upstream [NOTICE](NOTICE)
is retained. `.agents/skills/review-agent/` is a generated projection of the
upstream skill, including its explicit-invocation metadata. Run
`//.agents:write_skills` after changing the pin; do not edit the projection.
Repository routing in `AGENTS.md` requires this skill for change reviews.

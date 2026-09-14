---
title: Repository CI skill evaluations
---

The cases cover thin workflow extraction, reuse of the owning job command,
runner versus Vault trust boundaries, runtime credential handling, bounded
authentication diagnosis, and verification of an authorized exact candidate.
Skill-routing assertions are separate from behavioral rubrics.

The `eval_config_test` target validates the Promptfoo configuration, referenced
cases, and staged skill offline. It makes no model calls and does not prove
that an agent follows the instructions.

A live evaluation target is omitted. Representative behavior requires an
isolated repository with workflow and Bazel sources, a tool-capable agent,
and controlled Forgejo/Vault fixtures or explicitly authorized live services.
The staged skill alone does not provide that execution surface. Use such a
fixture before claiming behavioral coverage; this configuration records the
intended cases without adding a credentialed test to ordinary validation.

---
title: Repository Python skill evaluations
---

The offline target validates the Promptfoo configuration, cases, and staged
skill without model calls. Cases cover mandatory annotations, Python support
boundaries, namespaced imports, extras, dependency locking, and Gazelle regeneration.

Live evaluation is omitted because representative behavior requires a
fixture checkout and tool access to the pinned Bazel generators and Python
checks. Offline validation does not prove that an agent follows the skill.

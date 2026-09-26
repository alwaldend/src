---
title: Repository Go skill evaluations
---

The offline target validates the Promptfoo configuration, cases, and staged
skill without model calls. Cases exercise error causes, command boundaries,
local dependency discovery, and repository Go tooling.

A live target is omitted because representative evaluation needs a tool-capable
fixture repository, Bazel, and an isolated SSH endpoint. Offline validation
does not establish that an agent follows the skill. The publication change
uses its separate command-level SSH/HTTP fixture for implementation evidence.

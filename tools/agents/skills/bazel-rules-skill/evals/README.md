---
title: Bazel rules skill evaluations
---

# Bazel rules skill evaluations

This suite describes the required behavior for adding a repository skill. Its
offline Bazel target validates the Promptfoo configuration, referenced cases,
and skill staging without making a model call.

A live target is omitted because representative behavior requires creating
files and invoking Bazel. Configuration validation does not establish that an
agent follows the workflow; use a tool-capable, isolated fixture repository
before treating this as behavioral coverage.

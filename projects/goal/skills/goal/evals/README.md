---
title: Disabled goal skill evaluations
---

These cases check that the disabled skill directs new work to OpenSpec while
preserving the legacy CLI's record-inspection use. The offline
`eval_config_test` validates configuration, referenced cases, and skill
packaging without credentials or model calls. It does not establish the
model's routing behavior.

A live target is omitted because meaningful migration and continuation
checks need an isolated writable workspace and filesystem tools. The
retained API, store, and CLI tests cover legacy compatibility independently.

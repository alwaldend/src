---
title: Anki as code
linkTitle: Anki as code
description: A CLI project for managing Anki collections through editable TOML
layout: landing
statuses:
  - in_progress
languages:
  - go
tags:
  - cli
---

Anki as code aims to let you edit notes, decks, and card appearance in text
files, inspect proposed changes, and reconcile an existing Anki collection.

The project currently provides a CLI scaffold with export, plan, apply, build,
and generate-id commands. Collection operations are still placeholders and
return explicit errors. They do not change your collection.

[Read the documentation](/docs/projects/anki_as_code/)

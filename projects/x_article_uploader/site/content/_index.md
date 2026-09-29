---
title: X Article uploader
linkTitle: X Article uploader
description: Create an X Article draft from a Markdown blog post
layout: landing
statuses:
  - active
languages:
  - go
tags:
  - x
  - markdown
---

The X Article uploader creates an X Article draft from a blog post. It compiles
a post's Markdown body into the DraftJS `content_state` the Articles endpoint
accepts, then creates a draft from it; it never publishes.

## Features

- Offline, deterministic Markdown-to-`content_state` conversion
- Drafts only: no invocation can publish an article
- Images resolved through the media upload endpoint by content digest
- OAuth 1.0a credentials injected through the repository's Vault flow

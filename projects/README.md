---
title: Projects
description: Project tree
weight: 1
cascade:
  - categories:
      - project
---

This tree contains first-party products and reusable project code. All tracked
content follows the repository's public-source policy.

Each project owns its specifications and maintained changes in
`<project>/openspec/`. Use the [pinned OpenSpec workflow](../tools/openspec/README.md)
with that project selected. [`infra/src/`](../infra/src/README.md) owns repository evolution.

## Requirements

- Bazel targets MAY use public visibility when their owner intends external
  reuse.
- Project artifacts MAY be published through an explicit release workflow.
- Project targets MAY be used by production build targets.

## Project sites

Each project has a landing page in the main website at
`https://alwaldend.com/projects/<name>/`. Sites are published with the main
site rather than as separate Pages repositories, and repository reference
documentation, including each project README, lives under `/docs/`.

| Project                                                                    | Description                                                              | Landing page                                                                                                             |
| -------------------------------------------------------------------------- | ------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------ |
| [ActivityWatch ingester Android](activitywatch_ingester_android/README.md) | Android accessibility collector that sends browser tabs to ActivityWatch | [alwaldend.com/projects/activitywatch_ingester_android/](https://alwaldend.com/projects/activitywatch_ingester_android/) |
| [Agents](agents/README.md)                                                 | Repository-wide reusable agent skills                                    | [alwaldend.com/projects/agents/](https://alwaldend.com/projects/agents/)                                                 |
| [Al](al/README.md)                                                         | Repository command runner and Bazel configuration rules                  | [alwaldend.com/projects/al/](https://alwaldend.com/projects/al/)                                                         |
| [Alwaldend.com](alwaldend.com/README.md)                                   | Main website and project documentation                                   | [alwaldend.com](https://alwaldend.com/)                                                                                  |
| [Android launcher](android_launcher/README.md)                             | Text-only Android launcher built with Kotlin and Jetpack Compose         | [alwaldend.com/projects/android_launcher/](https://alwaldend.com/projects/android_launcher/)                             |
| [Ansible collection](ansible_collection/README.md)                         | Ansible collection alwaldend.main with Bazel source packaging            | [alwaldend.com/projects/ansible_collection/](https://alwaldend.com/projects/ansible_collection/)                         |
| [Autoscroll](autoscroll/README.md)                                         | Mouse-driven autoscroll CLI with reloadable configuration                | [alwaldend.com/projects/autoscroll/](https://alwaldend.com/projects/autoscroll/)                                         |
| [Bazel agent](bazel_agent/README.md)                                       | Bazel runner for repository agents                                       | [alwaldend.com/projects/bazel_agent/](https://alwaldend.com/projects/bazel_agent/)                                       |
| [Ci platform](ci_platform/README.md)                                       | Abandoned CI platform with a Go backend and Vue frontend                 | [alwaldend.com/projects/ci_platform/](https://alwaldend.com/projects/ci_platform/)                                       |
| [Dotfiles](dotfiles/README.md)                                             | Personal configuration files with installation and comparison commands   | [alwaldend.com/projects/dotfiles/](https://alwaldend.com/projects/dotfiles/)                                             |
| [Infinitime](infinitime/README.md)                                         | InfiniTime firmware fork with a text watchface and Pomodoro app          | [alwaldend.com/projects/infinitime/](https://alwaldend.com/projects/infinitime/)                                         |
| [Kustomization](kustomization/README.md)                                   | Kubernetes resources for Flux, Traefik, and cert-manager                 | [alwaldend.com/projects/kustomization/](https://alwaldend.com/projects/kustomization/)                                   |
| [Leetcode downloader](leetcode_downloader/README.md)                       | LeetCode submission export and documentation tools                       | [alwaldend.com/projects/leetcode_downloader/](https://alwaldend.com/projects/leetcode_downloader/)                       |
| [MCP Cordis](mcp_cordis/README.md)                                         | Workspace-local runtime packages behind a stable MCP server              | [alwaldend.com/projects/mcp_cordis/](https://alwaldend.com/projects/mcp_cordis/)                                         |
| [Nexus security plugin](nexus_security_plugin/README.md)                   | Security plugin for Sonatype Nexus 3                                     | [alwaldend.com/projects/nexus_security_plugin/](https://alwaldend.com/projects/nexus_security_plugin/)                   |
| [Renders](renders/README.md)                                               | Repository-owned render assets and their acceptance evidence             | [alwaldend.com/projects/renders/](https://alwaldend.com/projects/renders/)                                               |
| [Sri](sri/README.md)                                                       | Command-line Subresource Integrity calculator using OpenSSL              | [alwaldend.com/projects/sri/](https://alwaldend.com/projects/sri/)                                                       |
| [Tf modules](tf_modules/README.md)                                         | Reusable Terraform modules for Vault, virtual machines, and storage      | [alwaldend.com/projects/tf_modules/](https://alwaldend.com/projects/tf_modules/)                                         |
| [Useless QT GUI](useless_qt_gui/README.md)                                 | Desktop GUI application built with C++ and Qt                            | [alwaldend.com/projects/useless_qt_gui/](https://alwaldend.com/projects/useless_qt_gui/)                                 |

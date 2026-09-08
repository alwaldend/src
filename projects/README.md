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

These configured site addresses use project names with underscores replaced
by hyphens. The `alwaldend.com` project uses the main website.

| Project documentation                                        | Description                                                                 | Landing page                                                                            |
| ------------------------------------------------------------ | --------------------------------------------------------------------------- | --------------------------------------------------------------------------------------- |
| [Agents](agents/README.md)                                   | Repository-wide agent-system architecture and skills                        | [agents.alwaldend.com](https://agents.alwaldend.com/)                                   |
| [Al](al/README.md)                                           | Repository command runner and Bazel configuration rules                     | [al.alwaldend.com](https://al.alwaldend.com/)                                           |
| [Alwaldend.com](alwaldend.com/README.md)                     | Main website and project documentation                                      | [alwaldend.com](https://alwaldend.com/)                                                 |
| [Android launcher](android_launcher/README.md)               | Text-only Android launcher built with Kotlin and Jetpack Compose            | [android-launcher.alwaldend.com](https://android-launcher.alwaldend.com/)               |
| [Ansible collection](ansible_collection/README.md)           | Ansible collection alwaldend.main with Bazel source packaging               | [ansible-collection.alwaldend.com](https://ansible-collection.alwaldend.com/)           |
| [Autoscroll](autoscroll/README.md)                           | Mouse-driven autoscroll CLI with reloadable configuration                   | [autoscroll.alwaldend.com](https://autoscroll.alwaldend.com/)                           |
| [Bazel agent](bazel_agent/README.md)                         | Bazel runner for repository agents                                          | [bazel-agent.alwaldend.com](https://bazel-agent.alwaldend.com/)                         |
| [Ci platform](ci_platform/README.md)                         | Abandoned CI platform with a Go backend and Vue frontend                    | [ci-platform.alwaldend.com](https://ci-platform.alwaldend.com/)                         |
| [Dotfiles](dotfiles/README.md)                               | Personal configuration files with installation and comparison commands      | [dotfiles.alwaldend.com](https://dotfiles.alwaldend.com/)                               |
| [Goal](goal/README.md)                                       | Deprecated goal-record compatibility tooling; maintained work uses OpenSpec | [goal.alwaldend.com](https://goal.alwaldend.com/)                                       |
| [Hugo Landing](hugo_landing/README.md)                       | Reusable Hugo landing site for project pages                                | [hugo-landing.alwaldend.com](https://hugo-landing.alwaldend.com/)                       |
| [Infinitime](infinitime/README.md)                           | InfiniTime firmware fork with a text watchface and Pomodoro app             | [infinitime.alwaldend.com](https://infinitime.alwaldend.com/)                           |
| [Kustomization](kustomization/README.md)                     | Kubernetes resources for Flux, Traefik, and cert-manager                    | [kustomization.alwaldend.com](https://kustomization.alwaldend.com/)                     |
| [Leetcode downloader](leetcode_downloader/README.md)         | LeetCode submission export and documentation tools                          | [leetcode-downloader.alwaldend.com](https://leetcode-downloader.alwaldend.com/)         |
| [MCP Cordis](mcp_cordis/README.md)                           | Workspace-local runtime packages behind a stable MCP server                 | [mcp-cordis.alwaldend.com](https://mcp-cordis.alwaldend.com/)                           |
| [Nexus security plugin](nexus_security_plugin/README.md)     | Security plugin for Sonatype Nexus 3                                        | [nexus-security-plugin.alwaldend.com](https://nexus-security-plugin.alwaldend.com/)     |
| [Renders](renders/README.md)                                 | Repository-owned render assets and their acceptance evidence                | [renders.alwaldend.com](https://renders.alwaldend.com/)                                 |
| [Rules binary toolchain](rules_binary_toolchain/README.md)   | Bazel toolchains for packaged executable binaries                           | [rules-binary-toolchain.alwaldend.com](https://rules-binary-toolchain.alwaldend.com/)   |
| [Rules Dnscontrol](rules_dnscontrol/README.md)               | Bazel-aware DNSControl configuration packaging                              | [rules-dnscontrol.alwaldend.com](https://rules-dnscontrol.alwaldend.com/)               |
| [Rules docs](rules_docs/README.md)                           | Bazel documentation packaging rules                                         | [rules-docs.alwaldend.com](https://rules-docs.alwaldend.com/)                           |
| [Rules docs Gazelle](rules_docs_gazelle/README.md)           | Gazelle extension for Bazel documentation packaging rules                   | [rules-docs-gazelle.alwaldend.com](https://rules-docs-gazelle.alwaldend.com/)           |
| [Rules Hugo](rules_hugo/README.md)                           | Bazel rules for Hugo sites                                                  | [rules-hugo.alwaldend.com](https://rules-hugo.alwaldend.com/)                           |
| [rules_iso](rules_iso/README.md)                             | ISO image download and flash rules                                          | [rules-iso.alwaldend.com](https://rules-iso.alwaldend.com/)                             |
| [Rules Promptfoo](rules_promptfoo/README.md)                 | A pinned Bazel runner for Promptfoo skill evaluations                       | [rules-promptfoo.alwaldend.com](https://rules-promptfoo.alwaldend.com/)                 |
| [Rules Promptfoo Gazelle](rules_promptfoo_gazelle/README.md) | Gazelle extension for offline Promptfoo validation tests                    | [rules-promptfoo-gazelle.alwaldend.com](https://rules-promptfoo-gazelle.alwaldend.com/) |
| [Rules skill](rules_skill/README.md)                         | Bazel rules and validation for Codex skills                                 | [rules-skill.alwaldend.com](https://rules-skill.alwaldend.com/)                         |
| [Rules skill Gazelle](rules_skill_gazelle/README.md)         | Gazelle extension for Bazel skill libraries                                 | [rules-skill-gazelle.alwaldend.com](https://rules-skill-gazelle.alwaldend.com/)         |
| [Rules template](rules_template/README.md)                   | Bazel rules and a Go command for rendering template files                   | [rules-template.alwaldend.com](https://rules-template.alwaldend.com/)                   |
| [Sri](sri/README.md)                                         | Command-line Subresource Integrity calculator using OpenSSL                 | [sri.alwaldend.com](https://sri.alwaldend.com/)                                         |
| [Tf modules](tf_modules/README.md)                           | Reusable Terraform modules for Vault, virtual machines, and storage         | [tf-modules.alwaldend.com](https://tf-modules.alwaldend.com/)                           |
| [Useless QT GUI](useless_qt_gui/README.md)                   | Desktop GUI application built with C++ and Qt                               | [useless-qt-gui.alwaldend.com](https://useless-qt-gui.alwaldend.com/)                   |

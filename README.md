---
title: Src
description: Source code
cascade:
  - type: docs
    _target:
      path: /**
---

## Links

- Docs: https://alwaldend.github.io/src/
- Dev release: https://github.com/alwaldend/src/releases/tag/dev
- Dev release artifacts: https://github.com/alwaldend/src/releases/download/dev/alwaldend-src.tar.gz

## Licence

AGPLv3, see [LICENSE](./LICENSE.txt)

## Contents

<details open>
  <summary>Contents</summary>

<!-- README_CONTENTS START -->

- <a href=".github">.github</a>: Github config

  - <a href=".github/workflows">workflows</a>: Github workflows

- <a href="ansible">ansible</a>: Ansible

  - <a href="ansible/collections">collections</a>: Collections

  - <a href="ansible/playbooks">playbooks</a>: Playbooks

  - <a href="ansible/roles">roles</a>: Roles

    - <a href="ansible/roles/adguard">adguard</a>: Adguard

    - <a href="ansible/roles/caddy">caddy</a>: Caddy

    - <a href="ansible/roles/cifs">cifs</a>: CIFS

    - <a href="ansible/roles/consul">consul</a>: Consul

    - <a href="ansible/roles/consul_envoy">consul_envoy</a>: Consul envoy

    - <a href="ansible/roles/dns">dns</a>: DNS

    - <a href="ansible/roles/docker">docker</a>: Docker

    - <a href="ansible/roles/facts">facts</a>: Facts

    - <a href="ansible/roles/firewall">firewall</a>: Firewall

    - <a href="ansible/roles/hiddify_manager">hiddify_manager</a>: Hiddify manager

    - <a href="ansible/roles/hiddify_manager_host">hiddify_manager_host</a>: Hiddify manager with host

    - <a href="ansible/roles/host">host</a>: Host

    - <a href="ansible/roles/k3s">k3s</a>: K3s

    - <a href="ansible/roles/k3s_bootstrap">k3s_bootstrap</a>: K3s bootstrap

    - <a href="ansible/roles/k3s_cluster">k3s_cluster</a>: K3s cluster

    - <a href="ansible/roles/os">os</a>: OS

    - <a href="ansible/roles/pve_cluster">pve_cluster</a>: PVE cluster

    - <a href="ansible/roles/pve_vm">pve_vm</a>: PVE VM

    - <a href="ansible/roles/pve_vm_remove">pve_vm_remove</a>: PVE VM remove

    - <a href="ansible/roles/ssh">ssh</a>: SSH

    - <a href="ansible/roles/ssh_port_forward">ssh_port_forward</a>: SSH port forwarding

    - <a href="ansible/roles/ssh_update_known_hosts">ssh_update_known_hosts</a>: SSH update known hosts

    - <a href="ansible/roles/traefik">traefik</a>: Traefik

    - <a href="ansible/roles/update_all_packages">update_all_packages</a>: Update all packages

    - <a href="ansible/roles/users">users</a>: Users

    - <a href="ansible/roles/vault">vault</a>: Vault

    - <a href="ansible/roles/wireguard">wireguard</a>: Wireguard

    - <a href="ansible/roles/xray">xray</a>: Xray

- <a href="bzl">bzl</a>: Bazel

  - <a href="bzl/aspects">aspects</a>: Aspects

  - <a href="bzl/build-files">build-files</a>: Build files

  - <a href="bzl/configs">configs</a>: Configs

  - <a href="bzl/extensions">extensions</a>: Extensions

  - <a href="bzl/macros">macros</a>: Macros

  - <a href="bzl/providers">providers</a>: Providers

  - <a href="bzl/registry">registry</a>: Registry

    - <a href="bzl/registry/modules">modules</a>: Modules

      - <a href="bzl/registry/modules/com-github-aldanial-cloc">com-github-aldanial-cloc</a>: Cloc

      - <a href="bzl/registry/modules/com-github-fortawesome-font-awesome">com-github-fortawesome-font-awesome</a>: Font-Awesome

      - <a href="bzl/registry/modules/com-github-georgewfraser-java-language-server">com-github-georgewfraser-java-language-server</a>: Java languange server

      - <a href="bzl/registry/modules/com-github-google-docsy">com-github-google-docsy</a>: Docsy

      - <a href="bzl/registry/modules/com-github-twbs-bootstrap">com-github-twbs-bootstrap</a>: Bootstrap

      - <a href="bzl/registry/modules/com-nordicsemi-developer-nrfsdk">com-nordicsemi-developer-nrfsdk</a>: Nrfsdk

      - <a href="bzl/registry/modules/hedron_compile_commands">hedron_compile_commands</a>: Bazel-compile-commands-extractor

      - <a href="bzl/registry/modules/org-openssl-openssl">org-openssl-openssl</a>: Openssl

      - <a href="bzl/registry/modules/rules_haskell">rules_haskell</a>: Patched version of rules_haskell

      - <a href="bzl/registry/modules/us-nasm-nasm">us-nasm-nasm</a>: Netwide Assembler (NASM)

  - <a href="bzl/rules">rules</a>: Rules

  - <a href="bzl/toolchain-types">toolchain-types</a>: Toolchain types

  - <a href="bzl/toolchains">toolchains</a>: Toolchains

  - <a href="bzl/vars">vars</a>: Vars

- <a href="c">c</a>: C

  - <a href="c/misc">misc</a>: Misc

  - <a href="c/openssl">openssl</a>: Openssl

  - <a href="c/sri">sri</a>: SRI

- <a href="cfg">cfg</a>: Configs

  - <a href="cfg/dotfiles">dotfiles</a>: Dotfiles

  - <a href="cfg/leetcode-submissions">leetcode-submissions</a>: Leetcode submissions

  - <a href="cfg/misc">misc</a>: Misc

  - <a href="cfg/vial">vial</a>: Vial

- <a href="cpp">cpp</a>: C&#43;&#43;

  - <a href="cpp/infinitime">infinitime</a>: Infinitime

  - <a href="cpp/useless-qt-gui">useless-qt-gui</a>: Useless qt GUI

- <a href="data">data</a>: Data

  - <a href="data/misc">misc</a>: Misc

- <a href="drawio">drawio</a>: Drawio

  - <a href="drawio/diagrams">diagrams</a>: Diagrams

- <a href="go">go</a>: Go

  - <a href="go/bazel-shell-worker">bazel-shell-worker</a>: Bazel shell worker

  - <a href="go/file-installer">file-installer</a>: File installer

    - <a href="go/file-installer/cmd">cmd</a>: Cmd code for file-installer

  - <a href="go/leetcode-downloader">leetcode-downloader</a>: Leetcode downloader

    - <a href="go/leetcode-downloader/model">model</a>: Models for leetcode-downloader

  - <a href="go/readme-tree">readme-tree</a>: Readme tree

  - <a href="go/template-files">template-files</a>: Template files

  - <a href="go/utils">utils</a>: Utils

- <a href="hugo">hugo</a>: Hugo

  - <a href="hugo/sites">sites</a>: Sites

    - <a href="hugo/sites/docs">docs</a>: Docs

      - <a href="hugo/sites/docs/content">content</a>: Content

      - <a href="hugo/sites/docs/layouts">layouts</a>: Layouts

  - <a href="hugo/themes">themes</a>: Themes

- <a href="img">img</a>: Images

  - <a href="img/useless-qt-gui">useless-qt-gui</a>: Useless qt GUI

- <a href="java">java</a>: Java

- <a href="js">js</a>: Javascript

  - <a href="js/leetcode-downloader">leetcode-downloader</a>: Leetcode downloader

- <a href="kt">kt</a>: Kotlin

- <a href="lua">lua</a>: Lua

  - <a href="lua/nvim-config">nvim-config</a>: Neovim config

  - <a href="lua/nvim-lib">nvim-lib</a>: Neovim lib

- <a href="patch">patch</a>: Patch

  - <a href="patch/infinitime">infinitime</a>: Infinitime

  - <a href="patch/rules-haskell">rules-haskell</a>: Rules Haskell

- <a href="pl">pl</a>: Perl

- <a href="proto">proto</a>: Proto

  - <a href="proto/bazel-worker">bazel-worker</a>: Bazel worker

  - <a href="proto/leetcode-downloader">leetcode-downloader</a>: Leetcode downloader

- <a href="py">py</a>: Python

  - <a href="py/autoscroll">autoscroll</a>: Autoscroll

    - <a href="py/autoscroll/autoscroll">autoscroll</a>: Src

      - <a href="py/autoscroll/autoscroll/_internal">\_internal</a>: Autoscroll

  - <a href="py/bazel-python-shell">bazel-python-shell</a>: Bazel python shell

  - <a href="py/install-file">install-file</a>: Python scripts

  - <a href="py/replace-section">replace-section</a>: Replace section

- <a href="rs">rs</a>: Rust

  - <a href="rs/tools">tools</a>: Tools

- <a href="sh">sh</a>: Shell

  - <a href="sh/git-hooks">git-hooks</a>: Hooks

  - <a href="sh/scripts">scripts</a>: Scripts

- <a href="sql">sql</a>: Sql

- <a href="tools">tools</a>: Tools' /home/simeonwarren/Git/git.alwaldend.com/src/README.md

<!-- README_CONTENTS END -->

</details>

## Setup

- Install packages: `sudo dnf install clang clang-format java-latest-openjdk-devel rust mesa-libGL-devel go Xvfb`
- Install bazel: https://bazel.build/install/bazelisk
- Install nvm: https://github.com/nvm-sh/nvm
- Install node: `nvm install node`
- Install commandline tools to `${ANDROID_HOME}/cmdline-tools/latest`: https://developer.android.com/studio#command-tools
- Install android tools: `sdkmanager "platforms;android-36" "build-tools;36.0.0"`
- Install qt: `bazel run //starlark/bazel/qt:install`
- Install git hooks: `bazel run //shell/git-hooks:install`

## Development

- Build `compile_commands.json`: `bazel run :refresh_compile_commands`
- Run builds: `bazel build //...`
- Run tests: `bazel test //...`

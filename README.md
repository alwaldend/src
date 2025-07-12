---
title: Source code
cascade:
  - type: docs
    _target:
      path: "/**"
---

Source code

## Licence

AGPLv3, see [LICENSE](./LICENSE.txt)

## Contents

<!-- README_CONTENTS START -->

- <a href=".github">.github</a>: Github config

  - <a href=".github/workflows">workflows</a>: Github workflows

- <a href="bzl">bzl</a>: Bazel projects

  - <a href="bzl/aspects">aspects</a>: Bazel aspects

  - <a href="bzl/build-files">build-files</a>: BUILD.bazel files

  - <a href="bzl/configs">configs</a>: Bazel configs

  - <a href="bzl/extensions">extensions</a>: Module extensions

  - <a href="bzl/macros">macros</a>: Bazel macros

  - <a href="bzl/providers">providers</a>: Bazel providers

  - <a href="bzl/registry">registry</a>: Bazel registry

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

  - <a href="bzl/rules">rules</a>: Bazel rules

  - <a href="bzl/toolchain-types">toolchain-types</a>: Toolchain types

  - <a href="bzl/toolchains">toolchains</a>: Bazel rules

  - <a href="bzl/vars">vars</a>: Static bazel variables

- <a href="c">c</a>: C projects

  - <a href="c/misc">misc</a>: C utils

  - <a href="c/openssl">openssl</a>: Openssl build

  - <a href="c/sri">sri</a>: Cli app to calculate subresource integrity (SRI)

- <a href="cfg">cfg</a>: Configs

  - <a href="cfg/dotfiles">dotfiles</a>: Dotfile configs

- <a href="cpp">cpp</a>: C&#43;&#43; projects

  - <a href="cpp/infinitime">infinitime</a>: Fork of InfiniTimeOrg/InfiniTime

  - <a href="cpp/leetcode-submissions">leetcode-submissions</a>: Leetcode submissions (cpp)

  - <a href="cpp/useless-qt-gui">useless-qt-gui</a>: Useless qt GUI

- <a href="data">data</a>: Data

  - <a href="data/misc">misc</a>: Miscellaneous data

- <a href="drawio">drawio</a>: Drawio diagrams

  - <a href="drawio/diagrams">diagrams</a>: Drawio diagrams

- <a href="go">go</a>: Golang projects

  - <a href="go/bazel-shell-worker">bazel-shell-worker</a>: Bazel worker that runs shell commands

  - <a href="go/file-installer">file-installer</a>: Tool to install files

    - <a href="go/file-installer/cmd">cmd</a>: Cmd code for file-installer

  - <a href="go/leetcode-downloader">leetcode-downloader</a>: CLI app to download leetcode submissions

    - <a href="go/leetcode-downloader/model">model</a>: Models for leetcode-downloader

  - <a href="go/leetcode-submissions">leetcode-submissions</a>: Leetcode submissions (go)

  - <a href="go/readme-tree">readme-tree</a>: Tool to parse README.md files

  - <a href="go/template-files">template-files</a>: Cli tool to template files

  - <a href="go/utils">utils</a>: Random golang tools

- <a href="hugo">hugo</a>: Hugo projects

  - <a href="hugo/sites">sites</a>: Hugo sites

    - <a href="hugo/sites/docs">docs</a>: Miscellaneous knowledge

  - <a href="hugo/themes">themes</a>: Hugo themes

- <a href="img">img</a>: Images

  - <a href="img/useless-qt-gui">useless-qt-gui</a>: Assets for useless-qt-gui

- <a href="java">java</a>: Java projects

  - <a href="java/leetcode-submissions">leetcode-submissions</a>: Leetcode submissions (java)

- <a href="js">js</a>: Javascript projects

  - <a href="js/leetcode-downloader">leetcode-downloader</a>: Tampermonkey script to download leetcode submissions

- <a href="kt">kt</a>: Kotlin projects

- <a href="lua">lua</a>: Lua projects

  - <a href="lua/nvim-config">nvim-config</a>: Neovim config

  - <a href="lua/nvim-lib">nvim-lib</a>: Lua library for neovim

- <a href="md">md</a>: Markdown projects

  - <a href="md/misc">misc</a>: Miscellaneous knowledge

- <a href="patch">patch</a>: Patches

  - <a href="patch/infinitime">infinitime</a>: Git patches for InfiniTimeOrg/InfiniTime

  - <a href="patch/rules-haskell">rules-haskell</a>: Fixes for rules-haskell

- <a href="pl">pl</a>: Perl projects

- <a href="proto">proto</a>: Protobuf projects

  - <a href="proto/bazel-worker">bazel-worker</a>: Bazel worker protocol

  - <a href="proto/leetcode-downloader">leetcode-downloader</a>: Models for leetcode-downloader

- <a href="py">py</a>: Python projects

  - <a href="py/bazel-python-shell">bazel-python-shell</a>: Python shell allowing you to run shell commands in python environment

  - <a href="py/install-file">install-file</a>: Python scripts

  - <a href="py/leetcode-submissions">leetcode-submissions</a>: Leetcode submissions (python)

  - <a href="py/replace-section">replace-section</a>: Replace sections of files

- <a href="rs">rs</a>: Rust projects

  - <a href="rs/tools">tools</a>: Rust tools

- <a href="sh">sh</a>: Shell projects

  - <a href="sh/git-hooks">git-hooks</a>: Git hooks

  - <a href="sh/scripts">scripts</a>: Shell scripts

- <a href="tools">tools</a>: Tools

- <a href="vial">vial</a>: Vial configs

  - <a href="vial/keyboards">keyboards</a>: Keyboard configs

<!-- README_CONTENTS END -->

<!-- CLOC START -->

| Language         |    files |    blank |  comment |     code |
| :--------------- | -------: | -------: | -------: | -------: |
| JSON             |      966 |        2 |        0 |    25853 |
| Python           |      581 |     2340 |     1222 |     9588 |
| TOML             |        8 |      391 |        8 |     7259 |
| Go               |      231 |      718 |      617 |     6539 |
| Starlark         |      118 |      543 |      699 |     3903 |
| C++              |       87 |      138 |      199 |     2219 |
| Text             |        7 |      128 |        0 |     1788 |
| Lua              |        8 |       60 |       90 |      987 |
| Markdown         |       94 |      214 |        6 |      862 |
| YAML             |        5 |       80 |        2 |      606 |
| Bourne Shell     |       19 |       85 |      103 |      518 |
| C                |        7 |       52 |       91 |      407 |
| diff             |        4 |        7 |       51 |      319 |
| Java             |       11 |       37 |       15 |      250 |
| C/C++ Header     |        7 |       14 |        9 |      111 |
| JavaScript       |        1 |        7 |       10 |       78 |
| Protocol Buffers |        2 |       21 |       64 |       74 |
| SVG              |        3 |        3 |        4 |       55 |
| XML (Qt/GTK)     |        1 |        0 |        0 |       31 |
| ProGuard         |        1 |        9 |        8 |       25 |
| INI              |        1 |        2 |        0 |        9 |
| HCL              |        1 |        1 |        5 |        8 |
| Snakemake        |        1 |        2 |        6 |        4 |
| HTML             |        2 |        0 |        0 |        3 |
| Rust             |        1 |        0 |        0 |        3 |
| --------         | -------- | -------- | -------- | -------- |
| SUM:             |     2167 |     4854 |     3209 |    61499 |

<!-- CLOC END -->

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

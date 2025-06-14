# src

Source code

## Licence

AGPLv3, see [LICENSE](./LICENSE.txt)

## Contents

<details>
  <summary>Directory structure</summary>

<!-- README_CONTENTS START -->

  - [bzl](/bzl): Bazel projects

    - [aspects](/bzl/aspects): Bazel aspects

    - [build-files](/bzl/build-files): BUILD.bazel files

    - [configs](/bzl/configs): Bazel configs

    - [extensions](/bzl/extensions): Module extensions

    - [macros](/bzl/macros): Bazel macros

    - [providers](/bzl/providers): Bazel providers

    - [registry](/bzl/registry): Bazel registry

      - [modules](/bzl/registry/modules): Modules

        - [com-github-aldanial-cloc](/bzl/registry/modules/com-github-aldanial-cloc): Cloc

        - [com-github-georgewfraser-java-language-server](/bzl/registry/modules/com-github-georgewfraser-java-language-server): Java languange server

        - [com-nordicsemi-developer-nrfsdk](/bzl/registry/modules/com-nordicsemi-developer-nrfsdk): Nrfsdk

        - [hedron_compile_commands](/bzl/registry/modules/hedron_compile_commands): Bazel-compile-commands-extractor

        - [org-openssl-openssl](/bzl/registry/modules/org-openssl-openssl): Openssl

        - [rules_haskell](/bzl/registry/modules/rules_haskell): Patched version of rules_haskell

        - [us-nasm-nasm](/bzl/registry/modules/us-nasm-nasm): Netwide Assembler (NASM)

    - [rules](/bzl/rules): Bazel rules

    - [toolchains](/bzl/toolchains): Qt wrapper for bazel

    - [vars](/bzl/vars): Static bazel variables

  - [c](/c): C projects

    - [misc](/c/misc): C utils

    - [openssl](/c/openssl): Openssl build

    - [sri](/c/sri): Cli app to calculate subresource integrity (SRI)

  - [cfg](/cfg): Configs for different tools

    - [dotfiles](/cfg/dotfiles): Dotfile configs

  - [cpp](/cpp): C&#43;&#43; projects

    - [infinitime](/cpp/infinitime): Fork of InfiniTimeOrg/InfiniTime

    - [leetcode-submissions](/cpp/leetcode-submissions): Leetcode submissions

    - [useless-qt-gui](/cpp/useless-qt-gui): Useless qt GUI

  - [data](/data): Data

    - [misc](/data/misc): Miscellaneous data

  - [drawio](/drawio): Drawio diagrams

    - [diagrams](/drawio/diagrams): Drawio diagrams

  - [go](/go): Golang projects

    - [bazel-shell-worker](/go/bazel-shell-worker): Bazel worker that runs shell commands

    - [file-installer](/go/file-installer): Tool to install files

    - [leetcode-downloader](/go/leetcode-downloader): CLI app to download leetcode submissions

    - [readme-tree](/go/readme-tree): Tool to parse README.md files

    - [template-files](/go/template-files): Cli tool to template files

    - [utils](/go/utils): Random golang tools

  - [hugo](/hugo): Hugo projects

    - [misc](/hugo/misc): Miscellaneous knowledge

    - [projects](/hugo/projects): Project documentation

  - [img](/img): Images

    - [useless-qt-gui](/img/useless-qt-gui): Assets for useless-qt-gui

  - [java](/java): Java projects

    - [leetcode-submissions](/java/leetcode-submissions): Leetcode submissions

  - [js](/js): Javascript projects

    - [leetcode-downloader](/js/leetcode-downloader): Tampermonkey script to download leetcode submissions

  - [kt](/kt): Kotlin projects

  - [lua](/lua): Lua projects

    - [nvim-config](/lua/nvim-config): Neovim config

    - [nvim-lib](/lua/nvim-lib): Lua library for neovim

  - [md](/md): Markdown projects

    - [docs](/md/docs): Project documentation

    - [misc](/md/misc): Miscellaneous knowledge

  - [patch](/patch): Patches

    - [infinitime](/patch/infinitime): Git patches for InfiniTimeOrg/InfiniTime

    - [rules-haskell](/patch/rules-haskell): Fixes for rules-haskell

  - [pl](/pl): Perl projects

  - [proto](/proto): Protobuf projects

    - [bazel-worker](/proto/bazel-worker): Bazel worker protocol

    - [leetcode-downloader](/proto/leetcode-downloader): Models for leetcode-downloader

  - [py](/py): Python projects

    - [bazel-python-shell](/py/bazel-python-shell): Python shell allowing you to run shell commands in python environment

    - [install-file](/py/install-file): Python scripts

    - [leetcode-submissions](/py/leetcode-submissions): Leetcode submissions

    - [replace-section](/py/replace-section): Replace sections of files

  - [rs](/rs): Rust projects

    - [tools](/rs/tools): Rust tools

  - [sh](/sh): Shell projects

    - [git-hooks](/sh/git-hooks): Git hooks

    - [scripts](/sh/scripts): Shell scripts

  - [tools](/tools): Tools

  - [vial](/vial): Vial configs

    - [keyboards](/vial/keyboards): Keyboard configs
<!-- README_CONTENTS END -->

</details>

<details>
  <summary>Lines of code</summary>
<!-- CLOC START -->

Language|files|blank|comment|code
:-------|-------:|-------:|-------:|-------:
JSON|960|2|0|25803
Python|571|2092|1159|8513
Go|231|718|617|6521
Starlark|96|425|594|2839
C++|87|138|199|2219
Text|5|117|0|1766
Markdown|109|554|36|1726
Lua|8|60|90|987
TOML|6|54|0|545
Bourne Shell|19|85|103|518
YAML|3|77|0|504
C|7|52|91|407
diff|4|7|51|319
Java|11|37|15|250
C/C++ Header|7|14|9|111
JavaScript|1|7|10|78
Protocol Buffers|2|21|64|74
XML (Qt/GTK)|1|0|0|31
ProGuard|1|9|8|25
HCL|1|1|5|8
INI|1|2|0|7
Snakemake|1|2|6|4
Rust|1|0|0|3
--------|--------|--------|--------|--------
SUM:|2133|4474|3057|53258

<!-- CLOC END -->

</details>


## Setup

- Install packages: `sudo dnf install clang clang-format java-latest-openjdk-devel rust mesa-libGL-devel go`
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

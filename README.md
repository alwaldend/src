# src

Source code

## Licence

AGPLv3, see [LICENSE](./LICENSE.txt)

## Contents

<details>

<!-- README_CONTENTS START -->

  - [assets](/assets): Assets

    - [useless-qt-gui](/assets/useless-qt-gui): Assets for useless-qt-gui

  - [c](/c): C projects

    - [openssl](/c/openssl): Openssl build

    - [subresource-integrity-calculator](/c/subresource-integrity-calculator): Cli app to calculate subresource integrity (SRI)

  - [configs](/configs): Configs for different tools

    - [dotfiles](/configs/dotfiles): Dotfile configs

  - [cpp](/cpp): C&#43;&#43; projects

    - [infinitime](/cpp/infinitime): Fork of InfiniTimeOrg/InfiniTime

    - [leetcode-submissions](/cpp/leetcode-submissions): Leetcode submissions

    - [useless-qt-gui](/cpp/useless-qt-gui): Useless qt GUI

  - [drawio](/drawio): Drawio diagrams

    - [diagrams](/drawio/diagrams): Drawio diagrams

  - [go](/go): Golang projects

    - [bazel-shell-worker](/go/bazel-shell-worker): Bazel worker that runs shell commands

    - [file-installer](/go/file-installer): Tool to install files

    - [leetcode-downloader](/go/leetcode-downloader): CLI app to download leetcode submissions

    - [readme-tree](/go/readme-tree): Tool to parse README.md files

    - [utils](/go/utils): Random golang tools

  - [hugo](/hugo): Hugo projects

    - [knowledge](/hugo/knowledge): Wiki

    - [projects](/hugo/projects): Project documentation

  - [java](/java): Java projects

    - [leetcode-submissions](/java/leetcode-submissions): Leetcode submissions

  - [javascript](/javascript): Javascript projects

    - [leetcode-downloader](/javascript/leetcode-downloader): Tampermonkey script to download leetcode submissions

  - [lua](/lua): Lua projects

    - [nvim-config](/lua/nvim-config): Neovim config

    - [nvim-lib](/lua/nvim-lib): Lua library for neovim

  - [patches](/patches): Patches

    - [infinitime](/patches/infinitime): Git patches for InfiniTimeOrg/InfiniTime

  - [perl](/perl): Perl projects

  - [proto](/proto): Protobuf projects

    - [bazel-worker](/proto/bazel-worker): Bazel worker protocol

    - [leetcode-downloader](/proto/leetcode-downloader): Models for leetcode-downloader

  - [python](/python): Python projects

    - [bazel-python-shell](/python/bazel-python-shell): Python shell allowing you to run shell commands in python environment

    - [install-file](/python/install-file): Python scripts

    - [leetcode-submissions](/python/leetcode-submissions): Leetcode submissions

    - [replace-section](/python/replace-section): Replace sections of files

  - [rust](/rust): Rust projects

    - [tools](/rust/tools): Rust tools

  - [shell](/shell): Shell projects

    - [git-hooks](/shell/git-hooks): Git hooks

    - [scripts](/shell/scripts): Shell scripts

  - [starlark](/starlark): Starlark projects

    - [bazel](/starlark/bazel): Bazel projects

      - [aspects](/starlark/bazel/aspects): Bazel aspects

      - [build-files](/starlark/bazel/build-files): BUILD.bazel files

      - [configs](/starlark/bazel/configs): Bazel configs

      - [extensions](/starlark/bazel/extensions): Module extensions

      - [macros](/starlark/bazel/macros): Bazel macros

      - [providers](/starlark/bazel/providers): Bazel providers

      - [qt](/starlark/bazel/qt): Qt wrapper for bazel

      - [rules](/starlark/bazel/rules): Bazel rules

  - [vial](/vial): Vial configs

    - [keyboards](/vial/keyboards): Keyboard configs
<!-- README_CONTENTS END -->

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

- Run builds: `bazel build //...`
- Run tests: `bazel test //...`

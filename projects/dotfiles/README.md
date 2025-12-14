---
title: Dotfiles
description: Dotfiles
statuses:
  - maintenance
languages:
  - sh
  - lua
tags:
  - vial
  - dotfiles
  - make
  - udev
---

## Links

- Source code: https://github.com/alwaldend/src/tree/master/projects/dotfiles

## Features

- Archive with dotfiles that can be installed using make

## Bazel usage

```sh
bazel run //projects/dotfiles:help # Show help
bazel run //projects/dotfiles:diff # Show diff
bazel run //projects/dotfiles:install # Install files
```

## Archive usage

```sh
wget https://alwaldend.com/docs/projects/dotfiles/releases/head/dotfiles.tar
tar -xf dotfiles.tar
cd dotfiles
make help # Show help
make diff # Check diff between system files and those from the archive
make diff/nvim # Check diff for all nvim files
make diff/nvim/.config/nvim/lazy-lock.json # Diff specific file
make install # Install all files from the archive
make install/nvim # Install nvim files
make install/nvim/.config/nvim/lazy-lock.json # Install a specific file
```

## Help

```
help: Show help
install: Install files from the archive
diff: Show diff between archive files and system files
install/bin: Install bin files
diff/bin: Diff bin files
install/home: Install home files
diff/home: Diff home files
install/nvim: Install nvim files
diff/nvim: Diff nvim files
make: Nothing to be done for 'help'.
```

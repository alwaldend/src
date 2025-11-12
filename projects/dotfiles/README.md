---
title: Dotfiles
description: Dotfiles
statuses:
  - maintenance
languages:
  - sh
tags:
  - dotfiles
  - make
---

## Bazel usage

```sh
bazel run //projects/dotfiles:diff # Show diff
bazel run //projects/dotfiles:install # Install files
```

## Archive usage

```sh
tar -xf dotfiles.tar
cd dotfiles
make diff # Check diff between system files and those from the archive
make diff/nvim # Check diff for all nvim files
make diff/nvim/.config/nvim/lazy-lock.json # Diff specific file
make install # Install all files from the archive
make install/nvim # Install nvim files
make install/nvim/.config/nvim/lazy-lock.json # Install a specific file
```

## Help

```
install: Install files from the archive
diff: Show diff between archive files and system files
help: Show help
install/nvim: Install nvim files
diff/nvim: Diff nvim files
```

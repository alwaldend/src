---
title: Sri
description: Cli app for subresource integrity calculation
statuses:
  - finished
languages:
  - c
tags:
  - cli
---

## Downloads

See [./releases/head](./releases/head)

## Links

- https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity

## Usage

```
Usage: sri [OPTION...]
Generate sri of a file

Example:
    bazel run //projects/sri -- --digest sha256 --file ${PWD}/README.md

  -d, --digest=String        Digest type (sha256, for example)
  -f, --file=Path            Path to the file to parse
  -?, --help                 Give this help list
      --usage                Give a short usage message
```

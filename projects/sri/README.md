---
title: Sri
description: Subresource Integrity calculator
statuses:
  - finished
languages:
  - c
tags:
  - cli
---

## Links

- Source code: https://github.com/alwaldend/src/tree/master/projects/sri
- Docs: https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity

## Features

- Cli app that calculates SRI using [OpenSSL](https://www.openssl.org/)
- C, OpenSSL

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

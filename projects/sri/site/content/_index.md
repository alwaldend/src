---
title: Sri
linkTitle: Sri
description: Command-line Subresource Integrity calculator using OpenSSL
layout: landing
statuses:
  - finished
languages:
  - c
tags:
  - cli
---

Sri calculates [Subresource Integrity](https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity) hashes for files using OpenSSL. Point it at a file and a digest algorithm and it prints the hash to embed in a `integrity` attribute.

## Features

- Command-line hashing for Subresource Integrity attributes
- Implemented in C on top of OpenSSL
- Digest algorithm selected per invocation

## Usage

```text
Usage: sri [OPTION...]
Generate sri of a file

  -d, --digest=String        Digest type (sha256, for example)
  -f, --file=Path            Path to the file to parse
  -?, --help                 Give this help list
      --usage                Give a short usage message
```

[Documentation](/docs/projects/sri/)

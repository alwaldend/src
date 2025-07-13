---
title: SRI calculator
description: Cli app to calculate subresource integrity
---

## Links

- https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity

## Help

```
Usage: sri [OPTION...]
Generate sri of a file

Example:
    bazel run //c/sri:bin -- --digest sha256 --file ${PWD}/README.md

  -d, --digest=String        Digest type (sha256, for example)
  -f, --file=Path            Path to the file to parse
  -?, --help                 Give this help list
      --usage                Give a short usage message
```

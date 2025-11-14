---
title: Leetcode downloader
description: CLI app to download leetcode submissions
statuses:
  - in_progress
languages:
  - go
  - js
  - bzl
tags:
  - tampermonkey
  - bzl_rules
  - proto
---

{{< alwaldend/alert >}}
Download doesn't work anymore because of Cloudflare
{{< /alwaldend/alert >}}

## Downloads

See [./releases/head](./releases/head/)

## Usage

### Generate submission from a submission file

```sh
bazel run go/leetcode_downloader -- \
    --submissions-file "${PWD}/out/submissions.json" \
    --root-dir "${PWD}" \
    generate
```

## Help

```
Usage of flags:
  -base_url string
    	 (default "https://leetcode.com")
  -cookie string

  -limit uint
    	 (default 20)
  -offset uint

  -root-dir string
    	 (default "${PWD}")
  -submissions-file string
```

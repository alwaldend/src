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
---

{{< alwaldend/alert >}}
Download doesn't work anymore because of Cloudflare
{{< /alwaldend/alert >}}

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

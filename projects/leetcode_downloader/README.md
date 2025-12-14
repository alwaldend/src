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

## Links

- Source code: https://github.com/alwaldend/src/tree/master/projects/leetcode_downloader

## Features

- CLI and bzl code to generate submission docs
- [Tampermonkey](https://www.tampermonkey.net/) script to download submissions
  (Doesn't work because the API was removed)
- CLI to download submissions (Doesn't work because of bot protection)

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

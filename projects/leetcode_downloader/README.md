---
title: Leetcode downloader
description: LeetCode submission export and documentation tools
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

This project exports LeetCode submissions and generates documentation from
submission files using a CLI and Bazel rules. A Tampermonkey script provides
a browser-based download path; direct CLI downloads are currently blocked by
LeetCode bot protection.

## Links

- Source code: https://github.com/alwaldend/src/tree/master/projects/leetcode_downloader

## Features

- CLI and bzl code to generate submission docs
- [Tampermonkey](https://www.tampermonkey.net/) script to download submissions
- CLI to download submissions (Doesn't work because of bot protection)

## Usage

### Generate submission from a submission file

```sh
bazel run //projects/leetcode_downloader -- \
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

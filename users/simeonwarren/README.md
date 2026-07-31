---
title: simeonwarren
description: simeonwarren
---

## Generate a client cert

```sh
bazel run //users/simeonwarren:gen_client_cert -- --host phone1 --ttl "180d" --output_dir "${PWD}/certs"
```

## Run opencode

```sh
bazel run //users/simeonwarren:opencode -- "${PWD}"
```

## Apply terraform

```sh
bazel run //users/simeonwarren/tf:tf.apply
```

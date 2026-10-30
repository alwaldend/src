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

## Prepare opencode

```sh
ln -s "${PWD}/users/simeonwarren/opencode.json" ~/.config/opencode/opencode.json
ln -s "${PWD}/users/simeonwarren/tui.json" ~/.config/opencode/tui.json
```

## X profile cleanup

[`x_profile_cleaner.lua`](x_profile_cleaner.lua) selects posts, replies, and
reposts older than 30 days, while preserving the six named threads and their
connected own quote/reply branches. Open it for review with:

```sh
bazel_agent bazel run //projects/x_profile_cleaner -- \
  serve --config "${PWD}/users/simeonwarren/x_profile_cleaner.lua"
```

The policy leaves discovery unbounded to retain available thread context.
A seed absent from the available timeline causes planning to fail; see the
[cleaner's discovery boundary](../../projects/x_profile_cleaner/README.md#discovery-boundary).

[`x_profile_cleaner_single_post.lua`](x_profile_cleaner_single_post.lua) is a
separate one-post demonstration policy, fetching only the target day. Pass its
path to the same `serve`
command to review that target without loading the 30-day cleanup rules.

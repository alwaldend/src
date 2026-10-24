---
title: Browser filesystem provider
description: Shared caller-owned storage using Emscripten classic MEMFS
---

`createEmscriptenFilesystem(options?)` creates one initialized Node-compatible
provider over Emscripten classic MEMFS. The caller owns its lifetime and selects
the root shared by Go and Rust:

```js
import { createEmscriptenFilesystem } from "./emscripten.js";

const fs = await createEmscriptenFilesystem();
const root = "/workspace";
fs.mkdirSync(root);
```

Use this provider with the [browser backend API](../../../README.md#browser-command-backend).
`installGoFilesystem(fs)` installs Go's callback filesystem facade; the backend
factory captures the same `fs` and `root` for direct synchronous Rust filesystem
calls. Each provider has independent files and descriptors. The adapter installs
no globals, implements no backing store, and keeps the generated module alive
with the returned object. Each call copies options and lifecycle arrays before
passing them to its asynchronous factory, so concurrent callers can reuse options.

`//projects/anki_grpc/pkg/ankigrpc/browser:emscripten_files` supplies both
`emscripten.js` and the generated `emscripten_memfs.js` for consumer packaging.
The adapter supports string paths, byte-buffer I/O, Node callback and Sync APIs,
numeric/string open flags, directory options, and numeric/bigint Stats. Standard
descriptors 0–2 remain reserved. Unsupported hard links return `EMLINK`;
ownership changes return `ENOTSUP`. Metadata reflects upstream MEMFS, with
birthtime falling back to ctime. Synchronization completes only in memory;
durable snapshots remain caller-owned. The Rust artifact remains WASI and uses
`wasi.js`; this provider does not retarget or relink Rust with Emscripten.

## Generated artifact

`//projects/anki_grpc/pkg/ankigrpc/browser:emscripten_memfs_files` provides
`emscripten_memfs.js` as a DefaultInfo file and runfile. It is generated from
the pinned official Emscripten SDK and exports a default asynchronous factory:

```js
import createEmscriptenMemfs from "./emscripten_memfs.js";

const module = await createEmscriptenMemfs();
module.FS.mkdir("/workspace");
module.FS.writeFile("/workspace/example.txt", "Hello");
```

Each call creates independent classic MEMFS, descriptors and working-directory
state. `.FS` and `.ERRNO_CODES` are explicitly exported. Standard descriptors
0, 1 and 2 remain initialized. Emscripten owns filesystem semantics; this raw
artifact does not implement a Node-compatible projection. Go js/wasm and the
Rust WASI reactor remain independent consumers through separately owned adapters.
IndexedDB persistence remains explicit application snapshots, with no IDBFS.

The selected `WASM=0` output is **one JavaScript file**, with no `.wasm` sidecar,
no runtime import and no WebAssembly engine requirement. Emscripten emits its
ordinary modularized runtime and uses Binaryen's wasm2js during the build;
no source fragments are extracted. Its generated runtime allocates a 64 KiB
ArrayBuffer per instance, in addition to JavaScript objects and file contents.
The generated artifact E2E reports exact byte count, gzip size and SHA-256.

Upstream deprecates `WASM=0` and may remove it. The pinned 6.0.13 output is
tested; upgrades must rerun the artifact test and reassess output mode. The
build rule can declare a companion Wasm file if a future migration requires it.
Any such change must also update consumer packaging and the no-WebAssembly test.

The raw filesystem supports streams, fstat, ftruncate, fchmod, fchown, utime,
symlinks and the synchronous operations used by the provider projection.
Classic MEMFS has no hard-link operation; ownership-changing operations are
compatibility no-ops. `readlink` returns the stored link text in this release.

Run the raw artifact and Node-compatible provider E2Es with pinned Node:

```sh
bazel_agent bazel test //projects/anki_grpc/test/emscripten_memfs:emscripten_memfs_test \
  //projects/anki_grpc/test/browser_backend:emscripten_filesystem_test
```

The test runs two actual factory instances in an isolated ES-module realm inside
a bounded Node worker. It selects Emscripten's worker environment and asserts
the absence of Node, DOM, fetch and WebAssembly globals. It verifies binary and
Unicode operations, stream offsets, rename with an open descriptor, truncation,
instance isolation, expected errors and initialized stdio. Its undeclared output
`emscripten_memfs_report.json` records the exact artifact and exercised behaviors.
This is artifact-level evidence; application browser and native collection
round trips belong to the composed backend/provider/application tests.
The independent reusable-backend Chromium target is
`//projects/anki_grpc/test/browser_backend:browser_storage_test`; it uses an
explicit application-owned Emscripten fixture with production Go/Rust hosting.

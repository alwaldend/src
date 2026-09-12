## Why

`docs_filegroup` flattens every source to its basename, because it relies on
`pkg_files` with the default `strip_prefix`. A package whose glob reaches into
subdirectories therefore fails analysis as soon as two sources share a name:
`pkg_files` refuses to map both `spec.md` files to one destination.

`projects/alwaldend.com/openspec` is the first package to hit this. Archiving a
change there added a second `specs/<capability>/spec.md`, and
`//projects/alwaldend.com/openspec:docs` stopped analyzing while the baseline
still built.

## What Changes

- Add an opt-in `preserve_paths` parameter to `docs_filegroup` that keeps each
  source's path relative to its package instead of flattening it.
- Enable it for `projects/alwaldend.com/openspec`, the one package whose
  sources currently collide.

Flattening stays the default so existing consumers keep their published paths.

## Impact

- `projects/rules_docs/docs/defs.bzl` gains the parameter; default behavior is
  unchanged for every other consumer.
- `projects/alwaldend.com/openspec/BUILD.bazel` opts in.

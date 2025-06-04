#!/usr/bin/env sh

# See https://github.com/bazelbuild/rules_go/wiki/Editor-setup#3-editor-setup

if which bazel >/dev/null 2>/dev/null; then
    exec bazel run @rules_go//go/tools/gopackagesdriver -- "${@}"
fi

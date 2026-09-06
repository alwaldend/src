"""A separate pip extension identity for lazily fetched firmware lockfiles."""

load("@rules_python//python/extensions:pip.bzl", _pip = "pip")

pip = _pip

## Why

The reviewer requested removal of redundant \_module import alias suffixes.

## What Changes

Use natural module names in all affected imports and qualified references.
Teach repo-python to avoid redundant suffixes and aliases.

## Capabilities

Import naming correction; no runtime behavioral delta.

## Impact

Autoscroll import names and the Python skill and eval cases.

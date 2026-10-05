## Why

PR feedback requests explicit package exports instead of cross-package imports
of implementation files.

## What Changes

Expose the recommender model contracts, services, router factories and application
factory through package `__init__.py` files. Update consumers and Uvicorn's
factory path, and document this strong preference in repo-python.

## Capabilities

Package API organization correction without changes to HTTP behavior.

## Impact

Recommender Python source, generated BUILD declarations, and the Python skill.

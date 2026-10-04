## Decision

Use from parent.package import module for long module paths, retaining module
aliases and namespaced member references. Validate imported names as modules
using the tracked source inventory and available pinned Python modules.
Do not import packages merely to run validation.

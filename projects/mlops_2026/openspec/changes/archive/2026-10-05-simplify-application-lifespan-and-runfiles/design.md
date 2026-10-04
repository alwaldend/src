## Decision

FastAPI receives the application owner as its lifespan callable. `__call__`
returns that owner; `__aenter__` loads its model and `__aexit__` releases it.
Each factory invocation retains independent state and imports have no services.

One runfile helper owns resolver creation, repository mapping and existence
checks. CSV and model wrappers keep their existing messages. Preserve directory
and manifest runfiles, explicitly mapped external repository names and operation
from arbitrary working directories. Existing HTTP E2E checks cover startup
failures, recommendations, shutdown and per-instance routing.

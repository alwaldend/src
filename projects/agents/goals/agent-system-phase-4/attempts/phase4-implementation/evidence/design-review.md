# Phase 4 design review

The join between goal resume state, candidate-bound validation sets, typed
rewrite authorization, durable review references, and guarded release-ref
publication keeps each authority separate: goal, delivery, versioning,
release, and review consume typed references and never merge into one
orchestrator. Unsupported remote atomicity is an explicit refusal.

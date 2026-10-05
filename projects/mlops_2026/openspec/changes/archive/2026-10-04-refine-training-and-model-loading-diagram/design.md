## Decision

Trace the pinned upstream fit and serialization implementation. FastAPI calls
upstream load_model at startup, which calls joblib.load on the packaged file.
Use read/write edge labels to distinguish loading from artifact production.
Requests use the per-app loaded model without reopening the file.

## Decision

Use instance-owned application lifecycle and recommendation route classes, move HealthResponse into shared models, and expose service through the recommender package API.

Keep application and service objects per factory invocation. Preserve API paths, query validation, OpenAPI schemas, startup loading and shutdown cleanup. Existing HTTP and consumer E2E checks validate behavior. Review the aggregate PR for independent defects.

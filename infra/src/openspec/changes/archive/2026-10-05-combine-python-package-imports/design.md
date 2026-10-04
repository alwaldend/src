## Decision

Require combined imports from the same package, configure the existing import sorter to combine aliased imports, and document that different package paths require a common exported namespace before they can be combined.

Keep application and service objects per factory invocation. Preserve API paths, query validation, OpenAPI schemas, startup loading and shutdown cleanup. Existing HTTP and consumer E2E checks validate behavior. Review the aggregate PR for independent defects.

Classify the upstream `recommender` module as third-party explicitly in Ruff
and isort. A full checkout contains a same-named project package while Bazel
lint sandboxes expose only target sources; explicit classification keeps both
import sorting contexts consistent.

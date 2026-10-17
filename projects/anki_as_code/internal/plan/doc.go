// Package plan builds executable collection plans from supplied configuration,
// native snapshots and stored bindings. It owns reference resolution, resource
// comparison, ownership checks and operation preparation without I/O. Native
// values and prepared protobufs use private indexes. The protobuf plan is the
// executor boundary. Anki owns native validation, normalization and generated-card
// migrations during execution.
package plan

# health-api Specification

## Purpose

Provide a runnable HTTP foundation for the MLOps course with a health check,
an inspectable API description, and independent application instances.

## Requirements

### Requirement: Provide an operational health endpoint

The service SHALL respond to GET /health with HTTP 200 and the JSON object
{"status": "ok"}. The response SHALL indicate that the HTTP application is
running without invoking recommendation computation or external services.

#### Scenario: Check a running application

- **WHEN** a client sends GET /health to a running service
- **THEN** it receives HTTP 200 with application/json content and {"status": "ok"}

### Requirement: Publish the generated API description and documentation

The service SHALL provide a valid OpenAPI document at GET /openapi.json,
including the health endpoint and its successful response contract. It SHALL
provide the default documentation pages at /docs and /redoc, configured to use
that schema URL.

#### Scenario: Inspect the API contract

- **WHEN** a client requests GET /openapi.json
- **THEN** it receives HTTP 200 with a valid OpenAPI JSON document describing
  GET /health and its successful JSON response

#### Scenario: Open API documentation

- **WHEN** a client requests GET /docs or GET /redoc
- **THEN** it receives HTTP 200 with an HTML documentation page referring to
  /openapi.json

### Requirement: Preserve a versioned application API boundary

Future application endpoints SHALL use /api/v1 as their initial routing prefix.
Recommendations SHALL be exposed under that prefix; health and default API
documentation SHALL retain their agreed root paths.

#### Scenario: Request an unimplemented application endpoint

- **WHEN** a client requests GET /api/v1/predict
- **THEN** the service responds with HTTP 404 rather than exposing prediction
  or placeholder business behavior

### Requirement: Keep application instances independent

Constructing an application SHALL create fresh application and routing state.
Importing application modules SHALL NOT start a server, create shared mutable
application instances, or allocate external service resources. Changes to one
application's routing SHALL NOT affect another application's routing.

#### Scenario: Construct multiple applications

- **WHEN** two applications are constructed in one process and routing is
  subsequently added to one of them
- **THEN** the other application's routes remain unchanged

### Requirement: Provide a configurable local server entry point

The service SHALL be runnable through a repository Bazel target with configurable
host and port. It SHALL bind to 127.0.0.1:8000 by default, accept an explicitly
selected loopback address and available port, and shut down cleanly on termination.

#### Scenario: Use an explicitly selected port

- **WHEN** the server target starts with a loopback host and an available port
- **THEN** the health and OpenAPI endpoints are served on that address
- **AND** terminating the server releases its resources and listening socket

### Requirement: Serve model recommendations

GET /api/v1/recommend SHALL accept a required track_name query and an integer n
query defaulting to 5, constrained to 1 through 20. It SHALL return a
RecommenderRecommendationResponse with requested_track and recommendations of
RecommenderSong (track_name, artists, album_name). It SHALL use the upstream
load_model function to load the declared model.pkl during application lifespan,
with per-application state and no global model object.

#### Scenario: Recommend a known track

- **WHEN** a client requests a track in the trained model with n=3
- **THEN** it receives HTTP 200 and three typed song recommendations
- **AND** requested_track matches the query

#### Scenario: Reject an invalid recommendation count

- **WHEN** n is 0, 21, or nonnumeric, or track_name is absent
- **THEN** the response is HTTP 422

#### Scenario: Request an unknown track

- **WHEN** track_name does not exist in the model
- **THEN** requested_track is preserved and recommendations is empty

## ADDED Requirements

### Requirement: Required native decks retain their bindings

Planning SHALL reject deletion of a tracked native default deck with ID 1,
including after it has been renamed. Plan and apply SHALL reject the request
before native mutations or stored-binding changes and SHALL preserve input
configurations.

#### Scenario: Remove a tracked default deck declaration

- **WHEN** the configuration omits the tracked deck whose native ID is 1
- **THEN** planning and apply report that the required deck cannot be removed,
  without performing unrelated requested creations or discarding its binding

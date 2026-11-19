# Nexus security plugin

## Purpose

Provide configurable artifact checks for Sonatype Nexus 3. This baseline was
observed at repository revision `550d7e79` on 2026-09-08. Sources are the
[project README](../../../README.md),
[request handler](../../../main/java/main/RequestHandler.java),
[configuration update task](../../../main/java/capability/task/SecurityCapabilityUpdateTask.java),
and [Java build](../../../main/java/BUILD.bazel).
The handler's implemented checks cover eligible proxy-repository responses;
this specification does not claim scanning of every repository response.

## Requirements

### Requirement: Restrict scanning to eligible responses

The request handler SHALL first obtain the repository response and SHALL return
it without scanning when the security capability is inactive, the repository
is not a proxy, the payload is not Nexus `Content`, or the content has no asset.

#### Scenario: Capability is inactive

- **WHEN** a repository request reaches the handler while its capability is inactive
- **THEN** the handler returns the response from the next handler without invoking its scanners

#### Scenario: A response does not identify an asset

- **WHEN** an active capability receives proxy content without an asset attribute
- **THEN** the handler returns that response without scanning it

### Requirement: Apply enabled scanners in configured order

For an eligible response, the handler SHALL invoke enabled local scanning before
enabled remote scanning, stop at the first disallowing result, send accumulated
results to monitoring, and reject a disallowed artifact by raising an error.

#### Scenario: Local scanner disallows an artifact

- **WHEN** both scanners are enabled and the local scanner returns a disallowing result
- **THEN** the remote scanner is skipped and the handler raises an error identifying the artifact and recorded reason after invoking monitoring

#### Scenario: All invoked scanners allow an artifact

- **WHEN** every invoked scanner permits an eligible artifact
- **THEN** the handler returns the original repository response after invoking monitoring

### Requirement: Install only a successful remote configuration response

The configuration update task SHALL install a remote bundle configuration only
after a successful HTTP response with a non-null body and SHALL record the
update outcome in capability task status.

#### Scenario: Remote configuration is available

- **WHEN** the configured request returns a successful response with a configuration body
- **THEN** the task replaces the bundle configuration and records successful update status

#### Scenario: Remote configuration response is invalid

- **WHEN** the request fails or returns an unsuccessful response or null configuration
- **THEN** the task records failure status and raises the failure without installing that response

# CI Platform Specification

## Purpose

Preserve the source contract of an abandoned CI platform with Go backend
components and a Vue frontend. This baseline records checked-in behavior at
revision `550d7e79b1f5fdbc2b6017b75178471d6914082f`, observed on 2026-09-08.
It does not assert a functioning deployment or complete frontend integration.

Sources: [project README](../../../README.md),
[root command](../../../main/go/cmd/root_cmd.go),
[server command](../../../main/go/cmd/server_cmd.go),
[database command](../../../main/go/cmd/database_cmd.go),
and [frontend filesystem](../../../main/go/web/dist.go).

## Requirements

### Requirement: Explicit abandoned status

Project documentation SHALL identify the platform as abandoned. The baseline
SHALL distinguish the retained frontend source from an operational embedded
frontend, whose embed directive is commented out in the checked-in source.

#### Scenario: Assess supported operation

- **WHEN** a consumer reads this baseline to assess deployment readiness
- **THEN** the documented status SHALL remain abandoned and SHALL provide no
  guarantee of a complete working CI service.

### Requirement: Backend command dispatch

The retained CLI SHALL expose a persistent `--config` option, a `server run`
command that invokes the application's `Run` method, and a `database migrate`
command that invokes its `Migrate` method.

#### Scenario: Select a backend operation

- **WHEN** the CLI parses `server run` or `database migrate`
- **THEN** it SHALL dispatch to the corresponding application method through
  the retained command definitions.

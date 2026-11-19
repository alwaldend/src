# tool-cache Specification

## Purpose

Describe AL's cache-backed execution of Bazel-built tools declared in AL
configuration.

## Requirements

### Requirement: Resolve declared tools through AL configuration

`al tool` SHALL load tool declarations from the merged AL configuration. With
no `--config`, it SHALL use the repository-root `al.lua`; explicit `--config`
paths SHALL replace that default. After selecting a unique declared tool by
name, `al tool` SHALL run that tool with the remaining arguments. It SHALL fail
when the name is absent or ambiguous.

#### Scenario: Run a declared tool

- **WHEN** `al.lua` declares one tool named by the invocation
- **THEN** `al tool` starts the resolved tool executable and passes the remaining
  arguments unchanged.

#### Scenario: Reject missing or duplicate tools

- **WHEN** the requested name matches no tool or more than one tool
- **THEN** `al tool` does not build or execute anything and reports the failure.

### Requirement: Reuse validated executable cache entries

AL SHALL store executable cache entries under the user's XDG cache hierarchy by
default and SHALL allow `AL_TOOL_CACHE` and `--cache-root` to select an alternate
root. A source-keyed entry SHALL be runnable only when its metadata matches the
request and its selected executable is regular and executable. A valid entry
SHALL be executed without invoking Bazel.

#### Scenario: Execute a valid cache hit

- **WHEN** an entry exists for the exact request identity and passes validation
- **THEN** `al tool` executes the cached executable and does not build or install
  it.

### Requirement: Install misses through a configured Bazel method

For a cache miss, AL SHALL build the tool's configured Bazel label using Bazel,
install the built executable into the source-keyed cache entry, and execute that
executable. Installation SHALL be atomic, private to the user, and fail closed
if the selected executable is absent or not executable. Method and executable
location SHALL be supplied by the AL tool declaration, with Bazel as the default
build method.

#### Scenario: Build and install on a miss

- **WHEN** no validated cache entry exists
- **THEN** AL builds the configured Bazel label, publishes a complete cache entry
  for the exact request identity, and executes the installed executable.

#### Scenario: Refuse an unusable built output

- **WHEN** the configured method does not produce the declared executable or it
  lacks execute permission
- **THEN** AL does not publish a valid cache entry and reports failure.

### Requirement: Preserve tool execution identity

AL SHALL execute the resolved tool as a child process with inherited standard
streams and environment, preserving its process exit status. On cancellation it
SHALL stop the child and report that invocation as failed.

#### Scenario: Return child status

- **WHEN** the cached tool starts successfully and exits
- **THEN** `al tool` returns that tool's exit status as its own.

### Requirement: Remove the Bazel agent tool subcommand

`bazel_agent` SHALL accept only its Bazel execution and doctor subcommands. It
MUST NOT expose `tool run`, `tool warm`, `tool path`, or their cache and launcher
compatibility surface.

#### Scenario: Invoke the removed subcommand

- **WHEN** a caller runs `bazel_agent tool ...`
- **THEN** the runner reports usage failure without inspecting a tool cache.

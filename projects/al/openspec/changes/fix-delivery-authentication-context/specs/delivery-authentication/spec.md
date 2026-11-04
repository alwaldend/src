## Purpose

Guide delivery through GitHub credential-access failures without confusing API authentication, Git transport, and execution approval.

## ADDED Requirements

### Requirement: Actionable API authentication failure

AL SHALL retain a failed GitHub API preflight as a failure, preserve its error cause and sanitized diagnostic, and explain on HTTP 401 that Git transport selection does not repair API credential access. AL MUST NOT automatically retry, extract credentials, change login state, or escalate execution.

#### Scenario: API rejects authentication

- **WHEN** the GitHub CLI repository preflight reports HTTP 401
- **THEN** delivery stops before Git publication and explains how to compare credential access in the failing and approved execution contexts

#### Scenario: Another preflight failure

- **WHEN** repository preflight fails without HTTP 401
- **THEN** delivery preserves the failure without attributing it to inaccessible credentials

### Requirement: Recover through verified execution context

The delivery procedure SHALL distinguish API authentication from Git transport. When authentication fails in a sandbox, it SHALL compare sanitized status in the failing and explicitly approved contexts before asking for login changes. It MUST retain transport and receipt consistency, honor approval rejection, and avoid token disclosure.

#### Scenario: Stored login works outside the sandbox

- **WHEN** authentication fails in the sandbox but succeeds in an approved context
- **THEN** delivery uses that context through the execution approval mechanism with the same selected entry point and transport, without requesting a new login

#### Scenario: Authentication remains invalid

- **WHEN** the approved-context comparison also fails or cannot be authorized
- **THEN** the procedure reports the remaining authentication or approval blocker without bypassing it or claiming credential expiry from HTTP 401 alone

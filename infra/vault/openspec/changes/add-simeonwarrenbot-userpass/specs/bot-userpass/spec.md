## Purpose

Generate and store selected users' Vault passwords privately, with one explicit
read exception for the human operator of the bot account.

## ADDED Requirements

### Requirement: Selected userpass generation

An explicit selected-user list SHALL initially contain only simeonwarrenbot.
Each selected user SHALL receive a stable generated 72-character password,
stored with its username at `secrets/alwaldend.com/users/<username>/userpass`.
Unselected users SHALL receive no generated credential. Login updates SHALL
use the successfully stored password and preserve the existing alias/entity.
Secret values SHALL NOT enter source or normal command output.

#### Scenario: Reconcile the selected user

- **WHEN** the unchanged configuration is applied
- **THEN** only the bot userpass entry SHALL be managed
- **AND** its password SHALL remain unchanged and have length 72

#### Scenario: Failed storage

- **WHEN** storing a selected password fails
- **THEN** the login SHALL NOT be changed to that unstored password

### Requirement: No automatic user access

Users and the users group SHALL receive no policy granting access to userpass
entries. No explicit deny rules SHALL be added. The infrastructure writer SHALL
manage only selected entries, including metadata read/update without metadata
deletion. Existing unrelated permissions SHALL remain unchanged.

#### Scenario: Bot attempts to access its own userpass

- **WHEN** the bot reads or modifies its userpass data or metadata
- **THEN** the operation SHALL be denied by absence of an access grant

### Requirement: Exact human read exception

The existing simeonwarren identity SHALL receive read access to only the bot
userpass entry's data and metadata. That exception SHALL grant no writes,
namespace listing, or sibling reads.

#### Scenario: Read exception boundaries

- **WHEN** simeonwarren reads the bot userpass entry
- **THEN** it SHALL receive the matching username and 72-character password
- **AND** writes, bot namespace listing, and sibling reads SHALL remain denied

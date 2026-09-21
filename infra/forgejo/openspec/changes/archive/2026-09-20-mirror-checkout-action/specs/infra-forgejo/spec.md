## ADDED Requirements

### Requirement: Organization-owned CI action pull mirrors

Forgejo SHALL support catalog-declared organization-owned pull mirrors sourced
from each record's original upstream URL. Terraform SHALL own creation and
mirror synchronization settings. The checkout action mirror SHALL have its own
workflows disabled, and repository CI SHALL consume it at an immutable commit.
Mirror refreshes SHALL NOT implicitly update the workflow's action pin.

#### Scenario: Create the checkout mirror

- **WHEN** the authorized mirror deployment is applied
- **THEN** `alwaldend/com_github_actions_checkout` is a public pull mirror of `https://github.com/actions/checkout`
- **AND** the existing workflow action commit is available from the mirror
- **AND** the deployment does not recreate existing repositories or create copies on unselected forges

#### Scenario: Consume a mirrored action

- **WHEN** repository CI checks out its source revision
- **THEN** its checkout implementation comes from the Forgejo mirror at the declared immutable action commit
- **AND** the source revision remains the workflow event SHA

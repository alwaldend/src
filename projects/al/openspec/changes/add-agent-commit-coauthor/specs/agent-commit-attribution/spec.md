## ADDED Requirements

### Requirement: Configured agent commit attribution

AL delivery SHALL use the optional `al.deliveryCoAuthor` value from the
account's global Git configuration to add a `Co-authored-by` trailer.

#### Scenario: Prepare an attributed commit

- **WHEN** the account configures one valid co-author identity
- **THEN** the prepared commit contains that identity once in the final trailer block
- **AND** other co-authors are preserved and the LLM disclaimer remains the last trailer

#### Scenario: Prepare without configured attribution

- **WHEN** the setting is absent
- **THEN** delivery uses its existing commit message behavior

#### Scenario: Invalid global attribution

- **WHEN** the setting contains multiple values or a malformed identity
- **THEN** preparation refuses to create the commit

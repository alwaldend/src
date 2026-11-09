## ADDED Requirements

### Requirement: Team-owned token delivery

The existing Vast infrastructure identity SHALL publish its generated team token to the exact regular-user imagegen/vastai KV entry without using the recipient's authentication credentials.

#### Scenario: Destination access

- **WHEN** the team root manages the generated token's user KV entry
- **THEN** its existing writer identity SHALL have lifecycle access to that entry and read access to its metadata

#### Scenario: Unrelated user secret

- **WHEN** the team writer requests a sibling user secret or user Terraform backend
- **THEN** this policy SHALL grant no access to that path

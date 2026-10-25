## ADDED Requirements

### Requirement: Route the activated public root site through download hosting

After the authorized website rollout, the public apex SHALL resolve through
`download.alwaldend.com` using a 300-second CNAME declaration. The eight Pages
apex address records SHALL be absent. Local apex, mail, staging, and `www`
declarations SHALL retain their existing values and ownership.

#### Scenario: Activate the verified public website

- **WHEN** the authorized cutover follows verification of the published Yandex copy
- **THEN** public apex DNS resolves to the download host
- **AND** apex HTTPS serves that website and www redirects to the apex with trusted certificates
- **AND** unrelated records remain unchanged and a follow-up plan has no changes

## Purpose

Give public git.alwaldend.com a GitHub redirect while retaining the independent
local Forgejo endpoint and the existing authoritative DNS record ownership.

## ADDED Requirements

### Requirement: Public GitHub routing with local Forgejo preservation

The global git.alwaldend.com CNAME SHALL point to download.alwaldend.com.
The local A record SHALL continue to point to the Forgejo VM. Other DNS
records SHALL remain unchanged. The public download proxy SHALL redirect
ordinary HTTPS requests to github.com while preserving their path and query
string. Requests with go-get=1 SHALL instead return discovery metadata naming
git.alwaldend.com/alwaldend/src, Git, and https://github.com/alwaldend/src.git
for that repository root and its descendants. Unknown prefixes SHALL NOT
receive that metadata.

#### Scenario: Follow a public repository link

- **WHEN** a public client requests https://git.alwaldend.com/alwaldend/src?tab=readme
- **THEN** trusted HTTPS redirects to https://github.com/alwaldend/src?tab=readme

#### Scenario: Use local Forgejo

- **WHEN** a dc1 client resolves git.alwaldend.com
- **THEN** it continues to receive the Forgejo VM address

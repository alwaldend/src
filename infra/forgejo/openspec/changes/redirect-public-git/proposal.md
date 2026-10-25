## Why

The user requested that public git.alwaldend.com reach the public download
server and redirect to GitHub. Local Forgejo routing must remain independent.

## What Changes

- Update only the existing global git CNAME to download.alwaldend.com.
- Keep the local git A record and all other Forgejo records unchanged.
- Coordinate the Yandex-only HTTPS redirect with the download owner's
  Traefik configuration and verify the scoped DNS plan and redirect.

## Capabilities

### New Capabilities

- `public-git-routing`: Public GitHub redirection with preserved local Forgejo DNS.

### Modified Capabilities

None. Local Forgejo configuration retains its canonical hostname.

## Impact

The existing Forgejo DNS declaration and DNS-only Terraform workflow. The
download owner implements the redirect. No VM, Forgejo service, Vault role,
apex record, or private DNS change is needed.

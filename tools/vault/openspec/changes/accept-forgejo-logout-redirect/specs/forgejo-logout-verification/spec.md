## Purpose

Verify that Forgejo login cleanup invalidates the original authenticated session
across supported successful logout response formats.

## ADDED Requirements

### Requirement: Verify invalidation after HTTP 200 or 303 logout

The login wrapper SHALL accept logout HTTP 200 and 303 responses, then verify
that the retained original cookie cannot access authenticated settings. It
SHALL fail when that session still authenticates, regardless of logout status.

#### Scenario: Forgejo redirects after logout

- **WHEN** logout returns HTTP 303 and the retained cookie no longer authenticates
- **THEN** cleanup succeeds without following the logout redirect

#### Scenario: Logout leaves the original session valid

- **WHEN** logout returns a supported response but the retained cookie still authenticates
- **THEN** cleanup fails

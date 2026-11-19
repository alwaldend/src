## Context

Observed during the Forgejo src mirror setup: logout succeeds with HTTP 303,
but the wrapper accepts only HTTP 200. Token deletion precedes this step.

## Decision

Accept 200 and 303 without following logout redirects. Continue checking the
original retained session cookie against the authenticated applications page.
A redirect response alone never proves logout. Extend the existing fixture
before changing code to cover both responses and retained-session rejection.

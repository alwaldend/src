## Decision

Use stable 3.7.17, published 2026-10-30. Pin its existing Linux AMD64 archive
with SRI SHA-256. No new dependency or platform is introduced.

The migration guide's 3.7.17 item changes Ingress-NGINX generated names;
the local inference host uses the file provider. Verify native version,
shared packaging, strict client authentication, SNI, HTTP redirection, and
routing to actual T3 with disposable local credentials. This fixture does
not exercise live ACME enrollment or another deployment's full configuration.

Preserve the completed 3.7.16 extraction record. This new owner-local record
captures the independently requested follow-up without rewriting its evidence.

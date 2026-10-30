## Live acceptance: 2026-10-30

Both local and Yandex were explicitly authorized for namespace deployment and this schema publication. The existing Ansible wrapper applied only tagged publication tasks. Both hosts completed with zero failures; Nginx validated before reload. An initial trust failure changed no host state, and was corrected by using the existing scoped host CA in the owning wrapper.

Published /third_party/com_x_api_openapi/releases/2.173-4377d1791c01/openapi.json to both independent stores using the authenticated publisher. Anonymous HTTPS downloads from each host are 1,313,732 bytes with SHA-256 4377d1791c01632c9cf2860fc833a4cec76eff74f89d232fd7355e765e7bf0a9.

Root requests return HTTP 200 JSON listing only projects and third_party, without redirecting. Both namespace listings return 200. Range reads return 206. Root POST returns 405. Staging, state and sites paths return 404. Raw logs and machine-readable receipts remain under ignored out/x_schema_mirror.

## Review correction acceptance

RELEASE-SSH-AGENT-001 is fixture-tested: a built-CLI request using a public
identity and an in-memory private key first failed authentication, then passed
after restoring agent use by default. The explicit --ssh_disable_agent request
is rejected by the same local SSH endpoint. The test emits ssh-agent.json with
the public identity fingerprint and both outcomes. No private key is written.
The fixture reuses the existing golang.org/x/crypto dependency. Logs remain under
ignored out/x_sdk_review_fixes.

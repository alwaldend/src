"""Offline regressions for authentication boundaries and credential cleanup."""

import contextlib
import io
import json
import os
import unittest
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path
from types import SimpleNamespace
from unittest import mock

from infra.forgejo_runner.ci import smoke


class SmokeTest(unittest.TestCase):
    def setUp(self):
        self.stack = contextlib.ExitStack()
        self.addCleanup(self.stack.close)
        self.stack.enter_context(
            mock.patch.dict(
                os.environ,
                {
                    "CI_SERVER_URL": "https://forgejo.example.test",
                    "CI_DEFAULT_BRANCH": "main",
                    "CI_REF": "refs/heads/main",
                    "VAULT_ADDR": "https://vault.example.test",
                    "VAULT_AUTH_PATH": "forgejo_ci",
                    "VAULT_AUTH_ROLE": "secure",
                    "ACTIONS_ID_TOKEN_REQUEST_URL": (
                        "https://forgejo.example.test/api/actions/token?audience=old&job=1"
                    ),
                    "ACTIONS_ID_TOKEN_REQUEST_TOKEN": "dummy-issuance-credential",
                    "BUILD_WORKSPACE_DIRECTORY": "/fixture/workspace",
                },
                clear=True,
            )
        )
        self.stack.enter_context(
            mock.patch.object(os, "cpu_count", return_value=8)
        )
        read_text = Path.read_text

        def read_fixture(path, *args, **kwargs):
            if str(path) == "/proc/meminfo":
                return f"MemTotal: {16 * 1024**2} kB\n"
            return read_text(path, *args, **kwargs)

        self.stack.enter_context(
            mock.patch.object(Path, "read_text", read_fixture)
        )
        self.storage = self.stack.enter_context(
            mock.patch.object(
                os,
                "statvfs",
                return_value=SimpleNamespace(
                    f_blocks=500 * 1024**2,
                    f_frsize=1024,
                ),
            )
        )
        self.output = io.StringIO()
        self.stack.enter_context(contextlib.redirect_stdout(self.output))
        self.stack.enter_context(contextlib.redirect_stderr(self.output))
        self.revoked = []
        self.audiences = []
        self.metadata = {
            "policies": ["auth_token_lookup_self", "forgejo_ci_revoke_self"],
            "ttl": 300,
        }
        self.negative_error = smoke.AuthenticationError(
            400, "Vault negative login", "audience"
        )
        self.lookup_error = None
        self.http = self.stack.enter_context(
            mock.patch.object(
                smoke,
                "request_json",
                side_effect=self.exchange,
            )
        )

    def exchange(self, url, *, headers=None, payload=None, stage=None):
        endpoint = urllib.parse.urlsplit(url)
        if endpoint.netloc == "forgejo.example.test":
            query = urllib.parse.parse_qs(endpoint.query)
            self.assertEqual(["1"], query["job"])
            self.assertEqual(1, len(query["audience"]))
            self.assertEqual(
                "Bearer dummy-issuance-credential", headers["Authorization"]
            )
            audience = query["audience"][0]
            self.audiences.append(audience)
            return {
                "value": "dummy-negative-jwt"
                if audience.startswith("urn:")
                else "dummy-positive-jwt"
            }
        self.assertEqual("vault.example.test", endpoint.netloc)
        if endpoint.path == "/v1/auth/forgejo_ci/login":
            self.assertEqual("secure", payload["role"])
            if payload["jwt"] == "dummy-negative-jwt":
                if self.negative_error:
                    raise self.negative_error
                return {"auth": {"client_token": "dummy-unexpected-token"}}
            self.assertEqual("dummy-positive-jwt", payload["jwt"])
            return {"auth": {"client_token": "dummy-ci-token"}}
        if endpoint.path == "/v1/auth/token/lookup-self":
            self.assertEqual("dummy-ci-token", headers["X-Vault-Token"])
            if self.lookup_error:
                raise self.lookup_error
            return {"data": self.metadata}
        if endpoint.path == "/v1/auth/token/revoke-self":
            self.assertEqual({}, payload)
            self.revoked.append(headers["X-Vault-Token"])
            return {}
        self.fail(f"Unexpected fixture endpoint: {endpoint.path}")

    def test_success_uses_workspace_and_revokes_token(self):
        self.assertEqual(0, smoke.main())
        self.storage.assert_called_once_with("/fixture/workspace")
        self.assertEqual(
            [
                "urn:forgejo-runner:invalid-audience",
                "https://vault.example.test",
            ],
            self.audiences,
        )
        self.assertEqual(["dummy-ci-token"], self.revoked)
        self.assertIn("CI token revoked", self.output.getvalue())
        self.assertNotIn("dummy-", self.output.getvalue())

    def test_only_exact_configured_refs_authenticate(self):
        config = json.loads(
            (Path(smoke.__file__).parent.parent / "ci.json").read_text()
        )
        os.environ["CI_REF"] = "refs/heads/" + config["validation_branch"]
        self.assertEqual(0, smoke.main())
        self.http.reset_mock()
        for ref in (
            os.environ["CI_REF"] + "-other",
            "refs/tags/main",
            "refs/heads/feature",
        ):
            with self.subTest(ref=ref):
                os.environ["CI_REF"] = ref
                self.assertEqual(0, smoke.main())
                self.http.assert_not_called()

    def test_insufficient_resources_never_request_identity(self):
        with mock.patch.object(os, "cpu_count", return_value=4):
            self.assertEqual(1, smoke.main())
        self.http.assert_not_called()

    def test_unsafe_identity_endpoint_never_receives_credentials(self):
        for endpoint in (
            "http://forgejo.example.test/token",
            "https://foreign.example.test/token",
        ):
            with self.subTest(endpoint=endpoint):
                os.environ["ACTIONS_ID_TOKEN_REQUEST_URL"] = endpoint
                self.assertEqual(1, smoke.main())
                self.http.assert_not_called()

    def test_lookup_failures_still_revoke(self):
        self.lookup_error = KeyError("dummy-sensitive-response")
        self.assertEqual(1, smoke.main())
        self.assertEqual(["dummy-ci-token"], self.revoked)
        self.assertIn("KeyError", self.output.getvalue())
        self.assertNotIn("dummy-", self.output.getvalue())

    def test_policy_or_ttl_mismatch_still_revokes(self):
        for metadata in (
            {**self.metadata, "policies": ["default"]},
            {**self.metadata, "ttl": 301},
            {**self.metadata, "ttl": 0},
        ):
            with self.subTest(metadata=metadata):
                self.metadata = metadata
                self.revoked.clear()
                self.assertEqual(1, smoke.main())
                self.assertEqual(["dummy-ci-token"], self.revoked)

    def test_unexpected_negative_login_token_is_revoked(self):
        self.negative_error = None
        self.assertEqual(1, smoke.main())
        self.assertEqual(["dummy-unexpected-token"], self.revoked)
        self.assertEqual(1, len(self.audiences))

    def test_unrelated_login_failure_cannot_pass_negative_check(self):
        self.negative_error = smoke.AuthenticationError(
            400, "Vault negative login", "clock skew"
        )
        self.assertEqual(1, smoke.main())
        self.assertEqual(1, len(self.audiences))
        self.assertEqual([], self.revoked)


class TransportTest(unittest.TestCase):
    def test_redirects_are_rejected_for_every_supported_status(self):
        handler = smoke.RejectRedirects()
        for method in ("GET", "POST"):
            request = urllib.request.Request(
                "https://vault.example.test", method=method
            )
            for status in (301, 302, 303, 307, 308):
                with self.subTest(method=method, status=status):
                    with self.assertRaisesRegex(
                        RuntimeError, "redirect rejected"
                    ):
                        handler.redirect_request(
                            request,
                            None,
                            status,
                            "",
                            {},
                            "https://foreign.example.test",
                        )

    def test_http_error_reports_only_status_and_fixed_category(self):
        error = urllib.error.HTTPError(
            "https://vault.example.test",
            400,
            "dummy-sensitive-message",
            {},
            io.BytesIO(
                b'{"errors":["audience invalid: dummy-sensitive-claim"]}'
            ),
        )
        with mock.patch.object(smoke.opener, "open", side_effect=error):
            with self.assertRaises(smoke.AuthenticationError) as caught:
                smoke.request_json(
                    "https://vault.example.test", stage="Vault login"
                )
        self.assertEqual(
            "Vault login failed: HTTP 400 (audience)", str(caught.exception)
        )


if __name__ == "__main__":
    unittest.main()

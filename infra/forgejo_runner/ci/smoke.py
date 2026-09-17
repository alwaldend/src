"""Verify the runner allocation and its short-lived Forgejo/Vault identity."""

import json
import os
import sys
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path


def require(condition, message):
    if not condition:
        raise RuntimeError(message)


class AuthenticationError(RuntimeError):
    def __init__(self, status, stage, category):
        super().__init__(f"{stage} failed: HTTP {status} ({category})")
        self.status = status
        self.category = category


def error_category(content):
    # Only fixed categories leave memory; response text can contain claims.
    try:
        errors = json.loads(content).get("errors", [])
        message = " ".join(
            value for value in errors if isinstance(value, str)
        ).lower()
    except (ValueError, AttributeError, TypeError):
        message = ""
    if "not before (nbf)" in message or "issued at (iat)" in message:
        return "clock skew"
    if "expiration time (exp)" in message:
        return "expired identity"
    for name in (
        "audience",
        "issuer",
        "signature",
        "certificate",
        "oidc",
        "jwks",
    ):
        if name in message:
            return name
    for name in (
        "repository",
        "ref_type",
        "event_name",
        "workflow_ref",
        "run_id",
        "sha",
        "sub",
        "ref",
    ):
        if '"' + name + '"' in message or "'" + name + "'" in message:
            return "claim " + name
    if "claim" in message:
        return "claims"
    return "unclassified response"


class RejectRedirects(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, req, fp, code, msg, headers, newurl):
        raise RuntimeError("Authentication endpoint redirect rejected")


opener = urllib.request.build_opener(RejectRedirects())


def request_json(
    url, *, headers=None, payload=None, stage="Configuration read"
):
    body = None if payload is None else json.dumps(payload).encode()
    request = urllib.request.Request(
        url,
        data=body,
        headers={"Content-Type": "application/json", **(headers or {})},
    )
    try:
        with opener.open(request, timeout=30) as response:
            content = response.read(1048576)
            return json.loads(content) if content else {}
    except urllib.error.HTTPError as error:
        category = error_category(error.read(16384))
        raise AuthenticationError(error.code, stage, category) from None
    except urllib.error.URLError:
        raise RuntimeError(
            "Authentication endpoint could not be reached"
        ) from None


def verify():
    server_url = urllib.parse.urlsplit(os.environ["CI_SERVER_URL"])
    require(server_url.scheme == "https", "Unexpected Forgejo server URL")
    config = json.loads((Path(__file__).parent.parent / "ci.json").read_text())
    allowed_refs = {
        "refs/heads/" + os.environ["CI_DEFAULT_BRANCH"],
        "refs/heads/" + config["validation_branch"],
    }
    if os.environ["CI_REF"] not in allowed_refs:
        print(
            "Authentication smoke check does not apply to this release branch"
        )
        return

    require((os.cpu_count() or 0) >= 8, "Runner requires at least eight CPUs")
    memory = dict(
        line.split(":", 1)
        for line in Path("/proc/meminfo").read_text().splitlines()
    )
    memory_bytes = int(memory["MemTotal"].split()[0]) * 1024
    require(
        memory_bytes >= 15 * 1024**3, "Runner has less than 15 GiB usable RAM"
    )
    storage = os.statvfs(os.environ.get("BUILD_WORKSPACE_DIRECTORY", "."))
    require(
        storage.f_blocks * storage.f_frsize >= 450 * 1024**3,
        "Runner workspace has less than 450 GiB filesystem capacity",
    )
    print("Runner CPU, memory and workspace capacity checks passed")

    vault_url = os.environ["VAULT_ADDR"].rstrip("/")
    oidc_url = urllib.parse.urlsplit(
        os.environ["ACTIONS_ID_TOKEN_REQUEST_URL"]
    )
    require(
        oidc_url.scheme == "https" and oidc_url.netloc == server_url.netloc,
        "Unexpected Forgejo OIDC endpoint",
    )
    query = urllib.parse.parse_qsl(oidc_url.query, keep_blank_values=True)
    query = [(key, value) for key, value in query if key != "audience"]

    def identity(audience):
        return request_json(
            urllib.parse.urlunsplit(
                oidc_url._replace(
                    query=urllib.parse.urlencode(
                        [*query, ("audience", audience)]
                    )
                )
            ),
            headers={
                "Authorization": "Bearer "
                + os.environ["ACTIONS_ID_TOKEN_REQUEST_TOKEN"],
            },
            stage="Forgejo identity issuance",
        )["value"]

    login_url = f"{vault_url}/v1/auth/{os.environ['VAULT_AUTH_PATH']}/login"
    unauthorized_identity = identity("urn:forgejo-runner:invalid-audience")
    try:
        unexpected = request_json(
            login_url,
            payload={
                "role": os.environ["VAULT_AUTH_ROLE"],
                "jwt": unauthorized_identity,
            },
            stage="Vault negative login",
        )
    except AuthenticationError as error:
        require(
            error.status in (400, 403), "Unexpected audience rejection status"
        )
        require(error.category == "audience", str(error))
    else:
        request_json(
            f"{vault_url}/v1/auth/token/revoke-self",
            headers={"X-Vault-Token": unexpected["auth"]["client_token"]},
            payload={},
            stage="Unexpected token revocation",
        )
        raise RuntimeError("Vault accepted an unauthorized audience")
    print("Vault rejected an unauthorized audience")

    login = request_json(
        login_url,
        payload={
            "role": os.environ["VAULT_AUTH_ROLE"],
            "jwt": identity(vault_url),
        },
        stage="Vault login",
    )
    token = login["auth"]["client_token"]
    headers = {"X-Vault-Token": token}
    try:
        metadata = request_json(
            f"{vault_url}/v1/auth/token/lookup-self",
            headers=headers,
            stage="Vault token lookup",
        )["data"]
        require(
            set(metadata["policies"])
            == {"auth_token_lookup_self", "forgejo_ci_revoke_self"},
            "CI token has unexpected policies",
        )
        require(0 < metadata["ttl"] <= 300, "CI token has unexpected TTL")
        print("Vault OIDC authentication and restricted token checks passed")
    finally:
        request_json(
            f"{vault_url}/v1/auth/token/revoke-self",
            headers=headers,
            payload={},
            stage="Vault token revocation",
        )
    print("CI token revoked")


def main():
    try:
        verify()
    except Exception as error:
        # Do not print response bodies, request headers or credentials.
        if isinstance(error, RuntimeError):
            print(str(error), file=sys.stderr)
        else:
            print(
                f"Smoke check failed: {type(error).__name__}", file=sys.stderr
            )
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())

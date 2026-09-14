"""Publish a validation commit without rewriting history and observe CI."""

import argparse
import base64
import datetime
import json
import os
import re
import subprocess
import sys
import time
import urllib.error
import urllib.parse
import urllib.request


class ValidationError(Exception):
    """A controlled, credential-free diagnostic."""


class RejectRedirects(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, req, fp, code, msg, headers, newurl):
        raise ValidationError("API redirect rejected")


class Forgejo:
    def __init__(self, instance, token):
        parsed = urllib.parse.urlsplit(instance)
        if (
            parsed.scheme != "https"
            or not parsed.hostname
            or parsed.username
            or parsed.password
            or parsed.query
            or parsed.fragment
            or parsed.path not in ("", "/")
        ):
            raise ValidationError("Forgejo must have an HTTPS origin URL")
        self.origin = instance.rstrip("/")
        self.token = token
        self.opener = urllib.request.build_opener(RejectRedirects())

    def get(self, path, allow_missing=False, query=None):
        url = self.origin + "/api/v1/" + path
        if query:
            url += "?" + urllib.parse.urlencode(query)
        request = urllib.request.Request(
            url,
            headers={
                "Authorization": "token " + self.token,
                "Accept": "application/json",
            },
        )
        try:
            with self.opener.open(request, timeout=10) as response:
                body = response.read(1024 * 1024 + 1)
        except urllib.error.HTTPError as error:
            if allow_missing and error.code == 404:
                return None
            raise ValidationError(f"API returned HTTP {error.code}") from None
        except (urllib.error.URLError, TimeoutError):
            raise ValidationError("API connection failed") from None
        if len(body) > 1024 * 1024:
            raise ValidationError("API response exceeds size limit")
        try:
            return json.loads(body)
        except (ValueError, UnicodeError):
            raise ValidationError("API returned invalid JSON") from None


def git(workspace, *arguments, environment=None, timeout=30):
    result = subprocess.run(
        ["git", "-C", workspace, *arguments],
        env=environment,
        capture_output=True,
        text=True,
        timeout=timeout,
        check=False,
    )
    if result.returncode:
        raise ValidationError("Git operation failed; output suppressed")
    return result.stdout.strip()


def candidate_branch(workspace, candidate):
    if not re.fullmatch(r"[0-9a-f]{40}", candidate):
        raise ValidationError("Candidate must be a full commit SHA")
    resolved = git(workspace, "rev-parse", "--verify", candidate + "^{commit}")
    if resolved != candidate:
        raise ValidationError("Candidate is not an exact commit object")
    config = json.loads(
        git(workspace, "show", candidate + ":infra/forgejo_runner/ci.json")
    )
    branch = config.get("validation_branch", "")
    if not re.fullmatch(r"releases/[A-Za-z0-9_-]+", branch):
        raise ValidationError("Candidate has an invalid validation branch")
    git(
        workspace,
        "cat-file",
        "-e",
        candidate + ":.forgejo/workflows/secure.yaml",
    )
    return branch


def require_release_protection(rules):
    matching = [
        rule
        for rule in rules
        if rule.get("rule_name", rule.get("branch_name")) == "releases/*"
    ]
    if len(matching) != 1:
        raise ValidationError("Expected one releases/* protection rule")
    rule = matching[0]
    if not (
        rule.get("enable_push") is False
        or (
            rule.get("enable_push") is True
            and rule.get("enable_push_whitelist") is True
            and (
                rule.get("push_whitelist_teams")
                or rule.get("push_whitelist_usernames")
            )
        )
    ):
        raise ValidationError("Release protection does not restrict pushes")


def require_branch(published, branch, candidate):
    if (
        published.get("name") != branch
        or published.get("commit", {}).get("id") != candidate
        or published.get("protected") is not True
        or published.get("effective_branch_protection_name") != "releases/*"
    ):
        raise ValidationError("Published branch protection or SHA mismatch")


def push_candidate(
    workspace, candidate, branch, api, login, previous_commit=""
):
    environment = {
        key: value
        for key, value in os.environ.items()
        if not key.startswith(("GIT_TRACE", "GIT_CONFIG_"))
        and key != "GIT_CURL_VERBOSE"
    }
    authorization = base64.b64encode(
        (login + ":" + api.token).encode()
    ).decode()
    environment.update(
        {
            "GIT_TERMINAL_PROMPT": "0",
            "GIT_CONFIG_COUNT": "1",
            "GIT_CONFIG_KEY_0": "http." + api.origin + "/.extraHeader",
            "GIT_CONFIG_VALUE_0": "Authorization: Basic " + authorization,
        }
    )
    git(
        workspace,
        "-c",
        "credential.helper=",
        "-c",
        "core.hooksPath=/dev/null",
        "-c",
        "http.followRedirects=false",
        "push",
        "--porcelain",
        "--force-with-lease=refs/heads/" + branch + ":" + previous_commit,
        api.origin + "/" + api.repository + ".git",
        candidate + ":refs/heads/" + branch,
        environment=environment,
        timeout=60,
    )


def is_workflow(value):
    return value in ("secure.yaml", ".forgejo/workflows/secure.yaml")


def matching_tasks(api, prefix, candidate, branch, run_number):
    tasks = []
    for page in range(1, 5):
        response = api.get(
            prefix + "/actions/tasks", query={"page": page, "limit": 50}
        )
        entries = response["workflow_runs"]
        tasks.extend(
            {"id": task["id"], "status": task["status"]}
            for task in entries
            if task.get("head_sha") == candidate
            and task.get("head_branch") in (branch, "refs/heads/" + branch)
            and is_workflow(task.get("workflow_id"))
            and task.get("run_number") == run_number
        )
        if page * 50 >= response["total_count"] or len(entries) < 50:
            return tasks
    if not tasks:
        raise ValidationError("No matching jobs in the latest 200 tasks")
    return tasks


def observe_run(api, prefix, candidate, branch, timeout, receipt):
    deadline = time.monotonic() + timeout
    terminal = {"success", "failure", "cancelled", "skipped", "blocked"}
    while time.monotonic() < deadline:
        response = api.get(
            prefix + "/actions/runs",
            query={
                "head_sha": candidate,
                "ref": "refs/heads/" + branch,
                "workflow_id": "secure.yaml",
                "event": "push",
                "limit": 10,
            },
        )
        runs = response["workflow_runs"]
        if response["total_count"] > 1 or len(runs) > 1:
            raise ValidationError("More than one matching workflow run")
        if runs:
            run = runs[0]
            if (
                run.get("commit_sha") != candidate
                or run.get("prettyref") not in (branch, "refs/heads/" + branch)
                or not is_workflow(run.get("workflow_id"))
                or run.get("event") != "push"
            ):
                raise ValidationError("Workflow run does not match candidate")
            receipt.update(
                run_id=run["id"],
                run_number=run["index_in_repo"],
                run_status=run["status"],
                run_url=(
                    api.origin
                    + "/"
                    + api.repository
                    + "/actions/runs/"
                    + str(run["index_in_repo"])
                ),
            )
            if run["status"] in terminal:
                receipt["jobs"] = matching_tasks(
                    api, prefix, candidate, branch, run["index_in_repo"]
                )
                if run["status"] != "success":
                    raise ValidationError("Workflow did not succeed")
                if receipt["jobs"] and all(
                    job["status"] == "success" for job in receipt["jobs"]
                ):
                    return
        time.sleep(min(5, max(0, deadline - time.monotonic())))
    raise ValidationError("Timed out waiting for successful run and jobs")


def validate(arguments, receipt):
    token = os.environ.get("FORGEJO_API_TOKEN")
    if not token:
        raise ValidationError("FORGEJO_API_TOKEN is required")
    if not re.fullmatch(
        r"[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+", arguments.repository
    ):
        raise ValidationError("Invalid repository name")
    if not 30 <= arguments.timeout <= 300:
        raise ValidationError("Timeout must be between 30 and 300 seconds")
    branch = candidate_branch(arguments.workspace, arguments.commit)
    receipt.update(candidate=arguments.commit, branch=branch)
    api = Forgejo(arguments.instance, token)
    api.repository = arguments.repository
    prefix = "repos/" + api.repository
    branch_path = prefix + "/branches/" + urllib.parse.quote(branch, safe="")
    identity = api.get("user")
    if not identity.get("is_admin"):
        raise ValidationError("The controller identity must be a site admin")
    require_release_protection(api.get(prefix + "/branch_protections"))
    published = api.get(branch_path, allow_missing=True)
    previous_commit = arguments.previous_commit
    if previous_commit and not re.fullmatch(r"[0-9a-f]{40}", previous_commit):
        raise ValidationError("Previous commit must be a full commit SHA")
    if published is None and previous_commit:
        raise ValidationError("Expected validation branch no longer exists")
    advance = published is not None and (
        published.get("commit", {}).get("id") != arguments.commit
    )
    if advance:
        if not previous_commit:
            raise ValidationError("A different tip requires --previous-commit")
        require_branch(published, branch, previous_commit)
        try:
            git(
                arguments.workspace,
                "merge-base",
                "--is-ancestor",
                previous_commit,
                arguments.commit,
            )
        except ValidationError:
            raise ValidationError(
                "Previous commit is not a known ancestor of the candidate"
            ) from None
        receipt["previous_commit"] = previous_commit
    if published is None or advance:
        receipt["stage"] = "push"
        push_candidate(
            arguments.workspace,
            arguments.commit,
            branch,
            api,
            identity["login"],
            previous_commit=previous_commit,
        )
        receipt["pushed"] = True
        published = api.get(branch_path)
    receipt["stage"] = "branch_protection"
    require_branch(published, branch, arguments.commit)
    receipt["protected"] = True
    receipt["stage"] = "actions"
    observe_run(
        api, prefix, arguments.commit, branch, arguments.timeout, receipt
    )
    receipt["stage"] = "complete"


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--workspace", required=True)
    parser.add_argument("--commit", required=True)
    parser.add_argument("--previous-commit", default="")
    parser.add_argument("--instance", required=True)
    parser.add_argument("--repository", required=True)
    parser.add_argument("--timeout", type=int, default=240)
    arguments = parser.parse_args()
    receipt = {"stage": "preflight", "pushed": False, "protected": False}
    status = 0
    try:
        validate(arguments, receipt)
    except ValidationError as error:
        receipt["error"] = str(error)
        status = 1
    except Exception as error:
        receipt["error"] = "Validation failed: " + type(error).__name__
        status = 1
    receipt["observed_at"] = datetime.datetime.now(
        datetime.timezone.utc
    ).isoformat()
    print(json.dumps(receipt, sort_keys=True))
    return status


if __name__ == "__main__":
    sys.exit(main())

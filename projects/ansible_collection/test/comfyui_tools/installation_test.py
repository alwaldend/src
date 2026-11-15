"""Run the packaged ComfyUI asset role against disposable local files and HTTP."""

import argparse
import hashlib
import http.server
import json
import os
import pathlib
import subprocess
import sys
import tempfile
import threading

from python.runfiles import runfiles


class _Service(http.server.BaseHTTPRequestHandler):
    """Serve aggregate readiness and tiny public model fixtures only."""

    _busy: bool = False
    """Fixture queue state, without any job data."""

    def do_GET(self) -> None:
        """Return a fixture model or a ComfyUI aggregate endpoint."""
        payload: bytes
        if self.path == "/prompt":
            payload = json.dumps(
                {"exec_info": {"queue_remaining": int(self._busy)}}
            ).encode()
        elif self.path == "/object_info":
            payload = b'{"FixtureNode":{},"FixtureGitNode":{}}'
        else:
            payload = b"public model fixture"
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(payload)))
        self.end_headers()
        self.wfile.write(payload)

    def log_message(self, format: str, *args: object) -> None:
        """Suppress disposable request diagnostics."""


def _command(arguments: list[str], directory: pathlib.Path) -> str:
    """Run a bounded local fixture command and retain its failure diagnostics."""
    result: subprocess.CompletedProcess[str] = subprocess.run(
        arguments,
        cwd=directory,
        text=True,
        capture_output=True,
        timeout=120,
        check=True,
    )
    return result.stdout.strip()


def _exercise(ansible: str, root: pathlib.Path) -> None:
    """Verify installation, reapplication, update safety and failure containment."""
    application: pathlib.Path = root / "alternate-workspace/ComfyUI"
    sources: pathlib.Path = root / "sources"
    sources.mkdir()
    (application / "custom_nodes").mkdir(parents=True)
    (application / "main.py").write_text("# existing ComfyUI image\n")
    (sources / "node.py").write_text(
        "NODE_CLASS_MAPPINGS = {'FixtureNode': object}\n"
    )
    (sources / "web.js").write_text("// browser fixture\n")
    graph: str = json.dumps({"version": 0.4, "nodes": [], "links": []})
    (sources / "fixture.app.json").write_text(graph)
    repository: pathlib.Path = sources / "git-node"
    repository.mkdir()
    _command(["git", "init", "--quiet"], repository)
    (repository / "__init__.py").write_text("NODE_CLASS_MAPPINGS = {}\n")
    _command(["git", "add", "__init__.py"], repository)
    _command(
        [
            "git",
            "-c",
            "user.name=Fixture",
            "-c",
            "user.email=fixture@example.invalid",
            "commit",
            "--quiet",
            "-m",
            "Fixture node",
        ],
        repository,
    )
    revision: str = _command(["git", "rev-parse", "HEAD"], repository)
    restart: pathlib.Path = sources / "restart.py"
    restart.write_text(
        "import pathlib,sys\np=pathlib.Path(sys.argv[1]);p.write_text(p.read_text()+'restart\\n' if p.exists() else 'restart\\n')\n"
    )
    server: http.server.ThreadingHTTPServer = http.server.ThreadingHTTPServer(
        ("127.0.0.1", 0), _Service
    )
    thread: threading.Thread = threading.Thread(
        target=server.serve_forever, daemon=True
    )
    thread.start()
    endpoint: str = f"http://127.0.0.1:{server.server_port}"
    playbook: pathlib.Path = root / "playbook.yaml"
    playbook.write_text(
        "- hosts: localhost\n  connection: local\n  gather_facts: false\n  roles:\n    - alwaldend.main.comfyui_tools\n"
    )
    environment: dict[str, str] = dict(os.environ)
    environment.update(
        {
            "ANSIBLE_CONFIG": str(root / "ansible.cfg"),
            "ANSIBLE_COLLECTIONS_PATH": str(
                pathlib.Path(os.environ["RUNFILES_DIR"]) / "_main/collections"
            ),
            "ANSIBLE_LOCAL_TEMP": str(root / "controller-tmp"),
            "ANSIBLE_REMOTE_TEMP": str(root / "remote-tmp"),
            "ANSIBLE_NOCOLOR": "1",
            "GIT_CONFIG_NOSYSTEM": "1",
            "GIT_CONFIG_GLOBAL": "/dev/null",
        }
    )
    (root / "ansible.cfg").write_text(
        "[defaults]\nretry_files_enabled=False\nlocalhost_warning=False\n"
    )
    variables: dict[str, object] = {
        "ansible_python_interpreter": sys.executable,
        "comfyui_tools_directory": str(application),
        "comfyui_tools_state_directory": str(root / "state"),
        "comfyui_tools_environment_file": str(root / "environment"),
        "comfyui_tools_api_url": endpoint,
        "comfyui_tools_idle_retries": 1,
        "comfyui_tools_idle_delay": 0,
        "comfyui_tools_restart_command": [
            sys.executable,
            str(restart),
            str(root / "restarts"),
        ],
        "comfyui_tools_arguments": "--listen 127.0.0.1 --disable-all-custom-nodes --offline",
        "comfyui_tools_model_directories": {
            "checkpoint": "models/checkpoints",
            "lora": "models/loras",
            "controlnet": "models/controlnet",
        },
        "comfyui_tools_models": [
            {
                "filename": name + ".safetensors",
                "kind": name,
                "url": endpoint + "/" + name,
                "sha256": hashlib.sha256(b"public model fixture").hexdigest(),
            }
            for name in ("checkpoint", "lora", "controlnet")
        ],
        "comfyui_tools_git_nodes": {
            "FixtureGit": {"url": repository.as_uri(), "ref": revision}
        },
        "comfyui_tools_node_sources": {
            "fixture_node": str(sources / "node.py")
        },
        "comfyui_tools_node_web_scripts": {
            "fixture_node": str(sources / "web.js")
        },
        "comfyui_tools_workflows": [str(sources / "fixture.app.json")],
        "comfyui_tools_required_nodes": ["FixtureNode", "FixtureGitNode"],
    }

    def apply(*extra: str, success: bool = True) -> str:
        """Apply the actual role only to this fixture and check its exit status."""
        path: pathlib.Path = root / "variables.json"
        path.write_text(json.dumps(variables))
        result: subprocess.CompletedProcess[str] = subprocess.run(
            [
                ansible,
                "-i",
                "localhost,",
                str(playbook),
                "-e",
                "@" + str(path),
                *extra,
            ],
            cwd=root,
            env=environment,
            text=True,
            capture_output=True,
            timeout=180,
        )
        if (result.returncode == 0) != success:
            raise AssertionError(result.stdout + result.stderr)
        return result.stdout

    try:
        apply()
        for kind, directory in (
            ("checkpoint", "checkpoints"),
            ("lora", "loras"),
            ("controlnet", "controlnet"),
        ):
            assert (
                application / f"models/{directory}/{kind}.safetensors"
            ).read_bytes() == b"public model fixture"
        assert (
            application / "user/default/workflows/fixture.app.json"
        ).read_text() == graph
        assert (
            application / "custom_nodes/fixture_node/web/index.js"
        ).read_text() == "// browser fixture\n"
        assert (
            _command(
                ["git", "rev-parse", "HEAD"],
                application / "custom_nodes/FixtureGit",
            )
            == revision
        )
        assert (
            "--whitelist-custom-nodes FixtureGit fixture_node"
            in (root / "environment").read_text()
        )
        assert (root / "restarts").read_text() == "restart\n"
        apply()
        assert (root / "restarts").read_text() == "restart\n"
        workflow: pathlib.Path = (
            application / "user/default/workflows/fixture.app.json"
        )
        workflow.write_text("operator workflow edit\n")
        (sources / "node.py").write_text(
            "NODE_CLASS_MAPPINGS = {'FixtureNode': int}\n"
        )
        apply("--skip-tags", "workflows")
        assert workflow.read_text() == "operator workflow edit\n"
        assert (root / "restarts").read_text() == "restart\nrestart\n"
        _Service._busy = True
        (sources / "node.py").write_text(
            "BUSY_UPDATE_MUST_NOT_BE_INSTALLED = True\n"
        )
        apply(success=False)
        assert (
            "BUSY_UPDATE"
            not in (
                application / "custom_nodes/fixture_node/__init__.py"
            ).read_text()
        )
        _Service._busy = False
        variables["comfyui_tools_models"] = [
            {
                "filename": "invalid.safetensors",
                "kind": "lora",
                "url": endpoint + "/model",
                "sha256": "0" * 64,
            }
        ]
        apply(success=False)
        assert not (application / "models/loras/invalid.safetensors").exists()
        variables["comfyui_tools_models"] = [
            {
                "filename": "../escaped.safetensors",
                "kind": "lora",
                "url": endpoint + "/model",
            }
        ]
        apply(success=False)
        assert not (application / "models/escaped.safetensors").exists()
        evidence: pathlib.Path = (
            pathlib.Path(os.environ["TEST_UNDECLARED_OUTPUTS_DIR"])
            / "installation.json"
        )
        evidence.write_text(
            json.dumps(
                {
                    "alternate_workspace": True,
                    "model_kinds": ["checkpoint", "lora", "controlnet"],
                    "git_pin": revision,
                    "unchanged_restarts": 0,
                    "workflow_edits_preserved": True,
                    "busy_update_refused": True,
                    "checksum_failure_refused": True,
                    "path_escape_refused": True,
                    "live_operations": 0,
                },
                indent=2,
            )
            + "\n"
        )
    finally:
        server.shutdown()
        server.server_close()
        thread.join(timeout=5)


def _main() -> None:
    """Resolve the packaged launcher and own one disposable local installation."""
    parser: argparse.ArgumentParser = argparse.ArgumentParser()
    parser.add_argument("--ansible", required=True)
    options: argparse.Namespace = parser.parse_args()
    resolver: runfiles.Runfiles | None = runfiles.Create()
    assert resolver is not None
    ansible: str | None = resolver.Rlocation(options.ansible)
    assert ansible is not None
    with tempfile.TemporaryDirectory(
        prefix="comfyui-", dir=os.environ["TEST_TMPDIR"]
    ) as directory:
        _exercise(ansible, pathlib.Path(directory))


if __name__ == "__main__":
    _main()

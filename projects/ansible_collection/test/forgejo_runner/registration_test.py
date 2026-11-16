"""Exercise packaged runner registration and template rendering without host changes."""

import argparse
import json
import os
import pathlib
import pwd
import shutil
import subprocess
import sys
import tempfile

from python.runfiles import runfiles


def _exercise(
    ansible: pathlib.Path, playbook: pathlib.Path, root: pathlib.Path
) -> None:
    """Run success, repeat and failure cases with public dummy credentials only."""
    collection: pathlib.Path = (
        pathlib.Path(os.environ["RUNFILES_DIR"]) / "_main/collections"
    )
    packaged_role: pathlib.Path = (
        collection / "ansible_collections/alwaldend/main/roles/forgejo_runner"
    )
    local_collection: pathlib.Path = root / "collections"
    role: pathlib.Path = (
        local_collection
        / "ansible_collections/alwaldend/main/roles/forgejo_runner"
    )
    shutil.copytree(
        packaged_role, role, ignore=shutil.ignore_patterns("files")
    )
    handlers: pathlib.Path = role / "handlers/main.yaml"
    handlers.write_text(
        "- name: Restart the Forgejo runner service\n"
        "  ansible.builtin.command:\n"
        "    argv: ['{{ forgejo_runner_dir }}/forgejo-runner', '--restart']\n"
    )
    (root / "ansible.cfg").write_text(
        "[defaults]\nretry_files_enabled=False\n"
    )
    commands: pathlib.Path = root / "commands"
    commands.mkdir()
    (commands / "python3").symlink_to(pathlib.Path(sys.executable).resolve())
    environment: dict[str, str] = {
        key: value
        for key, value in os.environ.items()
        if key.startswith(("RUNFILES_", "TEST_"))
        or key in ("PATH", "LANG", "LC_ALL")
    }
    environment.update(
        {
            "HOME": str(root),
            "ANSIBLE_CONFIG": str(root / "ansible.cfg"),
            "ANSIBLE_COLLECTIONS_PATH": str(local_collection),
            "ANSIBLE_LOCAL_TEMP": str(root / "controller-tmp"),
            "ANSIBLE_REMOTE_TEMP": str(root / "remote-tmp"),
            "ANSIBLE_NOCOLOR": "1",
            "PATH": str(commands) + os.pathsep + os.environ["PATH"],
        }
    )
    variables: dict[str, object] = {
        "ansible_python_interpreter": sys.executable,
        "ansible_async_dir": str(root / "async"),
        "forgejo_runner_user": pwd.getpwuid(os.getuid()).pw_name,
        "forgejo_runner_instance": "https://forgejo.invalid",
        "forgejo_runner_token": "fixture-registration-token",
        "forgejo_runner_name": "fixture-worker",
        "fixture_role": str(role),
    }

    def install(name: str, mode: str = "success") -> pathlib.Path:
        """Prepare an executable that checks actual command arguments and stdin."""
        directory: pathlib.Path = root / name
        (directory / "work").mkdir(parents=True)
        executable: pathlib.Path = directory / "forgejo-runner"
        executable.write_text(
            "#!/usr/bin/env python3\n"
            "import json,pathlib,sys,time\n"
            "root=pathlib.Path(__file__).parent\n"
            "events=root/'events'\n"
            "def record(event):\n"
            "    with events.open('a') as output: output.write(event+'\\n')\n"
            "if sys.argv[1:]==['--restart']:\n"
            "    record('restart');sys.exit(0)\n"
            "assert sys.argv[1:]==['--config',str(root/'config.yaml'),'register']\n"
            "assert pathlib.Path.cwd()==root/'work'\n"
            "instance,token,name=sys.stdin.read().splitlines()\n"
            "assert instance=='https://forgejo.invalid'\n"
            "assert token=='fixture-registration-token' and name=='fixture-worker'\n"
            "mode=(root/'mode').read_text()\n"
            "if mode=='failure':\n"
            "    print(token);print(token,file=sys.stderr);sys.exit(1)\n"
            "if mode=='timeout': time.sleep(30)\n"
            "(root/'work/.runner').write_text(json.dumps({'token':token,'name':name}))\n"
            "record('register')\n"
        )
        executable.chmod(0o700)
        (directory / "mode").write_text(mode)
        variables["forgejo_runner_dir"] = str(directory)
        return directory

    def apply(success: bool = True) -> str:
        """Execute the real tasks and reject any disclosure of dummy credentials."""
        inputs: pathlib.Path = root / "variables.json"
        inputs.write_text(json.dumps(variables))
        result: subprocess.CompletedProcess[str] = subprocess.run(
            [
                str(ansible),
                "-i",
                "localhost,",
                str(playbook),
                "-e",
                "@" + str(inputs),
            ],
            cwd=root,
            env=environment,
            text=True,
            capture_output=True,
            timeout=90,
            check=False,
        )
        output: str = result.stdout + result.stderr
        assert "fixture-registration-token" not in output, (
            "Registration token disclosed"
        )
        assert (result.returncode == 0) == success, output
        return output

    directory: pathlib.Path = install("alternate-installation")
    apply()
    registration: pathlib.Path = directory / "work/.runner"
    original: bytes = registration.read_bytes()
    assert registration.stat().st_mode & 0o777 == 0o600
    assert (directory / "events").read_text() == "register\nrestart\n"
    unit: str = (directory / "forgejo-runner.service").read_text()
    assert (
        f"ExecStart={directory}/forgejo-runner --config {directory}/config.yaml daemon"
        in unit
    )
    assert f"WorkingDirectory={directory}/work" in unit
    assert f"Environment=HOME={directory}/home" in unit
    del variables["forgejo_runner_token"]
    del variables["forgejo_runner_instance"]
    apply()
    assert registration.read_bytes() == original
    assert (directory / "events").read_text() == "register\nrestart\n"
    variables.update(
        {
            "forgejo_runner_labels": ["custom:host"],
            "forgejo_runner_capacity": 2,
            "forgejo_runner_config": {
                "runner": {
                    "file": "/incorrect",
                    "labels": ["incorrect"],
                    "capacity": 99,
                    "timeout": "4h",
                },
                "cache": {"enabled": True},
                "host": {"workdir_parent": "/incorrect"},
            },
            "fixture_timeout": "4h",
            "fixture_cache": True,
        }
    )
    apply()
    assert registration.read_bytes() == original
    install("missing-inputs")
    apply(success=False)
    variables["forgejo_runner_instance"] = "https://forgejo.invalid"
    variables["forgejo_runner_token"] = "fixture-registration-token"
    failure: pathlib.Path = install("failed-registration", "failure")
    assert "Runner registration failed" in apply(success=False)
    assert not (failure / "work/.runner").exists()
    variables["forgejo_runner_registration_timeout"] = 1
    expired: pathlib.Path = install("expired-registration", "timeout")
    assert "Runner registration failed" in apply(success=False)
    assert not (expired / "work/.runner").exists()
    evidence: pathlib.Path = pathlib.Path(
        os.environ["TEST_UNDECLARED_OUTPUTS_DIR"]
    )
    for name in ("config.yaml", "forgejo-runner.service"):
        shutil.copyfile(directory / name, evidence / name)
    (evidence / "registration.json").write_text(
        json.dumps(
            {
                "initial_registration": True,
                "repeat_without_token": True,
                "mode": "0600",
                "managed_paths": True,
                "advanced_settings": True,
                "missing_inputs_rejected": True,
                "failure_output_suppressed": True,
                "registration_deadline": True,
                "live_operations": 0,
            },
            indent=2,
        )
        + "\n"
    )


def _main() -> None:
    """Resolve Bazel inputs and confine every operation to temporary test storage."""
    parser: argparse.ArgumentParser = argparse.ArgumentParser()
    parser.add_argument("--ansible", required=True)
    parser.add_argument("--playbook", required=True)
    options: argparse.Namespace = parser.parse_args()
    resolver: runfiles.Runfiles | None = runfiles.Create()
    assert resolver is not None
    ansible: str | None = resolver.Rlocation(options.ansible)
    playbook: str | None = resolver.Rlocation(options.playbook)
    assert ansible is not None and playbook is not None
    with tempfile.TemporaryDirectory(
        prefix="runner-", dir=os.environ["TEST_TMPDIR"]
    ) as directory:
        _exercise(
            pathlib.Path(ansible),
            pathlib.Path(playbook),
            pathlib.Path(directory),
        )


if __name__ == "__main__":
    _main()

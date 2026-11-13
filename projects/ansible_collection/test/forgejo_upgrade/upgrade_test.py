"""Exercise packaged Forgejo role transactions with real local SQLite instances."""

import argparse
import base64
import contextlib
import json
import os
import pathlib
import pwd
import secrets
import shutil
import signal
import socket
import sqlite3
import subprocess
import sys
import tarfile
import tempfile
import time
import unittest
import urllib.error
import urllib.request

from python.runfiles import runfiles

_OPTIONS: argparse.Namespace
"""Bazel-declared executables and fixture resources selected by main."""


class UpgradeTest(unittest.TestCase):
    """Verify preserved data and failure behavior through actual Ansible tasks."""

    def setUp(self) -> None:
        """Create an isolated account-free service tree and private test config."""
        self._root: pathlib.Path = pathlib.Path(
            tempfile.mkdtemp(prefix="forgejo-", dir=os.environ["TEST_TMPDIR"])
        )
        """Disposable fixture tree; no production files are reachable through its variables."""
        self._events: pathlib.Path = self._root / "events"
        """Sanitized service lifecycle event stream."""
        self._state: pathlib.Path = self._root / "instance"
        """Complete local Forgejo persistence root."""
        self._binary: pathlib.Path = self._root / "forgejo"
        """Installed fixture executable, distinct from either immutable Bazel input."""
        self._config: pathlib.Path = self._state / "configs/forgejo.ini"
        """Installed runtime configuration, including disposable generated keys."""
        self._backups: pathlib.Path = self._root / "backups"
        """Protected upgrade recovery directories retained within the fixture."""
        self._token: str = ""
        """Ephemeral local API credential; never included in evidence or command output."""
        self._private_values: list[str] = []
        """File-backed fixture secrets retained only for diagnostic redaction."""
        self._observed_versions: list[str] = []
        """Versions returned by accepted real service processes, for retained evidence."""
        self._observed_schemas: list[dict[str, int]] = []
        """Observed database migration markers; a release need not change them."""
        with socket.socket() as reservation:
            reservation.bind(("127.0.0.1", 0))
            self._port: int = reservation.getsockname()[1]
            """Loopback port selected for this fixture only."""
        self._environment: dict[str, str] = {
            key: value
            for key, value in os.environ.items()
            if key.startswith(("RUNFILES_", "TEST_"))
            or key in ("PATH", "LANG", "LC_ALL")
        }
        """Restricted subprocess environment without inherited production credentials."""
        for name in ("tmp", "home", "commands"):
            (self._root / name).mkdir(mode=0o700)
        fixture: pathlib.Path = self._root / "commands/systemctl"
        (fixture.parent / "python3").symlink_to(
            pathlib.Path(sys.executable).resolve()
        )
        fixture.write_text(
            "#!/usr/bin/env python3\n"
            + pathlib.Path(_OPTIONS.service_fixture).read_text()
        )
        fixture.chmod(0o700)
        self._environment.update(
            {
                "FORGEJO_TEST_ROOT": str(self._root),
                "HOME": str(self._root / "home"),
                "TMPDIR": str(self._root / "tmp"),
                "PATH": str(fixture.parent) + os.pathsep + os.environ["PATH"],
                "ANSIBLE_COLLECTIONS_PATH": _OPTIONS.collections,
                "ANSIBLE_LOCAL_TEMP": str(self._root / "tmp/controller"),
                "ANSIBLE_REMOTE_TEMP": str(self._root / "tmp/remote"),
                "ANSIBLE_NOCOLOR": "1",
                "ANSIBLE_CONFIG": str(self._root / "ansible.cfg"),
            }
        )
        (self._root / "ansible.cfg").write_text(
            "[defaults]\nretry_files_enabled=False\nlocalhost_warning=False\n"
        )
        (self._root / "playbook.yaml").write_text(
            "- hosts: localhost\n  gather_facts: false\n  connection: local\n  tasks:\n"
            "    - block:\n        - ansible.builtin.include_role:\n            name: alwaldend.main.forgejo\n"
            "      rescue:\n        - ansible.builtin.copy:\n"
            "            content: \"{{ forgejo_doctor_result.stdout | default('') }}\\n{{ forgejo_doctor_result.stderr | default('') }}\"\n"
            "            dest: \"{{ forgejo_fixture_root }}/doctor-output\"\n            mode: '0600'\n          no_log: true\n"
            "        - ansible.builtin.fail:\n            msg: Fixture role invocation failed.\n"
        )
        self._template: pathlib.Path = self._root / "forgejo.ini.j2"
        """Controller-side configuration with only runtime-generated fixture secrets."""
        self._template.write_text(
            f"APP_NAME = Fixture\nRUN_MODE = prod\nRUN_USER = {pwd.getpwuid(os.getuid()).pw_name}\n"
            f"WORK_PATH = {self._state}/work\n[server]\nHTTP_ADDR = 127.0.0.1\nHTTP_PORT = {self._port}\n"
            f"ROOT_URL = http://127.0.0.1:{self._port}/\nAPP_DATA_PATH = {self._state}/work/data\n"
            "DISABLE_SSH = true\nOFFLINE_MODE = true\n[database]\nDB_TYPE = sqlite3\n"
            f"PATH = {self._state}/work/data/forgejo.db\n[security]\nINSTALL_LOCK = true\n"
            f"SECRET_KEY = {secrets.token_hex(32)}\nINTERNAL_TOKEN = {secrets.token_urlsafe(64)}\n"
            f"[oauth2]\nJWT_SECRET = {secrets.token_urlsafe(32)}\n"
            f"[repository]\nROOT = {self._state}/work/data/forgejo-repositories\n"
            "[log]\nMODE = console\nLEVEL = Error\n[service]\nDISABLE_REGISTRATION = true\n"
        )
        self._variables: dict[str, object] = {
            "ansible_python_interpreter": sys.executable,
            "ansible_remote_tmp": str(self._root / "tmp/remote"),
            "forgejo_fixture_root": str(self._root),
            "forgejo_manage_packages": False,
            "forgejo_manage_user": False,
            "forgejo_owner": pwd.getpwuid(os.getuid()).pw_name,
            "forgejo_user": pwd.getpwuid(os.getuid()).pw_name,
            "forgejo_group": pwd.getpwuid(os.getuid()).pw_name,
            "forgejo_dir": str(self._state),
            "forgejo_bin_path": str(self._binary),
            "forgejo_config_local_path": str(self._template),
            "forgejo_service_path": str(self._root / "forgejo.service"),
            "forgejo_backup_dir": str(self._backups),
            "forgejo_health_url": f"http://127.0.0.1:{self._port}",
            "forgejo_health_retries": 5,
            "forgejo_health_delay": 1,
        }
        """Role inputs bound to disposable paths and the current unprivileged account."""

    def tearDown(self) -> None:
        """Stop only fixture processes and discard private runtime credentials."""
        pid_file: pathlib.Path = self._root / "pid"
        if pid_file.exists():
            try:
                os.killpg(int(pid_file.read_text()), signal.SIGKILL)
            except ProcessLookupError:
                pass
        shutil.rmtree(self._root)

    def _role(self, binary: str, success: bool = True) -> str:
        """Run the real packaged role without exposing its secret-bearing output."""
        self._variables["forgejo_bin_local_path"] = binary
        variables: pathlib.Path = self._root / "variables.json"
        variables.write_text(json.dumps(self._variables))
        process: subprocess.CompletedProcess[str] = subprocess.run(
            [
                _OPTIONS.ansible,
                "-i",
                "localhost,",
                "-e",
                "@" + str(variables),
                str(self._root / "playbook.yaml"),
            ],
            env=self._environment,
            text=True,
            capture_output=True,
            timeout=180,
            check=False,
        )
        if (process.returncode == 0) != success:
            # Fixtures contain generated credentials. Keep raw task output private
            # and return task names rather than exception or rendered config values.
            diagnostic: str = process.stdout[-7000:] + process.stderr[:1000]
            server_log: pathlib.Path = self._root / "server.log"
            if server_log.exists():
                diagnostic += (
                    "\nServer diagnostics:\n"
                    + "\n".join(
                        line
                        for line in server_log.read_text(
                            errors="replace"
                        ).splitlines()
                        if any(
                            marker in line.lower()
                            for marker in (
                                "[e]",
                                "[f]",
                                "error",
                                "fail",
                                "panic",
                                "fatal",
                            )
                        )
                    )[-4000:]
                )
            doctor_log: pathlib.Path = (
                pathlib.Path(
                    str(
                        self._variables.get(
                            "forgejo_log_dir", self._state / "logs"
                        )
                    )
                )
                / "upgrade-doctor.log"
            )
            if doctor_log.exists():
                messages: list[str] = [
                    line
                    for line in doctor_log.read_text(
                        errors="replace"
                    ).splitlines()
                    if any(
                        marker in line.lower()
                        for marker in ("[e]", "[w]", "error", "fail", "warn")
                    )
                ]
                diagnostic += (
                    "\nDoctor diagnostics:\n" + "\n".join(messages)[-5000:]
                )
            doctor_output: pathlib.Path = self._root / "doctor-output"
            if doctor_output.exists():
                diagnostic += (
                    "\nDoctor output:\n"
                    + doctor_output.read_text(errors="replace")[-5000:]
                )
            for line in self._template.read_text().splitlines():
                key, separator, value = line.partition(" = ")
                if separator and ("SECRET" in key or "TOKEN" in key):
                    diagnostic = diagnostic.replace(value, "<redacted>")
            if self._token:
                diagnostic = diagnostic.replace(self._token, "<redacted>")
            for value in self._private_values:
                diagnostic = diagnostic.replace(value, "<redacted>")
            self.fail(
                f"unexpected Ansible exit {process.returncode}: {diagnostic}"
            )
        if success:
            self.assertFalse((self._backups / "pending.json").exists())
            version_response: object = self._request("/version")
            self.assertIsInstance(version_response, dict)
            self._observed_versions.append(str(version_response["version"]))
        return process.stdout

    def _cli(self, *arguments: str) -> str:
        """Execute a local CLI operation and keep potential credentials private."""
        result: subprocess.CompletedProcess[str] = subprocess.run(
            [
                str(self._binary),
                "--config",
                str(self._config),
                "--work-path",
                str(self._state / "work"),
                *arguments,
            ],
            env=self._environment,
            text=True,
            capture_output=True,
            timeout=60,
            check=False,
        )
        self.assertEqual(
            result.returncode, 0, "local fixture CLI operation failed"
        )
        return result.stdout

    def _request(
        self, path: str, payload: dict[str, object] | None = None
    ) -> object:
        """Use the loopback fixture API without logging its access credential."""
        request: urllib.request.Request = urllib.request.Request(
            f"http://127.0.0.1:{self._port}/api/v1" + path,
            data=json.dumps(payload).encode() if payload is not None else None,
            headers={
                "Authorization": "token " + self._token,
                "Content-Type": "application/json",
            },
        )
        try:
            with urllib.request.urlopen(request, timeout=10) as response:
                return json.load(response)
        except urllib.error.HTTPError as error:
            error.close()
            raise

    def _database_versions(self) -> dict[str, int]:
        """Observe actual migration markers without reading user or credential rows."""
        with contextlib.closing(
            sqlite3.connect(self._state / "work/data/forgejo.db")
        ) as database:
            tables: list[tuple[str]] = database.execute(
                "SELECT name FROM sqlite_master WHERE type='table' AND name IN ('version', 'forgejo_version')"
            ).fetchall()
            schema: dict[str, int] = {
                table: int(
                    database.execute(
                        "SELECT version FROM " + table
                    ).fetchone()[0]
                )
                for (table,) in tables
            }
            self._observed_schemas.append(schema)
            return schema

    def _wrapper(
        self, name: str, source: str, version: str, rejected_command: str
    ) -> pathlib.Path:
        """Create a fixture executable that rejects only the requested CLI operation."""
        path: pathlib.Path = self._root / name
        path.write_text(
            f"#!/usr/bin/env python3\nimport os,sys\nif '--version' in sys.argv:\n print('forgejo version {version}')\nelif {rejected_command!r} and {rejected_command!r} in sys.argv:\n sys.exit(1)\nelse:\n os.execv({source!r}, [{source!r}]+sys.argv[1:])\n"
        )
        path.chmod(0o700)
        return path

    def _seed(self) -> None:
        """Create a real account and Git repository in the old SQLite instance."""
        self._cli(
            "admin",
            "user",
            "create",
            "--username",
            "fixture",
            "--email",
            "fixture@example.invalid",
            "--admin",
            "--random-password",
            "--must-change-password=false",
        )
        self._token = (
            self._cli(
                "admin",
                "user",
                "generate-access-token",
                "--username",
                "fixture",
                "--token-name",
                "fixture",
                "--scopes",
                "all",
                "--raw",
            )
            .splitlines()[-1]
            .strip()
        )
        self.assertTrue(
            len(self._token) == 40
            and all(
                character in "0123456789abcdef" for character in self._token
            ),
            "the fixture CLI did not return a raw access token",
        )
        self._request(
            "/user/repos",
            {
                "name": "preserved",
                "auto_init": True,
                "readme": "Default",
                "default_branch": "main",
            },
        )

    def _readme(self) -> bytes:
        """Read an actual Git blob through the service API for persistence acceptance."""
        content: object = self._request(
            "/repos/fixture/preserved/contents/README.md"
        )
        if not isinstance(content, dict) or not isinstance(
            content.get("content"), str
        ):
            self.fail("fixture repository did not return its committed README")
        return base64.b64decode(content["content"])

    def _actions(self) -> list[str]:
        """Read sanitized lifecycle operations emitted by the fixture manager."""
        return (
            [
                json.loads(line)["action"]
                for line in self._events.read_text().splitlines()
            ]
            if self._events.exists()
            else []
        )

    def _archives(self) -> list[pathlib.Path]:
        """Locate completed recovery archives, excluding staging directories."""
        return sorted(self._backups.glob("*/snapshot.tar.gz"))

    def _service(self, action: str) -> None:
        """Control only this fixture's service through its bounded systemctl shim."""
        subprocess.run(
            [
                str(self._root / "commands/systemctl"),
                action,
                "forgejo.service",
            ],
            env=self._environment,
            check=True,
            timeout=15,
        )

    def _restore_directory_metadata(self, metadata: object) -> None:
        """Apply the manifest's POSIX root metadata only inside this fixture."""
        if not isinstance(metadata, list):
            self.fail("recovery manifest directory metadata is not a list")
        for record in metadata:
            if not isinstance(record, dict):
                self.fail(
                    "recovery manifest contains invalid directory metadata"
                )
            path: object = record.get("path")
            uid: object = record.get("uid")
            gid: object = record.get("gid")
            mode: object = record.get("mode")
            if not (
                isinstance(path, str)
                and isinstance(uid, int)
                and isinstance(gid, int)
                and isinstance(mode, str)
            ):
                self.fail(
                    "recovery manifest contains invalid directory fields"
                )
            directory: pathlib.Path = pathlib.Path(path)
            self.assertTrue(
                directory.resolve().is_relative_to(self._root.resolve())
            )
            directory.mkdir(parents=True, exist_ok=True)
            os.chown(directory, uid, gid)
            directory.chmod(int(mode, 8))

    def _evidence(self, scenario: str) -> None:
        """Retain only scenario acceptance, versions, and lifecycle metadata."""
        output: pathlib.Path = pathlib.Path(
            os.environ["TEST_UNDECLARED_OUTPUTS_DIR"]
        )
        output.mkdir(parents=True, exist_ok=True)
        (output / (scenario + ".json")).write_text(
            json.dumps(
                {
                    "scenario": scenario,
                    "passed": True,
                    "actions": self._actions(),
                    "archive_count": len(self._archives()),
                    "observed_versions": self._observed_versions,
                    "observed_schemas": self._observed_schemas,
                },
                indent=2,
            )
            + "\n"
        )

    def test_upgrade_restore_and_unchanged(self) -> None:
        """Migrate actual old data, preserve a repository, and restore the old state."""
        external: dict[str, pathlib.Path] = {
            name: self._root / name
            for name in (
                "repositories",
                "lfs",
                "avatars",
                "repository-avatars",
                "logs",
                "git-home",
                "credentials",
                "empty-storage",
            )
        }
        for name, directory in external.items():
            directory.mkdir(mode=0o700 if name == "credentials" else 0o750)
            if name != "empty-storage":
                (directory / "retained-object").write_bytes(
                    b"preserved fixture object"
                )
        config_root: pathlib.Path = self._root / "configuration"
        self._config = config_root / "forgejo.ini"
        self._variables.update(
            {
                "forgejo_config_dir": str(config_root),
                "forgejo_config_path": str(self._config),
                "forgejo_log_dir": str(external["logs"]),
            }
        )
        configuration: str = self._template.read_text().replace(
            f"ROOT = {self._state}/work/data/forgejo-repositories",
            f'ROOT = "{external["repositories"]}"',
        ).replace(
            "[server]\n",
            f"[server]\nLFS_CONTENT_PATH = '{external['lfs']}'\n",
        ) + (
            f"[picture]\nAVATAR_UPLOAD_PATH = `{external['avatars']}`\n"
            f"REPOSITORY_AVATAR_UPLOAD_PATH = {external['repository-avatars']}\n"
        )
        configuration = (
            configuration.replace(
                "[log]\n", f"[log]\nROOT_PATH = {external['logs']}\n"
            )
            + f"[git]\nHOME_PATH = {external['git-home']}\n"
        )
        git_config: pathlib.Path = external["git-home"] / ".gitconfig"
        git_config.write_text("[fixture]\n\tretained = true\n")
        signing_state: pathlib.Path = external["git-home"] / ".gnupg"
        signing_state.mkdir(mode=0o700)
        signing_marker: pathlib.Path = signing_state / "fixture-state"
        signing_marker.write_bytes(b"synthetic signing state")
        token_file: pathlib.Path = external["credentials"] / "internal_token"
        token_line: str = next(
            line
            for line in configuration.splitlines()
            if line.startswith("INTERNAL_TOKEN = ")
        )
        token_value: str = token_line.partition(" = ")[2]
        token_file.write_text(token_value)
        token_file.chmod(0o600)
        self._private_values.append(token_value)
        configuration = configuration.replace(
            token_line,
            f"INTERNAL_TOKEN_URI = file:{token_file}",
        )
        self._template.write_text(configuration)
        backup_roots: list[str] = [
            str(self._state),
            str(self._binary),
            str(config_root),
            str(self._root / "forgejo.service"),
            *(str(path) for path in external.values()),
        ]
        self._variables["forgejo_backup_paths"] = backup_roots
        self._role(_OPTIONS.previous)
        self.assertFalse((self._state / 'work/"').exists())
        self._seed()
        previous_readme: bytes = self._readme()
        previous_schema: dict[str, int] = self._database_versions()
        self.assertTrue(previous_schema)
        signing_mode: int = signing_state.stat().st_mode & 0o7777
        original_metadata: dict[pathlib.Path, tuple[int, int, int]] = {}
        for root in backup_roots:
            directory: pathlib.Path = pathlib.Path(root)
            if directory.is_dir():
                info: os.stat_result = directory.stat()
                original_metadata[directory] = (
                    info.st_uid,
                    info.st_gid,
                    info.st_mode & 0o7777,
                )
        self._events.write_text("")
        for omitted in [
            str(self._binary),
            str(config_root),
            str(self._root / "forgejo.service"),
        ]:
            self._variables["forgejo_backup_paths"] = [
                root for root in backup_roots if root != omitted
            ]
            self._role(_OPTIONS.candidate, success=False)
            self.assertEqual(self._actions(), [])
            self.assertEqual(self._archives(), [])
        self._variables["forgejo_backup_paths"] = backup_roots
        self._role(_OPTIONS.candidate)
        self.assertTrue(self._database_versions())
        self.assertEqual(self._actions(), ["stop", "start"])
        self.assertIsInstance(self._request("/repos/fixture/preserved"), dict)
        self.assertEqual(self._readme(), previous_readme)
        archives: list[pathlib.Path] = self._archives()
        self.assertEqual(len(archives), 1)
        self.assertEqual(archives[0].stat().st_mode & 0o777, 0o600)
        self.assertEqual(archives[0].parent.stat().st_mode & 0o777, 0o700)
        with tarfile.open(archives[0]) as archive:
            names: list[str] = archive.getnames()
            self.assertTrue(any(name.endswith("forgejo.db") for name in names))
            self.assertTrue(any("preserved.git" in name for name in names))
        self._events.write_text("")
        self._binary.chmod(0o600)
        self._config.chmod(0o644)
        (self._root / "forgejo.service").chmod(0o600)
        self._role(_OPTIONS.candidate)
        self.assertEqual(self._actions(), [])
        self.assertEqual(self._archives(), archives)
        for path, mode in (
            (self._binary, 0o755),
            (self._config, 0o600),
            (self._root / "forgejo.service", 0o644),
        ):
            self.assertEqual(path.stat().st_mode & 0o7777, mode)
        self._role(_OPTIONS.previous, success=False)
        self.assertEqual(self._actions(), [])
        self.assertIsInstance(self._request("/repos/fixture/preserved"), dict)
        self._service("stop")
        manifest: dict[str, object] = json.loads(
            (archives[0].parent / "manifest.json").read_text()
        )
        shutil.rmtree(self._state)
        shutil.rmtree(config_root)
        for directory in external.values():
            shutil.rmtree(directory)
        self._binary.unlink()
        (self._root / "forgejo.service").unlink()
        with tarfile.open(archives[0]) as archive:
            # This archive comes only from validated test-private roots. Retain
            # directory modes, which the data extraction filter would omit.
            archive.extractall(str(manifest["archive_root"]), filter="tar")
        self._restore_directory_metadata(manifest.get("directory_roots", []))
        for directory, expected in original_metadata.items():
            info = directory.stat()
            self.assertEqual(
                (info.st_uid, info.st_gid, info.st_mode & 0o7777), expected
            )
        self.assertTrue(
            token_file.read_text() == token_value,
            "the restored internal-token file differs",
        )
        self.assertEqual(token_file.stat().st_mode & 0o7777, 0o600)
        self.assertIn("retained = true", git_config.read_text())
        self.assertEqual(
            signing_marker.read_bytes(), b"synthetic signing state"
        )
        self.assertEqual(signing_state.stat().st_mode & 0o7777, signing_mode)
        self._service("start")
        for _ in range(50):
            try:
                self.assertIsInstance(
                    self._request("/repos/fixture/preserved"), dict
                )
                break
            except (OSError, urllib.error.URLError):
                time.sleep(0.2)
        else:
            self.fail("restored prior service did not become ready")
        self.assertIn("15.0.3", self._cli("--version"))
        self.assertEqual(self._database_versions(), previous_schema)
        self.assertEqual(self._readme(), previous_readme)
        for name, directory in external.items():
            if name == "empty-storage":
                self.assertEqual(list(directory.iterdir()), [])
                continue
            self.assertEqual(
                (directory / "retained-object").read_bytes(),
                b"preserved fixture object",
            )
        self._request("/orgs", {"username": "restored-organization"})
        self.assertIsInstance(
            self._request(
                "/orgs/restored-organization/repos",
                {
                    "name": "restored",
                    "auto_init": True,
                    "readme": "Default",
                    "default_branch": "main",
                },
            ),
            dict,
        )
        self._evidence("upgrade-and-restore")

    def test_preflight_and_configuration_change(self) -> None:
        """Reject invalid declared backups without downtime and guard config changes."""
        self._role(_OPTIONS.candidate)
        self._events.write_text("")
        original: str = self._template.read_text()
        self._template.write_text(
            original.replace("DB_TYPE = sqlite3", "DB_TYPE = postgres")
        )
        self._role(_OPTIONS.candidate, success=False)
        self.assertEqual(self._actions(), [])
        self.assertEqual(self._archives(), [])
        self._template.write_text(
            original.replace(
                "APP_NAME = Fixture", "APP_NAME = Backup validation"
            )
        )
        self._variables["forgejo_backup_paths"] = [
            str(self._state),
            str(self._binary),
            str(self._root / "forgejo.service"),
            str(self._root / "missing-root"),
        ]
        self._role(_OPTIONS.candidate, success=False)
        self.assertEqual(self._actions(), [])
        self.assertEqual(self._archives(), [])
        self._variables.pop("forgejo_backup_paths")
        self._template.write_text(original)
        recovery_alias: pathlib.Path = self._root / "recovery-alias"
        recovery_alias.symlink_to(self._state, target_is_directory=True)
        self._variables["forgejo_backup_dir"] = str(
            recovery_alias / "recovery"
        )
        self._role(_OPTIONS.candidate, success=False)
        self.assertEqual(self._actions(), [])
        self.assertEqual(self._archives(), [])
        self.assertFalse((self._state / "recovery").exists())
        archive_storage: pathlib.Path = self._root / "archive-storage"
        (archive_storage / "input").mkdir(parents=True, mode=0o700)
        archive_alias: pathlib.Path = self._root / "archive-alias"
        archive_alias.symlink_to(archive_storage, target_is_directory=True)
        self._variables["forgejo_backup_dir"] = str(archive_storage)
        self._variables["forgejo_backup_paths"] = [
            str(self._state),
            str(self._binary),
            str(self._root / "forgejo.service"),
            str(archive_alias / "input"),
        ]
        self._role(_OPTIONS.candidate, success=False)
        self.assertEqual(self._actions(), [])
        self.assertEqual(self._archives(), [])
        self._variables.pop("forgejo_backup_paths")
        self._variables["forgejo_backup_dir"] = str(self._backups)
        self._template.write_text(original)
        self._variables["forgejo_database_path"] = str(
            self._state / "missing.db"
        )
        self._role(_OPTIONS.candidate, success=False)
        self.assertEqual(self._actions(), [])
        self._variables.pop("forgejo_database_path")
        self._template.write_text(original)
        for destination in (
            str(self._state / "recursive-backup"),
            "/" + str(self._state / "recursive-backup"),
        ):
            self._variables["forgejo_backup_dir"] = destination
            self._role(_OPTIONS.candidate, success=False)
            self.assertEqual(self._actions(), [])
        self._variables["forgejo_backup_dir"] = str(self._backups)
        self._template.write_text(
            original.replace("SECRET_KEY = ", "SECRET_KEY = changed-")
        )
        self._role(_OPTIONS.candidate, success=False)
        self.assertEqual(self._actions(), [])
        backup_volume: pathlib.Path = self._root / "backup-volume"
        backup_volume.mkdir(mode=0o700)
        backup_alias: pathlib.Path = self._root / "backup-volume-alias"
        backup_alias.symlink_to(backup_volume, target_is_directory=True)
        self._backups = backup_volume / "recovery"
        self._variables["forgejo_backup_dir"] = str(backup_alias / "recovery")
        self._template.write_text(
            original.replace(
                "APP_NAME = Fixture", "APP_NAME = Updated fixture"
            )
            + "[avatar]\nSTORAGE_TYPE = local\n"
        )
        self._role(_OPTIONS.candidate)
        self.assertEqual(self._actions(), ["stop", "start"])
        self.assertEqual(len(self._archives()), 1)
        with urllib.request.urlopen(
            f"http://127.0.0.1:{self._port}/", timeout=10
        ) as response:
            self.assertTrue(
                "Updated fixture" in response.read().decode(),
                "the configured application name was not served",
            )
        self._events.write_text("")
        self._role(_OPTIONS.candidate)
        self.assertEqual(self._actions(), [])
        self.assertEqual(len(self._archives()), 1)
        self._evidence("configuration-transaction")

    def test_rejected_candidate_stays_stopped(self) -> None:
        """Keep a backup and stop a candidate whose running version is unexpected."""
        self._role(_OPTIONS.previous)
        self._seed()
        self._events.write_text("")
        rejected: pathlib.Path = self._wrapper(
            "rejected", _OPTIONS.candidate, "18.0.0", ""
        )
        self._role(str(rejected), success=False)
        self.assertEqual(self._actions(), ["stop", "start", "stop"])
        self.assertEqual(len(self._archives()), 1)
        with self.assertRaises((OSError, urllib.error.URLError)):
            self._request("/version")
        self._evidence("rejected-candidate")

    def test_drain_failure_preserves_running_instance(self) -> None:
        """Reject an upgrade before stopping when the installed manager cannot drain."""
        self._role(_OPTIONS.previous)
        self._seed()
        self._wrapper(
            "installed-wrapper", _OPTIONS.previous, "15.0.3", "flush-queues"
        ).replace(self._binary)
        self._events.write_text("")
        self._role(_OPTIONS.candidate, success=False)
        self.assertEqual(self._actions(), [])
        self.assertEqual(self._archives(), [])
        self.assertIsInstance(self._request("/repos/fixture/preserved"), dict)
        self._evidence("drain-failure")

    def test_backup_failure_recovers_prior_process(self) -> None:
        """Do not install a candidate when recovery storage disappears after stop."""
        self._role(_OPTIONS.previous)
        self._seed()
        (self._root / "fail-backup").touch()
        self._events.write_text("")
        self._role(_OPTIONS.candidate, success=False)
        self.assertEqual(self._actions(), ["stop", "start"])
        self.assertEqual(self._archives(), [])
        self.assertIn("15.0.3", self._cli("--version"))
        for _ in range(50):
            try:
                self.assertIsInstance(
                    self._request("/repos/fixture/preserved"), dict
                )
                break
            except (OSError, urllib.error.URLError):
                time.sleep(0.2)
        else:
            self.fail("prior instance did not recover after backup failure")
        self._evidence("backup-failure")

    def test_marker_write_failure_recovers_prior_process(self) -> None:
        """Recover intact old inputs when the activation marker cannot be written."""
        self._role(_OPTIONS.previous)
        self._seed()
        previous_readme: bytes = self._readme()
        startup_files: list[pathlib.Path] = [
            self._binary,
            self._config,
            self._root / "forgejo.service",
        ]
        original_contents: list[bytes] = [
            path.read_bytes() for path in startup_files
        ]
        (self._root / "fail-marker").touch()
        self._events.write_text("")
        output: str = self._role(_OPTIONS.candidate, success=False)
        self.assertIn("can not use content with a dir as dest", output)
        self.assertEqual(self._actions(), ["stop", "start"])
        self.assertEqual(len(self._archives()), 1)
        self.assertTrue((self._backups / "pending.json").is_dir())
        self.assertTrue(
            all(
                path.read_bytes() == contents
                for path, contents in zip(startup_files, original_contents)
            ),
            "startup inputs changed before the marker was written",
        )
        for _ in range(50):
            try:
                version_response: object = self._request("/version")
                self.assertIsInstance(version_response, dict)
                self.assertEqual(
                    str(version_response["version"]).split("+", 1)[0],
                    "15.0.3",
                )
                self.assertEqual(self._readme(), previous_readme)
                break
            except (OSError, urllib.error.URLError):
                time.sleep(0.2)
        else:
            self.fail("prior instance did not recover after marker failure")
        self._evidence("marker-write-failure")

    def test_doctor_failure_stops_candidate(self) -> None:
        """Keep matching recovery data when a migrated candidate fails doctor checks."""
        self._role(_OPTIONS.previous)
        self._seed()
        self._events.write_text("")
        rejected: pathlib.Path = self._wrapper(
            "doctor-rejected", _OPTIONS.candidate, "17.0.1", "doctor"
        )
        self._role(str(rejected), success=False)
        self.assertEqual(self._actions(), ["stop", "start", "stop"])
        self.assertEqual(len(self._archives()), 1)
        self.assertTrue((self._backups / "pending.json").exists())
        self._role(str(rejected), success=False)
        self.assertEqual(self._actions(), ["stop", "start", "stop"])
        self.assertEqual(len(self._archives()), 1)
        with self.assertRaises((OSError, urllib.error.URLError)):
            self._request("/version")
        self._evidence("doctor-failure")


if __name__ == "__main__":
    _parser: argparse.ArgumentParser = argparse.ArgumentParser()
    """Command-line resource selector for the Bazel test entry point."""
    for _option in (
        "ansible",
        "candidate",
        "previous",
        "service-fixture",
        "collections",
    ):
        _parser.add_argument("--" + _option, required=True)
    _remaining: list[str]
    """Arguments passed through to unittest after resource resolution."""
    _OPTIONS, _remaining = _parser.parse_known_args()
    _resolver: runfiles.Runfiles | None = runfiles.Create()
    """Resolver for immutable files declared by the test target."""
    assert _resolver is not None
    for _option in (
        "ansible",
        "candidate",
        "previous",
        "service_fixture",
        "collections",
    ):
        setattr(
            _OPTIONS, _option, _resolver.Rlocation(getattr(_OPTIONS, _option))
        )
    unittest.main(argv=[sys.argv[0], *_remaining])

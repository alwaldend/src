"""HTTP acceptance checks; failure coverage is recorded before app implementation."""

import argparse
import contextlib
import hashlib
import http.client
import json
import os
import pathlib
import socket
import subprocess
import tempfile
import time
import typing
import urllib.error
import urllib.parse
import urllib.request
from collections import abc as collections_abc

from python.runfiles import runfiles


def resolve(path: str) -> pathlib.Path:
    resolver: runfiles.Runfiles | None = runfiles.Create()
    assert resolver is not None, "runfiles are unavailable"
    resolved: str | None = resolver.Rlocation(path)
    assert resolved is not None and pathlib.Path(resolved).is_file(), path
    return pathlib.Path(resolved)


def mapping(value: object) -> dict[str, object]:
    assert isinstance(value, dict), value
    return typing.cast(dict[str, object], value)


def request(port: int, path: str) -> tuple[int, str, bytes]:
    response: http.client.HTTPResponse
    error: urllib.error.HTTPError
    try:
        with urllib.request.urlopen(
            f"http://127.0.0.1:{port}{path}", timeout=2
        ) as response:
            return (
                response.status,
                response.headers["Content-Type"],
                response.read(),
            )
    except urllib.error.HTTPError as error:
        return error.code, error.headers["Content-Type"], error.read()


@contextlib.contextmanager
def running(
    executable: pathlib.Path, directory: pathlib.Path
) -> collections_abc.Iterator[int]:
    listener: socket.socket
    with socket.socket() as listener:
        listener.bind(("127.0.0.1", 0))
        port: int = listener.getsockname()[1]
    log: typing.BinaryIO
    with directory.joinpath(executable.name + ".log").open("wb") as log:
        process: subprocess.Popen[bytes] = subprocess.Popen(
            [str(executable), "--host", "127.0.0.1", "--port", str(port)],
            cwd=directory,
            stdout=log,
            stderr=subprocess.STDOUT,
        )
        try:
            deadline: float = time.monotonic() + 30
            while True:
                assert process.poll() is None, "server exited before readiness"
                try:
                    if request(port, "/health")[0] == 200:
                        break
                except (OSError, TimeoutError):
                    pass
                assert time.monotonic() < deadline, "readiness timed out"
                time.sleep(0.05)
            yield port
        finally:
            process.terminate()
            try:
                process.wait(timeout=10)
            except subprocess.TimeoutExpired:
                process.kill()
                process.wait(timeout=5)
                raise AssertionError(
                    "server failed graceful shutdown"
                ) from None
            assert process.returncode in (0, -15), process.returncode
    with socket.socket() as listener:
        listener.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        listener.bind(("127.0.0.1", port))


def main() -> None:
    parser: argparse.ArgumentParser = argparse.ArgumentParser()
    parser.add_argument(
        "--recommender",
        default="_main/projects/mlops_2026/cmd/recommender/recommender",
    )
    parser.add_argument(
        "--isolation", default="_main/projects/mlops_2026/test/e2e/isolation"
    )
    args: argparse.Namespace = parser.parse_args()
    executable: pathlib.Path = resolve(args.recommender)
    isolation: pathlib.Path = resolve(args.isolation)
    results: dict[str, object] = {}
    temporary: str
    port: int
    occupied: socket.socket
    value: object
    with tempfile.TemporaryDirectory(
        dir=os.environ.get("TEST_TMPDIR")
    ) as temporary:
        directory: pathlib.Path = pathlib.Path(temporary)
        with running(executable, directory) as port:
            status: int
            content_type: str
            body: bytes
            status, content_type, body = request(port, "/health")
            assert status == 200 and "application/json" in content_type
            assert json.loads(body) == {"status": "ok"}
            results["health"] = "passed"
            status, content_type, body = request(port, "/openapi.json")
            assert status == 200 and "application/json" in content_type
            schema: dict[str, object] = mapping(json.loads(body))
            assert str(schema["openapi"]).startswith("3.")
            paths: dict[str, object] = mapping(schema["paths"])
            assert set(paths) == {"/health", "/api/v1/recommend"}
            response: dict[str, object] = mapping(
                mapping(mapping(paths["/health"])["get"])["responses"]
            )
            assert "200" in response
            models: dict[str, object] = mapping(
                mapping(schema["components"])["schemas"]
            )
            health: dict[str, object] = mapping(models["HealthResponse"])
            assert "status" in mapping(health["properties"])
            assert "RecommenderRecommendationResponse" in models
            assert "RecommenderSong" in models
            results["openapi"] = "passed"
            query: str = urllib.parse.urlencode({"track_name": "Time", "n": 3})
            status, content_type, body = request(
                port, "/api/v1/recommend?" + query
            )
            assert status == 200
            recommendation_response: dict[str, object] = mapping(
                json.loads(body)
            )
            assert recommendation_response["requested_track"] == "Time"
            songs: object = recommendation_response["recommendations"]
            assert isinstance(songs, list) and len(songs) == 3
            song: object
            for song in songs:
                assert set(mapping(song)) == {
                    "track_name",
                    "artists",
                    "album_name",
                }
                for value in mapping(song).values():
                    assert isinstance(value, str)
            invalid: str
            for invalid in ("0", "21", "invalid"):
                query = urllib.parse.urlencode(
                    {"track_name": "Time", "n": invalid}
                )
                assert request(port, "/api/v1/recommend?" + query)[0] == 422
            assert request(port, "/api/v1/recommend")[0] == 422
            query = urllib.parse.urlencode(
                {"track_name": "__absent_course_track__"}
            )
            status, content_type, body = request(
                port, "/api/v1/recommend?" + query
            )
            assert status == 200
            assert json.loads(body) == {
                "requested_track": "__absent_course_track__",
                "recommendations": [],
            }
            query = urllib.parse.urlencode({"track_name": "Time"})
            status, content_type, body = request(
                port, "/api/v1/recommend?" + query
            )
            assert status == 200
            songs = mapping(json.loads(body))["recommendations"]
            assert isinstance(songs, list) and len(songs) == 5
            query = urllib.parse.urlencode({"track_name": "Time", "n": 20})
            status, content_type, body = request(
                port, "/api/v1/recommend?" + query
            )
            assert status == 200
            songs = mapping(json.loads(body))["recommendations"]
            assert isinstance(songs, list) and len(songs) == 20
            results["recommendations"] = "passed"
            path: str
            for path in ("/docs", "/redoc"):
                status, content_type, body = request(port, path)
                assert status == 200 and "text/html" in content_type
                assert b"/openapi.json" in body
            assert request(port, "/api/v1/predict")[0] == 404
            results["docs_and_version_boundary"] = "passed"
        with running(isolation, directory) as port:
            assert request(port, "/first/only-first")[0] == 200
            assert request(port, "/second/only-first")[0] == 404
            assert request(port, "/second/health")[0] == 200
            results["factory_isolation"] = "passed"
        with socket.socket() as occupied:
            occupied.bind(("127.0.0.1", 0))
            occupied.listen()
            blocked_port: int = occupied.getsockname()[1]
            completed: subprocess.CompletedProcess[bytes] = subprocess.run(
                [str(executable), "--port", str(blocked_port)],
                cwd=directory,
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT,
                timeout=15,
            )
            assert completed.returncode != 0, "port conflict must fail startup"
            results["port_conflict"] = "passed"
        runfiles_root: pathlib.Path = pathlib.Path(os.environ["TEST_SRCDIR"])
        without_model: pathlib.Path = directory / "without-model.manifest"
        manifest_lines: list[str] = []
        child: pathlib.Path
        for child in runfiles_root.iterdir():
            manifest_lines.append(f"{child.name} {child}\n")
        manifest_lines.append(
            f"_main/projects/mlops_2026/model.pkl {directory / 'absent.pkl'}\n"
        )
        without_model.write_text("".join(manifest_lines))
        environment: dict[str, str] = dict(os.environ)
        environment["RUNFILES_DIR"] = str(runfiles_root)
        environment["RUNFILES_MANIFEST_FILE"] = str(without_model)
        environment.pop("RUNFILES_MANIFEST_ONLY", None)
        completed = subprocess.run(
            [str(executable), "--port", "0"],
            cwd=directory,
            env=environment,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            timeout=15,
        )
        assert completed.returncode != 0
        assert (
            b"Recommendation model runfile is unavailable" in completed.stdout
        ), completed.stdout.decode(errors="replace")[-4000:]
        results["missing_model"] = "passed"
        name: str
        for name in ("data.csv", "dataset.csv"):
            manifest_lines[-1] = (
                f"+_repo_rules+com_github_yandex_practicum_mlops_freetrack/{name} "
                f"{directory / 'absent.csv'}\n"
            )
            without_model.write_text("".join(manifest_lines))
            completed = subprocess.run(
                [str(executable), "--port", "0"],
                cwd=directory,
                env=environment,
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT,
                timeout=15,
            )
            assert completed.returncode != 0
            assert (
                f"Course runfile is unavailable: {name}".encode()
                in completed.stdout
            ), completed.stdout.decode(errors="replace")[-4000:]
            results[f"missing_{name}"] = "passed"
        corrupt_model: pathlib.Path = directory / "corrupt.pkl"
        corrupt_model.write_bytes(b"\x00invalid model")
        manifest_lines[-1] = (
            f"_main/projects/mlops_2026/model.pkl {corrupt_model}\n"
        )
        without_model.write_text("".join(manifest_lines))
        completed = subprocess.run(
            [str(executable), "--port", "0"],
            cwd=directory,
            env=environment,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            timeout=15,
        )
        assert completed.returncode != 0
        assert b"Application startup failed" in completed.stdout
        results["corrupt_model"] = "passed"
    expected: dict[str, str] = {
        "dataset.csv": "b202fa49909b2d5cef71a04b1d21243cfeb36414535f2ca9272aa646721177bd",
        "data.csv": "c97c42349c2e97339ff908ab01163ec5a3264f02ec65e78cb961c85077ded7f8",
    }
    name: str
    digest: str
    for name, digest in expected.items():
        dataset: pathlib.Path = resolve(
            "com_github_yandex_practicum_mlops_freetrack/" + name
        )
        assert hashlib.sha256(dataset.read_bytes()).hexdigest() == digest
        results[name] = {"sha256": digest, "bytes": dataset.stat().st_size}
    output: pathlib.Path = pathlib.Path(
        os.environ["TEST_UNDECLARED_OUTPUTS_DIR"]
    )
    output.joinpath("recommender-e2e.json").write_text(
        json.dumps(results, indent=2) + "\n"
    )


if __name__ == "__main__":
    main()

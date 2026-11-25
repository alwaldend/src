"""Run the real job server against local Vast.ai, SSH, and tracking fixtures."""

import collections.abc
import concurrent.futures
import contextlib
import http.server
import json
import os
import pathlib
import shlex
import socket
import sqlite3
import subprocess
import sys
import tempfile
import threading
import time
import typing

import grpc
import vastai_job_pb2 as pb
import vastai_job_pb2_grpc as rpc


class Endpoint(http.server.BaseHTTPRequestHandler):
    """Serve fixture rental discovery and tracking requests."""

    states: typing.ClassVar[dict[str, str]] = {}
    update_attempts: typing.ClassVar[dict[str, int]] = {}
    rental_status: typing.ClassVar[str | None] = "running"
    _tracking_auth: typing.ClassVar[str] = "Bearer fixture-token"
    """Authorization required by the tracking routes."""

    def log_message(self, format: str, *args: object) -> None:
        """Keep the fixture's request logs quiet."""

    def _send(self, payload: dict[str, object]) -> None:
        """Return one successful JSON response."""
        self.send_response(200)
        self.end_headers()
        self.wfile.write(json.dumps(payload).encode())

    def do_GET(self) -> None:
        """Return a rental or the run's current tracked state."""
        if self.path == "/api/v0/instances/12/":
            if self.rental_status is None:
                self.send_error(404)
                return
            self._send(
                {
                    "instances": {
                        "id": 12,
                        "label": "gpu",
                        "actual_status": self.rental_status,
                        "ssh_host": "worker.example",
                        "ssh_port": 10022,
                    }
                }
            )
        elif self.path.startswith("/api/v1/instances/?"):
            self._send({"instances": [{"id": 12, "label": "gpu"}]})
        elif self.path.startswith("/api/2.0/mlflow/runs/get?run_id="):
            if self.headers.get("Authorization") != self._tracking_auth:
                self.send_error(401)
                return
            run_id: str = self.path.split("=", 1)[1]
            self._send(
                {
                    "run": {
                        "info": {"status": self.states.get(run_id, "RUNNING")}
                    }
                }
            )
        else:
            self.send_error(404)

    def do_POST(self) -> None:
        """Record a terminal tracking update and leave other routes closed."""
        assert self.path == "/api/2.0/mlflow/runs/update"
        if self.headers.get("Authorization") != self._tracking_auth:
            self.send_error(401)
            return
        record: dict[str, str] = json.loads(
            self.rfile.read(int(self.headers["Content-Length"]))
        )
        run_id: str = record["run_id"]
        self.update_attempts[run_id] = self.update_attempts.get(run_id, 0) + 1
        if run_id == "lost" and self.update_attempts[run_id] == 1:
            self.send_error(503)
            return
        self.states[run_id] = record["status"]
        self._send({})


def _ssh(root: pathlib.Path) -> None:
    """Simulate a worker with durable outcomes and one lost submit reply."""
    assert not {
        "VAST_API_KEY",
        "MLFLOW_TRACKING_TOKEN",
        "MLFLOW_TRACKING_TOKEN_FILE",
    }.intersection(os.environ)
    operation: str = sys.argv[-1].split()[-1]
    assert sys.argv[-2] == "worker.example"
    message: bytes = sys.stdin.buffer.read()
    reference: pb.JobReference = pb.JobReference.FromString(message)
    path: pathlib.Path = root / reference.job_id
    first: bool = not path.exists()
    if operation == "submit":
        with (root / f"{reference.job_id}.submissions").open(
            "a", encoding="utf-8"
        ) as requests:
            requests.write("submit\n")
    if operation == "submit" and first:
        path.write_text(str(pb.JOB_STATE_RUNNING), encoding="utf-8")
        if reference.job_id == "lost":
            raise SystemExit(255)
    if operation == "cancel":
        path.write_text(str(pb.JOB_STATE_KILLED), encoding="utf-8")
    result: pb.Job = pb.Job(
        job_id=reference.job_id,
        state=int(path.read_text(encoding="utf-8")),
    )
    if result.state == pb.JOB_STATE_FINISHED:
        result.exit_code = 0
    sys.stdout.buffer.write(result.SerializeToString())


@contextlib.contextmanager
def _server(
    args: list[str], address: str, environment: dict[str, str]
) -> collections.abc.Iterator[rpc.JobServiceStub]:
    """Start the job server and join it after the fixture finishes."""
    with subprocess.Popen(args, env=environment) as process:
        try:
            with grpc.insecure_channel(address) as channel:
                grpc.channel_ready_future(channel).result(timeout=15)
                yield rpc.JobServiceStub(channel)
        finally:
            process.terminate()
            assert process.wait(timeout=15) == 0


def _exercise(cli: pathlib.Path, root: pathlib.Path, endpoint: str) -> None:
    """Verify durable replay, tracking retry, cancellation, and rental loss."""
    with socket.socket() as sock:
        sock.bind(("127.0.0.1", 0))
        address: str = f"127.0.0.1:{sock.getsockname()[1]}"
    worker: pathlib.Path = root / "ssh"
    worker.write_text(
        "#!/bin/sh\nexec "
        + shlex.join(
            [
                sys.executable,
                str(pathlib.Path(__file__).resolve()),
                "--ssh",
                str(root),
            ]
        )
        + ' "$@"\n',
        encoding="utf-8",
    )
    worker.chmod(0o700)
    env: dict[str, str] = os.environ.copy()
    env["PATH"] = f"{root}{os.pathsep}{env['PATH']}"
    env["VAST_API_KEY"] = "fixture"
    env.pop("MLFLOW_TRACKING_TOKEN", None)
    token_file: pathlib.Path = root / "tracking-token"
    token_file.write_text("fixture-token\n", encoding="utf-8")
    env["MLFLOW_TRACKING_TOKEN_FILE"] = str(token_file)
    args: list[str] = [
        str(cli),
        "server",
        "serve",
        "--listen",
        address,
        "--state-dir",
        str(root / "state"),
        "--workers",
        "gpu",
        "--vast-endpoint",
        f"{endpoint}/api/v0",
        "--tracking-uri",
        endpoint,
    ]
    request: pb.SubmitJobRequest = pb.SubmitJobRequest(
        job_id="lost",
        worker="gpu",
        workload=pb.Workload(
            project_artifact_path="project.tar.gz",
            command="exit 0",
            tracking_uri=endpoint,
            experiment_id="0",
        ),
    )
    with _server(args, address, env) as service:
        try:
            service.SubmitJob(request, timeout=10)
        except grpc.RpcError as error:
            assert error.code() == grpc.StatusCode.UNAVAILABLE
        else:
            raise AssertionError("first SSH reply was lost")
        request.worker = "unrelated"
        try:
            service.SubmitJob(request, timeout=10)
        except grpc.RpcError as error:
            assert error.code() == grpc.StatusCode.ALREADY_EXISTS
        else:
            raise AssertionError("conflicting replay was accepted")
        request.worker = "gpu"
    with _server(args, address, env) as service:
        assert (
            service.SubmitJob(request, timeout=10).state
            == pb.JOB_STATE_RUNNING
        )
        (root / "lost").write_text(
            str(pb.JOB_STATE_FINISHED), encoding="utf-8"
        )
        deadline: float = time.monotonic() + 12
        while Endpoint.states.get("lost") != "FINISHED":
            assert time.monotonic() < deadline
            time.sleep(0.1)
        assert Endpoint.update_attempts["lost"] == 2
        job: pb.Job = service.GetJob(
            pb.JobReference(job_id="lost"), timeout=10
        )
        assert job.state == pb.JOB_STATE_FINISHED and job.exit_code == 0
        token_file.write_text("rotated-token\n", encoding="utf-8")
        Endpoint._tracking_auth = "Bearer rotated-token"
        request.job_id = "cancel"
        service.SubmitJob(request, timeout=10)
        assert (
            service.CancelJob(
                pb.JobReference(job_id="cancel"), timeout=10
            ).state
            == pb.JOB_STATE_KILLED
        )
        deadline = time.monotonic() + 12
        while Endpoint.states.get("cancel") != "KILLED":
            assert time.monotonic() < deadline
            time.sleep(0.1)
        request.job_id = "parallel"
        with concurrent.futures.ThreadPoolExecutor(max_workers=2) as pool:
            left: concurrent.futures.Future[pb.Job] = pool.submit(
                service.SubmitJob, request, timeout=10
            )
            right: concurrent.futures.Future[pb.Job] = pool.submit(
                service.SubmitJob, request, timeout=10
            )
            assert left.result() == right.result()
        assert (root / "parallel.submissions").read_text(
            encoding="utf-8"
        ) == "submit\n"
        parallel: pb.JobReference = pb.JobReference(job_id="parallel")
        Endpoint.rental_status = "stopped"
        unavailable: pb.Job = service.GetJob(parallel, timeout=10)
        assert (
            unavailable.state == pb.JOB_STATE_RUNNING
            and unavailable.observation_error
        )
        Endpoint.rental_status = "running"
        assert not service.GetJob(parallel, timeout=10).observation_error
        Endpoint.rental_status = None
        assert (
            service.GetJob(parallel, timeout=10).state == pb.JOB_STATE_UNKNOWN
        )

    artifact: pathlib.Path = pathlib.Path(
        os.environ["TEST_UNDECLARED_OUTPUTS_DIR"]
    )
    with (
        contextlib.closing(
            sqlite3.connect(root / "state/jobs.sqlite3")
        ) as source,
        contextlib.closing(
            sqlite3.connect(artifact / "coordinator.sqlite3")
        ) as destination,
    ):
        source.backup(destination)


def _main() -> None:
    """Run the server fixture or answer one simulated SSH invocation."""
    if sys.argv[1] == "--ssh":
        _ssh(pathlib.Path(sys.argv[2]))
        return
    with (
        tempfile.TemporaryDirectory(prefix="mlflow-server-") as directory,
        http.server.ThreadingHTTPServer(("127.0.0.1", 0), Endpoint) as server,
    ):
        thread: threading.Thread = threading.Thread(
            target=server.serve_forever
        )
        thread.start()
        try:
            _exercise(
                pathlib.Path(sys.argv[1]).absolute(),
                pathlib.Path(directory),
                f"http://127.0.0.1:{server.server_port}",
            )
        finally:
            server.shutdown()
            thread.join()


if __name__ == "__main__":
    _main()

"""Exercise the Linux worker through its CLI, gRPC socket, and MLflow artifacts."""

import collections.abc
import concurrent.futures
import contextlib
import os
import pathlib
import socket
import sqlite3
import subprocess
import sys
import tarfile
import tempfile
import time

import mlflow.tracking
import vastai_job_pb2


def _wait_until(
    predicate: collections.abc.Callable[[], bool], timeout: float = 35
) -> None:
    """Wait for an externally observable test event with a bounded deadline."""
    deadline: float = time.monotonic() + timeout
    while time.monotonic() < deadline:
        if predicate():
            return
        time.sleep(0.1)
    raise AssertionError("worker event did not occur before the deadline")


def _process_exited(pid: int) -> bool:
    """Report whether the background child exited, including before reaping."""
    try:
        status: str = pathlib.Path(f"/proc/{pid}/stat").read_text(
            encoding="utf-8"
        )
    except FileNotFoundError:
        return True
    return status.rpartition(") ")[2].startswith(("Z ", "X "))


def _socket_ready(path: pathlib.Path) -> bool:
    """Probe the active listener, including after a crash leaves an old inode."""
    with socket.socket(socket.AF_UNIX) as client:
        try:
            client.connect(str(path))
        except (FileNotFoundError, ConnectionRefusedError):
            return False
    return True


@contextlib.contextmanager
def _listener(
    argv: list[str], state: pathlib.Path, stderr: int | None = None
) -> collections.abc.Iterator[subprocess.Popen[bytes]]:
    """Start the service and clean up its listener at the end of the scenario."""
    with subprocess.Popen([*argv, "serve"], stderr=stderr) as service:
        try:
            _wait_until(lambda: _socket_ready(state / "worker.sock"))
            yield service
        finally:
            if service.poll() is None:
                service.terminate()
                assert service.wait(timeout=10) == 0


def _request(
    client: mlflow.tracking.MlflowClient,
    uri: str,
    experiment: str,
    archive: pathlib.Path,
    command: str,
) -> vastai_job_pb2.SubmitJobRequest:
    """Create the run and stage its project through the real tracking client."""
    run_id: str = client.create_run(experiment).info.run_id
    client.log_artifact(run_id, str(archive), artifact_path="vastai-source")
    return vastai_job_pb2.SubmitJobRequest(
        job_id=run_id,
        workload=vastai_job_pb2.Workload(
            project_artifact_path="vastai-source/project.tar.gz",
            command=command,
            tracking_uri=uri,
            experiment_id=experiment,
        ),
    )


def _call(
    argv: list[str],
    operation: str,
    request: vastai_job_pb2.SubmitJobRequest | vastai_job_pb2.JobReference,
) -> vastai_job_pb2.Job:
    """Send a protobuf command over the same streams used by SSH."""
    response: subprocess.CompletedProcess[bytes] = subprocess.run(
        [*argv, operation],
        input=request.SerializeToString(),
        stdout=subprocess.PIPE,
        check=True,
        timeout=20,
    )
    return vastai_job_pb2.Job.FromString(response.stdout)


def _wait_for_state(
    argv: list[str], reference: vastai_job_pb2.JobReference, expected: int
) -> vastai_job_pb2.Job:
    """Wait for a job state and fail promptly on an unexpected terminal result."""
    deadline: float = time.monotonic() + 40
    while time.monotonic() < deadline:
        job: vastai_job_pb2.Job = _call(argv, "status", reference)
        if job.state == expected:
            return job
        assert job.state == vastai_job_pb2.JOB_STATE_RUNNING, job
        time.sleep(0.15)
    raise AssertionError(f"job {reference.job_id} did not reach {expected}")


def _exercise(root: pathlib.Path, cli: pathlib.Path) -> None:
    """Check replay, completion, cancellation, and listener recovery."""
    state: pathlib.Path = root / "worker"
    argv: list[str] = [str(cli), "worker", "--state-dir", str(state)]
    uri: str = f"sqlite:///{root / 'tracking.sqlite3'}"
    client: mlflow.tracking.MlflowClient = mlflow.tracking.MlflowClient(
        tracking_uri=uri
    )
    experiment: str = client.create_experiment(
        "worker", artifact_location=(root / "artifacts").as_uri()
    )
    source: pathlib.Path = root / "input.txt"
    source.write_text("project input\n", encoding="utf-8")
    archive: pathlib.Path = root / "project.tar.gz"
    with tarfile.open(archive, "w:gz") as bundle:
        bundle.add(source, arcname="inputs/input.txt")

    submission: vastai_job_pb2.SubmitJobRequest = _request(
        client,
        uri,
        experiment,
        archive,
        "printf 'started\\n' >> attempts; "
        "while [ ! -f release ]; do sleep 0.1; done; "
        "cat inputs/input.txt; "
        'printf "%s\\n%s\\n%s\\n" "$MLFLOW_RUN_ID" '
        '"$MLFLOW_EXPERIMENT_ID" "$MLFLOW_TRACKING_URI"; '
        "printf 'training stderr\\n' >&2",
    )
    reference: vastai_job_pb2.JobReference = vastai_job_pb2.JobReference(
        job_id=submission.job_id
    )
    job_files: pathlib.Path = state / "jobs" / submission.job_id
    project: pathlib.Path = job_files / "project"
    worker_log: pathlib.Path = root / "worker.log"
    with (
        worker_log.open("wb") as log_file,
        _listener(argv, state, stderr=log_file.fileno()),
    ):
        with concurrent.futures.ThreadPoolExecutor(max_workers=2) as pool:
            first: concurrent.futures.Future[vastai_job_pb2.Job] = pool.submit(
                _call, argv, "submit", submission
            )
            second: concurrent.futures.Future[vastai_job_pb2.Job] = (
                pool.submit(_call, argv, "submit", submission)
            )
            for response in (first.result(), second.result()):
                assert response.job_id == reference.job_id
                assert response.state == vastai_job_pb2.JOB_STATE_RUNNING
        _wait_until((project / "attempts").is_file)
        assert (
            _call(argv, "status", reference).state
            == vastai_job_pb2.JOB_STATE_RUNNING
        )
        conflict: vastai_job_pb2.SubmitJobRequest = (
            vastai_job_pb2.SubmitJobRequest()
        )
        conflict.CopyFrom(submission)
        conflict.workload.command = "exit 99"
        rejected: subprocess.CompletedProcess[bytes] = subprocess.run(
            [*argv, "submit"],
            input=conflict.SerializeToString(),
            capture_output=True,
            check=False,
            timeout=20,
        )
        assert rejected.returncode and not rejected.stdout
        assert b"AlreadyExists" in rejected.stderr

        cancelled: vastai_job_pb2.SubmitJobRequest = _request(
            client,
            uri,
            experiment,
            archive,
            "touch started; sleep 25; touch escaped",
        )
        cancelled_ref: vastai_job_pb2.JobReference = (
            vastai_job_pb2.JobReference(job_id=cancelled.job_id)
        )
        cancelled_project: pathlib.Path = (
            state / "jobs" / cancelled.job_id / "project"
        )
        _call(argv, "submit", cancelled)
        _wait_for_state(argv, cancelled_ref, vastai_job_pb2.JOB_STATE_RUNNING)
        _wait_until((cancelled_project / "started").is_file)
        _call(argv, "cancel", cancelled_ref)
        _wait_for_state(argv, cancelled_ref, vastai_job_pb2.JOB_STATE_KILLED)
        assert not (cancelled_project / "escaped").exists()

        assert (
            _call(argv, "status", reference).state
            == vastai_job_pb2.JOB_STATE_RUNNING
        )

        # Make the first terminal write fail after the workload exits.
        with contextlib.closing(sqlite3.connect(state / "jobs.sqlite3")) as db:
            db.execute(
                "CREATE TRIGGER reject_completion BEFORE UPDATE OF phase ON jobs "
                "WHEN NEW.phase = 'finished' BEGIN "
                "SELECT RAISE(FAIL, 'test rejects completion'); END"
            )
            (project / "release").touch()
            _wait_until(
                lambda: b"test rejects completion" in worker_log.read_bytes(),
                timeout=10,
            )
            assert (
                _call(argv, "status", reference).state
                == vastai_job_pb2.JOB_STATE_RUNNING
            )
            db.execute("DROP TRIGGER reject_completion")
        completed: vastai_job_pb2.Job = _wait_for_state(
            argv, reference, vastai_job_pb2.JOB_STATE_FINISHED
        )
        assert completed.HasField("exit_code") and completed.exit_code == 0
        assert _call(argv, "submit", submission) == completed
        assert _call(argv, "cancel", reference) == completed
        assert (project / "attempts").read_text(
            encoding="utf-8"
        ) == "started\n"
        assert (job_files / "stdout.log").read_text(encoding="utf-8") == (
            f"project input\n{reference.job_id}\n{experiment}\n{uri}\n"
        )
        assert (job_files / "stderr.log").read_text(encoding="utf-8") == (
            "training stderr\n"
        )

        failed: vastai_job_pb2.SubmitJobRequest = _request(
            client,
            uri,
            experiment,
            archive,
            "sleep 60 & printf '%s\\n' \"$!\"; "
            "printf 'training failed\\n' >&2; exit 17",
        )
        _call(argv, "submit", failed)
        exit_result: vastai_job_pb2.Job = _wait_for_state(
            argv,
            vastai_job_pb2.JobReference(job_id=failed.job_id),
            vastai_job_pb2.JOB_STATE_FAILED,
        )
        assert (
            exit_result.HasField("exit_code") and exit_result.exit_code == 17
        )
        failed_files: pathlib.Path = state / "jobs" / failed.job_id
        child_pid: int = int(
            (failed_files / "stdout.log").read_text(encoding="utf-8")
        )
        _wait_until(lambda: _process_exited(child_pid), timeout=5)
        assert (failed_files / "stderr.log").read_text(
            encoding="utf-8"
        ) == "training failed\n"
        assert not (failed_files / "project").exists()
        assert not cancelled_project.exists()

        invalid_directory: pathlib.Path = root / "invalid"
        invalid_directory.mkdir()
        invalid_archive: pathlib.Path = invalid_directory / archive.name
        with tarfile.open(invalid_archive, "w:gz") as bundle:
            bundle.add(source, arcname="inputs/input.txt")
            unsupported: tarfile.TarInfo = tarfile.TarInfo("named_pipe")
            unsupported.type = tarfile.FIFOTYPE
            bundle.addfile(unsupported)
        invalid: vastai_job_pb2.SubmitJobRequest = _request(
            client, uri, experiment, invalid_archive, "exit 0"
        )
        _call(argv, "submit", invalid)
        invalid_result: vastai_job_pb2.Job = _wait_for_state(
            argv,
            vastai_job_pb2.JobReference(job_id=invalid.job_id),
            vastai_job_pb2.JOB_STATE_FAILED,
        )
        assert not invalid_result.HasField("exit_code")
        invalid_files: pathlib.Path = state / "jobs" / invalid.job_id
        assert not (invalid_files / "download").exists()
        assert not (invalid_files / "project").exists()

        interrupted: vastai_job_pb2.SubmitJobRequest = _request(
            client,
            uri,
            experiment,
            archive,
            "touch started; sleep 25; touch escaped",
        )
        interrupted_ref: vastai_job_pb2.JobReference = (
            vastai_job_pb2.JobReference(job_id=interrupted.job_id)
        )
        interrupted_project: pathlib.Path = (
            state / "jobs" / interrupted.job_id / "project"
        )
        _call(argv, "submit", interrupted)
        _wait_until((interrupted_project / "started").is_file)

    assert not (state / "worker.sock").exists()
    with _listener(argv, state) as service:
        assert (
            _call(argv, "status", interrupted_ref).state
            == vastai_job_pb2.JOB_STATE_KILLED
        )
        assert not interrupted_project.exists()
        assert _call(argv, "status", reference) == completed
        lost: vastai_job_pb2.SubmitJobRequest = _request(
            client,
            uri,
            experiment,
            archive,
            "touch started; while [ ! -f release ]; do sleep 0.1; done; touch exited",
        )
        lost_ref: vastai_job_pb2.JobReference = vastai_job_pb2.JobReference(
            job_id=lost.job_id
        )
        lost_project: pathlib.Path = state / "jobs" / lost.job_id / "project"
        _call(argv, "submit", lost)
        _wait_for_state(argv, lost_ref, vastai_job_pb2.JOB_STATE_RUNNING)
        _wait_until((lost_project / "started").is_file)
        try:
            service.kill()
            service.wait(timeout=10)
            with _listener(argv, state):
                lost_result: vastai_job_pb2.Job = _wait_for_state(
                    argv, lost_ref, vastai_job_pb2.JOB_STATE_UNKNOWN
                )
                assert (
                    lost_result.observation_error
                    and not lost_result.HasField("exit_code")
                )
                assert _call(argv, "submit", lost) == lost_result
        finally:
            (lost_project / "release").touch()
            _wait_until((lost_project / "exited").is_file)


def _main() -> None:
    """Run the service and MLflow CLI inside the Bazel test sandbox."""
    cli: pathlib.Path = pathlib.Path(sys.argv[1]).absolute()
    mlflow_cli: pathlib.Path = pathlib.Path(sys.argv[2]).absolute()
    os.environ["PATH"] = f"{mlflow_cli.parent}{os.pathsep}{os.environ['PATH']}"
    os.environ["MLFLOW_ENABLE_TELEMETRY"] = "false"
    # The Unix socket requires a short path. Bazel supplies a private /tmp.
    with tempfile.TemporaryDirectory(
        dir="/tmp", prefix="mlflow-worker-"
    ) as directory:
        root: pathlib.Path = pathlib.Path(directory)
        _exercise(root, cli)
        output: pathlib.Path = pathlib.Path(
            os.environ["TEST_UNDECLARED_OUTPUTS_DIR"]
        )
        with contextlib.closing(
            sqlite3.connect(root / "worker" / "jobs.sqlite3")
        ) as source:
            with contextlib.closing(
                sqlite3.connect(output / "jobs.sqlite3")
            ) as destination:
                source.backup(destination)


if __name__ == "__main__":
    _main()

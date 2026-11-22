"""Exercise MLflow plugin discovery and job RPCs against a local gRPC fixture."""

import collections.abc
import concurrent.futures
import functools
import json
import os
import pathlib
import random
import shlex
import subprocess
import sys
import tarfile
import tempfile
import unittest.mock

import grpc
import mlflow.entities
import mlflow.exceptions
import mlflow.projects
import mlflow.tracking
import mlflow_vastai_backend
import vastai_job_pb2
import vastai_job_pb2_grpc


class _JobFixture(vastai_job_pb2_grpc.JobServiceServicer):
    """Behave as the RPC peer and execute only the test's submitted workload."""

    def __init__(self, root: pathlib.Path) -> None:
        """Track fixture jobs and inject one ambiguous submit response."""
        self._root: pathlib.Path = root
        """Workspace where the fixture unpacks its test project."""
        self._requests: dict[str, vastai_job_pb2.SubmitJobRequest] = {}
        """Submitted requests indexed by idempotency key."""
        self._submission_peers: dict[str, str] = {}
        """Connection used for the initial submission and its retries."""
        self._states: dict[str, int] = {}
        """Workload states visible to the test client."""
        self._drops: int = 1
        """Number of submission replies to replace with a transport failure."""

    def drop_submit_replies(self, count: int) -> None:
        """Simulate reply loss on the next submission and its retries."""
        self._drops = count

    @property
    def submission_count(self) -> int:
        """Return the number of distinct jobs received by the fixture."""
        return len(self._requests)

    def SubmitJob(
        self,
        request: vastai_job_pb2.SubmitJobRequest,
        context: grpc.ServicerContext,
    ) -> vastai_job_pb2.Job:
        """Accept a job and lose its first response to exercise replay."""
        assert request.worker == "fixture-worker"
        assert request.ByteSize() < 4096
        assert (
            request.workload.project_artifact_path
            == "vastai-source/project.tar.gz"
        )
        client: mlflow.tracking.MlflowClient = mlflow.tracking.MlflowClient(
            tracking_uri=request.workload.tracking_uri
        )
        run: mlflow.entities.Run = client.get_run(request.job_id)
        assert run.info.experiment_id == request.workload.experiment_id
        assert run.info.status == "RUNNING"
        artifacts: list[mlflow.entities.FileInfo] = client.list_artifacts(
            request.job_id, "vastai-source"
        )
        assert len(artifacts) == 1
        assert artifacts[0].path == request.workload.project_artifact_path
        assert artifacts[0].file_size > 4 * 1024 * 1024
        assert (
            self._submission_peers.setdefault(request.job_id, context.peer())
            == context.peer()
        )
        if request.job_id in self._requests:
            assert self._requests[request.job_id] == request
        else:
            self._requests[request.job_id] = request
            self._states[request.job_id] = vastai_job_pb2.JOB_STATE_SCHEDULED
        if self._drops:
            self._drops -= 1
            context.abort(grpc.StatusCode.UNAVAILABLE, "fixture lost reply")
        return vastai_job_pb2.Job(
            job_id=request.job_id, state=self._states[request.job_id]
        )

    def GetJob(
        self,
        reference: vastai_job_pb2.JobReference,
        context: grpc.ServicerContext,
    ) -> vastai_job_pb2.Job:
        """Execute a scheduled fixture job once and confirm its result."""
        job_id: str = reference.job_id
        if self._states[job_id] == vastai_job_pb2.JOB_STATE_SCHEDULED:
            request: vastai_job_pb2.SubmitJobRequest = self._requests[job_id]
            client: mlflow.tracking.MlflowClient = (
                mlflow.tracking.MlflowClient(
                    tracking_uri=request.workload.tracking_uri
                )
            )
            with tempfile.TemporaryDirectory() as directory:
                archive_path: str = client.download_artifacts(
                    job_id, request.workload.project_artifact_path, directory
                )
                with tarfile.open(archive_path, mode="r:gz") as archive:
                    assert set(archive.getnames()) == {
                        "./MLproject",
                        "./train.py",
                        "./fixture-data.bin",
                    }
                    archive.extractall(self._root, filter="data")
            env: dict[str, str] = os.environ.copy()
            env.update(
                {
                    "MLFLOW_TRACKING_URI": request.workload.tracking_uri,
                    "MLFLOW_EXPERIMENT_ID": request.workload.experiment_id,
                    "MLFLOW_RUN_ID": job_id,
                }
            )
            env["PYTHONPATH"] = os.pathsep.join(sys.path)
            result: subprocess.CompletedProcess[bytes] = subprocess.run(
                ["/bin/sh", "-c", request.workload.command],
                cwd=self._root,
                env=env,
                timeout=15,
                check=True,
            )
            self._states[job_id] = vastai_job_pb2.JOB_STATE_FINISHED
            return vastai_job_pb2.Job(
                job_id=job_id,
                state=self._states[job_id],
                exit_code=result.returncode,
            )
        return vastai_job_pb2.Job(job_id=job_id, state=self._states[job_id])

    def CancelJob(
        self,
        reference: vastai_job_pb2.JobReference,
        context: grpc.ServicerContext,
    ) -> vastai_job_pb2.Job:
        """Confirm cancellation of a job waiting to run in the fixture."""
        assert (
            self._states[reference.job_id]
            == vastai_job_pb2.JOB_STATE_SCHEDULED
        )
        self._states[reference.job_id] = vastai_job_pb2.JOB_STATE_KILLED
        return vastai_job_pb2.Job(
            job_id=reference.job_id, state=vastai_job_pb2.JOB_STATE_KILLED
        )


def _main() -> None:
    """Run an MLflow project and verify its run ID, outputs, and cancellation."""
    os.environ["MLFLOW_DISABLE_AGENT_HINT"] = "1"
    os.environ["MLFLOW_ENABLE_TELEMETRY"] = "false"
    with (
        tempfile.TemporaryDirectory() as directory,
        concurrent.futures.ThreadPoolExecutor(max_workers=1) as executor,
    ):
        root: pathlib.Path = pathlib.Path(directory)
        source: pathlib.Path = root / "project"
        source.mkdir()
        (source / "fixture-data.bin").write_bytes(
            random.Random(0).randbytes(5 * 1024 * 1024)
        )
        (source / ".gitignore").write_text(".env\n")
        (source / ".env").write_text("FIXTURE=excluded\n")
        (source / "untracked.txt").write_text("excluded\n")
        (source / "MLproject").write_text(
            "name: fixture\nentry_points:\n  main:\n"
            "    parameters:\n      score: {type: float, default: 0}\n"
            "    command: "
            + json.dumps(f"{shlex.quote(sys.executable)} train.py {{score}}")
            + "\n"
        )
        (source / "train.py").write_text(
            "import sys\nimport mlflow\n"
            "with mlflow.start_run():\n"
            "    mlflow.log_metric('score', float(sys.argv[1]))\n"
            "    mlflow.log_text('training output\\n', 'result.txt')\n"
        )
        subprocess.run(["git", "init", "--quiet", str(source)], check=True)
        subprocess.run(
            ["git", "add", "--", "MLproject", "train.py", "fixture-data.bin"],
            cwd=source,
            check=True,
        )
        subprocess.run(
            [
                "git",
                "-c",
                "user.name=Fixture",
                "-c",
                "user.email=fixture@example.com",
                "-c",
                "commit.gpgsign=false",
                "commit",
                "--quiet",
                "-m",
                "Fixture",
            ],
            cwd=source,
            check=True,
        )
        subprocess.run(
            ["git", "tag", "training-version"], cwd=source, check=True
        )
        subprocess.run(
            ["git", "rm", "--cached", "--quiet", "--", "fixture-data.bin"],
            cwd=source,
            check=True,
        )
        fixture: _JobFixture = _JobFixture(root / "worker")
        server: grpc.Server = grpc.server(executor)
        vastai_job_pb2_grpc.add_JobServiceServicer_to_server(fixture, server)
        port: int = server.add_insecure_port("127.0.0.1:0")
        server.start()
        try:
            mlflow.set_tracking_uri(f"sqlite:///{root / 'tracking.db'}")
            config: dict[str, object] = {
                "server": f"127.0.0.1:{port}",
                "worker": "fixture-worker",
                "insecure": True,
                "project_files": ["fixture-data.bin"],
            }
            client: mlflow.tracking.MlflowClient = (
                mlflow.tracking.MlflowClient()
            )
            experiment_id: str = client.create_experiment(
                "fixture", artifact_location=(root / "artifacts").as_uri()
            )
            run_project: functools.partial[mlflow.projects.SubmittedRun] = (
                functools.partial(
                    mlflow.run,
                    str(source),
                    backend="vastai",
                    backend_config=config,
                )
            )
            run: mlflow.projects.SubmittedRun = run_project(
                parameters={"score": 42},
                experiment_id=experiment_id,
            )
            metrics: dict[str, float] = client.get_run(run.run_id).data.metrics
            assert metrics == {"score": 42.0}
            assert run.get_status() == "FINISHED"
            artifacts: list[str] = [
                item.path for item in client.list_artifacts(run.run_id)
            ]
            assert artifacts == ["result.txt", "vastai-source"]
            queued: mlflow.projects.SubmittedRun = run_project(
                experiment_id=experiment_id,
                synchronous=False,
            )
            queued.cancel()
            assert queued.get_status() == "KILLED"
            assert not queued.wait()
            blocked: pathlib.Path = root / "blocked-artifact-store"
            blocked.touch()
            failing_experiment: str = client.create_experiment(
                "upload-failure", artifact_location=blocked.as_uri()
            )
            try:
                run_project(
                    experiment_id=failing_experiment,
                )
            except OSError:
                pass
            else:
                raise AssertionError("Source upload unexpectedly succeeded")
            assert fixture.submission_count == 2
            failed_runs: list[mlflow.entities.Run] = client.search_runs(
                [failing_experiment]
            )
            assert [run.info.status for run in failed_runs] == ["FAILED"]
            channel_experiment: str = client.create_experiment(
                "channel-failure",
                artifact_location=(root / "channel-artifacts").as_uri(),
            )
            channel_error: OSError = OSError("fixture channel setup failed")
            with unittest.mock.patch.object(
                grpc, "insecure_channel", side_effect=channel_error
            ):
                try:
                    run_project(
                        experiment_id=channel_experiment,
                    )
                except OSError as error:
                    assert error is channel_error
                else:
                    raise AssertionError(
                        "Channel setup unexpectedly succeeded"
                    )
            assert fixture.submission_count == 2
            assert not client.search_runs([channel_experiment])

            ambiguous_experiment: str = client.create_experiment(
                "submit-replies-lost",
                artifact_location=(root / "ambiguous-artifacts").as_uri(),
            )
            fixture.drop_submit_replies(2)
            ambiguous_error: str
            try:
                run_project(
                    experiment_id=ambiguous_experiment,
                )
            except mlflow.exceptions.ExecutionException as error:
                assert isinstance(error.__cause__, grpc.RpcError)
                assert error.__cause__.code() == grpc.StatusCode.UNAVAILABLE
                ambiguous_error = str(error)
            else:
                raise AssertionError("A submission reply unexpectedly arrived")
            assert fixture.submission_count == 3
            ambiguous_runs: list[mlflow.entities.Run] = client.search_runs(
                [ambiguous_experiment]
            )
            assert [run.info.status for run in ambiguous_runs] == ["RUNNING"]
            assert ambiguous_runs[0].info.run_id in ambiguous_error
            assert "unconfirmed" in ambiguous_error

            tracking_uri: str = f"sqlite:///{root / 'other-tracking.db'}"
            tracking_client: mlflow.tracking.MlflowClient = (
                mlflow.tracking.MlflowClient(tracking_uri)
            )
            target_experiment: str = tracking_client.create_experiment(
                "separate-tracking",
                artifact_location=(root / "other-artifacts").as_uri(),
            )
            adapter: mlflow_vastai_backend.VastAIBackend = (
                mlflow_vastai_backend.VastAIBackend()
            )
            with (
                mlflow.start_run(experiment_id=experiment_id) as parent,
                unittest.mock.patch.object(
                    mlflow.projects.utils,
                    "_fetch_git_repo",
                    wraps=mlflow.projects.utils._fetch_git_repo,
                ) as fetch,
            ):
                separate_run: mlflow.projects.SubmittedRun = adapter.run(
                    str(source),
                    "main",
                    {"score": 7},
                    "training-version",
                    config,
                    tracking_uri,
                    target_experiment,
                )
            checkout: pathlib.Path = pathlib.Path(fetch.call_args.args[2])
            assert checkout != source
            assert not checkout.exists()
            assert subprocess.run(
                ["git", "symbolic-ref", "HEAD"],
                cwd=source,
                capture_output=True,
                check=True,
            ).stdout
            data: mlflow.entities.RunData = tracking_client.get_run(
                separate_run.run_id
            ).data
            assert float(data.params["score"]) == 7
            assert data.tags["mlflow.project.entryPoint"] == "main"
            assert "mlflow.parentRunId" not in data.tags
            assert {
                run.info.run_id for run in client.search_runs([experiment_id])
            } == {run.run_id, queued.run_id, parent.info.run_id}

            cleanup_experiment: str = client.create_experiment(
                "checkout-cleanup-failure",
                artifact_location=(root / "cleanup-artifacts").as_uri(),
            )
            cleanup_error: OSError = OSError("fixture checkout cleanup failed")
            cleanup: collections.abc.Callable[
                [tempfile.TemporaryDirectory[str]], None
            ] = tempfile.TemporaryDirectory.cleanup

            def fail_checkout_cleanup(
                directory: tempfile.TemporaryDirectory[str],
            ) -> None:
                """Inject a failure when releasing the versioned checkout."""
                cleanup(directory)
                if (
                    fetch.call_args
                    and directory.name == fetch.call_args.args[2]
                ):
                    raise cleanup_error

            with (
                unittest.mock.patch.object(
                    mlflow.projects.utils,
                    "_fetch_git_repo",
                    wraps=mlflow.projects.utils._fetch_git_repo,
                ) as fetch,
                unittest.mock.patch.object(
                    tempfile.TemporaryDirectory,
                    "cleanup",
                    autospec=True,
                    side_effect=fail_checkout_cleanup,
                ),
            ):
                try:
                    run_project(
                        version="training-version",
                        experiment_id=cleanup_experiment,
                    )
                except OSError as error:
                    assert error is cleanup_error
                else:
                    raise AssertionError(
                        "Checkout cleanup unexpectedly succeeded"
                    )
            assert fixture.submission_count == 4
            assert [
                run.info.status
                for run in client.search_runs([cleanup_experiment])
            ] == ["FAILED"]

            artifact_dir: str | None = os.environ.get(
                "TEST_UNDECLARED_OUTPUTS_DIR"
            )
            if artifact_dir:
                (pathlib.Path(artifact_dir) / "results.json").write_text(
                    json.dumps(
                        {
                            "metrics": metrics,
                            "artifacts": artifacts,
                            "cancelled": True,
                            "upload_failure_prevented_submission": True,
                            "channel_failure_prevented_run_creation": True,
                            "lost_replies_left_run_active": True,
                            "project_snapshot_boundary_verified": True,
                            "separate_tracking_server_verified": True,
                            "temporary_checkout_removed": True,
                            "cleanup_failure_prevented_submission": True,
                        }
                    )
                    + "\n"
                )
        finally:
            server.stop(None).wait()


if __name__ == "__main__":
    _main()

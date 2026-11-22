"""Stage and submit MLflow runs and invoke the generated job service."""

import contextlib
import types
import typing

import grpc
import mlflow.entities
import mlflow.exceptions
import vastai_job_pb2
import vastai_job_pb2_grpc

from mlflow_vastai_backend import _config, _tracking

_RPC_TIMEOUT: float = 10.0
"""Maximum time to wait for a single job service call."""


class ChannelFactory(typing.Protocol):
    """Create a channel for one MLflow operation."""

    def __call__(self, config: _config.BackendConfig, /) -> grpc.Channel:
        """Open a channel from the validated backend configuration."""


class JobServiceSession:
    """Call the generated gRPC service within one channel scope."""

    def __init__(self, channel: grpc.Channel) -> None:
        """Take ownership of the channel for one MLflow operation."""
        self._channel: grpc.Channel = channel
        """Channel shared by the RPCs in this operation."""

    def __enter__(self) -> typing.Self:
        """Initialize the stub, closing the channel if setup fails."""
        with contextlib.ExitStack() as setup:
            setup.push(self)
            self._stub: vastai_job_pb2_grpc.JobServiceStub = (
                vastai_job_pb2_grpc.JobServiceStub(self._channel)
            )
            """Stub shared by the RPCs in this operation."""
            setup.pop_all()
        return self

    def __exit__(
        self,
        error_type: type[BaseException] | None,
        error: BaseException | None,
        traceback: types.TracebackType | None,
    ) -> None:
        """Close the channel and retain any exception raised by the operation."""
        try:
            self._channel.close()
        except Exception as cleanup_error:
            if error is None:
                raise
            error.add_note(
                f"Closing the job service channel also failed ({cleanup_error!r})."
            )

    def submit_job(self, request: vastai_job_pb2.SubmitJobRequest) -> None:
        """Submit a job, retrying a transport failure once with the same job ID."""
        for attempt in (1, 2):
            try:
                self._call(self._stub.SubmitJob, request)
                return
            except (
                grpc.RpcError,
                mlflow.exceptions.ExecutionException,
            ) as error:
                code: grpc.StatusCode | None = (
                    error.code() if isinstance(error, grpc.RpcError) else None
                )
                if attempt == 1 and code in (
                    grpc.StatusCode.UNAVAILABLE,
                    grpc.StatusCode.DEADLINE_EXCEEDED,
                ):
                    continue
                detail: str = code.name if code is not None else str(error)
                raise mlflow.exceptions.ExecutionException(
                    f"Submission unconfirmed for MLflow run {request.job_id} ({detail})"
                ) from error

    def get_job(self, request: vastai_job_pb2.JobReference) -> int:
        """Read the confirmed state of an existing job."""
        return self._observe(self._stub.GetJob, request)

    def cancel_job(self, request: vastai_job_pb2.JobReference) -> int:
        """Request cancellation and return the confirmed state."""
        return self._observe(self._stub.CancelJob, request)

    def _observe(
        self,
        rpc: grpc.UnaryUnaryMultiCallable,
        request: vastai_job_pb2.JobReference,
    ) -> int:
        """Expose observation transport failures through MLflow."""
        try:
            return self._call(rpc, request)
        except grpc.RpcError as error:
            raise mlflow.exceptions.ExecutionException(
                f"Vast.ai job service RPC failed ({error.code().name})"
            ) from error

    def _call(
        self,
        rpc: grpc.UnaryUnaryMultiCallable,
        request: vastai_job_pb2.SubmitJobRequest | vastai_job_pb2.JobReference,
    ) -> int:
        """Invoke one generated RPC and validate its reply."""
        job: vastai_job_pb2.Job = rpc(request, timeout=_RPC_TIMEOUT)
        if job.job_id != request.job_id:
            raise mlflow.exceptions.ExecutionException(
                "Vast.ai job service returned a different job ID"
            )
        if job.observation_error or job.state in (
            vastai_job_pb2.JOB_STATE_UNSPECIFIED,
            vastai_job_pb2.JOB_STATE_UNKNOWN,
        ):
            raise mlflow.exceptions.ExecutionException(
                "Vast.ai job service could not confirm the workload status"
            )
        if (
            job.state
            not in vastai_job_pb2.JobState.DESCRIPTOR.values_by_number
        ):
            raise mlflow.exceptions.ExecutionException(
                "Vast.ai job service returned an unknown status"
            )
        return job.state


class JobService:
    """Submit MLflow runs and open connections for later job operations."""

    def __init__(
        self,
        run_config: _config.RunConfig,
        channel_factory: ChannelFactory,
        tracking: _tracking.TrackingService,
    ) -> None:
        """Bind the channel configuration and tracking service for this job."""
        self._run_config: _config.RunConfig = run_config
        """Validated settings for this job."""
        self._channel_factory: ChannelFactory = channel_factory
        """Factory for independent operation channels."""
        self._tracking: _tracking.TrackingService = tracking
        """Tracking service for staging project runs."""

    def submit_run(self) -> str:
        """Stage and submit the project, then return its MLflow run ID."""
        run_config: _config.RunConfig = self._run_config
        with self.connect() as session:
            run: mlflow.entities.RunInfo
            command: str
            run, command = self._tracking.prepare_project()
            request: vastai_job_pb2.SubmitJobRequest = (
                vastai_job_pb2.SubmitJobRequest(
                    job_id=run.run_id,
                    worker=run_config.backend_config.worker,
                    workload=vastai_job_pb2.Workload(
                        project_artifact_path=str(_tracking.PROJECT_ARTIFACT),
                        command=command,
                        tracking_uri=run_config.tracking_uri,
                        experiment_id=run.experiment_id,
                    ),
                )
            )
            session.submit_job(request)
        return request.job_id

    def connect(self) -> JobServiceSession:
        """Open an independently scoped channel to the service."""
        return JobServiceSession(
            self._channel_factory(self._run_config.backend_config)
        )

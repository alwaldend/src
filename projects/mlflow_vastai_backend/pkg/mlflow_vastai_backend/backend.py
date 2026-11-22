"""Expose the Vast.ai execution backend to MLflow Projects."""

import grpc
import mlflow.exceptions
import mlflow.projects
import mlflow.tracking
import pydantic

from mlflow_vastai_backend import (
    _config,
    _job_service,
    _submitted_run,
    _tracking,
)


def _open_channel(config: _config.BackendConfig) -> grpc.Channel:
    """Open a channel with the configured transport security."""
    if config.insecure:
        return grpc.insecure_channel(config.server)
    return grpc.secure_channel(config.server, grpc.ssl_channel_credentials())


class VastAIBackend(mlflow.projects.backend.AbstractBackend):
    """Prepare and submit projects for the registered Vast.ai backend."""

    def __init__(
        self,
        *,
        channel_factory: _job_service.ChannelFactory = _open_channel,
        tracking_client_factory: _tracking.TrackingClientFactory = mlflow.tracking.MlflowClient,
    ) -> None:
        """Accept client factories. Each channel returned transfers ownership."""
        self._channel_factory: _job_service.ChannelFactory = channel_factory
        """Factory for the channels used by the backend and its run handles."""
        self._tracking_client_factory: _tracking.TrackingClientFactory = (
            tracking_client_factory
        )
        """Factory for the artifact tracking client."""

    def run(
        self,
        project_uri: str,
        entry_point: str,
        params: _config.ProjectParameters | None,
        version: str | None,
        backend_config: dict[str, object],
        tracking_uri: str,
        experiment_id: str,
    ) -> mlflow.projects.submitted_run.SubmittedRun:
        """Create the tracking run, submit one prepared job, and return its handle."""
        try:
            run_config: _config.RunConfig = _config.RunConfig.model_validate(
                {
                    "project_uri": project_uri,
                    "entry_point": entry_point,
                    "params": params,
                    "version": version,
                    "backend_config": backend_config,
                    "tracking_uri": tracking_uri,
                    "experiment_id": experiment_id,
                }
            )
        except pydantic.ValidationError as error:
            raise mlflow.exceptions.ExecutionException(str(error)) from error
        tracking: _tracking.TrackingService = _tracking.TrackingService(
            run_config, self._tracking_client_factory(run_config.tracking_uri)
        )
        service: _job_service.JobService = _job_service.JobService(
            run_config,
            self._channel_factory,
            tracking,
        )
        run_id: str = service.submit_run()
        return _submitted_run.SubmittedRun(service, run_id)

"""Expose a submitted Vast.ai job through MLflow's run handle."""

import time

import mlflow.projects
import vastai_job_pb2

from mlflow_vastai_backend import _job_service

_POLL_INTERVAL: float = 1.0
"""Time between observations while waiting for a job."""


_TERMINAL_STATES: frozenset[int] = frozenset(
    (
        vastai_job_pb2.JOB_STATE_FINISHED,
        vastai_job_pb2.JOB_STATE_FAILED,
        vastai_job_pb2.JOB_STATE_KILLED,
    )
)
"""Statuses for which waiting can return the workload outcome."""


class SubmittedRun(mlflow.projects.submitted_run.SubmittedRun):
    """Translate MLflow lifecycle calls into job service requests."""

    def __init__(self, service: _job_service.JobService, run_id: str) -> None:
        """Bind the handle to one MLflow run."""
        self._service: _job_service.JobService = service
        """Service configured for the job's server."""
        self._job: vastai_job_pb2.JobReference = vastai_job_pb2.JobReference(
            job_id=run_id
        )
        """Reference carrying the MLflow run ID."""

    @property
    def run_id(self) -> str:
        """Return the tracking run ID recognized by MLflow."""
        return self._job.job_id

    def get_status(self) -> str:
        """Read the confirmed job status."""
        with self._service.connect() as session:
            state: int = session.get_job(self._job)
        return vastai_job_pb2.JobState.Name(state).removeprefix("JOB_STATE_")

    def _wait_for_terminal(
        self, session: _job_service.JobServiceSession, state: int
    ) -> int:
        """Poll until the service reports a terminal state."""
        while state not in _TERMINAL_STATES:
            time.sleep(_POLL_INTERVAL)
            state = session.get_job(self._job)
        return state

    def wait(self) -> bool:
        """Return whether the job finished successfully."""
        with self._service.connect() as session:
            state: int = self._wait_for_terminal(
                session, session.get_job(self._job)
            )
        return state == vastai_job_pb2.JOB_STATE_FINISHED

    def cancel(self) -> None:
        """Request cancellation and wait for a terminal state."""
        with self._service.connect() as session:
            self._wait_for_terminal(session, session.cancel_job(self._job))

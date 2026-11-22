"""Create MLflow project runs and stage snapshots in their artifact stores."""

import collections.abc
import contextlib
import os
import pathlib
import subprocess
import tarfile
import tempfile
import typing

import mlflow
import mlflow.entities
import mlflow.exceptions
import mlflow.projects
import mlflow.tracking

from mlflow_vastai_backend import _config

PROJECT_ARTIFACT: pathlib.PurePosixPath = pathlib.PurePosixPath(
    "vastai-source/project.tar.gz"
)
"""Run-relative location of the project snapshot."""


class TrackingClientFactory(typing.Protocol):
    """Supply a client for the selected tracking server."""

    def __call__(self, tracking_uri: str, /) -> mlflow.tracking.MlflowClient:
        """Provide the client used for all writes to a run."""


def _write_project_archive(
    work_dir: pathlib.Path, destination: pathlib.Path, project_files: list[str]
) -> None:
    """Package regular files selected by Git or explicitly listed by the user."""
    root: pathlib.Path = work_dir.resolve()
    files: set[str] = set(project_files)
    try:
        tracked: subprocess.CompletedProcess[bytes] = subprocess.run(
            ["git", "ls-files", "--cached", "-z", "--", "."],
            cwd=root,
            capture_output=True,
            check=False,
        )
        if tracked.returncode == 0:
            files.update(
                os.fsdecode(name)
                for name in tracked.stdout.split(b"\0")
                if name
            )
    except FileNotFoundError:
        pass
    if not files:
        raise mlflow.exceptions.ExecutionException(
            "No project files selected. Track them in Git or set project_files."
        )
    with tarfile.open(destination, mode="w:gz") as archive:
        for name in sorted(files):
            relative: pathlib.PurePosixPath = pathlib.PurePosixPath(name)
            path: pathlib.Path = root / relative
            if (
                relative.is_absolute()
                or ".." in relative.parts
                or ".git" in relative.parts
                or not path.is_relative_to(root)
                or path.resolve() != path
                or not path.is_file()
            ):
                raise mlflow.exceptions.ExecutionException(
                    f"Expected a regular file inside the project: {name}"
                )
            archive.add(path, arcname=f"./{relative}", recursive=False)


def _project_tags(
    config: _config.RunConfig, work_dir: str, tracking_uri: str
) -> dict[str, str]:
    """Use the same project provenance tags as MLflow Projects."""
    source_name: str = config.project_uri
    if mlflow.projects.utils._is_local_uri(source_name):
        source_name = (
            mlflow.tracking._tracking_service.utils._get_git_url_if_present(
                mlflow.projects.utils._expand_uri(source_name)
            )
        )
    tags: dict[str, str] = {
        mlflow.projects.utils.MLFLOW_USER: mlflow.projects.utils._get_user(),
        mlflow.projects.utils.MLFLOW_SOURCE_NAME: source_name,
        mlflow.projects.utils.MLFLOW_SOURCE_TYPE: mlflow.entities.SourceType.to_string(
            mlflow.entities.SourceType.PROJECT
        ),
        mlflow.projects.utils.MLFLOW_PROJECT_ENTRY_POINT: config.entry_point,
    }
    commit: str | None = mlflow.projects.utils.get_git_commit(work_dir)
    if commit is not None:
        tags[mlflow.projects.utils.MLFLOW_GIT_COMMIT] = commit
    parent: mlflow.entities.Run | None = (
        mlflow.active_run()
        if tracking_uri == mlflow.get_tracking_uri()
        else None
    )
    if parent is not None:
        tags[mlflow.projects.utils.MLFLOW_PARENT_RUN_ID] = parent.info.run_id
    repository: str | None = mlflow.projects.utils.get_git_repo_url(work_dir)
    if repository is not None:
        tags[mlflow.projects.utils.MLFLOW_GIT_REPO_URL] = repository
    if (
        config.version is not None
        and mlflow.projects.utils._is_valid_branch_name(
            work_dir, config.version
        )
    ):
        tags[mlflow.projects.utils.MLFLOW_GIT_BRANCH] = config.version
    return tags


def _project_subdirectory(root: pathlib.Path, subdirectory: str) -> str:
    """Resolve the project directory inside its source tree."""
    project: pathlib.Path = (root / subdirectory).resolve()
    if not project.is_dir() or not project.is_relative_to(root.resolve()):
        raise mlflow.exceptions.ExecutionException(
            f"Expected a project directory inside the source tree: {subdirectory}"
        )
    return str(project)


@contextlib.contextmanager
def _project_directory(
    config: _config.RunConfig,
) -> collections.abc.Iterator[str]:
    """Use local sources in place and own any fetched checkout or extraction."""
    uri: str
    subdirectory: str
    uri, subdirectory = mlflow.projects.utils._parse_subdirectory(
        config.project_uri
    )
    local: bool = mlflow.projects.utils._is_local_uri(uri)
    zipped: bool = mlflow.projects.utils._is_zip_uri(uri)
    if local:
        uri = mlflow.projects.utils._parse_file_uri(uri)
        if not zipped and config.version is None:
            yield _project_subdirectory(pathlib.Path(uri), subdirectory)
            return
        if not zipped and not mlflow.projects.utils._is_git_repo(uri):
            raise mlflow.exceptions.ExecutionException(
                "Setting a version is only supported for Git project URIs"
            )
        uri = str(pathlib.Path(uri).resolve())
    with tempfile.TemporaryDirectory() as directory:
        if zipped:
            mlflow.projects.utils._unzip_repo(
                uri if local else mlflow.projects.utils._fetch_zip_repo(uri),
                directory,
            )
        else:
            mlflow.projects.utils._fetch_git_repo(
                uri, config.version, directory
            )
        yield _project_subdirectory(pathlib.Path(directory), subdirectory)


def _load_entry_point(
    work_dir: str, entry_point: str
) -> mlflow.projects._project_spec.EntryPoint:
    """Load and validate the entry point for the worker's environment."""
    project: mlflow.projects._project_spec.Project = (
        mlflow.projects.utils.load_project(work_dir)
    )
    if project.docker_env:
        raise mlflow.exceptions.ExecutionException(
            "Vast.ai backend requires a preconfigured worker environment"
        )
    command: mlflow.projects._project_spec.EntryPoint = (
        project.get_entry_point(entry_point)
    )
    if any(
        parameter.type == "path" for parameter in command.parameters.values()
    ):
        raise mlflow.exceptions.ExecutionException(
            "Vast.ai backend does not support MLproject path parameters"
        )

    return command


class TrackingService:
    """Stage one project using a single MLflow tracking client."""

    def __init__(
        self,
        run_config: _config.RunConfig,
        client: mlflow.tracking.MlflowClient,
    ) -> None:
        """Bind the run configuration to its tracking client."""
        self._run_config: _config.RunConfig = run_config
        """Validated settings for the run being staged."""
        self._client: mlflow.tracking.MlflowClient = client
        """Client for project metadata and artifacts."""

    def prepare_project(self) -> tuple[mlflow.entities.RunInfo, str]:
        """Prepare and stage the project, then return its run and command."""
        run: mlflow.entities.Run | None = None
        try:
            with (
                _project_directory(self._run_config) as work_dir,
                tempfile.TemporaryDirectory() as directory,
            ):
                entry_point: mlflow.projects._project_spec.EntryPoint = (
                    _load_entry_point(work_dir, self._run_config.entry_point)
                )
                command: str = entry_point.compute_command(
                    self._run_config.params, storage_dir=None
                )
                final_params: dict[str, str]
                extra_params: dict[str, str]
                final_params, extra_params = entry_point.compute_parameters(
                    self._run_config.params, storage_dir=None
                )
                archive: pathlib.Path = (
                    pathlib.Path(directory) / PROJECT_ARTIFACT.name
                )
                _write_project_archive(
                    pathlib.Path(work_dir),
                    archive,
                    self._run_config.backend_config.project_files,
                )
                run = self._client.create_run(
                    experiment_id=self._run_config.experiment_id,
                    tags=_project_tags(
                        self._run_config, work_dir, self._client.tracking_uri
                    ),
                )
                run_id: str = run.info.run_id
                self._client.log_batch(
                    run_id,
                    params=[
                        mlflow.entities.Param(name, value)
                        for name, value in (
                            final_params | extra_params
                        ).items()
                    ],
                )
                self._client.log_artifact(
                    run_id, str(archive), str(PROJECT_ARTIFACT.parent)
                )
            return run.info, command
        except Exception as error:
            if run is not None:
                self._mark_failed(run.info.run_id, error)
            raise

    def _mark_failed(self, run_id: str, error: Exception) -> None:
        """Record a staging failure while retaining the original exception."""
        status: str = mlflow.entities.RunStatus.to_string(
            mlflow.entities.RunStatus.FAILED
        )
        try:
            self._client.set_terminated(run_id, status=status)
        except Exception as tracking_error:
            error.add_note(
                f"Could not mark MLflow run {run_id} as failed "
                f"({type(tracking_error).__name__})."
            )

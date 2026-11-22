"""Validate the settings supplied by MLflow Projects."""

import collections.abc
import typing

import pydantic

ProjectParameters: typing.TypeAlias = collections.abc.Mapping[str, object]
"""Project-defined names and values interpreted by MLflow."""


class BackendConfig(pydantic.BaseModel):
    """Job service settings received through MLflow's backend configuration."""

    model_config: typing.ClassVar[pydantic.ConfigDict] = pydantic.ConfigDict(
        frozen=True, strict=True
    )
    """Preserve validated settings and require their declared input types."""

    server: str = pydantic.Field(pattern=r"\S")
    """Nonblank job service gRPC address."""

    worker: str = ""
    """Optional worker selector configured on the job service."""

    insecure: bool = False
    """Allow cleartext for a trusted local job service."""

    project_files: list[str] = pydantic.Field(default_factory=list)
    """Additional project-relative files to upload, or the file list outside Git."""

    env_manager: typing.Literal["local"] | None = pydantic.Field(
        default=None, validation_alias="ENV_MANAGER"
    )
    """Environment manager supplied by MLflow Projects."""


class RunConfig(pydantic.BaseModel):
    """Validated project, tracking, and service settings for one run."""

    model_config: typing.ClassVar[pydantic.ConfigDict] = pydantic.ConfigDict(
        frozen=True, strict=True, hide_input_in_errors=True
    )
    """Validate the input types and report errors without their input values."""

    project_uri: str
    """Source URI understood by MLflow Projects."""

    entry_point: str
    """MLproject entry point to execute."""

    params: ProjectParameters | None
    """Parameters supplied to the entry point."""

    version: str | None
    """Optional source revision to fetch."""

    tracking_uri: str
    """MLflow tracking endpoint used by the adapter and workload."""

    experiment_id: str
    """Experiment that owns the tracking run."""

    backend_config: BackendConfig
    """Service and worker settings validated as part of the run."""

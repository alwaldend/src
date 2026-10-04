"""Shared health response contract."""

import typing

import pydantic


class HealthResponse(pydantic.BaseModel):
    status: typing.Literal["ok"] = "ok"

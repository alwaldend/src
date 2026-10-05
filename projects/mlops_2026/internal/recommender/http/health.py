"""Unversioned health route."""

import fastapi

from projects.mlops_2026.internal import models


def health() -> models.HealthResponse:
    return models.HealthResponse()


def create_health_router() -> fastapi.APIRouter:
    router: fastapi.APIRouter = fastapi.APIRouter()
    router.add_api_route(
        "/health",
        health,
        methods=["GET"],
        response_model=models.HealthResponse,
    )
    return router

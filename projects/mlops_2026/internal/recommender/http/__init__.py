"""Public unversioned HTTP API."""

from projects.mlops_2026.internal.recommender.http.health import (
    create_health_router,
)

__all__: tuple[str, ...] = ("create_health_router",)

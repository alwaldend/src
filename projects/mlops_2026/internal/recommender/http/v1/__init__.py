"""Public version-one HTTP API."""

from projects.mlops_2026.internal.recommender.http.v1.router import (
    RecommendationRouter,
    create_v1_router,
)

__all__: tuple[str, ...] = ("RecommendationRouter", "create_v1_router")

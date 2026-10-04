"""Public recommendation service API."""

from projects.mlops_2026.internal.recommender.service.recommendations import (
    RecommendationService,
)

__all__: tuple[str, ...] = ("RecommendationService",)

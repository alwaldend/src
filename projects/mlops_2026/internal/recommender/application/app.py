"""Assemble an independent application with an injected service."""

import fastapi

from projects.mlops_2026.internal.recommender import http, service
from projects.mlops_2026.internal.recommender.application import data
from projects.mlops_2026.internal.recommender.http import v1


class RecommenderApplication:
    def __init__(self, recommendations: service.RecommendationService) -> None:
        self._recommendations: service.RecommendationService = recommendations
        self._app: fastapi.FastAPI = fastapi.FastAPI(
            title="recommender", lifespan=self
        )
        self._app.include_router(http.create_health_router())
        self._app.include_router(v1.create_v1_router(recommendations))

    @property
    def app(self) -> fastapi.FastAPI:
        return self._app

    def __call__(self, app: fastapi.FastAPI) -> service.RecommendationService:
        return self._recommendations


def create_app() -> fastapi.FastAPI:
    recommendations: service.RecommendationService = (
        service.RecommendationService(data.model_path())
    )
    return RecommenderApplication(recommendations).app

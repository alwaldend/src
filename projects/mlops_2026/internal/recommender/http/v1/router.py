"""Versioned recommendation routes."""

import typing

import fastapi

from projects.mlops_2026.internal import models, recommender


class RecommendationRouter:
    def __init__(
        self, recommendations: recommender.service.RecommendationService
    ) -> None:
        self._recommendations: recommender.service.RecommendationService = (
            recommendations
        )
        self._router: fastapi.APIRouter = fastapi.APIRouter(prefix="/api/v1")
        self._router.add_api_route(
            "/recommend",
            self.recommend,
            methods=["GET"],
            response_model=models.RecommenderRecommendationResponse,
        )

    @property
    def router(self) -> fastapi.APIRouter:
        return self._router

    def recommend(
        self,
        track_name: typing.Annotated[
            str, fastapi.Query(description="Requested track title")
        ],
        n: typing.Annotated[
            int,
            fastapi.Query(
                alias="n",
                ge=1,
                le=20,
                description="Recommendation count (from 1 to 20)",
            ),
        ] = 5,
    ) -> models.RecommenderRecommendationResponse:
        return self._recommendations.recommend(track_name, n)


def create_v1_router(
    recommendations: recommender.service.RecommendationService,
) -> fastapi.APIRouter:
    return RecommendationRouter(recommendations).router

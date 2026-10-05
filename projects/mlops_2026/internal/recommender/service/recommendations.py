"""Own a recommendation model for one application lifespan."""

import pathlib
import types
import typing

import recommender

from projects.mlops_2026.internal import models


class RecommendationService:
    def __init__(self, model_path: pathlib.Path) -> None:
        self._model: recommender.RecommenderModel | None = None
        self._model_path: pathlib.Path = model_path

    async def __aenter__(self) -> None:
        self._model = typing.cast(
            recommender.RecommenderModel,
            recommender.load_model(str(self._model_path)),
        )

    async def __aexit__(
        self,
        exc_type: type[BaseException] | None,
        exc: BaseException | None,
        traceback: types.TracebackType | None,
    ) -> None:
        self.close()

    def __del__(self) -> None:
        self.close()

    def close(self) -> None:
        self._model = None

    def recommend(
        self, track_name: str, count: int
    ) -> models.RecommenderRecommendationResponse:
        model: recommender.RecommenderModel | None = self._model
        if model is None:
            raise RuntimeError("Recommendation model is not initialized")
        songs: list[models.RecommenderSong] = []
        record: dict[str, str]
        for record in model.recommend(track_name, N=count):
            songs.append(models.RecommenderSong.model_validate(record))
        return models.RecommenderRecommendationResponse(
            requested_track=track_name, recommendations=songs
        )

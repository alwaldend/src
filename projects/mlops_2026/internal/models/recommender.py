"""Shared recommendation response contracts."""

import pydantic


class RecommenderSong(pydantic.BaseModel):
    track_name: str = pydantic.Field(..., examples=["Time"])
    artists: str = pydantic.Field(..., examples=["Pink Floyd"])
    album_name: str = pydantic.Field(
        ..., examples=["The Dark Side of the Moon"]
    )


class RecommenderRecommendationResponse(pydantic.BaseModel):
    requested_track: str = pydantic.Field(..., examples=["Time"])
    recommendations: list[RecommenderSong]

"""Public application factory and packaged-data API."""

from projects.mlops_2026.internal.recommender.application.app import (
    RecommenderApplication,
    create_app,
)
from projects.mlops_2026.internal.recommender.application.data import (
    course_data_paths,
    model_path,
)

__all__: tuple[str, ...] = (
    "RecommenderApplication",
    "course_data_paths",
    "create_app",
    "model_path",
)

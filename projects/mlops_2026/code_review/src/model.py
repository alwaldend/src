"""Provide the lesson's typed, seeded taxi-fare estimator interface."""

import numpy
import numpy.typing
import pandas
import sklearn.ensemble


class TaxiFareModel:
    """Train and expose a gradient boosting regressor for artifact consumers."""

    def __init__(self, random_state: int) -> None:
        """Fix the estimator seed for reproducible training."""
        self._model: sklearn.ensemble.GradientBoostingRegressor = (
            sklearn.ensemble.GradientBoostingRegressor(
                random_state=random_state
            )
        )
        """Underlying estimator owned by this wrapper."""

    @property
    def model(self) -> sklearn.ensemble.GradientBoostingRegressor:
        """Expose the trained estimator for scoring and pickle serialization."""
        return self._model

    def fit(self, x: pandas.DataFrame, y: pandas.Series) -> None:
        """Fit the regressor to numeric features and fare targets."""
        self._model.fit(x, y)

    def predict(
        self, x: pandas.DataFrame
    ) -> numpy.typing.NDArray[numpy.float64]:
        """Predict fares in the same feature order used for training."""
        return self._model.predict(x)

"""Load taxi records and split numeric features from the fare target."""

import pathlib

import pandas
import sklearn.model_selection


def load_data(
    path: str | pathlib.Path, limit: int | None = None
) -> pandas.DataFrame:
    """Read all records, or an explicit prefix for a bounded exercise run."""
    return pandas.read_csv(path, nrows=limit)


def split_data(
    data: pandas.DataFrame,
    test_size: float,
    random_state: int,
) -> tuple[pandas.DataFrame, pandas.DataFrame, pandas.Series, pandas.Series]:
    """Exclude identifiers and the target, then make a seeded holdout split."""
    features: pandas.DataFrame = data.drop(
        columns=["fare_amount", "key", "Unnamed: 0"]
    )
    target: pandas.Series = data["fare_amount"]
    x_train: pandas.DataFrame
    x_test: pandas.DataFrame
    y_train: pandas.Series
    y_test: pandas.Series
    x_train, x_test, y_train, y_test = (
        sklearn.model_selection.train_test_split(
            features,
            target,
            test_size=test_size,
            random_state=random_state,
        )
    )
    return x_train, x_test, y_train, y_test

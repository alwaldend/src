"""Filter incomplete taxi records and derive numeric UTC time features."""

import pandas


def add_time_features(df: pandas.DataFrame) -> pandas.DataFrame:
    """Remove missing required values and zero locations without mutating input."""
    df = df.dropna(
        subset=[
            "fare_amount",
            "pickup_datetime",
            "pickup_longitude",
            "pickup_latitude",
            "dropoff_longitude",
            "dropoff_latitude",
            "passenger_count",
        ]
    ).copy()
    invalid_pickup: pandas.Series = (df["pickup_longitude"] == 0) & (
        df["pickup_latitude"] == 0
    )
    invalid_dropoff: pandas.Series = (df["dropoff_longitude"] == 0) & (
        df["dropoff_latitude"] == 0
    )
    df = df.loc[~(invalid_pickup | invalid_dropoff)].copy()
    df["pickup_datetime"] = pandas.to_datetime(df["pickup_datetime"], utc=True)
    df["hour"] = df["pickup_datetime"].dt.hour
    df["day_of_week"] = df["pickup_datetime"].dt.dayofweek
    return df.drop(columns=["pickup_datetime"])

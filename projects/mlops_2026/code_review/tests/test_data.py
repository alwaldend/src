import pandas as pd
from src.data import load_data

def test_load_data():
    df = load_data("data/uber.csv")
    assert not df.empty

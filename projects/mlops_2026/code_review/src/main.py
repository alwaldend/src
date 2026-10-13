import os
import pandas as pd
from src.data import load_data, split_data
from src.features import add_time_features
from src.model import TaxiFareModel

# Глобальная переменная
DATA_PATH = "data/main.csv"

# Загрузка и обработка данных
raw_data = load_data(DATA_PATH)
processed_data = add_time_features(raw_data)
X_train, X_test, y_train, y_test = split_data(processed_data)

# Обучение модели
model = TaxiFareModel()
model.fit(X_train, y_train)

# Оценка
score = model.model.score(X_test, y_test)
print(f"R²: {score:.2f}")
import os
import glob
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt


folder = "B:/"
csv_files = glob.glob(os.path.join(folder, "*.csv"))
print(f"📄 Найдено файлов: {len(csv_files)}")

# Параметры
SAMPLE_RATE_HZ = 10_000_000  # 10 МГц, должен соответствовать Go-программе

def classify_file(file_name):
    file_name = file_name.lower()
    if "afc" in file_name:
        return "afc"
    if "fft" in file_name or "spectrum" in file_name:
        return "spectrum"
    if "envelope" in file_name:
        return "envelope"
    if "saft" in file_name:
        return "saft"
    if "time" in file_name or "filtered" in file_name:
        return "time"
    return "unknown"

def load_and_plot(file_path):
    file_name = os.path.basename(file_path)
    signal_type = classify_file(file_name)

    try:
        if signal_type in ["afc", "spectrum"]:
            df = pd.read_csv(file_path, header=None, names=["FrequencyHz", "Amplitude"])
            df = df.dropna()

            plt.figure(figsize=(12, 4))
            plt.plot(df["FrequencyHz"], df["Amplitude"], label=file_name)
            plt.title(f"{signal_type.upper()} from {file_name}")
            plt.xlabel("Frequency (Hz)")
            plt.ylabel("Amplitude")
            plt.grid(True)
            plt.legend()
            plt.tight_layout()
            plt.show()

        elif signal_type in ["time", "envelope", "saft"]:
            df = pd.read_csv(file_path, header=None, names=["value"])
            df = df.dropna()

            num_samples = len(df)
            time_seconds = np.arange(num_samples) / SAMPLE_RATE_HZ

            plt.figure(figsize=(12, 4))
            plt.plot(time_seconds, df["value"], label=file_name)
            plt.title(f"{signal_type.upper()} signal from {file_name}")
            plt.xlabel("Time (s)")
            plt.ylabel("Value")
            plt.grid(True)
            plt.legend()
            plt.tight_layout()
            plt.show()

        else:
            print(f"⚠️ Пропущен неизвестный тип файла: {file_name}")

    except Exception as e:
        print(f"❌ Ошибка в файле {file_name}: {e}")

# Основной вызов
for file in sorted(csv_files):
    load_and_plot(file)

# Табличный просмотр первых строк
for file in sorted(csv_files):
    print(f"\n📄 Файл: {os.path.basename(file)}")
    try:
        df = pd.read_csv(file, sep=',', header=None)
        display(df.head())
    except Exception as e:
        print(f"❌ Ошибка при чтении: {e}")

import sys
import os
import pandas as pd
import matplotlib.pyplot as plt
from matplotlib.backends.backend_qt5agg import FigureCanvasQTAgg as FigureCanvas
from matplotlib.backends.backend_qt5agg import NavigationToolbar2QT as NavigationToolbar
from PyQt5.QtWidgets import (QApplication, QMainWindow, QWidget, QVBoxLayout,
                             QPushButton, QFileDialog, QHBoxLayout, QLabel,
                             QCheckBox, QSlider, QDoubleSpinBox, QTabWidget)
from PyQt5.QtCore import QTimer, Qt
import threading
import time
import numpy as np
import hashlib
from collections import deque
import queue
import logging

# Настройка логирования
logging.basicConfig(level=logging.DEBUG, format='%(asctime)s - %(levelname)s - %(message)s')

class CustomNavigationToolbar(NavigationToolbar):
    def __init__(self, canvas, parent, max_time_callback):
        super().__init__(canvas, parent)
        self.max_time_callback = max_time_callback

    def zoom(self, *args):
        super().zoom(*args)
        self._restrict_axes()

    def pan(self, *args):
        super().pan(*args)
        self._restrict_axes()

    def _restrict_axes(self):
        ax = self.canvas.figure.gca()
        ax.set_ylim(0, ax.get_ylim()[1])
        max_time = self.max_time_callback()
        current_xlim = ax.get_xlim()
        ax.set_xlim(max(0, current_xlim[0]), min(max_time, current_xlim[1]))
        self.canvas.draw()

class OscilloscopeGUI(QMainWindow):
    def __init__(self):
        super().__init__()
        logging.debug("Инициализация OscilloscopeGUI")
        self.setWindowTitle("Осциллограф")
        self.setGeometry(100, 100, 1200, 800)

        # Структуры данных
        self.data = {}
        self.history = {}
        self.lines = {}
        self.history_lines = {}
        self.colors = ['b', 'r', 'g', 'c', 'm', 'y', 'k']
        self.max_points = 10000
        self.history_max = 100000
        self.data_hashes = {}
        self.data_queue = queue.Queue()
        self.files = []

        # Переменные состояния
        self.running = False
        self.follow_realtime = True
        self.paused = False
        self.read_realtime = False
        self.max_time = 1.0
        self.window_size = 0.1  # Установлено значение по умолчанию
        self.x_position = 0
        self.trigger_level = 0
        self.time_step = 0.001
        self.data_lock = threading.Lock()

        # Настройка UI
        self.main_widget = QWidget()
        self.setCentralWidget(self.main_widget)
        self.main_layout = QVBoxLayout(self.main_widget)

        # Вкладки
        self.tabs = QTabWidget()
        self.live_widget = QWidget()
        self.history_widget = QWidget()
        self.tabs.addTab(self.live_widget, "Живой просмотр")
        self.tabs.addTab(self.history_widget, "История")
        self.main_layout.addWidget(self.tabs)

        # Панель управления
        self.control_panel = QWidget()
        self.control_layout = QHBoxLayout(self.control_panel)

        self.load_button = QPushButton("Загрузить CSV")
        self.load_button.clicked.connect(self.load_csv_files)
        self.control_layout.addWidget(self.load_button)

        self.follow_realtime_cb = QCheckBox("Отслеживать в реальном времени")
        self.follow_realtime_cb.setChecked(True)
        self.follow_realtime_cb.stateChanged.connect(self.toggle_realtime)
        self.control_layout.addWidget(self.follow_realtime_cb)

        self.realtime_read_cb = QCheckBox("Чтение CSV в реальном времени")
        self.realtime_read_cb.setChecked(False)
        self.realtime_read_cb.stateChanged.connect(self.toggle_realtime_read)
        self.control_layout.addWidget(self.realtime_read_cb)

        self.pause_button = QPushButton("Пауза")
        self.pause_button.clicked.connect(self.toggle_pause)
        self.control_layout.addWidget(self.pause_button)

        self.save_plot_button = QPushButton("Сохранить график")
        self.save_plot_button.clicked.connect(self.save_plot)
        self.control_layout.addWidget(self.save_plot_button)

        self.status_label = QLabel("Файлы не загружены")
        self.control_layout.addWidget(self.status_label)

        # Живой график
        self.live_layout = QVBoxLayout(self.live_widget)
        self.live_figure, self.live_ax = plt.subplots()
        self.live_canvas = FigureCanvas(self.live_figure)
        self.live_canvas.setMinimumSize(800, 400)
        self.live_toolbar = CustomNavigationToolbar(self.live_canvas, self, self.get_max_time)

        # График истории
        self.history_layout = QVBoxLayout(self.history_widget)
        self.history_figure, self.history_ax = plt.subplots()
        self.history_canvas = FigureCanvas(self.history_figure)
        self.history_canvas.setMinimumSize(800, 400)
        self.history_toolbar = CustomNavigationToolbar(self.history_canvas, self, self.get_max_time)

        # Элементы управления под графиком
        self.controls_below = QWidget()
        self.controls_layout = QHBoxLayout(self.controls_below)

        self.x_scale_spin = QDoubleSpinBox()
        self.x_scale_spin.setMinimum(0.001)
        self.x_scale_spin.setMaximum(1000)
        self.x_scale_spin.setValue(0.1)
        self.x_scale_spin.valueChanged.connect(self.update_x_scale)
        self.controls_layout.addWidget(QLabel("Масштаб X:"))
        self.controls_layout.addWidget(self.x_scale_spin)

        self.x_pos_slider = QSlider(Qt.Horizontal)
        self.x_pos_slider.setMinimum(0)
        self.x_pos_slider.setMaximum(1000)
        self.x_pos_slider.valueChanged.connect(self.update_x_position)
        self.controls_layout.addWidget(QLabel("Позиция X:"))
        self.controls_layout.addWidget(self.x_pos_slider)

        self.time_step_spin = QDoubleSpinBox()
        self.time_step_spin.setMinimum(0.001)
        self.time_step_spin.setMaximum(1000)
        self.time_step_spin.setValue(1)
        self.time_step_spin.valueChanged.connect(self.update_time_step)
        self.controls_layout.addWidget(QLabel("Шаг времени (мс):"))
        self.controls_layout.addWidget(self.time_step_spin)

        self.trigger_spin = QDoubleSpinBox()
        self.trigger_spin.setMinimum(-1e6)
        self.trigger_spin.setMaximum(1e6)
        self.trigger_spin.setValue(0)
        self.trigger_spin.valueChanged.connect(self.update_trigger)
        self.controls_layout.addWidget(QLabel("Триггер:"))
        self.controls_layout.addWidget(self.trigger_spin)

        # Сборка компоновок
        self.live_layout.addWidget(self.control_panel)
        self.live_layout.addWidget(self.live_toolbar)
        self.live_layout.addWidget(self.live_canvas)
        self.live_layout.addWidget(self.controls_below)
        self.history_layout.addWidget(self.history_toolbar)
        self.history_layout.addWidget(self.history_canvas)

        # Таймер
        self.update_timer = QTimer()
        self.update_timer.timeout.connect(self.refresh_plots)
        self.update_timer.start(200)

        logging.debug("Инициализация GUI завершена")

    def get_max_time(self):
        return self.max_time if self.max_time > 0 else 1

    def toggle_realtime(self, state):
        logging.debug(f"Переключение реального времени: {state}")
        self.follow_realtime = state == Qt.Checked
        if self.follow_realtime:
            self.x_pos_slider.setValue(self.x_pos_slider.maximum())
        self.refresh_plots()

    def toggle_realtime_read(self, state):
        logging.debug(f"Переключение чтения в реальном времени: {state}")
        self.read_realtime = state == Qt.Checked
        if self.read_realtime and self.files:
            self.start_file_monitors()

    def toggle_pause(self):
        logging.debug(f"Переключение паузы: {not self.paused}")
        self.paused = not self.paused
        self.pause_button.setText("Возобновить" if self.paused else "Пауза")

    def update_x_scale(self, value):
        logging.debug(f"Обновление масштаба X: {value}")
        self.window_size = max(0.001, min(value, self.max_time))  # Предотвращаем нулевой масштаб
        self.x_position = min(self.x_position, self.max_time - self.window_size)
        self.x_scale_spin.setValue(self.window_size)  # Синхронизируем значение
        self.refresh_plots()

    def update_x_position(self, value):
        logging.debug(f"Обновление позиции X: {value}")
        max_pos = max(0, self.max_time - self.window_size)
        self.x_position = (value / 1000) * max_pos
        self.refresh_plots()

    def update_time_step(self, value):
        logging.debug(f"Обновление шага времени: {value}")
        self.time_step = value / 1000
        with self.data_lock:
            for file_path in self.data:
                df = self.data[file_path]
                df['rel_time'] = np.arange(len(df)) * self.time_step
        self.refresh_plots()

    def update_trigger(self, value):
        logging.debug(f"Обновление триггера: {value}")
        self.trigger_level = value
        self.refresh_plots()

    def save_plot(self):
        logging.debug("Сохранение графика")
        file_name, _ = QFileDialog.getSaveFileName(self, "Сохранить график", "", "PNG (*.png)")
        if file_name:
            current_tab = self.tabs.currentWidget()
            figure = self.live_figure if current_tab == self.live_widget else self.history_figure
            figure.savefig(file_name)
            self.status_label.setText(f"График сохранен: {os.path.basename(file_name)}")

    def load_csv_files(self):
        logging.debug("Загрузка CSV файлов")
        files, _ = QFileDialog.getOpenFileNames(self, "Выберите CSV-файлы", "", "CSV-файлы (*.csv)")
        if files:
            self.files = files
            self.status_label.setText(f"Загружено {len(self.files)} файл(ов)")
            self.running = True
            self.data_hashes = {file: "" for file in self.files}
            self.history = {file: deque(maxlen=self.history_max) for file in self.files}
            self.initialize_plots()
            if self.read_realtime:
                self.start_file_monitors()

    def initialize_plots(self):
        logging.debug("Инициализация графиков")
        self.live_ax.clear()
        self.history_ax.clear()
        self.data = {}
        self.lines = {}
        self.history_lines = {}

        for idx, file in enumerate(self.files):
            try:
                df = pd.read_csv(file, nrows=1000, encoding='utf-8', on_bad_lines='skip')
                if df.empty:
                    logging.warning(f"Файл {file} пустой")
                    continue
                df['timestamp'] = pd.to_datetime(df['timestamp'], errors='coerce')
                if df['timestamp'].nunique() <= 1:
                    df['rel_time'] = np.arange(len(df)) * self.time_step
                else:
                    df['rel_time'] = (df['timestamp'] - df['timestamp'].iloc[0]).dt.total_seconds()

                self.data[file] = df
                color = self.colors[idx % len(self.colors)]

                line, = self.live_ax.plot([], [], color=color, label=os.path.basename(file))
                self.lines[file] = line

                h_line, = self.history_ax.plot([], [], color=color, label=os.path.basename(file))
                self.history_lines[file] = h_line

                self.history[file].extend(zip(df['rel_time'], df['value']))

            except Exception as e:
                self.status_label.setText(f"Ошибка чтения {file}: {str(e)}")
                logging.error(f"Ошибка чтения {file}: {str(e)}")

        self.format_axes()
        self.refresh_plots()

    def start_file_monitors(self):
        logging.debug("Запуск мониторинга файлов")
        for file in self.files:
            thread = threading.Thread(target=self.monitor_file, args=(file,))
            thread.daemon = True
            thread.start()

    def monitor_file(self, file_path):
        logging.debug(f"Мониторинг файла: {file_path}")
        last_pos = 0
        while self.running and file_path in self.files and self.read_realtime:
            try:
                with open(file_path, 'r', encoding='utf-8') as f:
                    f.seek(last_pos)
                    new_data = f.read()
                    if new_data:
                        df = pd.read_csv(pd.StringIO(new_data), on_bad_lines='skip')
                        if not df.empty:
                            df['timestamp'] = pd.to_datetime(df['timestamp'], errors='coerce')
                            if df['timestamp'].nunique() <= 1:
                                df['rel_time'] = np.arange(len(df)) * self.time_step + (len(self.data.get(file_path, [])) * self.time_step)
                            else:
                                df['rel_time'] = (df['timestamp'] - df['timestamp'].iloc[0]).dt.total_seconds()

                            data_hash = hashlib.md5(df.to_string().encode()).hexdigest()
                            if data_hash != self.data_hashes.get(file_path):
                                self.data_queue.put((file_path, df, data_hash))
                            last_pos = f.tell()
                time.sleep(0.2)
            except Exception as e:
                self.status_label.setText(f"Ошибка обновления {file_path}: {str(e)}")
                logging.error(f"Ошибка обновления {file_path}: {str(e)}")
                time.sleep(1)

    def refresh_plots(self):
        if self.paused:
            return

        logging.debug("Обновление графиков")
        if not self.data and not self.history:
            logging.debug("Нет данных для отображения")
            return

        start_time = time.time()
        while not self.data_queue.empty() and (time.time() - start_time) < 0.1:
            try:
                file_path, df, data_hash = self.data_queue.get_nowait()
                with self.data_lock:
                    if file_path in self.data:
                        self.data[file_path] = pd.concat([self.data[file_path], df]).drop_duplicates().reset_index(drop=True)
                    else:
                        self.data[file_path] = df
                    self.data_hashes[file_path] = data_hash
                    self.history[file_path].extend(zip(df['rel_time'], df['value']))
            except queue.Empty:
                break

        with self.data_lock:
            self.max_time = 0
            max_value = 0

            for file_path in self.data:
                if file_path in self.lines:
                    df = self.data[file_path].copy()
                    if df.empty:
                        continue
                    if len(df) > self.max_points:
                        step = len(df) // self.max_points
                        df = df.iloc[::step]

                    if self.trigger_level != 0:
                        trigger_idx = df[df['value'] >= self.trigger_level].index
                        if not trigger_idx.empty:
                            trigger_time = df.loc[trigger_idx[0], 'rel_time']
                            df['rel_time'] = df['rel_time'] - trigger_time

                    self.lines[file_path].set_data(df['rel_time'], df['value'])
                    self.max_time = max(self.max_time, df['rel_time'].max() if not df.empty else 0)
                    max_value = max(max_value, df['value'].max() if not df.empty else 0)

            for file_path in self.history:
                if file_path in self.history_lines:
                    times, values = zip(*self.history[file_path]) if self.history[file_path] else ([], [])
                    self.history_lines[file_path].set_data(times, values)
                    self.max_time = max(self.max_time, max(times, default=0))
                    max_value = max(max_value, max(values, default=0))

            if self.max_time <= 0:
                self.max_time = 1.0
            if self.window_size > self.max_time:
                self.window_size = self.max_time
                self.x_scale_spin.setValue(self.window_size)

            max_pos = max(0, self.max_time - self.window_size)
            if self.follow_realtime:
                self.x_position = max_pos
                self.x_pos_slider.setValue(1000)

            start_time = self.x_position
            end_time = min(self.x_position + self.window_size, self.max_time)

            self.live_ax.set_xlim(start_time, end_time)
            self.live_ax.set_ylim(0, max_value * 1.1 if max_value > 0 else 1)

            self.history_ax.set_xlim(0, self.max_time)
            self.history_ax.set_ylim(0, max_value * 1.1 if max_value > 0 else 1)

            self.format_axes()
            try:
                self.live_canvas.draw()
                self.history_canvas.draw()
            except Exception as e:
                logging.error(f"Ошибка отрисовки графиков: {str(e)}")

    def format_axes(self):
        for ax in [self.live_ax, self.history_ax]:
            ax.set_xlabel("Время (с)")
            ax.set_ylabel("Значение")
            ax.set_title("Данные осциллографа")
            ax.legend()
            ax.grid(True)

    def closeEvent(self, event):
        logging.debug("Закрытие приложения")
        self.running = False
        event.accept()

if __name__ == '__main__':
    logging.debug("Запуск приложения")
    app = QApplication(sys.argv)
    window = OscilloscopeGUI()
    window.show()
    sys.exit(app.exec_())
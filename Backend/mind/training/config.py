"""训练流水线配置"""
import os

# InfluxDB 3 连接参数
INFLUXDB_URL = os.getenv("INFLUXDB_URL", "http://localhost:8181")
INFLUXDB_TOKEN = os.getenv("INFLUXDB_TOKEN", "")
INFLUXDB_ORG = os.getenv("INFLUXDB_ORG", "my-org")
INFLUXDB_BUCKET = os.getenv("INFLUXDB_BUCKET", "sensor_data")

# 训练参数 — 必须与 Go 侧 DeltaT 一致
DELTA_T = 300  # 5分钟 = 300秒
STATE_DIM = 9
CONTROL_DIM = 5

# 状态列名（与 sensor 服务 InfluxDB measurement 对齐）
STATE_COLS = [
    "temperature", "ph", "oxygen", "ammonia", "water_level",
    "tds", "nitrate", "nitrite", "chloride",
]

# 控制列名（当 actuator 反馈未接入时用零填充）
CONTROL_COLS = ["heater_power", "aerator_power", "pump_power", "feeder_dispense", "light_intensity"]

# Trig 频率参数（必须与 Go 侧 ObservableConfig 一致，2个频率）
TRIG_FREQS = [0.5, 1.0]

# 模型输出路径
MODEL_OUTPUT_PATH = os.getenv("MODEL_OUTPUT_PATH", "../models/koopman_model.json")

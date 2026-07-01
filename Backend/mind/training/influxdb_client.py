"""从 InfluxDB 3 获取历史传感器数据用于模型训练"""
from influxdb_client_3 import InfluxDBClient3
import pandas as pd
import numpy as np

from config import (
    INFLUXDB_URL, INFLUXDB_TOKEN, INFLUXDB_ORG, INFLUXDB_BUCKET,
    STATE_COLS, DELTA_T,
)


def fetch_training_data(tank_id: str, hours: int = 720) -> pd.DataFrame:
    """获取指定鱼缸最近 N 小时的传感器数据。

    返回按时间排序的 DataFrame，列为 STATE_COLS。
    数据会被重采样到 DELTA_T 秒间隔并线性插值。
    """
    client = InfluxDBClient3(
        host=INFLUXDB_URL,
        token=INFLUXDB_TOKEN,
        org=INFLUXDB_ORG,
        database=INFLUXDB_BUCKET,
    )

    cols = ", ".join(STATE_COLS)
    query = f"""
    SELECT {cols}
    FROM sensor_data
    WHERE tank_id = '{tank_id}'
    AND time >= now() - interval '{hours} hours'
    ORDER BY time ASC
    """

    df = client.query(query, mode="pandas")

    # 确保规则时间间隔
    if df is not None and not df.empty:
        df = df.set_index("time")
        df = df.resample(f"{DELTA_T}s").mean().interpolate(method="linear")
        df = df.dropna()

    return df.reset_index() if df is not None and not df.empty else pd.DataFrame(columns=["time"] + STATE_COLS)


def fetch_control_data(tank_id: str, hours: int = 720) -> pd.DataFrame:
    """获取控制动作历史。

    如果 actuator 数据尚未写入 InfluxDB（初始部署的常见情况），
    返回全零 DataFrame，表示 EDMD 将学习自治 Koopman 算子。
    后续可扩展为从 actuator measurement 读取。
    """
    state_df = fetch_training_data(tank_id, hours)
    n = len(state_df)
    control_data = np.zeros((n, len(CONTROL_COLS := ["heater_power", "aerator_power", "pump_power", "feeder_dispense", "light_intensity"])))
    return pd.DataFrame(control_data, columns=CONTROL_COLS)

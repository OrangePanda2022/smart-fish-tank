#!/usr/bin/env python3
"""
AquaTank IoT MQTT 模拟器
向 NanoMQ 发布模拟传感器数据，供 sensor 服务消费。

用法:
    uv run mqtt_simulator.py --tank-id <UUID> --count 10
    uv run mqtt_simulator.py --broker 192.168.1.100 --port 1883 --tank-id tank-001 --count 20 --interval 3
"""

import json
import time
import argparse
import random
from datetime import datetime, timezone

import paho.mqtt.client as mqtt


def generate_sensor_data(tank_id: str, device_id: str, overrides: dict | None = None) -> dict:
    """生成模拟鱼缸传感器读数

    overrides: 可选字典，指定精确传感器值覆盖随机生成
    例如 {"oxygen": 4.0, "temperature": 20.0}
    """
    base = {
        "device_id": device_id,
        "tank_id": tank_id,
        "temperature": round(random.uniform(22.0, 29.0), 1),
        "ph": round(random.uniform(6.5, 8.0), 2),
        "oxygen": round(random.uniform(4.0, 9.0), 1),
        "ammonia": round(random.uniform(0.0, 1.5), 2),
        "water_level": round(random.uniform(60.0, 95.0), 1),
        "tds": round(random.uniform(100, 400), 0),
        "nitrate": round(random.uniform(0, 30), 1),
        "nitrite": round(random.uniform(0, 0.3), 2),
        "chloride": round(random.uniform(10, 150), 1),
        "timestamp": datetime.now(timezone.utc).isoformat(),
    }
    if overrides:
        for key, val in overrides.items():
            if key in base:
                base[key] = val
    return base


def main():
    parser = argparse.ArgumentParser(description="AquaTank MQTT 模拟器")
    parser.add_argument("--broker", default="localhost", help="MQTT broker 地址 (默认: localhost)")
    parser.add_argument("--port", type=int, default=1883, help="MQTT broker 端口 (默认: 1883)")
    parser.add_argument("--tank-id", default="019f1807-dbca-753d-afb8-6b2db160fa8d", help="目标鱼缸 ID (默认: 019f1807-dbca-753d-afb8-6b2db160fa8d)")
    parser.add_argument("--device-id", default="sensor-sim-001", help="模拟设备 ID (默认: sensor-sim-001)")
    parser.add_argument("--count", type=int, default=10, help="发送消息数量 (默认: 10)")
    parser.add_argument("--interval", type=float, default=2.0, help="消息间隔秒数 (默认: 2.0)")
    parser.add_argument("--oxygen", type=float, default=None, help="指定溶解氧值 (mg/L)")
    parser.add_argument("--temperature", type=float, default=None, help="指定温度值 (°C)")
    parser.add_argument("--ph", type=float, default=None, help="指定pH值")
    parser.add_argument("--ammonia", type=float, default=None, help="指定氨氮值 (mg/L)")
    parser.add_argument("--water-level", type=float, default=None, help="指定水位值 (%)")
    parser.add_argument("--tds", type=float, default=None, help="指定TDS值 (ppm)")
    parser.add_argument("--nitrate", type=float, default=None, help="指定硝酸根值 (mg/L)")
    parser.add_argument("--nitrite", type=float, default=None, help="指定亚硝酸盐值 (mg/L)")
    parser.add_argument("--chloride", type=float, default=None, help="指定氯离子值 (mg/L)")
    args = parser.parse_args()

    # 构建传感器值覆盖
    overrides = {}
    if args.oxygen is not None:
        overrides["oxygen"] = args.oxygen
    if args.temperature is not None:
        overrides["temperature"] = args.temperature
    if args.ph is not None:
        overrides["ph"] = args.ph
    if args.ammonia is not None:
        overrides["ammonia"] = args.ammonia
    if args.water_level is not None:
        overrides["water_level"] = args.water_level
    if args.tds is not None:
        overrides["tds"] = args.tds
    if args.nitrate is not None:
        overrides["nitrate"] = args.nitrate
    if args.nitrite is not None:
        overrides["nitrite"] = args.nitrite
    if args.chloride is not None:
        overrides["chloride"] = args.chloride

    # paho-mqtt v2 API: 使用 CallbackAPIVersion.VERSION2
    client = mqtt.Client(
        callback_api_version=mqtt.CallbackAPIVersion.VERSION2,
        client_id="aqua-simulator",
    )

    def on_connect(client, userdata, flags, rc, properties):
        if rc == 0:
            print(f"已连接到 MQTT broker {args.broker}:{args.port}")
        else:
            print(f"连接失败，返回码: {rc}")

    def on_disconnect(client, userdata, flags, rc, properties):
        if rc != 0:
            print(f"意外断开连接，返回码: {rc}")

    client.on_connect = on_connect
    client.on_disconnect = on_disconnect

    try:
        client.connect(args.broker, args.port, 60)
    except Exception as e:
        print(f"无法连接到 {args.broker}:{args.port} — {e}")
        return

    client.loop_start()

    # sensor 服务订阅的是 sensors/+/data，+ 是单层通配符
    topic = f"sensors/{args.device_id}/data"

    for i in range(args.count):
        payload = generate_sensor_data(args.tank_id, args.device_id, overrides if overrides else None)
        result = client.publish(topic, json.dumps(payload), qos=1)
        if result.rc == mqtt.MQTT_ERR_SUCCESS:
            print(f"[{i + 1}/{args.count}] 发送到 {topic}: {payload}")
        else:
            print(f"[{i + 1}/{args.count}] 发送失败: rc={result.rc}")
        time.sleep(args.interval)

    client.loop_stop()
    client.disconnect()
    print(f"模拟完成，共发送 {args.count} 条消息。")


if __name__ == "__main__":
    main()

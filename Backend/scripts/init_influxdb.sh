#!/bin/bash
# 初始化 InfluxDB 3 Core 数据库
# 在 docker compose up 之后运行

INFLUX_URL="http://localhost:8181"
DB_NAME="sensor_data"

echo "创建数据库: $DB_NAME"
curl -s -X POST "${INFLUX_URL}/api/v3/databases" \
  -H "Content-Type: application/json" \
  -d "{\"name\": \"${DB_NAME}\"}"

echo ""
echo "验证数据库..."
curl -s "${INFLUX_URL}/api/v3/databases"
echo ""
echo "InfluxDB 3 初始化完成。"

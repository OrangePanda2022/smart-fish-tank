# Sensor Service API 文档

## 基础信息

- **Base URL**: `/api/v1`
- **框架**: Gin
- **认证方式**: 无 (内部服务)

---

## 接口列表

### 1. 健康检查

**端点**: `GET /health`

**认证**: 否

**请求示例** (JavaScript):
```javascript
const response = await fetch('/health');
const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "status": "healthy",
  "service": "sensor-service"
}
```

---

### 2. 获取鱼缸最新传感器数据

**端点**: `GET /api/v1/tanks/:tankId/sensors/latest`

**认证**: 否 (内部服务调用)

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tankId | string | 是 | 鱼缸 ID |

**请求示例** (JavaScript):
```javascript
const tankId = 'tank-uuid-123';
const response = await fetch(`/api/v1/tanks/${tankId}/sensors/latest`);
const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "deviceId": "device-uuid-456",
  "tankId": "tank-uuid-123",
  "temperature": 25.5,
  "ph": 7.2,
  "oxygen": 8.5,
  "ammonia": 0.1,
  "waterLevel": 80.0,
  "timestamp": "2024-01-01T12:00:00Z"
}
```

**错误响应** (400):
```json
{
  "error": "tankId不能为空"
}
```

**错误响应** (404):
```json
{
  "error": "未找到数据"
}
```

---

### 3. 获取鱼缸传感器历史数据

**端点**: `GET /api/v1/tanks/:tankId/sensors/history`

**认证**: 否 (内部服务调用)

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tankId | string | 是 | 鱼缸 ID |

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| start | string | 否 | 开始时间 (ISO 8601 格式)，如 `2024-01-01T00:00:00Z` |
| end | string | 否 | 结束时间 (ISO 8601 格式)，如 `2024-01-02T00:00:00Z` |
| limit | integer | 否 | 返回数据条数限制，默认 100 |

**请求示例** (JavaScript):
```javascript
const tankId = 'tank-uuid-123';
const response = await fetch(
  `/api/v1/tanks/${tankId}/sensors/history?start=2024-01-01T00:00:00Z&end=2024-01-02T00:00:00Z&limit=50`
);
const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "data": [
    {
      "deviceId": "device-uuid-456",
      "tankId": "tank-uuid-123",
      "temperature": 25.5,
      "ph": 7.2,
      "oxygen": 8.5,
      "ammonia": 0.1,
      "waterLevel": 80.0,
      "timestamp": "2024-01-01T12:00:00Z"
    },
    {
      "deviceId": "device-uuid-456",
      "tankId": "tank-uuid-123",
      "temperature": 25.8,
      "ph": 7.1,
      "oxygen": 8.4,
      "ammonia": 0.15,
      "waterLevel": 79.5,
      "timestamp": "2024-01-01T12:05:00Z"
    }
  ],
  "count": 2
}
```

**字段说明**:
| 字段 | 类型 | 说明 |
|------|------|------|
| data | array | 传感器数据列表 |
| data[].deviceId | string | 设备 ID |
| data[].tankId | string | 鱼缸 ID |
| data[].temperature | float | 温度 (摄氏度) |
| data[].ph | float | pH 值 |
| data[].oxygen | float | 溶解氧 (mg/L) |
| data[].ammonia | float | 氨氮含量 (mg/L) |
| data[].waterLevel | float | 水位 (百分比) |
| data[].timestamp | string | 数据采集时间 (ISO 8601) |
| count | integer | 返回数据条数 |

**错误响应** (400):
```json
{
  "error": "tankId不能为空"
}
```

---

### 4. 获取设备传感器数据

**端点**: `GET /api/v1/devices/:deviceId/sensors`

**认证**: 否 (内部服务调用)

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| deviceId | string | 是 | 设备 ID |

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| limit | integer | 否 | 返回数据条数限制，默认 100 |

**请求示例** (JavaScript):
```javascript
const deviceId = 'device-uuid-456';
const response = await fetch(`/api/v1/devices/${deviceId}/sensors?limit=50`);
const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "data": [
    {
      "deviceId": "device-uuid-456",
      "tankId": "tank-uuid-123",
      "temperature": 25.5,
      "ph": 7.2,
      "oxygen": 8.5,
      "ammonia": 0.1,
      "waterLevel": 80.0,
      "timestamp": "2024-01-01T12:00:00Z"
    }
  ],
  "count": 1
}
```

**错误响应** (400):
```json
{
  "error": "deviceId不能为空"
}
```

---

## 通用响应格式

### 成功响应 (单条数据)
```json
{
  "deviceId": "string",
  "tankId": "string",
  "temperature": 25.5,
  "ph": 7.2,
  "oxygen": 8.5,
  "ammonia": 0.1,
  "waterLevel": 80.0,
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### 成功响应 (列表数据)
```json
{
  "data": [...],
  "count": 10
}
```

### 错误响应
```json
{
  "error": "错误描述"
}
```

---

## 数据模型

### SensorDataDTO (传感器数据)
```json
{
  "deviceId": "device-uuid",
  "tankId": "tank-uuid",
  "temperature": 25.5,
  "ph": 7.2,
  "oxygen": 8.5,
  "ammonia": 0.1,
  "waterLevel": 80.0,
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### SensorDataListDTO (传感器数据列表)
```json
{
  "data": [SensorDataDTO, ...],
  "count": 10
}
```

---

## 传感器指标说明

| 指标 | 类型 | 单位 | 正常范围 | 说明 |
|------|------|------|----------|------|
| temperature | float | °C | 22-28 | 水温 |
| ph | float | - | 6.5-8.0 | 酸碱度 |
| oxygen | float | mg/L | 5-9 | 溶解氧 |
| ammonia | float | mg/L | 0-0.5 | 氨氮 |
| waterLevel | float | % | 60-90 | 水位 |

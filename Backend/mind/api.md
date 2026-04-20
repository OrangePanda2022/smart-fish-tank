# Mind Service API 文档

## 基础信息

- **Base URL**: `/api/v1`
- **框架**: Hertz (字节跳动开源框架)
- **认证方式**: 无 (内部服务调用)

---

## 接口列表

### 1. 健康检查

**端点**: `GET /ping`

**认证**: 否

**请求示例** (JavaScript):
```javascript
const response = await fetch('/ping');
const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "message": "pong"
}
```

---

### 2. 分析鱼缸数据

**端点**: `GET /api/v1/analyse/:tank_id`

**认证**: 否 (内部服务调用)

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tank_id | string | 是 | 鱼缸 ID |

**请求体** (可选):
```json
{
  "data": [
    {
      "device_id": "device-uuid",
      "tank_id": "tank-uuid",
      "temperature": 25.5,
      "ph": 7.2,
      "oxygen": 8.5,
      "ammonia": 0.1,
      "water_level": 80.0,
      "timestamp": "2024-01-01T12:00:00Z"
    }
  ]
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| data | array | 否 | 传感器数据列表。如果为空，将使用数据库中存储的历史数据进行分析。 |
| data[].device_id | string | 否 | 设备 ID |
| data[].tank_id | string | 否 | 鱼缸 ID |
| data[].temperature | float | 否 | 温度 |
| data[].ph | float | 否 | pH 值 |
| data[].oxygen | float | 否 | 溶解氧 |
| data[].ammonia | float | 否 | 氨氮含量 |
| data[].water_level | float | 否 | 水位 |
| data[].timestamp | string | 否 | 时间戳 |

**请求示例** (JavaScript):
```javascript
const tankId = 'tank-uuid-123';
const response = await fetch(`/api/v1/analyse/${tankId}`, {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    data: [
      {
        device_id: 'device-uuid',
        tank_id: 'tank-uuid',
        temperature: 25.5,
        ph: 7.2,
        oxygen: 8.5,
        ammonia: 0.1,
        water_level: 80.0,
        timestamp: '2024-01-01T12:00:00Z'
      }
    ]
  })
});

const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "report_id": "report-uuid-456",
    "tank_id": "tank-uuid-123",
    "decision": {
      "status_score": 85,
      "summary": "水质状况良好，各项指标均在正常范围内",
      "actions": [
        {
          "device": " feeder",
          "action": "feed"
        },
        {
          "device": "light",
          "action": "dim"
        }
      ],
      "reasoning": "温度和pH值稳定，溶解氧充足，氨氮含量在安全范围内。建议正常喂食，保持当前光照。"
    }
  }
}
```

**错误响应** (400):
```json
{
  "code": 400,
  "message": "tank_id is required"
}
```

**错误响应** (500):
```json
{
  "code": 500,
  "message": "analyse failed: 错误详情"
}
```

---

## 响应字段说明

### AnalyseResponse
```json
{
  "report_id": "string",
  "tank_id": "string",
  "decision": Decision
}
```

### Decision
```json
{
  "status_score": 85,
  "summary": "string",
  "actions": [Action],
  "reasoning": "string"
}
```

### Action
```json
{
  "device": "string",
  "action": "string"
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| report_id | string | 分析报告 ID |
| tank_id | string | 鱼缸 ID |
| decision.status_score | integer | 状态评分 (0-100) |
| decision.summary | string | 状态总结 |
| decision.actions | array | 建议执行的动作列表 |
| decision.actions[].device | string | 设备类型 |
| decision.actions[].action | string | 建议动作 |
| decision.reasoning | string | 分析推理过程 |

---

## 设备类型

| device | 说明 |
|--------|------|
| feeder | 投喂器 |
| light | 灯光 |
| heater | 加热器 |
| pump | 水泵 |
| aerator | 增氧机 |

---

## 建议动作

| action | 说明 |
|--------|------|
| feed | 投喂 |
| dim | 调暗灯光 |
| brighten | 调亮灯光 |
| heat | 加热 |
| cool | 降温 |
| increase_aeration | 增加曝气 |
| change_water | 换水 |
| add_oxygen | 加氧 |
| maintain | 保持现状 |

---

## 注意事项

1. **框架差异**: Mind 服务使用 Hertz 框架，而其他服务使用 Gin。考虑统一框架以降低维护成本。

2. **数据来源**: 如果请求体中不提供 `data` 字段，服务将从数据库获取历史传感器数据进行分析。

3. **分析逻辑**: 具体的分析算法和决策逻辑在 `mind/internal/service/analyse.go` 中实现。

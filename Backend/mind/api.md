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

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| stream | string | 否 | 设为 `true` 时启用 SSE 流式响应 |

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
          "device": "feeder",
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

### 3. 预测鱼缸未来状态

基于 Koopman 算子模型预测鱼缸水质参数的未来走势，并使用 MPC（模型预测控制）计算最优设备控制建议。只需传入 `tank_id`，系统自动获取最新传感器数据执行预测。

**端点**: `GET /api/v1/predict/:tank_id`

**认证**: 否 (内部服务调用)

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tank_id | string | 是 | 鱼缸 ID |

**查询参数**:
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| horizon | integer | 否 | 12 | 预测步数（1步=5分钟），范围 1-60，默认12步即1小时 |

**请求示例** (JavaScript):
```javascript
const tankId = 'tank-001';

// 默认预测1小时
const response = await fetch(`/api/v1/predict/${tankId}`);

// 自定义预测30分钟
const response30min = await fetch(`/api/v1/predict/${tankId}?horizon=6`);

const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "tank_id": "tank-001",
    "predictions": [
      {
        "state": [25.3, 7.1, 6.9, 0.28, 76.5],
        "time": "2026-06-30T12:05:00Z",
        "is_anomaly": false
      },
      {
        "state": [25.5, 7.1, 6.8, 0.30, 76.0],
        "time": "2026-06-30T12:10:00Z",
        "is_anomaly": false
      },
      {
        "state": [25.8, 7.0, 6.7, 0.35, 75.5],
        "time": "2026-06-30T12:15:00Z",
        "is_anomaly": false
      },
      {
        "state": [26.1, 6.9, 6.5, 0.42, 75.0],
        "time": "2026-06-30T12:20:00Z",
        "is_anomaly": false
      },
      {
        "state": [26.5, 6.8, 6.3, 0.51, 74.5],
        "time": "2026-06-30T12:25:00Z",
        "is_anomaly": true
      }
    ],
    "mpc_actions": [
      {
        "device": "heater",
        "action": "cool",
        "value": 0.18
      },
      {
        "device": "aerator",
        "action": "on",
        "value": 0.75
      },
      {
        "device": "pump",
        "action": "on",
        "value": 0.62
      },
      {
        "device": "feeder",
        "action": "maintain",
        "value": 0.45
      },
      {
        "device": "light",
        "action": "maintain",
        "value": 0.52
      }
    ],
    "mpc_predicted": [
      {
        "state": [25.2, 7.1, 7.0, 0.22, 76.8],
        "time": "2026-06-30T12:05:00Z",
        "is_anomaly": false
      }
    ],
    "anomaly_step": 4,
    "anomaly_warning": "第 5 步预测异常 (1500秒后)，建议立即采取行动",
    "confidence": 0.85,
    "generated_at": "2026-06-30T12:00:00Z"
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
  "message": "prediction failed: 获取当前状态失败: 鱼缸 tank-xxx 无传感器数据"
}
```

**错误响应** (503):
```json
{
  "code": 503,
  "message": "prediction service unavailable: model not loaded"
}
```

> **说明**: 当 Koopman 模型未加载时（例如初始部署、训练前或模型文件损坏），此端点返回 503。已有的 `/analyse` 端点不受影响，仍可正常使用。

---

### 4. 预测流式推送（SSE）

以 Server-Sent Events 方式持续推送鱼缸预测结果。客户端连接后立即收到一次初始预测，后续可周期性接收更新。

**端点**: `GET /api/v1/predict/:tank_id/stream`

**认证**: 否 (内部服务调用)

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tank_id | string | 是 | 鱼缸 ID |

**查询参数**:
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| horizon | integer | 否 | 12 | 预测步数，范围 1-60 |

**请求示例** (JavaScript):
```javascript
const tankId = 'tank-001';
const eventSource = new EventSource(`/api/v1/predict/${tankId}/stream?horizon=12`);

eventSource.onmessage = (event) => {
  const prediction = JSON.parse(event.data);
  console.log('新预测:', prediction);
};

eventSource.onerror = () => {
  console.error('SSE 连接断开');
  eventSource.close();
};
```

**SSE 响应格式**:
```
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive
Transfer-Encoding: chunked

data: {"tank_id":"tank-001","predictions":[...],"mpc_actions":[...],"mpc_predicted":[...],"anomaly_step":4,"anomaly_warning":"...","confidence":0.85,"generated_at":"2026-06-30T12:00:00Z"}

data: {"tank_id":"tank-001","predictions":[...],"mpc_actions":[...],"mpc_predicted":[...],"anomaly_step":-1,"anomaly_warning":"","confidence":0.88,"generated_at":"2026-06-30T12:00:30Z"}
```

**错误事件** (首次预测失败时):
```
event: error
data: prediction failed: 预测模型未加载
```

---

### 5. 模型状态查询

查询 Koopman 预测模型的加载状态和基本信息，用于健康检查和运维监控。

**端点**: `GET /api/v1/predict/status`

**认证**: 否 (内部服务调用)

**请求示例** (JavaScript):
```javascript
const response = await fetch('/api/v1/predict/status');
const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "model_loaded": true,
    "version": "1.0.0",
    "trained_at": "2026-06-30T10:00:00Z"
  }
}
```

**模型未加载时** (200):
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "model_loaded": false,
    "version": "",
    "trained_at": ""
  }
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

### PredictResponse
```json
{
  "tank_id": "string",
  "predictions": [PredictionPoint],
  "mpc_actions": [ControlAction],
  "mpc_predicted": [PredictionPoint],
  "anomaly_step": "integer",
  "anomaly_warning": "string",
  "confidence": "float",
  "generated_at": "string (ISO 8601)"
}
```

### PredictionPoint
```json
{
  "state": [25.5, 7.2, 7.0, 0.25, 75.0],
  "time": "string (ISO 8601)",
  "is_anomaly": "boolean"
}
```

### ControlAction
```json
{
  "device": "string",
  "action": "string",
  "value": "float (0-1)"
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| tank_id | string | 鱼缸 ID |
| predictions | array | 自由演化预测轨迹（不做任何控制干预下的未来状态序列） |
| predictions[].state | float[5] | 预测时刻的5维物理状态，顺序: [温度(°C), pH, 溶解氧(mg/L), 氨氮(mg/L), 水位(%)] |
| predictions[].time | string | 预测时刻的时间戳 (ISO 8601) |
| predictions[].is_anomaly | boolean | 该预测时刻是否有参数超出安全范围 |
| mpc_actions | array | MPC 计算的最优控制建议（第一步的动作） |
| mpc_actions[].device | string | 目标设备: `heater` / `aerator` / `pump` / `feeder` / `light` |
| mpc_actions[].action | string | 建议动作（见下方动作表） |
| mpc_actions[].value | float | 归一化控制强度 [0, 1]，0=最低/关，1=最高/开 |
| mpc_predicted | array | 施加 MPC 最优控制后的预测轨迹（格式同 predictions） |
| anomaly_step | integer | 自由演化轨迹中首次出现异常的步数索引（从0开始），-1 表示全程安全 |
| anomaly_warning | string | 异常预警文本，为空字符串表示未来安全无风险 |
| confidence | float | 模型置信度 [0, 1]，越接近 1 越可靠 |
| generated_at | string | 预测结果的生成时间 (ISO 8601) |

---

### ModelStatusResponse
```json
{
  "model_loaded": "boolean",
  "version": "string",
  "trained_at": "string (ISO 8601)"
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| model_loaded | boolean | Koopman 模型是否已成功加载 |
| version | string | 模型版本号（来自模型 JSON 文件） |
| trained_at | string | 模型训练完成时间 |

---

## 设备类型

| device | 说明 |
|--------|------|
| heater | 加热器 |
| aerator | 增氧机 |
| pump | 水泵 |
| feeder | 投喂器 |
| light | 灯光 |

---

## 建议动作

### Analyse 端点动作

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

### Predict 端点 MPC 动作

MPC 使用三档归一化控制策略（基于 `value` 值域），动作名称与设备强绑定：

| device | value < 0.3 | 0.3 ≤ value ≤ 0.7 | value > 0.7 |
|--------|-------------|---------------------|-------------|
| heater | `cool` 降温 | `maintain` 保持 | `heat` 加热 |
| aerator | `off` 关闭 | `maintain` 保持 | `on` 开启 |
| pump | `off` 关闭 | `maintain` 保持 | `on` 开启 |
| feeder | `hold` 停止投喂 | `maintain` 保持 | `feed` 投喂 |
| light | `dim` 调暗 | `maintain` 保持 | `brighten` 调亮 |

---

## 水质参数安全范围

| 参数 | 单位 | 最小值 | 最大值 | 说明 |
|------|------|--------|--------|------|
| temperature | °C | 22 | 28 | 水温 |
| ph | — | 6.5 | 8.0 | 酸碱度 |
| oxygen | mg/L | 5 | 9 | 溶解氧 |
| ammonia | mg/L | 0 | 0.5 | 氨氮（有毒） |
| water_level | % | 60 | 90 | 水位 |

> 当 `predictions[].is_anomaly` 为 `true` 时，表示该预测步中至少有一个参数将超出上述安全范围。`anomaly_step` 指示首次越界的步数，`anomaly_warning` 提供人类可读的预警文本。

---

## 预测原理

### Koopman 算子

将鱼缸的非线性水质动态 `x_{k+1} = f(x_k, u_k)` 通过可观测量提升函数 ψ(x) 映射到高维空间，使其近似线性：

```
ψ(x_{k+1}) ≈ K · ψ(x_k) + B · u_k
```

- `K`：Koopman 矩阵（lift_dim × lift_dim），通过 EDMD 从历史数据学习
- `B`：控制矩阵（lift_dim × 5），描述执行器对提升态的影响
- `ψ(x)`：可观测量函数，由 identity + quadratic + cross + trigonometric 组成

### MPC 控制器

由于 Koopman 模型是线性的，MPC 优化问题归约为凸二次规划（QP），可高效求解全局最优解：

```
min  Σ[(z_k - z_ref)ᵀ Q (z_k - z_ref) + u_kᵀ R u_k]
s.t. ψ_{k+1} = K·ψ_k + B·u_k
     0 ≤ u_k ≤ 1
     z_min ≤ z_k ≤ z_max
```

- 代价矩阵 Q 中氨氮权重加倍（有毒物质优先控制）
- 代价矩阵 R 中投喂器代价减半（鼓励适量喂食）
- 氨氮设定点取下界 0（目标是最小化有毒物质）

---

## 注意事项

1. **框架差异**: Mind 服务使用 Hertz 框架，而其他服务使用 Gin。考虑统一框架以降低维护成本。

2. **数据来源**: 如果请求体中不提供 `data` 字段，服务将从数据库获取历史传感器数据进行分析。

3. **分析逻辑**: 具体的分析算法和决策逻辑在 `mind/internal/service/analyse.go` 中实现。

4. **预测服务可选**: 预测端点依赖 Koopman 模型文件（`./models/koopman_model.json`）。模型未加载时预测端点返回 503，分析端点不受影响。预测是增强功能，不是核心依赖。

5. **模型热重载**: 模型支持运行时热重载，替换磁盘上的模型文件后调用热重载接口即可生效，无需重启服务。

6. **NATS 状态缓存**: 预测服务通过后台订阅 NATS `sensor.update` 事件维护鱼缸最新状态缓存，减少对 sensor 服务的实时查询依赖。

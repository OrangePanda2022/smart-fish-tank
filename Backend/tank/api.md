# Tank Service API 文档

## 基础信息

- **Base URL**: `/api/v1`
- **框架**: Gin
- **认证方式**: Session (通过 Cookie 或 Authorization Header)
- **用户标识**: 需要认证后从会话中获取 user_id

---

## 认证说明

### 认证方式

请求中需要包含有效的会话信息，通过以下方式之一传递：

1. **Cookie 方式** (自动处理):
```javascript
credentials: 'include'  // fetch 请求会自动发送 Cookie
```

2. **Header 方式**:
```javascript
headers: {
  'Authorization': 'Bearer <session_token>'
}
```

### 获取用户 ID

当前实现中，用户 ID 从会话中间件获取。如果未登录或会话无效，接口返回错误。

---

## 接口列表

### 1. 创建鱼缸

**端点**: `PATCH /api/v1/tank/`

**认证**: 是

**注意**: 当前路由使用 `PATCH` 方法，但根据 RESTful 规范，创建资源应使用 `POST`。

**请求体**:
```json
{
  "tank_name": "我的鱼缸",
  "tank_size": 100
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tank_name | string | 是 | 鱼缸名称 |
| tank_size | integer | 是 | 鱼缸容量 (升) |

**请求示例** (JavaScript):
```javascript
const response = await fetch('/api/v1/tank/', {
  method: 'PATCH',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer <session_token>'
  },
  credentials: 'include',
  body: JSON.stringify({
    tank_name: '我的鱼缸',
    tank_size: 100
  })
});

const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "msg": "ok"
}
```

**错误响应**:
```json
{
  "msg": "error",
  "error": "错误描述"
}
```

---

### 2. 获取当前用户的所有鱼缸

**端点**: `GET /api/v1/tank/`

**认证**: 是

**请求示例** (JavaScript):
```javascript
const response = await fetch('/api/v1/tank/', {
  method: 'GET',
  headers: {
    'Authorization': 'Bearer <session_token>'
  },
  credentials: 'include'
});

const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "msg": "ok",
  "tanks": [
    {
      "tank_id": "tank-uuid-123",
      "tank_name": "我的鱼缸",
      "tank_size": 100,
      "user_id": "user-uuid",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

**响应字段说明**:
| 字段 | 类型 | 说明 |
|------|------|------|
| msg | string | 状态消息 "ok" |
| tanks | array | 鱼缸列表 |
| tanks[].tank_id | string | 鱼缸 ID |
| tanks[].tank_name | string | 鱼缸名称 |
| tanks[].tank_size | integer | 鱼缸容量 (升) |
| tanks[].user_id | string | 所属用户 ID |
| tanks[].created_at | string | 创建时间 |
| tanks[].updated_at | string | 更新时间 |

---

### 3. 获取指定鱼缸详情

**端点**: `GET /api/v1/tank/:tank_id`

**认证**: 是

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tank_id | string | 是 | 鱼缸 ID |

**请求示例** (JavaScript):
```javascript
const tankId = 'tank-uuid-123';
const response = await fetch(`/api/v1/tank/${tankId}`, {
  method: 'GET',
  headers: {
    'Authorization': 'Bearer <session_token>'
  },
  credentials: 'include'
});

const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "msg": "ok",
  "tank": {
    "tank_id": "tank-uuid-123",
    "tank_name": "我的鱼缸",
    "tank_size": 100,
    "user_id": "user-uuid",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### 4. 更新鱼缸

**端点**: `PUT /api/v1/tank/:tank_id`

**认证**: 是

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tank_id | string | 是 | 鱼缸 ID |

**请求体**:
```json
{
  "tank_name": "新鱼缸名称",
  "tank_size": 150
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tank_name | string | 否 | 新的鱼缸名称 |
| tank_size | integer | 否 | 新的鱼缸容量 (升) |

**请求示例** (JavaScript):
```javascript
const tankId = 'tank-uuid-123';
const response = await fetch(`/api/v1/tank/${tankId}`, {
  method: 'PUT',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer <session_token>'
  },
  credentials: 'include',
  body: JSON.stringify({
    tank_name: '新鱼缸名称',
    tank_size: 150
  })
});

const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "msg": "ok",
  "tank": {
    "tank_id": "tank-uuid-123",
    "tank_name": "新鱼缸名称",
    "tank_size": 150,
    "user_id": "user-uuid",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-02T00:00:00Z"
  }
}
```

**注意**: 当前实现中，更新后会返回更新后的完整鱼缸对象。

---

### 5. 删除鱼缸

**端点**: `DELETE /api/v1/tank/:tank_id`

**认证**: 是

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tank_id | string | 是 | 鱼缸 ID |

**请求示例** (JavaScript):
```javascript
const tankId = 'tank-uuid-123';
const response = await fetch(`/api/v1/tank/${tankId}`, {
  method: 'DELETE',
  headers: {
    'Authorization': 'Bearer <session_token>'
  },
  credentials: 'include'
});

const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "msg": "ok",
  "tank": {
    "tank_id": "tank-uuid-123",
    "tank_name": "我的鱼缸",
    "tank_size": 100,
    "user_id": "user-uuid",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-02T00:00:00Z"
  }
}
```

**说明**: 删除后返回被删除的鱼缸数据。

---

### 6. 推送直播帧

**端点**: `POST /api/v1/tank/:tank_id/stream/frame`

**认证**: 暂无（MVP阶段，TODO: X-Device-Key 设备认证）

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tank_id | string | 是 | 鱼缸 ID |

**请求体**: raw JPEG 二进制数据（`Content-Type: image/jpeg`）

**ESP32-CAM 示例** (Arduino):
```cpp
camera_fb_t* fb = esp_camera_fb_get();
HTTPClient http;
http.begin("http://<host>:8084/api/v1/tank/<tank_id>/stream/frame");
http.addHeader("Content-Type", "image/jpeg");
http.POST(fb->buf, fb->len);
http.end();
esp_camera_fb_return(fb);
delay(1000); // 1 FPS
```

**成功响应** (200):
```json
{
  "msg": "ok"
}
```

**错误响应** (400):
```json
{
  "msg": "error",
  "error": "not a JPEG"
}
```

---

### 7. 获取直播流

**端点**: `GET /api/v1/tank/:tank_id/stream`

**认证**: 是

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tank_id | string | 是 | 鱼缸 ID |

**响应格式**: `multipart/x-mixed-replace`（MJPEG 流）

**浏览器使用**:
```html
<img src="http://<host>:8084/api/v1/tank/<tank_id>/stream" />
```

无需任何 JavaScript，`<img>` 标签即可直接显示实时画面。

**说明**:
- 长连接，浏览器关闭标签或导航离开后服务端自动清理资源
- 1-3秒帧间隔下延迟约100-500ms
- 每个 tank 支持同时1-3个观看者
- 当前 `getUserID()` 为硬编码，所有权验证需等 auth middleware 接入

---

## 通用响应格式

### 成功响应
```json
{
  "msg": "ok",
  "tank": { ... }
}
```

或

```json
{
  "msg": "ok",
  "tanks": [...]
}
```

### 错误响应
```json
{
  "msg": "error",
  "error": "错误描述"
}
```

---

## 数据模型

### Tank (鱼缸)
```json
{
  "tank_id": "uuid-string",
  "tank_name": "string",
  "tank_size": 100,
  "user_id": "uuid-string",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### CreateTankRequest (创建请求)
```json
{
  "tank_name": "string",
  "tank_size": 100
}
```

### UpdateTankRequest (更新请求)
```json
{
  "tank_name": "string",
  "tank_size": 100
}
```

---

## 注意事项

1. **用户隔离**: 当前实现中，鱼缸操作会从会话中获取 user_id，用户只能操作自己的鱼缸。

2. **路由方法问题**: 创建鱼缸使用 `PATCH` 方法，建议改为 `POST` 以符合 RESTful 规范。

3. **错误处理**: 当前实现使用 `fmt.Println` 输出错误而非结构化日志，且部分错误未返回合适的 HTTP 状态码。

4. **直播流网关**: MJPEG 流是长连接，不应经过 gateway 代理（会被超时杀死），浏览器应直连 tank 服务 8084 端口。

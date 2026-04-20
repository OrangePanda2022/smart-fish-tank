# Gateway Service API 文档

## 基础信息

- **Base URL**: `/`
- **框架**: Gin
- **认证方式**: 透传下游服务认证

---

## 架构说明

Gateway 是系统的 API 网关，负责：
1. **路由转发**: 将请求转发到对应的后端微服务
2. **负载均衡**: 使用 Sticky Balancing 确保同一用户的请求转发到同一节点
3. **聚合请求**: 支持合并多个服务的响应

---

## 路由配置

Gateway 通过路径前缀匹配后端服务：

| 路径前缀 | 目标服务 |
|----------|----------|
| `/api/v1/auth/*` | auth-service |
| `/api/v1/users/*` | auth-service |
| `/api/v1/admin/*` | auth-service |
| `/api/v1/oauth2/*` | auth-service |
| `/api/v1/tanks/*` | sensor-service |
| `/api/v1/devices/*` | sensor-service |
| `/api/v1/tank/*` | tank-service |
| `/api/v1/analyse/*` | mind-service |

---

## 接口列表

### 1. 代理请求

**端点**: `GET /<any-path>`

**认证**: 透传下游服务认证

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| aggregate | boolean | 否 | 是否启用聚合模式 |

**请求示例** (JavaScript):
```javascript
// 普通代理请求
const response = await fetch('/api/v1/tank/', {
  method: 'GET',
  headers: {
    'Authorization': 'Bearer <token>'
  },
  credentials: 'include'
});

const data = await response.json();
console.log(data);
```

**成功响应**: 透传下游服务响应

**错误响应** (未找到路由):
```json
{
  "error": "路由未找到或服务不可用"
}
```

---

### 2. 聚合请求

**端点**: `GET /<any-path>?aggregate=true`

**认证**: 透传下游服务认证

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| aggregate | boolean | 是 | 必须设置为 `true` |
| services | string | 是 | 逗号分隔的服务名列表，如 `user,order` |

**请求示例** (JavaScript):
```javascript
const response = await fetch('/api/v1/dashboard?aggregate=true&services=user,tank,sensor', {
  method: 'GET',
  headers: {
    'Authorization': 'Bearer <token>'
  },
  credentials: 'include'
});

const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "user": {
    "code": 0,
    "message": "success",
    "data": { ... }
  },
  "tank": {
    "code": 0,
    "message": "success",
    "data": { ... }
  },
  "sensor": {
    "code": 0,
    "message": "success",
    "data": { ... }
  }
}
```

**响应说明**:
- 响应是一个对象，key 是请求中指定的服务名
- 每个服务的响应是独立返回的
- 如果某个服务失败，不影响其他服务

**错误响应** (聚合失败):
```json
{
  "error": "Aggregation failed",
  "details": "错误详情"
}
```

---

## 负载均衡策略

### Sticky Balancing

Gateway 使用 Sticky Balancing 确保同一用户的请求在短时间内被转发到同一个后端节点。

**实现方式**:
1. 从请求中提取用户标识或生成随机数
2. 使用 Jump Hash 算法选择目标节点
3. 保证相同标识的请求始终选择相同节点

**优势**:
- 提高缓存命中率
- 保证 session 一致性
- 减少跨节点数据同步开销

---

## 路由匹配规则

### 匹配算法

```go
// 简化逻辑
func Match(path string) (serviceName string, exists bool) {
    for prefix, serviceName := range routes {
        if strings.HasPrefix(path, prefix) {
            return serviceName, true
        }
    }
    return "", false
}
```

### 匹配示例

| 请求路径 | 匹配前缀 | 目标服务 |
|----------|----------|----------|
| `/api/v1/users/me` | `/api/v1/users` | auth-service |
| `/api/v1/tank/` | `/api/v1/tank` | tank-service |
| `/api/v1/analyse/123` | `/api/v1/analyse` | mind-service |
| `/api/v1/tanks/123/sensors/latest` | `/api/v1/tanks` | sensor-service |

### 潜在问题

**路由冲突**: 简单前缀匹配可能导致路由冲突。

例如：
- `/api/v1/users/me` → auth-service
- `/api/v1/users/123` → auth-service

如果后续添加 `/api/v1/user-actions/*` 路由，由于使用 `HasPrefix` 匹配，可能产生意外行为。

**建议**: 使用更精确的路由匹配算法，如最长前缀匹配或正则表达式匹配。

---

## 中间件链

请求经过以下中间件处理：

1. **RequestID**: 生成并传递请求追踪 ID
2. **Auth**: JWT Token 验证 (仅验证格式，不验证用户)
3. **Proxy**: 转发到后端服务

---

## 错误处理

### 当前实现

Gateway 的错误处理尚未统一：

```go
// TODO 统一错误返回
if h.strict {
    return
}

// TODO 统一返回 统一错误处理
c.JSON(500, gin.H{"error": "Aggregation failed", "details": err.Error()})
```

### 建议改进

1. 实现统一的错误响应格式
2. 根据下游服务响应状态码决定 Gateway 响应码
3. 添加错误日志和监控

---

## 健康检查

Gateway 本身不提供独立健康检查端点。健康检查通过负载均衡器的节点健康检测实现。

**检测机制**:
1. 定期向各节点发送探测请求
2. 标记不响应或响应超时的节点为不健康
3. 不健康节点不参与负载均衡

---

## 请求流程

```
Client Request
      │
      ▼
┌─────────────┐
│   Gateway   │
└─────────────┘
      │
      ▼
 Match Route ──────────────────┐
      │                        │
      ▼                        ▼
┌─────────────┐         ┌─────────────┐
│   找到路由   │         │  未找到路由  │
└─────────────┘         └─────────────┘
      │                        │
      ▼                        ▼
 Check aggregate ──────► Return 404
      │
      ▼
   Yes │ No
      │   │
      ▼   ▼
┌─────────────┐    ┌─────────────┐
│  Aggregate  │    │   Select    │
│  Requests   │    │   Node      │
└─────────────┘    └─────────────┘
      │                   │
      └─────────┬─────────┘
                ▼
         ┌─────────────┐
         │   Forward   │
         │   Request   │
         └─────────────┘
                │
                ▼
         ┌─────────────┐
         │  Return     │
         │  Response   │
         └─────────────┘
```

---

## 配置说明

### 路由配置

路由表在 `gateway/internal/domain/route/route.go` 中定义：

```go
type Route struct {
    ID       RouteID
    Name     string
    Path     string
    Method   string
    Upstream upstream.UpstreamID
    Plugins  []string
    Enabled  bool
}
```

### 上游服务配置

每个服务可以有多个节点（实例），支持水平扩展。

---

## 注意事项

1. **严格模式**: 当前存在 `strict` 模式但功能未完成，可能导致意外行为。

2. **聚合功能**: 聚合请求的路径拼接存在 TODO 标注，需要验证拼接逻辑是否正确。

3. **会话保持**: 依赖下游服务处理会话，Gateway 仅做透传。

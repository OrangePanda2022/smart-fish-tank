# Auth Service API 文档

## 基础信息

- **Base URL**: `/api/v1`
- **框架**: Gin
- **认证方式**: Session Cookie + JWT Bearer Token

---

## 认证说明

### 需要认证的接口

需要在上方请求中包含以下 Header：

```
Authorization: Bearer <jwt_token>
Cookie: session=<session_id>
```

### 公开接口

无需认证即可访问的接口。

---

## 接口列表

### 1. 用户注册

**端点**: `POST /api/v1/auth/register`

**认证**: 否

**请求体**:
```json
{
  "email": "user@example.com",
  "password": "yourpassword123",
  "user_name": "JohnDoe"
}
```

**请求示例** (JavaScript):
```javascript
const response = await fetch('/api/v1/auth/register', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    email: 'user@example.com',
    password: 'yourpassword123',
    user_name: 'JohnDoe'
  })
});

const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "profile": {
      "user_id": "uuid-string",
      "user_email": "user@example.com",
      "user_name": "JohnDoe",
      "user_roles": ["user"],
      "user_status": "active",
      "created_at": "2024-01-01T00:00:00Z"
    },
    "session_id": "session-uuid-string"
  }
}
```

**错误响应** (400):
```json
{
  "code": 400,
  "message": "bad request",
  "error": "具体错误信息"
}
```

---

### 2. 用户登录

**端点**: `POST /api/v1/auth/login`

**认证**: 否

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| login_challenge | string | 是 | Hydra 登录挑战码 |

**请求体**:
```json
{
  "email": "user@example.com",
  "password": "yourpassword123"
}
```

**请求示例** (JavaScript):
```javascript
const loginChallenge = 'your-login-challenge-token';
const response = await fetch(`/api/v1/auth/login?login_challenge=${loginChallenge}`, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  credentials: 'include',
  body: JSON.stringify({
    email: 'user@example.com',
    password: 'yourpassword123'
  })
});

const data = await response.json();
console.log(data);
// 成功后会设置 Cookie: session=<session_id>
```

**成功响应** (200):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "profile": {
      "user_id": "uuid-string",
      "user_email": "user@example.com",
      "user_name": "JohnDoe",
      "user_roles": ["user"],
      "user_status": "active"
    },
    "session_id": "session-uuid-string",
    "redirect_to": "https://client-app.com/callback?code=..."
  }
}
```

---

### 3. TOTP 验证

**端点**: `POST /api/v1/auth/totp/verify`

**认证**: 否

**请求体**:
```json
{
  "email": "user@example.com",
  "code": "123456"
}
```

**请求示例** (JavaScript):
```javascript
const response = await fetch('/api/v1/auth/totp/verify', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    email: 'user@example.com',
    code: '123456'
  })
});

const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "reset_token": "one-time-reset-token"
  }
}
```

---

### 4. 请求密码重置

**端点**: `POST /api/v1/auth/password/reset/request`

**认证**: 否

**请求体**:
```json
{
  "email": "user@example.com"
}
```

**请求示例** (JavaScript):
```javascript
const response = await fetch('/api/v1/auth/password/reset/request', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    email: 'user@example.com'
  })
});

const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "next_step": ""
  }
}
```

**说明**: 接口本身不暴露用户是否存在，即使邮箱不存在也返回成功。

---

### 5. 确认密码重置

**端点**: `POST /api/v1/auth/password/reset/confirm`

**认证**: 否

**请求体**:
```json
{
  "reset_token": "one-time-reset-token-from-verify-totp",
  "new_password": "newpassword123"
}
```

**请求示例** (JavaScript):
```javascript
const response = await fetch('/api/v1/auth/password/reset/confirm', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    reset_token: 'one-time-reset-token-from-verify-totp',
    new_password: 'newpassword123'
  })
});

const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "password reset success"
  }
}
```

---

### 6. OAuth2 获取登录请求

**端点**: `GET /api/v1/oauth2/login`

**认证**: 否

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| login_challenge | string | 是 | Hydra 登录挑战码 |

**请求示例** (JavaScript):
```javascript
const challenge = 'hydra-login-challenge-token';
const response = await fetch(`/api/v1/oauth2/login?login_challenge=${challenge}`, {
  method: 'GET',
  credentials: 'include'
});

const data = await response.json();
console.log(data);
```

**成功响应** (200) - 需要展示登录页:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "login_challenge": "hydra-login-challenge-token",
    "skip": false,
    "client_id": "client-app-id",
    "requested_scope": ["openid", "profile", "email"],
    "request_url": "https://hydra/oauth2/auth?...",
    "oidc_context": {},
    "requested_aud": ["api"],
    "login_session_id": "session-id"
  }
}
```

**成功响应** (200) - 可跳过登录:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "redirect_to": "https://hydra/...",
    "skip": true
  }
}
```

---

### 7. OAuth2 获取同意请求

**端点**: `GET /api/v1/oauth2/consent`

**认证**: 可选（已登录用户可跳过同意）

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| consent_challenge | string | 是 | Hydra 同意挑战码 |

**请求示例** (JavaScript):
```javascript
const challenge = 'hydra-consent-challenge-token';
const response = await fetch(`/api/v1/oauth2/consent?consent_challenge=${challenge}`, {
  method: 'GET',
  credentials: 'include'
});

const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "consent_challenge": "hydra-consent-challenge-token",
    "skip": false,
    "subject": "user-uuid",
    "client_id": "client-app-id",
    "client_name": "My Application",
    "requested_scope": ["openid", "profile"],
    "requested_aud": ["api"],
    "request_url": "https://hydra/oauth2/auth?..."
  }
}
```

---

### 8. OAuth2 提交同意决策

**端点**: `POST /api/v1/oauth2/consent`

**认证**: 是

**请求体**:
```json
{
  "consent_challenge": "hydra-consent-challenge-token",
  "accept": true,
  "grant_scope": ["openid", "profile", "email"],
  "grant_audience": ["api"],
  "remember": true,
  "remember_for": 3600
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| consent_challenge | string | 是 | 同意挑战码 |
| accept | boolean | 是 | 是否同意 |
| grant_scope | array | 否 | 授权的权限范围 |
| grant_audience | array | 否 | 授权的受众 |
| remember | boolean | 否 | 是否记住同意决策 |
| remember_for | integer | 否 | 记住时长（秒） |

**请求示例** (JavaScript):
```javascript
const response = await fetch('/api/v1/oauth2/consent', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  credentials: 'include',
  body: JSON.stringify({
    consent_challenge: 'hydra-consent-challenge-token',
    accept: true,
    grant_scope: ['openid', 'profile', 'email'],
    grant_audience: ['api'],
    remember: true,
    remember_for: 3600
  })
});

const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "redirect_to": "https://client-app.com/callback?code=...",
    "accepted": true
  }
}
```

---

### 9. 用户登出

**端点**: `POST /api/v1/users/logout`

**认证**: 是

**请求体**:
```json
{
  "refresh_token": "user-refresh-token"
}
```

**请求示例** (JavaScript):
```javascript
const response = await fetch('/api/v1/users/logout', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer <jwt_token>'
  },
  credentials: 'include',
  body: JSON.stringify({
    refresh_token: 'user-refresh-token'
  })
});

const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "logout success"
  }
}
```

---

### 10. 获取当前用户资料

**端点**: `GET /api/v1/users/me`

**认证**: 是

**请求示例** (JavaScript):
```javascript
const response = await fetch('/api/v1/users/me', {
  method: 'GET',
  headers: {
    'Authorization': 'Bearer <jwt_token>'
  },
  credentials: 'include'
});

const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": "uuid-string",
    "user_email": "user@example.com",
    "user_name": "JohnDoe",
    "user_roles": ["user"],
    "user_status": "active",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### 11. 修改昵称

**端点**: `PATCH /api/v1/users/name`

**认证**: 是

**请求体**:
```json
{
  "user_name": "NewName"
}
```

**请求示例** (JavaScript):
```javascript
const response = await fetch('/api/v1/users/name', {
  method: 'PATCH',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer <jwt_token>'
  },
  credentials: 'include',
  body: JSON.stringify({
    user_name: 'NewName'
  })
});

const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "user name updated"
  }
}
```

---

### 12. 修改密码

**端点**: `PATCH /api/v1/users/password`

**认证**: 是

**请求体**:
```json
{
  "old_password": "oldpassword123",
  "new_password": "newpassword123"
}
```

**请求示例** (JavaScript):
```javascript
const response = await fetch('/api/v1/users/password', {
  method: 'PATCH',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer <jwt_token>'
  },
  credentials: 'include',
  body: JSON.stringify({
    old_password: 'oldpassword123',
    new_password: 'newpassword123'
  })
});

const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "password updated"
  }
}
```

**说明**: 修改密码后会使所有历史刷新令牌失效。

---

### 13. 删除用户（软删除）

**端点**: `DELETE /api/v1/users/me`

**认证**: 是

**请求示例** (JavaScript):
```javascript
const response = await fetch('/api/v1/users/me', {
  method: 'DELETE',
  headers: {
    'Authorization': 'Bearer <jwt_token>'
  },
  credentials: 'include'
});

const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "user deleted"
  }
}
```

---

### 14. 管理员 - 设置用户状态

**端点**: `PATCH /api/v1/admin/users/:user_id/status`

**认证**: 是 (管理员)

**路径参数**:
| 参数 | 类型 | 说明 |
|------|------|------|
| user_id | string | 目标用户 ID |

**请求体**:
```json
{
  "enabled": false
}
```

**请求示例** (JavaScript):
```javascript
const response = await fetch('/api/v1/admin/users/target-user-uuid/status', {
  method: 'PATCH',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer <admin_jwt_token>'
  },
  credentials: 'include',
  body: JSON.stringify({
    enabled: false
  })
});

const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "user status updated"
  }
}
```

---

### 15. 管理员 - 设置用户角色

**端点**: `PATCH /api/v1/admin/users/:user_id/role`

**认证**: 是 (管理员)

**路径参数**:
| 参数 | 类型 | 说明 |
|------|------|------|
| user_id | string | 目标用户 ID |

**请求体**:
```json
{
  "role": ["admin"]
}
```

**请求示例** (JavaScript):
```javascript
const response = await fetch('/api/v1/admin/users/target-user-uuid/role', {
  method: 'PATCH',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer <admin_jwt_token>'
  },
  credentials: 'include',
  body: JSON.stringify({
    role: ['admin']
  })
});

const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "user role updated"
  }
}
```

---

### 16. 健康检查

**端点**: `GET /healthz`

**认证**: 否

**请求示例** (JavaScript):
```javascript
const response = await fetch('/healthz');
const data = await response.json();
console.log(data);
```

**成功响应** (200):
```json
{
  "status": "ok"
}
```

---

## 通用响应格式

### 成功响应
```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

### 错误响应
```json
{
  "code": <error_code>,
  "message": "error",
  "error": "错误详情描述"
}
```

### 常见错误码
| code | 说明 |
|------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或认证过期 |
| 403 | 无权限访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 数据模型

### UserProfile
```json
{
  "user_id": "uuid-string",
  "user_email": "user@example.com",
  "user_name": "JohnDoe",
  "user_roles": ["user", "admin"],
  "user_status": "active",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### 用户角色
| 角色 | 说明 |
|------|------|
| user | 普通用户 |
| admin | 管理员 |

### 用户状态
| 状态 | 说明 |
|------|------|
| active | 正常 |
| disabled | 已禁用 |
| deleted | 已删除 |

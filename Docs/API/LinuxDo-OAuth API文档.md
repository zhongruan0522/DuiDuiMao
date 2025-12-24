# LinuxDo Connect OAuth2 接入文档

## 概述

LinuxDo Connect 是 LinuxDo 论坛提供的 OAuth2 认证服务，允许第三方应用使用 LinuxDo 账号登录。

## 端点信息

| 用途 | URL |
|------|-----|
| 授权端点 | `https://connect.linux.do/oauth2/authorize` |
| Token端点 | `https://connect.linux.do/oauth2/token` |
| 用户信息端点 | `https://connect.linux.do/api/user` |

## 授权流程

### 1. 跳转授权页面

```
GET https://connect.linux.do/oauth2/authorize
```

**参数：**

| 参数 | 必填 | 说明 |
|------|------|------|
| client_id | 是 | 应用ID |
| response_type | 是 | 固定值 `code` |
| redirect_uri | 是 | 回调地址 |
| state | 推荐 | 防CSRF随机字符串 |

**示例：**
```
https://connect.linux.do/oauth2/authorize?client_id=YOUR_CLIENT_ID&response_type=code&redirect_uri=https://example.com/callback&state=random_string
```

### 2. 获取 Access Token

用户授权后会携带 `code` 跳转到回调地址，用 `code` 换取 `access_token`。

```
POST https://connect.linux.do/oauth2/token
```

**认证方式：** HTTP Basic Auth
- Username: `client_id`
- Password: `client_secret`

**请求体（form-urlencoded）：**

| 参数 | 必填 | 说明 |
|------|------|------|
| grant_type | 是 | 固定值 `authorization_code` |
| code | 是 | 授权码 |
| redirect_uri | 是 | 回调地址（需与授权时一致） |

**响应示例：**
```json
{
    "access_token": "xxxxxxxxxxxxxxxx",
    "token_type": "Bearer",
    "expires_in": 3600
}
```

### 3. 获取用户信息

```
GET https://connect.linux.do/api/user
```

**请求头：**
```
Authorization: Bearer {access_token}
```

**响应示例：**
```json
{
    "id": 124,
    "username": "Bee",
    "name": "(  ⩌   ˰ ⩌)",
    "active": true,
    "trust_level": 2,
    "silenced": false
}
```

**字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 用户唯一ID |
| username | string | 用户名 |
| name | string | 昵称 |
| active | bool | 账号是否激活 |
| trust_level | int | 信任等级（0-4） |
| silenced | bool | 是否被禁言 |

## 信任等级说明

| 等级 | 说明 |
|------|------|
| 0 | 新用户 |
| 1 | 基本用户 |
| 2 | 成员 |
| 3 | 活跃用户 |
| 4 | 领导者 |

## 注意事项

1. `state` 参数务必使用并校验，防止 CSRF 攻击
2. `redirect_uri` 必须与申请应用时填写的一致
3. `access_token` 有过期时间，注意处理过期情况
4. 建议根据 `trust_level` 设置不同的兑换权限/额度
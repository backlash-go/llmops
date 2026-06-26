# 登录模块

## 功能

实现:

- Keycloak 登录

```
我要实现我的 前端应用    后端接口 myapp
1. 浏览器访问 点击keylocak  登录
2. myapp 生成随机 state
3. myapp 把 state 相关信息存入 我的前端应用domain path 浏览器的 域 Cookie 设置一个时间
4. myapp 通过 URL 参数把 state 带给 Keycloak
5. Keycloak 认证结束带着CODE STATE 等等给我 放进回调 URL
7. myapp 验证 URL state 

```

- token 解析 之后

```


{
  "acr": "1",
  "at_hash": "5E184nPuJFJ36nvtrJrqrA",
  "aud": "xmit-llmops",
  "auth_time": 1782452389,
  "azp": "xmit-llmops",
  "email": "xianbin.xi@luxshare-ict.com",
  "email_verified": true,
  "exp": 1782452690,
  "family_name": "席贤斌",
  "given_name": "Xianbin.Xi",
  "iat": 1782452390,
  "iss": "https://smart-auto.luxshare-ict.com/keycloak/realms/xmit-luxshare",
  "jti": "cd61c664-0daa-7fdb-9eb2-4fd9d06311e6",
  "name": "Xianbin.Xi 席贤斌",
  "preferred_username": "31070182",
  "roles": [
    "llmops-admin"
  ],
  "sid": "gOSMZuyxOzAzQmirOAMKOdme",
  "sub": "dc4663db-5a90-4333-8471-db0a88b5520d",
  "typ": "ID"
}

```

```

拿到 Keycloak token claims
→ 先查 user_identity(provider + issuer + subject)
→ 不存在：
    创建 user，拿 user.ID
    再创建 user_identity
→ 存在：
    查 user
    对比 email / first_name / last_name
    不同才更新 user
    同步更新 user_identity 的 provider 信息/raw_profile
```
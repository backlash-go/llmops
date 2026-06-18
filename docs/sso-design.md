
前端

https://smart-auto.com/ai-platform

服务端

https://smart-auto.com/ai-platform/api/v1



用户浏览器
    │
    ├─→ 前端: https://smart-auto.com/ai-platform
    │       (React/Vue SPA)
    │
    ├─→ 后端: https://smart-auto.com/ai-platform/api/v1
    │       (API Server)
    │
    └─→ Keycloak: https://<keycloak-domain>/realms/<realm>
            (身份认证中心)




┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   前端 SPA   │     │   后端 API   │     │   Keycloak   │
└──────┬───────┘     └──────┬───────┘     └──────┬───────┘
       │                     │                     │
       │  1. 访问页面，检查token                    │
       │─────────────────────────────────────────→│
       │  2. 无token/过期 → 重定向到Keycloak登录    │
       │                     │                     │
       │  3. 如果已SSO登录过，Keycloak直接回调      │
       │     (无需再次输入密码)                     │
       │←─────────────────────────────────────────│
       │  4. 携带 authorization_code               │
       │                     │                     │
       │  5. code换token ────→│                    │
       │                     │─── code+secret ───→│
       │                     │←── access_token ───│
       │  6. 返回JWT token ←─│                    │
       │                     │                     │
       │  7. 后续请求携带token │                    │
       │────────────────────→│  验证token签名      │
       │                     │                     │



SSO 关键点
首次登录：用户被重定向到 Keycloak 登录页，输入账号密码验证
SSO 免登录：如果用户已在同一 Keycloak Realm 下的其他系统登录过，Keycloak 的 session cookie 仍然有效，会直接返回 authorization code，无需再次输入密码







            

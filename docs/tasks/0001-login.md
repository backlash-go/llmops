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
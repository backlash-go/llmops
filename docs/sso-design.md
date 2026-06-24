
前端

https://smart-auto.com/ai-platform

服务端

https://smart-auto.com/ai-platform/api/v1










SSO 关键点
首次登录：用户被重定向到 Keycloak 登录页，输入账号密码验证
SSO 免登录：如果用户已在同一 Keycloak Realm 下的其他系统登录过，Keycloak 的 session cookie 仍然有效，会直接返回 authorization code，无需再次输入密码





```

1. 测试发现文档（无需登录）

 1007  curl https://keycloak.kuaifuinfo.com/realms/test-realm-1/.well-known/openid-configuration | jq
 
2. 获取公钥（无需登录） 
 1008  curl https://keycloak.kuaifuinfo.com/realms/test-realm-1/protocol/openid-connect/certs | jq

https://keycloak.kuaifuinfo.com/realms/test-realm-1/protocol/openid-connect/auth?client_id=grafana-oauth-kc&response_type=code&scope=openid%20email%20profile%20offline_access%20roles&state=abc123&redirect_uri=https://sgrafana.kuaifuinfo.com/login/generic_oauth


```

3  认真成功 带着 code 请求  https://sgrafana.kuaifuinfo.com/login/generic_oauth 

```
https://smart-auto.luxshare-ict.com/keycloak/realms/xmit-luxshare/protocol/openid-connect/auth?client_id=k8s-grafana-oauth&redirect_uri=https://smart-auto.luxshare-ict.com/grafana-k8s/login/generic_oauth&response_type=code&scope=openid profile email offline_access roles&state=2CFzzKHhWDPb5sU9zxxvuMwt2r-CQyigLj2mXRa6nSk=

```

####创建一个普通用户
test-realm-user-1
ahjKLjsj1Hl12

4. 用 code 换 token（这是服务端第一件真正要做的事）

```
curl -X POST \
https://smart-auto.luxshare-ict.com/keycloak/realms/xmit-luxshare/protocol/openid-connect/token \
-H "Content-Type: application/x-www-form-urlencoded" \
-d "grant_type=authorization_code" \
-d "client_id=k8s-grafana-oauth" \
-d "client_secret=4rmxWEjvj5v8VsfmXulQ3IE4Hzja94dj" \
-d "code=9565024b-6f33-db50-8f61-a557ca0e31eb.yeHHtq9iyZYJgB7WU8LrXosF.39ea55c9-478c-4607-a3fa-b4894b4324e1" \
-d "redirect_uri=https://smart-auto.luxshare-ict.com/grafana-k8s/login/generic_oauth"
```

response

```json
{
  "access_token":"eyJ...",
  "refresh_token":"eyJ...",
  "id_token":"eyJ..."
}
```


5. 用 access_token 获取用户信息
```shell
ACCESS_TOKEN=eyJ....

curl \
-H "Authorization: Bearer $ACCESS_TOKEN" \
https://keycloak.kuaifuinfo.com/realms/test-realm-1/protocol/openid-connect/userinfo | jq

```

```json
{
  "sub":"...",
  "preferred_username":"admin",
  "email":"admin@test.com",
  "name":"Administrator"
}
```


6. 测试 refresh token


```shell
curl -X POST \
https://keycloak.kuaifuinfo.com/realms/test-realm-1/protocol/openid-connect/token \
-H "Content-Type: application/x-www-form-urlencoded" \
-d "grant_type=refresh_token" \
-d "client_id=ai-platform-oauth-kc" \
-d "client_secret=xxx" \
-d "refresh_token=xxx"

```

```json
{
  "access_token":"new",
  "refresh_token":"new"
}
```






repeat test


```
https://keycloak.kuaifuinfo.com/realms/test-realm-1/protocol/openid-connect/auth
?client_id=ai-platform
&response_type=code
&scope=openid email profile
&redirect_uri=https://sgrafana.com/login/generic_oauth
&state=abc123

https://keycloak.kuaifuinfo.com/realms/test-realm-1/protocol/openid-connect/auth?client_id=ai-platform&response_type=code&scope=openid%20email%20profile%20offline_access%20roles&redirect_uri=https://httpbin.org/get

####创建一个普通用户
test-realm-user-1
ahjKLjsj1Hl12
https://keycloak.kuaifuinfo.com/admin/test-realm-1/console/
```

```


```


```shell


curl -X POST \
https://keycloak.kuaifuinfo.com/realms/test-realm-1/protocol/openid-connect/token \
-H "Content-Type: application/x-www-form-urlencoded" \
-d "grant_type=authorization_code" \
-d "client_id=ai-platform" \
-d "client_secret=cAAvlWPNXq031vzjEj077rdfeeJudxFy" \
-d "code=62c0135f-f96c-42fc-bc23-41c776420293.50018b43-5e3c-48cb-b4e8-3e056ee78d46.24013730-8c4a-4508-85ec-3637dde737ed" \
-d "redirect_uri=https://httpbin.org/get"

```


```
expired time type session
AUTH_SESSION_ID YTgzMzM3OGYtZWY5YS00MTEzLWFkYjAtYWFkODg2N2ZiNTYwLmJZV2xMY3BsRklmaFhvWlJjNVdCTmxfOVlHb3VCVUtUVi05SGNkZ2UxUzJESlZxaHUwaTdQMEJLTnJQdUpjNHVxWUNJRWlIUXdHWHdvU2V2bEE0a1VB
expired time type session
KEYCLOAK_IDENTITY eyJhbGciOiJIUzUxMiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICIxOTc1MjZjOC0wYjMyLTQ0YWItOTQxMy0wYTZkYmI3ZTViNjAifQ.eyJleHAiOjE3ODE4OTU0MjksImlhdCI6MTc4MTg1OTQyOSwianRpIjoiMDQ2ZGQ0ZDktOTBkYy1lOWMyLWVjOWQtYmMyMTg4MmZiMjYwIiwiaXNzIjoiaHR0cHM6Ly9rZXljbG9hay5rdWFpZnVpbmZvLmNvbS9yZWFsbXMvbWFzdGVyIiwic3ViIjoiNzVkYWVkMDItODg2MS00NzA1LWFmNDItNGRhYzczNTU0MTFlIiwidHlwIjoiU2VyaWFsaXplZC1JRCIsInNpZCI6ImE4MzMzNzhmLWVmOWEtNDExMy1hZGIwLWFhZDg4NjdmYjU2MCIsInN0YXRlX2NoZWNrZXIiOiJxYVpFOTc2NTctZlFGSW5YN3RlYk9xclJZRWExdXdnTW5jOS1wQ0pzS0xRIn0.6-IdPI_GsXa1U9UmrcJDm8SeIC6B9N7OoZYru9GG3_QmocuX1K0dwNNTaIVUqKnhniqdmQxQ5fDDgSFmO-nPLw
KEYCLOAK_SESSION gBL6BkfqSSWjpzatNbOtTG8FiCUjV6LKhs4c9AxwYLo    2026-06-19T20:40:14.965Z 
```


```
curl -vk 
-H 'Cookie: KEYCLOAK_SESSION=YmE0NGYwYjEtYmRkYy00YjNiLTk0NWQtYTEzNDdlZTQwMGQyLlFoVU1xM0tXTU0zUl84UnJFVlk4NGJTazZsY0hFWllxT2VqSU9DWjJXbXRFV2I1bVhfNDVSOTlrVHNTSEVmVkRKcWt2Sjl4NVlSSlNSdVJ3UjJRWTFn' 
'https://keycloak.kuaifuinfo.com/realms/test-realm-1/protocol/openid-connect/auth?client_id=ai-platform&response_type=code&scope=openid%20email%20profile%20offline_access%20roles&redirect_uri=https%3A%2F%2Fhttpbin.org%2Fget&state=test123&prompt=none'


```
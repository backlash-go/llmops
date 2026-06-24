# Repository Guidelines

## 目标与适用范围

本文以用户创建接口 `POST /ops/api/v1/users` 为标准样板，约束后续 HTTP API 的目录、分层、命名、参数校验、错误处理和数据访问方式。新增资源接口时，应保持以下调用链：

```text
Router → Controller → Service → Store → MySQL
                    ↘ API DTO   ↘ Model
```

`Create` 接口可作为分层参考，但不要直接复制其成功响应：请求中含有 `password`，敏感字段不得返回给客户端。

## 目录与职责

以 `user` 资源为例：

- `internal/apiserver/router/user.go`：注册路径、HTTP 方法和 Controller。
- `internal/apiserver/controller/v1/user/`：解析请求、调用 Service、写入响应。
- `internal/apiserver/service/v1/user.go`：业务逻辑、DTO 到 Model 的转换、业务错误映射。
- `internal/apiserver/store/store.go`：聚合各资源 Store。
- `internal/apiserver/store/mysql/user.go`：定义存储接口并实现数据库操作。
- `pkg/api/llmops/v1/user.go`：请求与响应 DTO。
- `internal/pkg/model/user.go`：数据库模型及 GORM 映射。
- `internal/pkg/code/`：统一业务错误码。

新增资源时沿用相同目录，不允许 Controller 直接访问数据库，也不要在 Store 中编写 HTTP 或业务响应逻辑。

## 新增接口的实现顺序

### 1. 定义 API DTO

在 `pkg/api/llmops/v1/<resource>.go` 中定义独立请求和响应：

```go
type CreateUserRequest struct {
    Username string `json:"username" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8,max=32"`
}

type CreateUserResponse struct {
    ID       uint64 `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}
```

请求 DTO、数据库 Model、响应 DTO 必须分离。密码、令牌、内部错误、软删除字段等不得出现在响应 DTO 中。JSON 字段统一使用 `snake_case`。

> **待补充：** 分页 DTO、通用响应字段及时间格式。

### 2. 定义数据库 Model

Model 放在 `internal/pkg/model/`。字段必须明确声明 `gorm` 与 `json` 标签；表名通过 `TableName()` 固定。唯一约束、非空约束和默认值应同时与数据库建表语句保持一致。

DTO 转 Model 可使用 `copier.Copy`，但必须检查字段是否确实存在于 Model 中。密码等敏感数据应先加密，再写入专用字段；禁止保存明文。

### 3. 扩展 Store

在资源 Store 接口中声明方法：

```go
Create(ctx context.Context, user *model.User) error
```

MySQL 实现只负责持久化，并透传底层错误：

```go
return u.db.WithContext(ctx).Create(user).Error
```

所有数据库操作必须绑定 `context.Context`。Store 不负责构造 HTTP 错误码。

### 4. 实现 Service

Service 接口使用 API DTO 作为输入，使用响应 DTO 或错误作为输出。职责包括：

- 转换 DTO 与 Model；
- 执行业务规则；
- 调用 Store；
- 将重复数据、数据不存在等底层错误映射为 `internal/pkg/code` 中的业务错误；
- 避免把 SQL 文本直接暴露给客户端。

新增错误码时，在 `internal/pkg/code/` 中按现有分段规则定义并注册。数据库未知错误统一映射为 `code.ErrDatabase`。

> **待补充：** 事务边界、重复键判断方式及错误文案规范。

### 5. 实现 Controller

每个动作单独放置文件，如 `create.go`、`get.go`、`update.go`。Controller 只执行以下流程：

1. 使用 `log.L(c)` 记录入口或必要上下文；
2. 使用 `c.ShouldBindJSON(&request)` 绑定并校验请求；
3. 绑定失败时返回 `code.ErrBind` 或 `code.ErrValidation`；
4. 调用 `u.srv.Users().Create(c, &request)`；
5. 使用 `core.WriteResponse(c, err, data)` 统一写响应；
6. 每次写入错误响应后立即 `return`。

优先使用 DTO 的 `binding` 标签完成字段校验，跨字段或业务校验放入 Service，避免在 Controller 重复判断 `Username == ""`。

### 6. 注册路由

资源路由集中在 `internal/apiserver/router/<resource>.go`：

```go
users := v.Group("/users")
users.POST("", userController.Create)
```

路径使用复数、小写、短横线风格；版本前缀由上层统一提供，例如 `/ops/api/v1`。Controller 通过构造函数注入 `store.Factory`，禁止使用临时数据库连接。

## 编码约定

- Go 文件必须符合 `gofmt`/`goimports`；本项目导入路径 `llmops/...` 单独分组。
- 导出标识符使用 `PascalCase`，内部标识符使用 `camelCase`，包名简短小写。
- 接口使用职责名称，如 `UserSrv`、`UserStore`；实现使用非导出类型，如 `userService`、`users`。
- 方法接收者保持简短且一致：`u *UserController`、`u *userService`。
- 注释说明“做什么”，避免重复代码本身；导出类型和方法应有完整注释。
- 不忽略可能失败的错误。若确实忽略，必须写明原因。
- 日志中不得记录密码、令牌、Cookie、Authorization Header 或完整个人敏感信息。

## 接口完成检查清单

新增接口完成后逐项确认：

- [ ] 已定义独立 Request/Response DTO
- [ ] 参数标签和业务校验位置正确
- [ ] Router、Controller、Service、Store 调用链完整
- [ ] 数据库操作携带 Context
- [ ] 底层错误已映射为统一业务错误码
- [ ] 响应未包含密码或内部字段
- [ ] 日志未包含敏感数据
- [ ] API 路径、JSON 字段和 Go 命名符合规范

## 后续约定预留

<!--
### 约定名称

- 适用场景：
- 强制要求：
- 推荐写法：
- 禁止写法：
- 参考文件：
-->

> **待补充：** 查询列表、详情、更新、删除、事务、权限校验及审计日志规范。

## OAuth 登录接口

`GET /ops/login/generic_oauth` 同时承担登录入口和 Keycloak 回调：

- 首次访问不含 `code`、`state`、`error`，后端生成 32 字节随机 `state`，写入 `HttpOnly`、`SameSite=Lax` Cookie，并返回 `302 Location` 跳转 Keycloak。
- Keycloak 回调携带 `code` 和 `state`；浏览器自动携带原 Cookie。
- Controller 使用常量时间比较 URL 与 Cookie 中的 `state`。无论成功或失败，Cookie 都必须立即删除，保证单次使用。
- `error` 回调也必须先校验 `state`，再返回认证失败。
- Keycloak 地址、Client ID、回调地址和 Cookie Secure 开关统一从 `oauth` 配置读取。
- Controller 不记录或输出 Cookie、Token 等敏感数据。

当前回调完成 `state` 校验后返回授权码。接入 Token Endpoint 后，应改为由服务端直接使用授权码换取 Token，不再把授权码返回前端。

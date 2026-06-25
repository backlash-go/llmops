# Repository Architecture Guidelines

## 目标与适用范围

本文定义本项目后续开发应遵守的通用架构规则，覆盖 HTTP API、配置初始化、中间件、认证登录、业务分层、数据访问、错误处理和编码约定。

新增业务能力时优先保持以下调用链：

```text
Command/App → Options/Config → Server → Router → Controller → Service → Store → Database/External System  API DTO   ↘ Model
```

各层只处理自己的职责。禁止跨层直接访问，例如 Controller 不直接访问数据库，Router 不处理业务逻辑，Store 不构造 HTTP 响应。

## 顶层目录职责

- `cmd/`：程序入口，只负责启动应用，不写逻辑。
- `configs/llmops-apiserver.yaml`：配置文件，字段必须能映射到 `options`。
- `internal/apiserver/`：API Server 主体，包含启动、路由、Controller、Service、Store。
- `internal/pkg/`：仅供本项目内部使用的公共能力，如错误码、中间件、配置项、server 封装、数据库模型。
- `pkg/`：可被外部复用的公共包，如 API DTO、日志、应用框架、存储客户端、工具函数。
- `api/`：协议定义或生成产物。
- `docs/`、`tasks/`、`design/`：文档任务和需求描述

`internal` 下的包不应被项目外部依赖；`pkg` 下的包应保持更稳定、通用、低业务耦合。

## 启动与配置规则

启动链路应保持清晰：

```text
cmd/<app>/main.go → internal/<app>/app.go → options.NewOptions() → config.CreateConfigFromOptions()
  → Run(cfg)
  → create server
  → PrepareRun()
  → Run()
```

配置规则：

- 新配置定义在 `internal/pkg/options/xxx.go` 字段必须使用 `mapstructure` 标签明确映射。
- 默认值在 `NewXxxOptions()` 中设置。
- 参数校验放在 `Validate()` 中
- `ApplyTo()` 只负责把 options 应用到运行配置，不写业务逻辑。
- 新配置 放到 s *apiServer  在 PrepareRun()  实现 s.initxxx() 方法

。

## Server 与中间件规则

通用 HTTP Server 由 `internal/pkg/server` 负责创建和运行，业务 Server 只负责组装配置与注册业务路由。

中间件配置链路：



新增中间件时：

- 在internal/pkg/middleware 新创建的增加一个文件  在函数defaultMiddlewares 映射 配置文件里面 可选择配置 是否注入

## HTTP API 分层规则

标准资源接口应遵守：

```text
Router → Controller → Service → Store → MySQL
                    ↘ API DTO   ↘ Model
```

各层职责：

- Router：只注册路径、HTTP 方法和 Controller，不读取数据库，不处理业务逻辑。
- Controller：绑定请求、触发参数校验、调用 Service、统一写响应。
- Service：执行业务规则、DTO/Model 转换、调用 Store、映射业务错误。
- Store：只负责持久化或读取数据，透传底层错误，不构造 HTTP 响应。
- Model：描述数据库结构和 ORM 映射。
- API DTO：描述外部请求和响应，禁止直接暴露内部 Model。

新增资源时建议目录：

```text
internal/apiserver/router/<resource>.go
internal/apiserver/controller/v1/<resource>/
internal/apiserver/service/v1/<resource>.go
internal/apiserver/store/mysql/<resource>.go
pkg/api/llmops/v1/<resource>.go
internal/pkg/model/<resource>.go
```

路径规则：

- REST 资源路径使用复数、小写、短横线风格。
- API 版本前缀由上层统一提供，例如 `/ops/api/v1`。
- 登录、回调、健康检查、metrics 等非资源路径可以独立注册，但仍应保持 Controller 边界清晰。

## DTO 与 Model 规则

请求 DTO、响应 DTO、数据库 Model 必须分离。

- 请求 响应的 DTO 放在 `pkg/api/llmops/v1/`。
- 响应 DTO 不得直接复用含敏感字段的请求 DTO。
- 响应 DTO 不得包含密码、Token、Cookie、内部错误、软删除字段等敏感或内部字段。
- JSON 字段统一使用 `snake_case`。
- 参数校验优先使用 `binding` 标签。
- Model 放在 `internal/pkg/model/`，字段必须声明明确的 `gorm` 与 `json` 标签。
- Model 必须通过 `TableName()` 固定表名。





DTO 转 Model 可以使用 `copier.Copy`，但必须确认字段匹配

## Controller 规则

Controller 每个动作建议独立文件，例如 `create.go`、`get.go`、`update.go`、`delete.go`。

Controller 标准流程：

```text
记录必要日志
  → ShouldBindJSON / ShouldBindQuery / ShouldBindUri
  → 绑定或校验失败立即返回
  → 调用 Service
  → core.WriteResponse(c, err, data)
  → 写错误响应后立即 return
```

Controller 不应：

- 直接访问数据库；
- 构造 SQL 或 GORM 查询；
- 写复杂业务规则；
- 重复 DTO binding 已能完成的必填校验；
- 返回包含敏感字段的请求对象；


## Service 规则

Service 接口使用 API DTO 作为输入，使用响应 DTO 或错误作为输出。

Service 职责包括：

- DTO 与 Model 转换；
- 执行业务规则；
- 调用 Store；
- 编排事务；
- 将底层错误映射为 `internal/pkg/code` 中的统一业务错误；



## Store 与数据库规则

Store 接口定义资源级数据访问能力，MySQL 实现放在 `internal/apiserver/store/mysql/`。

所有数据库操作必须携带 `context.Context`：

```go
return u.db.WithContext(ctx).Create(user).Error
```

Store 只负责数据访问：

- 不构造 HTTP 响应；
- 不依赖 Gin；
- 不写业务错误码；
- 不记录敏感 SQL 参数；
- 不吞掉底层错误。


## 错误码规则

统一错误码放在 `internal/pkg/code/`。

新增错误码时：

- 按现有分段规则定义；
- 注册 HTTP 状态码和外部可见文案；
- Service 负责把底层错误映射成业务错误；
- Controller 只负责把错误交给 `core.WriteResponse`。


## 编码约定

- Go 文件必须符合 `gofmt`/`goimports`。
- 本项目导入路径 `llmops/...` 单独分组。
- 导出标识符使用 `PascalCase`，内部标识符使用 `camelCase`。
- 包名简短、小写、无下划线。
- 接口使用职责名称，如 `UserSrv`、`UserStore`。
- 实现使用非导出类型，如 `userService`、`users`。
- 方法接收者保持简短且一致，如 `u *UserController`、`u *userService`。
- 注释说明“做什么”，避免重复代码本身。
- 导出类型、导出方法应有完整注释。
- 不忽略可能失败的错误。若确实忽略，必须写明原因。






```
cmd/
├── llmops-apiserver/   服务端入口
└── llmopsctl/          命令行客户端

internal/apiserver/
├── controller/v1/      Gin HTTP 控制器
├── service/v1/         业务逻辑
├── store/              存储接口
│   ├── mysql/          MySQL 实现
│   ├── etcd/           etcd 实现
│   └── fake/           测试内存实现
├── options/            配置项
├── router.go           路由注册
└── server.go           HTTP、gRPC、MySQL、Redis 初始化

internal/pkg/
├── middleware/         认证、校验等中间件
├── options/            通用配置结构
└── server/             Gin HTTP 通用服务器

pkg/
├── app/                Cobra + Viper 应用框架
├── db/                 数据库连接
├── log/                日志系统
├── shutdown/           优雅关闭
└── storage/            Redis
```


```
HTTP 请求
↓
Gin Router
↓
Controller
↓
Service
↓
Store Interface
↓
MySQL / etcd / fake
```
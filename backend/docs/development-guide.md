# 开发指南

本文档提供了 WaitingToDo 项目的开发环境搭建、开发流程和编码规范的说明。
## 1. 项目概览与目录结构

```
WaitingToDo/
├─ backend/
│  ├─ cmd/                  # 程序入口（main.go）
│  ├─ config/               # 应用配置（config.yaml）
│  ├─ docs/                 # 文档（开发/部署/API）
│  ├─ internal/
│  │  ├─ configs            # （可能包含配置初始化）
│  │  ├─ handlers           # 控制器/处理器
│  │  ├─ middlewares        # 中间件（含 JWT、安全、限流）
│  │  ├─ models             # 模型
│  │  ├─ repositories       # 数据访问层
│  │  ├─ routers            # 路由注册（Gin）
│  │  └─ services           # 业务服务（含 MQ 消费者、定时任务）
│  ├─ pkg/                  # 公共包（logger/redisContent/minioContent 等）
│  ├─ .air.toml             # Air 热重载配置
│  └─ Dockerfile            # 后端镜像构建
└─ frontend/vue/            # 前端（Vue3 + Vite）
   ├─ .env.development      # 开发环境变量
   ├─ .env.production       # 生产环境变量
   ├─ package.json          # 前端依赖与脚本
   └─ vite.config.js        # Vite 配置
```


## 2. 环境要求

- 后端：
  - Go 1.24+
  - MySQL 8.0+
  - Redis 6.0+
  - RabbitMQ 3.8+
  - MinIO（最新稳定版）
- 前端：
  - Node.js 18+
  - npm 9+（或使用 pnpm/yarn 亦可）

## 3. 必备依赖服务准备

请确保以下服务在本地可用，并与后端配置保持一致（见 backend/config/config.yaml）：
- MySQL：默认示例 DSN 为 user:password@tcp(127.0.0.1:3306)/WaitingToDo
- Redis：127.0.0.1:6379（可设密码）
- RabbitMQ：amqp://guest:guest@localhost:5672/
- MinIO：Endpoint 127.0.0.1:9000（Admin 账号建议本地自定义）

MinIO 启动后请在控制台创建用于存储图片/文件的 bucket，并配置合适的访问策略（开发环境可先开放读权限以便调试，生产请务必收紧权限）。

## 4. 后端开发与运行

### 4.1 获取与安装依赖

```bash
# Windows PowerShell
cd backend
# 安装依赖
go mod tidy
```

### 4.2 配置后端
- mysql.dsn：数据库连接串（库名建议 WaitingToDo）
- redis.addr / redis.password
- rabbitmq.dsn：队列连接（项目会启动好友/团队相关消费者）
- minio.endpoint / access_key_id / secret_access_key
- log.*：日志路径与级别

提示：首次运行前请在 MySQL 创建数据库（若使用自动建表，首次运行应用会创建表结构）。

### 4.3 启动方式（二选一）

- 直接运行：
```bash
# 需先启动 MySQL / Redis / RabbitMQ / MinIO
cd backend
go run cmd/main.go
```

- 热重载（推荐，已提供 .air.toml）：
```bash
# 一次性安装
go install github.com/cosmtrek/air@latest
# 在 backend 目录运行
air
```

运行后默认监听 8080 端口。应用启动时会：
- 注册路由与中间件（CORS、Security、限流、JWT）
- 启动 MQ 消费者（好友/团队相关）
- 启动定时任务服务（TickerNotify）

### 4.4 常用开发命令

```bash
# 代码格式化/静态检查（建议）
go fmt ./...
go vet ./...

# 运行单元测试（若后续补充测试）
go test ./...
```

## 5. 前端开发与运行

### 5.1 安装依赖

```bash
cd frontend/vue
npm install
```

### 5.2 开发环境变量

```
VITE_APP_API_BASE_URL=http://localhost:8080
VITE_PIC_BASE_URL=http://127.0.0.1:9000
```

### 5.3 本地启动

```bash
npm run dev
# Vite 默认 http://localhost:5173
```

前端将通过 VITE_APP_API_BASE_URL 访问后端接口。除 /auth 外的接口均需携带 JWT：
```
Authorization: Bearer <token>
```

## 6. 接口与鉴权约定（开发联调必读）

- 放行路由：/auth/login、/auth/register、/auth/forget、/auth/captcha（免 JWT）
- 受保护路由：/user、/upload、/task、/friend、/message（需 JWT）

## 7. 日志与排查

- 支持文件大小/切分/压缩/控制台输出等。
- 常见问题：
  - 启动报 MQ/Redis/MinIO 连接失败：请先确保依赖服务已启动且账号/地址一致。
  - 接口 401：检查是否调用了保护路由且未携带 Authorization 头，或 Token 过期。
  - 跨域失败：核对前端访问来源是否在 CORS AllowOrigins 列表中。
  - 静态资源/图片无法访问：核对 MinIO endpoint 与前端 VITE_PIC_BASE_URL 一致，Bucket/策略正确。

## 8. 代码规范与提交

- 后端：
  - Go 官方风格 + go fmt/go vet，建议引入 golangci-lint（可选）。
- 前端：
  - 建议使用 ESLint + Prettier（按团队规范落地）。
- 提交信息遵循 Conventional Commits：
```
<type>(<scope>): <description>
```
- **feat**: 新功能
- **fix**: 修复 Bug
- **docs**: 文档变更
- **style**: 代码风格调整
- **refactor**: 代码重构
- **test**: 测试相关变更
- **chore**: 构建过程或辅助工具的变动

## 9. 开发流程建议

1) 创建/更新配置（config.yaml）并启动依赖服务
2) 启动后端（Air 热重载或 go run）
3) 启动前端（npm run dev）
4) 通过 /auth/login 获取 token，完成后续受保护接口调试
5) 若需新增接口：在 routers 下注册路由 -> handlers/ services/ repositories 衔接实现 -> 自测与联调

---
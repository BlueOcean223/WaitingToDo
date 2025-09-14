# WaitingToDo

[English](./backend/docs/README_EN.md) | 中文

一个功能丰富的现代化待办事项管理平台，支持个人任务管理、团队协作、好友系统等多种功能。

## 📖 项目简介

WaitingToDo 是一个基于 Vue 3 + Go 开发的全栈待办事项管理应用。它不仅提供了传统的个人任务管理功能，还集成了团队协作、好友系统、实时消息通知、文件管理等现代化功能，旨在为用户提供一个高效、便捷的任务管理解决方案。

## ✨ 主要特性

### 🔐 用户系统
- 用户注册、登录、密码重置
- JWT 身份认证
- 用户资料管理
- 头像上传

### 📝 任务管理
- 创建、编辑、删除任务
- 任务状态管理（待办/已完成）
- 任务截止日期提醒
- 任务标签分类

### 👥 团队协作
- 创建和管理团队
- 邀请码加入团队
- 协作任务进度跟踪

### 🤝 好友系统
- 添加好友
- 好友请求管理
- 好友任务分享

### 🔔 消息通知
- 实时消息推送
- 任务提醒通知
- 团队协作通知
- 好友请求通知

### 📁 文件管理
- 任务附件上传
- 文件预览
- 文件下载

### 📊 数据统计
- 任务完成率统计
- 个人效率分析
- 团队协作数据
- 可视化图表展示


## 🚀 快速开始指南

### 环境要求

#### 后端环境
- Go 1.19+
- MySQL 8.0+
- Redis 6.0+
- MinIO（对象存储）
- RabbitMQ（消息队列）

#### 前端环境
- Node.js 16+
- npm 8+ 或 yarn 1.22+

### 安装步骤

#### 1. 克隆项目
```bash
git clone https://github.com/yourusername/WaitingToDo.git
cd WaitingToDo
```

#### 2. 后端设置

```bash
# 进入后端目录
cd backend

# 安装依赖
go mod download

# 复制配置文件
cp config/config.example.yaml config/config.yaml

# 编辑配置文件，设置数据库连接等信息
vim config/config.yaml

# 运行数据库迁移
go run main.go migrate

# 启动后端服务
go run main.go
```

#### 3. 前端设置

```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install
# 或使用 yarn
yarn install

# 启动开发服务器
npm run dev
# 或使用 yarn
yarn dev
```

#### 4. 数据库设置

```sql
-- 创建数据库
CREATE DATABASE waitingtodo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户（可选）
CREATE USER 'waitingtodo'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON waitingtodo.* TO 'waitingtodo'@'localhost';
FLUSH PRIVILEGES;
```

#### 5. 使用 Docker（推荐）

```bash
# 使用 Docker Compose 一键启动
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 访问应用

- 前端应用：http://localhost:3000
- 后端 API：http://localhost:8080
- API 文档：http://localhost:8080/swagger/index.html

## 📚 使用说明

### 基本使用流程

1. **注册账号**：访问应用首页，点击注册按钮创建新账号
2. **登录系统**：使用注册的邮箱和密码登录
3. **创建任务**：点击"添加任务"按钮，填写任务信息
4. **管理任务**：在任务列表中查看、编辑、完成或删除任务
5. **团队协作**：创建团队，邀请成员，分配团队任务
6. **添加好友**：搜索用户，发送好友请求，与好友分享任务

### 高级功能

- **任务筛选**：使用状态、优先级、标签等条件筛选任务
- **批量操作**：选择多个任务进行批量删除或状态更新
- **数据导出**：导出任务数据为 CSV 或 PDF 格式
- **API 集成**：使用 RESTful API 与第三方应用集成

### 移动端使用

应用采用响应式设计，在移动设备上也能获得良好的使用体验：

- 触摸友好的界面设计
- 手势操作支持
- 离线数据缓存
- 推送通知支持

## 🛠️ 技术栈

### 前端技术
- **框架**：Vue 3 + TypeScript
- **构建工具**：Vite
- **状态管理**：Pinia
- **路由**：Vue Router 4
- **UI 组件**：Element Plus
- **样式**：SCSS
- **HTTP 客户端**：Axios
- **图表**：ECharts

### 后端技术
- **语言**：Go 1.19+
- **框架**：Gin
- **数据库**：MySQL 8.0
- **缓存**：Redis
- **对象存储**：MinIO
- **消息队列**：RabbitMQ
- **认证**：JWT
- **API 文档**：Swagger

### 开发工具
- **版本控制**：Git
- **容器化**：Docker + Docker Compose
- **代码质量**：ESLint + Prettier（前端），golangci-lint（后端）
- **测试**：Jest（前端），Go testing（后端）

## 📁 项目结构

```
WaitingToDo/
├── frontend/                 # 前端项目
│   ├── src/
│   │   ├── components/       # 组件
│   │   ├── views/           # 页面
│   │   ├── store/           # 状态管理
│   │   ├── router/          # 路由配置
│   │   ├── api/             # API 接口
│   │   └── utils/           # 工具函数
│   ├── public/              # 静态资源
│   └── package.json         # 依赖配置
├── backend/                 # 后端项目
│   ├── api/                 # API 路由
│   ├── config/              # 配置文件
│   ├── internal/            # 内部模块
│   │   ├── handler/         # 处理器
│   │   ├── service/         # 业务逻辑
│   │   ├── repository/      # 数据访问
│   │   └── model/           # 数据模型
│   ├── pkg/                 # 公共包
│   └── main.go              # 入口文件
├── docs/                    # 项目文档
├── docker-compose.yml       # Docker 编排
└── README.md               # 项目说明
```

## 🤝 贡献指南

我们欢迎所有形式的贡献！无论是报告 bug、提出新功能建议，还是提交代码改进。

### 如何贡献

1. **Fork 项目**
   ```bash
   # 点击 GitHub 页面右上角的 Fork 按钮
   ```

2. **创建功能分支**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **提交更改**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

4. **推送到分支**
   ```bash
   git push origin feature/your-feature-name
   ```

5. **创建 Pull Request**
   - 在 GitHub 上创建 Pull Request
   - 详细描述你的更改
   - 等待代码审查

### 代码规范

#### 前端代码规范
- 使用 ESLint + Prettier 进行代码格式化
- 遵循 Vue 3 Composition API 最佳实践
- 组件命名使用 PascalCase
- 文件命名使用 kebab-case

#### 后端代码规范
- 遵循 Go 官方代码规范
- 使用 golangci-lint 进行代码检查
- 函数和变量命名使用驼峰命名法
- 包名使用小写字母

#### 提交信息规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
type(scope): description

[optional body]

[optional footer]
```

类型说明：
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

### 报告问题

如果你发现了 bug 或有功能建议，请：

1. 检查 [Issues](https://github.com/yourusername/WaitingToDo/issues) 中是否已有相关问题
2. 如果没有，创建新的 Issue
3. 详细描述问题或建议
4. 如果是 bug，请提供复现步骤


---

⭐ 如果这个项目对你有帮助，请给我们一个 Star！

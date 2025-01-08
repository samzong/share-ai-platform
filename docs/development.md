# Share AI Platform 开发指南

## 目录
- [开发环境要求](#开发环境要求)
- [项目设置](#项目设置)
- [开发工作流](#开发工作流)
- [常用命令](#常用命令)
- [Docker 开发环境](#docker-开发环境)
- [故障排除](#故障排除)

## 开发环境要求

- Node.js 18+
- Go 1.21+
- Docker & Docker Compose
- Make

## 项目设置

1. 克隆项目
```bash
git clone https://github.com/samzong/share-ai-platform.git
cd share-ai-platform
```

2. 安装依赖
```bash
# 安装所有依赖
make install

# 或者分别安装
make install-frontend  # 只安装前端依赖
make install-backend   # 只安装后端依赖
```

## 开发工作流

### 1. 启动开发环境

完整开发环境（包含数据库）：
```bash
make dev
```

分别启动各个服务：
```bash
# 只启动数据库服务（PostgreSQL 和 Redis）
make db

# 只启动后端服务
make backend

# 只启动前端开发服务器
make frontend
```

### 2. 开发过程

前端开发：
- 访问 http://localhost:3000
- 修改 `frontend/src` 下的文件会自动热重载
- 前端 API 请求默认代理到 http://localhost:8080/api

后端开发：
- API 服务运行在 http://localhost:8080
- 修改 Go 文件会自动重新编译和重启服务
- API 文档访问 http://localhost:8080/swagger/index.html

数据库：
- PostgreSQL 运行在 localhost:5432
- Redis 运行在 localhost:6379

### 3. 构建项目

```bash
# 构建整个项目
make build

# 分别构建
make build-frontend  # 只构建前端
make build-backend   # 只构建后端
```

### 4. 运行测试

```bash
# 运行所有测试
make test

# 分别测试
make test-frontend  # 只运行前端测试
make test-backend   # 只运行后端测试
```

## 常用命令

### 开发相关命令

```bash
# 启动完整开发环境
make dev

# 只启动数据库服务
make db

# 只启动后端服务
make backend

# 只启动前端服务
make frontend

# 启动前后端（不包含数据库）
make start
```

### 构建相关命令

```bash
# 构建整个项目
make build

# 构建 Docker 镜像
make docker

# 清理构建文件
make clean
```

### 停止服务

```bash
# 停止所有服务
make stop

# 分别停止服务
make stop-frontend  # 停止前端
make stop-backend   # 停止后端
make stop-services  # 停止数据库服务
```

## Docker 开发环境

项目提供了完整的 Docker 开发环境，包含热重载支持：

```bash
# 启动 Docker 开发环境
make docker-dev

# 启动生产环境
make prod
```

Docker 开发环境特点：
- 前端热重载
- 后端热重载（使用 air）
- 数据持久化
- 环境隔离
- 统一的开发环境

## 故障排除

1. 端口占用问题
```bash
# 停止所有服务
make stop

# 检查端口占用
lsof -i :3000  # 检查前端端口
lsof -i :8080  # 检查后端端口
lsof -i :5432  # 检查 PostgreSQL 端口
lsof -i :6379  # 检查 Redis 端口
```

2. 数据库连接问题
```bash
# 重启数据库服务
make stop-services
make db
```

3. 依赖问题
```bash
# 清理并重新安装依赖
make clean
make install
```

4. Docker 问题
```bash
# 清理 Docker 资源
docker-compose down -v
docker system prune -f
```

## 项目结构

```
.
├── frontend/                # 前端项目目录
│   ├── src/
│   │   ├── components/     # 可复用组件
│   │   ├── pages/         # 页面组件
│   │   ├── services/      # API 服务
│   │   ├── utils/         # 工具函数
│   │   ├── hooks/         # 自定义 Hooks
│   │   └── styles/        # 样式文件
│   └── public/            # 静态资源
│
├── backend/                # 后端项目目录
│   ├── cmd/               # 主程序入口
│   ├── internal/          # 内部包
│   │   ├── api/          # API 处理器
│   │   ├── middleware/   # 中间件
│   │   ├── models/       # 数据模型
│   │   ├── services/     # 业务逻辑
│   │   └── database/     # 数据库操作
│   ├── configs/          # 配置文件
│   ├── scripts/          # 脚本文件
│   └── docs/             # 文档
│
└── deploy/               # 部署相关文件 
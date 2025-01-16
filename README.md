# Share AI Platform

开源的 AI 镜像分享平台，致力于简化 AI 应用的部署和分享过程。

## 功能特点

- 镜像管理：浏览、搜索和收藏 AI 相关镜像
- 用户系统：支持用户注册、登录和个人信息管理
- 一键部署：快速部署 AI 应用到支持的云平台
- 分类系统：预留算法、模型、数据集等分类入口

## 技术栈

### 前端
- React
- Ant Design
- Context API
- Axios

### 后端
- Go (Golang)
- Gin Framework
- PostgreSQL
- Redis

## 快速开始

### 前端开发
```bash
cd frontend
npm install
npm start
```

### 后端开发
```bash
cd backend
go mod tidy
go run cmd/main.go
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


## Issue

### DB Migration

```bash
cd backend
make migrate # init db
```


### run dev

```bash
# at git root path
make dev
```

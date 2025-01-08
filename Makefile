.PHONY: help install start build clean test dev prod docker docker-dev frontend backend stop db

# 默认目标
help:
	@echo "可用的命令："
	@echo "make install     - 安装前后端依赖"
	@echo "make frontend   - 只启动前端开发服务器"
	@echo "make backend    - 只启动后端服务"
	@echo "make start      - 启动开发环境（前后端）"
	@echo "make build      - 构建前后端项目"
	@echo "make clean      - 清理构建文件"
	@echo "make test       - 运行测试"
	@echo "make dev        - 启动开发环境（包含数据库和Redis）"
	@echo "make db         - 只启动数据库服务"
	@echo "make prod       - 启动生产环境"
	@echo "make docker     - 构建Docker镜像"
	@echo "make docker-dev - 启动Docker开发环境"
	@echo "make stop       - 停止所有服务"

# 安装依赖
install-frontend:
	@echo "安装前端依赖..."
	cd frontend && npm install

install-backend:
	@echo "安装后端依赖..."
	cd backend && go mod tidy && go install github.com/cosmtrek/air@latest

install: install-frontend install-backend

# 启动前端
frontend:
	@echo "启动前端开发服务器..."
	cd frontend && npm start

# 启动后端
backend:
	@echo "启动后端服务..."
	cd backend && (command -v air >/dev/null 2>&1 && air -c .air.toml || go run cmd/main.go)

# 启动数据库服务
db:
	@echo "启动数据库和Redis..."
	docker-compose up -d postgres redis
	@echo "等待数据库启动..."
	sleep 5

# 启动开发环境
start:
	@echo "启动后端服务..."
	cd backend && go run cmd/main.go &
	@echo "启动前端开发服务器..."
	cd frontend && npm start

# 构建项目
build-frontend:
	@echo "构建前端..."
	cd frontend && npm run build

build-backend:
	@echo "构建后端..."
	cd backend && go build -o bin/main cmd/main.go

build: build-frontend build-backend

# 清理构建文件
clean-frontend:
	@echo "清理前端构建文件..."
	rm -rf frontend/build
	rm -rf frontend/node_modules

clean-backend:
	@echo "清理后端构建文件..."
	rm -rf backend/bin
	rm -rf backend/vendor

clean: clean-frontend clean-backend

# 运行测试
test-frontend:
	@echo "运行前端测试..."
	cd frontend && npm test

test-backend:
	@echo "运行后端测试..."
	cd backend && go test ./...

test: test-frontend test-backend

# 启动开发环境（包含数据库和Redis）
dev: db
	@echo "启动后端服务..."
	cd backend && go run cmd/main.go &
	@echo "启动前端开发服务器..."
	cd frontend && npm start

# 启动生产环境
prod:
	docker-compose up --build -d

# 构建Docker镜像
docker:
	docker build -t share-ai-platform .

# 启动Docker开发环境
docker-dev:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build

# 停止所有服务
stop-services:
	@echo "停止 Docker 服务..."
	docker-compose down

stop-backend:
	@echo "停止后端服务..."
	pkill -f "go run cmd/main.go" || true

stop-frontend:
	@echo "停止前端服务..."
	pkill -f "node" || true

stop: stop-services stop-backend stop-frontend
 
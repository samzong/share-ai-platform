.PHONY: dev build test clean db db-stop migrate swagger

# 开发模式，使用 air 热重载
dev:
	air

# 构建项目
build:
	go build -o bin/main cmd/main.go

# 运行测试
test:
	go test ./... -v

# 清理临时文件
clean:
	rm -rf tmp/
	rm -rf bin/

# 初始化开发环境
init:
	go mod download
	go mod tidy

# 启动数据库
db:
	docker-compose up -d postgres redis

# 停止数据库
db-stop:
	docker-compose stop postgres redis

# 运行数据库迁移
migrate:
	go run cmd/migrate/main.go

# 生成 Swagger 文档
swagger:
	swag init -g cmd/main.go -o docs

.DEFAULT_GOAL := dev 
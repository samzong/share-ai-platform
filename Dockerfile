# 构建前端
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ .
RUN npm run build

# 构建后端
FROM golang:1.21-alpine AS backend-builder
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

# 最终镜像
FROM alpine:latest
WORKDIR /app

# 复制前端构建产物
COPY --from=frontend-builder /app/frontend/build /app/frontend/build

# 复制后端二进制文件和配置
COPY --from=backend-builder /app/backend/main /app/backend/
COPY --from=backend-builder /app/backend/configs /app/backend/configs

EXPOSE 8080
EXPOSE 3000

CMD ["/app/backend/main"] 
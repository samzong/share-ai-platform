# 构建前端
FROM node:18-alpine as frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# 构建后端
FROM golang:1.21-alpine as backend-builder
WORKDIR /app
COPY backend/ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# 最终镜像
FROM alpine:latest
WORKDIR /app
COPY --from=backend-builder /app/main .
COPY --from=backend-builder /app/configs ./configs
COPY --from=frontend-builder /app/frontend/build ./frontend/build

EXPOSE 8080
CMD ["./main"] 
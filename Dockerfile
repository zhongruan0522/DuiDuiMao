# ===== 阶段1: 构建前端 =====
FROM node:20-alpine AS frontend-builder

WORKDIR /app/web

# 复制前端依赖文件
COPY web/package*.json ./

# 安装前端依赖
RUN npm install

# 复制前端源代码
COPY web/ ./

# 构建前端
RUN npm run build


# ===== 阶段2: 构建后端 =====
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app

# 复制Go依赖文件
COPY go.mod go.sum ./

# 下载Go依赖
RUN go mod download

# 复制后端源代码
COPY cmd/ ./cmd/
COPY internal/ ./internal/

# 构建后端
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server


# ===== 阶段3: 运行时镜像 =====
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件和静态文件
COPY --from=backend-builder /app/server .
COPY --from=frontend-builder /app/web/dist ./web/dist

# 复制配置文件模板
COPY config.example.yaml ./

# 创建临时目录（dev模式CSV存储）
RUN mkdir -p Temp

# 暴露端口
EXPOSE 3001

# 启动服务
CMD ["./server"]

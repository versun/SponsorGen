FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制go模块文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -o sponsorgen .

# 使用精简的alpine镜像
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata imagemagick librsvg && \
    mkdir -p /app/output /app/cache /app/assets

WORKDIR /app

# 复制编译好的二进制文件
COPY --from=builder /app/sponsorgen .
COPY assets ./assets

# 暴露默认端口
EXPOSE 5000

# 设置默认环境变量
ENV OUTPUT_DIR="/app/output" \
    CACHE_DIR="/app/cache" \
    DEFAULT_AVATAR="/app/assets/default_avatar.svg" \
    REFRESH_MINUTES="60"

# 启动服务
CMD ["./sponsorgen", "-port", "5000"]
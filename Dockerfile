# --- 第一阶段：构建可执行文件 ---
FROM golang:1.23 AS builder
ARG service_name

# 设置工作目录
WORKDIR /app

# 1. 复制模块定义，下载依赖
COPY go.mod go.sum ./
RUN go env -w GOPROXY="https://goproxy.cn,direct" \
 && go mod download

# 2. 复制指定微服务目录下的源码
#    这样可以利用上一步的依赖缓存
COPY microservice/${service_name} ./microservice/${service_name}

# 3. 切到服务目录并编译
WORKDIR /app/microservice/${service_name}
RUN go build -o /app/bin/${service_name} .

# --- 第二阶段：运行时镜像 ---
FROM debian:bookworm-slim
ARG service_name

# 安装认证、时区
RUN apt-get update && apt-get install -y \
    tzdata ca-certificates \
 && rm -rf /var/lib/apt/lists/* \
 && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
 && echo "Asia/Shanghai" > /etc/timezone

# 拷入可执行文件
COPY --from=builder /app/bin/${service_name} /app/${service_name}
RUN chmod +x /app/${service_name}

# 最终工作目录
WORKDIR /app

# 启动时直接执行对应二进制
ENTRYPOINT ["./"]
CMD ["./${service_name}"]

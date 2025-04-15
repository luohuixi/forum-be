# forum-be
木犀论坛后端仓库

mysql

websocket: 聊天

拆分成微服务microservice，使用grpc（protobuf）通信：
- gateway
- post
- user
- chat
- feed

casbin: 鉴权

redis: 点赞; 聊天; tag缓存; tag排行榜; 消息队列(pub-sub); hot post

log: zap

CI: .drone.yml

Kafka: 动态功能的消息队列

# 部署

在根目录下对每个微服务进行docker build即可

```
docker buildx build --platform linux/amd64 -t registry.cn-shenzhen.aliyuncs.com/muxi/forum-backend:$(tag) .
```

相关中间件:

1. redis
2. etcd
3. kafka
4. mysql



注册中心的相关配置是从环境变量中读取的

```
MICRO_REGISTRY_ADDRESS=localhost:2379
```

生成相关的proto代码需要先下载如下插件
```
go install github.com/go-micro/generator/cmd/protoc-gen-micro@latest
```

生成proto代码
```
protoc --proto_path=. --go_out=:. --micro_out=. path/your.proto
```
module forum-chat

replace forum => ../../

go 1.16

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/micro/go-micro => github.com/Lofanmi/go-micro v1.16.1-0.20210804063523-68bbf601cfa4 // to go 1.16

require (
	forum v0.0.0-00010101000000-000000000000
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins/wrapper/trace/opentracing v0.0.0-20200119172437-4fe21aa238fd
	github.com/opentracing/opentracing-go v1.2.0
	github.com/spf13/viper v1.12.0
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	golang.org/x/net v0.0.0-20220520000938-2e3eb7b945c2
)

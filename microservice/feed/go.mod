module forum-feed

replace forum => ../../

replace forum-user => ../user

go 1.16

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/micro/go-micro => github.com/Lofanmi/go-micro v1.16.1-0.20210804063523-68bbf601cfa4 // to go 1.16

require (
	forum v0.0.0-00010101000000-000000000000
	forum-user v0.0.0-00010101000000-000000000000
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/jinzhu/gorm v1.9.16
	github.com/micro/cli v0.2.0
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins/registry/kubernetes v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/go-plugins/wrapper/trace/opentracing v0.0.0-20200119172437-4fe21aa238fd
	github.com/opentracing/opentracing-go v1.2.0
	github.com/spf13/viper v1.12.0
	golang.org/x/net v0.0.0-20220812174116-3211cb980234
)

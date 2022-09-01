module feed-client

replace forum => ../../

replace forum-feed => ../../microservice/feed

replace forum-user => ../../microservice/user

go 1.16

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/micro/go-micro => github.com/Lofanmi/go-micro v1.16.1-0.20210804063523-68bbf601cfa4 // to go 1.16

require (
	forum v0.0.0-00010101000000-000000000000
	forum-feed v0.0.0-00010101000000-000000000000
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins/wrapper/trace/opentracing v0.0.0-20200119172437-4fe21aa238fd
	github.com/opentracing/opentracing-go v1.2.0
)

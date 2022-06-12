module forum-user

replace forum => ../../

go 1.16

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/micro/go-micro => github.com/Lofanmi/go-micro v1.16.1-0.20210804063523-68bbf601cfa4 // to go 1.16

require (
	forum v0.0.0-00010101000000-000000000000
	github.com/ShiinaOrez/GoSecurity v0.0.0-20191118072239-d06064a9edd6
	github.com/golang/protobuf v1.5.2
	github.com/jinzhu/gorm v1.9.16
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins/wrapper/trace/opentracing v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/protobuf v0.0.0-20180321161605-ebd3be6d4fdb // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/satori/go.uuid v1.2.0
	github.com/smartystreets/goconvey v1.7.2
	github.com/spf13/viper v1.11.0
	golang.org/x/net v0.0.0-20220421235706-1d1ef9303861
	google.golang.org/grpc v1.46.0 // indirect
)

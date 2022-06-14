module github.com/forum-be

replace forum => ./

go 1.16

require (
	forum v0.0.0-00010101000000-000000000000
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/casbin/casbin/v2 v2.47.2
	github.com/casbin/gorm-adapter/v3 v3.7.2
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.6.0
	github.com/jinzhu/gorm v1.9.16
	github.com/micro/go-micro v1.18.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/smartystreets/goconvey v1.7.2
	github.com/spf13/viper v1.12.0
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.uber.org/zap v1.21.0
)

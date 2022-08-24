package main

import (
	"flag"
	"github.com/micro/go-micro"
	"github.com/opentracing/opentracing-go"
	"log"

	"forum-feed/dao"
	pb "forum-feed/proto"
	s "forum-feed/service"
	"forum/config"
	logger "forum/log"
	"forum/pkg/handler"
	"forum/pkg/tracer"

	_ "github.com/micro/go-plugins/registry/kubernetes"

	"github.com/micro/cli"
	opentracingWrapper "github.com/micro/go-plugins/wrapper/trace/opentracing"
	"github.com/spf13/viper"
)

// 使用--sub运行subscribe服务
// 否则默认运行feed服务
var subFg = flag.Bool("sub", false, "use subscribe service mode")

// InitRpcClient ... 记录 bug: go-micro 框架的命令行参数和 go 标准库的 flag 冲突了
func initRpcClient() {
	s.UserInit()
}

// 包含两个服务：feed服务和subscribe服务
// subscribe服务 --> 异步将feed数据写入数据库
func main() {
	flag.Parse()

	var err error

	// init config
	if !*subFg {
		// feed-service
		initRpcClient()
		err = config.Init("./conf/config.yaml", "FORUM_FEED")
	} else {
		// sub-service
		err = config.Init("./conf/config_sub.yaml", "FORUM_SUB")
	}

	if err != nil {
		panic(err)
	}

	t, io, err := tracer.NewTracer(viper.GetString("local_name"), viper.GetString("tracing.jager"))
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	defer logger.SyncLogger()

	dao.Init()

	// set var t to Global Tracer (opentracing single instance mode)
	opentracing.SetGlobalTracer(t)

	if *subFg {
		// sub-service

		logger.Info("Subscribe service start...")
		s.SubServiceRun()
		return
	}

	// feed-service

	srv := micro.NewService(
		micro.Name(viper.GetString("local_name")),
		micro.WrapHandler(
			opentracingWrapper.NewHandlerWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapHandler(handler.ServerErrorHandlerWrapper()),
		micro.Flags(cli.BoolFlag{
			Name:   "sub",
			Usage:  "use subscribe service mode",
			Hidden: false,
		}),
	)

	// Init will parse the command line flags.
	srv.Init()

	pb.RegisterFeedServiceHandler(srv.Server(), s.New(dao.GetDao()))

	// Run the server
	if err := srv.Run(); err != nil {
		logger.Error(err.Error())
	}
}

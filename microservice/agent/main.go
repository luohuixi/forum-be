package main

import (
	"forum-agent/async"
	"forum-agent/core"
	"forum-agent/dao"
	pb "forum-agent/proto"
	"forum-agent/service"
	"log"

	"forum-agent/tool"
	"forum/client"
	"forum/config"
	logger "forum/log"
	"forum/pkg/handler"
	"forum/pkg/identity"
	"forum/pkg/tracer"

	"github.com/go-micro/plugins/v4/registry/etcd"
	opentracingWrapper "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	if err := config.Init("conf/agent.yaml", "FORUM_AGENT"); err != nil {
		panic(err)
	}
	core.ChatModelInit()
	core.EmbedderInit()
	core.ESClientInit()
	core.KafkaInit()
	core.AgentInit(tool.ToolList())
	dao.Init()
	worker := async.NewAsyncManager(dao.GetDao(), core.GetReActAgent())
	worker.Begin()

	traceAddr := "http://" + configValue("tracing.jager") + "/api/traces"
	t, io, err := tracer.NewTracer(configValue("local_name"), traceAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := io.Close(); err != nil {
			log.Println("close tracer io failed:", err)
		}
	}()
	defer logger.SyncLogger()

	// set var t to Global Tracer (opentracing single instance mode)
	opentracing.SetGlobalTracer(t)
	r := etcd.NewRegistry(
		registry.Addrs(configValue("etcd.addr")),
		etcd.Auth(configValue("etcd.username"), configValue("etcd.password")),
	)

	srv := micro.NewService(
		micro.Name(identity.Prefix()+configValue("local_name")),
		micro.WrapHandler(
			opentracingWrapper.NewHandlerWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapHandler(handler.ServerErrorHandlerWrapper()),

		micro.WrapClient(
			opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapCall(handler.ClientErrorHandlerWrapper()),

		micro.Registry(r),
	)

	// Init will parse the command line flags.
	srv.Init()
	client.UserInit(srv)
	dao.Init()

	// Register handler
	if err := pb.RegisterAgentServiceHandler(srv.Server(), service.New(dao.GetDao(), core.GetReActAgent())); err != nil {
		panic(err)
	}

	// Run the server
	if err := srv.Run(); err != nil {
		logger.Error(err.Error())
	}
}

func configValue(key string) string {
	return viper.GetString(key)
}

//const input = `[
//  {
//    "Id": 508,
//    "Content": "## **克隆仓库**
//git clone
//## **提交**
//1. git add
//> git add <文件名>：将指定的文件添加到暂存区。
//git add .：将所有修改过的文件添加到暂存区。
//git add -A：将所有修改过的文件和新文件（包括未跟踪的文件）添加到暂存区
//2. git commit -m
//> ***git commit*** ：这将打开文本编辑器，让你输入提交信息。完成信息编写后保存并关闭编辑器，提交就会完成。
//***git commit -m*** ：这是一种快速提交的方式，允许你直接在命令行中提供提交信息。例如，git commit -m "修复了登录功能的 bug" 会创建一个提交，其信息是“修复了登录功能的 bug”。
//git commit 只影响本地仓库，并不会更改远程仓库（如 GitHub 上的仓库）。要将这些更改推送到远程仓库，你需要使用 git push 命令。
//3. git push
//> ***git push <远程仓库名> <分支名>*** ：这个命令会将指定的本地分支推送到指定的远程仓库。例如，git push origin master 会将本地的 master 分支推送到名为 origin 的远程仓库。
//***git push*** ：如果已经设置了本地分支和远程分支之间的跟踪关系，可以直接使用这个命令来推送更改。Git 会自动推送到之前配置的远程分支。
//***git push -u <远程仓库名> <分支名>*** ：除了推送更改外，这个命令还会设置本地分支和远程分支之间的跟踪关系。在首次推送分支时常用这个命令。例如，git push -u origin feature 会将本地的 feature 分支推送到远程仓库，并设置跟踪关系。",
//    "Title": "git 基本操作"
//  }
//]`

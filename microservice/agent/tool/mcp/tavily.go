package mcp

import (
	"context"
	"os"
	"sync"
	"time"

	forumlog "forum/log"

	eino_mcp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

var (
	tavilyTools []tool.BaseTool
	tavilyOnce  sync.Once
	tavilyErr   error
)

func TavilyTools() ([]tool.BaseTool, error) {
	tavilyOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		tavilyURL := "https://mcp.tavily.com/mcp/?tavilyApiKey=" + os.Getenv("TAVILY_API_KEY")
		cli, err := client.NewStreamableHttpClient(tavilyURL)
		if err != nil {
			forumlog.Error("创建 MCP 客户端失败", forumlog.String(err.Error()))
			tavilyErr = err
			return
		}

		if err := cli.Start(ctx); err != nil {
			forumlog.Error("连接 MCP Server 失败", forumlog.String(err.Error()))
			tavilyErr = err
			return
		}

		initReq := mcp.InitializeRequest{}
		initReq.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
		initReq.Params.ClientInfo = mcp.Implementation{
			Name: "Forum Agent",
		}

		if _, err := cli.Initialize(ctx, initReq); err != nil {
			forumlog.Error("MCP 初始化失败", forumlog.String(err.Error()))
			tavilyErr = err
			return
		}

		tavilyTools, tavilyErr = eino_mcp.GetTools(ctx, &eino_mcp.Config{Cli: cli})
	})

	return tavilyTools, tavilyErr
}

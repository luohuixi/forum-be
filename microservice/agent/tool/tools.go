package tool

import (
	"forum-agent/tool/mcp"

	einotool "github.com/cloudwego/eino/components/tool"
)

func ToolList() []einotool.BaseTool {
	tools, err := mcp.TavilyTools()
	if err != nil {
		tools = nil
	}
	tools = append(tools, ToolListRAG()...)
	tools = append(tools, ToolListKafka()...)
	return tools
}

package core

import (
	"context"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func CommentPrompt(content string) ([]*schema.Message, error) {
	input := prompt.FromMessages(
		schema.GoTemplate,
		schema.SystemMessage("你是一个论坛回答助手，你的任务是调用工具查找资料并结合自己的思考回答用户的问题，然后将回答发送到Kafka消息队列(有提供相应工具)"),
		schema.SystemMessage("优先查询es向量数据库，其余搜索工具作为兜底和完善，返回的回答如果参考了es向量数据库，要将参考的元信息如: post_id, title 也要返回，且查询es向量数据库时如果用户的疑问过于复杂，应拆分成多个关键词搜索相关信息，而不是直接将用户问题作为查询输入"),
		schema.UserMessage(content),
	)

	return input.Format(context.Background(), nil)
}

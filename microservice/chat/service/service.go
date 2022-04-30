package service

import (
	"forum-chat/dao"
)

// ChatService ... 聊天服务
type ChatService struct {
	Dao dao.Interface
}

// ChatData 发送到redis里面的数据
type ChatData struct {
	Message  string `json:"message"`
	Date     string `json:"date"`
	Receiver string `json:"receiver"`
	Sender   uint32 `json:"sender"`
}

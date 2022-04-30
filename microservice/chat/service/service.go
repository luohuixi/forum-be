package service

import (
	"crypto/sha256"
	"encoding/hex"
	"forum-chat/dao"
	"strconv"
)

// ChatService ... 聊天服务
type ChatService struct {
	Dao dao.Interface
}

func sum256(data uint32) string {
	hash := sha256.Sum256([]byte(strconv.Itoa(int(data))))
	return hex.EncodeToString(hash[:])
}

// ChatData 发送到redis里面的数据
type ChatData struct {
	Message  string `json:"message"`
	Date     string `json:"date"`
	Receiver uint32 `json:"receiver"`
	Sender   uint32 `json:"sender"`
}

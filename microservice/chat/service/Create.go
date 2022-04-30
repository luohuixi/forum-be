package service

import (
	"context"
	"encoding/json"
	errno "forum-chat/errno"
	pb "forum-chat/proto"
	e "forum/pkg/err"
	"strconv"
	"time"
)

// Create 发送消息
func (s *ChatService) Create(ctx context.Context, req *pb.PushRequest, resp *pb.Response) error {

	queue := sum256(req.UserId)
	data := ChatData{
		Message:  req.Message,
		Date:     strconv.FormatInt(time.Now().UnixNano(), 10),
		Receiver: req.TargetUserId,
		Sender:   req.UserId,
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return e.ServerErr(errno.ErrDatabase, err.Error())
	}

	// logger.With("queue", queue).With("data", req.Message).With("userID", req.TokenId).Debugln("用户发送消息")
	err = s.Dao.Post(queue, dataBytes)

	if err != nil {
		return e.ServerErr(errno.ErrDatabase, err.Error())
	}

	return nil
}

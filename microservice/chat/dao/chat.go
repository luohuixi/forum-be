package dao

import (
	"encoding/json"
	pb "forum-chat/proto"
	"forum/log"
	"strconv"
	"time"
)

func (d *Dao) Create(data *ChatData) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return d.Redis.LPush("chat:"+strconv.Itoa(int(data.Receiver)), msg).Err()
}

func (d *Dao) GetList(id uint32, expiration time.Duration) ([]string, error) {
	key := "chat:" + strconv.Itoa(int(id))

	if d.Redis.LLen(key).Val() == 0 {
		msg, err := d.Redis.BRPop(expiration, key).Result() // 阻塞
		if err != nil {
			return nil, err
		}

		return msg[1:], nil // first ele is key
	}

	var list []string
	for d.Redis.LLen(key).Val() != 0 {
		msg, err := d.Redis.RPop(key).Result()
		if err != nil {
			return nil, err
		}

		list = append(list, msg)
	}

	return list, nil
}

// Rewrite 未成功发送的消息逆序放回list的Right
func (d *Dao) Rewrite(id uint32, list []string) error {
	log.Info("Rewrite")

	for i := len(list); i > 0; i-- {
		if err := d.Redis.RPush("chat:"+strconv.Itoa(int(id)), list[i-1]).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (d *Dao) ListHistory(userId, otherUserId, offset, limit uint32, pagination bool) ([]*pb.Message, error) {

	if otherUserId < userId {
		otherUserId, userId = userId, otherUserId
	}
	key := "history:" + strconv.Itoa(int(userId)) + "-" + strconv.Itoa(int(otherUserId)) // history:min_id-max_id

	var start int64 = 0
	var end int64 = -1

	if pagination {
		start = int64(offset)
		end = int64(offset + limit)
	}

	list, err := d.Redis.LRange(key, start, end).Result() // DESC
	if err != nil {
		return nil, err
	}

	histories := make([]*pb.Message, len(list))
	for i, history := range list {
		var msg pb.Message
		if err := json.Unmarshal([]byte(history), &msg); err != nil {
			return nil, err
		}
		histories[i] = &msg
	}

	return histories, nil
}

func (d *Dao) CreateHistory(userId uint32, list []string) error {
	log.Info("CreateHistory")

	for i := len(list); i > 0; i-- {
		var msg ChatData
		if err := json.Unmarshal([]byte(list[i-1]), &msg); err != nil {
			return err
		}

		min := userId
		if min > msg.Sender {
			min, msg.Sender = msg.Sender, min
		}

		if err := d.Redis.LPush("history:"+strconv.Itoa(int(min))+"-"+strconv.Itoa(int(msg.Sender)), list[i-1]).Err(); err != nil {
			return err
		}
	}

	return nil
}

package dao

import (
	"encoding/json"
	"time"
)

func (d *Dao) Create(data *ChatData) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return d.Redis.LPush(data.Receiver, msg).Err()
}

func (d *Dao) GetList(id string, expiration time.Duration) ([]string, error) {
	if d.Redis.LLen(id).Val() == 0 {
		msg, err := d.Redis.BRPop(expiration, id).Result()
		if err != nil {
			return nil, err
		}
		return msg, nil
	}

	var list []string
	for d.Redis.LLen(id).Val() != 0 {
		msg, err := d.Redis.RPop(id).Result()
		if err != nil {
			return nil, err
		}
		list = append(list, msg)
	}

	return list, nil
}

// Rewrite 未成功发送的消息逆序放回list的Right
func (d *Dao) Rewrite(id string, list []string) error {
	for i := len(list); i > 0; i-- {
		if err := d.Redis.RPush(id, list[i-1]).Err(); err != nil {
			return err
		}
	}
	return nil
}

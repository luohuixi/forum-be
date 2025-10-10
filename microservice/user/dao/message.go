package dao

import (
	"encoding/json"
	"strconv"
)

const (
	RedisKeyPrefix = "TeaHouse"
	Message        = "messages"
)

func GetKey(userId uint32) string {
	if userId == 0 {
		return RedisKeyPrefix + ":" + Message
	}
	return RedisKeyPrefix + ":" + Message + ":" + strconv.Itoa(int(userId))
}

func (d Dao) ListMessage() ([]string, error) {
	return d.Redis.LRange(GetKey(0), 0, -1).Result()
}

func (d Dao) ListPrivateMessage(userId uint32) ([]string, error) {
	return d.Redis.LRange(GetKey(userId), 0, -1).Result()
}

func (d Dao) CreateMessage(userId uint32, message string) error {
	return d.Redis.LPush(GetKey(userId), message).Err()
}

func (d Dao) DeleteOneMessage(userId uint32, uid string) error {
	message, err := d.ListPrivateMessage(userId)
	if err != nil {
		return err
	}

	var targetMessage string
	for _, msg := range message {
		var data map[string]string
		if err := json.Unmarshal([]byte(msg), &data); err == nil && data["id"] == uid {
			targetMessage = msg
			break
		}
	}

	if targetMessage != "" {
		_, err = d.Redis.LRem(GetKey(userId), 0, targetMessage).Result()
		return err
	}

	return nil
}

func (d Dao) DeleteMessage(userId uint32) error {
	return d.Redis.Del(GetKey(userId)).Err()
}

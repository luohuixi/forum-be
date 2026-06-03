package dao

import (
	"encoding/json"
	"strconv"
	"time"
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
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(message), &data); err == nil {
		if _, ok := data["read"]; !ok {
			data["read"] = false
		}
		if _, ok := data["created_at"]; !ok {
			data["created_at"] = time.Now().Format("2006-01-02 15:04:05")
		}
		if patched, err := json.Marshal(data); err == nil {
			message = string(patched)
		}
	}
	return d.Redis.LPush(GetKey(userId), message).Err()
}

func (d Dao) MarkOneMessageRead(userId uint32, uid string) error {
	messages, err := d.ListPrivateMessage(userId)
	if err != nil {
		return err
	}

	found := false
	for i, msg := range messages {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(msg), &data); err == nil && data["id"] == uid {
			data["read"] = true
			patched, err := json.Marshal(data)
			if err != nil {
				return err
			}
			messages[i] = string(patched)
			found = true
			break
		}
	}
	if !found {
		return nil
	}

	key := GetKey(userId)
	pipe := d.Redis.TxPipeline()
	pipe.Del(key)
	for i := len(messages) - 1; i >= 0; i-- {
		pipe.LPush(key, messages[i])
	}
	_, err = pipe.Exec()
	return err
}

func (d Dao) MarkAllMessageRead(userId uint32) error {
	messages, err := d.ListPrivateMessage(userId)
	if err != nil {
		return err
	}
	if len(messages) == 0 {
		return nil
	}

	key := GetKey(userId)
	pipe := d.Redis.TxPipeline()
	pipe.Del(key)
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(msg), &data); err == nil {
			data["read"] = true
			if patched, err := json.Marshal(data); err == nil {
				msg = string(patched)
			}
		}
		pipe.LPush(key, msg)
	}
	_, err = pipe.Exec()
	return err
}

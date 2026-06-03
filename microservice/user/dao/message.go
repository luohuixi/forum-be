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
	message = normalizeMessage(message, time.Now())
	return d.Redis.LPush(GetKey(userId), message).Err()
}

func normalizeMessage(message string, now time.Time) string {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(message), &data); err == nil {
		if _, ok := data["read"]; !ok {
			data["read"] = false
		}
		if _, ok := data["created_at"]; !ok {
			data["created_at"] = now.Format("2006-01-02 15:04:05")
		}
		if patched, err := json.Marshal(data); err == nil {
			return string(patched)
		}
	}
	return message
}

func (d Dao) CreateOrUpdateInteractionMessage(userId uint32, message string) error {
	now := time.Now()
	message = normalizeMessage(message, now)

	var next map[string]interface{}
	if err := json.Unmarshal([]byte(message), &next); err != nil {
		return d.CreateMessage(userId, message)
	}
	if !isDeduplicatedInteraction(next) {
		return d.CreateMessage(userId, message)
	}

	messages, err := d.ListPrivateMessage(userId)
	if err != nil {
		return err
	}

	merged := false
	next["read"] = false
	next["created_at"] = now.Format("2006-01-02 15:04:05")
	for i, msg := range messages {
		var current map[string]interface{}
		if err := json.Unmarshal([]byte(msg), &current); err != nil {
			continue
		}
		if sameInteractionMessage(current, next) {
			if id, ok := current["id"]; ok {
				next["id"] = id
			}
			patched, err := json.Marshal(next)
			if err != nil {
				return err
			}
			messages = append(messages[:i], messages[i+1:]...)
			messages = append([]string{string(patched)}, messages...)
			merged = true
			break
		}
	}

	if !merged {
		return d.Redis.LPush(GetKey(userId), message).Err()
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

func isDeduplicatedInteraction(data map[string]interface{}) bool {
	messageType := interfaceString(data["type"])
	return messageType == "like" || messageType == "collection"
}

func sameInteractionMessage(current, next map[string]interface{}) bool {
	return interfaceString(current["type"]) == interfaceString(next["type"]) &&
		interfaceString(current["send_user_id"]) == interfaceString(next["send_user_id"]) &&
		interfaceString(current["post_id"]) == interfaceString(next["post_id"]) &&
		interfaceString(current["comment_id"]) == interfaceString(next["comment_id"])
}

func interfaceString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatUint(uint64(v), 10)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	default:
		return ""
	}
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

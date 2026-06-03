package dao

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-redis/redis"
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

func (d Dao) ListDeduplicatedPrivateMessage(userId uint32) ([]string, error) {
	messages, err := d.ListPrivateMessage(userId)
	if err != nil {
		return nil, err
	}
	return deduplicateInteractionMessages(messages), nil
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

	next["read"] = false
	next["created_at"] = now.Format("2006-01-02 15:04:05")
	key := GetKey(userId)

	for attempt := 0; attempt < 5; attempt++ {
		err := d.Redis.Watch(func(tx *redis.Tx) error {
			messages, err := tx.LRange(key, 0, -1).Result()
			if err != nil {
				return err
			}
			messages, err = mergeInteractionIntoMessages(messages, next)
			if err != nil {
				return err
			}
			return replaceMessages(tx, key, messages)
		}, key)
		if err == redis.TxFailedErr {
			continue
		}
		return err
	}

	return redis.TxFailedErr
}

func deduplicateInteractionMessages(messages []string) []string {
	result := make([]string, 0, len(messages))
	seen := make(map[string]int)

	for _, msg := range messages {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(msg), &data); err != nil || !isDeduplicatedInteraction(data) {
			result = append(result, msg)
			continue
		}

		key := interactionKey(data)
		existingIndex, exists := seen[key]
		if !exists {
			seen[key] = len(result)
			result = append(result, msg)
			continue
		}

		var existing map[string]interface{}
		if err := json.Unmarshal([]byte(result[existingIndex]), &existing); err != nil {
			continue
		}
		result[existingIndex] = mergeInteractionMessage(existing, data)
	}

	return result
}

func mergeInteractionIntoMessages(messages []string, next map[string]interface{}) ([]string, error) {
	merged := copyMessageMap(next)
	rest := make([]string, 0, len(messages))
	preservedID := false

	for _, msg := range messages {
		var current map[string]interface{}
		if err := json.Unmarshal([]byte(msg), &current); err != nil {
			rest = append(rest, msg)
			continue
		}
		if sameInteractionMessage(current, merged) {
			if !preservedID && interfaceString(current["id"]) != "" {
				merged["id"] = current["id"]
				preservedID = true
			}
			continue
		}
		rest = append(rest, msg)
	}

	patched, err := json.Marshal(merged)
	if err != nil {
		return nil, err
	}
	return append([]string{string(patched)}, rest...), nil
}

func replaceMessages(tx *redis.Tx, key string, messages []string) error {
	_, err := tx.TxPipelined(func(pipe redis.Pipeliner) error {
		pipe.Del(key)
		for i := len(messages) - 1; i >= 0; i-- {
			pipe.LPush(key, messages[i])
		}
		return nil
	})
	return err
}

func copyMessageMap(data map[string]interface{}) map[string]interface{} {
	next := make(map[string]interface{}, len(data))
	for key, value := range data {
		next[key] = value
	}
	return next
}

func isDeduplicatedInteraction(data map[string]interface{}) bool {
	messageType := interfaceString(data["type"])
	return messageType == "like" || messageType == "collection"
}

func interactionKey(data map[string]interface{}) string {
	return interfaceString(data["type"]) + ":" +
		interfaceString(data["send_user_id"]) + ":" +
		interfaceString(data["post_id"]) + ":" +
		interfaceString(data["comment_id"])
}

func sameInteractionMessage(current, next map[string]interface{}) bool {
	return interactionKey(current) == interactionKey(next)
}

func mergeInteractionMessage(first, second map[string]interface{}) string {
	winner := first
	loser := second
	if messageTime(second).After(messageTime(first)) {
		winner = second
		loser = first
	}

	if interfaceString(winner["id"]) == "" {
		winner["id"] = loser["id"]
	}
	winner["read"] = boolValue(winner["read"]) && boolValue(loser["read"])

	patched, err := json.Marshal(winner)
	if err != nil {
		return mustMarshalMessage(first)
	}
	return string(patched)
}

func mustMarshalMessage(data map[string]interface{}) string {
	patched, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(patched)
}

func messageTime(data map[string]interface{}) time.Time {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", interfaceString(data["created_at"]), time.Local)
	return t
}

func boolValue(value interface{}) bool {
	v, ok := value.(bool)
	return ok && v
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

	var target map[string]interface{}
	for i, msg := range messages {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(msg), &data); err == nil && data["id"] == uid {
			target = data
			data["read"] = true
			patched, err := json.Marshal(data)
			if err != nil {
				return err
			}
			messages[i] = string(patched)
			break
		}
	}
	if target == nil {
		return nil
	}

	if isDeduplicatedInteraction(target) {
		for i, msg := range messages {
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(msg), &data); err != nil {
				continue
			}
			if sameInteractionMessage(data, target) {
				data["read"] = true
				patched, err := json.Marshal(data)
				if err != nil {
					return err
				}
				messages[i] = string(patched)
			}
		}
	}
	messages = deduplicateInteractionMessages(messages)

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

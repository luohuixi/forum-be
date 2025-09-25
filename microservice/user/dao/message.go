package dao

import "strconv"

var t = []string{"like", "comment", "collection", "reply_comment"}

const (
	RedisKeyPrefix = "TeaHouse"
	Message        = "messages"
)

func GetKey(userId uint32, t string) string {
	if userId == 0 {
		return RedisKeyPrefix + ":" + Message
	}
	return RedisKeyPrefix + ":" + Message + ":" + strconv.Itoa(int(userId)) + ":" + t
}

func (d Dao) ListMessage() ([]string, error) {
	publicMessage, err := d.Redis.LRange(GetKey(0, ""), 0, -1).Result()
	if err != nil {
		return nil, err
	}

	return publicMessage, nil
}

func (d Dao) ListPrivateMessage(userId uint32) ([]string, error) {
	var messages []string

	for _, str := range t {
		userMessages, err := d.Redis.LRange(GetKey(userId, str), 0, -1).Result()
		if err != nil {
			return nil, err
		}
		messages = append(messages, userMessages...)
	}

	return messages, nil
}

func (d Dao) CreateMessage(userId uint32, t, message string) error {
	return d.Redis.LPush(GetKey(userId, t), message).Err()
}

func (d Dao) DeleteMessage(userId uint32) error {
	for _, str := range t {
		if err := d.Redis.Del(GetKey(userId, str)).Err(); err != nil {
			return err
		}
	}
	return nil
}

package dao

import "strconv"

var t = []string{"like", "comment", "collection", "reply_comment"}

func (d Dao) ListMessage() ([]string, error) {
	publicMessage, err := d.Redis.LRange("messages:", 0, -1).Result()
	if err != nil {
		return nil, err
	}

	return publicMessage, nil
}

func (d Dao) ListPrivateMessage(userId uint32) ([]string, error) {
	var messages []string

	for _, str := range t {
		userMessages, err := d.Redis.LRange("messages:"+strconv.Itoa(int(userId))+":"+str, 0, -1).Result()
		if err != nil {
			return nil, err
		}
		messages = append(messages, userMessages...)
	}

	return messages, nil
}

func (d Dao) CreateMessage(userId uint32, t, message string) error {
	key := "messages:"
	if userId != 0 {
		key += strconv.Itoa(int(userId))
	}
	if t != "" {
		key += ":" + t
	}

	return d.Redis.LPush(key, message).Err()
}

func (d Dao) DeleteMessage(userId uint32) error {
	key := "messages:"
	if userId != 0 {
		key += strconv.Itoa(int(userId))
	}

	for _, str := range t {
		if err := d.Redis.Del(key + ":" + str).Err(); err != nil {
			return err
		}
	}
	return nil
}

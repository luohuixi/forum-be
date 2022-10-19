package dao

import "strconv"

func (d Dao) ListMessage(userId uint32) ([]string, error) {
	var messages []string
	publicMessage, err := d.Redis.LRange("messages:", 0, -1).Result()
	if err != nil {
		return nil, err
	}

	messages = append(messages, publicMessage...)

	userMessage, err := d.Redis.LRange("messages:"+strconv.Itoa(int(userId)), 0, -1).Result()
	if err != nil {
		return nil, err
	}

	messages = append(messages, userMessage...)

	return messages, nil
}

func (d Dao) CreateMessage(userId uint32, message string) error {
	key := "messages:"
	if userId != 0 {
		key += strconv.Itoa(int(userId))
	}

	return d.Redis.LPush(key, message).Err()
}

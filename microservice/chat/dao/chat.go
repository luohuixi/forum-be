package dao

import (
	"encoding/json"
	"fmt"
)

func (d *Dao) Create(data *ChatData) error {
	fmt.Println("dao.Create")
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return d.Redis.LPush(data.Receiver, msg).Err()
}

func (d *Dao) GetList(id string) ([]string, error) {
	var list []string

	for message, err := d.Redis.RPop(id).Result(); message != ""; {
		if err != nil {
			return nil, err
		}
		list = append(list, message)
	}

	return list, nil
	// return d.Redis.LRange(strconv.Itoa(int(id)), -1, 0).Result()
}

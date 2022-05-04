package dao

import (
	"encoding/json"
	"log"
	"time"
)

func (d *Dao) Create(data *ChatData) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return d.Redis.LPush(data.Receiver, msg).Err()
}

func (d *Dao) GetList(id string) ([]string, error) {
	list, err := d.Redis.BRPop(time.Hour, id).Result()
	if err != nil {
		return nil, err
	}
	log.Println("return", list)
	return list, nil
}

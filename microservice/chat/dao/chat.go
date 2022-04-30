package dao

import (
	"strconv"
)

func (d *Dao) Create(data *ChatData) error {
	return d.Redis.LPush(data.Receiver, strconv.Itoa(int(data.Sender))+"/"+data.Date+"/"+data.Message).Err()
}

func (d *Dao) GetList(id uint32) ([]string, error) {
	var list []string

	for message, err := d.Redis.RPop(strconv.Itoa(int(id))).Result(); message != ""; {
		if err != nil {
			return nil, err
		}
		list = append(list, message)
	}

	return list, nil
}

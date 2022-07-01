package dao

import (
	"github.com/go-redis/redis"
	"strconv"
)

func (d *Dao) AddLike(userId uint32, target Item) error {
	id := strconv.Itoa(int(target.Id))

	pipe := d.Redis.TxPipeline()

	pipe.SAdd("like:"+d.Map[target.TypeId]+"_list:"+id, userId)

	pipe.SAdd("like:user:"+strconv.Itoa(int(userId)), strconv.Itoa(int(target.TypeId))+":"+id)

	_, err := pipe.Exec()
	return err
}

func (d *Dao) RemoveLike(userId uint32, target Item) error {
	id := strconv.Itoa(int(target.Id))

	pipe := d.Redis.TxPipeline()

	pipe.SRem("like:"+d.Map[target.TypeId]+"_list:"+id, userId)

	pipe.SRem("like:user:"+strconv.Itoa(int(userId)), strconv.Itoa(int(target.TypeId))+":"+id)

	_, err := pipe.Exec()
	return err
}

func (d *Dao) GetLikedNum(target Item) (int64, error) {
	return d.Redis.SCard("like:" + d.Map[target.TypeId] + "_list:" + strconv.Itoa(int(target.Id))).Result()
	// res, err := d.Redis.Get("like:" + d.Map[target.TypeId] + ":" + strconv.Itoa(int(target.Id))).Result()
	// if err == redis.Nil {
	// 	return 0, nil
	// } else if err != nil {
	// 	return 0, err
	// }
}

func (d *Dao) IsUserHadLike(userId uint32, target Item) (bool, error) {
	return d.Redis.SIsMember("like:"+d.Map[target.TypeId]+"_list:"+strconv.Itoa(int(target.Id)), userId).Result()
}

func (d *Dao) ListUserLike(userId uint32) ([]*Item, error) {
	res, err := d.Redis.SMembers("like:user:" + strconv.Itoa(int(userId))).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var list []*Item
	for _, s := range res { // eg:  1:666
		typeId, err := strconv.Atoi(string(s[0]))
		if err != nil {
			panic(err)
		}

		id, err := strconv.Atoi(s[2:])
		if err != nil {
			panic(err)
		}

		list = append(list, &Item{
			Id:     uint32(id),
			TypeId: uint8(typeId),
		})
	}
	return list, nil
}

type Item struct {
	Id     uint32
	TypeId uint8
}

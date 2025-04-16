package dao

import (
	"github.com/go-redis/redis"
	"strconv"
	"strings"
)

type Item struct {
	Id       uint32
	TypeName string
}

func (d *Dao) AddLike(userId uint32, target Item) error {
	id := strconv.Itoa(int(target.Id))

	pipe := d.Redis.TxPipeline()

	pipe.SAdd("like:"+target.TypeName+"_list:"+id, userId)

	pipe.SAdd("like:user:"+strconv.Itoa(int(userId)), target.TypeName+":"+id)

	_, err := pipe.Exec()
	return err
}

func (d *Dao) RemoveLike(userId uint32, target Item) error {
	id := strconv.Itoa(int(target.Id))

	pipe := d.Redis.TxPipeline()

	pipe.SRem("like:"+target.TypeName+"_list:"+id, userId)

	pipe.SRem("like:user:"+strconv.Itoa(int(userId)), target.TypeName+":"+id)

	_, err := pipe.Exec()
	return err
}

func (d *Dao) GetLikedNum(target Item) (int64, error) {
	return d.Redis.SCard("like:" + target.TypeName + "_list:" + strconv.Itoa(int(target.Id))).Result()
	// res, err := d.Redis.Get("like:" + target.TypeName + ":" + strconv.Itoa(int(target.Id))).Result()
	// if err == redis.Nil {
	// 	return 0, nil
	// } else if err != nil {
	// 	return 0, err
	// }
}

func (d *Dao) IsUserHadLike(userId uint32, target Item) (bool, error) {
	return d.Redis.SIsMember("like:"+target.TypeName+"_list:"+strconv.Itoa(int(target.Id)), userId).Result()
}

func (d *Dao) ListUserLike(userId uint32) ([]*Item, error) {
	res, err := d.Redis.SMembers("like:user:" + strconv.Itoa(int(userId))).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var list []*Item
	for _, s := range res { // eg:  comment:666
		res := strings.Split(s, ":")

		id, err := strconv.Atoi(res[1])
		if err != nil {
			panic(err)
		}

		list = append(list, &Item{
			Id:       uint32(id),
			TypeName: res[0],
		})
	}
	return list, nil
}

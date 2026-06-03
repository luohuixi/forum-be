package dao

import (
	"encoding/json"
	"errors"
	"fmt"
	pb "forum-chat/proto"
	"forum/log"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

const (
	RedisPrefixKey = "TeaHouse"
	Chat           = "chat"
)

func GetKey(id uint32) string {
	return RedisPrefixKey + ":" + Chat + ":" + strconv.Itoa(int(id))
}

func (d *Dao) ensureMessageReadColumn() {
	if d.DB == nil || d.DB.Migrator().HasColumn(&DBdata{}, "read") {
		return
	}
	if err := d.DB.Exec("ALTER TABLE messages ADD COLUMN `read` tinyint(1) NOT NULL DEFAULT 1").Error; err != nil && !strings.Contains(strings.ToLower(err.Error()), "duplicate column") {
		log.Error("ensure message read column failed", zap.Error(err))
	}
}

func (d *Dao) Create(data *ChatData) error {
	//序列化为string后推入数列左侧(越靠左的消息越新)
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return d.Redis.LPush(GetKey(data.Receiver), msg).Err()
}

func (d *Dao) CreateMessage(msg *ChatData) error {
	if msg.Sender == 0 {
		msg.Sender = msg.LegacySender
	}
	var count int64
	err := d.DB.Table("messages").
		Where("receiver_id = ? AND sender_id = ? AND content = ? AND type_name = ? AND time = ?", msg.Receiver, msg.Sender, msg.Content, msg.TypeName, msg.Time).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	data := DBdata{
		ReceiverID: msg.Receiver,
		SenderID:   msg.Sender,
		Content:    msg.Content,
		TypeName:   msg.TypeName,
		Time:       msg.Time,
		Read:       false,
	}

	return d.DB.Table("messages").Create(&data).Error
}

func (d *Dao) GetList(id uint32, expiration time.Duration, wait bool) ([]string, error) {
	t := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("GetList Done(Receiver:%d)", id), zap.Duration("duration", time.Since(t)))
	}()
	// 使用用户的id创建一个key
	key := GetKey(id)
	// 如果数列里面为空的话则阻塞等待
	if d.Redis.LLen(key).Val() == 0 {
		if !wait {
			return nil, nil
		}

		msg, err := d.Redis.BRPop(expiration, key).Result() // 阻塞
		if err != nil {
			if errors.Is(err, redis.Nil) {
				// 超时但没拿到值，是正常现象
				return nil, nil
			}
			return nil, err
		}

		return msg[1:], nil // first ele is key
	}

	// 逐个读取整个队列,从右侧取(先取旧的消息)
	var list []string
	for d.Redis.LLen(key).Val() != 0 {
		msg, err := d.Redis.RPop(key).Result()
		if err != nil {
			return nil, err
		}

		list = append(list, msg)
	}

	return list, nil
}

// Rewrite 未成功发送的消息逆序放回list的Right
func (d *Dao) Rewrite(id uint32, list []string) error {
	log.Info("Rewrite")

	//从右侧将消息回写,这个地方为了保持消息的顺序与原来保持一致需要倒序写入
	for i := len(list); i > 0; i-- {
		if err := d.Redis.RPush(GetKey(id), list[i-1]).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (d *Dao) SyncPendingHistory(userId uint32) error {
	list, err := d.Redis.LRange(GetKey(userId), 0, -1).Result()
	if err != nil {
		return err
	}
	return d.CreateHistory(userId, list)
}

//func (d *Dao) ListHistory(userId, otherUserId, offset, limit uint32, pagination bool) ([]*pb.Message, error) {
//	//这个地方为了保证顺序要进行交换
//	if otherUserId < userId {
//		otherUserId, userId = userId, otherUserId
//	}
//
//	//读取这两个人的历史消息
//	key := "history:" + strconv.Itoa(int(userId)) + "-" + strconv.Itoa(int(otherUserId)) // history:min_id-max_id
//
//	var start int64 = 0
//	var end int64 = -1
//
//	if pagination {
//		start = int64(offset)
//		end = int64(offset + limit)
//	}
//
//	//批量读取消息,这个地方要注意索引从小到大是从新到旧
//	list, err := d.Redis.LRange(key, start, end).Result() // DESC
//	if err != nil {
//		return nil, err
//	}
//
//	histories := make([]*pb.Message, len(list))
//	for i, history := range list {
//		var msg pb.Message
//		if err := json.Unmarshal([]byte(history), &msg); err != nil {
//			return nil, err
//		}
//		histories[i] = &msg
//	}
//
//	return histories, nil
//}

func (d *Dao) ListHistory(userId, otherUserId, offset, limit uint32, pagination bool) ([]*pb.Message, error) {
	var history []*pb.Message

	if pagination {
		err := d.DB.Table("messages").Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", userId, otherUserId, otherUserId, userId).Order("time DESC").Offset(int(offset)).Limit(int(limit)).Find(&history).Error
		if err != nil {
			return nil, err
		}
	} else {
		if err := d.DB.Table("messages").Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", userId, otherUserId, otherUserId, userId).Order("time DESC").Find(&history).Error; err != nil {
			return nil, err
		}
	}

	return history, nil
}

//func (d *Dao) CreateHistory(userId uint32, list []string) error {
//	log.Info("CreateHistory")
//
//	for i := len(list); i > 0; i-- {
//		var msg ChatData
//		if err := json.Unmarshal([]byte(list[i-1]), &msg); err != nil {
//			return err
//		}
//
//		//调整为统一顺序,为什么要这么做呢?因为这么做首先可以保证一致性,其次可以减半存储空间
//		minId := userId
//		if minId > msg.Sender {
//			minId, msg.Sender = msg.Sender, minId
//		}
//
//		//推送到这两个人的历史消息里面去
//		if err := d.Redis.LPush("history:"+strconv.Itoa(int(minId))+"-"+strconv.Itoa(int(msg.Sender)), list[i-1]).Err(); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}

func (d *Dao) CreateHistory(userId uint32, list []string) error {
	for i := len(list); i > 0; i-- {
		var msg ChatData
		if err := json.Unmarshal([]byte(list[i-1]), &msg); err != nil {
			return err
		}
		if msg.Sender == 0 {
			msg.Sender = msg.LegacySender
		}

		if msg.Receiver == 0 {
			msg.Receiver = userId
		}
		if err := d.CreateMessage(&msg); err != nil {
			return err
		}
	}
	return nil
}

type ConversationSummary struct {
	OtherID         uint32 `gorm:"column:other_id"`
	LastMessageTime string `gorm:"column:last_message_time"`
	LastMessage     string `gorm:"column:last_message"`
	UnreadCount     uint32 `gorm:"column:unread_count"`
}

func (d *Dao) GetUserList(userId uint32, limit, page int) ([]*pb.UserStatus, error) {
	var summaries []ConversationSummary

	// 未读会话优先，再按最近一条消息排序。
	err := d.DB.Table("messages").
		Select(`
			CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END AS other_id,
			MAX(time) AS last_message_time,
			SUBSTRING_INDEX(GROUP_CONCAT(content ORDER BY time DESC SEPARATOR '\n'), '\n', 1) AS last_message,
			SUM(CASE WHEN receiver_id = ? AND `+"`read`"+` = 0 THEN 1 ELSE 0 END) AS unread_count
		`, userId, userId).
		Where("sender_id = ? OR receiver_id = ?", userId, userId).
		Group("other_id").
		Order("unread_count DESC, MAX(time) DESC").
		Offset(page * limit).
		Limit(limit).
		Scan(&summaries).Error

	if err != nil {
		return nil, err
	}

	usersStatus := make([]*pb.UserStatus, 0, len(summaries))
	for _, summary := range summaries {
		var userStatus pb.UserStatus
		err := d.DB.Table("users").Where("id = ?", summary.OtherID).Find(&userStatus).Error
		if err != nil {
			return nil, err
		}
		userStatus.LastMessageTime = &summary.LastMessageTime
		userStatus.LastMessage = &summary.LastMessage
		userStatus.UnreadCount = summary.UnreadCount

		usersStatus = append(usersStatus, &userStatus)
	}

	return usersStatus, nil
}

func (d *Dao) MarkRead(userId, otherUserId uint32) error {
	return d.DB.Table("messages").
		Where("receiver_id = ? AND sender_id = ? AND `read` = ?", userId, otherUserId, false).
		Update("read", true).Error
}

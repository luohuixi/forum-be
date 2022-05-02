package chat

import (
	"context"
	"fmt"
	pb "forum-chat/proto"
	. "forum-gateway/handler"
	"forum-gateway/log"
	"forum-gateway/pkg/errno"
	"forum-gateway/service"
	"forum-gateway/util"
	m "forum/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	l "log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	UserId string
	Socket *websocket.Conn
	Open   bool
}

// WsHandler socket 连接 中间件 作用:升级协议,用户验证,自定义信息等
func WsHandler(c *gin.Context) {
	log.Info("Chat WsHandler function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var upGrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		SendError(c, errno.ErrWebsocket, nil, err.Error(), GetLine())
		return
	}

	id := c.DefaultQuery("id", "20")
	userId, ok, err := m.GetStringFromRedis(id)
	if !ok {
		l.Println("not ok")
		conn.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
		return
	}

	client := &Client{
		UserId: userId,
		Socket: conn,
		Open:   true,
	}

	go client.Read()
	go client.Write()
}

// Read 从client接收消息
func (c *Client) Read() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			l.Println("client close connect")
			c.Open = false
			break
		}

		index := strings.IndexByte(string(message), '-')
		if index == -1 || index == 0 {
			c.Socket.WriteMessage(websocket.CloseMessage, []byte("format error, eg. 5-外比巴卜"))
			break
		}
		targetUserId := string(message)[:index]
		if _, err := strconv.Atoi(targetUserId); err != nil {
			c.Socket.WriteMessage(websocket.CloseMessage, []byte("format error, eg. 5-外比巴卜"))
			break
		}

		if targetUserId == c.UserId {
			c.Socket.WriteMessage(websocket.CloseMessage, []byte("error: can't message yourself"))
			break
		}

		createReq := &pb.CreateRequest{
			UserId:       c.UserId,
			TargetUserId: targetUserId,
			Message:      string(message)[index+1:],
		}

		if _, err := service.ChatClient.Create(context.Background(), createReq); err != nil {
			fmt.Println(err)
			c.Socket.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
			break
		}
	}
}

// Write 返回client收到的消息
func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		if !c.Open {
			break
		}

		res, err := m.RedisDB.Self.BRPop(time.Hour, c.UserId).Result()
		if err != nil {
			log.Error(err.Error())
			c.Socket.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
			break
		}

		for _, msg := range res {
			l.Println(msg)
			if err != nil {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
				break
			}
			c.Socket.WriteMessage(websocket.TextMessage, []byte(msg))
		}
	}
}

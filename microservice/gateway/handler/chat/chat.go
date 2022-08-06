package chat

import (
	"context"
	"encoding/json"
	"fmt"
	pb "forum-chat/proto"
	. "forum-gateway/handler"
	"forum-gateway/service"
	"forum-gateway/util"
	"forum/log"
	m "forum/model"
	"forum/pkg/errno"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Client struct {
	UserId string
	Socket *websocket.Conn
	Close  chan struct{}
}

// WsHandler ... socket 连接 中间件 作用:升级协议,用户验证,自定义信息等
// @Summary WebSocket
// @Description 建立 WebSocket 连接
// @Tags chat
// @Param id query string true "uuid"
// @Router /chat/ws [get]
func WsHandler(c *gin.Context) {
	log.Info("Chat WsHandler function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var upGrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		SendError(c, errno.ErrWebsocket, nil, err.Error(), GetLine())
		return
	}

	id := c.DefaultQuery("id", "0")
	userId, ok, err := m.GetStringFromRedis(id)
	if !ok {
		log.Error("not ok")
		conn.WriteMessage(websocket.CloseMessage, []byte("error"))
		return
	}

	client := &Client{
		UserId: userId,
		Socket: conn,
		Close:  make(chan struct{}),
	}

	go client.Read()
	go client.Write()
}

// Read 从client接收消息
func (c *Client) Read() {
	defer func() {
		close(c.Close)
		c.Socket.Close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			log.Info("client close connect")
			break
		}

		var req pb.CreateRequest
		if err := json.Unmarshal(message, &req); err != nil {
			log.Error(err.Error())
			c.Socket.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			break
		}

		req.UserId = c.UserId
		// index := strings.IndexByte(string(message), '-')
		// if index == -1 || index == 0 {
		// 	log.Error("index wrong")
		// 	c.Socket.WriteMessage(websocket.TextMessage, []byte("format error, eg. 5-外比巴卜"))
		// 	break
		// }
		// targetUserId := string(message)[:index]
		// if _, err := strconv.Atoi(targetUserId); err != nil {
		// 	log.Error(err.Error())
		// 	c.Socket.WriteMessage(websocket.TextMessage, []byte("format error, eg. 5-外比巴卜"))
		// 	break
		// }
		//
		if req.TargetUserId == c.UserId {
			log.Error("error: can't message yourself")
			c.Socket.WriteMessage(websocket.TextMessage, []byte("error: can't message yourself"))
			break
		}
		//
		// createReq := &pb.CreateRequest{
		// 	UserId:       c.UserId,
		// 	TargetUserId: targetUserId,
		// 	Content:      string(message)[index+1:],
		// }

		if _, err := service.ChatClient.Create(context.Background(), &req); err != nil {
			log.Error(err.Error())
			c.Socket.WriteMessage(websocket.TextMessage, []byte(err.Error()))
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
		getListRequest := &pb.GetListRequest{
			UserId: c.UserId,
		}
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour)) // set rpc expiration to 1 Hour
		go func() {
			<-c.Close // cancel the request when client close connect
			fmt.Println("cancel", c.UserId)
			cancel()
		}()

		resp, err := service.ChatClient.GetList(ctx, getListRequest)
		if err != nil {
			log.Error(err.Error())
			c.Socket.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return
		}

		for _, msg := range resp.List {
			fmt.Println("msg", msg)
			c.Socket.WriteMessage(websocket.TextMessage, []byte(msg))
		}
	}
}

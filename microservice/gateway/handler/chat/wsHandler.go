package chat

import (
	"context"
	pb "forum-chat/proto"
	. "forum-gateway/handler"
	"forum-gateway/log"
	"forum-gateway/pkg/errno"
	"forum-gateway/service"
	"forum-gateway/util"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	UserId uint32
	Socket *websocket.Conn
}

// WsHandler TestHandler socket 连接 中间件 作用:升级协议,用户验证,自定义信息等
func WsHandler(c *gin.Context) {
	log.Info("Chat create function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	id := c.MustGet("userId").(uint32)
	var req CreateRequest
	if err := c.ShouldBindJSON(req); err != nil {
		SendBadRequest(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	var upGrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		SendError(c, errno.ErrWebsocket, nil, err.Error(), GetLine())
		return
	}

	client := &Client{
		UserId: id,
		Socket: conn,
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
		c.Socket.PongHandler()
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			c.Socket.Close()
			break
		}

		res := strings.Split(string(message), "-")
		if _, err := strconv.Atoi(res[0]); err != nil {
			c.Socket.Close()
			break
		}

		createReq := &pb.CreateRequest{
			UserId:       c.UserId,
			TargetUserId: res[0],
			Message:      res[1],
		}
		if _, err := service.ChatClient.Create(context.Background(), createReq); err != nil {
			c.Socket.Close()
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
		getListReq := &pb.GetListRequest{
			UserId: c.UserId,
		}

		res, err := service.ChatClient.GetList(context.Background(), getListReq)
		if err != nil {
			c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
			c.Socket.Close()
			break
		}

		if len(res.List) == 0 {
			time.Sleep(time.Second)
			continue
		}

		for _, msg := range res.List {
			c.Socket.WriteMessage(websocket.TextMessage, []byte(msg))
		}
	}
}

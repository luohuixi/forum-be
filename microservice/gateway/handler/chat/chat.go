package chat

import (
	"context"
	"encoding/json"
	pb "forum-chat/proto"
	. "forum-gateway/handler"
	"forum-gateway/service"
	"forum-gateway/util"
	"forum/log"
	"forum/pkg/errno"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Client struct {
	UserId uint32
	Socket *websocket.Conn
	Close  chan struct{}
}

// WsHandler 建立 WebSocket 连接
// @Summary WebSocket 连接
// @Description 通过 WebSocket 进行实时通信。请使用 ws:// 或 wss:// 连接此接口，并在 Apifox 选择 WebSocket 请求。
// @Tags chat
// @Param id query string true "用户 UUID"
// @Success 101 {string} string "WebSocket 连接成功"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 500 {object} ErrorResponse "服务器错误"
// @Router /chat/ws [get]
func WsHandler(c *gin.Context) {
	log.Info("Chat WsHandler function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var upGrader = websocket.Upgrader{
		CheckOrigin:  func(r *http.Request) bool { return true },
		Subprotocols: []string{c.Request.Header.Get("Sec-WebSocket-Protocol")},
	}

	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		SendError(c, errno.ErrWebsocket, nil, err.Error(), GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)

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

		//读取写入的message
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

		if req.TargetUserId == c.UserId {
			log.Error("error: can't message yourself")
			c.Socket.WriteMessage(websocket.TextMessage, []byte("error: can't message yourself"))
			break
		}

		if req.TargetUserId == 0 {
			log.Error("error: wrong target_user_id")
			c.Socket.WriteMessage(websocket.TextMessage, []byte("error: wrong target_user_id"))
			break
		}
		//创建聊天记录
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
		//获取聊天记录
		getListRequest := &pb.GetListRequest{
			UserId: c.UserId,
		}

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour)) // set rpc expiration to 1 Hour
		go func() {
			<-c.Close // cancel the request when client close connect
			cancel()
		}()
		// 死循环获取,直到客户端断开连接
		resp, err := service.ChatClient.GetList(ctx, getListRequest)
		if err != nil {
			log.Error(err.Error())
			c.Socket.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return
		}

		for _, msg := range resp.List {
			c.Socket.WriteMessage(websocket.TextMessage, []byte(msg))
		}
	}
}

type Message struct {
	Content  string `json:"content"`
	Time     string `json:"time"`
	Sender   uint32 `json:"sender"`
	TypeName string `json:"type_name"`
}

package chat

import (
	"context"
	"encoding/json"
	pb "forum-chat/proto"
	. "forum-gateway/handler"
	"forum-gateway/util"
	"forum/log"
	"forum/pkg/errno"
	"net/http"
	"strings"
	"time"

	"forum/client"

	mclient "go-micro.dev/v4/client"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Client struct {
	UserId uint32
	Socket *websocket.Conn
	ctx    context.Context
	cancel context.CancelFunc
}

// WsHandler 建立 WebSocket 连接
// @Summary 建立 WebSocket 连接
// @Description 通过 WebSocket 实现客户端与服务器之间的实时通信。
// @Description 使用 `ws://` 或 `wss://` 协议访问此接口，连接成功后可进行双向通信。
// @Description 客户端连接后，请使用 JSON 格式发送消息，结构如下：
// @Description
// @Description ```json
// @Description {
// @Description   "target_user_id": 123,
// @Description   "content": "你好",
// @Description   "type_name": "text",
// @Description   "time": "2025-07-20 12:00:00"
// @Description }
// @Description ```
// @Tags chat
// @Param Sec-WebSocket-Protocol header string false "子协议，一般用于身份校验或版本协商"
// @Success 101 {string} string "WebSocket 连接成功"
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

	newCtx, cancel := context.WithCancel(context.WithoutCancel(c.Request.Context()))
	client := &Client{
		UserId: userId,
		Socket: conn,
		// 设置一个不受取消影响的上下文，WsHandler 结束时不会取消这个上下文，用于继承链路
		ctx:    newCtx,
		cancel: cancel,
	}

	go client.Read()
	go client.Write()
}

// Read 从client接收消息
func (c *Client) Read() {
	defer func() {
		c.cancel()
		c.Socket.Close()
	}()

	for {

		//读取写入的message
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			// 正常情况由 read 控制 write 的退出
			// 因为 ReadMessage 会一直阻塞直到断开连接报错或读取成功
			log.Info("client close connect", zap.Error(err))
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
		if _, err := client.ChatClient.Create(c.ctx, &req); err != nil {
			log.Error(err.Error())
			c.Socket.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			break
		}
	}
}

// Write 返回client收到的消息
func (c *Client) Write() {
	defer func() {
		c.cancel()
		c.Socket.Close()
	}()

	for {
		if c.ctx.Err() != nil {
			return
		}

		// 获取聊天记录
		getListRequest := &pb.GetListRequest{
			UserId: c.UserId,
			Wait:   true,
		}

		ctx, cancel := context.WithDeadline(c.ctx, time.Now().Add(time.Hour)) // set rpc expiration to 1 Hour
		// 死循环获取,直到客户端断开连接
		resp, err := client.ChatClient.GetList(ctx, getListRequest,
			mclient.WithRequestTimeout(time.Hour),
			withConnectionTimeout(time.Hour))

		cancel()

		if err != nil {
			if !strings.Contains(err.Error(), context.DeadlineExceeded.Error()) && !strings.Contains(err.Error(), context.Canceled.Error()) {
				log.Error("chatClient getList error", zap.Error(err))
			}
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

// go-microV4只设置RequestTimeout不能改变客户端超时，要设置ConnectionTimeout
func withConnectionTimeout(d time.Duration) mclient.CallOption {
	return func(o *mclient.CallOptions) {
		o.ConnectionTimeout = d
	}
}

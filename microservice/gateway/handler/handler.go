package handler

import (
	"bytes"
	"encoding/json"
	"forum-gateway/util"
	"forum/log"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"net/http"
	"runtime"
	"strconv"

	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Response 请求响应
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
} // @name Response

func GetLine() string {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return "forum-gateway/handler/handler.go:30"
	}
	return file + ":" + strconv.Itoa(line)
}

func SendMicroServiceResponse(c *gin.Context, err error, m proto.Message, data any) {
	jsonpbMarshaler := &jsonpb.Marshaler{
		EmitDefaults: true, // 是否将字段值为空的渲染到JSON结构中
		OrigName:     true, // 是否使用原生的proto协议中的字段
	}

	var buffer bytes.Buffer
	if err := jsonpbMarshaler.Marshal(&buffer, m); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(buffer.Bytes(), &data); err != nil {
		panic(err)
	}

	code, message := errno.DecodeErr(err)
	log.Info(message, zap.String("X-Request-PostId", util.GetReqID(c)))

	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func SendResponse(c *gin.Context, err error, data interface{}) {
	code, message := errno.DecodeErr(err)
	log.Info(message, zap.String("X-Request-PostId", util.GetReqID(c)))

	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func SendError(c *gin.Context, err error, data interface{}, cause string, source string) {
	code, message := errno.DecodeErr(err)
	log.Error(message,
		zap.String("X-Request-PostId", util.GetReqID(c)),
		zap.String("cause", cause),
		zap.String("source", source))

	var responseCode int
	switch {
	case code == http.StatusNotFound:
		responseCode = http.StatusNotFound
	case code > 20000:
		responseCode = http.StatusBadRequest
	default:
		responseCode = http.StatusInternalServerError
	}

	c.JSON(responseCode, Response{
		Code:    code,
		Message: message + " " + cause,
		Data:    data,
	})
}

package handler

import (
	"forum-gateway/util"
	"forum/log"
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
		return "forum-gateway/handler/handler.go:26"
	}
	return file + ":" + strconv.Itoa(line)
}

func SendResponse(c *gin.Context, err error, data interface{}) {
	code, message := errno.DecodeErr(err)
	log.Info(message, zap.String("X-Request-Id", util.GetReqID(c)))

	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func SendError(c *gin.Context, err error, data interface{}, cause string, source string) {
	code, message := errno.DecodeErr(err)
	log.Error(message,
		zap.String("X-Request-Id", util.GetReqID(c)),
		zap.String("cause", cause),
		zap.String("source", source))

	var responseCode = http.StatusInternalServerError

	if code == http.StatusNotFound {
		responseCode = http.StatusNotFound
	} else if code > 20000 {
		responseCode = http.StatusBadRequest
	}

	c.JSON(responseCode, Response{
		Code:    code,
		Message: message + " " + cause,
		Data:    data,
	})
}

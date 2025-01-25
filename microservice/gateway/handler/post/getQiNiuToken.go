package post

import (
	. "forum-gateway/handler"
	"forum-gateway/util"
	"forum/log"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// GetQiNiuToken ...
// @Summary  获取七牛云token
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Success 200 {object} QiNiuToken
// @Router /post/qiniu_token [get]
func (a *Api) GetQiNiuToken(c *gin.Context) {
	log.Info("Post GetQiNiuToken function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	token := getToken()
	resp := QiNiuToken{
		Token: token,
	}

	SendResponse(c, nil, resp)
}

var (
	accessKey, secretKey, bucketName, domainName, upToken string
)

func initOSS() {
	accessKey = viper.GetString("oss.access_key")
	secretKey = viper.GetString("oss.secret_key")
	bucketName = viper.GetString("oss.bucket_name")
	domainName = viper.GetString("oss.domain_name")
}

func getToken() string {
	var maxInt uint64 = 1 << 32
	initOSS()
	putPolicy := storage.PutPolicy{
		Scope:   bucketName,
		Expires: maxInt,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken = putPolicy.UploadToken(mac)
	return upToken
}

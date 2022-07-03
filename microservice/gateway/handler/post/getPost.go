package post

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/service"
	pb "forum-post/proto"
	"github.com/gin-gonic/gin"
)

func (a *Api) Test(c *gin.Context) {
	res, err := service.PostClient.GetPost(context.Background(), &pb.Request{
		Id: 2,
	})

	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}
	SendResponse(c, nil, res.Id)
}

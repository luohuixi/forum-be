package post

import (
	"context"
	"fmt"
	. "forum-gateway/handler"
	"forum-gateway/service"
	pb "forum-post/proto"
	"github.com/gin-gonic/gin"
)

func (a *Api) Test(c *gin.Context) {
	fmt.Println("123: ", 123)

	_, err := service.PostClient.GetPost(context.Background(), &pb.Request{
		Id: 1,
	})

	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}
}

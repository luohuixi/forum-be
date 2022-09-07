package main

import (
	"context"
	pb "forum-post/proto"
	"github.com/micro/go-micro"
)

func main() { // TODO
	service := micro.NewService(micro.Name("forum.cli.post"))

	service.Init()

	client := pb.NewPostServiceClient("forum.service.post", service.Client())

	_, err := client.CreatePost(context.TODO(), &pb.CreatePostRequest{
		UserId:  2,
		Content: "外比巴卜",
		// TypeId:     1, // 默认为1
		Title:    "first post",
		Category: "2",
	})
	// _, err = client.CreateComment(context.TODO(), &pb.CreateCommentRequest{
	// 	PostId: 1,
	// 	// TypeId:    2,
	// 	FatherId:  1,
	// 	Content:   "first comment to comment",
	// 	CreatorId: 2,
	// })
	// _, err = client.UpdatePostInfo(context.TODO(), &pb.UpdatePostInfoRequest{
	// 	Id:         1,
	// 	Content:    "",
	// 	Title:      "",
	// 	Category: 1,
	// })

	panic(err)
	//
	// fmt.Println("post:", post.List[0].Category)
	//
	// if err != nil {
	// 	fmt.Println("err: ", err)
	// }
}

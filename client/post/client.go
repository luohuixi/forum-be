package main

import (
	"context"
	"fmt"
	pb "forum-post/proto"
	micro "go-micro.dev/v4"
)

func main() { // TODO
	service := micro.NewService(micro.Name("forum.cli.post"))

	service.Init()

	client := pb.NewPostService("forum.service.post", service.Client())

	resp, err := client.ListMainPost(context.TODO(), &pb.ListMainPostRequest{
		UserId:     2,
		TypeName:   "normal",
		Limit:      10,
		Offset:     0,
		Pagination: true,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("----- : ", resp.Posts[0].IsCollection, " -----")
	fmt.Println("----- : ", len(resp.Posts), " -----")

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

	//
	// fmt.Println("post:", post.List[0].Category)
	//
	// if err != nil {
	// 	fmt.Println("err: ", err)
	// }
}

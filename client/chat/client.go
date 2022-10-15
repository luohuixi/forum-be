package main

import (
	"context"
	"fmt"
	pb "forum-chat/proto"
	micro "go-micro.dev/v4"
)

func main() {
	service := micro.NewService(micro.Name("forum.cli.chat"))

	service.Init()

	client := pb.NewChatService("forum.service.chat", service.Client())

	resp, err := client.ListHistory(context.TODO(), &pb.ListHistoryRequest{
		UserId:      5,
		OtherUserId: 4,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Histories)
}

package service

import (
	"context"

	pb "forum-user/proto"
	"forum/client"
)

func GetUserProfile(ctx context.Context, id uint32, viewerID uint32) (*pb.UserProfile, error) {
	return client.UserClient.GetProfile(ctx, &pb.GetRequest{Id: id, ViewerId: viewerID})
}

func IsCollectionAndLikePublic(ctx context.Context, userID uint32) (bool, error) {
	profile, err := GetUserProfile(ctx, userID, 0)
	if err != nil {
		return false, err
	}
	return profile.GetIsPublicCollectionAndLike(), nil
}

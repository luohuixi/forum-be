package user

import (
	"context"

	pb "forum-user/proto"
	"forum/client"
)

func IsCollectionAndLikePublic(ctx context.Context, userID uint32) (bool, error) {
	profile, err := client.UserClient.GetProfile(ctx, &pb.GetRequest{Id: userID})
	if err != nil {
		return false, err
	}
	return profile.GetIsPublicCollectionAndLike(), nil
}

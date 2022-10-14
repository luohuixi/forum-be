package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	"forum-user/pkg/role"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"forum/pkg/handler"

	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	opentracingWrapper "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"

	pbu "forum-user/proto"

	micro "go-micro.dev/v4"
)

// PostService ... 帖子服务
type PostService struct {
	Dao dao.Interface
}

func New(i dao.Interface) *PostService {
	service := new(PostService)
	service.Dao = i
	return service
}

var UserClient pbu.UserService

func UserInit() {
	service := micro.NewService(micro.Name("forum.cli.user"),
		micro.WrapClient(
			opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapCall(handler.ClientErrorHandlerWrapper()),
	)

	service.Init()

	UserClient = pbu.NewUserService("forum.service.user", service.Client())
}

func (s PostService) processComments(userId uint32, commentInfos []*dao.CommentInfo) []*pb.CommentInfo {
	comments := make([]*pb.CommentInfo, len(commentInfos))

	for i, comment := range commentInfos {
		item := dao.Item{
			Id:       comment.Id,
			TypeName: constvar.Comment,
		}

		num, err := s.Dao.GetLikedNum(item)
		if err != nil {
			logger.Error(errno.ErrRedis.Error(), logger.String(err.Error()))
		}

		isLiked, err := s.Dao.IsUserHadLike(userId, item)
		if err != nil {
			logger.Error(errno.ErrRedis.Error(), logger.String(err.Error()))
		}

		comments[i] = &pb.CommentInfo{
			Id:            comment.Id,
			TypeName:      comment.TypeName,
			Content:       comment.Content,
			FatherId:      comment.FatherId,
			Time:          comment.CreateTime,
			CreatorId:     comment.CreatorId,
			CreatorName:   comment.CreatorName,
			CreatorAvatar: comment.CreatorAvatar,
			LikeNum:       uint32(num),
			IsLiked:       isLiked,
			ImgUrls:       comment.ImgUrls,
		}
	}

	return comments
}

func (s PostService) getPostInfo(postId uint32, userId uint32) (bool, bool, uint32, []string, uint32, uint32) {
	item := dao.Item{
		Id:       postId,
		TypeName: constvar.Post,
	}

	isLiked, err := s.Dao.IsUserHadLike(userId, item)
	if err != nil {
		logger.Error(errno.ErrRedis.Error(), logger.String(err.Error()))
	}

	isCollection, err := s.Dao.IsUserCollectionPost(userId, postId)
	if err != nil {
		logger.Error(errno.ErrDatabase.Error(), logger.String(err.Error()))
	}

	likeNum, err := s.Dao.GetLikedNum(item)
	if err != nil {
		logger.Error(errno.ErrRedis.Error(), logger.String(err.Error()))
	}

	tags, err := s.Dao.ListTagsByPostId(postId)
	if err != nil {
		logger.Error(errno.ErrDatabase.Error(), logger.String(err.Error()))
	}

	commentNum, err := s.Dao.GetCommentNumByPostId(postId)
	if err != nil {
		logger.Error(errno.ErrDatabase.Error(), logger.String(err.Error()))
	}

	collectionNum, err := s.Dao.GetCollectionNumByPostId(postId)
	if err != nil {
		logger.Error(errno.ErrDatabase.Error(), logger.String(err.Error()))
	}

	return isLiked, isCollection, uint32(likeNum), tags, commentNum, collectionNum
}

func (s PostService) GetUserDomain(userId uint32) (string, error) {
	getResp, err := UserClient.GetProfile(context.TODO(), &pbu.GetRequest{Id: userId})
	if err != nil {
		return "", err
	}

	return role.Role2Domain(getResp.Role), nil
}

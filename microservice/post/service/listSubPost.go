package service

//
// import (
// 	"context"
// 	"forum-post/dao"
// 	pb "forum-post/proto"
// 	logger "forum/log"
// 	"forum/pkg/errno"
// 	"go.uber.org/zap"
// )
//
// func (s *PostService) ListSubPost(_ context.Context, req *pb.ListSubPostRequest, resp *pb.ListPostResponse) error {
// 	logger.Info("PostService ListSubPost")
//
// 	filter := &dao.PostModel{
// 		TypeName: req.TypeName,
// 		// MainPostId: req.MainPostId, FIXME
// 	}
//
// 	mainPost, err := s.Dao.GetPost(req.MainPostId)
// 	if err != nil {
// 		return errno.ServerErr(errno.ErrDatabase, err.Error())
// 	}
//
// 	var posts []*dao.PostModel
//
// 	subPosts, err := s.Dao.ListPost(filter, req.Offset, req.Limit, req.LastId, req.Pagination)
// 	if err != nil {
// 		return errno.ServerErr(errno.ErrDatabase, err.Error())
// 	}
//
// 	posts = append(posts, mainPost)
// 	posts = append(posts, subPosts...)
// 	resp.List = make([]*pb.Post, len(posts))
// 	for i, post := range posts {
// 		comments, err := s.Dao.ListCommentByPostId(post.Id)
// 		if err != nil {
// 			logger.Error(err.Error(), zap.Error(errno.ErrDatabase))
// 		}
//
// 		s.processComments(req.UserId, comments)
//
// 		isLiked, isCollection, likeNum, tags, commentNum := s.getPostInfo(post.Id, req.UserId)
//
// 		if likeNum != 0 {
// 			post.LikeNum = likeNum
// 		}
//
// 		resp.List[i] = &pb.Post{
// 			Id:            post.Id,
// 			Title:         post.Title,
// 			Time:          post.LastEditTime,
// 			Category:      post.Category,
// 			CreatorId:     post.CreatorId,
// 			CreatorName:   post.CreatorName,
// 			CreatorAvatar: post.CreatorAvatar,
// 			Content:       post.Content,
// 			LikeNum:       post.LikeNum,
// 			CommentNum:    commentNum,
// 			IsLiked:       isLiked,
// 			IsCollection:  isCollection,
// 			Comments:      comments,
// 			Tags:          tags,
// 		}
// 	}
//
// 	return nil
// }

package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"forum/pkg/pagetoken"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *PostService) ListComments(_ context.Context, req *pb.ListCommentRequest, resp *pb.ListCommentResponse) error {
	logger.Info("PostService ListComments")

	// 按 father_id 列出子评论（二级评论分页）
	if req.GetFatherId() != 0 {
		return s.listSubComments(req, resp)
	}

	// 按 target_id + target_type 列出一级评论
	if req.GetTargetId() == 0 || req.GetTargetType() == "" {
		return errno.ServerErr(errno.ErrBadRequest, "target_id and target_type required")
	}

	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		pageSize = constvar.DefaultPageSize
	} else if pageSize > constvar.MaxPageSize {
		pageSize = constvar.MaxPageSize
	}
	limit := pageSize + 1

	var pageToken pb.CommentPageToken
	err := pagetoken.DecodePageToken(req.GetPageToken(), &pageToken)
	if err != nil {
		return errno.ServerErr(errno.ErrBadRequest, "invalid page token")
	}

	var pageTokenPtr *pb.CommentPageToken
	if req.GetPageToken() != "" {
		pageTokenPtr = &pageToken

		if pageTokenPtr.GetSortType() != req.GetSortType() {
			return errno.ServerErr(errno.ErrBadRequest, "page token sort type mismatch")
		}
	}

	switch req.GetSortType() {
	case constvar.SortByNewest:
		return s.listCommentsNewest(req.GetTargetId(), req.GetTargetType(), pageTokenPtr, limit, resp)
	case constvar.SortByHottest:
		return s.listCommentsHottest(req.GetTargetId(), req.GetTargetType(), pageTokenPtr, limit, resp)
	default:
		return errno.ServerErr(errno.ErrBadRequest, "sort type not legal")
	}
}

// listSubComments 按 father_id 分页获取子评论（二级评论"查看更多"）
func (s *PostService) listSubComments(req *pb.ListCommentRequest, resp *pb.ListCommentResponse) error {
	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		pageSize = constvar.DefaultPageSize
	} else if pageSize > constvar.MaxPageSize {
		pageSize = constvar.MaxPageSize
	}
	limit := pageSize + 1

	var pageToken pb.CommentPageToken
	err := pagetoken.DecodePageToken(req.GetPageToken(), &pageToken)
	if err != nil {
		return errno.ServerErr(errno.ErrBadRequest, "invalid page token")
	}

	var pageTokenPtr *pb.CommentPageToken
	if req.GetPageToken() != "" {
		pageTokenPtr = &pageToken
		if pageTokenPtr.GetSortType() != req.GetSortType() {
			return errno.ServerErr(errno.ErrBadRequest, "page token sort type mismatch")
		}
	}

	switch req.GetSortType() {
	case constvar.SortByNewest:
		return s.listSubCommentsNewest(req.GetFatherId(), pageTokenPtr, limit, resp)
	case constvar.SortByHottest:
		return s.listSubCommentsHottest(req.GetFatherId(), pageTokenPtr, limit, resp)
	default:
		return errno.ServerErr(errno.ErrBadRequest, "sort type not legal")
	}
}

func (s *PostService) listSubCommentsCommon(
	fatherID uint32, pageToken *pb.CommentPageToken, limit uint32, resp *pb.ListCommentResponse,
	fetch func(token *pb.CommentPageToken, limit uint32) ([]*dao.CommentInfo, error),
	nextToken func(last *dao.CommentInfo) *pb.CommentPageToken,
) error {
	comments, err := fetch(pageToken, limit)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.HasMore = false
	if uint32(len(comments)) > limit-1 {
		resp.HasMore = true
		comments = comments[:limit-1]
	}

	if len(comments) == 0 {
		resp.PageToken = ""
		resp.Comments = nil
		return nil
	}

	if resp.HasMore {
		tokenStr, err := pagetoken.EncodePageToken(nextToken(comments[len(comments)-1]))
		if err != nil {
			return errno.ServerErr(errno.InternalServerError, "failed to encode page token")
		}
		resp.PageToken = tokenStr
	} else {
		resp.PageToken = ""
	}

	resp.Comments = make([]*pb.CommentInfo, len(comments))
	for i, c := range comments {
		resp.Comments[i] = commentInfoToPB(c, 0, nil)
	}
	return nil
}

func (s *PostService) listSubCommentsNewest(fatherID uint32, pageToken *pb.CommentPageToken, limit uint32, resp *pb.ListCommentResponse) error {
	return s.listSubCommentsCommon(fatherID, pageToken, limit, resp,
		func(token *pb.CommentPageToken, limit uint32) ([]*dao.CommentInfo, error) {
			if token == nil {
				return s.Dao.ListSubCommentsNewest(fatherID, limit)
			}
			return s.Dao.ListSubCommentsNewestWithCursor(fatherID, token.GetId(), token.GetCreateTime().AsTime(), limit)
		},
		func(last *dao.CommentInfo) *pb.CommentPageToken {
			return &pb.CommentPageToken{
				Id:         last.Id,
				CreateTime: timestamppb.New(last.CreatedAt),
				SortType:   constvar.SortByNewest,
			}
		},
	)
}

func (s *PostService) listSubCommentsHottest(fatherID uint32, pageToken *pb.CommentPageToken, limit uint32, resp *pb.ListCommentResponse) error {
	return s.listSubCommentsCommon(fatherID, pageToken, limit, resp,
		func(token *pb.CommentPageToken, limit uint32) ([]*dao.CommentInfo, error) {
			if token == nil {
				return s.Dao.ListSubCommentsHottest(fatherID, limit)
			}
			return s.Dao.ListSubCommentsHottestWithCursor(fatherID, token.GetId(), token.GetLikeNum(), limit)
		},
		func(last *dao.CommentInfo) *pb.CommentPageToken {
			return &pb.CommentPageToken{
				Id:       last.Id,
				LikeNum:  last.LikeNum,
				SortType: constvar.SortByHottest,
			}
		},
	)
}

func (s *PostService) listCommentsCommon(
	targetID uint32, targetType string, pageToken *pb.CommentPageToken, limit uint32, resp *pb.ListCommentResponse,
	fetch func(token *pb.CommentPageToken, limit uint32) ([]*dao.CommentInfo, error),
	nextToken func(last *dao.CommentInfo) *pb.CommentPageToken,
) error {

	comments, err := fetch(pageToken, limit)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.HasMore = false
	if uint32(len(comments)) > limit-1 {
		resp.HasMore = true
		comments = comments[:limit-1]
	}

	if len(comments) == 0 {
		resp.PageToken = ""
		resp.Comments = nil
		return nil
	}

	if resp.HasMore {
		tokenStr, err := pagetoken.EncodePageToken(nextToken(comments[len(comments)-1]))
		if err != nil {
			return errno.ServerErr(errno.InternalServerError, "failed to encode page token")
		}
		resp.PageToken = tokenStr
	} else {
		resp.PageToken = ""
	}

	// 批量获取每个主评论的子评论数和前3条子评论
	commentIDs := make([]uint32, len(comments))
	for i, c := range comments {
		commentIDs[i] = c.Id
	}

	subCommentsMap, err := s.Dao.BatchListCommentsByFatherIDs(commentIDs, 3)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	subNumsMap, err := s.Dao.BatchGetCommentNumByFatherIDs(commentIDs)
	if err != nil {
		logger.Error("batch get sub comment num error", logger.String(err.Error()))
	}

	resp.Comments = make([]*pb.CommentInfo, len(comments))
	for i, c := range comments {
		resp.Comments[i] = commentInfoToPB(c, subNumsMap[c.Id], subCommentsMap[c.Id])
	}
	return nil
}

func (s *PostService) listCommentsNewest(targetID uint32, targetType string, pageToken *pb.CommentPageToken, limit uint32, resp *pb.ListCommentResponse) error {
	return s.listCommentsCommon(targetID, targetType, pageToken, limit, resp,
		func(token *pb.CommentPageToken, limit uint32) ([]*dao.CommentInfo, error) {
			if token == nil {
				return s.Dao.ListPrimaryCommentsNewest(targetID, targetType, limit)
			}
			return s.Dao.ListPrimaryCommentsNewestWithCursor(
				targetID, targetType,
				token.GetId(),
				token.GetCreateTime().AsTime(),
				limit,
			)
		},
		func(last *dao.CommentInfo) *pb.CommentPageToken {
			return &pb.CommentPageToken{
				Id:         last.Id,
				CreateTime: timestamppb.New(last.CreatedAt),
				SortType:   constvar.SortByNewest,
			}
		},
	)
}

func (s *PostService) listCommentsHottest(targetID uint32, targetType string, pageToken *pb.CommentPageToken, limit uint32, resp *pb.ListCommentResponse) error {
	return s.listCommentsCommon(targetID, targetType, pageToken, limit, resp,
		func(token *pb.CommentPageToken, limit uint32) ([]*dao.CommentInfo, error) {
			if token == nil {
				return s.Dao.ListPrimaryCommentsHottest(targetID, targetType, limit)
			}
			return s.Dao.ListPrimaryCommentsHottestWithCursor(
				targetID, targetType,
				token.GetId(),
				token.GetLikeNum(),
				limit,
			)
		},
		func(last *dao.CommentInfo) *pb.CommentPageToken {
			return &pb.CommentPageToken{
				Id:       last.Id,
				LikeNum:  last.LikeNum,
				SortType: constvar.SortByHottest,
			}
		},
	)
}

func commentInfoToPB(c *dao.CommentInfo, subNum uint32, subComments []*dao.CommentInfo) *pb.CommentInfo {
	pbSubComments := make([]*pb.CommentInfo, len(subComments))
	for i, sc := range subComments {
		pbSubComments[i] = &pb.CommentInfo{
			Id:            sc.Id,
			TypeName:      sc.TypeName,
			Content:       sc.Content,
			FatherId:      sc.FatherId,
			CreateTime:    timestamppb.New(sc.CreatedAt),
			CreatorId:     sc.CreatorId,
			CreatorName:   sc.CreatorName,
			CreatorAvatar: sc.CreatorAvatar,
			LikeNum:       sc.LikeNum,
			ImgUrl:        sc.ImgUrl,
			TargetId:      sc.TargetID,
			TargetType:    sc.TargetType,
		}
	}

	return &pb.CommentInfo{
		Id:            c.Id,
		TypeName:      c.TypeName,
		Content:       c.Content,
		FatherId:      c.FatherId,
		CreateTime:    timestamppb.New(c.CreatedAt),
		CreatorId:     c.CreatorId,
		CreatorName:   c.CreatorName,
		CreatorAvatar: c.CreatorAvatar,
		LikeNum:       c.LikeNum,
		ImgUrl:        c.ImgUrl,
		TargetId:      c.TargetID,
		TargetType:    c.TargetType,
		SubNum:        subNum,
		SubComments:   pbSubComments,
	}
}

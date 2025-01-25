package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"sync"
)

func (s *PostService) ListMainPost(_ context.Context, req *pb.ListMainPostRequest, resp *pb.ListPostResponse) error {

	logger.Info("PostService ListMainPost")

	filter := &dao.PostModel{
		Category: req.Category,
	}

	if req.Domain != constvar.AllDomain {
		filter.Domain = req.Domain
	}

	tag, err := s.Dao.GetTagByContent(req.Tag)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	posts, err := s.Dao.ListMainPost(filter, req.Filter, req.Offset, req.Limit, req.LastId, req.Pagination, req.SearchContent, tag.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.Posts = make([]*pb.Post, len(posts))

	// 使用 WaitGroup 来等待所有 goroutine 完成
	var wg sync.WaitGroup
	// 用于存放结果的一个临时切片，以便保持顺序
	result := make([]*pb.Post, len(posts))

	for i, post := range posts {
		wg.Add(1)

		// 启动一个 goroutine 处理每个帖子
		go func(i int, post *dao.PostInfo) {
			defer wg.Done()

			// 获取每个帖子的附加信息
			isLiked, isCollection, likeNum, tags, commentNum, collectionNum := s.getPostInfo(post.Id, req.UserId)

			// 更新点赞数（如果有）
			if likeNum != 0 {
				post.LikeNum = likeNum
			}

			// 将结果存入临时结果切片
			result[i] = &pb.Post{
				Id:            post.Id,
				Title:         post.Title,
				Time:          post.LastEditTime,
				Category:      post.Category,
				CreatorId:     post.CreatorId,
				CreatorName:   post.CreatorName,
				CreatorAvatar: post.CreatorAvatar,
				LikeNum:       post.LikeNum,
				CommentNum:    commentNum,
				IsLiked:       isLiked,
				IsCollection:  isCollection,
				Tags:          tags,
				ContentType:   post.ContentType,
				Summary:       post.Summary,
				CollectionNum: collectionNum,
				Domain:        post.Domain,
			}
		}(i, post) // 将i和post传递给goroutine

	}

	// 等待所有 goroutine 完成
	wg.Wait()

	// 将结果从临时切片复制到最终响应
	resp.Posts = result

	return nil
}

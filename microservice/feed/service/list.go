package service

import (
	"context"
	"forum-feed/dao"
	pb "forum-feed/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
)

// List ... feed列表
func (s *FeedService) List(_ context.Context, req *pb.ListRequest, res *pb.ListResponse) error {
	logger.Info("FeedService List")

	// 普通用户，只能返回有权限访问的 projects
	if req.Role == constvar.Normal {
		// projectIds, err = GetFilterFromProjectService(req.UserId)
		// if err != nil {
		// 	return errno.ServerErr(errno.ErrGetDataFromRPC, err.Error())
		// }
	}

	// 筛选条件
	var filter = &dao.FeedModel{
		UserId: req.UserId,
		// GroupId:    req.Filter.GroupId,
		// ProjectIds: projectIds,
	}

	feeds, err := s.Dao.List(filter, req.Offset, req.Limit, req.LastId, req.Pagination)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	// 数据格式化
	list, err := FormatListData(feeds)
	if err != nil {
		// return errno.ServerErr(errno.ErrFormatList, err.Error())
	}

	res.List = list

	return nil
}

func FormatListData(list []*dao.FeedModel) ([]*pb.FeedItem, error) {
	var result []*pb.FeedItem
	var time string
	var sourceId uint32

	for index, feed := range list {

		data := &pb.FeedItem{
			Id:          feed.Id,
			Action:      feed.Action,
			ShowDivider: false,
			CreateTime:  feed.CreateTime,
			User: &pb.User{
				Name:      feed.UserName,
				Id:        feed.UserId,
				AvatarUrl: feed.UserAvatar,
			},
			Source: &pb.Source{
				TypeName: feed.SourceTypeName,
				Id:       feed.SourceObjectId,
				Name:     feed.SourceObjectName,
			},
		}

		// showDivider --> 分割线
		// 需要分割的情况
		// 1.第一条数据 2.不同日期 3.不同项目
		if index == 0 {
			time = data.CreateTime
			sourceId = data.Source.Id
			data.ShowDivider = true
		} else if time != data.CreateTime {
			time = data.CreateTime
			data.ShowDivider = true
		} else if sourceId != data.Source.Id {
			sourceId = data.Source.Id
			data.ShowDivider = true
		}

		result = append(result, data)
	}

	return result, nil
}

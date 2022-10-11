package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"forum/util"
	"gorm.io/gorm"
	"strconv"
)

func (s *PostService) CreateReport(_ context.Context, req *pb.CreateReportRequest, _ *pb.Response) error {
	logger.Info("PostService CreateReport")

	isReported, err := s.Dao.IsUserHadReportPost(req.UserId, req.PostId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errno.NotFoundErr(errno.ErrItemNotFound, "post-"+strconv.Itoa(int(req.PostId)))
		}

		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if isReported {
		return errno.ServerErr(errno.ErrRepeatReport, "")
	}

	report := &dao.ReportModel{
		UserId:     req.UserId,
		CreateTime: util.GetCurrentTime(),
		PostId:     req.PostId,
		TypeName:   req.TypeName,
		Cause:      req.Cause,
	}

	if err := s.Dao.CreateReport(report); err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	// auto ban when count >= constvar.BanNumber
	go func() {
		count, err := s.Dao.GetReportNumByPostId(req.PostId)
		if err != nil {
			logger.Error(errno.ErrDatabase.Error(), logger.String(err.Error()))
		}

		if count >= constvar.BanNumber {
			post, err := s.Dao.GetPost(req.PostId)
			if err != nil {
				logger.Error(errno.ErrDatabase.Error(), logger.String(err.Error()))
			}

			post.IsReport = true
			if err := post.Save(); err != nil {
				logger.Error(errno.ErrDatabase.Error(), logger.String(err.Error()))
			}
		}
	}()

	return nil
}

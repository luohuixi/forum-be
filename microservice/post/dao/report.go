package dao

import (
	pb "forum-post/proto"
	"forum/pkg/constvar"
)

type ReportModel struct {
	Id         uint32
	UserId     uint32
	CreateTime string
	PostId     uint32
	TypeName   string
	Cause      string
}

func (ReportModel) TableName() string {
	return "reports"
}

// Create ...
func (r *ReportModel) Create() error {
	return dao.DB.Create(r).Error
}

func (r *ReportModel) Delete() error {
	return dao.DB.Delete(r).Error
}

func (Dao) CreateReport(report *ReportModel) error {
	return report.Create()
}

func (d Dao) GetReport(id uint32) (*ReportModel, error) {
	var report ReportModel
	err := d.DB.Table("reports").Where("id = ?", id).First(&report).Error

	return &report, err
}

func (d Dao) ValidReport(postId uint32) error {
	tx := d.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Table("reports").Where("post_id = ?", postId).Delete(&ReportModel{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := d.DeletePost(postId); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (d Dao) InValidReport(id, postId uint32) error {
	if err := d.DB.Table("reports").Where("id = ?", id).Delete(&ReportModel{}).Error; err != nil {
		return err
	}

	count, err := d.GetReportNumByPostId(postId)
	if err != nil {
		return err
	}

	if count == constvar.BanNumber-1 { // cancel auto ban
		post, err := d.GetPost(postId)
		if err != nil {
			return err
		}

		post.IsReport = false
		return post.Save(dao.DB)
	}

	return nil
}

func (d *Dao) ListReport(offset, limit, lastId uint32, pagination bool) ([]*pb.Report, error) {
	var reports []*pb.Report
	query := d.DB.Table("reports").Select("reports.id id, reports.post_id, reports.user_id, reports.create_time, cause, reports.type_name, u.name user_name, u.avatar user_avatar, u2.id be_reported_user_id, u2.name be_reported_user_name, p.title post_title").Joins("join users u on u.id = reports.user_id").Joins("join posts p on p.id = reports.post_id").Joins("join users u2 on u2.id = p.creator_id")

	if pagination {
		if limit == 0 {
			limit = constvar.DefaultLimit
		}

		query = query.Offset(int(offset)).Limit(int(limit))

		if lastId != 0 {
			query = query.Where("reports.id < ?", lastId)
		}
	}

	err := query.Scan(&reports).Error

	return reports, err
}

func (d *Dao) GetReportNumByPostId(postId uint32) (uint32, error) {
	var count int64
	err := d.DB.Model(&ReportModel{}).Where("post_id = ?", postId).Count(&count).Error
	return uint32(count), err
}

func (d *Dao) IsUserHadReportPost(userId uint32, postId uint32) (bool, error) {
	_, err := d.GetPost(postId)
	if err != nil {
		return false, err
	}

	var count int64
	if err := d.DB.Table("reports").Where("user_id = ? AND post_id = ?", userId, postId).Count(&count).Error; err != nil {
		return false, err
	}

	return count != 0, nil
}

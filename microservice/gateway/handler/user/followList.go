package user

import (
	. "forum-gateway/handler"
	"forum-gateway/util"
	"forum/log"
	"forum/model"
	"forum/pkg/errno"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type followCountRow struct {
	Id    uint32
	Count uint32
}

func parseFollowListQuery(c *gin.Context) (uint32, int, int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		SendError(c, errno.ErrPathParam, nil, "id must be a positive integer", GetLine())
		return 0, 0, 0, false
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if err != nil {
		SendError(c, errno.ErrQuery, nil, err.Error(), GetLine())
		return 0, 0, 0, false
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		SendError(c, errno.ErrQuery, nil, err.Error(), GetLine())
		return 0, 0, 0, false
	}
	if page < 0 {
		page = 0
	}

	return uint32(id), limit, page * limit, true
}

func listFollowUsers(c *gin.Context, following bool) {
	targetID, limit, offset, ok := parseFollowListQuery(c)
	if !ok {
		return
	}

	viewerID := c.MustGet("userId").(uint32)
	db := model.GetSelfDB()
	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close()
	}

	users := make([]FollowListUser, 0)
	var query string
	if following {
		query = `
SELECT u.id, u.name, COALESCE(u.avatar, '') AS avatar, u.role, COALESCE(u.signature, '') AS signature
FROM user_follows f
JOIN users u ON u.id = f.followee_id AND COALESCE(u.re, 0) = 0
WHERE f.follower_id = ?
ORDER BY f.created_at DESC, f.id DESC
LIMIT ? OFFSET ?`
	} else {
		query = `
SELECT u.id, u.name, COALESCE(u.avatar, '') AS avatar, u.role, COALESCE(u.signature, '') AS signature
FROM user_follows f
JOIN users u ON u.id = f.follower_id AND COALESCE(u.re, 0) = 0
WHERE f.followee_id = ?
ORDER BY f.created_at DESC, f.id DESC
LIMIT ? OFFSET ?`
	}
	if err := db.Raw(query, targetID, limit, offset).Scan(&users).Error; err != nil {
		SendError(c, errno.ErrDatabase, nil, err.Error(), GetLine())
		return
	}

	if len(users) == 0 {
		SendResponse(c, nil, FollowListResponse{Users: users})
		return
	}

	userIDs := make([]uint32, 0, len(users))
	for _, item := range users {
		userIDs = append(userIDs, item.Id)
	}

	followingCounts := map[uint32]uint32{}
	followerCounts := map[uint32]uint32{}
	var rows []followCountRow
	if err := db.Raw(`
SELECT follower_id AS id, COUNT(*) AS count
FROM user_follows
WHERE follower_id IN ?
GROUP BY follower_id`, userIDs).Scan(&rows).Error; err != nil {
		SendError(c, errno.ErrDatabase, nil, err.Error(), GetLine())
		return
	}
	for _, row := range rows {
		followingCounts[row.Id] = row.Count
	}
	rows = nil
	if err := db.Raw(`
SELECT followee_id AS id, COUNT(*) AS count
FROM user_follows
WHERE followee_id IN ?
GROUP BY followee_id`, userIDs).Scan(&rows).Error; err != nil {
		SendError(c, errno.ErrDatabase, nil, err.Error(), GetLine())
		return
	}
	for _, row := range rows {
		followerCounts[row.Id] = row.Count
	}

	followingMap := map[uint32]bool{}
	if viewerID != 0 {
		rows = nil
		if err := db.Raw(`
SELECT followee_id AS id, 1 AS count
FROM user_follows
WHERE follower_id = ? AND followee_id IN ?`, viewerID, userIDs).Scan(&rows).Error; err != nil {
			SendError(c, errno.ErrDatabase, nil, err.Error(), GetLine())
			return
		}
		for _, row := range rows {
			followingMap[row.Id] = true
		}
	}

	for idx := range users {
		users[idx].FollowingCount = followingCounts[users[idx].Id]
		users[idx].FollowerCount = followerCounts[users[idx].Id]
		users[idx].IsFollowing = followingMap[users[idx].Id]
	}

	SendResponse(c, nil, FollowListResponse{Users: users})
}

// ListFollowing ... 获取用户关注列表
// @Summary 获取用户关注列表 api
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param id path int true "user_id"
// @Param limit query string false "limit"
// @Param page query string false "page"
// @Success 200 {object} FollowListResponse
// @Router /user/following/{id} [get]
func ListFollowing(c *gin.Context) {
	log.Info("User ListFollowing function called.", zap.String("X-Request-Id", util.GetReqID(c)))
	listFollowUsers(c, true)
}

// ListFollowers ... 获取用户粉丝列表
// @Summary 获取用户粉丝列表 api
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param id path int true "user_id"
// @Param limit query string false "limit"
// @Param page query string false "page"
// @Success 200 {object} FollowListResponse
// @Router /user/followers/{id} [get]
func ListFollowers(c *gin.Context) {
	log.Info("User ListFollowers function called.", zap.String("X-Request-Id", util.GetReqID(c)))
	listFollowUsers(c, false)
}

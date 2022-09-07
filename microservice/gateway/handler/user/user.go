package user

// TeamLoginRequest login 请求
type TeamLoginRequest struct {
	OauthCode string `json:"oauth_code"`
} // @name TeamLoginRequest

// StudentLoginRequest StudentLogin 请求
type StudentLoginRequest struct {
	StudentId string `json:"student_id" binding:"required"`
	Password  string `json:"password" binding:"required"`
} // @name StudentLoginRequest

// TeamLoginResponse login 请求响应
type TeamLoginResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
} // @name TeamLoginResponse

// StudentLoginResponse login 请求响应
type StudentLoginResponse struct {
	Token string `json:"token"`
} // @name StudentLoginResponse

// GetInfoRequest 获取 info 请求
type GetInfoRequest struct {
	Ids []uint32 `json:"ids" binding:"required"`
} // @name GetInfoRequest

type userInfo struct {
	Id        uint32 `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
	Signature string `json:"signature"`
}

// GetInfoResponse 获取 info 响应
type GetInfoResponse struct {
	List []userInfo `json:"list"`
} // @name GetInfoResponse

// GetProfileRequest 获取 profile 请求
type GetProfileRequest struct {
	Id uint32 `json:"id"`
} // @name GetProfileRequest

// UserProfile 获取 profile 响应
type UserProfile struct {
	Id     uint32 `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Email  string `json:"email"`
	Role   string `json:"role"`
} // @name UserProfile

// ListRequest 获取 userList 请求
type ListRequest struct {
	Team  uint32 `json:"team"`
	Group uint32 `json:"group"`
} // @name ListRequest

type user struct {
	Id     uint32 `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	Role   string `json:"role"`
} // @name user

// ListResponse 获取 userList 响应
type ListResponse struct {
	Count uint32 `json:"count"`
	List  []user `json:"list"`
} // @name ListResponse

// UpdateInfoRequest 更新 userInfo 请求
type UpdateInfoRequest struct {
	userInfo
} // @name UpdateInfoRequest

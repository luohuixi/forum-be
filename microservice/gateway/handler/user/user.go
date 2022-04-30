package user

// teamLoginRequest login 请求
type teamLoginRequest struct {
	OauthCode string `json:"oauth_code"`
} // @name teamLoginRequest

// studentLoginRequest StudentLogin 请求
type studentLoginRequest struct {
	StudentId string `json:"student_id"`
	Password  string `json:"password"`
} // @name studentLoginRequest

// teamLoginResponse login 请求响应
type teamLoginResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
} // @name teamLoginResponse

// studentLoginResponse login 请求响应
type studentLoginResponse struct {
	Token string `json:"token"`
} // @name studentLoginResponse

// GetInfoRequest 获取 info 请求
type getInfoRequest struct {
	Ids []uint32 `json:"ids" binding:"required"`
} // @name GetInfoRequest

type userInfo struct {
	Id        uint32 `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}

// GetInfoResponse 获取 info 响应
type getInfoResponse struct {
	List []userInfo `json:"list"`
} // @name getInfoResponse

// getProfileRequest 获取 profile 请求
type getProfileRequest struct {
	Id uint32 `json:"id"`
} // @name getProfileRequest

// userProfile 获取 profile 响应
type userProfile struct {
	Id     uint32 `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Email  string `json:"email"`
	Role   uint32 `json:"role"`
} // @name userProfile

// listRequest 获取 userList 请求
type listRequest struct {
	Team  uint32 `json:"team"`
	Group uint32 `json:"group"`
} // @name listRequest

type user struct {
	Id     uint32 `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	Role   uint32 `json:"role"`
} // @name user

// listResponse 获取 userList 响应
type listResponse struct {
	Count uint32 `json:"count"`
	List  []user `json:"list"`
} // @name listResponse

// updateInfoRequest 更新 userInfo 请求
type updateInfoRequest struct {
	userInfo
} // @name updateInfoRequest

// updateTeamGroupRequest
type updateTeamGroupRequest struct {
	Ids   []uint32 `json:"ids"`
	Value uint32   `json:"value"`
	Kind  uint32   `json:"kind"`
} // @name updateTeamGroupRequest

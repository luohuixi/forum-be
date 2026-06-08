package user

// TeamLoginRequest login 请求
type TeamLoginRequest struct {
	OauthCode string `json:"oauth_code"`
} // @name TeamLoginRequest

// StudentLoginRequest StudentLogin 请求
type StudentLoginRequest struct {
	StudentId        string `json:"student_id"`
	Password         string `json:"password"`
	Action           string `json:"action"`
	SessionId        string `json:"session_id"`
	Captcha          string `json:"captcha"`
	SecondAuthMethod string `json:"second_auth_method"`
	SecondAuthCode   string `json:"second_auth_code"`
	Provider         string `json:"provider"`
	OauthCode        string `json:"oauth_code"`
	CallbackURL      string `json:"callback_url"`
} // @name StudentLoginRequest

// TeamLoginResponse login 请求响应
type TeamLoginResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
} // @name TeamLoginResponse

// StudentLoginResponse login 请求响应
type StudentLoginResponse struct {
	RedirectURL                string   `json:"redirect_url"`
	Token                      string   `json:"token"`
	SessionId                  string   `json:"session_id"`
	Status                     string   `json:"status"`
	Message                    string   `json:"message"`
	CaptchaImageBase64         string   `json:"captcha_image_base64"`
	AvailableSecondAuthMethods []string `json:"available_second_auth_methods"`
	CurrentSecondAuthMethod    string   `json:"current_second_auth_method"`
	SecondAuthSMSTarget        string   `json:"second_auth_sms_target"`
	SecondAuthEmailTarget      string   `json:"second_auth_email_target"`
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
	Id                        uint32 `json:"id"`
	Name                      string `json:"name"`
	Avatar                    string `json:"avatar"`
	Email                     string `json:"email"`
	Role                      string `json:"role"`
	Signature                 string `json:"signature"`
	IsPublicFeed              bool   `json:"is_public_feed"`
	IsPublicCollectionAndLike bool   `json:"is_public_collection_and_like"`
	FollowingCount            uint32 `json:"following_count"`
	FollowerCount             uint32 `json:"follower_count"`
	IsFollowing               bool   `json:"is_following"`
} // @name UserProfile

// MyProfile 获取 my profile 响应
type MyProfile struct {
	Id                        uint32 `json:"id"`
	Name                      string `json:"name"`
	Avatar                    string `json:"avatar"`
	Email                     string `json:"email"`
	StudentId                 string `json:"student_id"`
	Role                      string `json:"role"`
	Signature                 string `json:"signature"`
	IsPublicFeed              bool   `json:"is_public_feed"`
	IsPublicCollectionAndLike bool   `json:"is_public_collection_and_like"`
	FollowingCount            uint32 `json:"following_count"`
	FollowerCount             uint32 `json:"follower_count"`
	IsFollowing               bool   `json:"is_following"`
} // @name MyProfile

type FollowRequest struct {
	TargetUserID uint32 `json:"target_user_id" binding:"required"`
} // @name FollowRequest

type FollowResponse struct {
	IsFollowing    bool   `json:"is_following"`
	FollowingCount uint32 `json:"following_count"`
	FollowerCount  uint32 `json:"follower_count"`
} // @name FollowResponse

type FollowListUser struct {
	Id             uint32 `json:"id"`
	Name           string `json:"name"`
	Avatar         string `json:"avatar"`
	Role           string `json:"role"`
	Signature      string `json:"signature"`
	FollowingCount uint32 `json:"following_count"`
	FollowerCount  uint32 `json:"follower_count"`
	IsFollowing    bool   `json:"is_following"`
} // @name FollowListUser

type FollowListResponse struct {
	Users []FollowListUser `json:"users"`
} // @name FollowListResponse

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
	Name                      string `json:"name"`
	AvatarURL                 string `json:"avatar_url"`
	Signature                 string `json:"signature"`
	IsPublicFeed              bool   `json:"is_public_feed"`
	IsPublicCollectionAndLike bool   `json:"is_public_collection_and_like"`
} // @name UpdateInfoRequest

type ListMessageResponse struct {
	Messages []string `json:"messages"`
}

type CreateMessageRequest struct {
	Message string `json:"message" binding:"required"`
}

type CreatePrivateMessageRequest struct {
	ReceiveUserid  uint32 `json:"receive_userid" binding:"required"`
	Type           string `json:"type" binding:"required"` // comment/like/collection/reply_comment
	Content        string `json:"content"`
	PostId         uint32 `json:"post_id" binding:"required"`
	CommentId      uint32 `json:"comment_id"`
	PostTitle      string `json:"post_title" binding:"required"`
	CommentContent string `json:"comment_content"`
}

type AddRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

package models

// TikTok OAuth2 Token Response Data
type TokenResponseData struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int64  `json:"expires_in"`
	OpenID           string `json:"open_id"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
}

// TikTok OAuth2 Token Response
type TokenResponse struct {
	Data        TokenResponseData `json:"data"`
	ErrorCode   int               `json:"error_code"`
	Description string            `json:"description"`
}

// TikTok OAuth2 Error Response
type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	LogID            string `json:"log_id"`
}

// Auth Request State (CSRF koruması için)
type AuthState struct {
	State   string `json:"state"`
	Created int64  `json:"created"`
}

// TikTok User Object (from User Info API)
type UserInfo struct {
	OpenID           string `json:"open_id"`
	UnionID          string `json:"union_id"`
	AvatarURL        string `json:"avatar_url"`
	AvatarURL100     string `json:"avatar_url_100"`
	AvatarLargeURL   string `json:"avatar_large_url"`
	DisplayName      string `json:"display_name"`
	BioDescription   string `json:"bio_description"`
	ProfileDeepLink  string `json:"profile_deep_link"`
	IsVerified       bool   `json:"is_verified"`
	Username         string `json:"username"`
	FollowerCount    int64  `json:"follower_count"`
	FollowingCount   int64  `json:"following_count"`
	LikesCount       int64  `json:"likes_count"`
	VideoCount       int64  `json:"video_count"`
}

// TikTok User Info API Response
type UserInfoResponse struct {
	Data        UserInfo `json:"data"`
	ErrorCode   int      `json:"error_code"`
	Description string   `json:"description"`
}

// Combined response with token and user info
type AuthResponse struct {
	Token    TokenResponseData `json:"token"`
	UserInfo UserInfo          `json:"user_info"`
}

// API Response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

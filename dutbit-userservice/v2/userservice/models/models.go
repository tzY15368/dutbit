package models

var Test = 1

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type UpdateRequest struct {
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
type RegisterDocument struct {
	Username        string                 `json:"username" binding:"required"`
	Email           string                 `json:"email" binding:"required"`
	Password        string                 `json:"password" binding:"required"`
	Role            int                    `json:"role"`
	Site            map[string]interface{} `json:"site"`
	Created_at      int64                  `json:"created_at"`
	Ip              string                 `json:"ip"`
	Last_login_time int64                  `json:"last_login_time"`
	Last_login_ip   string                 `json:"last_login_ip"`
	Confirmation    map[string]interface{} `json:"confirmation"`
}

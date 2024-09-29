package service

type UserLoginRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
}
type UserRegisterRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
	NickName string `json:"nickname"`
	Role     uint32 `json:"role"`
}

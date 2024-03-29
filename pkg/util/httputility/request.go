package httputility

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ChangePasswordRequest struct {
	Oldpassword string `json:"oldpassword"`
	Newpassword string `json:"newpassword"`
}

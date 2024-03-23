package httputility

type SendMessage struct {
	Message string `json:"message"`
}

type IDResponse struct {
	UUID string `json:"uuid"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type UpdateResponse struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

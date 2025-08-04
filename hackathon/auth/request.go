package auth

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	UserId string `json:"userId"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresAt   string `json:"expiresAt"`
	UserId      string `json:"userId"`
}

type VerifyRequest struct {
	AccessToken string
}

type VerifyResponse struct {
	UserId string
}

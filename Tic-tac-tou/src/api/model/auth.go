package model

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	ID    string `json:"id"`
	Login string `json:"login"`
}

type JwtRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type JwtResponse struct {
	Type         string `json:"type"`
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type RefreshJwtRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Login string `json:"login"`
}

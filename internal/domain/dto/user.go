package dto

type SignInInput struct {
	Code string `json:"code" binding:"required"`
}

type AccessTokenOutput struct {
	Token string `json:"token"`
}

package dto

type SignInInput struct {
	Code string `json:"code"`
}

type AccessTokenOutput struct {
	Token string `json:"token"`
}

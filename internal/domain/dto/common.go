package dto

type UserInfo struct {
	ID         string `json:"id" binding:"uuid"`
	Username   string `json:"username"`
	ProfileURL string `json:"profileURL"`
	Role       string `json:"role"`
}

type IDInput struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type IDOutput struct {
	ID string `json:"id" binding:"uuid"`
}

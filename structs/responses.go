package structs

type MessageResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type DefaultResponse[T any] struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type BearerStruct struct {
	UserId       string `json:"userId"`
	Type         string `json:"type"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Exp          string `json:"expired"`
}

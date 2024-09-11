package structs

type MessageResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type DefaultResponseWithData[T any] struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type DefaultResponseMessageOnly struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type BearerStruct struct {
	UserId       string `json:"userId"`
	Type         string `json:"type"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Exp          string `json:"expired"`
}

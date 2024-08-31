package structs

import "github.com/gin-gonic/gin"

type SuccessMessage struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    gin.H  `json:"data"`
}

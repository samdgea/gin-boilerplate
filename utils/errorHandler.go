package utils

import "github.com/gin-gonic/gin"

func ThrowError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"error":   true,
		"message": message,
	})
}

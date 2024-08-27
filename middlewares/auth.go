package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samdgea/gin-boilerplate/db"
	"github.com/samdgea/gin-boilerplate/models"
	"github.com/samdgea/gin-boilerplate/utils"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strings"
)

func IsAuth(c *gin.Context) {
	var token string
	var user models.UserModel

	token = c.GetHeader("Authorization")
	if token == "" {
		utils.ThrowError(c, http.StatusUnauthorized, "Need Header Auth")
		c.Abort()
		return
	}

	header := strings.Split(token, " ")
	tokenHeader := header[0]
	if tokenHeader != "Bearer" {
		utils.ThrowError(c, http.StatusUnauthorized, "Invalid Header Auth")
		c.Abort()
		return
	}

	tokenHeader = header[1]
	if tokenHeader == "" {
		utils.ThrowError(c, http.StatusUnauthorized, "Invalid Header Auth")
		c.Abort()
		return
	}

	decodeToken, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		utils.ThrowError(c, http.StatusUnauthorized, "Invalid Token")
		c.Abort()
		return
	}

	userId := decodeToken.Claims.(jwt.MapClaims)["userId"].(float64)

	if err = db.DB.First(&user, int(userId)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ThrowError(c, http.StatusUnauthorized, "User not found")
		} else {
			utils.ThrowError(c, http.StatusInternalServerError, "Failed to get user data")
		}
		c.Abort()
		return
	}

	c.Set("user", user)
	c.Set("userId", userId)
	c.Next()
}

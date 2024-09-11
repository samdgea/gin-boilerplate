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
	"time"
)

func IsAuth(c *gin.Context) {
	var token string
	var user models.UserModel
	var tokenModel models.TokenModel

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

	userId := decodeToken.Claims.(jwt.MapClaims)["userId"].(string)
	tokenId := decodeToken.Claims.(jwt.MapClaims)["tokenId"].(string)

	if err = db.DB.First(&user, "id = ?", userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ThrowError(c, http.StatusUnauthorized, "User not found")
		} else {
			utils.ThrowError(c, http.StatusInternalServerError, "Failed to get user data")
		}
		c.Abort()
		return
	}

	if err = db.DB.First(&tokenModel, "user_id = ? AND token_id = ?", userId, tokenId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ThrowError(c, http.StatusUnauthorized, "Incorrect credentials")
		} else {
			utils.ThrowError(c, http.StatusInternalServerError, err.Error())
		}
		c.Abort()
		return
	}

	currentTime := time.Now()

	if !tokenModel.IsActive || currentTime.After(tokenModel.ExpiresAt) {
		db.DB.Model(&tokenModel).Update("IsActive", false)

		utils.ThrowError(c, http.StatusUnauthorized, "Access Token is revoked or expired")
		c.Abort()
		return
	}

	c.Set("user", user)
	c.Set("userId", userId)
	c.Next()
}

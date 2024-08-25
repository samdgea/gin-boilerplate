package controllers

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samdgea/gin-boilerplate/db"
	"github.com/samdgea/gin-boilerplate/models"
	"github.com/samdgea/gin-boilerplate/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

func Login(c *gin.Context) {
	type LoginRequest struct {
		Username string `json:"userName" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var loginRequest LoginRequest
	var user models.UserModel

	err := c.ShouldBindJSON(&loginRequest)
	if err != nil {
		utils.ThrowError(c, http.StatusBadRequest, "Invalid Request")
		return
	}

	if err = db.DB.QueryRow("SELECT id, username, password, isActive FROM users WHERE username = $1", loginRequest.Username).Scan(&user.Id, &user.Username, &user.Password, &user.IsActive); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ThrowError(c, http.StatusUnauthorized, "Incorrect credentials")
		} else {
			utils.ThrowError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Compare BCrypt hash password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		utils.ThrowError(c, http.StatusUnauthorized, "Incorrect credentials")
		return
	}

	// Check if user is Active
	if !user.IsActive {
		utils.ThrowError(c, http.StatusUnauthorized, "Account was disabled, contact your Administrator")
		return
	}

	sign := jwt.New(jwt.SigningMethodHS256)
	claims := sign.Claims.(jwt.MapClaims)
	claims["userId"] = user.Id
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // 1 Day

	token, err := sign.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		utils.ThrowError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Login Success",
		"data": gin.H{
			"token":   token,
			"type":    "Bearer",
			"expired": "1 Day",
		},
	})
}

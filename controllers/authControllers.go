package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samdgea/gin-boilerplate/db"
	"github.com/samdgea/gin-boilerplate/models"
	"github.com/samdgea/gin-boilerplate/structs"
	"github.com/samdgea/gin-boilerplate/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

	// Use Gorm to find the user by username
	if err = db.DB.Where("username = ?", loginRequest.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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
	claims["userId"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // 1 Day

	token, err := sign.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		utils.ThrowError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := structs.SuccessMessage{
		Error:   false,
		Message: "Login Success",
		Data: gin.H{
			"token":   token,
			"type":    "Bearer",
			"expired": "1 Day",
		},
	}

	c.JSON(http.StatusOK, response)
}

func Me(c *gin.Context) {
	user := c.MustGet("user").(models.UserModel)
	userId := c.GetString("userId")

	response := structs.SuccessMessage{
		Error:   false,
		Message: "Success",
		Data: gin.H{
			"userId":   userId,
			"fullName": user.FirstName + " " + user.LastName,
		},
	}

	c.JSON(http.StatusOK, response)
}

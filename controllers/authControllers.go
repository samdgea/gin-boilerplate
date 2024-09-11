package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/samdgea/gin-boilerplate/db"
	"github.com/samdgea/gin-boilerplate/models"
	"github.com/samdgea/gin-boilerplate/structs"
	"github.com/samdgea/gin-boilerplate/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strings"
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

	accessToken := jwt.New(jwt.SigningMethodHS256)
	atClaims := accessToken.Claims.(jwt.MapClaims)
	atClaims["userId"] = user.ID
	atClaims["tokenId"] = uuid.New()
	atClaims["exp"] = time.Now().Add(time.Hour * 1).Unix() // 1 Hour

	// Add entry to TokenModel
	db.DB.Create(&models.TokenModel{
		UserID:    atClaims["userId"].(uuid.UUID),
		TokenID:   atClaims["tokenId"].(uuid.UUID),
		ExpiresAt: time.Unix(atClaims["exp"].(int64), 0),
	})

	aToken, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		utils.ThrowError(c, http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["userId"] = user.ID
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix() // 24 Hours

	rToken, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		utils.ThrowError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := structs.DefaultResponseWithData[structs.BearerStruct]{
		Error:   false,
		Message: "Login Success",
		Data: structs.BearerStruct{
			UserId:       user.ID.String(),
			Type:         "Bearer",
			AccessToken:  aToken,
			RefreshToken: rToken,
			Exp:          time.Unix(atClaims["exp"].(int64), 0).Format(time.RFC1123),
		},
	}

	c.JSON(http.StatusOK, response)
}

func RefreshToken(c *gin.Context) {
	type tokenReqBody struct {
		RefreshToken string `json:"refresh_token"`
	}

	var tokenReq tokenReqBody
	var user models.UserModel

	err := c.ShouldBindJSON(&tokenReq)
	if err != nil {
		utils.ThrowError(c, http.StatusBadRequest, "Invalid Request")
		return
	}

	decodeToken, err := jwt.Parse(tokenReq.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		fmt.Println(err)
		utils.ThrowError(c, http.StatusUnauthorized, "Invalid Token")
		c.Abort()
		return
	}

	if claims, ok := decodeToken.Claims.(jwt.MapClaims); ok && decodeToken.Valid {
		if err = db.DB.First(&user, "id = ?", claims["userId"]).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.ThrowError(c, http.StatusUnauthorized, "User not found")
			} else {
				utils.ThrowError(c, http.StatusInternalServerError, "Failed to get user data")
			}
			c.Abort()
			return
		}
	}

	accessToken := jwt.New(jwt.SigningMethodHS256)
	newAccessToken := accessToken.Claims.(jwt.MapClaims)
	newAccessToken["userId"] = user.ID
	newAccessToken["tokenId"] = uuid.New()
	newAccessToken["exp"] = time.Now().Add(time.Hour * 1).Unix()

	// Add entry to TokenModel
	db.DB.Create(&models.TokenModel{
		UserID:    newAccessToken["userId"].(uuid.UUID),
		TokenID:   newAccessToken["tokenId"].(uuid.UUID),
		ExpiresAt: time.Unix(newAccessToken["exp"].(int64), 0),
	})

	aToken, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		utils.ThrowError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := structs.DefaultResponseWithData[structs.BearerStruct]{
		Error:   false,
		Message: "New Access Token has been issued",
		Data: structs.BearerStruct{
			UserId:       user.ID.String(),
			Type:         "Bearer",
			AccessToken:  aToken,
			RefreshToken: tokenReq.RefreshToken,
			Exp:          time.Unix(newAccessToken["exp"].(int64), 0).Format(time.RFC1123),
		},
	}

	c.JSON(http.StatusOK, response)
}

func Me(c *gin.Context) {
	user := c.MustGet("user").(models.UserModel)

	type DataResponse struct {
		UserId   uuid.UUID `json:"userId"`
		FullName string    `json:"fullName"`
	}

	data := DataResponse{
		UserId:   user.ID,
		FullName: user.FirstName + " " + user.LastName,
	}

	response := structs.DefaultResponseWithData[DataResponse]{
		Error:   false,
		Message: "Success",
		Data:    data,
	}

	c.JSON(http.StatusOK, response)
}

func Logout(c *gin.Context) {
	var token string
	var tokenModel models.TokenModel

	token = c.GetHeader("Authorization")
	header := strings.Split(token, " ")

	decodeToken, err := jwt.Parse(header[1], func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		utils.ThrowError(c, http.StatusUnauthorized, "Invalid Token")
		c.Abort()
		return
	}

	tokenId := decodeToken.Claims.(jwt.MapClaims)["tokenId"].(string)

	if err = db.DB.First(&tokenModel, "token_id = ?", tokenId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ThrowError(c, http.StatusUnauthorized, "Access Token Session was not found")
		} else {
			utils.ThrowError(c, http.StatusInternalServerError, "Failed to get Access Token Session")
		}
		c.Abort()
		return
	}

	// Revoke Access Token
	db.DB.Model(&tokenModel).Update("IsActive", false)

	c.JSON(http.StatusOK, structs.DefaultResponseMessageOnly{
		Error:   false,
		Message: "Access Token has been revoked successfully",
	})
}

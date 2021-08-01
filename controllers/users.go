package controllers

import (
	"context"
	"mail-service/service"
	"mail-service/structs"
	"mail-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func signUp(c *gin.Context) {
	user := structs.User{}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	hashedPassword, err := utils.HashedPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed created password", err.Error(),
		})
		return
	}

	user.Passwrod = hashedPassword
	user.TokenHash = utils.GenerateRandomString(15)

	if err = initApp.mainDB.Collection("users").InsertOne(context.TODO(), structs.User{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed created password", err.Error(),
		})
		return
	}

	from := "vikisquarez@gmail.com"
	to := []string{user.Email}
	subject := "Email Verification for Bookite"
	mailType := service.MailConfirmation
	mailData := &service.MailData{
		Username: user.Username,
		Code:     utils.GenerateRandomString(8),
	}
}

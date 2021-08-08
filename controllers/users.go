package controllers

import (
	"context"
	"github.com/google/uuid"
	"mail-service/global"
	"mail-service/service"
	"mail-service/structs"
	"mail-service/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	mailService service.MailService
}

func SignUp(c *gin.Context) {
	var ah *AuthHandler
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
			"message": "Failed created password" + err.Error(),
		})
		return
	}

	id := uuid.New().String()
	user.ID = id
	user.Password = hashedPassword
	user.TokenHash = utils.GenerateRandomString(15)
	user.CreatedAt = time.Now().Unix()

	_, err = global.MainDB.Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed created password" + err.Error(),
		})
		return
	}

	from := "testingtrash13@gmail.com"
	to := []string{user.Email}
	subject := "Email Verification for Bookite"
	mailType := service.MailConfirmation
	mailData := &service.MailData{
		Username: user.Username,
		Code:     utils.GenerateRandomString(8),
	}

	mailReq := ah.mailService.NewMail(from, to, subject, mailType, mailData)
	err = ah.mailService.SendMail(mailReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error to send mail",
			"error":   err.Error(),
			"data":    "",
		})
		return
	}

	verification := &structs.VerificationData{
		UserID:    id,
		Email:     user.Email,
		Code:      mailData.Code,
		Type:      structs.MailConfirmation,
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)),
	}

	_, err = global.MainDB.Collection("users_verification").InsertOne(context.TODO(), verification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error to store mail verification data",
			"error":   err.Error(),
			"data":    "",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Please verify your email account using the confirmation code send to your mail",
		"error":   "",
		"data":    verification,
	})
}

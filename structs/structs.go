package structs

import "time"

type User struct {
	ID         string `json:"id" bson:"id"`
	Email      string `json:"email" bson:"email"`
	Password   string `json:"password" bson:"password"`
	Username   string `json:"username" bson:"username"`
	TokenHash  string `json:"token_hash" bson:"token_hash"`
	IsVerified bool   `json:"is_verified" bson:"is_verified"`
	CreatedAt  int64  `json:"created_at" bson:"created_at"`
	UpdatedAt  int64  `json:"updated_at" bson:"updated_at"`
}

type VerificationData struct {
	UserID    string               `json:"user_id" bson:"user_id"`
	Email     string               `json:"email" bson:"email"`
	Code      string               `json:"code" bson:"code"`
	ExpiresAt time.Time            `json:"expiresat" bson:"expiresat"`
	Type      VerificationDataType `json:"type" bson:"type"`
}

type VerificationDataType int

const (
	MailConfirmation VerificationDataType = iota + 1
	PassReset
)

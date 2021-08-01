package structs

type User struct {
	ID         string `json:"id" bson:"id"`
	Email      string `json:"email" bson:"email"`
	Password   string `json:"password" bson:"password"`
	Username   string `json:"username" bson:"username"`
	TokenHash  string `json:"tokenhash" bson:"tokenhash"`
	IsVerified bool   `json:"isverified" bson:"isverified"`
	CreatedAt  int64  `json:"createdat" bson:"createdat"`
	UpdatedAt  int64  `json:"updatedat" bson:"updatedat"`
}

type VerificationData struct {
	Email     string               `json:"email" bson:"email"`
	Code      string               `json:"code" bson:"code"`
	ExpiresAt int64                `json:"expiresat" bson:"expiresat"`
	Type      VerificationDataType `json:"type" bson:"type"`
}

type VerificationDataType int

const (
	MailConfirmation VerificationDataType = iota + 1
	PassReset
)

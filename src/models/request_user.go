package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type UserLoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=3,max=20"`
}

type UserCustomClaims struct {
	*User
	jwt.RegisteredClaims
}

type UserLoginResponse struct {
	*User
	Token    string `json:"token"`
	ExpireAt int64  `json:"expiresAt"`
}

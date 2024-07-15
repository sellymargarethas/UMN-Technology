package models

import "github.com/golang-jwt/jwt"

type Login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Id           int    `json:"id"`
	Nama         string `json:"nama"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type JWTClaim struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

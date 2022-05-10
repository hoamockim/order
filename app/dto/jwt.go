package dto

import "github.com/dgrijalva/jwt-go"

type AuthClaims struct {
	jwt.StandardClaims
	UserId string `json:"userId"`
}

package models

import "github.com/golang-jwt/jwt"

type AppClaims struct {
	UserId             string `json:"userid"`
	jwt.StandardClaims        // jwt for user login response
}

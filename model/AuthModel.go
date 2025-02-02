package model

import "github.com/dgrijalva/jwt-go"

type User struct {
	Nik      string `json:"nik"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Token    string `json:"token"`
}

type TokenClaim struct {
	Nik  string `json:"nik"`
	Role string `json:"role"`
	jwt.StandardClaims
}

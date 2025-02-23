package model

import (
	"github.com/golang-jwt/jwt/v4"
)

type QRClaims struct {
	UUID   string
	Random string
	jwt.RegisteredClaims
}

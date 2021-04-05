package contract

import (
	"github.com/dgrijalva/jwt-go"
)

type JWT interface {
	GenerateToken(appKey, appSecret string) (string, error)
	CheckToken(token string) (*CustomClaims, error)
}

type CustomClaims struct {
	jwt.StandardClaims
	AppKey    string
	AppSecret string
}

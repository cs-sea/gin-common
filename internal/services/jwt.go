package services

import (
	"time"

	"github.com/pkg/errors"

	"github.com/dgrijalva/jwt-go"

	"github.com/cs-sea/gin-common/contract"
)

type JWTConfig struct {
	encryptType string
	Secret      string
	Expires     time.Duration
}

var _ contract.JWT = &JWTService{}

type JWTService struct {
	config *JWTConfig
}

func NewJWTService(config *JWTConfig) contract.JWT {
	return &JWTService{
		config: config,
	}
}

func (j *JWTService) GenerateToken(appKey, appSecret string) (string, error) {
	now := time.Now()

	expires := now.Add(j.config.Expires)

	claims := &contract.CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
		},
		AppKey:    appKey,
		AppSecret: appSecret,
	}

	tokenClaims := jwt.NewWithClaims(j.GetJwtEncryptAlgorithm(), claims)

	token, err := tokenClaims.SignedString(j.config.Secret)

	return token, err
}

func (j *JWTService) CheckToken(token string) (*contract.CustomClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &contract.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.config.Secret, nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "解析token 失败")
	}

	if tokenClaims == nil {
		return nil, errors.New("payload is null")
	}

	if claims, ok := tokenClaims.Claims.(*contract.CustomClaims); ok {
		return claims, nil
	}

	return nil, errors.New("token 断言失败")
}

func (j *JWTService) GetJwtEncryptAlgorithm() jwt.SigningMethod {
	return jwt.GetSigningMethod(j.config.encryptType)
}

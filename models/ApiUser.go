package models

import "gorm.io/gorm"

type ApiUser struct {
	gorm.Model
	AppKey    string
	AppSecret string
}

func (a *ApiUser) TableName() string {
	return "api_users"
}

package config

import (
	"p4_web/model"
	"p4_web/tools/env"
	"time"
)

type AuthConfigStruct struct {
	PrimaryKey string
	Uids       []string
	Model      *model.User
	Keep       time.Duration
	SecretKey  string
}

var AuthConfig *AuthConfigStruct

func init() {
	AuthConfig = &AuthConfigStruct{
		PrimaryKey: "Id",
		Uids:       []string{"mail"},
		Keep:       time.Hour * 24 * 7,
		SecretKey:  env.Get("SECRET_KEY"),
	}
}

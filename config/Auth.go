package config

import (
	"p4_web/model"
	"p4_web/tools/env"
	"time"
)

type Guard struct {
	PrimaryKey string
	Uids       []string
	Model      *model.User
	Keep       time.Duration
	SecretKey  string
}

type AuthConfigStruct struct {
	Guard  string
	Guards map[string]Guard
}

func (a AuthConfigStruct) Guarder(guard ...string) Guard {
	var currentGuard string
	if len(guard) == 0 {
		currentGuard = a.Guard
	} else {
		currentGuard = guard[0]
	}
	return a.Guards[currentGuard]
}

var AuthConfig *AuthConfigStruct

func init() {
	AuthConfig = &AuthConfigStruct{
		Guard: "user",
		Guards: map[string]Guard{
			"user": Guard{
				PrimaryKey: "Id",
				Uids:       []string{"mail"},
				Keep:       time.Hour * 24 * 7,
				SecretKey:  env.Get("SECRET_KEY"),
			},
		},
	}
}

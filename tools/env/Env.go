package env

import (
	"os"
	"p4_web/tools/exception"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(exception.ApiException{
			Code:    []int{-2},
			Message: "Error loading .env file",
		})
	}
}

func Get(key string) string {
	return os.Getenv(key)
}

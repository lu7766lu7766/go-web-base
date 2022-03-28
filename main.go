package main

import (
	"p4_web/router"
)

var host = ":8080"

func main() {
	r := router.InitRouter()
	r.Run(host)
}

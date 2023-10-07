package main

import (
	"meow/flags"
	"meow/handling"

	"github.com/gin-gonic/gin"
)

func main() {
	domain := flags.Domain
	if *domain == "" || *domain == "YOUR_DOMAIN_HERE" {
		panic("You have not specified a domain!")
	}
	r := gin.Default()
	r.GET("/t", handling.HandleTikTokRequest)
	r.GET("/", handling.HandleIndex)

	r.Run(":4232")
}

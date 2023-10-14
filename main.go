package main

import (
	"meow/config"
	"meow/handling"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	domain := config.Domain
	if domain == "" || domain == "YOUR_DOMAIN_HERE" {
		panic("You have not specified a domain!")
	}
	println("Starting server on port " + config.Port)
	println("Domain: " + domain)
	println("Public?: " + strconv.FormatBool(config.Public))

	r := gin.Default()
	r.GET("/t", handling.HandleRequest)
	r.GET("/s", handling.HandleSoundCollageRequest)
	r.GET("/", handling.HandleIndex)
	r.GET("/collage-:id", handling.HandleDirectCollage)
	r.GET("/video-:id", handling.HandleDirectVideo)

	r.Run(config.Port)
}

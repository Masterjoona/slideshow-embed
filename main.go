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
	println("Public: " + strconv.FormatBool(config.Public))
	println("Sound route: " + strconv.FormatBool(config.IsffmpegInstalled))

	r := gin.Default()
	r.GET("/", handling.HandleIndex)

	r.GET("/t", handling.HandleRequest)
	r.GET("/collage-:id", handling.HandleDirectCollage)

	if config.IsffmpegInstalled {
		r.GET("/s", handling.HandleSoundCollageRequest)
		r.GET("/video-:id", handling.HandleDirectVideo)
	}

	r.Run(config.Port)
}

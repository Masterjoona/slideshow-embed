package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	InitEnvs()
	if Domain == "" || Domain == "YOUR_Domain_HERE" {
		panic("You have not specified a Domain!")
	}
	println("Starting server on port " + Port)
	println("Domain: " + Domain)
	println("Public: " + strconv.FormatBool(Public))
	println("Limit public amount: " + strconv.Itoa(LimitPublicAmount))
	println("Sound route: " + strconv.FormatBool(IsffmpegInstalled))
	println("Fancy Slideshow: " + strconv.FormatBool(FancySlideshow))

	r := gin.Default()
	r.GET("/", HandleIndex)

	r.GET("/t", HandleRequest)
	r.GET("/collage-:id", HandleDirectFile)

	if IsffmpegInstalled {
		r.GET("/s", HandleSoundCollageRequest)
		r.GET("/video-:id", HandleDirectFile)
	}
	if IsffmpegInstalled && FancySlideshow {
		r.GET("/f", HandleFancySlideshowRequest)
		r.GET("/slide-:id", HandleDirectFile)
	}

	r.Run(Port)
}

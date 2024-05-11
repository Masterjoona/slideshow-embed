// go:build scrape
package main

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	InitEnvs()
	if Domain == "" {
		panic("You have not specified a Domain!")
	}
	if len(InstallIds) == 0 && !Scraping {
		panic("You are trying to use tiktok api without any install ids!")
	}
	println("Starting server on port " + Port)
	println("Domain: " + Domain)
	println("Scraping: " + strconv.FormatBool(Scraping))
	println("Install IDs: " + strings.Join(InstallIds, ", ") + Ternary(Scraping, " (Ignored)", ""))
	println("Public: " + strconv.FormatBool(Public))
	println("Limit public amount: " + strconv.Itoa(LimitPublicAmount))
	println("Sound route: " + strconv.FormatBool(IsffmpegInstalled))
	println("Fancy Slideshow: " + strconv.FormatBool(FancySlideshow))

	r := gin.Default()
	r.GET("/", HandleIndex)

	r.GET(PathCollage, HandleRequest)
	r.GET("/collage-:id", HandleDirectFile("collage"))

	if IsffmpegInstalled {
		r.GET(PathCollageSound, HandleSoundCollageRequest)
		r.GET("/video-:id", HandleDirectFile("video"))
	}
	if IsffmpegInstalled && FancySlideshow {
		r.GET(PathSlide, HandleFancySlideshowRequest)
		r.GET("/slide-:id", HandleDirectFile("slide"))
	}
	if Downloader {
		r.GET(PathDownloader, HandleDownloader)
	}

	r.Run(Port)
}

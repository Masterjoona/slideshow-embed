// go:build scrape
package main

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	InitEnvs()
	if Domain == "" {
		panic("You have not specified a Domain!")
	}

	println("Starting server on port " + Port)
	println("Domain: " + Domain)
	println("Provider: " + TiktokProvider)
	println("Public: " + strconv.FormatBool(Public))
	println("Limit public amount: " + strconv.Itoa(LimitPublicAmount))
	println("Sound route: " + strconv.FormatBool(IsffmpegInstalled))
	println("Fancy Slideshow: " + strconv.FormatBool(FancySlideshow))

	r := gin.Default()
	r.GET("/", HandleIndex)

	r.GET(PathCollage, HandleRequest)
	r.GET("/collage-:id", HandleDirectFile("collage"))

	r.GET(PathJson, HandleJsonRequest)
	r.GET(PathVideoProxy, HandleVideoProxy)

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

	if Subtitler {
		r.GET(PathSubs, HandleSubtitleVideo)
		r.GET("/subs-:id", HandleDirectFile("subs"))
	}

	go func() {
		for {
			RecentTiktokReqs.Flush()
			ShortURLCache.Flush()
			time.Sleep(5 * time.Minute)
		}
	}()

	r.Run(Port)
}

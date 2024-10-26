package main

import (
	"meow/pkg/config"
	"meow/pkg/handler"
	"meow/pkg/net"
	"meow/pkg/providers"
	"meow/pkg/vars"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitEnvs()
	if config.Domain == "" {
		panic("You have not specified a Domain!")
	}

	println("Starting server on port " + config.Port)
	println("Domain: " + config.Domain)
	println("Provider: " + config.TiktokProvider)
	println("Public: " + strconv.FormatBool(config.Public))
	println("Limit public amount: " + strconv.Itoa(config.LimitPublicAmount))
	println("Sound route: " + strconv.FormatBool(config.IsffmpegInstalled))
	println("Fancy Slideshow: " + strconv.FormatBool(config.FancySlideshow))

	r := gin.Default()
	r.GET("/", handler.HandleIndex)

	r.GET(vars.PathCollage, handler.HandleRequest)
	r.GET("/collage-:id", handler.HandleDirectFile("collage"))

	r.GET(vars.PathJson, handler.HandleJsonRequest)
	r.GET(vars.PathVideoProxy+"/:id", handler.HandleVideoProxy)

	if config.IsffmpegInstalled {
		r.GET(vars.PathCollageSound, handler.HandleSoundCollageRequest)
		r.GET("/video-:id", handler.HandleDirectFile("video"))
	}
	if config.IsffmpegInstalled && config.FancySlideshow {
		r.GET(vars.PathSlide, handler.HandleFancySlideshowRequest)
		r.GET("/slide-:id", handler.HandleDirectFile("slide"))
	}
	if config.Downloader {
		r.GET(vars.PathDownloader, handler.HandleDownloader)
	}

	if config.Subtitler {
		r.GET(vars.PathSubs, handler.HandleSubtitleVideo)
		r.GET("/subs-:id", handler.HandleDirectFile("subs"))
	}

	go func() {
		for {
			providers.RecentTiktokReqs.Flush()
			net.ShortURLCache.Flush()
			time.Sleep(5 * time.Minute)
		}
	}()

	r.Run(config.Port)
}

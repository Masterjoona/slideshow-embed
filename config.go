package main

import (
	"math/rand"
	"os"
	"strconv"
	"time"
)

var Domain string
var Port string
var LocalStats Stats

var LimitPublicAmount int
var Public bool
var IsffmpegInstalled bool
var FancySlideshow bool

func addTrailingSlash(s string) string {
	if s != "" && s[len(s)-1] != '/' {
		return s + "/"
	}
	return s
}

func InitEnvs() {
	rand.NewSource(time.Now().UnixNano())
	Domain = addTrailingSlash(os.Getenv("DOMAIN"))
	Public = os.Getenv("PUBLIC") == "true"
	IsffmpegInstalled = os.Getenv("FFMPEG") == "true"
	FancySlideshow = os.Getenv("FANCY_SLIDESHOW") == "true"

	if LimitPublicAmount = 0; os.Getenv("LIMIT_PUBLIC_AMOUNT") != "" {
		LimitPublicAmount, _ = strconv.Atoi(os.Getenv("LIMIT_PUBLIC_AMOUNT"))
	}

	if Port = os.Getenv("PORT"); Port == "" {
		Port = "4232"
	}
	if Port[0] != ':' {
		Port = ":" + Port
	}

	UpdateLocalStats()

}

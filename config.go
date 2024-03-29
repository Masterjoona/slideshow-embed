package main

import (
	"math/rand"
	"os"
	"strconv"
	"time"
)

var (
	Domain            string
	Port              string
	LocalStats        Stats
	InstallIds        []string
	LimitPublicAmount int
	Public            bool
	IsffmpegInstalled bool
	FancySlideshow    bool
)

func addTrailingSlash(s string) string {
	if s != "" && s[len(s)-1] != '/' {
		return s + "/"
	}
	return s
}

func checkEnvOrDefault(env string, def string) string {
	if val := os.Getenv(env); val != "" {
		return val
	}
	return def
}

func InitEnvs() {
	rand.NewSource(time.Now().UnixNano())
	Domain = addTrailingSlash(os.Getenv("DOMAIN"))
	Public = os.Getenv("PUBLIC") == "true"
	IsffmpegInstalled = os.Getenv("FFMPEG") == "true"
	FancySlideshow = os.Getenv("FANCY_SLIDESHOW") == "true"

	LimitPublicAmount, _ = strconv.Atoi(os.Getenv("LIMIT_PUBLIC_AMOUNT"))

	Port = checkEnvOrDefault("PORT", ":4232")

	if installId := os.Getenv("INSTALL_ID"); installId != "" {
		InstallIds = []string{installId}
	} else {
		// thanks yt-dlp love you <3 (and tiktxk)
		InstallIds = []string{
			"7351144126450059040",
			"7351149742343391009",
			"7351153174894626592",
		}
	}

	UpdateLocalStats()
}

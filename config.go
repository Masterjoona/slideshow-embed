package main

import (
	"os"
	"os/exec"
	"strconv"
)

var Domain string
var Port string
var LocalStats Stats

var SlideIndex bool
var LimitPublicAmount int
var Public bool
var IsffmpegInstalled bool
var FancySlideshow bool

func checkBinary(bin string) bool {
	_, err := exec.LookPath(bin)
	return err == nil
}

func addTrailingSlash(s string) string {
	if s != "" && s[len(s)-1] != '/' {
		return s + "/"
	}
	return s
}

func InitEnvs() {
	Domain = addTrailingSlash(os.Getenv("DOMAIN"))
	Public = os.Getenv("PUBLIC") == "true"
	FancySlideshow = os.Getenv("FANCY_SLIDESHOW") == "true"
	// SlideIndex = os.Getenv("SLIDE_INDEX") == "true" iw as gonna do this before they nuked the html

	if LimitPublicAmount = 0; os.Getenv("LIMIT_PUBLIC_AMOUNT") != "" {
		LimitPublicAmount, _ = strconv.Atoi(os.Getenv("LIMIT_PUBLIC_AMOUNT"))
	}

	if Port = os.Getenv("PORT"); Port == "" {
		Port = "4232"
	}
	if Port[0] != ':' {
		Port = ":" + Port
	}

	IsffmpegInstalled = checkBinary("/usr/bin/ffmpeg") || checkBinary("/usr/local/bin/ffmpeg")

	if os.Getenv("FFMPEG") == "true" {
		IsffmpegInstalled = true
	} else if os.Getenv("FFMPEG") == "false" {
		IsffmpegInstalled = false
	}

	UpdateLocalStats()

}

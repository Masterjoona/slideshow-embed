package config

import (
	"os"
)

var Domain string
var Public bool
var Port string
var IsffmpegInstalled bool

func init() {
	Domain = os.Getenv("DOMAIN")
	Public = os.Getenv("PUBLIC") == "true"
	if Port = os.Getenv("PORT"); Port == "" {
		Port = "4232"
	}
	if Port[0] != ':' {
		Port = ":" + Port
	}

	_, err := os.Stat("/usr/bin/ffmpeg")
	if err == nil {
		IsffmpegInstalled = true
	} else {
		IsffmpegInstalled = false
	}

	if os.Getenv("FFMPEG") == "true" {
		IsffmpegInstalled = true
	} else if os.Getenv("FFMPEG") == "false" {
		IsffmpegInstalled = false
	}
}

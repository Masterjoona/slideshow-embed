package config

import (
	"os"
	"os/exec"
)

var Domain string
var Public bool
var Port string
var IsffmpegInstalled bool

func checkBinary(bin string) bool {
	_, err := exec.LookPath(bin)
	return err == nil
}

func init() {
	Domain = os.Getenv("DOMAIN")
	Public = os.Getenv("PUBLIC") == "true"
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
}

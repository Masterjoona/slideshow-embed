package main

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func addString(s string, r string, trailing bool) string {
	if (trailing && !strings.HasSuffix(s, r)) || (!trailing && !strings.HasPrefix(s, r)) {
		if trailing {
			return s + r
		}
		return r + s
	}
	return s
}

func checkEnvOrDefault(env string, def string) string {
	if val := os.Getenv(env); val != "" {
		return val
	}
	return def
}

func commaSeparatedToStrs(s string) []string {
	if s == "" {
		return []string{}
	}
	strs := strings.Split(s, ",")
	for i, str := range strs {
		strs[i] = strings.TrimSpace(str)
	}
	return strs
}

func InitEnvs() {
	rand.NewSource(time.Now().UnixNano())
	Domain = addString(os.Getenv("DOMAIN"), "/", true)
	Public = os.Getenv("PUBLIC") == "true"
	Downloader = os.Getenv("DOWNLOADER") == "true"
	IsffmpegInstalled = os.Getenv("FFMPEG") == "true"
	FancySlideshow = os.Getenv("FANCY_SLIDESHOW") == "true"

	LimitPublicAmount, _ = strconv.Atoi(os.Getenv("LIMIT_PUBLIC_AMOUNT"))

	Port = addString(checkEnvOrDefault("PORT", ":4232"), ":", false)

	InstallIds = commaSeparatedToStrs(checkEnvOrDefault(
		"INSTALL_IDS",
		"",
	))

	DeviceIds = commaSeparatedToStrs(checkEnvOrDefault(
		"DEVICE_IDS",
		"7351044760062330401",
	))
	// my own - I use my own iid as well
	// while testing, two iids returned different responses?

	UpdateLocalStats()
}

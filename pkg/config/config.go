package config

import (
	"math/rand"
	"meow/pkg/files"
	"meow/pkg/util"
	"os"
	"strconv"
	"time"
)

var (
	Domain             string
	Port               string
	Downloader         bool
	TemporaryDirectory = util.Ternary(isDocker, "/tmp", "tmp")
	LocalStats         Stats
	LimitPublicAmount  int
	Public             bool
	IsffmpegInstalled  bool
	FancySlideshow     bool
	Subtitler          bool
	TiktokProvider     string

	fileSize, _  = files.GetFileSize("/.dockerenv")
	isDocker     = fileSize > -1
	PythonServer = "http://" + util.Ternary(isDocker, "photo_collager", "localhost") + ":9700"
)

func InitEnvs() {
	rand.NewSource(time.Now().UnixNano())

	Domain = addString(os.Getenv("DOMAIN"), "/", true)
	Public = os.Getenv("PUBLIC") == "true"
	Downloader = os.Getenv("DOWNLOADER") == "true"
	IsffmpegInstalled = os.Getenv("FFMPEG") == "true"
	FancySlideshow = os.Getenv("FANCY_SLIDESHOW") == "true"
	Subtitler = os.Getenv("SUBTITLER") == "true"
	TiktokProvider = os.Getenv("TIKTOK_PROVIDER")

	LimitPublicAmount, _ = strconv.Atoi(os.Getenv("LIMIT_PUBLIC_AMOUNT"))

	Port = addString(checkEnvOrDefault("PORT", ":4232"), ":", false)

	UpdateLocalStats()
}

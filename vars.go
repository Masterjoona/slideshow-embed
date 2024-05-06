package main

import "regexp"

const (
	million  = 1000000
	thousand = 1000
	digits   = "0123456789"
	hexChars = "0123456789abcdef"

	KB = 1 << 10
	MB = 1 << 20
	GB = 1 << 30

	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"

	PathCollage      = "/t"
	PathCollageSound = "/s"
	PathSlide        = "/f"
)

var (
	ErrorImages = []string{
		"https://media.discordapp.net/attachments/961445186280509451/980132677338423316/fuckmedaddyharderohyeailovecokcimsocissyfemboy.gif",
		"https://media.discordapp.net/attachments/901959319719936041/996927812927750264/chrome_2WOKI6Jm3v.gif",
		"https://cdn.discordapp.com/attachments/749030987530502197/980338691706880051/79587A35-FD36-41D3-8232-7A29B46D2543.gif",
	}
	ErrorImagesIndex = 0
)

var PythonServer = "http://" + Ternary(isDocker(), "photo_collager", "localhost") + ":9700"

var (
	longLinkRe     = regexp.MustCompile(`https:\/\/(?:www.)?(?:vxtiktok|tiktok|tiktxk|)\.com\/@.{2,32}\/(?:photo|video)\/(\d+)`)
	shortLinkRe    = regexp.MustCompile(`https:\/\/.{1,3}\.(?:(?:vx|)tikt(?:x|o)k)\.com/(?:.{1,2}/|)(.{5,12})`)
	AudioSrcRe     = regexp.MustCompile(`<a href="(.*)" onclick="bdl\(this, event\)" type="audio`)
	VideoSrcLinkRe = regexp.MustCompile(`<a href="(.*)" onclick="bdl\(this, event\)" rel`)
)

var (
	Domain            string
	Port              string
	LocalStats        Stats
	InstallIds        []string
	DeviceIds         []string
	LimitPublicAmount int
	Public            bool
	IsffmpegInstalled bool
	FancySlideshow    bool
)

package vars

import (
	"net/http"
	"time"
)

const (
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"

	PathCollage      = "/t"
	PathCollageSound = "/s"
	PathSlide        = "/f"
	PathDownloader   = "/d"
	PathSubs         = "/subs"
	PathJson         = "/json"
	PathVideoProxy   = "/vproxy"
	PathTest         = "/test"

	SubtitlesHost       = "https://api16-normal-c-useast2a.tiktokv.com/tiktok/cla/subtitle_translation/get/v1/?"
	SubtitlesHostBackup = "https://www.tiktok.com/tiktok/cla/subtitle_translation/get/v1/?subtitle_id=7441071229936864000&device_platform=web_pc"
)

var HttpClient = &http.Client{Timeout: time.Second * 10}

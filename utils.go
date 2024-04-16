package main

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
)

var errorImages = []string{
	"https://media.discordapp.net/attachments/961445186280509451/980132677338423316/fuckmedaddyharderohyeailovecokcimsocissyfemboy.gif",
	"https://media.discordapp.net/attachments/901959319719936041/996927812927750264/chrome_2WOKI6Jm3v.gif",
	"https://cdn.discordapp.com/attachments/749030987530502197/980338691706880051/79587A35-FD36-41D3-8232-7A29B46D2543.gif",
}
var errorImagesIndex = 0

func ErrorImage() string {
	if errorImagesIndex == 2 {
		errorImagesIndex = 0
	} else {
		errorImagesIndex++
	}
	return errorImages[errorImagesIndex]
}

func validateURL(url string) bool {
	if url == "" {
		return false
	}
	if !strings.Contains(url, ".tiktxk.com") && !strings.Contains(url, ".tiktok.com") {
		return false
	}
	return true
}

func EscapeString(input string) string {
	decoded, err := url.QueryUnescape(input)
	if err != nil {
		return input
	}
	return decoded
}

func SplitURLAndIndex(URL string) (string, string, bool) {
	lastInd := strings.LastIndex(URL, "/")
	index := URL[lastInd+1:]
	if index == "" {
		index = "1"
	}
	sound := strings.HasSuffix(index, "s")
	if sound {
		index = strings.Replace(index, "s", "", 1)
	}
	return URL[:lastInd], index, sound
}

func UpdateLocalStats() {
	collageFiles, err := os.ReadDir("collages")
	if err != nil {
		println("Error while updating local stats: " + err.Error())
		return
	}
	count := 0
	sort.Slice(collageFiles, func(i, j int) bool {
		fileI, err1 := collageFiles[i].Info()
		fileJ, err2 := collageFiles[j].Info()
		if err1 != nil || err2 != nil {
			return collageFiles[i].Name() < collageFiles[j].Name()
		}
		return fileI.ModTime().After(fileJ.ModTime())
	})

	var fileLinks []FileLink

	for _, file := range collageFiles {
		fileLinks = append(fileLinks, FileLink{
			Name: file.Name(),
			Path: Domain + file.Name(),
		})
		count++
	}
	countString := strconv.Itoa(count)
	if LimitPublicAmount > 0 && len(fileLinks) > LimitPublicAmount {
		fileLinks = fileLinks[:LimitPublicAmount]
		countString += fmt.Sprintf(" (limited to %d)", LimitPublicAmount)
	}

	bytes, err := GetDirectorySize("collages")
	if err != nil {
		println("Error while getting size " + err.Error())
		return
	}
	size := FormatSize(bytes)

	LocalStats = Stats{
		FilePaths: fileLinks,
		FileCount: countString,
		TotalSize: size,
	}
}

func isDocker() bool {
	_, err := os.Stat("/.dockerenv")
	return !os.IsNotExist(err)
}

func Ternary(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}

const (
	million  = 1000000
	thousand = 1000
	digits   = "0123456789"
	hexChars = "0123456789abcdef"

	KB = 1 << 10
	MB = 1 << 20
	GB = 1 << 30

	appName            = "musical_ly"
	appId              = 0 //_AID = 0 # aweme = 1128, trill = 1180, musical_ly = 1233, universal = 0;
	appVersion         = "34.1.2"
	appManifestVersion = "2023401020"
	UserAgent          = "com.zhiliaoapp.musically" + appVersion + " (Linux; U; Android 13; en_US; Pixel 7; Build/TD1A.220804.031; Cronet/58.0.2991.0)"
)

var PythonServer = "http://" + Ternary(isDocker(), "photo_collager", "localhost") + ":9700"

func FormatLargeNumbers(numberString string) string {
	number, err := strconv.Atoi(numberString)
	if err != nil {
		return "0"
	}
	switch {
	case number >= million:
		return fmt.Sprintf("%.1fM", float64(number)/million)
	case number >= thousand:
		return fmt.Sprintf("%.1fK", float64(number)/thousand)
	default:
		return fmt.Sprintf("%d", number)
	}
}

func ReverseString(s string) string {
	var reversed strings.Builder
	for i := len(s) - 1; i >= 0; i-- {
		reversed.WriteByte(s[i])
	}
	return reversed.String()
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

/*
	func randomBigInt(min, max int64) int64 {
		return min + rand.Int63n(max-min)
	}
*/
func GenerateRandomHex() string {
	var b strings.Builder
	for i := 0; i < 16; i++ {
		b.WriteByte(hexChars[rand.Intn(len(hexChars))])
	}
	return b.String()
}

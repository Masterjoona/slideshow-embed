package main

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
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

type FileLink struct {
	Name string
	Path string
}

type Stats struct {
	FilePaths []FileLink
	FileCount string
	TotalSize string
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

func ternary(cond bool, a, b string) string {
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

	UserAgent = "com.ss.android.ugc.trill/494+Mozilla/5.0+(Linux;+Android+12;+2112123G+Build/SKQ1.211006.001;+wv)+AppleWebKit/537.36+(KHTML,+like+Gecko)+Version/4.0+Chrome/107.0.5304.105+Mobile+Safari/537.36"
)

var PythonServer = "http://" + ternary(isDocker(), "photo_collager", "localhost") + ":9700"

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

func GenerateRandomString(hexify bool) string {
	var b strings.Builder
	var charset string
	if hexify {
		charset = hexChars
	} else {
		charset = digits
	}
	for i := 0; i < 6; i++ {
		b.WriteByte(charset[rand.Intn(len(charset))])
	}
	return b.String()
}

func GetTimestamp(precision string) string {
	return strconv.FormatInt(time.Now().Unix(), 10) + precision
}

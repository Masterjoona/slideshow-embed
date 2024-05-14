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

func ErrorImage() string {
	if ErrorImagesIndex == 2 {
		ErrorImagesIndex = 0
	} else {
		ErrorImagesIndex++
	}
	return ErrorImages[ErrorImagesIndex]
}

func validateURL(url string) bool {
	if url == "" {
		return false
	}
	if !strings.Contains(url, ".tiktxk.com") && !strings.Contains(url, ".tiktok.com") &&
		!strings.Contains(url, ".vxtiktok.com") {
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
	if LimitPublicAmount > -1 && len(fileLinks) > LimitPublicAmount {
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

/*
	func randomInt(min, max int) int {
		return min + rand.Intn(max-min)
	}

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

func IsAwemeBeingRendered(id string) bool {
	_, ok := CurrentlyRenderingAwemes[id]
	return ok
}

func AddAwemeToRendering(id string) {
	CurrentlyRenderingAwemes[id] = struct{}{}
}

func RemoveAwemeFromRendering(id string) {
	delete(CurrentlyRenderingAwemes, id)
}

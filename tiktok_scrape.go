//go:build scrape

package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

const Scraping = true

func getHash(url string) string {
	reversed := ReverseString(url)

	base64Encoded := base64.StdEncoding.EncodeToString([]byte(reversed))

	hash := md5.Sum([]byte(base64Encoded))

	hashString := fmt.Sprintf("%x", hash)

	return ReverseString(hashString)
}

func fetchTTSave(tiktokUrl, mode, hash string) (*string, error) {
	client := &http.Client{}
	var data = strings.NewReader(fmt.Sprintf(`{"id":"%s","hash":"%s","mode":"%s","locale":"en","loading_indicator_url":"https://ttsave.app/images/slow-down.gif","unlock_url":"https://ttsave.app/en/unlock"}`, tiktokUrl, hash, mode))
	req, err := http.NewRequest("POST", "https://api.ttsave.app/", data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	text := string(bodyText)
	return &text, nil
}

func getData(body *string) (string, string, Counts) {
	authorRe := regexp.MustCompile(`mb-2">(.*)</a>`)
	match := authorRe.FindStringSubmatch(*body)
	author := match[1]

	nicknameRe := regexp.MustCompile(`text-center">(.*)</h2>`)
	match = nicknameRe.FindStringSubmatch(*body)
	author = match[1] + " (" + author + ")"

	captionRe := regexp.MustCompile(`oneliner">(.*)<\/p>`)
	match = captionRe.FindStringSubmatch(*body)
	caption := match[1]

	countRe := regexp.MustCompile(`text-gray-500">(.*)</span>`)
	matches := countRe.FindAllStringSubmatch(*body, -1)

	return author, caption, Counts{
		Views:     matches[0][1],
		Likes:     matches[1][1],
		Comments:  matches[2][1],
		Favorites: matches[3][1],
		Shares:    matches[4][1],
	}
}

func getMediaLink(body *string) string {
	re := regexp.MustCompile(`<a href="(.*)" on`)
	match := re.FindStringSubmatch(*body)
	return match[1]
}

func getSlideLinks(body *string) []string {
	re := regexp.MustCompile(`<img src="([^"]*)">`)
	matches := re.FindAllStringSubmatch(*body, -1)
	var slideLinks []string
	for _, match := range matches {
		slideLinks = append(slideLinks, match[1])
	}
	return slideLinks
}

func getVideoDimensionsFromUrl(videoURL string) (width, height string, err error) {
	cmd := exec.Command(
		"ffprobe",
		"-v",
		"error",
		"-select_streams",
		"v:0",
		"-show_entries",
		"stream=width,height",
		"-of",
		"csv=s=x:p=0",
		videoURL,
	)

	output, err := cmd.Output()
	if err != nil {
		return "0", "0", err
	}

	dimensions := strings.Split(string(output), "x")
	if len(dimensions) != 2 {
		return "0", "0", fmt.Errorf("unexpected output format")
	}

	fmt.Sscanf(dimensions[0], "%s", &width)
	fmt.Sscanf(dimensions[1], "%s", &height)

	return width, height, nil
}

func FetchTiktokData(videoId string) (SimplifiedData, error) {
	url := "https://www.tiktok.com/@placeholder/video/" + videoId
	hash := getHash(url)

	data, err := fetchTTSave(url, "video", hash)
	if err != nil {
		return SimplifiedData{}, err
	}

	slideLinks := getSlideLinks(data)
	author, caption, stats := getData(data)

	if len(slideLinks) == 0 {
		// must be a video?
		isVideo := true
		video := getMediaLink(data)
		width, height, err := getVideoDimensionsFromUrl(video)
		if err != nil {
			return SimplifiedData{}, err
		}
		return SimplifiedData{
			Author:  author,
			Caption: caption,
			Details: stats,
			IsVideo: isVideo,
			Video:   SimplifiedVideo{Url: video, Width: width, Height: height},
		}, nil

	}

	return SimplifiedData{
		Author:     author,
		Caption:    caption,
		Details:    stats,
		SoundUrl:   getMediaLink(data),
		ImageLinks: slideLinks,
		IsVideo:    false,
		Video:      SimplifiedVideo{},
	}, nil
}

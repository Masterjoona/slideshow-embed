package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

const Scraping = "ttsave"

func _(url string) string {
	reversed := ReverseString(url)

	base64Encoded := base64.StdEncoding.EncodeToString([]byte(reversed))

	hash := md5.Sum([]byte(base64Encoded))

	hashString := fmt.Sprintf("%x", hash)

	return ReverseString(hashString)
}

func fetchTTSave(tiktokUrl string) (*string, error) {
	client := &http.Client{}
	var data = strings.NewReader(fmt.Sprintf(`{"language_id":"1","query":"%s"}`, tiktokUrl))
	req, err := http.NewRequest("POST", "https://ttsave.app/download", data)
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

func getData(body *string) (string, string, Counts, error) {
	authorRe := regexp.MustCompile(`mb-2">(.*)</a>`)
	match := authorRe.FindStringSubmatch(*body)
	if len(match) == 0 {
		return "", "", Counts{}, fmt.Errorf("could not find author")
	}
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
	}, nil
}

func getMediaLink(body *string, video bool) string {
	if video {
		match := VideoSrcLinkRe.FindStringSubmatch(*body)
		return match[1]
	}
	match := AudioSrcRe.FindStringSubmatch(*body)
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

func FetchTiktokDataTTSave(videoId string) (SimplifiedData, error) {
	url := "https://www.tiktok.com/@placeholder/video/" + videoId

	data, err := fetchTTSave(url)
	if err != nil {
		return SimplifiedData{}, err
	}

	slideLinks := getSlideLinks(data)
	author, caption, stats, err := getData(data)
	if err != nil {
		return SimplifiedData{}, err
	}

	if len(slideLinks) == 0 {
		// must be a video?
		video := getMediaLink(data, true)
		width, height, err := GetVideoDimensionsFromUrl(video)
		if err != nil {
			return SimplifiedData{}, err
		}
		return SimplifiedData{
			Author:  author,
			Caption: caption,
			Details: stats,
			VideoID: videoId,
			Video:   SimplifiedVideo{Url: video, Width: width, Height: height},
		}, nil

	}

	return SimplifiedData{
		Author:     author,
		Caption:    caption,
		Details:    stats,
		VideoID:    videoId,
		SoundLink:  getMediaLink(data, false),
		ImageLinks: slideLinks,
		Video:      SimplifiedVideo{},
	}, nil
}

//go:build !tikwm && !ttsave

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const Scraping = "false"

func (t *TiktokHTMLScript) getVideoDetails() Counts {
	return Counts{
		Likes:     FormatLargeNumbers(strconv.Itoa(t.Stats.DiggCount)),
		Comments:  FormatLargeNumbers(strconv.Itoa(t.Stats.CommentCount)),
		Shares:    FormatLargeNumbers(strconv.Itoa(t.Stats.ShareCount)),
		Views:     FormatLargeNumbers(strconv.Itoa(t.Stats.PlayCount)),
		Favorites: FormatLargeNumbers(t.Stats.CollectCount),
	}
}

func FetchTiktokData(videoId string) (SimplifiedData, error) {
	postAweme, err := fetch("https://tiktok.com/@placeholder/video/" + videoId)
	if err != nil {
		return SimplifiedData{}, err
	}
	videoUrl := postAweme.Video.DownloadAddr
	isVideo := !strings.Contains(videoUrl, "music")
	imageLinks := []string{}
	if !isVideo {
		imageLinks = postAweme.getImageLinks()
	}
	return SimplifiedData{
		Author: EscapeString(
			postAweme.Author.Nickname,
		) + " (@" + postAweme.Author.UniqueID + ")",
		Caption:    postAweme.Desc, // + "\n\n" + postAweme.Music.Title + " - " + postAweme.Music.Author + "🎵"
		VideoID:    videoId,
		Details:    postAweme.getVideoDetails(),
		ImageLinks: imageLinks,
		SoundLink:  postAweme.Music.PlayURL,
		Video: SimplifiedVideo{
			Url:    videoUrl,
			Width:  strconv.Itoa(postAweme.Video.Width),
			Height: strconv.Itoa(postAweme.Video.Height),
		},
	}, nil
}

func fetch(apiURL string) (TiktokHTMLScript, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return TiktokHTMLScript{}, err
	}

	req.Header.Set("accept-language", "fi-FI,fi;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("origin", "https://www.tiktok.com")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("referer", "https://www.tiktok.com/")
	req.Header.Set("sec-ch-ua", `"Chromium";v="122", "Not(A:Brand";v="24", "Google Chrome";v="122"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return TiktokHTMLScript{}, err
	}
	defer resp.Body.Close()

	textByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return TiktokHTMLScript{}, err
	}

	text := string(textByte)
	text = strings.Split(text, `"itemStruct":`)[1]
	text = strings.Split(text, `},"shareMeta":`)[0]
	fmt.Println(text)
	var responseScript TiktokHTMLScript
	err = json.Unmarshal([]byte(text), &responseScript)
	if err != nil {
		fmt.Println("Error unmarshalling script:", err)
		return TiktokHTMLScript{}, err
	}
	return responseScript, nil
}

func (t *TiktokHTMLScript) getImageLinks() []string {
	images := t.ImagePost.Images
	imageLinks := make([]string, 0, len(images))
	for _, image := range images {
		imageLinks = append(imageLinks, image.ImageURL.URLList[0])
	}
	return imageLinks
}

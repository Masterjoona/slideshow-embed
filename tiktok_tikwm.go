//go:build tikwm && !ttsave

package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const Scraping = "tikwm"
const baseUrl = "https://tikwm.com"

func fetchTikwm(videoId string) (TikWmResp, error) {
	client := &http.Client{}
	var data = strings.NewReader("url=https://www.tiktok.com/@placeholder/video/" + videoId)
	req, err := http.NewRequest("POST", baseUrl+"/api/", data)
	if err != nil {
		return TikWmResp{}, err
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err := client.Do(req)
	if err != nil {
		return TikWmResp{}, err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return TikWmResp{}, err
	}

	var tikWmResp TikWmResp
	err = json.Unmarshal(bodyText, &tikWmResp)
	if err != nil {
		return TikWmResp{}, err
	}

	return tikWmResp, nil
}

func FetchTiktokData(videoId string) (SimplifiedData, error) {
	tikWmApiResp, err := fetchTikwm(videoId)
	data := tikWmApiResp.Data
	if err != nil {
		return SimplifiedData{}, err
	}
	stats := Counts{
		Likes:     FormatLargeNumbers(strconv.Itoa(data.DiggCount)),
		Comments:  FormatLargeNumbers(strconv.Itoa(data.CommentCount)),
		Shares:    FormatLargeNumbers(strconv.Itoa(data.ShareCount)),
		Views:     FormatLargeNumbers(strconv.Itoa(data.PlayCount)),
		Favorites: FormatLargeNumbers(strconv.Itoa(data.CollectCount)),
	}
	author := data.Author.Nickname + " (@" + data.Author.UniqueID + ")"
	caption := data.Title
	if data.Duration != 0 {
		videoUrl := data.Play
		width, height, err := GetVideoDimensionsFromUrl(videoUrl)
		if err != nil {
			return SimplifiedData{}, err
		}

		return SimplifiedData{
			Author:  author,
			Caption: caption,
			Details: stats,
			Video: SimplifiedVideo{
				Url:    videoUrl,
				Width:  width,
				Height: height,
			},
		}, nil
	}
	return SimplifiedData{
		Author:     author,
		Caption:    caption,
		Details:    stats,
		SoundLink:  data.MusicInfo.Play,
		ImageLinks: data.Images,
		Video:      SimplifiedVideo{},
	}, nil
}

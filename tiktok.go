//go:build !scrape

package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const Scraping = false

func getVideoDetails(aweme Aweme) Counts {
	return Counts{
		Likes:     FormatLargeNumbers(strconv.Itoa(aweme.Statistics.DiggCount)),
		Comments:  FormatLargeNumbers(strconv.Itoa(aweme.Statistics.CommentCount)),
		Shares:    FormatLargeNumbers(strconv.Itoa(aweme.Statistics.ShareCount)),
		Views:     FormatLargeNumbers(strconv.Itoa(aweme.Statistics.PlayCount)),
		Favorites: FormatLargeNumbers(strconv.Itoa(aweme.Statistics.CollectCount)),
	}
}

func FetchTiktokData(videoId string) (SimplifiedData, error) {
	queryString := buildQueryUrl(videoId)
	apiResponse, err := fetch(queryString)
	if err != nil {
		return SimplifiedData{}, err
	}
	postAweme := apiResponse.AwemeList[0]

	isVideo := !strings.Contains(postAweme.Video.PlayAddr.URLList[0], "music")
	imageLinks := []string{}
	if !isVideo {
		imageLinks = getImageLinks(postAweme)
	}
	return SimplifiedData{
		Author: EscapeString(
			postAweme.Author.Nickname,
		) + " (@" + postAweme.Author.UniqueID + ")",
		Caption:    postAweme.Desc, // + "\n\n" + postAweme.Music.Title + " - " + postAweme.Music.Author + "🎵"
		VideoID:    videoId,
		Details:    getVideoDetails(postAweme),
		ImageLinks: imageLinks,
		SoundLink:  postAweme.Music.PlayURL.URI,
		IsVideo:    isVideo,
		Video: SimplifiedVideo{
			Url:    postAweme.Video.PlayAddr.URLList[0],
			Width:  strconv.Itoa(postAweme.Video.Width),
			Height: strconv.Itoa(postAweme.Video.Height),
		},
	}, nil
}

func fetch(apiURL string) (TikTokAPIResponse, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return TikTokAPIResponse{}, err
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
		return TikTokAPIResponse{}, err
	}
	defer resp.Body.Close()

	var responseStruct TikTokAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&responseStruct)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		println(apiURL)
		return TikTokAPIResponse{}, err
	}

	return responseStruct, nil
}

func getImageLinks(aweme Aweme) []string {
	imageLinks := make([]string, 0, len(aweme.ImagePostInfo.Images))
	for _, image := range aweme.ImagePostInfo.Images {
		var url string
		if aweme.ImagePostInfo.Images[0].BitrateImages != nil {
			url = image.BitrateImages[0].BitrateImage.URLList[1]
		} else {
			url = image.Thumbnail.URLList[1]
		}
		imageLinks = append(imageLinks, url)
	}

	return imageLinks
}

func setParam(key, value string) string {
	return key + "=" + value + "&"
}

func getNextId(ids []string) string {
	return ids[rand.Intn(len(ids))]
}

func buildQueryUrl(videoId string) string {
	query := "https://api22-normal-c-alisg.tiktokv.com/aweme/v1/feed/?"
	query += setParam("iid", getNextId(InstallIds))
	query += setParam(
		"device_id",
		getNextId(DeviceIds),
		//strconv.FormatInt(randomBigInt(7250000000000000000, 7351147085025500000), 10),
	)

	// https://github.com/Evil0ctal/Douyin_TikTok_Download_API
	query += setParam("channel", "googleplay")
	query += setParam("app_name", "musical_ly")
	query += setParam("version_code", "300904")
	query += setParam("device_platform", "android")
	query += setParam("device_type", "SM-ASUS_Z01QD")
	query += setParam("os_version", "9")

	return query + setParam("aweme_id", videoId)
}

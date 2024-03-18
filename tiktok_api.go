package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// https://github.com/Britmoji/tiktxk/blob/main/src/util/tiktok.ts

const (
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"
)

type Counts struct {
	Likes     string
	Comments  string
	Shares    string
	Views     string
	Favorites string
	Downloads string
}

func GetVideoDetails(aweme Aweme) Counts {
	return Counts{
		Likes:     FormatLargeNumbers(strconv.Itoa(aweme.Statistics.DiggCount)),
		Comments:  FormatLargeNumbers(strconv.Itoa(aweme.Statistics.CommentCount)),
		Shares:    FormatLargeNumbers(strconv.Itoa(aweme.Statistics.ShareCount)),
		Views:     FormatLargeNumbers(strconv.Itoa(aweme.Statistics.PlayCount)),
		Favorites: FormatLargeNumbers(strconv.Itoa(aweme.Statistics.CollectCount)),
		Downloads: FormatLargeNumbers(strconv.Itoa(aweme.Statistics.DownloadCount)),
	}
}

func GetLongVideoId(videoUrl string) (string, error) {
	if !validateURL(videoUrl) {
		return "", errors.New("invalid URL")
	}
	if strings.Contains(videoUrl, "/photo/") {
		return strings.Split(videoUrl, "/photo/")[1], nil
	}

	if strings.Contains(videoUrl, "/video/") {
		return strings.Split(videoUrl, "/video/")[1], nil
	}

	videoUrl = strings.ReplaceAll(videoUrl, "tiktxk.com", "tiktok.com")
	resp, err := http.Head(videoUrl)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("failed to fetch the slideshow")
	}

	queryless := strings.Split(resp.Request.URL.String(), "?")[0]
	return strings.Split(queryless, "/")[5], nil

}

func FetchTiktokData(videoId string) (SimplifiedData, error) {
	queryString := "https://api22-normal-c-useast2a.tiktokv.com/aweme/v1/feed/?aid=1180&version_name=26.1.3&version_code=260103&aweme_id=" + videoId + "&build_number=26.1.3&manifest_version_code=260103&update_version_code=260103&opeudid=c07d3f637cde7535&uuid=6973698106620498&_rticket=1710507415580&ts=1710507415&device_brand=Google&device_type=Pixel+4&device_platform=android&resolution=1080*1920&dpi=420&os_version=10&os_api=29&carrier_region=US&sys_region=US&region=US&app_name=trill&app_language=en&language=en&timezone_name=America%2FNew_York&timezone_offset=-14400&channel=googleplay&ac=wifi&mcc_mnc=310260&is_my_cn=0&ssmix=a&as=a1qwert123&cp=cbfhckdckkde1"
	apiResponse, err := fetch(queryString)
	if err != nil {
		return SimplifiedData{}, err
	}
	postAweme := apiResponse.AwemeList[0]

	isVideo := !strings.Contains(postAweme.Video.PlayAddr.URLList[0], "music")
	imageLinks := []string{}
	if !isVideo {
		imageLinks = GetImageLinks(postAweme)
	}
	return SimplifiedData{
		Author: EscapeString(
			postAweme.Author.Nickname,
		) + " (@" + postAweme.Author.UniqueID + ")",
		Caption:    postAweme.Desc + "\n\n" + postAweme.Music.Title + " - " + postAweme.Music.Author + "🎵",
		VideoID:    videoId,
		Details:    GetVideoDetails(postAweme),
		ImageLinks: imageLinks,
		SoundUrl:   postAweme.Music.PlayURL.URI,
		IsVideo:    isVideo,
		Video:      postAweme.Video,
	}, nil
}

func fetch(apiURL string) (TikTokAPIResponse, error) {
	client := &http.Client{
		Timeout: time.Second * 4,
	}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return TikTokAPIResponse{}, err
	}

	req.Header.Set("authority", "api22-normal-c-useast2a.tiktokv.com")
	req.Header.Set("accept", "*/*")
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
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("User-Agent", UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return TikTokAPIResponse{}, err
	}
	defer resp.Body.Close()

	if resp.ContentLength == 0 {
		fmt.Println("Response is empty")
		return TikTokAPIResponse{}, nil
	}

	var responseStruct TikTokAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&responseStruct)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return TikTokAPIResponse{}, err
	}

	return responseStruct, nil
}

func GetImageLinks(aweme Aweme) []string {
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

func DownloadImage(url, outputPath string) error {
	url = EscapeString(url)
	client := &http.Client{
		Timeout: time.Second * 4,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)

	return err

}

func DownloadImages(links []string, outputDir string) error {
	var wg sync.WaitGroup
	CreateDirectory(outputDir)
	for index, link := range links {
		wg.Add(1)
		go func(i int, url string) {
			defer wg.Done()
			if err := DownloadImage(url, fmt.Sprintf("%s/%d.jpg", outputDir, i+1)); err != nil {
				log.Printf("error downloading image %s: %v\n", url, err)
			}
		}(index, link)
	}
	wg.Wait()
	return nil
}

func DownloadAudio(link string, outputDir string) error {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("range", "bytes=0-")
	req.Header.Set("referer", "https://www.tiktok.com/")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		return errors.New("failed to fetch the audio")
	}

	out, err := os.Create(fmt.Sprintf("%s/%s", outputDir, "audio.mp3"))
	if err != nil {
		fmt.Println("Error creating the file:", err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Error writing the file:", err)
		return err
	}
	return nil
}

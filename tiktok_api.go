package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// https://github.com/Britmoji/tiktxk/blob/main/src/util/tiktok.ts

const (
	UserAgent      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"
	MaxRetry       = 3
	RetryDelaySecs = 2
)

func PostDetails(videoId string) (TikTokAPIResponse, error) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	queryString := "https://api22-normal-c-useast2a.tiktokv.com/aweme/v1/feed/?aid=1180&version_name=26.1.3&version_code=260103&aweme_id=" + videoId + "&build_number=26.1.3&manifest_version_code=260103&update_version_code=260103&opeudid=c07d3f637cde7535&uuid=6973698106620498&_rticket=1710507415580&ts=1710507415&device_brand=Google&device_type=Pixel+4&device_platform=android&resolution=1080*1920&dpi=420&os_version=10&os_api=29&carrier_region=US&sys_region=US&region=US&app_name=trill&app_language=en&language=en&timezone_name=America%2FNew_York&timezone_offset=-14400&channel=googleplay&ac=wifi&mcc_mnc=310260&is_my_cn=0&ssmix=a&as=a1qwert123&cp=cbfhckdckkde1"
	respStruct, err := fetch(queryString)
	if err != nil {
		return TikTokAPIResponse{}, err
	}
	respStruct.AwemeList = respStruct.AwemeList[:1]
	return respStruct, nil
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
	req.Header.Set(
		"User-Agent",
		UserAgent,
	)

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

func DownloadImage(url, outputPath string) error {
	url = EscapeString(url)
	client := &http.Client{
		Timeout: time.Second * 4,
	}

	var err error
	for retry := 0; retry < MaxRetry; retry++ {
		var resp *http.Response
		resp, err = client.Get(url)
		if err != nil {
			time.Sleep(RetryDelaySecs * time.Second)
			continue
		}
		defer resp.Body.Close()

		if !strings.HasPrefix(resp.Header.Get("Content-Type"), "image/") {
			println("not an image")
			time.Sleep(RetryDelaySecs * time.Second)
			continue
		}

		out, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("error creating file: %v", err)
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return fmt.Errorf("error copying file: %v", err)
		}

		return nil
	}

	if err != nil {
		return fmt.Errorf("maximum retry count reached: %v", err)
	}

	return fmt.Errorf("maximum retry count reached")
}

type Counts struct {
	Likes     string
	Comments  string
	Shares    string
	Views     string
	Favorites string
	Downloads string
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
	aweme, err := http.Head(videoUrl)
	if err != nil {
		return "", err
	}
	defer aweme.Body.Close()
	if aweme.StatusCode != 200 {
		return "", errors.New("failed to fetch the slideshow")
	}
	queryless := strings.Split(aweme.Request.URL.String(), "?")[0]
	return strings.Split(queryless, "/")[5], nil

}

func GetImageLinks(aweme Aweme) []string {
	imageLinks := []string{}
	for _, image := range aweme.ImagePostInfo.Images {
		imageLinks = append(imageLinks, image.BitrateImages[0].BitrateImage.URLList[1])
	}
	return imageLinks
}

func GetAuthor(aweme Aweme) string {
	return EscapeString(aweme.Author.Nickname) + " (@" + EscapeString(aweme.Author.UniqueID) + ")"
}

func GetVideoDetails(aweme Aweme) Counts {
	return Counts{
		Likes:     formatNumber(strconv.Itoa(aweme.Statistics.DiggCount)),
		Comments:  formatNumber(strconv.Itoa(aweme.Statistics.CommentCount)),
		Shares:    formatNumber(strconv.Itoa(aweme.Statistics.ShareCount)),
		Views:     formatNumber(strconv.Itoa(aweme.Statistics.PlayCount)),
		Favorites: formatNumber(strconv.Itoa(aweme.Statistics.CollectCount)),
		Downloads: formatNumber(strconv.Itoa(aweme.Statistics.DownloadCount)),
	}
}

func formatNumber(numberString string) string {
	const (
		million  = 1000000
		thousand = 1000
	)
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

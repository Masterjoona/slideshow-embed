package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

/*
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

	func FetchTiktokData(videoId string) (SimplifiedData, error) {
		queryString := "https://api22-normal-c-useast2a.tiktokv.com/aweme/v1/feed/?aid=1180&version_name=26.1.3&version_code=260103&aweme_id=" + videoId + "&build_number=26.1.3&manifest_version_code=260103&update_version_code=260103&opeudid=c07d3f637cde7535&uuid=6973698106620498&_rticket=1710507415580&ts=1710507415&device_brand=Google&device_type=Pixel+4&device_platform=android&resolution=1080*1920&dpi=420&os_version=10&os_api=29&carrier_region=US&sys_region=US&region=US&app_name=trill&app_language=en&language=en&timezone_name=America%2FNew_York&timezone_offset=-14400&channel=googleplay&ac=wifi&mcc_mnc=310260&is_my_cn=0&ssmix=a&as=a1qwert123&cp=cbfhckdckkde1"
		apiResponse, err := fetchUntilData(queryString)
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

	func fetchUntilData(url string) (TikTokAPIResponse, error) {
		// Yes, it works.
		iterations := 0
		for {
			if iterations > 150 { // one req took me 250 reqs lol
				println("Took too many iterations to fetch the data")
				return TikTokAPIResponse{}, errors.New("failed to fetch the data")
			}
			apiResponse, length, err := fetch(url)
			if err != nil {
				return TikTokAPIResponse{}, err
			}
			if length == 0 {
				iterations++
				continue
			}
			println("Took", iterations, "iterations to fetch the data")
			return apiResponse, nil
		}
	}

	func fetch(apiURL string) (TikTokAPIResponse, int, error) {
		client := &http.Client{
			Timeout: time.Second * 4,
		}
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return TikTokAPIResponse{}, 0, err
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
			return TikTokAPIResponse{}, -1, err
		}
		defer resp.Body.Close()

		if resp.ContentLength == 0 {
			return TikTokAPIResponse{}, 0, nil
		}

		var responseStruct TikTokAPIResponse
		err = json.NewDecoder(resp.Body).Decode(&responseStruct)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return TikTokAPIResponse{}, int(resp.ContentLength), err
		}

		return responseStruct, int(resp.ContentLength), nil
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
*/

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
		return "", errors.New("failed to fetch the tiktok")
	}

	queryless := strings.Split(resp.Request.URL.String(), "?")[0]
	return strings.Split(queryless, "/")[5], nil

}

func DownloadImage(url string) ([]byte, error) {
	url = EscapeString(url)
	client := &http.Client{
		Timeout: time.Second * 4,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return imageBytes, nil
}

func DownloadImages(links []string) (*[][]byte, error) {
	var wg sync.WaitGroup
	var imagesInted []ImageWithIndex

	for i, link := range links {
		wg.Add(1)
		go func(url string, index int) {
			defer wg.Done()
			if imgBytes, err := DownloadImage(url); err == nil {
				imagesInted = append(imagesInted, ImageWithIndex{Bytes: imgBytes, Index: index})
			} else {
				log.Printf("error downloading image %s: %v\n", url, err)
			}
		}(link, i)
	}
	wg.Wait()

	sort.Slice(imagesInted, func(i, j int) bool {
		return imagesInted[i].Index < imagesInted[j].Index
	})

	images := make([][]byte, 0, len(imagesInted))
	for _, img := range imagesInted {
		images = append(images, img.Bytes)
	}
	return &images, nil
}

func DownloadAudio(link string) (*[]byte, error) {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("range", "bytes=0-")
	req.Header.Set("referer", "https://www.tiktok.com/")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		return nil, errors.New("failed to fetch the audio")
	}

	audioBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &audioBytes, nil
}

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

func downloadImage(url string) ([]byte, error) {
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

	if !strings.Contains(resp.Header.Get("Content-Type"), "image") {
		/*
			thats weird?
			https://vm.tiktok.com/ZGeH3Covr/
			{ "code": 4404,"error": "fail to get resource"}
		*/
		return nil, errors.New("tiktok returned a non-image response")
	}

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return imageBytes, nil
}

func DownloadImages(links []string) (*[][]byte, int) {
	var wg sync.WaitGroup
	var imagesInted []ImageWithIndex
	failedCount := 0

	for i, link := range links {
		wg.Add(1)
		go func(url string, index int) {
			defer wg.Done()
			if imgBytes, err := downloadImage(url); err == nil {
				imagesInted = append(imagesInted, ImageWithIndex{Bytes: imgBytes, Index: index})
			} else {
				log.Printf("error downloading image %s: %v\n", url, err)
				failedCount += 1
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
	return &images, failedCount
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

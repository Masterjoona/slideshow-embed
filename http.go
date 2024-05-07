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

	matches := longLinkRe.FindStringSubmatch(videoUrl)
	if len(matches) > 1 {
		return matches[1], nil
	}

	matches = shortLinkRe.FindStringSubmatch(videoUrl)
	if len(matches) > 1 {
		resp, err := http.Head("https://vm.tiktok.com/" + matches[1])
		if err != nil {
			return "", err
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return "", errors.New("failed to fetch the tiktok")
		}

		return longLinkRe.FindStringSubmatch(resp.Request.URL.String())[1], nil
	}

	return "", errors.New("failed to extract the video id")
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

func (t *SimplifiedData) DownloadImages() int {
	var wg sync.WaitGroup
	var indexedImages []ImageWithIndex

	var failedCount int = 0 // Max images in a tiktok is 35, so we can expect fails to be <=35

	for i, link := range t.ImageLinks {
		wg.Add(1)
		go func(url string, index int) {
			defer wg.Done()
			if imgBytes, err := downloadImage(url); err == nil {
				indexedImages = append(indexedImages, ImageWithIndex{Bytes: imgBytes, Index: index})
			} else {
				log.Printf("error downloading image %s: %v\n", url, err)
				failedCount += 1
			}
		}(link, i)
	}
	wg.Wait()

	sort.Slice(indexedImages, func(i, j int) bool {
		return indexedImages[i].Index < indexedImages[j].Index
	})

	t.ImageBuffers = make([][]byte, 0, len(indexedImages))
	for _, img := range indexedImages {
        t.ImageBuffers = append(t.ImageBuffers, img.Bytes)
    }
	return failedCount
}

func (t *SimplifiedData) DownloadSound() error {
	req, err := http.NewRequest("GET", t.SoundLink, nil)
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

	t.SoundBuffer, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

package net

import (
	"errors"
	"io"
	"meow/pkg/vars"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var ShortURLCache = NewCache[string, ShortLinkInfo]()
var longLinkRe = regexp.MustCompile(`https:\/\/(?:www.)?(?:vxtiktok|tiktok|tiktxk|)\.com\/(@.{2,32})\/(?:photo|video)\/(\d+)`)
var shortLinkRe = regexp.MustCompile(`https:\/\/.{1,3}\.(?:(?:vx|)tikt(?:x|o)k)\.com/(?:.{1,2}/|)(.{5,12})\/`)

func GetLongVideoId(videoUrl string) (string, error) {
	if !validateURL(videoUrl) {
		if _, err := strconv.Atoi(videoUrl); err == nil {
			// i guess we can trust the user that this is a valid video id
			return videoUrl, nil
		}
		return "", errors.New("invalid URL")
	}

	matches := longLinkRe.FindStringSubmatch(videoUrl)
	if len(matches) > 2 {
		return matches[2], nil
	}

	matches = shortLinkRe.FindStringSubmatch(videoUrl)
	if len(matches) > 1 {
		if cachedInfo, ok := ShortURLCache.Get(matches[1]); ok {
			return cachedInfo.VideoId, nil
		}
		resp, err := http.Head("https://vm.tiktok.com/" + matches[1])
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		finalUrl := resp.Request.URL.String()
		matches = longLinkRe.FindStringSubmatch(finalUrl)
		if len(matches) > 2 {
			ShortURLCache.Set(matches[2], ShortLinkInfo{VideoId: matches[2], UniqueUserId: matches[1]})
			return matches[2], nil
		}
	}

	return "", errors.New("failed to extract the video id")
}

func DownloadMedia(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", vars.UserAgent)

	resp, err := vars.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !strings.Contains(resp.Header.Get("Content-Type"), "image") && !strings.Contains(resp.Header.Get("Content-Type"), "video") {
		return nil, errors.New("tiktok returned a non-media response")
	}

	mediaBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return mediaBytes, nil
}

func DownloadImages(videoId string, imgLinks []string) [][]byte {
	var wg sync.WaitGroup
	var indexedImages []ImageWithIndex

	// var failedCount int

	for i, link := range imgLinks {
		wg.Add(1)
		go func(url string, index int) {
			defer wg.Done()
			if imgBytes, err := DownloadMedia(url); err == nil {
				indexedImages = append(indexedImages, ImageWithIndex{Bytes: imgBytes, Index: index})
			} else {
				println("error downloading image on: %v\n", videoId, err)
				// failedCount += 1
			}
		}(link, i)
	}
	wg.Wait()

	sort.Slice(indexedImages, func(i, j int) bool {
		return indexedImages[i].Index < indexedImages[j].Index
	})

	imageBuffers := make([][]byte, 0, len(indexedImages))

	for _, img := range indexedImages {
		imageBuffers = append(imageBuffers, img.Bytes)
	}

	return imageBuffers //, failedCount
}

func DownloadSound(soundLink string) ([]byte, error) {
	req, err := http.NewRequest("GET", soundLink, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", vars.UserAgent)
	req.Header.Set("range", "bytes=0-")
	req.Header.Set("referer", "https://www.tiktok.com/")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		return nil, errors.New("failed to fetch the audio")
	}

	soundBuffer, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return soundBuffer, nil
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

func GetLongVideoId(videoUrl string) (string, string, error) {
	if !validateURL(videoUrl) {
		return "", "", errors.New("invalid URL")
	}

	matches := longLinkRe.FindStringSubmatch(videoUrl)
	if len(matches) > 1 {
		return matches[1], matches[2], nil
	}

	matches = shortLinkRe.FindStringSubmatch(videoUrl)
	if len(matches) > 1 {
		if cachedInfo, ok := ShortURLCache.Get(matches[1]); ok {
			return cachedInfo.VideoId, cachedInfo.UniqueUserId, nil
		}
		resp, err := http.Head("https://vm.tiktok.com/" + matches[1])
		if err != nil {
			return "", "", err
		}

		defer resp.Body.Close()

		finalUrl := resp.Request.URL.String()
		if strings.Contains(finalUrl, "@") {
			matches = longLinkRe.FindStringSubmatch(finalUrl)
			if len(matches) > 1 {
				ShortURLCache.Put(matches[2], ShortLinkInfo{VideoId: matches[2], UniqueUserId: matches[1]})
				return matches[1], matches[2], nil
			}
		}
	}

	return "", "", errors.New("failed to extract the video id")
}

func downloadMedia(url string) ([]byte, error) {
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

	if !strings.Contains(resp.Header.Get("Content-Type"), "image") && !strings.Contains(resp.Header.Get("Content-Type"), "video") {
		/*
			thats weird?
			https://www.tiktok.com/@f3l1xfromvenus/photo/7360247887340637472
			some image had
			{ "code": 4404,"error": "fail to get resource"}
		*/
		return nil, errors.New("tiktok returned a non-media response")
	}

	mediaBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return mediaBytes, nil
}

func (t *SimplifiedData) DownloadImages() int {
	var wg sync.WaitGroup
	var indexedImages []ImageWithIndex

	var failedCount int

	for i, link := range t.ImageLinks {
		wg.Add(1)
		go func(url string, index int) {
			defer wg.Done()
			if imgBytes, err := downloadMedia(url); err == nil {
				indexedImages = append(indexedImages, ImageWithIndex{Bytes: imgBytes, Index: index})
			} else {
				println("error downloading image on: %v\n", t.VideoID, err)
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

func (t *SimplifiedData) DownloadVideoAndSubtitles(lang string) error {
	url := SubtitlesHost + "subtitle_id=02981317794434464&target_language=" + lang + "&item_id=" + t.VideoID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var subtitlesResp SubtitlesResp
	err = json.Unmarshal(bodyText, &subtitlesResp)
	if err != nil {
		return err
	}
	subtitles := subtitlesResp.WebvttSubtitle
	if subtitles == "" {
		return errors.New("no subtitles found")
	}

	err = CreateDirectory(TemporaryDirectory + "/collages/" + t.VideoID)
	if err != nil {
		return err
	}
	err = os.WriteFile(TemporaryDirectory+"/collages/"+t.VideoID+"/subtitles.vtt", []byte(subtitles), 0644)
	if err != nil {
		return err
	}

	videoUrl := t.Video.Url
	videoBytes, err := downloadMedia(videoUrl)
	if err != nil {
		return err
	}
	return os.WriteFile(TemporaryDirectory+"/collages/"+t.VideoID+"/video.mp4", videoBytes, 0644)
}

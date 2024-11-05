package net

import (
	"encoding/json"
	"errors"
	"io"
	"meow/pkg/config"
	"meow/pkg/files"
	"meow/pkg/vars"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
)

var ShortURLCache = NewCache[ShortLinkInfo](20)
var longLinkRe = regexp.MustCompile(`https:\/\/(?:www.)?(?:vxtiktok|tiktok|tiktxk|)\.com\/(@.{2,32})\/(?:photo|video)\/(\d+)`)
var shortLinkRe = regexp.MustCompile(`https:\/\/.{1,3}\.(?:(?:vx|)tikt(?:x|o)k)\.com/(?:.{1,2}/|)(.{5,12})\/`)

var tmpDir = config.TemporaryDirectory + "/collages/"

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

func DownloadVideoAndSubtitles(videoId, videoUrl, lang string) error {
	url := vars.SubtitlesHost + "subtitle_id=02981317794434464&target_language=" + lang + "&item_id=" + videoId
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", vars.UserAgent)

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

	videoTmpDir := tmpDir + videoId

	err = files.CreateDirectory(videoTmpDir)
	if err != nil {
		return err
	}
	err = os.WriteFile(videoTmpDir+"/subtitles.vtt", []byte(subtitles), 0644)
	if err != nil {
		return err
	}

	videoBytes, err := DownloadMedia(videoUrl)
	if err != nil {
		return err
	}
	return os.WriteFile(videoTmpDir+"/video.mp4", videoBytes, 0644)
}

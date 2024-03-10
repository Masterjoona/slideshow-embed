package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type Counts struct {
	Likes    string
	Comments string
	Shares   string
	Views    string
}

var proxitokInstances = []string{
	"https://proxitok.pabloferreiro.es/",
	"https://proxitok.belloworld.it/",
	"https://proxitok.privacydev.net/",
	// "https://tok.artemislena.eu/",
	"https://tok.adminforge.de/",
	"https://tik.hostux.net/",
	"https://tt.vern.cc/",
	"https://tik.hostux.net/",
	"https://proxitok.lunar.icu/",
	"https://cringe.datura.network/",
	"https://tt.opnxng.com/",
	"https://tiktok.wpme.pl/",
	"https://proxitok.r4fo.com/",
}

var instanceIndex = 0

func getNextInstanceURL() string {
	if instanceIndex == len(proxitokInstances)-1 {
		instanceIndex = 0
	} else {
		instanceIndex++
	}
	return proxitokInstances[instanceIndex]
}

func getInstanceURL() string {
	if ProxiTokInstance != "" && ProxiTokInstance != "/" {
		return ProxiTokInstance
	}
	return getNextInstanceURL()
}

func GetLongVideoId(videoUrl string) (string, error) {
	if !validateURL(videoUrl) {
		return "", errors.New("invalid URL")
	}
	if strings.Contains(videoUrl, "/photo/") {
		return strings.Split(videoUrl, "/photo/")[1], nil
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

func FetchProxiTokVideo(videoID string) (string, error) {
	instanceUrl := getInstanceURL() + "@placeholder/video/" + videoID
	println("fetching from " + instanceUrl)
	resp, err := http.Get(instanceUrl)
	if err != nil {
		println("error while fetching the slideshow")
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		println("wrong status code")
		return "", errors.New("failed to fetch the slideshow")
	}
	respContent, err := io.ReadAll(resp.Body)
	if err != nil {
		println("error while reading the response")
		return "", err
	}
	content := string(respContent)
	if strings.Contains(content, "author_secret") {
		return "private", nil
	}
	if strings.Contains(content, "There was an error processing your request!") {
		return "", errors.New("error processing the request")
	}
	return content, nil
}

func GetImageLinks(content string) ([]string, error) {
	imageLinkRegex := regexp.MustCompile(`<img src="\/stream\?url=(.*)">`)
	matches := imageLinkRegex.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return nil, errors.New("no images found")
	}

	imageLinks := make([]string, len(matches))
	for i, match := range matches {
		imageLinks[i] = EscapeString(match[1])
	}
	return imageLinks[:len(imageLinks)-1], nil
}

func GetAudioLink(input string) (string, error) {
	soundRegex := regexp.MustCompile(`preload="none" src="\/stream\?url=(.*)">`)
	soundMatch := soundRegex.FindStringSubmatch(input)
	if len(soundMatch) == 0 {
		return "", fmt.Errorf("no audio link found in input")
	}
	return EscapeString(soundMatch[1]), nil
}

func GetAuthorAndCaption(content string) (string, string, error) {
	nicknameRegex := regexp.MustCompile(`<strong>(.*)</`)
	authorNameRegex := regexp.MustCompile(`">(@.*)</a`)
	captionRegex := regexp.MustCompile(`twitter:description" content="(.*)">`)

	nicknameMatch := nicknameRegex.FindStringSubmatch(content)
	authorNameMatch := authorNameRegex.FindStringSubmatch(content)
	captionMatch := captionRegex.FindStringSubmatch(content)

	if len(nicknameMatch) == 0 || len(authorNameMatch) == 0 || len(captionMatch) == 0 {
		return "", "", errors.New("no author or caption found")
	}

	stylized := EscapeString(nicknameMatch[1]) + " (" + EscapeString(authorNameMatch[1]) + ")"
	return stylized, EscapeString(captionMatch[1]), nil
}

func GetVideoDetails(content string) Counts {
	detailsRegex := regexp.MustCompile(
		`<span class="icon">\s*<i class="gg-(?:eye|heart|comment|share)"><\/i>\s*<\/span>\s*<span>(.*?)<\/span>`,
	)
	detailsMatches := detailsRegex.FindAllStringSubmatch(content, -1)
	if len(detailsMatches) < 2 {
		println("incorrect details matches")
		return Counts{}
	}
	return Counts{
		Likes:    detailsMatches[1][1],
		Comments: detailsMatches[2][1],
		Shares:   detailsMatches[3][1],
		Views:    detailsMatches[0][1],
	}
}

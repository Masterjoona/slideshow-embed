package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ExtractImageLinks(input string) ([]string, error) {
	imagePostRegex := regexp.MustCompile(`imagePost":(.*?),"title"`)
	directUrlRegex := regexp.MustCompile(`"urlList":\["(.*?)"`)

	imagePostMatch := imagePostRegex.FindStringSubmatch(input)
	if len(imagePostMatch) == 0 {
		return nil, fmt.Errorf("no image links found in input")
	}

	escapedUrls := strings.ReplaceAll(imagePostMatch[1], "\\u002F", "/")
	imageLinksMatches := directUrlRegex.FindAllStringSubmatch(escapedUrls, -1)
	if len(imageLinksMatches) == 0 {
		fmt.Println(escapedUrls)
		return nil, fmt.Errorf("no image links found in input")
	}

	imageLinks := make([]string, len(imageLinksMatches))
	for index, link := range imageLinksMatches {
		imageLinks[index] = link[1]
	}

	return imageLinks, nil
}

func ExtractAudioLink(input string) (string, error) {
	soundRegex := regexp.MustCompile(`playUrl":"(.*?)"`)
	soundMatch := soundRegex.FindStringSubmatch(input)
	if len(soundMatch) == 0 {
		return "", fmt.Errorf("no audio link found in input")
	}
	escapedUrl := strings.ReplaceAll(soundMatch[1], "\\u002F", "/")
	return escapedUrl, nil
}

func ExtractVideoID(url string) (string, error) {
	parts := strings.Split(url, "/")
	if len(parts) < 4 {
		return "", fmt.Errorf("invalid URL: %s", url)
	}
	return parts[3], nil
}

func GetVideoAuthorAndCaption(url string, videoID string) (string, string, string, error) {
	responseBody, err := FetchResponseBody(url)

	if err != nil {
		return "", "", "", err
	}

	authorAndNickRegex := regexp.MustCompile(`uniqueId":"(.*?)","nickname":"(.*?)","avatar`)
	captionRegex := regexp.MustCompile(`"id": "\d+",\s*"desc": ".*?",`)
	possibleTitleRegex := regexp.MustCompile(`},"title":"(.*?)"},"locationCreated":`)

	authorNick := authorAndNickRegex.FindStringSubmatch(responseBody)
	if len(authorNick) == 0 {
		return "", "", "", fmt.Errorf("no author name found in response")
	}

	author := authorNick[2] + " (@" + authorNick[1] + ")"
	caption := captionRegex.FindStringSubmatch(responseBody)
	possibleTitle := possibleTitleRegex.FindStringSubmatch(responseBody)

	captionText := ""
	possibleTitleText := ""

	if len(caption) != 0 {
		captionText = caption[1]
	} else {
		captionText = ""
	}

	if len(possibleTitle) != 0 {
		possibleTitleText = possibleTitle[1]
	} else {
		possibleTitleText = ""
	}

	captionText = possibleTitleText + " " + captionText
	return author, captionText, responseBody, nil
}

type Counts struct {
	Likes     string
	Comments  string
	Favorited string
	Shares    string
	Views     string
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

func GetVideoDetails(responseBody string) Counts {
	detailsRegex := regexp.MustCompile(`{"diggCount":(\d+),"shareCount":(\d+),"commentCount":(\d+),"playCount":(\d+),"collectCount":"(\d+)"}`)
	detailsMatch := detailsRegex.FindStringSubmatch(responseBody)
	if len(detailsMatch) == 0 {
		return Counts{
			Likes:     "0",
			Comments:  "0",
			Favorited: "0",
			Shares:    "0",
			Views:     "0",
		}
	}

	return Counts{
		Likes:     formatNumber(detailsMatch[1]),
		Comments:  formatNumber(detailsMatch[3]),
		Favorited: formatNumber(detailsMatch[5]),
		Shares:    formatNumber(detailsMatch[2]),
		Views:     formatNumber(detailsMatch[4]),
	}
}

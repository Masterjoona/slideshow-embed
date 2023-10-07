package extracting

import (
	"fmt"
	"meow/httputil"
	"regexp"
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

func ExtractVideoID(url string) (string, error) {
	parts := strings.Split(url, "/")
	if len(parts) < 4 {
		return "", fmt.Errorf("invalid URL: %s", url)
	}
	return parts[3], nil
}

func GetVideoAuthorAndCaption(url string, videoID string) (string, string, string, error) {
	responseBody, err := httputil.FetchResponseBody(url)

	if err != nil {
		return "", "", "", err
	}

	authorNameRegex := regexp.MustCompile(`"nickname":"(.*?)"`)
	captionRegex := regexp.MustCompile(`"contents":\[{"desc":"(.*?)",`)
	possibleTitleRegex := regexp.MustCompile(`},"title":"(.*?)"},"locationCreated":`)

	authorName := authorNameRegex.FindStringSubmatch(responseBody)[1]
	caption := captionRegex.FindStringSubmatch(responseBody)
	possibleTitle := possibleTitleRegex.FindStringSubmatch(responseBody)

	captionText := ""
	possibleTitleText := ""

	if len(caption) != 0 {
		captionText = caption[1]
	} else {
		captionText = "No caption"
	}

	if len(possibleTitle) != 0 {
		possibleTitleText = possibleTitle[1]
	} else {
		possibleTitleText = ""
	}
	captionText = possibleTitleText + " " + captionText
	return authorName, captionText, responseBody, nil
}

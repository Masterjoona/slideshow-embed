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
	responseBody, err := httputil.FetchResponseBody(url)

	if err != nil {
		return "", "", "", err
	}

	authorNameRegex := regexp.MustCompile(`{"@type":"Thing","@id":"https:\/\/www\.tiktok\.com\/@(?:.*?)","name":"(.*?) \| TikTok"}},{"@type":`)
	captionRegex := regexp.MustCompile(`"contents":\[{"desc":"(.*?)",`)
	possibleTitleRegex := regexp.MustCompile(`},"title":"(.*?)"},"locationCreated":`)

	authorName := authorNameRegex.FindStringSubmatch(responseBody)
	if len(authorName) == 0 {
		return "", "", "", fmt.Errorf("no author name found in response")
	}
	authorNameText := authorName[1]

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
	return authorNameText, captionText, responseBody, nil
}

func extractCount(responseBody, dataE2E string) string {
	regexPattern := fmt.Sprintf(`<strong data-e2e="%s" class="(?:.*?)">([^<]+)<\/strong>`, dataE2E)
	regex := regexp.MustCompile(regexPattern)
	matches := regex.FindStringSubmatch(responseBody)
	if len(matches) > 1 {
		return matches[1]
	}
	return "0"
}

type Counts struct {
	Likes     string
	Comments  string
	Favorited string
	Shares    string
}

func GetVideoDetails(responseBody string) Counts {
	likesCount := extractCount(responseBody, "like-count")
	commentsCount := extractCount(responseBody, "comment-count")
	favoritedCount := extractCount(responseBody, "undefined-count") // undefined-count lmao
	sharesCount := extractCount(responseBody, "share-count")

	return Counts{
		Likes:     likesCount,
		Comments:  commentsCount,
		Favorited: favoritedCount,
		Shares:    sharesCount,
	}
}
